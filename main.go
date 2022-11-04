package main

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

var (
	Location string
)

func main() {

	pulumi.Run(NewDC)
}
