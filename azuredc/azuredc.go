package azuredc

import (
	"fmt"
	"strconv"
	"strings"
)

// VNETInfo contains details on the VNET we want to create.
// For e.g. Address Prefix, Subnets etc.
type VNETInfo struct {
	Name          string
	AddressPrefix string
	SubnetsInfo   []SubnetInfo
}

// SubnetInfo contains details on the SubnetInfo we want to create.
// For e.g. Name, Cidr etc.
type SubnetInfo struct {
	Name          string
	AddressPrefix string
	Description   string
	Tags          map[string]string
}

func (v *VNETInfo) GenerateSpoke(offset int) {

	// Add offset to second octet to generate 2nd octet number
	//octetStartStr := strconv.Itoa(octetStart)
	newOctetStr := strconv.Itoa(octetStart + offset)

	// Generate VNET Address Prefix
	v.AddressPrefix = strings.Replace(ReferenceSpokeVNET.AddressPrefix, "x", newOctetStr, 1)

	for _, subnet := range ReferenceSpokeVNET.SubnetsInfo {
		nSubnet := SubnetInfo{
			Name:          subnet.Name,
			Description:   subnet.Description,
			AddressPrefix: strings.Replace(subnet.AddressPrefix, "x", newOctetStr, 1),
		}
		v.SubnetsInfo = append(v.SubnetsInfo, nSubnet)
	}
}

// Returns an address prefix and subnet CIDRs used to create a vnet. The second octet of the
// generated CIDRs are "0" by default and can changed via the "offset" parameter
func NewSpokeVnetTemplate(name string, offset ...int) *VNETInfo {
	vnet := &VNETInfo{}

	if len(offset) == 0 {
		vnet.GenerateSpoke(0)
	} else {
		vnet.GenerateSpoke(offset[0])
	}

	vnet.Name = name

	return vnet
}

func (v *VNETInfo) Dump() {

	fmt.Println("VNet:\n" + v.AddressPrefix)

	fmt.Println("Subnets:")

	for _, subnet := range v.SubnetsInfo {
		fmt.Println(subnet.Name + ":" + subnet.AddressPrefix)
	}
}
