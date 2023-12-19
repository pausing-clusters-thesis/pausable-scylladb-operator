package scylladbdatacenterpool

import (
	"context"
	"fmt"
	"maps"
	"slices"
	"time"

	pausingv1alpha1 "github.com/pausing-clusters-thesis/pausable-scylladb-operator/pkg/api/pausing/v1alpha1"
	"github.com/pausing-clusters-thesis/pausable-scylladb-operator/pkg/controllerhelpers"
	"github.com/pausing-clusters-thesis/pausable-scylladb-operator/pkg/naming"
	scyllav1alpha1 "github.com/scylladb/scylla-operator/pkg/api/scylla/v1alpha1"
	socontrollerhelpers "github.com/scylladb/scylla-operator/pkg/controllerhelpers"
	soslices "github.com/scylladb/scylla-operator/pkg/helpers/slices"
	sonaming "github.com/scylladb/scylla-operator/pkg/naming"
	"github.com/scylladb/scylla-operator/pkg/pointer"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
)

func (sdcpc *Controller) syncScyllaDBDatacenters(
	ctx context.Context,
	key string,
	sdcp *pausingv1alpha1.ScyllaDBDatacenterPool,
	scyllaDBDatacenterClaims []*pausingv1alpha1.ScyllaDBDatacenterClaim,
	scyllaDBDatacenters map[string]*scyllav1alpha1.ScyllaDBDatacenter,
) ([]metav1.Condition, error) {
	var progressingConditions []metav1.Condition
	var err error

	startTime := time.Now()
	klog.V(4).InfoS("Started syncing ScyllaDBDatacenters", "ScyllaDBDatacenterPool", klog.KObj(sdcp), "startTime", startTime)
	defer func() {
		klog.V(4).InfoS("Finished syncing ScyllaDBDatacenters", "ScyllaDBDatacenterPool", klog.KObj(sdcp), "duration", time.Since(startTime))
	}()

	unclaimedScyllaDBDatacentersMap := maps.Clone(scyllaDBDatacenters)
	for _, sdcc := range scyllaDBDatacenterClaims {
		if sdcc.Status.ScyllaDBDatacenterName == nil || len(*sdcc.Status.ScyllaDBDatacenterName) == 0 {
			continue
		}

		delete(unclaimedScyllaDBDatacentersMap, sdcc.Name)
	}

	sortScyllaDBDatacentersByStatusFunc := func(first, second *scyllav1alpha1.ScyllaDBDatacenter) int {
		firstRolledOut, err := socontrollerhelpers.IsScyllaDBDatacenterRolledOut(first)
		if err == nil && firstRolledOut {
			return -1
		}

		secondRolledOut, err := socontrollerhelpers.IsScyllaDBDatacenterRolledOut(first)
		if err == nil && secondRolledOut {
			return 1
		}

		// Neither is rolled out nor do we know their statuses, we try comparing creation timestamps as best effort.
		return first.CreationTimestamp.Compare(second.CreationTimestamp.Time)
	}

	unclaimedScyllaDBDatacenters := slices.Collect(maps.Values(scyllaDBDatacenters))
	slices.SortFunc(unclaimedScyllaDBDatacenters, sortScyllaDBDatacentersByStatusFunc)

	for _, sdcc := range scyllaDBDatacenterClaims {
		if sdcc.DeletionTimestamp != nil {
			continue
		}

		if controllerhelpers.HasAnnotation(sdcc, naming.ScyllaDBDatacenterClaimBindCompletedAnnotation) {
			// TODO: sync bound
		} else {
			klog.V(5).InfoS("Synchronising unbound ScyllaDBDatacenterClaim", "ScyllaDBDatacenterPool", klog.KObj(sdcp), "ScyllaDBDatacenterClaim", klog.KObj(sdcc))

			if sdcc.Status.ScyllaDBDatacenterName == nil || len(*sdcc.Status.ScyllaDBDatacenterName) == 0 {
				// No ScyllaDBDatacenter has been claimed yet.
				klog.V(5).InfoS("Synchronising unbound ScyllaDBDatacenterClaim: claiming a new ScyllaDBDatacenter", "ScyllaDBDatacenterPool", klog.KObj(sdcp), "ScyllaDBDatacenterClaim", klog.KObj(sdcc))

				if len(unclaimedScyllaDBDatacenters) == 0 {
					// Create on demand.
					klog.V(5).InfoS("Synchronising unbound ScyllaDBDatacenterClaim: no unclaimed ScyllaDBDatacenters in pool, creating on demand", "ScyllaDBDatacenterPool", klog.KObj(sdcp), "ScyllaDBDatacenterClaim", klog.KObj(sdcc))

					sdc := makeScyllaDBDatacenter(sdcp)
					sdc = claimScyllaDBDatacenter(sdc, sdcc)

					// We use a naming scheme in this case to know the name of the ScyllaDBDatacenter to be created beforehand.
					sdc.Name = sdcc.Name

					sdccCopy := sdcc.DeepCopy()
					sdccCopy.Status.ScyllaDBDatacenterName = pointer.Ptr(sdc.Name)

					sdccCopy, err = sdcpc.pausingClient.ScyllaDBDatacenterClaims(sdcc.Namespace).UpdateStatus(ctx, sdccCopy, metav1.UpdateOptions{})
					if err != nil {
						return progressingConditions, fmt.Errorf("can't update ScyllaDBDatacenterClaim %q status: %w", sonaming.ObjRef(sdccCopy), err)
					}

					_, err = sdcpc.scyllaClient.ScyllaDBDatacenters(sdcc.Namespace).Create(ctx, sdc, metav1.CreateOptions{})
					if err != nil {
						return progressingConditions, fmt.Errorf("can't create ScyllaDBDatacenter %q: %w", naming.ManualRef(sdcc.Namespace, sdc.Name), err)
					}

					sdccCopy = sdccCopy.DeepCopy()
					metav1.SetMetaDataAnnotation(&sdccCopy.ObjectMeta, naming.ScyllaDBDatacenterClaimBindCompletedAnnotation, "")
					_, err = sdcpc.pausingClient.ScyllaDBDatacenterClaims(sdcc.Namespace).Update(ctx, sdccCopy, metav1.UpdateOptions{})
					if err != nil {
						return progressingConditions, fmt.Errorf("can't update ScyllaDBDatacenterClaim %q status: %w", sonaming.ObjRef(sdccCopy), err)
					}
				} else {
					sdc := unclaimedScyllaDBDatacenters[0]
					unclaimedScyllaDBDatacenters = unclaimedScyllaDBDatacenters[1:]

					klog.V(5).InfoS("Synchronising unbound ScyllaDBDatacenterClaim: retrieving unclaimed ScyllaDBDatacenter from pool", "ScyllaDBDatacenterPool", klog.KObj(sdcp), "ScyllaDBDatacenterClaim", klog.KObj(sdcc), "ScyllaDBDatacenter", klog.KObj(sdc))

					sdcCopy := claimScyllaDBDatacenter(sdc, sdcc)

					sdccCopy := sdcc.DeepCopy()
					sdccCopy.Status.ScyllaDBDatacenterName = pointer.Ptr(sdcCopy.Name)

					sdccCopy, err = sdcpc.pausingClient.ScyllaDBDatacenterClaims(sdcc.Namespace).UpdateStatus(ctx, sdccCopy, metav1.UpdateOptions{})
					if err != nil {
						return progressingConditions, fmt.Errorf("can't update ScyllaDBDatacenterClaim %q status: %w", sonaming.ObjRef(sdccCopy), err)
					}

					_, err = sdcpc.scyllaClient.ScyllaDBDatacenters(sdcCopy.Namespace).Update(ctx, sdcCopy, metav1.UpdateOptions{})
					if err != nil {
						return progressingConditions, fmt.Errorf("can't update ScyllaDBDatacenter %q: %w", naming.ObjRef(sdc), err)
					}

					sdccCopy = sdccCopy.DeepCopy()
					metav1.SetMetaDataAnnotation(&sdccCopy.ObjectMeta, naming.ScyllaDBDatacenterClaimBindCompletedAnnotation, "")
					_, err = sdcpc.pausingClient.ScyllaDBDatacenterClaims(sdcc.Namespace).Update(ctx, sdccCopy, metav1.UpdateOptions{})
					if err != nil {
						return progressingConditions, fmt.Errorf("can't update ScyllaDBDatacenterClaim %q status: %w", sonaming.ObjRef(sdccCopy), err)
					}
				}
			} else {
				klog.V(5).InfoS("Synchronising unbound ScyllaDBDatacenterClaim: ScyllaDBDatacenter already claimed in previous iteration", "ScyllaDBDatacenterPool", klog.KObj(sdcp), "ScyllaDBDatacenterClaim", klog.KObj(sdcc), "ScyllaDBDatacenter", klog.KRef(sdcc.Namespace, *sdcc.Status.ScyllaDBDatacenterName))

				// SDC was already claimed in previous iterations, continue binding.
				sdc, ok := scyllaDBDatacenters[*sdcc.Status.ScyllaDBDatacenterName]
				if ok {
					sdcCopy := claimScyllaDBDatacenter(sdc, sdcc)

					sdccCopy := sdcc.DeepCopy()
					sdccCopy.Status.ScyllaDBDatacenterName = pointer.Ptr(sdcCopy.Name)

					sdccCopy, err = sdcpc.pausingClient.ScyllaDBDatacenterClaims(sdcc.Namespace).UpdateStatus(ctx, sdccCopy, metav1.UpdateOptions{})
					if err != nil {
						return progressingConditions, fmt.Errorf("can't update ScyllaDBDatacenterClaim %q status: %w", sonaming.ObjRef(sdccCopy), err)
					}

					_, err = sdcpc.scyllaClient.ScyllaDBDatacenters(sdcCopy.Namespace).Update(ctx, sdcCopy, metav1.UpdateOptions{})
					if err != nil {
						return progressingConditions, fmt.Errorf("can't update ScyllaDBDatacenter %q: %w", naming.ObjRef(sdc), err)
					}

					sdccCopy = sdccCopy.DeepCopy()
					metav1.SetMetaDataAnnotation(&sdccCopy.ObjectMeta, naming.ScyllaDBDatacenterClaimBindCompletedAnnotation, "")
					_, err = sdcpc.pausingClient.ScyllaDBDatacenterClaims(sdcc.Namespace).Update(ctx, sdccCopy, metav1.UpdateOptions{})
					if err != nil {
						return progressingConditions, fmt.Errorf("can't update ScyllaDBDatacenterClaim %q status: %w", sonaming.ObjRef(sdccCopy), err)
					}
				} else {
					// Use client to make sure we haven't hit a stale cache.
					sdc, err = sdcpc.scyllaClient.ScyllaDBDatacenters(sdcc.Namespace).Get(ctx, *sdcc.Status.ScyllaDBDatacenterName, metav1.GetOptions{})
					if err != nil {
						if !errors.IsNotFound(err) {
							return progressingConditions, fmt.Errorf("can't get ScyllaDBDatacenter %q: %w", naming.ManualRef(sdcc.Namespace, *sdcc.Status.ScyllaDBDatacenterName), err)
						}

						// ScyllaDBDatacenterClaim claimed a ScyllaDBDatacenter which doesn't exist.
						// Clear the reference on claim.
						sdccCopy := sdcc.DeepCopy()
						sdccCopy.Status.ScyllaDBDatacenterName = nil

						_, err = sdcpc.pausingClient.ScyllaDBDatacenterClaims(sdcc.Namespace).UpdateStatus(ctx, sdccCopy, metav1.UpdateOptions{})
						if err != nil {
							return progressingConditions, fmt.Errorf("can't update ScyllaDBDatacenterClaim %q status: %w", sonaming.ObjRef(sdccCopy), err)
						}

						// Claim will be bound in the next interation.
						continue
					}

					sdcControllerRef := metav1.GetControllerOfNoCopy(sdc)
					if sdcControllerRef == nil || sdcControllerRef.UID != sdcc.UID {
						// ScyllaDBDatacenterClaim claimed a ScyllaDBDatacenter which is controlled by someone else.
						// Clear the reference on claim.
						sdccCopy := sdcc.DeepCopy()
						sdccCopy.Status.ScyllaDBDatacenterName = nil

						_, err = sdcpc.pausingClient.ScyllaDBDatacenterClaims(sdcc.Namespace).UpdateStatus(ctx, sdccCopy, metav1.UpdateOptions{})
						if err != nil {
							return progressingConditions, fmt.Errorf("can't update ScyllaDBDatacenterClaim %q status: %w", sonaming.ObjRef(sdccCopy), err)
						}

						// Claim will be bound in the next interation.
						continue
					} else {
						// SDC is owned by claim. Complete binding.
						sdccCopy := sdcc.DeepCopy()
						metav1.SetMetaDataAnnotation(&sdccCopy.ObjectMeta, naming.ScyllaDBDatacenterClaimBindCompletedAnnotation, "")
						_, err = sdcpc.pausingClient.ScyllaDBDatacenterClaims(sdcc.Namespace).Update(ctx, sdccCopy, metav1.UpdateOptions{})
						if err != nil {
							return progressingConditions, fmt.Errorf("can't update ScyllaDBDatacenterClaim %q status: %w", sonaming.ObjRef(sdccCopy), err)
						}
					}
				}
			}
		}
	}

	// Set a progressing condition and queue immediately to let the released ScyllaDBDatacenters to propagate into cache and minimize fluctuations of the pool size.
	if len(unclaimedScyllaDBDatacenters) < len(scyllaDBDatacenters) {
		progressingConditions = append(progressingConditions, metav1.Condition{
			Type:               scyllaDBDatacenterControllerProgressingCondition,
			Status:             metav1.ConditionTrue,
			ObservedGeneration: sdcp.Generation,
			Reason:             "AwaitingReleasedScyllaDBDatacentersCacheSync",
			Message:            "Waiting for the released ScyllaDBDatacenters to propagate into cache.",
		})

		sdcpc.queue.Add(key)
		return progressingConditions, nil
	}

	capacityDiff := int(sdcp.Spec.Capacity) - len(unclaimedScyllaDBDatacenters)
	limitDiff := int(sdcp.Spec.Limit) - len(unclaimedScyllaDBDatacenters)

	if capacityDiff > 0 {
		klog.V(4).InfoS("Too few ScyllaDBDatacenters in the pool", "ScyllaDBDatacenterPool", klog.KObj(sdcp), "capacity", sdcp.Spec.Capacity, "got", len(unclaimedScyllaDBDatacenters), "creating", capacityDiff)

		scyllaDBDatacenterTemplate := makeScyllaDBDatacenter(sdcp)

		// TODO: create concurrently
		// TODO: create in bursts
		for range capacityDiff {
			sdc, err := sdcpc.scyllaClient.ScyllaDBDatacenters(sdcp.Namespace).Create(ctx, scyllaDBDatacenterTemplate, metav1.CreateOptions{})
			if err != nil {
				return progressingConditions, fmt.Errorf("can't create ScyllaDBDatacenter in ScyllaDBDatacenterPool %q: %w", sonaming.ObjRef(sdcp), err)
			}
			socontrollerhelpers.AddGenericProgressingStatusCondition(&progressingConditions, scyllaDBDatacenterControllerProgressingCondition, sdc, "create", sdcp.Generation)
		}
	} else if limitDiff < 0 {
		limitDiff *= -1
		klog.V(4).InfoS("Too many ScyllaDBDatacenters in the pool", "ScyllaDBDatacenterPool", klog.KObj(sdcp), "limit", sdcp.Spec.Limit, "got", len(unclaimedScyllaDBDatacenters), "pruning", limitDiff)

		// TODO: delete concurrently
		for _, sdc := range unclaimedScyllaDBDatacenters[:limitDiff] {
			err = sdcpc.scyllaClient.ScyllaDBDatacenters(sdcp.Namespace).Delete(ctx, sdc.Name, metav1.DeleteOptions{})
			if err != nil && !errors.IsNotFound(err) {
				return progressingConditions, fmt.Errorf("can't delete ScyllaDBDatacenter %q in ScyllaDBDatacenterPool %q: %w", sonaming.ObjRef(sdc), sonaming.ObjRef(sdcp), err)
			}
		}
	}

	return progressingConditions, nil
}

