package main

import (
	// "fmt"

// "github.com/mo7amedgom3a/arch-visualizer/backend/pkg/aws/compute/ec2"
// "github.com/mo7amedgom3a/arch-visualizer/backend/pkg/aws/networking"
// "github.com/mo7amedgom3a/arch-visualizer/backend/pkg/aws/compute/alb"
// "github.com/mo7amedgom3a/arch-visualizer/backend/pkg/aws/compute/autoscaling"
// "github.com/mo7amedgom3a/arch-visualizer/backend/pkg/aws/storage/s3"
"github.com/mo7amedgom3a/arch-visualizer/backend/pkg/aws/iam"
// "github.com/mo7amedgom3a/arch-visualizer/backend/pkg/usecases/scenario1_basic_web_app"
// "github.com/mo7amedgom3a/arch-visualizer/backend/pkg/usecases/scenario2_high_availability"

)
func main() {	
	// ec2.EC2Runner()
	// fmt.Println("******************************************************************")
	// networking.NetworkingRunner()
	// alb.ALBRunner()
	// autoscaling.ASGRunner()
	// s3.S3Runner()

	// To fetch and save S3 policies, uncomment the following:
	iam.FetchAndSaveS3Policies()
}
