# CashLenX CLI Reference

**Version**: 2.0.0
**Last Updated**: 2025-01-26

**See your money clearly**

Command-line interface for managing personal finances with CashLenX.

## Overview

The CashLenX CLI provides a complete interface to manage your finances from the command line, organized by access level:

- **`cashlenx open`** - Public commands (no authentication required)
- **`cashlenx admin`** - Admin commands (requires admin privileges)
- **`cashlenx cash`** - Cash flow management (requires authentication)
- **`cashlenx category`** - Category management (requires authentication)

## Quick Start

```bash
# Start the server
cashlenx open start -p 8080

# Check server health
cashlenx open health

# Add expense
cashlenx cash expense -c "Food" -a 45.50 -d "Lunch"

# Add income
cashlenx cash income -c "Salary" -a 5000

# View today's transactions
cashlenx cash query -b $(date +%Y-%m-%d)

# Create backup
cashlenx admin backup -o backup.json
```

## Installation

```bash
cd backend
go build -o cashlenx main.go
sudo mv cashlenx /usr/local/bin/
```

## Environment Setup

```bash
# MongoDB
export DB_TYPE=mongodb
export MONGO_DB_URI="mongodb://cashlenx:cashlenx123@localhost:27017/cashlenx?authSource=admin"
export DB_NAME=cashlenx

# MySQL
export DB_TYPE=mysql
export MYSQL_DB_URI="cashlenx:cashlenx123@tcp(localhost:3306)/cashlenx"
export DB_NAME=cashlenx

# JWT Secret (required for authentication)
export JWT_SECRET="your-secret-key"
```

---

## Command Structure

```
cashlenx
├── open                # Public commands (no auth)
│   ├── health         Check system health
│   ├── version        Show version info
│   └── start          Start API server
├── admin              # Admin commands (admin role required)
│   ├── backup         Create database backup
│   ├── restore        Restore from backup
│   ├── export         Export to Excel (TODO: move to statistic)
│   └── import         Import from Excel (TODO: move to statistic)
├── cash               # Cash flow commands (user auth required)
│   ├── income         Add income
│   ├── expense        Add expense
│   ├── update         Update transaction
│   ├── delete         Delete transaction
│   ├── query          Query transactions
│   ├── list           List all transactions
│   ├── range          Query date range
│   └── summary        Show summary
└── category           # Category commands (user auth required)
    ├── create         Create category
    ├── update         Update category
    ├── delete         Delete category
    ├── query          Query categories
    ├── list           List all categories
    └── tree           Show category tree
```

---

## Public Commands (`cashlenx open`)

### open health
Check if the API server is running and healthy.

```bash
cashlenx open health
```

**Output**:
```
✅ System is healthy
Status Code: 200 OK
```

**Requirements**: Server must be running at localhost:8080

---

### open version
Show version information.

```bash
cashlenx open version
```

**Output**:
```
CashLenX v2.0.0
Build Time: 2024-01-15T10:00:00Z
Git Commit: abc1234
Go Version: go1.21.5
OS/Arch: linux/amd64
```

---

### open start
Start the API server.

```bash
cashlenx open start -p 8080
```

**Flags**:
- `-p, --port` - Server port (default: 8080)

**Environment variables required**:
- `MONGO_DB_URI` or `MYSQL_DB_URI` - Database connection
- `DB_TYPE` - Database type (mongodb/mysql)
- `DB_NAME` - Database name
- `JWT_SECRET` - Secret for JWT tokens

---

## Admin Commands (`cashlenx admin`)

**Note**: All admin commands require admin privileges and the ADMIN_TOKEN environment variable.

### admin backup
Create a backup of all database data (all users).

```bash
# Auto-generated filename
cashlenx admin backup

# Custom filename
cashlenx admin backup -o backup_20240115.json
```

**Flags**:
- `-o, --output` - Backup file path (optional, default: cashlenx_backup_TIMESTAMP.json)
- `-t, --admin-token` - Admin token for dangerous operations

**Output Statistics**:
- Users: success/failed counts
- Categories: success/failed counts
- Cash Flows: success/failed counts

**Environment**:
- `ADMIN_TOKEN` - Required for verification

---

### admin restore
Restore database from a backup file.

