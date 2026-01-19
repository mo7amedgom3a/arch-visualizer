package sdk

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
)

// DecodePolicyDocument decodes a URL-encoded policy document string
func DecodePolicyDocument(encoded string) (string, error) {
	if encoded == "" {
		return "", nil
	}

	// Try to decode URL-encoded string
	decoded, err := url.QueryUnescape(encoded)
	if err != nil {
		// If decoding fails, return original (might already be decoded)
		return encoded, fmt.Errorf("failed to decode policy document: %w", err)
	}

	// Validate that the decoded string is valid JSON
	var jsonDoc interface{}
	if err := json.Unmarshal([]byte(decoded), &jsonDoc); err != nil {
		// If it's not valid JSON, it might already be decoded, return as-is
		return encoded, fmt.Errorf("decoded string is not valid JSON: %w", err)
	}

	return decoded, nil
}

// DecodeAllPolicyFiles processes all service-specific policy JSON files and decodes their policy_document fields
func DecodeAllPolicyFiles() error {
	// List available services
	services, err := ListAvailableServices()
	if err != nil {
		return fmt.Errorf("failed to list services: %w", err)
	}

	fmt.Printf("Found %d services to process\n", len(services))

	// Process each service
	for _, service := range services {
		fmt.Printf("\nProcessing service: %s\n", service)
		if err := decodeServicePolicies(service); err != nil {
			fmt.Printf("Warning: Failed to decode policies for service %s: %v\n", service, err)
			continue
		}
		fmt.Printf("✓ Successfully decoded policies for service: %s\n", service)
	}

	return nil
}

// decodeServicePolicies decodes policy documents in a specific service's JSON file
func decodeServicePolicies(service string) error {
	// Load policies from service file
	entries, err := LoadServicePolicies(service)
	if err != nil {
		return fmt.Errorf("failed to load service policies: %w", err)
	}

	fmt.Printf("  Found %d policies\n", len(entries))

	// Decode each policy document
	decodedCount := 0
	for i := range entries {
		if entries[i].PolicyDocument != "" {
			decoded, err := DecodePolicyDocument(entries[i].PolicyDocument)
			if err == nil && decoded != entries[i].PolicyDocument {
				// Only update if decoding was successful and changed the value
				entries[i].PolicyDocument = decoded
				decodedCount++
			}
		}
	}

	fmt.Printf("  Decoded %d policy documents\n", decodedCount)

	// Save back to file
	dataDir := getDataDirectory()
	serviceDir := filepath.Join(dataDir, service)
	jsonPath := filepath.Join(serviceDir, "policies.json")

	jsonData, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal decoded policies: %w", err)
	}

	if err := os.WriteFile(jsonPath, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write decoded policies: %w", err)
	}

	return nil
}
