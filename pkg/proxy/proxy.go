package proxy

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
	"time"

	pausingv1alpha1 "github.com/pausing-clusters-thesis/pausable-scylladb-operator/pkg/api/pausing/v1alpha1"
	pausingv1alpha1client "github.com/pausing-clusters-thesis/pausable-scylladb-operator/pkg/client/pausing/clientset/versioned/typed/pausing/v1alpha1"
	pausingv1alpha1informers "github.com/pausing-clusters-thesis/pausable-scylladb-operator/pkg/client/pausing/informers/externalversions/pausing/v1alpha1"
	pausingv1alpha1listers "github.com/pausing-clusters-thesis/pausable-scylladb-operator/pkg/client/pausing/listers/pausing/v1alpha1"
	"github.com/pausing-clusters-thesis/pausable-scylladb-operator/pkg/controllerhelpers"
	"github.com/pausing-clusters-thesis/pausable-scylladb-operator/pkg/naming"
	"github.com/pausing-clusters-thesis/pausable-scylladb-operator/pkg/readonlyconn"
	"github.com/pausing-clusters-thesis/pausable-scylladb-operator/pkg/routing"
	socontrollerhelpers "github.com/scylladb/scylla-operator/pkg/controllerhelpers"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog/v2"
)

const (
	clientHelloReadDeadline = 3 * time.Second
)

type Proxy struct {
	pausingClient pausingv1alpha1client.PausingV1alpha1Interface

	pausableScyllaDBDatacenterLister pausingv1alpha1listers.PausableScyllaDBDatacenterLister

	cachesToSync []cache.InformerSynced

	routingTable *routing.Table
}

func NewProxy(
	pausingClient pausingv1alpha1client.PausingV1alpha1Interface,
	pausableScyllaDBDatacenterInformer pausingv1alpha1informers.PausableScyllaDBDatacenterInformer,
	routingTable *routing.Table,
) *Proxy {
	return &Proxy{
		pausingClient: pausingClient,

		pausableScyllaDBDatacenterLister: pausableScyllaDBDatacenterInformer.Lister(),

		cachesToSync: []cache.InformerSynced{
			pausableScyllaDBDatacenterInformer.Informer().HasSynced,
		},

		routingTable: routingTable,
	}
}

func (p *Proxy) Serve(ctx context.Context, listener net.Listener) error {
	if !cache.WaitForNamedCacheSync("Proxy", ctx.Done(), p.cachesToSync...) {
		return fmt.Errorf("error waiting for caches to sync")
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			var netErr net.Error
			ok := errors.As(err, &netErr)
			if ok && netErr.Temporary() {
				klog.V(4).InfoS("can't accept connection temporarily: %w", netErr)
				continue
			}

			return err
		}

		go p.handleConnection(ctx, conn)
	}
}