```bash
cashlenx admin restore -i backup_20240115.json

# Skip confirmation
cashlenx admin restore -i backup_20240115.json -f
```

**Flags**:
- `-i, --input` - Backup file path (required)
- `-f, --force` - Skip confirmation prompt
- `-t, --admin-token` - Admin token for dangerous operations

**Output Statistics**:
- Users: success/failed counts
- Categories: success/failed counts
- Cash Flows: success/failed counts

⚠️ **WARNING**: This replaces all existing data!

---

### admin export
Export data to Excel.

```bash
# Export all data
cashlenx admin export -o data.xlsx

# Export date range
cashlenx admin export -f 2024-01-01 -t 2024-01-31 -o january.xlsx
```

**Flags**:
- `-o, --output` - Output file path (default: ./export.xlsx)
- `-f, --from` - Start date (optional, YYYY-MM-DD)
- `-t, --to` - End date (optional, YYYY-MM-DD)

**TODO**: This command will move to `cashlenx statistic export` with user data isolation.

---

### admin import
Import data from Excel.

```bash
cashlenx admin import -i data.xlsx
```

**Flags**:
- `-i, --input` - Input file path (required)

**TODO**: This command will move to `cashlenx statistic import` with user data isolation.

---

## Cash Flow Commands (`cashlenx cash`)

**Note**: All cash flow commands enforce user data isolation - users can only access their own transactions.

### cash income
Add new income transaction.

```bash
cashlenx cash income -c "Salary" -a 5000 -d "Monthly salary"
cashlenx cash income -c "Freelance" -a 1500 -b 2024-01-15
```

**Flags**:
- `-c, --category` - Category name (required)
- `-a, --amount` - Amount (required)
- `-b, --date` - Transaction date (optional, default: today)
- `-d, --description` - Description (optional)

---

### cash expense
Add new expense transaction.

```bash
cashlenx cash expense -c "Food & Dining" -a 45.50 -d "Lunch"
cashlenx cash expense -c "Transportation" -a 20 -b 2024-01-15
```

**Flags**:
- `-c, --category` - Category name (required)
- `-a, --amount` - Amount (required)
- `-b, --date` - Transaction date (optional, default: today)
- `-d, --description` - Description (optional)

---

### cash update
Update existing transaction.

```bash
cashlenx cash update -i 507f1f77bcf86cd799439011 -a 50.00
cashlenx cash update -i 507f1f77bcf86cd799439011 -c "Groceries" -d "Updated"
```

**Flags**:
- `-i, --id` - Transaction ID (required)
- `-a, --amount` - New amount (optional)
- `-c, --category` - New category (optional)
- `-b, --date` - New date (optional)
- `-d, --description` - New description (optional)

**Note**: Can only update your own transactions.

---

### cash delete
Delete transaction(s).

```bash
# Delete by ID
cashlenx cash delete -i 507f1f77bcf86cd799439011

# Delete all transactions on a date
cashlenx cash delete -b 2024-01-15
```

**Flags**:
- `-i, --id` - Transaction ID
- `-b, --date` - Date (YYYY-MM-DD)

**Note**: Can only delete your own transactions.

---

### cash query
Query transactions by filters.

```bash
# Query by ID
cashlenx cash query -i 507f1f77bcf86cd799439011

# Query by date
cashlenx cash query -b 2024-01-15

# Query by exact description
cashlenx cash query -e "Lunch"

# Query by fuzzy description
cashlenx cash query -f "lunch"
```

**Flags**:
- `-i, --id` - Query by ID
- `-b, --date` - Query by date
- `-e, --exact` - Query by exact description
- `-f, --fuzzy` - Query by fuzzy description

**Note**: Only returns your own transactions.

---

### cash list
List all transactions with pagination.

```bash
# List all
cashlenx cash list

# List with limit
cashlenx cash list -l 20

# List with offset
cashlenx cash list -l 20 -o 40

# Filter by type
cashlenx cash list -t income
cashlenx cash list -t expense
```

**Flags**:
- `-l, --limit` - Maximum records to return (default: 50)
- `-o, --offset` - Number of records to skip (default: 0)
- `-t, --type` - Filter by type (income/expense)

