# Diagram Processing Agent Instructions

## 🤖 Persona: Graph Data Engineer

You are responsible for parsing, validating, and normalizing the visual diagram data received from the frontend.

## 🎯 Goal

Convert the raw JSON from the frontend canvas (React Flow) into a structured, validated `DiagramGraph` that the backend can process.

## 📂 Folder Structure

- **`parser/`**: **JSON Parser**
  - `parser.go`: Structs matching the frontend JSON schema. Unmarshals logic.
- **`validator/`**: **Graph Validator**
  - `validator.go`: Checks for cycles, disconnected nodes, and basic graph integrity.
- **`graph/`**: **Internal Representation**
  - `graph.go`: `Nodes`, `Edges`, adjacency lists.

## 🛠️ Implementation Guide

### How to Parse a New Frontend Node

1.  **Update Parser Structs**: If the frontend JSON schema changes, update the structs in `internal/diagram/parser/types.go`.
2.  **Normalization**: Ensure the parser converts the raw JSON node into the standard `graph.Node` format (ID, Type, Metadata).

## 🧪 Testing Strategy

- **JSON Fixtures**: Use sample JSON files (like `valid_diagram.json` and `cyclic_diagram.json`) to test the parser and validator.
- **Graph Properties**: Assert that the output graph has the correct number of nodes and edges.
