#!/bin/bash

# Database connection details
DB_USER="dreamuser"
DB_NAME="dreamdb"
DB_HOST="localhost"
DB_PORT="5432"

echo "Checking database tables..."

# List all tables
echo -e "\nTables in database:"
docker exec dream-postgres-1 psql -U $DB_USER -d $DB_NAME -c "\dt"

# Count records in each table
echo -e "\nRecord counts:"
docker exec dream-postgres-1 psql -U $DB_USER -d $DB_NAME -c "
SELECT 'users' as table_name, COUNT(*) as count FROM users
UNION ALL
SELECT 'linux_processes' as table_name, COUNT(*) as count FROM linux_process
UNION ALL
SELECT 'windows_processes' as table_name, COUNT(*) as count FROM windows_process;
"

# Show sample data from each table
echo -e "\nSample data from users table:"
docker exec dream-postgres-1 psql -U $DB_USER -d $DB_NAME -c "SELECT * FROM users LIMIT 5;"

echo -e "\nSample data from linux_processes table:"
docker exec dream-postgres-1 psql -U $DB_USER -d $DB_NAME -c "SELECT * FROM linux_processes LIMIT 5;"

echo -e "\nSample data from windows_process table:"
docker exec dream-postgres-1 psql -U $DB_USER -d $DB_NAME -c "SELECT * FROM windows_processes LIMIT 5;" 