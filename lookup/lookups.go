package lookup

import (
	"fmt"
	"strconv"
	"strings"
)

type VnetTemplate struct {
	AddressPrefix string
	Subnets       []SubnetTemplate
}

type SubnetTemplate struct {
	Name          string
	AddressPrefix string
	Description   string
}

func (v *VnetTemplate) Generate(offset int) {

	// Add offset to second octet to generate 2nd octet number
	//octetStartStr := strconv.Itoa(octetStart)
	newOctetStr := strconv.Itoa(octetStart + offset)

	// Generate VNET Address Prefix
	v.AddressPrefix = strings.Replace(ReferenceVNET.AddressPrefix, "x", newOctetStr, 1)

	for _, subnet := range ReferenceVNET.Subnets {
		nSubnet := SubnetTemplate{
			Name:          subnet.Name,
			Description:   subnet.Description,
			AddressPrefix: strings.Replace(subnet.AddressPrefix, "x", newOctetStr, 1),
		}
		v.Subnets = append(v.Subnets, nSubnet)
	}
}

func NewVnetTemplate(offset ...int) *VnetTemplate {
	vnet := &VnetTemplate{}

	if len(offset) == 0 {
		vnet.Generate(0)
	} else {
		vnet.Generate(offset[0])
	}

	return vnet
}

func (v *VnetTemplate) Dump() {

	fmt.Println("VNet:\n" + v.AddressPrefix)

	fmt.Println("Subnets:")

	for _, subnet := range v.Subnets {
		fmt.Println(subnet.Name + ":" + subnet.AddressPrefix)
	}
}
