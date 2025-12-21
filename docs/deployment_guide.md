# CashLenX API - Local Docker Deployment & Testing Guide

**Quick guide to deploy and test the CashLenX API locally**

---

## 📋 Prerequisites

- Docker and Docker Compose installed
- Go 1.21+ installed
- Terminal/Command line access
- curl or Postman for API testing

---

## 🚀 Step 1: Start MongoDB with Docker

### Option A: Using Docker Compose (Recommended)

```bash
# Navigate to project root
cd /home/user/cashlenx_server

# Start MongoDB container
docker-compose up -d mongodb

# Verify MongoDB is running
docker-compose ps

# Expected output:
# NAME                 STATUS              PORTS
# cashlenx-mongodb     Up X seconds        0.0.0.0:27017->27017/tcp
```

### Option B: Using Docker CLI Directly

```bash
docker run -d \
  --name cashlenx-mongodb \
  -p 27017:27017 \
  -e MONGO_INITDB_ROOT_USERNAME=cashlenx \
  -e MONGO_INITDB_ROOT_PASSWORD=cashlenx123 \
  -e MONGO_INITDB_DATABASE=cashlenx \
  mongo:7.0
```

### Verify MongoDB Connection

```bash
# Test connection (should return "ping: 1")
docker exec cashlenx-mongodb mongosh \
  --username cashlenx \
  --password cashlenx123 \
  --authenticationDatabase admin \
  --eval "db.adminCommand('ping')"
```

---

## 🔧 Step 2: Build the Backend

```bash
# Navigate to backend directory
cd backend

# Download dependencies
go mod download

# Build the binary
go build -o cashlenx .

# Verify build
ls -lh cashlenx
```

**Expected**: You should see a `cashlenx` binary (size ~20-30 MB)

---

## ⚙️ Step 3: Configure Environment

### Set Environment Variables

```bash
# Database type (mongodb or mysql)
export DB_TYPE=mongodb

# MongoDB connection string
export MONGO_DB_URI="mongodb://cashlenx:cashlenx123@localhost:27017/cashlenx?authSource=admin"

# Database name
export DB_NAME=cashlenx

# Optional: Log file location
export LOG_FOLDER="./logs"
```

### Create a .env file (Alternative)

```bash
cat > .env <<'EOF'
DB_TYPE=mongodb
MONGO_DB_URI=mongodb://cashlenx:cashlenx123@localhost:27017/cashlenx?authSource=admin
DB_NAME=cashlenx
LOG_FOLDER=./logs
EOF

# Load environment
source .env
```

---

## 🧪 Step 4: Test Database Connection

```bash
# Test database connectivity
./cashlenx db connect
```

**Expected output**:
```
Successfully connected to MongoDB!
Database: cashlenx
Collections: (empty or existing collections)
```

**If connection fails**:
- Check MongoDB is running: `docker ps | grep mongodb`
- Check connection string is correct
- Verify credentials: `cashlenx:cashlenx123`

---

## 🎯 Step 5: Seed Test Data (Optional)

```bash
# Create demo categories and transactions
./cashlenx db seed
```

**This will create**:
- 5 categories (Food, Transportation, Salary, etc.)
- 10 sample transactions (mix of income and expenses)

---

## 🚀 Step 6: Start API Server

### Start the server

```bash
# Start on port 8080
./cashlenx server start -p 8080
```

**Expected output**:
```
API server is running on http://localhost:8080
```

### Keep server running in background (Alternative)

```bash
# Start in background
nohup ./cashlenx server start -p 8080 > server.log 2>&1 &

# Get process ID
echo $! > server.pid

# Check it's running
cat server.log
```

---

## ✅ Step 7: Verify API is Running

### Test 1: Health Check

```bash
curl http://localhost:8080/api/system/health
```

**Expected response**:
```json
{
  "code": "OK",
  "message": "",
  "data": {
    "status": "healthy",
    "service": "cashlenx-api",
    "message": "API is running"
  },
  "meta": {},
  "extra": {},
  "errors": []
}
```

### Test 2: Version Info

