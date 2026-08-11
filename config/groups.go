// SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: CC0-1.0

package config

import (
	"strings"

	"github.com/crossplane/upjet/v2/pkg/config"
	"github.com/crossplane/upjet/v2/pkg/types/name"
)

// GroupKindOverrides overrides the group and kind of the resource if it matches
// any entry in the GroupMap.
func GroupKindOverrides() config.ResourceOption {
	return func(r *config.Resource) {
		if f, ok := GroupMap[r.Name]; ok {
			r.ShortGroup, r.Kind = f(r.Name)
		}
	}
}

// GroupKindCalculator returns the correct group and kind name for given TF
// resource.
type GroupKindCalculator func(resource string) (string, string)

// ReplaceGroupWords uses given group as the group of the resource and removes
// a number of words in resource name before calculating the kind of the resource.
func ReplaceGroupWords(group string, count int) GroupKindCalculator {
	return func(resource string) (string, string) {
		// "vultr_database_connection_pool" -> (database, ConnectionPool)
		words := strings.Split(strings.TrimPrefix(resource, "vultr_"), "_")
		snakeKind := strings.Join(words[count:], "_")
		return group, name.NewFromSnake(snakeKind).Camel
	}
}

// KnownGroupKind returns a GroupKindCalculator that assigns the given static
// group and kind regardless of the resource name.
func KnownGroupKind(group, kind string) GroupKindCalculator {
	return func(string) (string, string) { return group, kind }
}

// GroupMap assigns every Vultr resource to one of a small set of service-area
// API groups (compute, database, iam, network, storage, vke). Kinds carry the
// full resource noun since the group name alone does not disambiguate; in the
// database group the TF prefix doubles as the group and is dropped from the
// kind. The deprecated vultr_vpc2 is intentionally absent.
var GroupMap = map[string]GroupKindCalculator{
	// compute
	"vultr_bare_metal_server": KnownGroupKind("compute", "BareMetalServer"),
	"vultr_inference":         KnownGroupKind("compute", "Inference"),
	"vultr_instance":          KnownGroupKind("compute", "Instance"),
	"vultr_instance_ipv4":     KnownGroupKind("compute", "InstanceIPv4"),
	"vultr_iso_private":       KnownGroupKind("compute", "ISOPrivate"),
	"vultr_snapshot":          KnownGroupKind("compute", "Snapshot"),
	"vultr_snapshot_from_url": KnownGroupKind("compute", "SnapshotFromURL"),
	"vultr_ssh_key":           KnownGroupKind("compute", "SSHKey"),
	"vultr_startup_script":    KnownGroupKind("compute", "StartupScript"),

	// database
	"vultr_database":                 KnownGroupKind("database", "Database"),
	"vultr_database_connection_pool": KnownGroupKind("database", "ConnectionPool"),
	"vultr_database_connector":       KnownGroupKind("database", "Connector"),
	"vultr_database_db":              KnownGroupKind("database", "DB"),
	"vultr_database_quota":           KnownGroupKind("database", "Quota"),
	"vultr_database_replica":         KnownGroupKind("database", "Replica"),
	"vultr_database_topic":           KnownGroupKind("database", "Topic"),
	"vultr_database_user":            KnownGroupKind("database", "User"),

	// iam
	"vultr_oidc_issuer":                          KnownGroupKind("iam", "OIDCIssuer"),
	"vultr_oidc_provider":                        KnownGroupKind("iam", "OIDCProvider"),
	"vultr_oidc_token":                           KnownGroupKind("iam", "OIDCToken"),
	"vultr_organization":                         KnownGroupKind("iam", "Organization"),
	"vultr_organization_group":                   KnownGroupKind("iam", "OrganizationGroup"),
	"vultr_organization_invitation":              KnownGroupKind("iam", "OrganizationInvitation"),
	"vultr_organization_policy":                  KnownGroupKind("iam", "OrganizationPolicy"),
	"vultr_organization_policy_group_attachment": KnownGroupKind("iam", "OrganizationPolicyGroupAttachment"),
	"vultr_organization_policy_user_attachment":  KnownGroupKind("iam", "OrganizationPolicyUserAttachment"),
	"vultr_organization_role":                    KnownGroupKind("iam", "OrganizationRole"),
	"vultr_organization_role_group_attachment":   KnownGroupKind("iam", "OrganizationRoleGroupAttachment"),
	"vultr_organization_role_policy_attachment":  KnownGroupKind("iam", "OrganizationRolePolicyAttachment"),
	"vultr_organization_role_session":            KnownGroupKind("iam", "OrganizationRoleSession"),
	"vultr_organization_role_trust":              KnownGroupKind("iam", "OrganizationRoleTrust"),
	"vultr_user":                                 KnownGroupKind("iam", "User"),

	// network
	"vultr_dns_domain":                       KnownGroupKind("network", "DNSDomain"),
	"vultr_dns_record":                       KnownGroupKind("network", "DNSRecord"),
	"vultr_firewall_group":                   KnownGroupKind("network", "FirewallGroup"),
	"vultr_firewall_rule":                    KnownGroupKind("network", "FirewallRule"),
	"vultr_load_balancer":                    KnownGroupKind("network", "LoadBalancer"),
	"vultr_nat_gateway":                      KnownGroupKind("network", "NATGateway"),
	"vultr_nat_gateway_firewall_rule":        KnownGroupKind("network", "NATGatewayFirewallRule"),
	"vultr_nat_gateway_port_forwarding_rule": KnownGroupKind("network", "NATGatewayPortForwardingRule"),
	"vultr_reserved_ip":                      KnownGroupKind("network", "ReservedIP"),
	"vultr_reverse_ipv4":                     KnownGroupKind("network", "ReverseIPv4"),
	"vultr_reverse_ipv6":                     KnownGroupKind("network", "ReverseIPv6"),
	"vultr_vpc":                              KnownGroupKind("network", "VPC"),

	// storage
	"vultr_block_storage":               KnownGroupKind("storage", "BlockStorage"),
	"vultr_container_registry":          KnownGroupKind("storage", "ContainerRegistry"),
	"vultr_object_storage":              KnownGroupKind("storage", "ObjectStorage"),
	"vultr_object_storage_bucket":       KnownGroupKind("storage", "ObjectStorageBucket"),
	"vultr_virtual_file_system_storage": KnownGroupKind("storage", "VirtualFileSystemStorage"),

	// vke
	"vultr_kubernetes":            KnownGroupKind("vke", "Kubernetes"),
	"vultr_kubernetes_node_pools": KnownGroupKind("vke", "KubernetesNodePool"),
}
