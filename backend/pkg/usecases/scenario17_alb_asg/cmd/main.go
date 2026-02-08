package main

import (
	"context"
	"fmt"
	"os"

	"github.com/mo7amedgom3a/arch-visualizer/backend/pkg/usecases/scenario17_alb_asg"
)

func main() {
	ctx := context.Background()
	if err := scenario17_alb_asg.Run(ctx); err != nil {
		fmt.Printf("Error running scenario: %v\n", err)
		os.Exit(1)
	}
}
