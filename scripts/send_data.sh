#!/bin/bash -x

# Check if file argument and OS type are provided
if [ -z "$1" ] || [ -z "$2" ]; then
    echo "Usage: $0 <filename> <os_type>"
    echo "Example: $0 testdata/windows_tasklist.txt windows"
    echo "         $0 testdata/linux_ps.txt linux"
    exit 1
fi

FILE="$1"
OS_TYPE="$2"

# Set fields based on OS type
if [ "$OS_TYPE" = "windows" ]; then
    MACHINE_ID="win-pc-001"
    MACHINE_NAME="Windows Workstation"
    OS_VERSION="Windows"
    COMMAND_TYPE="tasklist"
    USER_NAME="jane.smith"
    RAND_STR=$(date +%s%N | md5 | head -c 8)
    USER_ID="$RAND_STR@med.edu"
    FACULTY="Medicine"
elif [ "$OS_TYPE" = "linux" ]; then
    MACHINE_ID="linux-pc-001"
    MACHINE_NAME="Linux Workstation"
    OS_VERSION="Ubuntu"
    COMMAND_TYPE="ps"
    USER_NAME="john.doe"
    RAND_STR=$(date +%s%N | md5 | head -c 8)
    USER_ID="$RAND_STR@law.edu"
    FACULTY="Law"
else
    echo "Invalid OS type: $OS_TYPE. Must be 'windows' or 'linux'."
    exit 1
fi

# Read the file content
FILE_CONTENT=$(cat "$FILE" | jq -Rs .)

# Create the JSON payload
JSON_PAYLOAD=$(cat <<EOF
{
  "machine_id": "$MACHINE_ID",
  "machine_name": "$MACHINE_NAME",
  "os_version": "$OS_VERSION",
  "timestamp": "$(date -u +"%Y-%m-%dT%H:%M:%SZ")",
  "command_type": "$COMMAND_TYPE",
  "user_name": "$USER_NAME",
  "user_id": "$USER_ID",
  "faculty": "$FACULTY",
  "command_output": $FILE_CONTENT
}
EOF
)

# Send the request
curl -X POST http://localhost:8080/upload \
-H "Content-Type: application/json" \
-d "$JSON_PAYLOAD" -v