**Note**: Only returns your own transactions.

---

### cash range
Query transactions by date range.

```bash
cashlenx cash range -f 2024-01-01 -t 2024-01-31
cashlenx cash range --from 2024-01-01 --to 2024-01-31
```

**Flags**:
- `-f, --from` - Start date (YYYY-MM-DD) (required)
- `-t, --to` - End date (YYYY-MM-DD) (required)

**Output includes**:
- All transactions in range
- Total income
- Total expense
- Balance

**Note**: Only returns your own transactions.

---

### cash summary
Show financial summary.

```bash
# Monthly summary
cashlenx cash summary -p monthly -d 2024-01

# Daily summary (if implemented)
cashlenx cash summary -p daily -d 2024-01-15

# Yearly summary (if implemented)
cashlenx cash summary -p yearly -d 2024
```

**Flags**:
- `-p, --period` - Period type (daily/monthly/yearly) (required)
- `-d, --date` - Date for summary (required)
  - Daily: YYYY-MM-DD
  - Monthly: YYYY-MM
  - Yearly: YYYY

**Output includes**:
- Total income
- Total expense
- Balance
- Transaction count

**Note**: Only includes your own transactions.

---

## Category Commands (`cashlenx category`)

**Note**: All category commands enforce user data isolation - users can only access their own categories.

### category create
Create new category.

```bash
cashlenx category create -n "Food & Dining"
cashlenx category create -n "Groceries" -p 507f1f77bcf86cd799439011 -t expense
```

**Flags**:
- `-n, --name` - Category name (required)
- `-p, --parent` - Parent category ID (optional)
- `-t, --type` - Category type: income or expense (optional)
- `-r, --remark` - Description or notes (optional)

**Note**: Category is created for your user account.

---

### category update
Update existing category.

```bash
cashlenx category update -i 507f1f77bcf86cd799439011 -n "New Name"
cashlenx category update -i 507f1f77bcf86cd799439011 -p 507f1f77bcf86cd799439012
```

**Flags**:
- `-i, --id` - Category ID (required)
- `-n, --name` - New name (optional)
- `-t, --type` - New type (optional)
- `-p, --parent` - New parent ID (optional)
- `-r, --remark` - New description (optional)

**Note**: Can only update your own categories.

---

### category delete
Delete category.

```bash
cashlenx category delete -i 507f1f77bcf86cd799439011
```

**Flags**:
- `-i, --id` - Category ID (required)

**Note**: Can only delete your own categories.

---

### category query
Query categories by filters.

```bash
# Query by ID
cashlenx category query -i 507f1f77bcf86cd799439011

# Query by name
cashlenx category query -n "Food & Dining"

# Query children of a parent
cashlenx category query -p 507f1f77bcf86cd799439011
```

**Flags**:
- `-i, --id` - Query by ID
- `-n, --name` - Query by name
- `-p, --parent` - Query by parent ID

**Note**: Only searches your own categories.

---

### category list
List all categories.

```bash
# List all
cashlenx category list

# Filter by type
cashlenx category list -t expense
cashlenx category list -t income

# With pagination
cashlenx category list -l 20 -o 40
```

**Flags**:
- `-l, --limit` - Maximum records to return (default: 50)
- `-o, --offset` - Number of records to skip (default: 0)
- `-t, --type` - Filter by type (income/expense)

**Note**: Only returns your own categories.

---

### category tree
Show category tree structure.

```bash
# All categories
cashlenx category tree

# Filter by type
cashlenx category tree -t expense
```

**Flags**:
- `-t, --type` - Filter by type (income/expense)

**Note**: Only shows your own categories.

---

## 🚧 Planned: Statistic Commands (`cashlenx statistic`)

These commands will be available to all authenticated users with proper data isolation.

### statistic export
Export your own data to Excel.

```bash
cashlenx statistic export -o mydata.xlsx
cashlenx statistic export -f 2024-01-01 -t 2024-12-31 -o year2024.xlsx
```

**Note**: Only exports your own data.

### statistic import
Import data to your account.

```bash
cashlenx statistic import -i mydata.xlsx
```

**Note**: Only imports to your account.

### statistic summary
Advanced summaries and analytics.

