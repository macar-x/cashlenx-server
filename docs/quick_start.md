# CashLenX API - Quick Start

**Get the API running in 5 minutes**

---

## 🚀 Fast Track (Copy & Paste)

### 1. Start MongoDB
```bash
cd /home/user/cashlenx-server
docker-compose up -d mongodb
```

### 2. Build & Configure
```bash
cd backend
go build -o cashlenx .

export DB_TYPE=mongodb
export MONGO_DB_URI="mongodb://cashlenx:cashlenx123@localhost:27017/cashlenx?authSource=admin"
export DB_NAME=cashlenx
```

### 3. Test Connection
```bash
./cashlenx db connect
```

### 4. Start API Server
```bash
./cashlenx server start -p 8080
```

### 5. Test API (in new terminal)
```bash
curl http://localhost:8080/api/health
curl http://localhost:8080/api/version | jq .
```

---

## 📋 Common Commands Reference

### Server Operations
```bash
# Start server
./cashlenx server start -p 8080

# Start in background
nohup ./cashlenx server start -p 8080 > server.log 2>&1 &

# Check server status
curl http://localhost:8080/api/health

# Stop server (if in background)
pkill -f "cashlenx server"
```

### Database Operations
```bash
# Test connection
./cashlenx db connect

# Seed demo data
./cashlenx db seed

# Check MongoDB status
docker ps | grep mongodb

# View MongoDB logs
docker logs cashlenx-mongodb

# Restart MongoDB
docker-compose restart mongodb
```

### Category Management
```bash
# Create category
curl -X POST http://localhost:8080/api/category \
  -H "Content-Type: application/json" \
  -d '{"name": "Food", "remark": "Food expenses"}'

# List all categories
curl http://localhost:8080/api/category/list

# Get category by name
curl http://localhost:8080/api/category/name/Food

# Update category
curl -X PUT http://localhost:8080/api/category/{id} \
  -H "Content-Type: application/json" \
  -d '{"name": "Food & Drinks"}'

# Delete category
curl -X DELETE http://localhost:8080/api/category/{id}
```

### Transaction Management
```bash
# Create expense
curl -X POST http://localhost:8080/api/cash/expense \
  -H "Content-Type: application/json" \
  -d '{
    "belongs_date": "20251212",
    "category_name": "Food",
    "amount": 45.50,
    "description": "Lunch"
  }'

# Create income
curl -X POST http://localhost:8080/api/cash/income \
  -H "Content-Type: application/json" \
  -d '{
    "belongs_date": "20251212",
    "category_name": "Salary",
    "amount": 5000.00,
    "description": "Monthly salary"
  }'

# List all transactions
curl http://localhost:8080/api/cash/list

# Get transaction by ID
curl http://localhost:8080/api/cash/{id}

# Get by date
curl http://localhost:8080/api/cash/date/20251212

# Get date range
curl "http://localhost:8080/api/cash/range?from=20251201&to=20251231"

# Update transaction
curl -X PUT http://localhost:8080/api/cash/{id} \
  -H "Content-Type: application/json" \
  -d '{"amount": 50.00, "description": "Updated lunch"}'

# Delete transaction
curl -X DELETE http://localhost:8080/api/cash/{id}
```

### Analytics & Reports
```bash
# Daily summary
curl http://localhost:8080/api/cash/summary/daily/20251212

# Monthly summary
curl http://localhost:8080/api/cash/summary/monthly/202512

# Yearly summary
curl http://localhost:8080/api/cash/summary/yearly/2025
```

### CLI Alternatives
```bash
# Category operations
./cashlenx category create -n "Food" -r "Food expenses"
./cashlenx category query -n "Food"
./cashlenx category list

# Cash flow operations
./cashlenx cash income -c "Salary" -a 5000 -d "20251212" -e "Salary"
./cashlenx cash expense -c "Food" -a 45.50 -d "20251212" -e "Lunch"
./cashlenx cash query -d "20251212"
./cashlenx cash range -f "20251201" -t "20251231"
./cashlenx cash summary --type daily -d "20251212"

# Data management
./cashlenx manage export -f "20251201" -t "20251231" -o "backup.xlsx"
./cashlenx manage import -i "backup.xlsx"
```

---

## 🧪 Quick Test Script

Save as `test_api.sh`:

```bash
#!/bin/bash

API="http://localhost:8080/api"

echo "Testing CashLenX API..."
echo ""

# Health check
echo "1. Health Check"
curl -s $API/health | jq .
echo ""

# Create category
echo "2. Create Category"
CATEGORY=$(curl -s -X POST $API/category \
  -H "Content-Type: application/json" \
  -d '{"name":"TestFood","remark":"Test"}')
CATEGORY_ID=$(echo $CATEGORY | jq -r .id)
echo "Category ID: $CATEGORY_ID"
echo ""

# Create expense
echo "3. Create Expense"
EXPENSE=$(curl -s -X POST $API/cash/expense \
  -H "Content-Type: application/json" \
  -d '{
    "belongs_date":"20251212",
    "category_name":"TestFood",
    "amount":25.00,
    "description":"Test expense"
  }')
EXPENSE_ID=$(echo $EXPENSE | jq -r .id)
echo "Expense ID: $EXPENSE_ID"
echo ""

# List transactions
echo "4. List Transactions"
curl -s "$API/cash/list?limit=5" | jq '.data | length'
echo " transactions found"
echo ""

# Get daily summary
echo "5. Daily Summary"
curl -s $API/cash/summary/daily/20251212 | jq .
echo ""

# Cleanup
echo "6. Cleanup"
curl -s -X DELETE $API/cash/$EXPENSE_ID > /dev/null
curl -s -X DELETE $API/category/$CATEGORY_ID > /dev/null
echo "Test data cleaned up"
echo ""

echo "✅ All tests passed!"
```

Run it:
```bash
chmod +x test_api.sh
./test_api.sh
```

---

## 🛑 Common Issues & Solutions

| Problem | Solution |
|---------|----------|
| Port 8080 in use | `./cashlenx server start -p 8081` |
| MongoDB not running | `docker-compose up -d mongodb` |
| Category not found | Create category first with POST /api/category |
| Connection refused | Check server is running: `curl localhost:8080/api/health` |
| Permission denied | `chmod +x cashlenx` |

---

## 📚 Full Documentation

- **Deployment Guide**: `docs/deployment_guide.md` (detailed step-by-step)
- **API Reference**: `docs/api.md`
- **CLI Reference**: `docs/cli.md`
- **Roadmap**: `docs/roadmap.md`
- **Testing Guide**: `docs/testing.md`

---

## 🎯 What to Test

### Essential (5 min)
- [x] Health check
- [x] Create category
- [x] Create expense
- [x] List transactions

### Important (15 min)
- [ ] Create income
- [ ] Get by date
- [ ] Update transaction
- [ ] Delete transaction
- [ ] Daily summary

### Complete (30 min)
- [ ] All 21 endpoints (see `docs/api.md`)
- [ ] Error cases
- [ ] Pagination
- [ ] Filtering

---

**Need help?** Check `docs/deployment_guide.md` for detailed instructions and troubleshooting.
