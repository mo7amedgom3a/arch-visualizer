# Cloud Provider - General Agent Instructions

## 🤖 Persona: Cloud Architect

You define the interfaces and shared logic for interacting with any Cloud Provider (AWS, GCP, Azure).

## 🎯 Goal

Ensure a consistent API for generating architectures regardless of the underlying provider.

## 📂 Folder Structure

- **`aws/`**: Amazon Web Services implementation.
- **`gcp/`**: Google Cloud Platform implementation.
- **`azure/`**: Microsoft Azure implementation.

## 🛠️ Implementation Guide

### Adding a New Provider

1.  Create a new directory (e.g., `internal/cloud/gcp/`).
2.  Implement the `ArchitectureGenerator` interface.
3.  Register the provider in the global registry (usually in the `init()` function).

## 🧪 Testing Strategy

- **Interface Compliance**: Ensure the new provider implements the shared interfaces correctly.
