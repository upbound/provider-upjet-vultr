package clients

import (
	"context"
	"encoding/json"
	"time"

	"github.com/crossplane/crossplane-runtime/v2/pkg/errors"
	"github.com/crossplane/crossplane-runtime/v2/pkg/resource"
	xpv2 "github.com/crossplane/crossplane/apis/v2/core/v2"
	terraformsdk "github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/vultr/govultr/v3"
	"github.com/vultr/terraform-provider-vultr/vultr"
	"golang.org/x/oauth2"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/crossplane/upjet/v2/pkg/terraform"

	clusterv1beta1 "github.com/upbound/provider-vultr/apis/cluster/v1beta1"
	namespacedv1beta1 "github.com/upbound/provider-vultr/apis/namespaced/v1beta1"
)

const (
	// error messages
	errNoProviderConfig     = "no providerConfigRef provided"
	errGetProviderConfig    = "cannot get referenced ProviderConfig"
	errTrackUsage           = "cannot track ProviderConfig usage"
	errExtractCredentials   = "cannot extract credentials"
	errUnmarshalCredentials = "cannot unmarshal vultr credentials as JSON"
	errNoAPIKey             = `credentials secret has no "api_key" key`

	userAgent = "crossplane-provider-vultr"

	keyAPIKey     = "api_key"
	keyRateLimit  = "rate_limit"
	keyRetryLimit = "retry_limit"
)

// TerraformSetupBuilder builds a terraform.SetupFn that configures the
// in-process Vultr Terraform provider (no-fork, no terraform CLI).
func TerraformSetupBuilder() terraform.SetupFn {
	return func(ctx context.Context, client client.Client, mg resource.Managed) (terraform.Setup, error) {
		creds, pcSpec, err := providerCredentials(ctx, client, mg)
		if err != nil {
			return terraform.Setup{}, err
		}

		cfg, err := buildConfiguration(creds, pcSpec.RateLimit, pcSpec.RetryLimit)
		if err != nil {
			return terraform.Setup{}, errors.Wrap(err, "cannot build provider configuration")
		}

		ps := terraform.Setup{Configuration: cfg}
		return ps, errors.Wrap(configureProviderMeta(ctx, &ps), "cannot configure the vultr provider client")
	}
}

// providerCredentials resolves the ProviderConfig referenced by the managed
// resource and extracts the Vultr credentials it points to.
func providerCredentials(ctx context.Context, client client.Client, mg resource.Managed) (map[string]string, *namespacedv1beta1.ProviderConfigSpec, error) {
	pcSpec, err := resolveProviderConfig(ctx, client, mg)
	if err != nil {
		return nil, nil, errors.Wrap(err, "cannot resolve provider config")
	}

	data, err := resource.CommonCredentialExtractor(ctx, pcSpec.Credentials.Source, client, pcSpec.Credentials.CommonCredentialSelectors)
	if err != nil {
		return nil, nil, errors.Wrap(err, errExtractCredentials)
	}

	creds := map[string]string{}
	if err := json.Unmarshal(data, &creds); err != nil {
		return nil, nil, errors.Wrap(err, errUnmarshalCredentials)
	}
	return creds, pcSpec, nil
}

// NewVultrClient returns a Vultr API client configured from the ProviderConfig
// referenced by the managed resource. It serves the controllers that talk to
// the Vultr API directly instead of going through the Terraform provider.
func NewVultrClient(ctx context.Context, client client.Client, mg resource.Managed) (*govultr.Client, error) {
	creds, pcSpec, err := providerCredentials(ctx, client, mg)
	if err != nil {
		return nil, err
	}
	apiKey := creds[keyAPIKey]
	if apiKey == "" {
		return nil, errors.New(errNoAPIKey)
	}

	httpClient := oauth2.NewClient(context.WithoutCancel(ctx), oauth2.StaticTokenSource(&oauth2.Token{AccessToken: apiKey}))
	vultrClient := govultr.NewClient(httpClient)
	vultrClient.SetUserAgent(userAgent)
	if pcSpec.RateLimit != nil && *pcSpec.RateLimit != 0 {
		vultrClient.SetRateLimit(time.Duration(*pcSpec.RateLimit) * time.Millisecond)
	}
	if pcSpec.RetryLimit != nil && *pcSpec.RetryLimit != 0 {
		vultrClient.SetRetryLimit(int(*pcSpec.RetryLimit))
	}
	return vultrClient, nil
}

func buildConfiguration(creds map[string]string, rateLimit, retryLimit *int64) (map[string]any, error) {
	apiKey := creds[keyAPIKey]
	if apiKey == "" {
		return nil, errors.New(errNoAPIKey)
	}
	cfg := map[string]any{keyAPIKey: apiKey}
	// The SDKv2 provider schema declares these as TypeInt; pass plain ints so
	// the field reader's weak decoding never has to guess.
	if rateLimit != nil {
		cfg[keyRateLimit] = int(*rateLimit)
	}
	if retryLimit != nil {
		cfg[keyRetryLimit] = int(*retryLimit)
	}
	return cfg, nil
}

