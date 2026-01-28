#!/bin/bash

# Script to test the diagram processing API endpoint
# This simulates an API request with the diagram JSON

API_URL="http://localhost:8080/api/diagrams/process"
DIAGRAM_JSON_FILE="json-request-diagram-valid.json"

# Check if the diagram JSON file exists
if [ ! -f "$DIAGRAM_JSON_FILE" ]; then
    echo "Error: $DIAGRAM_JSON_FILE not found"
    exit 1
fi

# Make the API request
echo "Sending diagram processing request to $API_URL"
echo "Using diagram file: $DIAGRAM_JSON_FILE"
echo ""

curl -X POST "$API_URL?project_name=Test%20Project&iac_tool_id=1&user_id=00000000-0000-0000-0000-000000000001" \
  -H "Content-Type: application/json" \
  -d @"$DIAGRAM_JSON_FILE" \
  -w "\n\nHTTP Status: %{http_code}\n" \
  -v

echo ""
echo "Done!"
