package run

import (
	"context"
	"fmt"
	"sync"
	"time"

	pausingversionedclient "github.com/pausing-clusters-thesis/pausable-scylladb-operator/pkg/client/pausing/clientset/versioned"
	pausinginformers "github.com/pausing-clusters-thesis/pausable-scylladb-operator/pkg/client/pausing/informers/externalversions"
	"github.com/pausing-clusters-thesis/pausable-scylladb-operator/pkg/controller/pausablescylladbdatacenter"
	"github.com/pausing-clusters-thesis/pausable-scylladb-operator/pkg/controller/scylladbdatacenterclaim"
	"github.com/pausing-clusters-thesis/pausable-scylladb-operator/pkg/controller/scylladbdatacenterpool"
	"github.com/pausing-clusters-thesis/pausable-scylladb-operator/pkg/version"
	scyllaversionedclient "github.com/scylladb/scylla-operator/pkg/client/scylla/clientset/versioned"
	scyllainformers "github.com/scylladb/scylla-operator/pkg/client/scylla/informers/externalversions"
	"github.com/scylladb/scylla-operator/pkg/crypto"
	sogenericclioptions "github.com/scylladb/scylla-operator/pkg/genericclioptions"
	soleaderelection "github.com/scylladb/scylla-operator/pkg/leaderelection"
	sosignals "github.com/scylladb/scylla-operator/pkg/signals"
	"github.com/spf13/cobra"
	apierrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/cli-runtime/pkg/genericiooptions"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	cliflag "k8s.io/component-base/cli/flag"
	"k8s.io/klog/v2"
)

const (
	resyncPeriod                  = 12 * time.Hour
	cryptoKeyBufferSizeMaxFlagKey = "crypto-key-buffer-size-max"
)

type Options struct {
	sogenericclioptions.ClientConfig
	sogenericclioptions.InClusterReflection
	sogenericclioptions.LeaderElection

	kubeClient    kubernetes.Interface
	pausingClient pausingversionedclient.Interface
	scyllaClient  scyllaversionedclient.Interface

	ConcurrentSyncs int

	CryptoKeyBufferSizeMin int
	CryptoKeyBufferSizeMax int
	CryptoKeyBufferDelay   time.Duration
}

func NewOptions(streams genericiooptions.IOStreams) *Options {
	return &Options{
		ClientConfig:        sogenericclioptions.NewClientConfig("pausable-scylladb-operator"),
		InClusterReflection: sogenericclioptions.InClusterReflection{},
		LeaderElection:      sogenericclioptions.NewLeaderElection(),

		ConcurrentSyncs: 50,

		CryptoKeyBufferSizeMin: 2,
		CryptoKeyBufferSizeMax: 6,
		CryptoKeyBufferDelay:   200 * time.Millisecond,
	}
}

