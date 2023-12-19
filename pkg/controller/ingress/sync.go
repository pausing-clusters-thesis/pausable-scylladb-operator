package ingress

import (
	"context"
	"fmt"
	"time"

	"github.com/pausing-clusters-thesis/pausable-scylladb-operator/pkg/naming"
	podapiutil "github.com/pausing-clusters-thesis/pausable-scylladb-operator/pkg/thirdparty/k8s.io/kubernetes/pkg/api/v1/pod"
	podutil "github.com/pausing-clusters-thesis/pausable-scylladb-operator/pkg/thirdparty/k8s.io/kubernetes/pkg/util/pod"
	"github.com/scylladb/scylla-operator/pkg/helpers/slices"
	"github.com/scylladb/scylla-operator/pkg/internalapi"
	sonaming "github.com/scylladb/scylla-operator/pkg/naming"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/labels"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/klog/v2"
)

type serviceReference struct {
	name, namespace string
}

func (ic *Controller) sync(ctx context.Context) error {
	startTime := time.Now()
	klog.V(4).InfoS("Started sync", "startTime", startTime)
	defer func() {
		klog.V(4).InfoS("Finished sync", "duration", time.Since(startTime))
	}()

	ingressClasses, err := ic.ingressClassLister.List(labels.Everything())
	if err != nil {
		return fmt.Errorf("can't list ingress classes: %w", err)
	}

	getIngressClassName := func(ingressClass *networkingv1.IngressClass) string {
		return ingressClass.GetName()
	}
	ingressClassNamesProvisionedByUs := sets.New(slices.ConvertSlice(slices.Filter(ingressClasses, isIngressClassProvisionedByUs), getIngressClassName)...)

	ingresses, err := ic.ingressLister.List(labels.Everything())
	if err != nil {
		return fmt.Errorf("can't list ingresses: %w", err)
	}

	var errs []error
	var services []serviceReference
	routes := map[string][]string{}
	for _, ingress := range ingresses {
		if ingress.DeletionTimestamp != nil {
			continue
		}

		if ingress.Spec.IngressClassName == nil || !ingressClassNamesProvisionedByUs.Has(*ingress.Spec.IngressClassName) {
			continue
		}

		var defaultBackendServicePort int32
		var defaultBackendServiceName string
		hasSupportedDefaultBackend := ingress.Spec.DefaultBackend != nil && ingress.Spec.DefaultBackend.Service != nil
		if hasSupportedDefaultBackend {
			defaultBackendServiceName = ingress.Spec.DefaultBackend.Service.Name
			defaultBackendServicePort = ingress.Spec.DefaultBackend.Service.Port.Number
			if defaultBackendServicePort == 0 {
				var defaultBackendService *corev1.Service
				defaultBackendService, err = ic.serviceLister.Services(ingress.Namespace).Get(ingress.Spec.DefaultBackend.Service.Name)
				if err != nil {
					return fmt.Errorf("can't get default backend service %q: %w", naming.ManualRef(ingress.Namespace, ingress.Spec.DefaultBackend.Service.Name), err)
				}

				defaultBackendServicePort, err = getServicePortByName(defaultBackendService, ingress.Spec.DefaultBackend.Service.Port.Name)
				if err != nil {
					return fmt.Errorf("can't get default backend service port by name: %w", err)
				}
			}
		}

		for _, rule := range ingress.Spec.Rules {
			host := rule.Host
			if len(host) == 0 {
				// TODO: event instead?
				klog.V(4).InfoS("Ingress rule is missing a host, skipping.", "Ingress", klog.KObj(ingress), "Rule", rule)
				continue
			}

			if rule.HTTP == nil {
				if !hasSupportedDefaultBackend {
					// TODO: event instead?
					klog.V(4).InfoS("Ingress rule is missing httpIngressRuleValue and ingress does not have a supported definition of defaultBackend, skipping.", "Ingress", klog.KObj(ingress), "Rule", rule)
					continue
				}

				backendHost := fmt.Sprintf("%s.%s.svc.cluster.local:%d", defaultBackendServiceName, ingress.Namespace, defaultBackendServicePort)
				routes[host] = append(routes[host], backendHost)
				continue
			}

			for _, path := range rule.HTTP.Paths {
				// TODO: verify path type?

				if path.Path != "/" && path.Path != "/*" {
					// TODO: event instead?
					klog.V(4).InfoS("Ingress rule's path specifies unsupported path, skipping.", "Ingress", klog.KObj(ingress), "Rule", rule, "Path", path)
					continue
				}

				if path.Backend.Service == nil {
					// TODO: event instead?
					klog.V(4).InfoS("Ingress rule's path does not reference service as backend, skipping.", "Ingress", klog.KObj(ingress), "Rule", rule, "Path", path)
					continue
				}

				svcName := path.Backend.Service.Name
				svcPort := path.Backend.Service.Port.Number
				if svcPort == 0 {
					var svc *corev1.Service
					svc, err = ic.serviceLister.Services(ingress.Namespace).Get(svcName)
					if err != nil {
						errs = append(errs, fmt.Errorf("can't get backend service %q: %w", naming.ManualRef(ingress.Namespace, svcName), err))
						continue
					}

					svcPort, err = getServicePortByName(svc, path.Backend.Service.Port.Name)
					if err != nil {
						errs = append(errs, fmt.Errorf("can't get rule path service port by name: %w", err))
						continue
					}
				}

				backendHost := fmt.Sprintf("%s.%s.svc.cluster.local:%d", svcName, ingress.Namespace, svcPort)
				routes[host] = append(routes[host], backendHost)
				services = append(services, serviceReference{
					name:      svcName,
					namespace: ingress.Namespace,
				})
			}
		}
	}

	sanitizedRoutes := map[string]string{}
	for host, backendHosts := range routes {
		if len(backendHosts) > 1 {
			// TODO: create an event with information about the duplicated ingresses
			errs = append(errs, fmt.Errorf("duplicate routes for host %q: %v, skipping", host, backendHosts))
			delete(routes, host)
			continue
		}

		sanitizedRoutes[host] = backendHosts[0]
	}

	klog.V(5).InfoS("Saving routes to the routing table.", "Routes", sanitizedRoutes)
	ic.routingTable.Store(sanitizedRoutes)

	// TODO: the condition should be updated on recreation of the routing table
	for _, svcRef := range services {
		svc, err := ic.serviceLister.Services(svcRef.namespace).Get(svcRef.name)
		if err != nil {
			errs = append(errs, fmt.Errorf("can't get service %q: %w", sonaming.ManualRef(svcRef.namespace, svcRef.name), err))
			continue
		}

		scyllaServiceTypeLabel := svc.GetLabels()[sonaming.ScyllaServiceTypeLabel]
		if len(scyllaServiceTypeLabel) == 0 || sonaming.ScyllaServiceType(scyllaServiceTypeLabel) != sonaming.ScyllaServiceTypeMember {
			continue
		}

		pods, err := ic.podLister.Pods(svcRef.namespace).List(labels.SelectorFromSet(svc.Spec.Selector))
		if err != nil {
			errs = append(errs, fmt.Errorf("can't list pods for service %q: %w", sonaming.ObjRef(svc), err))
			continue
		}

		for _, pod := range pods {
			podStatusCopy := pod.Status.DeepCopy()
			podConditionUpdated := podapiutil.UpdatePodCondition(podStatusCopy, &corev1.PodCondition{
				Type:    naming.IngressControllerScyllaDBMemberPodConditionType,
				Status:  corev1.ConditionTrue,
				Reason:  internalapi.AsExpectedReason,
				Message: "",
			})

			if !podConditionUpdated {
				continue
			}

			_, _, _, err = podutil.PatchPodStatus(ctx, ic.kubeClient, pod.GetNamespace(), pod.GetName(), pod.GetUID(), pod.Status, *podStatusCopy)
			if err != nil {
				errs = append(errs, fmt.Errorf("can't patch pod status for pod %q: %w", sonaming.ObjRef(pod), err))
			}
		}
	}

	return utilerrors.NewAggregate(errs)
}

func getServicePortByName(svc *corev1.Service, portName string) (int32, error) {
	p, _, ok := slices.Find(svc.Spec.Ports, func(p corev1.ServicePort) bool {
		return p.Name == portName
	})
	if !ok {
		return 0, fmt.Errorf("service %q does not have a port %q", naming.ObjRef(svc), portName)
	}

	return p.Port, nil
}
