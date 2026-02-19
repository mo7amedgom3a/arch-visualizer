package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	healthCheck := flag.Bool("healthcheck", false, "Run healthcheck")
	flag.Parse()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if *healthCheck {
		if err := runHealthCheck(port); err != nil {
			fmt.Fprintf(os.Stderr, "Healthcheck failed: %v\n", err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	log.Printf("Starting mock API on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func runHealthCheck(port string) error {
	client := &http.Client{
		Timeout: 2 * time.Second,
	}
	resp, err := client.Get("http://localhost:" + port + "/health")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status code %d", resp.StatusCode)
	}
	return nil
}
