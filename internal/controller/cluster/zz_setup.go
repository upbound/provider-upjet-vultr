// SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/upjet/v2/pkg/controller"

	baremetalserver "github.com/upbound/provider-vultr/internal/controller/cluster/compute/baremetalserver"
	inference "github.com/upbound/provider-vultr/internal/controller/cluster/compute/inference"
	instance "github.com/upbound/provider-vultr/internal/controller/cluster/compute/instance"
	instanceipv4 "github.com/upbound/provider-vultr/internal/controller/cluster/compute/instanceipv4"
	isoprivate "github.com/upbound/provider-vultr/internal/controller/cluster/compute/isoprivate"
	snapshot "github.com/upbound/provider-vultr/internal/controller/cluster/compute/snapshot"
	snapshotfromurl "github.com/upbound/provider-vultr/internal/controller/cluster/compute/snapshotfromurl"
	sshkey "github.com/upbound/provider-vultr/internal/controller/cluster/compute/sshkey"
	startupscript "github.com/upbound/provider-vultr/internal/controller/cluster/compute/startupscript"
	connectionpool "github.com/upbound/provider-vultr/internal/controller/cluster/database/connectionpool"
	connector "github.com/upbound/provider-vultr/internal/controller/cluster/database/connector"
	database "github.com/upbound/provider-vultr/internal/controller/cluster/database/database"
	db "github.com/upbound/provider-vultr/internal/controller/cluster/database/db"
	quota "github.com/upbound/provider-vultr/internal/controller/cluster/database/quota"
	replica "github.com/upbound/provider-vultr/internal/controller/cluster/database/replica"
	topic "github.com/upbound/provider-vultr/internal/controller/cluster/database/topic"
	user "github.com/upbound/provider-vultr/internal/controller/cluster/database/user"
	oidcissuer "github.com/upbound/provider-vultr/internal/controller/cluster/iam/oidcissuer"
	oidcprovider "github.com/upbound/provider-vultr/internal/controller/cluster/iam/oidcprovider"
	oidctoken "github.com/upbound/provider-vultr/internal/controller/cluster/iam/oidctoken"
	organization "github.com/upbound/provider-vultr/internal/controller/cluster/iam/organization"
	organizationgroup "github.com/upbound/provider-vultr/internal/controller/cluster/iam/organizationgroup"
	organizationinvitation "github.com/upbound/provider-vultr/internal/controller/cluster/iam/organizationinvitation"
	organizationpolicy "github.com/upbound/provider-vultr/internal/controller/cluster/iam/organizationpolicy"
	organizationpolicygroupattachment "github.com/upbound/provider-vultr/internal/controller/cluster/iam/organizationpolicygroupattachment"
	organizationpolicyuserattachment "github.com/upbound/provider-vultr/internal/controller/cluster/iam/organizationpolicyuserattachment"
	organizationrole "github.com/upbound/provider-vultr/internal/controller/cluster/iam/organizationrole"
	organizationrolegroupattachment "github.com/upbound/provider-vultr/internal/controller/cluster/iam/organizationrolegroupattachment"
	organizationrolepolicyattachment "github.com/upbound/provider-vultr/internal/controller/cluster/iam/organizationrolepolicyattachment"
	organizationrolesession "github.com/upbound/provider-vultr/internal/controller/cluster/iam/organizationrolesession"
	organizationroletrust "github.com/upbound/provider-vultr/internal/controller/cluster/iam/organizationroletrust"
	useriam "github.com/upbound/provider-vultr/internal/controller/cluster/iam/user"
	dnsdomain "github.com/upbound/provider-vultr/internal/controller/cluster/network/dnsdomain"
	dnsrecord "github.com/upbound/provider-vultr/internal/controller/cluster/network/dnsrecord"
	firewallgroup "github.com/upbound/provider-vultr/internal/controller/cluster/network/firewallgroup"
	firewallrule "github.com/upbound/provider-vultr/internal/controller/cluster/network/firewallrule"
	loadbalancer "github.com/upbound/provider-vultr/internal/controller/cluster/network/loadbalancer"
	natgateway "github.com/upbound/provider-vultr/internal/controller/cluster/network/natgateway"
	natgatewayfirewallrule "github.com/upbound/provider-vultr/internal/controller/cluster/network/natgatewayfirewallrule"
	natgatewayportforwardingrule "github.com/upbound/provider-vultr/internal/controller/cluster/network/natgatewayportforwardingrule"
	reservedip "github.com/upbound/provider-vultr/internal/controller/cluster/network/reservedip"
	reverseipv4 "github.com/upbound/provider-vultr/internal/controller/cluster/network/reverseipv4"
	reverseipv6 "github.com/upbound/provider-vultr/internal/controller/cluster/network/reverseipv6"
	vpc "github.com/upbound/provider-vultr/internal/controller/cluster/network/vpc"
	providerconfig "github.com/upbound/provider-vultr/internal/controller/cluster/providerconfig"
	blockstorage "github.com/upbound/provider-vultr/internal/controller/cluster/storage/blockstorage"
	containerregistry "github.com/upbound/provider-vultr/internal/controller/cluster/storage/containerregistry"
	objectstorage "github.com/upbound/provider-vultr/internal/controller/cluster/storage/objectstorage"
	objectstoragebucket "github.com/upbound/provider-vultr/internal/controller/cluster/storage/objectstoragebucket"
	virtualfilesystemstorage "github.com/upbound/provider-vultr/internal/controller/cluster/storage/virtualfilesystemstorage"
	kubernetes "github.com/upbound/provider-vultr/internal/controller/cluster/vke/kubernetes"
	kubernetesnodepool "github.com/upbound/provider-vultr/internal/controller/cluster/vke/kubernetesnodepool"
)

