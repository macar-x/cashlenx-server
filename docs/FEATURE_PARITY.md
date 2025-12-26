# API vs CLI Feature Parity

**Last Updated**: 2025-01-26
**Purpose**: Track feature parity between API and CLI interfaces

## Current Implementation Status

### ✅ Public Features (No Authentication)

| Feature | API Endpoint | CLI Command | Status |
|---------|-------------|-------------|---------|
| Health Check | `GET /api/open/health` | `cashlenx open health` | ✅ Both |
| Version Info | `GET /api/open/version` | `cashlenx open version` | ✅ Both |
| Server Start | N/A | `cashlenx open start` | ✅ CLI only |
| User Login | `POST /api/open/auth/login` | N/A | ✅ API only |
| User Registration | `POST /api/open/auth/register` | N/A | ✅ API only |

### ✅ User Features (Authentication Required)

#### Cash Flow Operations
| Feature | API Endpoint | CLI Command | Status | User Isolation |
|---------|-------------|-------------|---------|----------------|
| Create Income | `POST /api/cash/income` | `cashlenx cash income` | ✅ Both | ✅ Yes |
| Create Expense | `POST /api/cash/expense` | `cashlenx cash expense` | ✅ Both | ✅ Yes |
| Query by ID | `GET /api/cash/{id}` | `cashlenx cash query -i {id}` | ✅ Both | ✅ Yes |
| Query by Date | `GET /api/cash/date/{date}` | `cashlenx cash query -b {date}` | ✅ Both | ✅ Yes |
| Update by ID | `PUT /api/cash/{id}` | `cashlenx cash update` | ✅ Both | ✅ Yes |
| Delete by ID | `DELETE /api/cash/{id}` | `cashlenx cash delete -i {id}` | ✅ Both | ✅ Yes |
| Delete by Date | `DELETE /api/cash/date/{date}` | `cashlenx cash delete -b {date}` | ✅ Both | ✅ Yes |
| List All | `GET /api/cash?limit=N&offset=M` | `cashlenx cash list` | ✅ Both | ✅ Yes |
| Query Range | `GET /api/cash/range?from=X&to=Y` | `cashlenx cash range` | ✅ Both | ✅ Yes |
| Monthly Summary | `GET /api/cash/summary/monthly/{yyyymm}` | `cashlenx cash summary` | ✅ Both | ✅ Yes |

#### Category Operations
| Feature | API Endpoint | CLI Command | Status | User Isolation |
|---------|-------------|-------------|---------|----------------|
| Create Category | `POST /api/category` | `cashlenx category create` | ✅ Both | ✅ Yes |
| List Categories | `GET /api/category?limit=N&offset=M` | `cashlenx category list` | ✅ Both | ✅ Yes |
| Query by ID | `GET /api/category/{id}` | `cashlenx category query -i {id}` | ✅ Both | ✅ Yes |
| Query by Name | `GET /api/category/name/{name}` | `cashlenx category query -n {name}` | ✅ Both | ✅ Yes |
| Get Child Categories | `GET /api/category/{id}/children` | `cashlenx category query -p {id}` | ✅ Both | ✅ Yes |
| Update Category | `PUT /api/category/{id}` | `cashlenx category update` | ✅ Both | ✅ Yes |
| Delete Category | `DELETE /api/category/{id}` | `cashlenx category delete` | ✅ Both | ✅ Yes |
| Get Category Tree | `GET /api/category/tree` | `cashlenx category tree` | ✅ Both | ✅ Yes |

### ✅ Admin Features (Admin Only)

| Feature | API Endpoint | CLI Command | Status | Notes |
|---------|-------------|-------------|---------|-------|
| Create User | `POST /api/admin/user` | N/A | ✅ API only | Admin management |
| List Users | `GET /api/admin/user` | N/A | ✅ API only | Admin management |
| Query User | `GET /api/admin/user/{id}` | N/A | ✅ API only | Admin management |
| Update User | `PUT /api/admin/user/{id}` | N/A | ✅ API only | Admin management |
| Delete User | `DELETE /api/admin/user/{id}` | N/A | ✅ API only | Admin management |
| Database Backup | `GET /api/admin/manage/dump` | `cashlenx admin backup` | ✅ Both | Creates backup file |
| Database Restore | `POST /api/admin/manage/restore` | `cashlenx admin restore` | ✅ Both | Restores from backup |
| Export to Excel | `GET /api/admin/manage/export` | `cashlenx admin export` | ✅ Both | **TODO**: Move to user statistic module |
| Import from Excel | `POST /api/admin/manage/import` | `cashlenx admin import` | ✅ Both | **TODO**: Move to user statistic module |