```bash
curl http://localhost:8080/api/system/version
```

**Expected response**:
```json
{
  "code": "OK",
  "message": "",
  "data": {
    "version": "1.0.0",
    "name": "CashLenX API",
    "description": "Personal finance management API",
    "endpoints": {
      "cash_flow": [...],
      "category": [...],
      "health": [...]
    }
  },
  "meta": {},
  "extra": {},
  "errors": []
}
```

### Test 3: Get Version Info (Formatted)

```bash
curl -s http://localhost:8080/api/system/version | jq .
```

---

## 🧪 Step 8: Run Basic API Tests

### Create a Category

```bash
curl -X POST http://localhost:8080/api/category \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Food",
    "type": "expense",
    "remark": "Food and dining expenses"
  }'
```

**Expected response**:
```json
{
  "code": "OK",
  "message": "",
  "data": {
    "id": "507f1f77bcf86cd799439011",
    "parent_id": "000000000000000000000000",
    "name": "Food",
    "type": "expense",
    "remark": "Food and dining expenses",
    "create_time": "2025-12-12T...",
    "modify_time": "2025-12-12T..."
  },
  "meta": {},
  "extra": {},
  "errors": []
}
```

**Save the category ID** - you'll need it for next steps.

### List All Categories

```bash
curl "http://localhost:8080/api/category?limit=50&offset=0"
```

**Expected response**:
```json
{
  "code": "OK",
  "message": "",
  "data": [
    {
      "id": "507f1f77bcf86cd799439011",
      "parent_id": "000000000000000000000000",
      "name": "Food",
      "type": "expense",
      "remark": "Food and dining expenses",
      "create_time": "2025-12-12T...",
      "modify_time": "2025-12-12T..."
    }
  ],
  "meta": {
    "total_count": 1,
    "limit": 50,
    "offset": 0
  },
  "extra": {},
  "errors": []
}
```

### Create an Expense

```bash
curl -X POST http://localhost:8080/api/cash/expense \
  -H "Content-Type: application/json" \
  -d '{
    "belongs_date": "20251212",
    "category_name": "Food",
    "amount": 45.50,
    "description": "Lunch at restaurant"
  }'
```

**Expected response**:
```json
{
  "code": "OK",
  "message": "",
  "data": {
    "id": "507f1f77bcf86cd799439012",
    "category_id": "507f1f77bcf86cd799439011",
    "belongs_date": "2025-12-12T00:00:00Z",
    "flow_type": "EXPENSE",
    "amount": 45.5,
    "description": "Lunch at restaurant",
    "remark": "",
    "create_time": "2025-12-12T...",
    "modify_time": "2025-12-12T..."
  },
  "meta": {},
  "extra": {},
  "errors": []
}
```

### Create an Income

```bash
curl -X POST http://localhost:8080/api/cash/income \
  -H "Content-Type: application/json" \
  -d '{
    "belongs_date": "20251212",
    "category_name": "Salary",
    "amount": 5000.00,
    "description": "December salary"
  }'
```

**Note**: You may need to create "Salary" category first if it doesn't exist.

### List All Transactions

```bash
curl "http://localhost:8080/api/cash?limit=20&offset=0"
```

### Get Daily Summary

```bash
curl http://localhost:8080/api/cash/summary/daily/20251212
```

**Expected response**:
```json
{
  "code": "OK",
  "message": "",
  "data": {
    "total_income": 5000.0,
    "total_expense": 45.5,
    "balance": 4954.5,
    "transaction_count": 2,
    "category_breakdown": {
      "Food": 45.5,
      "Salary": 5000.0
    }
  },
  "meta": {},
  "extra": {},
  "errors": []
}
```

---

## 📚 Step 9: Run Comprehensive Tests

For complete testing, follow the detailed guide:

```bash
# View the complete test guide
cat docs/api.md

# Or open in browser/editor
code docs/api.md
```

**The guide includes**:
- All 21 endpoint tests
- Phase-by-phase testing strategy
- Expected responses
- Error handling tests
- Performance tests

---

