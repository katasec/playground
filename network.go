package main

import (
	"github.com/katasec/playground/azuredc"
	"github.com/katasec/playground/utils"
	containerinstance "github.com/pulumi/pulumi-azure-native/sdk/go/azure/containerinstance"
	network "github.com/pulumi/pulumi-azure-native/sdk/go/azure/network"
	"github.com/pulumi/pulumi-azure-native/sdk/go/azure/resources"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	// github.com/pulumi/pulumi-azure/sdk/v4/go/azure
	// github.com/pulumi/pulumi-azure/sdk/v4/go/azure
)

// NewDC creates a new data centre based on a reference azuredc
func NewDC(ctx *pulumi.Context) error {

	// Create hub resource group and VNET
	hubrg, err := resources.NewResourceGroup(ctx, "rg-play-hub-", &resources.ResourceGroupArgs{})
	utils.ExitOnError(err)
	hubVnet, firewall := CreateHub(ctx, hubrg, &azuredc.ReferenceHubVNET)

	// Create Firewall in Hub
	//firewall := createFirewall(ctx, hubrg, hubVnet)

	// Create nprod resource group and VNET
	nprodrg, err := resources.NewResourceGroup(ctx, "rg-play-nprod-", &resources.ResourceGroupArgs{})
	utils.ExitOnError(err)

	// Create nprod route to firewall
	nprdRoute := createFWRoute(ctx, nprodrg, "rt-nprod", firewall)

	// Create Spoke VNET with nprod route
	nprodCidrs := azuredc.NewSpokeVnetTemplate("nprod")
	nprodVnet := CreateVNET(ctx, nprodrg, nprodCidrs, nprdRoute)

	// Create prod resource group and VNET
	prodrg, err := resources.NewResourceGroup(ctx, "rg-play-prod-", &resources.ResourceGroupArgs{})
	utils.ExitOnError(err)

	// Create prod route to firewall
	prdRoute := createFWRoute(ctx, prodrg, "rt-prod", firewall)

	// Create Spoke VNET with prod route
	prodCidrs := azuredc.NewSpokeVnetTemplate("prod", 1)
	prodVnet := CreateVNET(ctx, prodrg, prodCidrs, prdRoute)

	// Peer hub to nprod
	peerNetworks(ctx, "hub-to-nprod", hubrg, hubVnet, nprodVnet)
	peerNetworks(ctx, "nprod-to-hub", nprodrg, nprodVnet, hubVnet)

	// Peer hub to prod
	peerNetworks(ctx, "hub-to-prod", hubrg, hubVnet, prodVnet)
	peerNetworks(ctx, "prod-to-hub", prodrg, prodVnet, hubVnet)

	// Spin up a container instance
	//StartContainers(ctx, hubrg, hubVnet)
	return err
}

func peerNetworks(ctx *pulumi.Context, urn string, srcRg *resources.ResourceGroup, src *network.VirtualNetwork, dst *network.VirtualNetwork) {
	name := pulumi.Sprintf("%s-to-%s", src.Name, dst.Name)
	_, err := network.NewVirtualNetworkPeering(ctx, urn, &network.VirtualNetworkPeeringArgs{
		Name:                      name,
		VirtualNetworkPeeringName: name,
		ResourceGroupName:         srcRg.Name,
		VirtualNetworkName:        src.Name,

		AllowForwardedTraffic:     pulumi.Bool(true),
		AllowGatewayTransit:       pulumi.Bool(false),
		AllowVirtualNetworkAccess: pulumi.Bool(true),
		RemoteVirtualNetwork: &network.SubResourceArgs{
			Id: dst.ID(),
		},
		UseRemoteGateways: pulumi.Bool(false),
	})
	utils.ExitOnError(err)
}

