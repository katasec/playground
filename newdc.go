package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/katasec/playground/azuredc"
	"github.com/katasec/playground/utils"
	network "github.com/pulumi/pulumi-azure-native/sdk/go/azure/network"
	"github.com/pulumi/pulumi-azure-native/sdk/go/azure/resources"
	computec "github.com/pulumi/pulumi-azure/sdk/v5/go/azure/compute"
	networkc "github.com/pulumi/pulumi-azure/sdk/v5/go/azure/network"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

var (
	launchVmFlag      = true
	launchBastionFlag = true
	launchK8sFlag     = true
)

// NewDC creates a new data centre based on a reference azuredc
func NewDC(ctx *pulumi.Context) error {

	// Create hub resource group and VNET
	hubrg, err := resources.NewResourceGroup(ctx, "rg-play-hub-", &resources.ResourceGroupArgs{})
	utils.ExitOnError(err)
	hubVnet, firewall := CreateHub(ctx, hubrg, &azuredc.ReferenceHubVNET)

	// Create some spokes
	rg1, vnet1 := AddSpoke(ctx, "nprod", hubrg, hubVnet, firewall, 0)
	// rg2, vnet2 := AddSpoke(ctx, "prod", hubrg, hubVnet, firewall, 1)
	AddSpoke(ctx, "prod", hubrg, hubVnet, firewall, 1)
	// AddSpoke(ctx, "nprod2", hubrg, hubVnet, firewall, 2)

	if launchBastionFlag {
		launchBastion(ctx, hubrg, hubVnet)
	}

	// Launch some vms
	if launchVmFlag {
		launchVM(ctx, vnet1, rg1, "snet-tier2-vm", "vm01") // <- nprod vm
		//launchVM(ctx, hubVnet, hubrg, "snet-test", "vm03") // <- hub vm
		//launchVM(ctx, vnet2, rg2, "snet-tier2-vm", "vm02") // <- prod vm
	}

	if launchK8sFlag {
		launchK8s(ctx, rg1, vnet1)
	}

	return err
}

// Creates an Azure Virtual Network and subnets using the provided VNETInfo
func CreateHub(ctx *pulumi.Context, rg *resources.ResourceGroup, vnetInfo *azuredc.VNETInfo) (*network.VirtualNetwork, *network.AzureFirewall) {

	// Generate list of subnets to create
	subnets := network.SubnetTypeArray{}
	for _, subnet := range vnetInfo.SubnetsInfo {
		subnets = append(subnets, network.SubnetTypeArgs{
			AddressPrefix: pulumi.String(subnet.AddressPrefix),
			Name:          pulumi.String(subnet.Name),
		})
	}

	// Create VNET + subnets
	vnet, err := network.NewVirtualNetwork(ctx, vnetInfo.Name, &network.VirtualNetworkArgs{
		AddressSpace: &network.AddressSpaceArgs{
			AddressPrefixes: pulumi.StringArray{
				pulumi.String(vnetInfo.AddressPrefix),
			},
		},
		ResourceGroupName: rg.Name,
		Subnets:           &subnets,
	})

	// Create Firewall
	firewall := createFirewall(ctx, rg, vnet)
	utils.ExitOnError(err)

	return vnet, firewall
}

func AddSpoke(ctx *pulumi.Context, spokeName string, hubrg *resources.ResourceGroup, hubVnet *network.VirtualNetwork, firewall *network.AzureFirewall, offset int) (*resources.ResourceGroup, *network.VirtualNetwork) {

	// Create a resource group
	rgName := fmt.Sprintf("%s-%s-", azuredc.RgPrefix, spokeName)
	rg, err := resources.NewResourceGroup(ctx, rgName, &resources.ResourceGroupArgs{})
	utils.ExitOnError(err)

	// Create a route to firewall
	routeName := fmt.Sprintf("rt-%s", spokeName)
	spokeRoute := createFWRoute(ctx, rg, routeName, firewall)

	// Generate CIDRs
	spokeCidrs := azuredc.NewSpokeVnetTemplate(spokeName, offset)

	// Create VNET using generated ResourceGroup, Route & CIDRs
	vnet := CreateVNET(ctx, rg, spokeCidrs, spokeRoute)

	// Peer hub with spoke
	pulumiUrn1 := fmt.Sprintf("hub-to-%s", spokeName)
	peerNetworks(ctx, pulumiUrn1, hubrg, hubVnet, vnet)

	// Peer spoke with hub
	pulumiUrn2 := fmt.Sprintf("%s-to-hub", spokeName)
	peerNetworks(ctx, pulumiUrn2, rg, vnet, hubVnet)

	return rg, vnet

}