### ❌ Removed Features (Rarely Used)

| Feature | Reason | Alternative |
|---------|--------|-------------|
| DB Connect | Rarely needed | Use application logs |
| DB Seed/Init | Development only | Manual data creation |
| DB Stats | Rarely used | Use monitoring tools |
| DB Reset/Truncate | Dangerous, rarely used | Manual database operations |
| DB Indexes | Rarely needed | Migrations handle this |

## 🚧 Planned: User Statistic Module

**Purpose**: Allow all authenticated users to analyze their own data with proper isolation

### Planned Features
| Feature | API Endpoint | CLI Command | User Isolation |
|---------|-------------|-------------|----------------|
| Export User Data | `GET /api/statistic/export` | `cashlenx statistic export` | ✅ Yes |
| Import User Data | `POST /api/statistic/import` | `cashlenx statistic import` | ✅ Yes |
| Daily Summary | `GET /api/statistic/summary/daily/{date}` | `cashlenx statistic summary -p daily` | ✅ Yes |
| Monthly Summary | `GET /api/statistic/summary/monthly/{yyyymm}` | `cashlenx statistic summary -p monthly` | ✅ Yes |
| Yearly Summary | `GET /api/statistic/summary/yearly/{yyyy}` | `cashlenx statistic summary -p yearly` | ✅ Yes |
| Category Breakdown | `GET /api/statistic/category-breakdown` | `cashlenx statistic breakdown` | ✅ Yes |
| Spending Trends | `GET /api/statistic/trends` | `cashlenx statistic trends` | ✅ Yes |
| Top Expenses | `GET /api/statistic/top-expenses?limit=N` | `cashlenx statistic top -n N` | ✅ Yes |

**Migration Plan**:
1. Move `export` and `import` from admin to statistic module
2. Add user-specific data isolation to export/import
3. Implement summary and analytics features
4. Deprecate admin export/import endpoints

## Architecture Notes

### User Data Isolation (Implemented)
- **Three-layer architecture**: Mapper → Service → Controller
- **Mapper layer**: `*AndUser()` methods enforce database-level isolation
- **Service layer**: `*ForUser()` methods provide business logic with user context
- **Controller layer**: Extract `userId` from JWT and pass to services
- **Applied to**: Cash flows, Categories, Backup/Restore

### Route Organization (Implemented)
- **Public routes**: `/api/open/*` - No authentication required
- **Admin routes**: `/api/admin/*` - Admin role required
- **User routes**: `/api/cash/*`, `/api/category/*` - Authentication required, user-specific data

### CLI Organization (Implemented)
- **Public commands**: `cashlenx open` - No server needed (health/version need server, start doesn't)
- **Admin commands**: `cashlenx admin` - Admin privileges required
- **User commands**: `cashlenx cash`, `cashlenx category` - Authentication required

## Testing Checklist

### User Isolation Testing
- [ ] User A cannot access User B's cash flows
- [ ] User A cannot access User B's categories
- [ ] User A can only export their own data
- [ ] User A can only import to their own account
- [ ] Backup includes all users' data (admin only)
- [ ] Restore properly isolates data by user

### Feature Parity Testing
- [ ] All API endpoints have CLI equivalents (or documented reason why not)
- [ ] All CLI commands have API equivalents (or documented reason why not)
- [ ] Response formats are consistent between API and CLI
- [ ] Error handling is consistent between API and CLI

## Known Gaps

### API-only Features
- User management (login, register, user CRUD) - Makes sense, CLI users are already authenticated

### CLI-only Features
- Server start command - Makes sense, can't start server via API

### Missing Features (both API and CLI)
- User statistic module (planned)
- User-specific export/import with isolation (planned)
- Advanced analytics (planned)

## Version History

### v2.0.0 (Current)
- ✅ Implemented user data isolation for cash flows and categories
- ✅ Reorganized API routes into /open and /admin
- ✅ Reorganized CLI into open and admin structure
- ✅ Removed rarely-used commands
- ✅ Added backup/restore with user support
- 🚧 Planning statistic module

### v1.0.0 (Previous)
- Basic cash flow and category CRUD
- No user isolation
- No authentication
