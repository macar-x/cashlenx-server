# Backend API Implementation TODO

## Completed ✅

### Infrastructure
- [x] CORS middleware
- [x] Logging middleware  
- [x] Health check endpoint (`GET /api/health`)
- [x] Version info endpoint (`GET /api/version`)

### Cash Flow API
- [x] `POST /api/cash/expense` - Create expense
- [x] `POST /api/cash/income` - Create income
- [x] `GET /api/cash/{id}` - Query by ID
- [x] `GET /api/cash/date/{date}` - Query by date
- [x] `DELETE /api/cash/{id}` - Delete by ID
- [x] `DELETE /api/cash/date/{date}` - Delete by date

### Manage API
- [x] `GET /api/manage/dump` - Download database dump (requires ADMIN_TOKEN)
- [x] `POST /api/manage/restore` - Restore database from dump (requires ADMIN_TOKEN)
- [x] `POST /api/manage/truncate` - Truncate database (requires ADMIN_TOKEN)
- [x] `GET /api/manage/export` - Export data to Excel
- [x] `POST /api/manage/import` - Import data from Excel

## To Implement 🚧

### Cash Flow API Extensions
- [ ] `PUT /api/cash/{id}` - Update cash flow record
- [ ] `GET /api/cash/range?from={date}&to={date}` - Query by date range
- [ ] `GET /api/cash/summary/daily?date={date}` - Daily summary
- [ ] `GET /api/cash/summary/monthly?year={year}&month={month}` - Monthly summary
- [ ] `GET /api/cash/summary/yearly?year={year}` - Yearly summary

### Category API
- [x] `POST /api/category` - Create category
- [x] `GET /api/category` - List all categories
- [x] `GET /api/category/{id}` - Get category by ID
- [x] `PUT /api/category/{id}` - Update category
- [x] `DELETE /api/category/{id}` - Delete category
- [ ] `GET /api/category/{id}/stats` - Category statistics

### Statistics API
- [ ] `GET /api/stats/overview?period={period}` - Financial overview
- [ ] `GET /api/stats/trends?period={period}` - Spending trends
- [ ] `GET /api/stats/category-breakdown?period={period}` - Category breakdown
- [ ] `GET /api/stats/income-vs-expense?period={period}` - Income vs expense
- [ ] `GET /api/stats/top-expenses?limit={n}&period={period}` - Top N expenses

### Import/Export API
- [ ] `POST /api/export` - Export to Excel
- [ ] `POST /api/import` - Import from Excel
- [ ] `GET /api/export/csv` - Export to CSV

## Implementation Guide

### 1. Update Cash Flow Record

**Endpoint**: `PUT /api/cash/{id}`

**Request Body**:
```json
{
  "amount": 100.50,
  "date": "2024-01-15",
  "category": "Food & Dining",
  "description": "Updated description"
}
```

**Implementation**:
```go
// In cash_flow_controller package
func UpdateById(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]
    
    var cashFlow model.CashFlow
    json.NewDecoder(r.Body).Decode(&cashFlow)
    
    // Update in database
    err := service.UpdateCashFlow(id, cashFlow)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(cashFlow)
}
```

### 2. Date Range Query

**Endpoint**: `GET /api/cash/range?from=2024-01-01&to=2024-01-31`

**Response**:
```json
{
  "from": "2024-01-01",
  "to": "2024-01-31",
  "count": 15,
  "total_income": 5000.00,
  "total_expense": 2500.00,
  "balance": 2500.00,
  "transactions": [...]
}
```

### 3. Summary Endpoints

**Daily**: `GET /api/cash/summary/daily?date=2024-01-15`
**Monthly**: `GET /api/cash/summary/monthly?year=2024&month=1`
**Yearly**: `GET /api/cash/summary/yearly?year=2024`

**Response Format**:
```json
{
  "period": "2024-01",
  "income": 5000.00,
  "expense": 2500.00,
  "balance": 2500.00,
  "transaction_count": 15,
  "categories": {
    "Food & Dining": 500.00,
    "Transportation": 200.00
  }
}
```

### 4. Category Management

Create new controller: `backend/controller/category_controller/`

**Files needed**:
- `create.go` - Create category
- `query.go` - List and get by ID
- `update.go` - Update category
- `delete.go` - Delete category
- `stats.go` - Category statistics

### 5. Statistics API

Create new controller: `backend/controller/stats_controller/`

**Aggregation queries needed**:
- Group by date for trends
- Group by category for breakdown
- Sum by type for income vs expense
- Sort and limit for top expenses

## Testing

### Manual Testing

```bash
# Start MongoDB
docker compose up -d mongodb

# Start backend
cd backend
export MONGO_DB_URI="mongodb://cashlenx:cashlenx123@localhost:27017/cashlenx?authSource=admin"
go run main.go server start -p 8080

# Test health
curl http://localhost:8080/api/health

# Test version
curl http://localhost:8080/api/version

# Test CORS
curl -H "Origin: http://localhost:3000" \
     -H "Access-Control-Request-Method: GET" \
     -H "Access-Control-Request-Headers: Content-Type" \
     -X OPTIONS \
     http://localhost:8080/api/health

# Test cash flow
curl http://localhost:8080/api/cash/date/$(date +%Y-%m-%d)
```

### Integration Testing

Create test files:
- `backend/controller/server_test.go`
- `backend/controller/cash_flow_controller/controller_test.go`

## Database Schema

### MongoDB Collections

**cash_flows**:
```javascript
{
  _id: ObjectId,
  amount: Number,
  date: String (YYYY-MM-DD),
  category: String,
  type: String (income/outcome),
  description: String,
  created_at: Date,
  updated_at: Date
}
```

**categories**:
```javascript
{
  _id: ObjectId,
  name: String,
  icon: String,
  color: String,
  type: String (income/expense),
  created_at: Date
}
```

### MySQL Tables

Already defined in `docker/mysql/init-mysql.sql`

## Priority Order

1. **High Priority** (Core functionality):
   - PUT /api/cash/{id} - Update records
   - GET /api/cash/range - Date range queries
   - Category CRUD operations

2. **Medium Priority** (Enhanced features):
   - Summary endpoints
   - Statistics API
   - Category statistics

3. **Low Priority** (Nice to have):
   - Import/Export API
   - Backup/Restore

## Notes

- All endpoints should return proper HTTP status codes
- Use consistent error response format
- Add input validation
- Consider pagination for list endpoints
- Add query parameters for filtering and sorting
- Document all endpoints with examples