// Setup creates all controllers with the supplied logger and adds them to
// the supplied manager.
func Setup(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		baremetalserver.Setup,
		inference.Setup,
		instance.Setup,
		instanceipv4.Setup,
		isoprivate.Setup,
		snapshot.Setup,
		snapshotfromurl.Setup,
		sshkey.Setup,
		startupscript.Setup,
		connectionpool.Setup,
		connector.Setup,
		database.Setup,
		db.Setup,
		quota.Setup,
		replica.Setup,
		topic.Setup,
		user.Setup,
		oidcissuer.Setup,
		oidcprovider.Setup,
		oidctoken.Setup,
		organization.Setup,
		organizationgroup.Setup,
		organizationinvitation.Setup,
		organizationpolicy.Setup,
		organizationpolicygroupattachment.Setup,
		organizationpolicyuserattachment.Setup,
		organizationrole.Setup,
		organizationrolegroupattachment.Setup,
		organizationrolepolicyattachment.Setup,
		organizationrolesession.Setup,
		organizationroletrust.Setup,
		useriam.Setup,
		dnsdomain.Setup,
		dnsrecord.Setup,
		firewallgroup.Setup,
		firewallrule.Setup,
		loadbalancer.Setup,
		natgateway.Setup,
		natgatewayfirewallrule.Setup,
		natgatewayportforwardingrule.Setup,
		reservedip.Setup,
		reverseipv4.Setup,
		reverseipv6.Setup,
		vpc.Setup,
		providerconfig.Setup,
		blockstorage.Setup,
		containerregistry.Setup,
		objectstorage.Setup,
		objectstoragebucket.Setup,
		virtualfilesystemstorage.Setup,
		kubernetes.Setup,
		kubernetesnodepool.Setup,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}