func (p *Proxy) handleConnection(ctx context.Context, conn net.Conn) {
	var err error
	defer conn.Close()

	err = conn.SetReadDeadline(time.Now().Add(clientHelloReadDeadline))
	if err != nil {
		klog.ErrorS(err, "can't set read deadline to client connection")
		return
	}

	clientHello, clientReader, err := readonlyconn.PeekClientHello(conn)
	if err != nil {
		klog.ErrorS(err, "can't peek client hello")
		return
	}

	klog.V(5).InfoS("Read client hello", "ClientHello", clientHello)

	err = conn.SetDeadline(time.Time{})
	if err != nil {
		klog.ErrorS(err, "can't set read deadline to client connection")
		return
	}

	sniHost := clientHello.ServerName
	if len(sniHost) == 0 {
		klog.ErrorS(nil, "invalid server name in client hello", "ServerName", sniHost)
		return
	}

	backendHost, ok := p.routingTable.Get(sniHost)
	if !ok {
		klog.V(4).InfoS("Can't find a mapping for server name in the routing table", "ServerName", sniHost)

		// TODO: validate domain?
		pausableScyllaDBDatacenterDNSDomainPrefix := strings.TrimSuffix(sniHost, fmt.Sprintf(".%s", naming.PausableScyllaDBDatacenterDNSDomainSuffix))
		if pausableScyllaDBDatacenterDNSDomainPrefix == sniHost {
			klog.V(5).InfoS("Server name does not have a suffix referring to PausableScyllaDBClusters.", "ServerName", sniHost)

			return
		}

		prefixSubdomains := strings.Split(pausableScyllaDBDatacenterDNSDomainPrefix, ".")
		prefixSubdomainsLen := len(prefixSubdomains)
		if prefixSubdomainsLen < 2 {
			klog.V(5).InfoS("Server name does not have correct prefix format to reference a PausableScyllaDBDatacenter.", "ServerName", sniHost)

			// DNS subdomains do not match the expected format.
			return
		}

		name := prefixSubdomains[prefixSubdomainsLen-2]
		namespace := prefixSubdomains[prefixSubdomainsLen-1]

		klog.V(5).InfoS("Connection refers to a PausableScyllaDBDatacenter, will try to unpause it.", "ServerName", sniHost, "Ref", klog.KRef(namespace, name))

		var psc *pausingv1alpha1.PausableScyllaDBDatacenter
		psc, err = p.pausableScyllaDBDatacenterLister.PausableScyllaDBDatacenters(namespace).Get(name)
		if err != nil {
			if apierrors.IsNotFound(err) {
				klog.ErrorS(err, "can't get PausableScyllaDBDatacenter", "ServerName", sniHost, "Ref", klog.KRef(namespace, name))
				return
			}

			// Use client in case cache is stale.
			psc, err = p.pausingClient.PausableScyllaDBDatacenters(namespace).Get(ctx, name, metav1.GetOptions{})
			if err != nil {
				if !apierrors.IsNotFound(err) {
					klog.ErrorS(err, "can't get PausableScyllaDBDatacenter", "ServerName", sniHost, "Ref", klog.KRef(namespace, name))
					return
				}

				klog.ErrorS(err, "requested PausableScyllaDBDatacenter does not exist", "ServerName", sniHost, "Ref", klog.KRef(namespace, name))
				return
			}
		}

		klog.V(5).InfoS("Patching PausableScyllaDBDatacenter to unpause it.", "ServerName", sniHost, "Ref", klog.KObj(psc))

		_, err = p.pausingClient.PausableScyllaDBDatacenters(namespace).Patch(
			ctx,
			name,
			types.JSONPatchType,
			[]byte(`[{"op": "replace", "path": "/spec/paused", "value": false}]`),
			metav1.PatchOptions{},
		)
		if err != nil {
			klog.ErrorS(err, "can't patch PausableScyllaDBDatacenter", "Ref", klog.KObj(psc))
			return
		}

		klog.V(5).InfoS("Waiting for PausableScyllaDBDatacenter to roll out.", "ServerName", sniHost, "Ref", klog.KObj(psc))

		_, err = controllerhelpers.WaitForPausableScyllaDBDatacenterState(
			ctx,
			p.pausingClient.PausableScyllaDBDatacenters(namespace),
			name,
			socontrollerhelpers.WaitForStateOptions{
				TolerateDelete: false,
			},
			controllerhelpers.IsPausableScyllaDBDatacenterRolledOut,
		)
		if err != nil {
			klog.ErrorS(err, "can't wait for PausableScyllaDBDatacenter to roll out", "ServerName", sniHost, "Ref", klog.KObj(psc))
			return
		}

		klog.V(5).InfoS("PausableScyllaDBDatacenter rolled out, trying to find a route.", "ServerName", sniHost, "Ref", klog.KObj(psc))

		// TODO: support default backends
		backendHost, ok = p.routingTable.Get(sniHost)
		if !ok {
			klog.ErrorS(nil, "Can't find a mapping for server name in the routing table despite PausableScyllaDBDatacenter being available. This should never happen.", "ServerName", sniHost, "Ref", klog.KObj(psc))
			return
		}

		klog.V(5).InfoS("Route found for PausableScyllaDBDatacenter, will try to proxy the connection.", "ServerName", sniHost, "Ref", klog.KObj(psc), "BackendHost", backendHost)
	}

	backendConn, err := net.Dial("tcp", backendHost)
	if err != nil {
		klog.ErrorS(err, "can't connect to backend server", "ServerName", sniHost, "BackendHost", backendHost)
		return
	}
	defer backendConn.Close()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		io.Copy(conn, backendConn)
		conn.(*net.TCPConn).CloseWrite()
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		io.Copy(backendConn, clientReader)
		backendConn.(*net.TCPConn).CloseWrite()
		wg.Done()
	}()

	wg.Wait()
	return
}
