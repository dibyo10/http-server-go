# 🔐 Go Auth API (In-Memory)

A minimal authentication and user management service built using Go's standard library.

This project implements signup, signin, update, and delete operations with secure password hashing and concurrent-safe in-memory storage.

---

## 🚀 What This Project Does

| Endpoint | Method | Description |
|---|---|---|
| `/signup` | `POST` | Create a new user (password hashed using bcrypt) |
| `/signin` | `POST` | Authenticate user credentials |
| `/users/{id}` | `PUT` | Update user details |
| `/users/{id}` | `DELETE` | Delete a user |
| `/` | `GET` | Health check |

Runs on `localhost:8080`

---

## 🧠 Key Engineering Concepts I Practiced

### 1. Building REST APIs Using Only `net/http`

No frameworks. No Gin. No shortcuts.

Used:
- `http.NewServeMux`
- Pattern-based routing (Go 1.22+)
- Proper HTTP status codes
- JSON encoding/decoding

This forced me to understand how routing actually works under the hood.

### 2. Password Security with bcrypt

```go
bcrypt.GenerateFromPassword()
bcrypt.CompareHashAndPassword()
```

Learned:
- Never store plaintext passwords.
- Always hash before persistence.
- Authentication is hash comparison, not string equality.

### 3. Concurrency Safety with RWMutex

Implemented an in-memory cache:

```go
var userCache = make(map[int]User)
```

Protected with:

```go
sync.RWMutex
```

Used:
- `RLock()` for reads
- `Lock()` for writes

Learned:
- Why maps are not thread-safe in Go.
- The difference between read locks and write locks.
- How race conditions can silently corrupt state.

### 4. Route Parameters (Go 1.22 Feature)

```go
r.PathValue("id")
```

This clarified:
- How Go's new ServeMux supports method-aware routing.
- How path parameters are extracted natively without third-party routers.

### 5. Safe State Mutation

Pattern used in updates:
1. Read under `RLock`
2. Modify local copy
3. Write back under `Lock`

Understood:
- Why structs are copied by value.
- Why you must reassign updated structs into the map.

### 6. API Input Validation

Implemented guard clauses for:
- Missing fields
- Invalid IDs
- Incorrect credentials

Learned:
- Fail fast.
- Return meaningful HTTP status codes.
- Avoid ambiguous responses.

---

## 🧩 Architecture Decisions

- In-memory storage for simplicity
- Mutex-protected shared state
- No persistence layer
- No authentication tokens (stateless demo)

This project focuses on **correctness**, **concurrency safety**, and **clean handler design** — not on production scalability.

---

## 🛠 How to Run

```bash
go run main.go
```

Then test with:

```bash
curl -X POST http://localhost:8080/signup \
  -H "Content-Type: application/json" \
  -d '{"name":"Dibyo","email":"dibyo@mail.com","password":"123456"}'
```

---

## ⚠️ Limitations

- Data resets on restart
- No JWT/session management
- No email uniqueness enforcement
- Linear search for signin (not optimized)
- No database persistence

This is intentionally minimal.

---

## 📈 What This Project Improved in Me

- Comfort with Go's standard library over frameworks
- Stronger understanding of race conditions
- Better mental model of value vs pointer semantics
- Practical API design discipline
- Writing thread-safe stateful services