## 🛠️ Troubleshooting

### Problem: Can't connect to MongoDB

**Solution 1**: Check if MongoDB is running
```bash
docker ps | grep mongodb
```

**Solution 2**: Restart MongoDB
```bash
docker-compose restart mongodb
```

**Solution 3**: Check logs
```bash
docker logs cashlenx-mongodb
```

### Problem: API returns "category does not exist"

**Solution**: Create the category first
```bash
curl -X POST http://localhost:8080/api/category \
  -H "Content-Type: application/json" \
  -d '{"name": "Food"}'
```

### Problem: Port 8080 already in use

**Solution**: Use different port
```bash
./cashlenx server start -p 8081
```

Then test with: `curl http://localhost:8081/api/system/health`

### Problem: Permission denied when running ./cashlenx

**Solution**: Make binary executable
```bash
chmod +x cashlenx
```

### Problem: Connection refused

**Solution**: Check firewall and ensure server is running
```bash
# Check if server process is running
ps aux | grep cashlenx

# Check if port is listening
netstat -tuln | grep 8080
# or
lsof -i :8080
```

---

## 🧹 Step 10: Cleanup (When Done Testing)

### Stop API Server

```bash
# If running in foreground: Ctrl+C

# If running in background:
kill $(cat server.pid)
rm server.pid
```

### Stop MongoDB

```bash
# Using docker-compose
docker-compose down mongodb

# Or stop container directly
docker stop cashlenx-mongodb
docker rm cashlenx-mongodb
```

### Keep MongoDB Data (for next session)

```bash
# Just stop, don't remove
docker-compose stop mongodb
```

### Remove All Data (fresh start)

```bash
# Stop and remove containers + volumes
docker-compose down -v
```

---

## 📊 Testing Checklist

After deployment, verify these work:

- [ ] MongoDB is running and accessible
- [ ] Backend builds successfully
- [ ] Database connection works (`./cashlenx db connect`)
- [ ] API server starts on port 8080
- [ ] Health endpoint returns 200 OK
- [ ] Version endpoint lists all endpoints
- [ ] Can create category
- [ ] Can list categories
- [ ] Can create expense
- [ ] Can create income
- [ ] Can list transactions
- [ ] Can get daily summary
- [ ] Can update transaction
- [ ] Can delete transaction

---

## 🎯 Next Steps

1. **Run all 21 endpoint tests** - See `docs/api.md`
2. **Try the CLI commands** - See `docs/cli.md`
3. **Import test data** - Use `./cashlenx manage import`
4. **Export data** - Use `./cashlenx manage export`
5. **Integrate with Flutter** - Start with dashboard API calls

---

## 📖 Additional Resources

- **API Reference**: `docs/api.md`
- **API Reference**: `docs/api.md`
- **CLI Reference**: `docs/cli.md`
- **Environment Setup**: `docs/environment.md`
- **Docker Setup**: `docs/docker.md`

---

## 💡 Quick Tips

### Test with Formatted JSON Output

```bash
# Install jq for pretty JSON
# Ubuntu/Debian: sudo apt install jq
# macOS: brew install jq

# Use with curl
curl -s http://localhost:8080/api/system/version | jq .
```

### Save API Responses

```bash
# Save to file
curl http://localhost:8080/api/category > categories.json

# View formatted
cat categories.json | jq .
```

### Test Multiple Requests

```bash
# Create 10 test expenses
for i in {1..10}; do
  curl -X POST http://localhost:8080/api/cash/expense \
    -H "Content-Type: application/json" \
    -d "{
      \"belongs_date\": \"20251212\",
      \"category_name\": \"Food\",
      \"amount\": $((RANDOM % 100 + 1)),
      \"description\": \"Test expense $i\"
    }"
  echo ""
done
```

### Monitor API Logs

```bash
# If running in background
tail -f server.log

# Or check cashlenx.log
tail -f cashlenx.log
```

---

**Ready to test!** 🚀

Start with Step 1 and work through each step. If you encounter any issues, check the troubleshooting section or refer to the detailed documentation.
