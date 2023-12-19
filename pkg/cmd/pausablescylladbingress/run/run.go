package run

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"

	pausingversionedclient "github.com/pausing-clusters-thesis/pausable-scylladb-operator/pkg/client/pausing/clientset/versioned"
	pausinginformers "github.com/pausing-clusters-thesis/pausable-scylladb-operator/pkg/client/pausing/informers/externalversions"
	"github.com/pausing-clusters-thesis/pausable-scylladb-operator/pkg/controller/ingress"
	"github.com/pausing-clusters-thesis/pausable-scylladb-operator/pkg/proxy"
	"github.com/pausing-clusters-thesis/pausable-scylladb-operator/pkg/routing"
	"github.com/pausing-clusters-thesis/pausable-scylladb-operator/pkg/version"
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
	netutils "k8s.io/utils/net"
)

const resyncPeriod = 12 * time.Hour

type Options struct {
	sogenericclioptions.ClientConfig
	sogenericclioptions.InClusterReflection
	sogenericclioptions.LeaderElection

	kubeClient    kubernetes.Interface
	pausingClient pausingversionedclient.Interface

	ConcurrentSyncs int

	IPv4Address string
	HTTPSPort   string
}

func NewOptions(streams genericiooptions.IOStreams) *Options {
	return &Options{
		ClientConfig:        sogenericclioptions.NewClientConfig("proxy"),
		InClusterReflection: sogenericclioptions.InClusterReflection{},
		LeaderElection:      sogenericclioptions.NewLeaderElection(),

		ConcurrentSyncs: 50,

		IPv4Address: "0.0.0.0",
		HTTPSPort:   "8443",
	}
}

func NewCommand(streams genericiooptions.IOStreams) *cobra.Command {
	o := NewOptions(streams)

	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run ingress.",
		Long:  `Run ingress.`,
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
	cmd.Flags().StringVarP(&o.IPv4Address, "ipv4-address", "", o.IPv4Address, "IPv4 address that the proxy server listens on.")
	cmd.Flags().StringVarP(&o.HTTPSPort, "https-port", "", o.HTTPSPort, "Secure port that the proxy server listens on.")

	return cmd
}

func (o *Options) Validate() error {
	var err error
	var errs []error

	errs = append(errs, o.ClientConfig.Validate())
	errs = append(errs, o.InClusterReflection.Validate())
	errs = append(errs, o.LeaderElection.Validate())

	ip := netutils.ParseIPSloppy(o.IPv4Address)
	if ip == nil || ip.To4() == nil {
		errs = append(errs, fmt.Errorf("ipv4-address must be a valid IPv4 address"))
	}

	_, err = netutils.ParsePort(o.HTTPSPort, false)
	if err != nil {
		errs = append(errs, fmt.Errorf("https-port must be a valid, non-zero port number: %w", err))
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

	return nil
}

func (o *Options) Run(ctx context.Context, streams genericiooptions.IOStreams, cmd *cobra.Command) error {
	// Lock names cannot be changed, because it may lead to two leaders during rolling upgrades.
	const lockName = "pausable-scylladb-ingress-lock"

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
	kubeInformers := informers.NewSharedInformerFactory(o.kubeClient, resyncPeriod)
	pausingInformers := pausinginformers.NewSharedInformerFactory(o.pausingClient, resyncPeriod)

	routingTable := routing.NewTable()

	ic, err := ingress.NewController(
		o.kubeClient,
		kubeInformers.Networking().V1().Ingresses(),
		kubeInformers.Networking().V1().IngressClasses(),
		kubeInformers.Core().V1().Services(),
		kubeInformers.Core().V1().Pods(),
		routingTable,
	)
	if err != nil {
		return fmt.Errorf("can't create ingress controller: %w", err)
	}

	lc := net.ListenConfig{}
	listener, err := lc.Listen(ctx, "tcp4", fmt.Sprintf("%s:%s", o.IPv4Address, o.HTTPSPort))
	if err != nil && !errors.Is(err, net.ErrClosed) {
		return fmt.Errorf("can't create listener: %w", err)
	}

	p := proxy.NewProxy(
		o.pausingClient.PausingV1alpha1(),
		pausingInformers.Pausing().V1alpha1().PausableScyllaDBDatacenters(),
		routingTable,
	)

	var wg sync.WaitGroup
	defer wg.Wait()

	/* Start informers */
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

	/* Start controllers */
	wg.Add(1)
	go func() {
		defer wg.Done()
		ic.Run(ctx, o.ConcurrentSyncs)
	}()

	/* Start proxy */
	wg.Add(1)
	go func() {
		defer wg.Done()

		klog.InfoS("Starting Proxy server")
		err = p.Serve(ctx, listener)
		if err != nil {
			klog.ErrorS(err, "Proxy server failed")
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		<-ctx.Done()
		klog.InfoS("Shutting down Proxy server")
		defer klog.InfoS("Proxy server shut down")

		listener.Close()
	}()

	<-ctx.Done()

	return nil
}