func makeScyllaDBDatacenter(sdcp *pausingv1alpha1.ScyllaDBDatacenterPool) *scyllav1alpha1.ScyllaDBDatacenter {
	return &scyllav1alpha1.ScyllaDBDatacenter{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: fmt.Sprintf("%s-", sdcp.Name),
			Namespace:    sdcp.Namespace,
			Labels: func() map[string]string {
				ls := map[string]string{}

				maps.Copy(ls, naming.ScyllaDBDatacenterPoolSelectorLabels(sdcp))

				return ls
			}(),
			Annotations: map[string]string{
				sonaming.DelayedVolumeMountAnnotation: "",
			},
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(sdcp, pausingv1alpha1.ScyllaDBDatacenterPoolControllerGVK),
			},
		},
		Spec: func() scyllav1alpha1.ScyllaDBDatacenterSpec {
			specCopy := *sdcp.Spec.Template.Spec.DeepCopy()

			for i := range specCopy.Racks {
				if specCopy.Racks[i].ScyllaDB == nil {
					specCopy.Racks[i].ScyllaDB = &scyllav1alpha1.ScyllaDBTemplate{}
				}

				if specCopy.Racks[i].ScyllaDB.Storage == nil {
					specCopy.Racks[i].ScyllaDB.Storage = &scyllav1alpha1.StorageOptions{}
				}

				specCopy.Racks[i].ScyllaDB.Storage.StorageClassName = pointer.Ptr(sdcp.Spec.ProxyStorageClassName)
			}

			return specCopy
		}(),
	}
}

func claimScyllaDBDatacenter(sdc *scyllav1alpha1.ScyllaDBDatacenter, sdcc *pausingv1alpha1.ScyllaDBDatacenterClaim) *scyllav1alpha1.ScyllaDBDatacenter {
	sdcCopy := sdc.DeepCopy()

	for _, k := range naming.ScyllaDBDatacenterPoolSelectorLabelKeys() {
		delete(sdcCopy.Labels, k)
	}

	maps.Copy(sdcCopy.Labels, naming.ScyllaDBDatacenterClaimSelectorLabels(sdcc))

	sdcCopy.OwnerReferences = soslices.FilterOut(sdcCopy.OwnerReferences, func(ref metav1.OwnerReference) bool {
		return ref.Kind == pausingv1alpha1.ScyllaDBDatacenterPoolControllerGVK.Kind
	})

	sdcCopy.OwnerReferences = append(sdcCopy.OwnerReferences, *metav1.NewControllerRef(sdcc, pausingv1alpha1.ScyllaDBDatacenterClaimControllerGVK))

	return sdcCopy
}
