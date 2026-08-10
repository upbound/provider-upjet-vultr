// SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: CC0-1.0

package config

import (
	// Note(turkenh): we are importing this to embed provider schema document
	"context"
	_ "embed"

	ujconfig "github.com/crossplane/upjet/v2/pkg/config"
	"github.com/vultr/terraform-provider-vultr/vultr"

	"github.com/upbound/provider-vultr/config/cluster"
	"github.com/upbound/provider-vultr/config/templates"
)

const (
	resourcePrefix = "vultr"
	modulePath     = "github.com/upbound/provider-vultr"
	versionV1Beta1 = "v1beta1"
)

//go:embed schema.json
var providerSchema string

//go:embed provider-metadata.yaml
var providerMetadata string

// GetProvider returns provider configuration
func GetProvider(_ context.Context) (*ujconfig.Provider, error) {
	sdkProvider := vultr.Provider()

	defaultResourceOptions := []ujconfig.ResourceOption{
		GroupKindOverrides(),
		ExternalNameConfigurations(),
	}

	pc := ujconfig.NewProvider(
		[]byte(providerSchema), resourcePrefix, modulePath, []byte(providerMetadata),
		ujconfig.WithRootGroup("vultr.upbound.io"),
		ujconfig.WithIncludeList([]string{}),
		ujconfig.WithControllerTemplate(templates.ControllerTemplate),
		ujconfig.WithTerraformPluginSDKIncludeList(ExternalNameConfigured()),
		ujconfig.WithFeaturesPackage("internal/features"),
		ujconfig.WithTerraformProvider(sdkProvider),
		ujconfig.WithSchemaTraversers(&ujconfig.SingletonListEmbedder{}),
		ujconfig.WithDefaultResourceOptions(defaultResourceOptions...),
	)

	// add custom config functions
	for _, configure := range cluster.ProviderConfiguration {
		configure(pc)
	}

	pc.ConfigureResources()

	registerTFConversions(pc)

	return pc, nil
}
