package config

import (
	"github.com/crossplane/upjet/v2/pkg/config"
)

// ExternalNameConfigs contains all external name configurations for this
// provider.
var ExternalNameConfigs = map[string]config.ExternalName{
	// compute
	"vultr_bare_metal_server": config.IdentifierFromProvider,
	"vultr_inference":         config.IdentifierFromProvider,
	"vultr_instance":          config.IdentifierFromProvider,
	"vultr_instance_ipv4":     config.IdentifierFromProvider,
	"vultr_iso_private":       config.IdentifierFromProvider,
	"vultr_snapshot":          config.IdentifierFromProvider,
	"vultr_snapshot_from_url": config.IdentifierFromProvider,
	"vultr_ssh_key":           config.IdentifierFromProvider,
	"vultr_startup_script":    config.IdentifierFromProvider,

	// database
	"vultr_database":                 config.IdentifierFromProvider,
	"vultr_database_connection_pool": config.IdentifierFromProvider,
	"vultr_database_connector":       config.IdentifierFromProvider,
	"vultr_database_db":              config.IdentifierFromProvider,
	"vultr_database_quota":           config.IdentifierFromProvider,
	"vultr_database_replica":         config.IdentifierFromProvider,
	"vultr_database_topic":           config.IdentifierFromProvider,
	"vultr_database_user":            config.IdentifierFromProvider,

	// iam
	"vultr_oidc_issuer":                          config.IdentifierFromProvider,
	"vultr_oidc_provider":                        config.IdentifierFromProvider,
	"vultr_oidc_token":                           config.IdentifierFromProvider,
	"vultr_organization":                         config.IdentifierFromProvider,
	"vultr_organization_group":                   config.IdentifierFromProvider,
	"vultr_organization_invitation":              config.IdentifierFromProvider,
	"vultr_organization_policy":                  config.IdentifierFromProvider,
	"vultr_organization_policy_group_attachment": config.IdentifierFromProvider,
	"vultr_organization_policy_user_attachment":  config.IdentifierFromProvider,
	"vultr_organization_role":                    config.IdentifierFromProvider,
	"vultr_organization_role_group_attachment":   config.IdentifierFromProvider,
	"vultr_organization_role_policy_attachment":  config.IdentifierFromProvider,
	"vultr_organization_role_session":            config.IdentifierFromProvider,
	"vultr_organization_role_trust":              config.IdentifierFromProvider,
	"vultr_user":                                 config.IdentifierFromProvider,

	// network
	"vultr_dns_domain":                       config.ParameterAsIdentifier("domain"),
	"vultr_dns_record":                       config.IdentifierFromProvider,
	"vultr_firewall_group":                   config.IdentifierFromProvider,
	"vultr_firewall_rule":                    config.IdentifierFromProvider,
	"vultr_load_balancer":                    config.IdentifierFromProvider,
	"vultr_nat_gateway":                      config.IdentifierFromProvider,
	"vultr_nat_gateway_firewall_rule":        config.IdentifierFromProvider,
	"vultr_nat_gateway_port_forwarding_rule": config.IdentifierFromProvider,
	"vultr_reserved_ip":                      config.IdentifierFromProvider,
	"vultr_reverse_ipv4":                     config.IdentifierFromProvider,
	"vultr_reverse_ipv6":                     config.IdentifierFromProvider,
	"vultr_vpc":                              config.IdentifierFromProvider,

	// storage
	"vultr_block_storage":               config.IdentifierFromProvider,
	"vultr_container_registry":          config.IdentifierFromProvider,
	"vultr_object_storage":              config.IdentifierFromProvider,
	"vultr_object_storage_bucket":       config.IdentifierFromProvider,
	"vultr_virtual_file_system_storage": config.IdentifierFromProvider,

	// vke
	"vultr_kubernetes":            config.IdentifierFromProvider,
	"vultr_kubernetes_node_pools": config.IdentifierFromProvider,
}

// ExternalNameConfigurations applies all external name configs listed in the
// table ExternalNameConfigs and sets the version of those resources to v1beta1
// assuming they will be tested.
func ExternalNameConfigurations() config.ResourceOption {
	return func(r *config.Resource) {
		e, configured := ExternalNameConfigs[r.Name]
		if !configured {
			return
		}
		r.ExternalName = e
		r.Version = versionV1Beta1
	}
}

// ExternalNameConfigured returns the list of all resources whose external name
// is configured manually.
func ExternalNameConfigured() []string {
	l := make([]string, len(ExternalNameConfigs))
	i := 0
	for name := range ExternalNameConfigs {
		// $ is added to match the exact string since the format is regex.
		l[i] = name + "$"
		i++
	}
	return l
}
