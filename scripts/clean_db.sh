#!/bin/bash

# Database connection details
DB_USER="dreamuser"
DB_NAME="dreamdb"
DB_HOST="localhost"
DB_PORT="5432"

echo "Cleaning database tables..."

# Connect to PostgreSQL and truncate tables
docker exec dream-postgres-1 psql -U $DB_USER -d $DB_NAME -c "TRUNCATE users, linux_processes, windows_processes CASCADE;"

if [ $? -eq 0 ]; then
    echo "Database tables cleaned successfully"
else
    echo "Error cleaning database tables"
    exit 1
fi 