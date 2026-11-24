// Package clusterauth contains the controller that publishes the kubeconfig
// of a VKE cluster as the connection details of a ClusterAuth.
package clusterauth

import (
	"context"
	"encoding/base64"
	"net/http"
	"time"

	"github.com/crossplane/crossplane-runtime/v2/pkg/errors"
	"github.com/crossplane/crossplane-runtime/v2/pkg/event"
	"github.com/crossplane/crossplane-runtime/v2/pkg/meta"
	"github.com/crossplane/crossplane-runtime/v2/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/v2/pkg/resource"
	xpv2 "github.com/crossplane/crossplane/apis/v2/core/v2"
	tjcontroller "github.com/crossplane/upjet/v2/pkg/controller"
	ujresource "github.com/crossplane/upjet/v2/pkg/resource"
	"github.com/vultr/govultr/v3"
	"k8s.io/client-go/tools/clientcmd"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/upbound/provider-vultr/apis/cluster/vke/v1beta1"
	"github.com/upbound/provider-vultr/internal/clients"
	"github.com/upbound/provider-vultr/internal/features"
)

const (
	errNotClusterAuth   = "managed resource is not a ClusterAuth custom resource"
	errNoClusterID      = "the ID of the cluster is not set"
	errNewClient        = "cannot create a Vultr API client"
	errGetKubeconfig    = "cannot get the kubeconfig of the cluster"
	errDecodeKubeconfig = "cannot base64-decode the kubeconfig of the cluster"
	errParseKubeconfig  = "cannot parse the kubeconfig of the cluster"
	errNoCurrentContext = "the kubeconfig of the cluster has no usable current context"

	// The kubeconfig is only read, so a short timeout is enough to surface a
	// stuck Vultr API call without holding a reconcile worker.
	reconcileTimeout = 1 * time.Minute

	keyKubeconfig = "kubeconfig"
	keyEndpoint   = "endpoint"
	keyClusterCA  = "clusterCA"
	keyClientCert = "clientCert"
	keyClientKey  = "clientKey"
)

// SetupWebhookWithManager registers the conversion webhook for ClusterAuth.
func SetupWebhookWithManager(mgr ctrl.Manager) error {
	if err := ctrl.NewWebhookManagedBy(mgr, &v1beta1.ClusterAuth{}).
		Complete(); err != nil {
		return errors.Wrap(err, "cannot register webhook for the kind v1beta1.ClusterAuth")
	}
	return nil
}

// Setup adds a controller that reconciles ClusterAuth.
func Setup(mgr ctrl.Manager, o tjcontroller.Options) error {
	name := managed.ControllerName(v1beta1.ClusterAuth_GroupKind)

	opts := []managed.ReconcilerOption{
		managed.WithExternalConnector(&connector{
			kube:        mgr.GetClient(),
			newClientFn: newKubeconfigGetter,
		}),
		managed.WithTimeout(reconcileTimeout),
		managed.WithPollInterval(o.PollInterval),
		managed.WithLogger(o.Logger.WithValues("controller", name)),
		managed.WithRecorder(event.NewAPIRecorder(mgr.GetEventRecorderFor(name))),
	}
	if o.Features.Enabled(features.EnableBetaManagementPolicies) {
		opts = append(opts, managed.WithManagementPolicies())
	}

	return ctrl.NewControllerManagedBy(mgr).
		Named(name).
		WithOptions(o.ForControllerRuntime()).
		For(&v1beta1.ClusterAuth{}).
		Complete(managed.NewReconciler(mgr,
			resource.ManagedKind(v1beta1.ClusterAuth_GroupVersionKind),
			opts...))
}

// SetupGated adds a controller that reconciles ClusterAuth once its CRD is
// established.
func SetupGated(mgr ctrl.Manager, o tjcontroller.Options) error {
	o.Gate.Register(func() {
		if err := Setup(mgr, o); err != nil {
			mgr.GetLogger().Error(err, "unable to setup reconciler", "gvk", v1beta1.ClusterAuth_GroupVersionKind.String())
		}
	}, v1beta1.ClusterAuth_GroupVersionKind)
	return nil
}