func NewCommand(streams genericiooptions.IOStreams) *cobra.Command {
	o := NewOptions(streams)

	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run pausable-scylladb-operator.",
		Long:  `Run pausable-scylladb-operator.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ver := version.Get()
			klog.InfoS("Binary info",
				"Command", cmd.Name(),
				"Revision", version.OptionalToString(ver.Revision),
				"RevisionTime", version.OptionalToString(ver.RevisionTime),
				"Modified", version.OptionalToString(ver.Modified),
				"GoVersion", version.OptionalToString(ver.GoVersion),
			)
			cliflag.PrintFlags(cmd.Flags())

			stopCh := sosignals.StopChannel()
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			go func() {
				<-stopCh
				cancel()
			}()

			err := o.Validate()
			if err != nil {
				return err
			}

			err = o.Complete(ctx)
			if err != nil {
				return err
			}

			err = o.Run(ctx, streams, cmd)
			if err != nil {
				return err
			}

			return nil
		},

		SilenceErrors: true,
		SilenceUsage:  true,
	}

	o.ClientConfig.AddFlags(cmd)
	o.InClusterReflection.AddFlags(cmd)
	o.LeaderElection.AddFlags(cmd)

	cmd.Flags().IntVarP(&o.ConcurrentSyncs, "concurrent-syncs", "", o.ConcurrentSyncs, "The number of objects that are allowed to sync concurrently for each controller.")
	cmd.Flags().IntVarP(&o.CryptoKeyBufferSizeMin, "crypto-key-buffer-size-min", "", o.CryptoKeyBufferSizeMin, "Minimal number of pre-generated crypto keys that are used for quick certificate issuance. The minimum size is 1.")
	cmd.Flags().IntVarP(&o.CryptoKeyBufferSizeMax, cryptoKeyBufferSizeMaxFlagKey, "", o.CryptoKeyBufferSizeMax, "Maximum number of pre-generated crypto keys that are used for quick certificate issuance. The minimum size is 1. If not set, it will adjust to be at least the size of crypto-key-buffer-size-min.")
	cmd.Flags().DurationVarP(&o.CryptoKeyBufferDelay, "crypto-key-buffer-delay", "", o.CryptoKeyBufferDelay, "Delay is the time to wait when generating next certificate in the (min, max) range. Certificate generation bellow the min threshold is not affected.")

	return cmd
}

func (o *Options) Validate() error {
	var errs []error

	errs = append(errs, o.ClientConfig.Validate())
	errs = append(errs, o.InClusterReflection.Validate())
	errs = append(errs, o.LeaderElection.Validate())

	if o.CryptoKeyBufferSizeMin < 1 {
		errs = append(errs, fmt.Errorf("crypto-key-buffer-size-min (%d) has to be at least 1", o.CryptoKeyBufferSizeMin))
	}

	if o.CryptoKeyBufferSizeMax < 1 {
		errs = append(errs, fmt.Errorf("crypto-key-buffer-size-max (%d) has to be at least 1", o.CryptoKeyBufferSizeMax))
	}

	if o.CryptoKeyBufferSizeMax < o.CryptoKeyBufferSizeMin {
		errs = append(errs, fmt.Errorf(
			"crypto-key-buffer-size-max (%d) can't be lower then crypto-key-buffer-size-min (%d)",
			o.CryptoKeyBufferSizeMax,
			o.CryptoKeyBufferSizeMin,
		))
	}

	return apierrors.NewAggregate(errs)
}

func (o *Options) Complete(ctx context.Context) error {
	err := o.ClientConfig.Complete()
	if err != nil {
		return err
	}

	err = o.InClusterReflection.Complete()
	if err != nil {
		return err
	}

	err = o.LeaderElection.Complete()
	if err != nil {
		return err
	}

	o.kubeClient, err = kubernetes.NewForConfig(o.ProtoConfig)
	if err != nil {
		return fmt.Errorf("can't build kubernetes clientset: %w", err)
	}

	o.pausingClient, err = pausingversionedclient.NewForConfig(o.RestConfig)
	if err != nil {
		return fmt.Errorf("can't build pausing clientset: %w", err)
	}

	o.scyllaClient, err = scyllaversionedclient.NewForConfig(o.RestConfig)
	if err != nil {
		return fmt.Errorf("can't build scylla clientset: %w", err)
	}

	return nil
}

func (o *Options) Run(ctx context.Context, streams genericiooptions.IOStreams, cmd *cobra.Command) error {
	// Lock names cannot be changed, because it may lead to two leaders during rolling upgrades.
	const lockName = "pausable-scylladb-operator-lock"

	return soleaderelection.Run(
		ctx,
		cmd.Name(),
		lockName,
		o.Namespace,
		o.kubeClient,
		o.LeaderElectionLeaseDuration,
		o.LeaderElectionRenewDeadline,
		o.LeaderElectionRetryPeriod,
		func(ctx context.Context) error {
			return o.run(ctx, streams)
		},
	)
}

func (o *Options) run(ctx context.Context, streams genericiooptions.IOStreams) error {
	rsaKeyGenerator, err := crypto.NewRSAKeyGenerator(
		o.CryptoKeyBufferSizeMin,
		o.CryptoKeyBufferSizeMax,
		o.CryptoKeyBufferDelay,
	)
	if err != nil {
		return fmt.Errorf("can't create rsa key generator: %w", err)
	}
	defer rsaKeyGenerator.Close()

	kubeInformers := informers.NewSharedInformerFactory(o.kubeClient, resyncPeriod)
	pausingInformers := pausinginformers.NewSharedInformerFactory(o.pausingClient, resyncPeriod)
	scyllaInformers := scyllainformers.NewSharedInformerFactory(o.scyllaClient, resyncPeriod)

	psdcc, err := pausablescylladbdatacenter.NewController(
		o.kubeClient,
		o.pausingClient.PausingV1alpha1(),
		pausingInformers.Pausing().V1alpha1().PausableScyllaDBDatacenters(),
		pausingInformers.Pausing().V1alpha1().ScyllaDBDatacenterClaims(),
		pausingInformers.Pausing().V1alpha1().ScyllaDBDatacenterPools(),
		kubeInformers.Core().V1().PersistentVolumeClaims(),
		scyllaInformers.Scylla().V1alpha1().ScyllaDBDatacenters(),
		kubeInformers.Core().V1().Secrets(),
		kubeInformers.Core().V1().ConfigMaps(),
		rsaKeyGenerator,
	)
	if err != nil {
		return fmt.Errorf("can't create PausableScyllaDBDatacenter controller: %w", err)
	}

	sdcpc, err := scylladbdatacenterpool.NewController(
		o.kubeClient,
		o.pausingClient.PausingV1alpha1(),
		o.scyllaClient.ScyllaV1alpha1(),
		pausingInformers.Pausing().V1alpha1().ScyllaDBDatacenterPools(),
		pausingInformers.Pausing().V1alpha1().ScyllaDBDatacenterClaims(),
		scyllaInformers.Scylla().V1alpha1().ScyllaDBDatacenters(),
	)
	if err != nil {
		return fmt.Errorf("can't create ScyllaDBDatacenterPool controller: %w", err)
	}

	sdccc, err := scylladbdatacenterclaim.NewController(
		o.kubeClient,
		o.pausingClient.PausingV1alpha1(),
		o.scyllaClient.ScyllaV1alpha1(),
		pausingInformers.Pausing().V1alpha1().ScyllaDBDatacenterClaims(),
		scyllaInformers.Scylla().V1alpha1().ScyllaDBDatacenters(),
	)
	if err != nil {
		return fmt.Errorf("can't create ScyllaDBDatacenterClaim controller: %w", err)
	}

	var wg sync.WaitGroup
	defer wg.Wait()

	/* Start informers */

	wg.Add(1)
	go func() {
		defer wg.Done()
		rsaKeyGenerator.Run(ctx)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		kubeInformers.Start(ctx.Done())
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		pausingInformers.Start(ctx.Done())
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		scyllaInformers.Start(ctx.Done())
	}()

	/* Start controllers */
	wg.Add(1)
	go func() {
		defer wg.Done()
		psdcc.Run(ctx, o.ConcurrentSyncs)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		sdcpc.Run(ctx, o.ConcurrentSyncs)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		sdccc.Run(ctx, o.ConcurrentSyncs)
	}()

	<-ctx.Done()

	return nil
}
