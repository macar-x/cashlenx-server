# Statistic Module Implementation Plan

**Status**: 🚧 In Progress
**Version**: 2.1.0

## Overview

The statistic module provides user-specific analytics and data management features with proper data isolation. This allows all authenticated users to analyze and export their own financial data.

## Architecture

### Three-Layer Architecture (Same as cash/category)

```
Controller Layer (extract userId from JWT)
    ↓
Service Layer (business logic with user context)
    ↓
Mapper Layer (database queries with user filtering)
```

### Module Structure

```
CLI:
cmd/statistic_cmd/
├── root.go           # Main statistic command
├── export.go         # Export user data to Excel
├── import.go         # Import user data from Excel
├── summary.go        # Financial summaries (daily/monthly/yearly)
├── breakdown.go      # Category breakdown analysis
├── trends.go         # Spending trends over time
└── top.go            # Top N expenses

API:
controller/statistic_controller/
├── export.go         # GET /api/statistic/export
├── import.go         # POST /api/statistic/import
├── summary.go        # GET /api/statistic/summary/{period}/{date}
├── breakdown.go      # GET /api/statistic/breakdown
├── trends.go         # GET /api/statistic/trends
└── top.go            # GET /api/statistic/top

service/statistic_service/
├── export.go         # Export logic with user filtering
├── import.go         # Import logic with user association
├── summary.go        # Calculate summaries for user
├── breakdown.go      # Analyze category breakdown
├── trends.go         # Analyze spending trends
└── top.go            # Get top expenses
```

## Features

### 1. Export (Migrated from admin)

**CLI**: `cashlenx statistic export -o mydata.xlsx`
**API**: `GET /api/statistic/export?from=2024-01-01&to=2024-12-31`

**User Isolation**: ✅ Only exports the authenticated user's data

**Changes from admin version**:
- Add userId parameter to service layer
- Filter cash flows by userId
- Filter categories by userId
- Update service to use *ForUser methods

### 2. Import (Migrated from admin)

**CLI**: `cashlenx statistic import -i mydata.xlsx`
**API**: `POST /api/statistic/import`

**User Isolation**: ✅ Only imports to the authenticated user's account

**Changes from admin version**:
- Add userId parameter to service layer
- Associate all imported records with userId
- Validate categories belong to user
- Use *ForUser methods for creation

### 3. Summary (New Feature)

**CLI**:
```bash
cashlenx statistic summary -p daily -d 2024-01-15
cashlenx statistic summary -p monthly -d 2024-01
cashlenx statistic summary -p yearly -d 2024
```

**API**:
```
GET /api/statistic/summary/daily/2024-01-15
GET /api/statistic/summary/monthly/202401
GET /api/statistic/summary/yearly/2024
```

**Response**:
```json
{
  "period": "2024-01",
  "period_type": "monthly",
  "income": 5000.00,
  "expense": 2500.00,
  "balance": 2500.00,
  "transaction_count": 15,
  "income_count": 2,
  "expense_count": 13,
  "average_transaction": 166.67,
  "categories": {
    "Food & Dining": 500.00,
    "Transportation": 200.00
  }
}
```

**User Isolation**: ✅ Only includes the authenticated user's transactions

### 4. Category Breakdown (New Feature)

**CLI**: `cashlenx statistic breakdown -p month -d 2024-01`

**API**: `GET /api/statistic/breakdown?period=month&date=2024-01`

**Response**:
```json
{
  "period": "2024-01",
  "total_expense": 2500.00,
  "total_income": 5000.00,
  "expense_categories": [
    {
      "category": "Food & Dining",
      "amount": 500.00,
      "percentage": 20.0,
      "count": 8
    },
    {
      "category": "Transportation",
      "amount": 200.00,
      "percentage": 8.0,
      "count": 3
    }
  ],
  "income_categories": [
    {
      "category": "Salary",
      "amount": 5000.00,
      "percentage": 100.0,
      "count": 1
    }
  ]
}
```

**User Isolation**: ✅ Only analyzes the authenticated user's data

### 5. Trends Analysis (New Feature)

**CLI**: `cashlenx statistic trends -p year -d 2024`

**API**: `GET /api/statistic/trends?period=year&date=2024`