// kubeconfigGetter is the part of the Vultr Kubernetes API a ClusterAuth needs.
type kubeconfigGetter interface {
	GetKubeConfig(ctx context.Context, vkeID string) (*govultr.KubeConfig, *http.Response, error)
}

func newKubeconfigGetter(ctx context.Context, kube client.Client, mg resource.Managed) (kubeconfigGetter, error) {
	c, err := clients.NewVultrClient(ctx, kube, mg)
	if err != nil {
		return nil, err
	}
	return c.Kubernetes, nil
}

type connector struct {
	kube        client.Client
	newClientFn func(ctx context.Context, kube client.Client, mg resource.Managed) (kubeconfigGetter, error)
}

func (c *connector) Connect(ctx context.Context, mg resource.Managed) (managed.ExternalClient, error) {
	getter, err := c.newClientFn(ctx, c.kube, mg)
	if err != nil {
		return nil, errors.Wrap(err, errNewClient)
	}
	return &external{client: getter}, nil
}

type external struct {
	client kubeconfigGetter
}

// Observe fetches the kubeconfig on every reconcile and publishes it. There
// is no external resource behind a ClusterAuth, so it exists for as long as
// the managed resource is not being deleted, and it is never out of date.
func (e *external) Observe(ctx context.Context, mg resource.Managed) (managed.ExternalObservation, error) {
	cr, ok := mg.(*v1beta1.ClusterAuth)
	if !ok {
		return managed.ExternalObservation{}, errors.New(errNotClusterAuth)
	}
	if meta.WasDeleted(cr) {
		return managed.ExternalObservation{ResourceExists: false}, nil
	}
	clusterID := cr.Spec.ForProvider.ClusterID
	if clusterID == nil || *clusterID == "" {
		return managed.ExternalObservation{}, errors.New(errNoClusterID)
	}

	kc, _, err := e.client.GetKubeConfig(ctx, *clusterID)
	if err != nil {
		return managed.ExternalObservation{}, errors.Wrap(err, errGetKubeconfig)
	}
	conn, endpoint, err := connectionDetails(kc.KubeConfig)
	if err != nil {
		return managed.ExternalObservation{}, err
	}

	cr.Status.AtProvider.Endpoint = &endpoint
	cr.Status.SetConditions(xpv2.Available())
	ujresource.SetUpToDateCondition(mg, true)
	return managed.ExternalObservation{
		ResourceExists:    true,
		ResourceUpToDate:  true,
		ConnectionDetails: conn,
	}, nil
}

// connectionDetails turns the base64-encoded kubeconfig returned by the Vultr
// API into connection details, and returns the API server endpoint.
func connectionDetails(encoded string) (managed.ConnectionDetails, string, error) {
	raw, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, "", errors.Wrap(err, errDecodeKubeconfig)
	}
	cfg, err := clientcmd.Load(raw)
	if err != nil {
		return nil, "", errors.Wrap(err, errParseKubeconfig)
	}
	kctx, ok := cfg.Contexts[cfg.CurrentContext]
	if !ok {
		return nil, "", errors.New(errNoCurrentContext)
	}
	cluster, ok := cfg.Clusters[kctx.Cluster]
	if !ok {
		return nil, "", errors.New(errNoCurrentContext)
	}
	user, ok := cfg.AuthInfos[kctx.AuthInfo]
	if !ok {
		return nil, "", errors.New(errNoCurrentContext)
	}

	return managed.ConnectionDetails{
		keyKubeconfig: raw,
		keyEndpoint:   []byte(cluster.Server),
		keyClusterCA:  cluster.CertificateAuthorityData,
		keyClientCert: user.ClientCertificateData,
		keyClientKey:  user.ClientKeyData,
	}, cluster.Server, nil
}

func (e *external) Create(_ context.Context, _ resource.Managed) (managed.ExternalCreation, error) {
	return managed.ExternalCreation{}, nil
}

func (e *external) Update(_ context.Context, _ resource.Managed) (managed.ExternalUpdate, error) {
	return managed.ExternalUpdate{}, nil
}

func (e *external) Delete(_ context.Context, _ resource.Managed) (managed.ExternalDelete, error) {
	return managed.ExternalDelete{}, nil
}

func (e *external) Disconnect(_ context.Context) error {
	return nil
}