func launchVM(ctx *pulumi.Context, vnet *network.VirtualNetwork, rg *resources.ResourceGroup, subnetName string, vmName string) {

	// Look up the firewall subnet
	subnet := network.LookupSubnetOutput(ctx, network.LookupSubnetOutputArgs{
		ResourceGroupName:  rg.Name,
		SubnetName:         pulumi.String(subnetName),
		VirtualNetworkName: vnet.Name,
	})

	nic, err := networkc.NewNetworkInterface(ctx, fmt.Sprintf("%s-nic", vmName), &networkc.NetworkInterfaceArgs{
		ResourceGroupName: rg.Name,
		IpConfigurations: networkc.NetworkInterfaceIpConfigurationArray{
			&networkc.NetworkInterfaceIpConfigurationArgs{
				Name:                       pulumi.String("internal"),
				SubnetId:                   subnet.Id(),
				PrivateIpAddressAllocation: pulumi.String("Dynamic"),
			},
		},
	})
	utils.ExitOnError(err)

	_, err = computec.NewLinuxVirtualMachine(ctx, vmName, &computec.LinuxVirtualMachineArgs{
		ResourceGroupName: rg.Name,
		Size:              pulumi.String("Standard_B2s"),
		AdminUsername:     pulumi.String("adminuser"),
		NetworkInterfaceIds: pulumi.StringArray{
			nic.ID(),
		},
		AdminSshKeys: computec.LinuxVirtualMachineAdminSshKeyArray{
			&computec.LinuxVirtualMachineAdminSshKeyArgs{
				Username: pulumi.String("adminuser"),
				//Export an env vairable here in your local system which contains your public key path
				PublicKey: readFileOrPanic(os.Getenv("SSH_ADMIN_KEY")),
				//PublicKey: os.Getenv("sshkey"),
			},
		},
		OsDisk: &computec.LinuxVirtualMachineOsDiskArgs{
			Caching:            pulumi.String("ReadWrite"),
			StorageAccountType: pulumi.String("Standard_LRS"),
		},
		SourceImageReference: &computec.LinuxVirtualMachineSourceImageReferenceArgs{
			Publisher: pulumi.String("Canonical"),
			Offer:     pulumi.String("UbuntuServer"),
			Sku:       pulumi.String("16.04-LTS"),
			Version:   pulumi.String("latest"),
		},
	})
	utils.ExitOnError(err)

}

func readFileOrPanic(path string) pulumi.StringInput {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err.Error())
	}
	return pulumi.String(string(data))
}

func launchBastion(ctx *pulumi.Context, rg *resources.ResourceGroup, vnet *network.VirtualNetwork) *computec.BastionHost {

	subnet := network.LookupSubnetOutput(ctx, network.LookupSubnetOutputArgs{
		ResourceGroupName:  rg.Name,
		SubnetName:         pulumi.String("AzureBastionSubnet"),
		VirtualNetworkName: vnet.Name,
	})

	subnetId := subnet.Id().ApplyT(func(subnetId *string) string { return *subnetId }).(pulumi.StringOutput)

	// Create an Management IP for the Basic firewall for Azure Service Traffic
	bastionIp, _ := network.NewPublicIPAddress(ctx, "basition-ip", &network.PublicIPAddressArgs{
		ResourceGroupName:        rg.Name,
		PublicIPAllocationMethod: pulumi.String("Static"),
		Sku: &network.PublicIPAddressSkuArgs{
			Name: pulumi.String("Standard"),
			Tier: pulumi.String("Regional"),
		},
	})

	bastion, err := computec.NewBastionHost(ctx, "bastion", &computec.BastionHostArgs{
		ResourceGroupName: rg.Name,
		IpConfiguration: computec.BastionHostIpConfigurationArgs{
			Name:              pulumi.String("ipconfig"),
			SubnetId:          subnetId,
			PublicIpAddressId: bastionIp.ID(),
		},
	})
	utils.ExitOnError(err)

	return bastion
}
