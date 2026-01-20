package main

import (
	// "fmt"

	// "github.com/mo7amedgom3a/arch-visualizer/backend/pkg/aws/compute/ec2"
	// "github.com/mo7amedgom3a/arch-visualizer/backend/pkg/aws/networking"
	// "github.com/mo7amedgom3a/arch-visualizer/backend/pkg/aws/compute/alb"
	// "github.com/mo7amedgom3a/arch-visualizer/backend/pkg/aws/compute/autoscaling"
	"github.com/mo7amedgom3a/arch-visualizer/backend/pkg/usecases/scenario1_basic_web_app"
)

func main() {
	// ec2.EC2Runner()
	// fmt.Println("******************************************************************")
	// networking.NetworkingRunner()
	// alb.ALBRunner()
	// autoscaling.ASGRunner()
	scenario1_basic_web_app.BasicWebAppRunner()
}
