package pausablescylladbdatacenter

import (
	"context"
	"fmt"
	"time"

	pausingv1alpha1 "github.com/pausing-clusters-thesis/pausable-scylladb-operator/pkg/api/pausing/v1alpha1"
	"github.com/pausing-clusters-thesis/pausable-scylladb-operator/pkg/naming"
	sokubecrypto "github.com/scylladb/scylla-operator/pkg/kubecrypto"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apierrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/klog/v2"
)

func (psdcc *Controller) syncCerts(
	ctx context.Context,
	psc *pausingv1alpha1.PausableScyllaDBDatacenter,
	secrets map[string]*corev1.Secret,
) ([]metav1.Condition, []metav1.Condition, error) {
	var err error
	var errs []error
	var availableConditions, progressingConditions []metav1.Condition

	cm := sokubecrypto.NewCertificateManager(
		psdcc.keyGetter,
		psdcc.kubeClient.CoreV1(),
		psdcc.secretLister,
		psdcc.kubeClient.CoreV1(),
		psdcc.configMapLister,
		psdcc.eventRecorder,
	)

	// TODO: extend with standard labels
	labels := getSelectorLabels(psc)

	klog.V(4).InfoS("Building client CA", "PausableScyllaDBDatacenter", klog.KObj(psc))
	_, err = cm.ManageSigningTLSSecret(
		ctx,
		time.Now,
		&psc.ObjectMeta,
		pausableScyllaDBDatacenterControllerGVK,
		&sokubecrypto.CAConfig{
			MetaConfig: sokubecrypto.MetaConfig{
				Name:   naming.GetPausableScyllaDBDatacenterLocalClientCAName(psc.Name),
				Labels: labels,
			},
			Validity: 10 * 365 * 24 * time.Hour,
			Refresh:  8 * 365 * 24 * time.Hour,
		},
		secrets,
	)
	if err != nil {
		errs = append(errs, fmt.Errorf("can't manage signing TLS secret for client CA: %w", err))
	}

	klog.V(4).InfoS("Building serving CA", "PausableScyllaDBDatacenter", klog.KObj(psc))
	_, err = cm.ManageSigningTLSSecret(
		ctx,
		time.Now,
		&psc.ObjectMeta,
		pausableScyllaDBDatacenterControllerGVK,
		&sokubecrypto.CAConfig{
			MetaConfig: sokubecrypto.MetaConfig{
				Name:   naming.GetPausableScyllaDBDatacenterLocalServingCAName(psc.Name),
				Labels: labels,
			},
			Validity: 10 * 365 * 24 * time.Hour,
			Refresh:  8 * 365 * 24 * time.Hour,
		},
		secrets,
	)
	if err != nil {
		errs = append(errs, fmt.Errorf("can't manage signing TLS secret for client CA: %w", err))
	}

	return availableConditions, progressingConditions, apierrors.NewAggregate(errs)
}
