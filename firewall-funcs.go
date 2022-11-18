package main

import (
	"github.com/katasec/playground/utils"
	network "github.com/pulumi/pulumi-azure-native/sdk/go/azure/network"
	"github.com/pulumi/pulumi-azure-native/sdk/go/azure/resources"
	networkc "github.com/pulumi/pulumi-azure/sdk/v5/go/azure/network"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

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
func createFirewall(ctx *pulumi.Context, rg *resources.ResourceGroup, vnet *network.VirtualNetwork) *network.AzureFirewall {

	// Create Firewall Policy for later assignment
	fwPolicy := createFirewallPolicy(ctx, rg)

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
		FirewallPolicy: network.SubResourceArgs{
			Id: fwPolicy.ID(),
		},
	}, pulumi.DependsOn([]pulumi.Resource{vnet}))
	utils.ExitOnError(err)

	return firewall
}

func createFirewallPolicy(ctx *pulumi.Context, rg *resources.ResourceGroup) *network.FirewallPolicy {

	// Arguments to create a fwpolicy
	args := &network.FirewallPolicyArgs{
		ResourceGroupName: rg.Name,
		Sku: network.FirewallPolicySkuArgs{
			Tier: pulumi.String("Basic"),
		},
	}

	// Create fw policy
	firewallPolicy, err := network.NewFirewallPolicy(ctx, "hub-policy-", args)
	if err != nil {
		panic(err)
	}

	_, err = networkc.NewFirewallPolicyRuleCollectionGroup(ctx, "apprules-", &networkc.FirewallPolicyRuleCollectionGroupArgs{
		Priority:         pulumi.Int(101),
		FirewallPolicyId: firewallPolicy.ID(),
		ApplicationRuleCollections: networkc.FirewallPolicyRuleCollectionGroupApplicationRuleCollectionArray{
			networkc.FirewallPolicyRuleCollectionGroupApplicationRuleCollectionArgs{
				Name:     pulumi.String("Allowed_HTTPS"),
				Action:   pulumi.String("Allow"),
				Priority: pulumi.Int(100),
				Rules: networkc.FirewallPolicyRuleCollectionGroupApplicationRuleCollectionRuleArray{
					networkc.FirewallPolicyRuleCollectionGroupApplicationRuleCollectionRuleArgs{
						Name: pulumi.String("allow_ghcr.io"),
						SourceAddresses: pulumi.StringArray{
							pulumi.String("*"),
						},
						DestinationFqdns: pulumi.StringArray{
							pulumi.String("ghcr.io"),
						},
						Protocols: networkc.FirewallPolicyRuleCollectionGroupApplicationRuleCollectionRuleProtocolArray{
							&networkc.FirewallPolicyRuleCollectionGroupApplicationRuleCollectionRuleProtocolArgs{
								Port: pulumi.Int(443),
								Type: pulumi.String("Https"),
							},
						},
					},
				},
			},
		},
	})
	utils.ExitOnError(err)

	return firewallPolicy
}