func createFirewall(ctx *pulumi.Context, rg *resources.ResourceGroup, vnet *network.VirtualNetwork) *network.AzureFirewall {

	// Create an Management IP for the Basic firewall for Azure Service Traffic
	managementIp, _ := network.NewPublicIPAddress(ctx, "fw-mgmt-ip", &network.PublicIPAddressArgs{
		ResourceGroupName:        rg.Name,
		PublicIPAllocationMethod: pulumi.String("Static"),
		Sku: &network.PublicIPAddressSkuArgs{
			Name: pulumi.String("Standard"),
			Tier: pulumi.String("Regional"),
		},
	})

	// Create a public IP for Firewall for inbound/outbound traffic
	publicIp, _ := network.NewPublicIPAddress(ctx, "fwip", &network.PublicIPAddressArgs{
		ResourceGroupName:        rg.Name,
		PublicIPAllocationMethod: pulumi.String("Static"),
		Sku: &network.PublicIPAddressSkuArgs{
			Name: pulumi.String("Standard"),
			Tier: pulumi.String("Regional"),
		},
	})

	// Look up the firewall subnet
	fwSubnet := network.LookupSubnetOutput(ctx, network.LookupSubnetOutputArgs{
		ResourceGroupName:  rg.Name,
		SubnetName:         pulumi.String("AzureFirewallSubnet"),
		VirtualNetworkName: vnet.Name,
	})

	// Look up the mgmt subnet subnet
	mgmtfwSubnet := network.LookupSubnetOutput(ctx, network.LookupSubnetOutputArgs{
		ResourceGroupName:  rg.Name,
		SubnetName:         pulumi.String("AzureFirewallManagementSubnet"),
		VirtualNetworkName: vnet.Name,
	})

	// Create a firewall
	firewall, err := network.NewAzureFirewall(ctx, "hubfirewall", &network.AzureFirewallArgs{
		ResourceGroupName: rg.Name,
		Sku: &network.AzureFirewallSkuArgs{
			Name: pulumi.String("AZFW_VNet"),
			Tier: pulumi.String("Basic"),
		},
		IpConfigurations: &network.AzureFirewallIPConfigurationArray{
			network.AzureFirewallIPConfigurationArgs{
				Name: pulumi.String("configuration"),
				PublicIPAddress: network.SubResourceArgs{
					Id: publicIp.ID(),
				},
				Subnet: network.SubResourceArgs{
					Id: fwSubnet.Id(),
				},
			},
		},
		ManagementIpConfiguration: &network.AzureFirewallIPConfigurationArgs{
			Name: pulumi.String("mgmt-configuration"),
			PublicIPAddress: network.SubResourceArgs{
				Id: managementIp.ID(),
			},
			Subnet: network.SubResourceArgs{
				Id: mgmtfwSubnet.Id(),
			},
		},
	}, pulumi.DependsOn([]pulumi.Resource{vnet}))
	utils.ExitOnError(err)

	return firewall
}

func createFWRoute(ctx *pulumi.Context, rg *resources.ResourceGroup, tableName string, firewall *network.AzureFirewall) *network.RouteTable {

	// Create Table
	routeTable, err := network.NewRouteTable(ctx, tableName, &network.RouteTableArgs{
		ResourceGroupName: rg.Name,
		RouteTableName:    pulumi.String(tableName),
	})
	utils.ExitOnError(err)

	// Create route to firewall
	_, err = network.NewRoute(ctx, tableName+"-firewall-route", &network.RouteArgs{
		AddressPrefix:     pulumi.String("0.0.0.0/0"),
		NextHopType:       pulumi.String("VirtualAppliance"),
		ResourceGroupName: rg.Name,
		RouteName:         pulumi.String("firewall-route"),
		RouteTableName:    routeTable.Name,
		NextHopIpAddress:  firewall.IpConfigurations.Index(pulumi.Int(0)).PrivateIPAddress(),
	})
	utils.ExitOnError(err)

	return routeTable
}

func StartContainers(ctx *pulumi.Context, rg *resources.ResourceGroup, hub *network.VirtualNetwork) {
	// containerinstance.NewContainerGroup(ctx, "bash", &containerinstance.ContainerGroupArgs{
	// 	ContainerGroupName: pulumi.String("bash"),
	// 	ResourceGroupName:  rg.Name,
	// 	Containers: containerinstance.ContainerArray{
	// 		containerinstance.ContainerArgs{
	// 			Image: pulumi.String("registry.hub.docker.com/library/bash"),
	// 			Command: pulumi.StringArray{
	// 				pulumi.String("/usr/local/bin/bash"),
	// 				pulumi.String("-c"),
	// 				pulumi.String("echo hello; sleep 100000"),
	// 			},
	// 		},
	// 	},
	// })

	_, err := containerinstance.NewContainerGroup(ctx, "bash", &containerinstance.ContainerGroupArgs{
		OsType:             pulumi.String("Linux"),
		ContainerGroupName: pulumi.String("bash"),
		ResourceGroupName:  rg.Name,
		Containers: containerinstance.ContainerArray{
			containerinstance.ContainerArgs{
				Name: pulumi.String("hello-world"),
				Resources: containerinstance.ResourceRequirementsArgs{
					Limits: containerinstance.ResourceLimitsArgs{
						Cpu:        pulumi.Float64(0.5),
						MemoryInGB: pulumi.Float64(1.5),
					},
					Requests: containerinstance.ResourceRequestsArgs{
						Cpu:        pulumi.Float64(0.5),
						MemoryInGB: pulumi.Float64(1.5),
					},
				},
				Image: pulumi.String("registry.hub.docker.com/library/bash"),
				Command: pulumi.StringArray{
					pulumi.String("/usr/local/bin/bash"),
					pulumi.String("-c"),
					pulumi.String("echo hello; sleep 100000"),
				},
			},
		},
		NetworkProfile: containerinstance.ContainerGroupNetworkProfileArgs{
			Id: pulumi.String("/subscriptions/174c6cc1-faef-4e40-91f4-1bef3a703153/resourceGroups/rg-play-hub-cfe74535/providers/Microsoft.Network/virtualNetworks/hub60b0830a/subnets/snet-hub/aci"),
		},
	})
	utils.ExitOnError(err)

}
