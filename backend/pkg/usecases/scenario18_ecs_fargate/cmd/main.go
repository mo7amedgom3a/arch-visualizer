package main

import (
	"context"
	"fmt"
	"os"

	"github.com/mo7amedgom3a/arch-visualizer/backend/pkg/usecases/scenario18_ecs_fargate"
)

func main() {
	ctx := context.Background()
	if err := scenario18_ecs_fargate.Run(ctx); err != nil {
		fmt.Printf("Error running scenario: %v\n", err)
		os.Exit(1)
	}
}