// configureProviderMeta configures a fresh in-process SDKv2 provider instance
// with the resolved configuration and stores the resulting provider meta (the
// Vultr API client) in the Setup, as required by upjet's TerraformPluginSDK
// connectors. A fresh instance is used per call because provider meta is
// per-ProviderConfig state.
func configureProviderMeta(ctx context.Context, ps *terraform.Setup) error {
	p := vultr.Provider()
	diag := p.Configure(context.WithoutCancel(ctx), &terraformsdk.ResourceConfig{Config: ps.Configuration})
	if diag != nil && diag.HasError() {
		return errors.Errorf("failed to configure the vultr provider: %v", diag)
	}
	ps.Meta = p.Meta()
	return nil
}

func toSharedPCSpec(pc *clusterv1beta1.ProviderConfig) (*namespacedv1beta1.ProviderConfigSpec, error) {
	if pc == nil {
		return nil, nil
	}
	data, err := json.Marshal(pc.Spec)
	if err != nil {
		return nil, err
	}

	var mSpec namespacedv1beta1.ProviderConfigSpec
	err = json.Unmarshal(data, &mSpec)
	return &mSpec, err
}

func resolveProviderConfig(ctx context.Context, crClient client.Client, mg resource.Managed) (*namespacedv1beta1.ProviderConfigSpec, error) {
	switch managed := mg.(type) {
	case resource.LegacyManaged:
		return resolveLegacy(ctx, crClient, managed)
	case resource.ModernManaged:
		return resolveModern(ctx, crClient, managed)
	default:
		return nil, errors.New("resource is not a managed resource")
	}
}

func resolveLegacy(ctx context.Context, client client.Client, mg resource.LegacyManaged) (*namespacedv1beta1.ProviderConfigSpec, error) {
	configRef := mg.GetProviderConfigReference()
	if configRef == nil {
		return nil, errors.New(errNoProviderConfig)
	}
	pc := &clusterv1beta1.ProviderConfig{}
	if err := client.Get(ctx, types.NamespacedName{Name: configRef.Name}, pc); err != nil {
		return nil, errors.Wrap(err, errGetProviderConfig)
	}

	t := resource.NewLegacyProviderConfigUsageTracker(client, &clusterv1beta1.ProviderConfigUsage{})
	if err := t.Track(ctx, mg); err != nil {
		return nil, errors.Wrap(err, errTrackUsage)
	}

	return toSharedPCSpec(pc)
}

// resolveNamespacedSpec converts a namespaced ProviderConfig spec into the
// internal, fully-resolved ProviderConfigSpec. The namespaced spec omits the
// secret namespace, so it is resolved to the namespace of the referencing
// managed resource here.
func resolveNamespacedSpec(spec namespacedv1beta1.NamespacedProviderConfigSpec, namespace string) namespacedv1beta1.ProviderConfigSpec {
	resolved := namespacedv1beta1.ProviderConfigSpec{
		ReconciliationPolicy: spec.ReconciliationPolicy,
		RateLimit:            spec.RateLimit,
		RetryLimit:           spec.RetryLimit,
		Credentials: namespacedv1beta1.ProviderCredentials{
			Source: spec.Credentials.Source,
			CommonCredentialSelectors: xpv2.CommonCredentialSelectors{
				Fs:  spec.Credentials.Fs,
				Env: spec.Credentials.Env,
			},
		},
	}
	if spec.Credentials.SecretRef != nil {
		resolved.Credentials.SecretRef = spec.Credentials.SecretRef.ToSecretKeySelector(namespace)
	}
	return resolved
}

func resolveModern(ctx context.Context, crClient client.Client, mg resource.ModernManaged) (*namespacedv1beta1.ProviderConfigSpec, error) {
	configRef := mg.GetProviderConfigReference()
	if configRef == nil {
		return nil, errors.New(errNoProviderConfig)
	}

	pcRuntimeObj, err := crClient.Scheme().New(namespacedv1beta1.SchemeGroupVersion.WithKind(configRef.Kind))
	if err != nil {
		return nil, errors.Wrap(err, "unknown GVK for ProviderConfig")
	}
	pcObj, ok := pcRuntimeObj.(client.Object)
	if !ok {
		// This indicates a programming error, types are not properly generated
		return nil, errors.New("runtime object is not a client.Object")
	}

	// Namespace will be ignored if the PC is a cluster-scoped type
	if err := crClient.Get(ctx, types.NamespacedName{Name: configRef.Name, Namespace: mg.GetNamespace()}, pcObj); err != nil {
		return nil, errors.Wrap(err, errGetProviderConfig)
	}

	var pcSpec namespacedv1beta1.ProviderConfigSpec
	pcu := &namespacedv1beta1.ProviderConfigUsage{}
	switch pc := pcObj.(type) {
	case *namespacedv1beta1.ProviderConfig:
		pcSpec = resolveNamespacedSpec(pc.Spec, mg.GetNamespace())
	case *namespacedv1beta1.ClusterProviderConfig:
		pcSpec = pc.Spec
	default:
		return nil, errors.New("unknown provider config type")
	}
	t := resource.NewProviderConfigUsageTracker(crClient, pcu)
	if err := t.Track(ctx, mg); err != nil {
		return nil, errors.Wrap(err, errTrackUsage)
	}
	return &pcSpec, nil
}
