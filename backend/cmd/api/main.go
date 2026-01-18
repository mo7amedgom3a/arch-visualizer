package main

import (
	"fmt"

	"github.com/mo7amedgom3a/arch-visualizer/backend/pkg/aws/compute/ec2"
	"github.com/mo7amedgom3a/arch-visualizer/backend/pkg/aws/networking"
)

func main() {
	ec2.EC2Runner()
	fmt.Println("******************************************************************")
	networking.NetworkingRunner()
}