**Response**:
```json
{
  "period": "2024",
  "period_type": "year",
  "data_points": [
    {
      "date": "2024-01",
      "income": 5000.00,
      "expense": 2500.00,
      "balance": 2500.00
    },
    {
      "date": "2024-02",
      "income": 5000.00,
      "expense": 2800.00,
      "balance": 2200.00
    }
  ],
  "trends": {
    "income_trend": "stable",
    "expense_trend": "increasing",
    "average_monthly_expense": 2650.00
  }
}
```

**User Isolation**: ✅ Only analyzes the authenticated user's data

### 6. Top Expenses (New Feature)

**CLI**: `cashlenx statistic top -n 10 -p month -d 2024-01`

**API**: `GET /api/statistic/top?limit=10&period=month&date=2024-01`

**Response**:
```json
{
  "period": "2024-01",
  "limit": 10,
  "total_expense": 2500.00,
  "expenses": [
    {
      "id": "507f1f77bcf86cd799439011",
      "date": "2024-01-15",
      "category": "Electronics",
      "amount": 500.00,
      "description": "New laptop",
      "percentage": 20.0
    },
    {
      "id": "507f1f77bcf86cd799439012",
      "date": "2024-01-10",
      "category": "Food & Dining",
      "amount": 120.00,
      "description": "Restaurant",
      "percentage": 4.8
    }
  ]
}
```

**User Isolation**: ✅ Only includes the authenticated user's expenses

## Implementation Steps

### Phase 1: CLI Structure
1. Create `cmd/statistic_cmd/` directory
2. Create root.go with StatisticCmd
3. Create export.go (migrate from admin with user context)
4. Create import.go (migrate from admin with user context)
5. Create summary.go (new feature)
6. Create breakdown.go (new feature)
7. Create trends.go (new feature)
8. Create top.go (new feature)
9. Update `cmd/root.go` to register StatisticCmd

### Phase 2: Service Layer
1. Create `service/statistic_service/` directory
2. Implement ExportForUser() - migrate from manage_service
3. Implement ImportForUser() - migrate from manage_service
4. Implement GetSummaryForUser()
5. Implement GetBreakdownForUser()
6. Implement GetTrendsForUser()
7. Implement GetTopExpensesForUser()

### Phase 3: API Controller
1. Create `controller/statistic_controller/` directory
2. Implement export.go (extract userId from JWT)
3. Implement import.go (extract userId from JWT)
4. Implement summary.go (extract userId from JWT)
5. Implement breakdown.go (extract userId from JWT)
6. Implement trends.go (extract userId from JWT)
7. Implement top.go (extract userId from JWT)
8. Register routes in `controller/server.go`

### Phase 4: Testing & Cleanup
1. Test all CLI commands with user isolation
2. Test all API endpoints with user isolation
3. Verify user A cannot access user B's statistics
4. Update documentation
5. Deprecate admin export/import (keep for backward compatibility)

## Migration Strategy

### Admin Export/Import → Statistic Export/Import

**Keep admin version** for backward compatibility (will access all data, admin only)
**Add statistic version** with user isolation (user-specific data)

This allows:
- Admins can still export/import all data via admin commands
- Users can export/import their own data via statistic commands
- Gradual migration without breaking existing workflows

## Success Criteria

- ✅ All statistic commands enforce user data isolation
- ✅ Users can only export/import their own data
- ✅ All analytics only show user's own transactions
- ✅ Feature parity between API and CLI
- ✅ Comprehensive error handling
- ✅ Documentation updated

## Testing Plan

```bash
# Test export with user isolation
TOKEN_USER_A="..." # User A's token
TOKEN_USER_B="..." # User B's token

# User A exports their data
cashlenx statistic export -o user_a_data.xlsx
curl -H "Authorization: Bearer $TOKEN_USER_A" \
  http://localhost:8080/api/statistic/export > user_a_api.xlsx

# User B exports their data
cashlenx statistic export -o user_b_data.xlsx
curl -H "Authorization: Bearer $TOKEN_USER_B" \
  http://localhost:8080/api/statistic/export > user_b_api.xlsx

# Verify files contain different data
# Verify User A cannot see User B's transactions in exports
```

## Timeline

- Phase 1 (CLI): ~30 minutes
- Phase 2 (Service): ~45 minutes
- Phase 3 (API): ~30 minutes
- Phase 4 (Testing): ~15 minutes

**Total**: ~2 hours