// SetupGated creates all controllers with the supplied logger and adds them to
// the supplied manager gated.
func SetupGated(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		baremetalserver.SetupGated,
		inference.SetupGated,
		instance.SetupGated,
		instanceipv4.SetupGated,
		isoprivate.SetupGated,
		snapshot.SetupGated,
		snapshotfromurl.SetupGated,
		sshkey.SetupGated,
		startupscript.SetupGated,
		connectionpool.SetupGated,
		connector.SetupGated,
		database.SetupGated,
		db.SetupGated,
		quota.SetupGated,
		replica.SetupGated,
		topic.SetupGated,
		user.SetupGated,
		oidcissuer.SetupGated,
		oidcprovider.SetupGated,
		oidctoken.SetupGated,
		organization.SetupGated,
		organizationgroup.SetupGated,
		organizationinvitation.SetupGated,
		organizationpolicy.SetupGated,
		organizationpolicygroupattachment.SetupGated,
		organizationpolicyuserattachment.SetupGated,
		organizationrole.SetupGated,
		organizationrolegroupattachment.SetupGated,
		organizationrolepolicyattachment.SetupGated,
		organizationrolesession.SetupGated,
		organizationroletrust.SetupGated,
		useriam.SetupGated,
		dnsdomain.SetupGated,
		dnsrecord.SetupGated,
		firewallgroup.SetupGated,
		firewallrule.SetupGated,
		loadbalancer.SetupGated,
		natgateway.SetupGated,
		natgatewayfirewallrule.SetupGated,
		natgatewayportforwardingrule.SetupGated,
		reservedip.SetupGated,
		reverseipv4.SetupGated,
		reverseipv6.SetupGated,
		vpc.SetupGated,
		providerconfig.SetupGated,
		blockstorage.SetupGated,
		containerregistry.SetupGated,
		objectstorage.SetupGated,
		objectstoragebucket.SetupGated,
		virtualfilesystemstorage.SetupGated,
		kubernetes.SetupGated,
		kubernetesnodepool.SetupGated,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}

// SetupWebhookWithManager registers conversion webhooks for all resource kinds in the group.
func SetupWebhookWithManager(mgr ctrl.Manager) error {
	for _, setup := range []func(ctrl.Manager) error{
		baremetalserver.SetupWebhookWithManager,
		inference.SetupWebhookWithManager,
		instance.SetupWebhookWithManager,
		instanceipv4.SetupWebhookWithManager,
		isoprivate.SetupWebhookWithManager,
		snapshot.SetupWebhookWithManager,
		snapshotfromurl.SetupWebhookWithManager,
		sshkey.SetupWebhookWithManager,
		startupscript.SetupWebhookWithManager,
		connectionpool.SetupWebhookWithManager,
		connector.SetupWebhookWithManager,
		database.SetupWebhookWithManager,
		db.SetupWebhookWithManager,
		quota.SetupWebhookWithManager,
		replica.SetupWebhookWithManager,
		topic.SetupWebhookWithManager,
		user.SetupWebhookWithManager,
		oidcissuer.SetupWebhookWithManager,
		oidcprovider.SetupWebhookWithManager,
		oidctoken.SetupWebhookWithManager,
		organization.SetupWebhookWithManager,
		organizationgroup.SetupWebhookWithManager,
		organizationinvitation.SetupWebhookWithManager,
		organizationpolicy.SetupWebhookWithManager,
		organizationpolicygroupattachment.SetupWebhookWithManager,
		organizationpolicyuserattachment.SetupWebhookWithManager,
		organizationrole.SetupWebhookWithManager,
		organizationrolegroupattachment.SetupWebhookWithManager,
		organizationrolepolicyattachment.SetupWebhookWithManager,
		organizationrolesession.SetupWebhookWithManager,
		organizationroletrust.SetupWebhookWithManager,
		useriam.SetupWebhookWithManager,
		dnsdomain.SetupWebhookWithManager,
		dnsrecord.SetupWebhookWithManager,
		firewallgroup.SetupWebhookWithManager,
		firewallrule.SetupWebhookWithManager,
		loadbalancer.SetupWebhookWithManager,
		natgateway.SetupWebhookWithManager,
		natgatewayfirewallrule.SetupWebhookWithManager,
		natgatewayportforwardingrule.SetupWebhookWithManager,
		reservedip.SetupWebhookWithManager,
		reverseipv4.SetupWebhookWithManager,
		reverseipv6.SetupWebhookWithManager,
		vpc.SetupWebhookWithManager,
		providerconfig.SetupWebhookWithManager,
		blockstorage.SetupWebhookWithManager,
		containerregistry.SetupWebhookWithManager,
		objectstorage.SetupWebhookWithManager,
		objectstoragebucket.SetupWebhookWithManager,
		virtualfilesystemstorage.SetupWebhookWithManager,
		kubernetes.SetupWebhookWithManager,
		kubernetesnodepool.SetupWebhookWithManager,
	} {
		if err := setup(mgr); err != nil {
			return err
		}
	}
	return nil
}