```bash
cashlenx statistic summary -p daily -d 2024-01-15
cashlenx statistic summary -p monthly -d 2024-01
cashlenx statistic summary -p yearly -d 2024
```

### statistic breakdown
Category breakdown analysis.

```bash
cashlenx statistic breakdown --period month --date 2024-01
```

### statistic trends
Spending trend analysis.

```bash
cashlenx statistic trends --period year --date 2024
```

### statistic top
Top N expenses.

```bash
cashlenx statistic top -n 10 --period month --date 2024-01
```

---

## Data Isolation

### User Data Isolation (Implemented)

All user commands (`cash` and `category`) enforce strict data isolation:

- Users can **only** access their own transactions
- Users can **only** access their own categories
- Attempts to access another user's data return "not found" errors

### Admin Data Access

Admin commands can access data across all users:

- `backup` includes all users' data
- `restore` properly isolates data by user
- `export/import` currently access all data (will be moved to user statistic module)

---

## Examples

### Daily Workflow

```bash
# Add morning coffee
cashlenx cash expense -c "Food & Dining" -a 4.50 -d "Morning coffee"

# Add lunch
cashlenx cash expense -c "Food & Dining" -a 12.00 -d "Lunch"

# Check today's transactions
cashlenx cash query -b $(date +%Y-%m-%d)

# Get monthly summary
cashlenx cash summary -p monthly -d $(date +%Y-%m)
```

### Category Management

```bash
# Create main categories
cashlenx category create -n "Food & Dining" -t expense
cashlenx category create -n "Transportation" -t expense
cashlenx category create -n "Salary" -t income

# Create subcategories
FOOD_ID=$(cashlenx category query -n "Food & Dining" | grep ID | cut -d: -f2)
cashlenx category create -n "Groceries" -p $FOOD_ID -t expense
cashlenx category create -n "Restaurants" -p $FOOD_ID -t expense

# View category tree
cashlenx category tree -t expense
```

### Admin Data Management

```bash
# Create backup before major changes
cashlenx admin backup -o backup_before_cleanup.json

# Restore from backup if needed
cashlenx admin restore -i backup_before_cleanup.json -f
```

---

## Implementation Status

### ✅ Implemented
- **Public**: health, version, start
- **Admin**: backup, restore, export, import
- **Cash**: income, expense, query, delete, list, update, range, summary (monthly)
- **Category**: create, query, delete, update, list, tree

### 🚧 Planned
- **Statistic module**: User-specific export/import with data isolation
- **Advanced analytics**: Category breakdown, spending trends, top expenses
- **Daily/yearly summaries**: Extended summary periods

### ❌ Removed (Rarely Used)
- DB connect - Rarely needed
- DB stats - Use monitoring tools
- DB init/seed - Development only
- DB reset/truncate - Dangerous
- DB indexes - Handled by migrations

---

## Building for Production

```bash
# Build with version info
cd backend
go build -ldflags "\
  -X github.com/macar-x/cashlenx-server/cmd/open_cmd.Version=2.0.0 \
  -X github.com/macar-x/cashlenx-server/cmd/open_cmd.BuildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ) \
  -X github.com/macar-x/cashlenx-server/cmd/open_cmd.GitCommit=$(git rev-parse --short HEAD)" \
  -o cashlenx main.go

# Install globally
sudo mv cashlenx /usr/local/bin/
```

---

## Troubleshooting

### Authentication Issues

```bash
# Ensure JWT_SECRET is set
echo $JWT_SECRET

# Check server logs for auth errors
cashlenx open start -p 8080  # Check output for errors
```

### Database Connection Issues

```bash
# Test server health
cashlenx open health

# Check environment variables
echo $MONGO_DB_URI
echo $DB_TYPE
echo $DB_NAME
```

### Permission Denied

```bash
# Make binary executable
chmod +x cashlenx

# Check if it's in PATH
which cashlenx

# Use full path if needed
./cashlenx open version
```

---

## See Also

- [API Documentation](./api.md)
- [Feature Parity Matrix](./FEATURE_PARITY.md)
- [Quick Start Guide](./quick_start.md)
- [Deployment Guide](./deployment_guide.md)

---

## Support

For issues and feature requests:
https://github.com/macar-x/cashlenx-server/issues
