📘 Migration & Server Commands Cheat Sheet

## 🟦 1. Menjalankan Server

```go
go run ./cmd/server
```

Atau (jika kamu membuat Makefile nanti):

```go
make server
```

Endpoint yang tersedia:

**Endpoint Deskripsi**

```go
GET /health Cek status server
POST /api/users/register Registrasi user
POST /api/users/login Login & mendapatkan JWT
GET /api/protected/profile Cek profile (perlu JWT)
GET /api/admin/stats Akses admin-only
```

🟧 2. Migration Commands

Semua migration command berjalan via:

```go
run ./cmd/migrate
```

Dengan flags -action, -name, -steps, dan -version.

## 🟩 2.1 Create / Generate Migration File

Membuat file .up.sql dan .down.sql otomatis.

```go
go run ./cmd/migrate -action=create -name add_users_table
```

Output:

```
migrations/1736438001_add_users_table.up.sql
migrations/1736438001_add_users_table.down.sql
```

## 🟩 2.2 Apply All Up Migrations

Menjalankan semua migration naik.

```
go run ./cmd/migrate -action=up
```

## 🟩 2.3 Rollback Migration

Rollback 1 langkah:

```
go run ./cmd/migrate -action=down -steps=1
```

Rollback semua:

```
go run ./cmd/migrate -action=down
```

## 🟩 2.4 Lihat Version Migration Saat Ini

```
go run ./cmd/migrate -action=version
```

Output contoh:

```
current version: 1, dirty: false
```

## 🟩 2.5 Force Version (Recovery Mode)

Jika migration pernah FAIL atau “dirty”, gunakan:

```
go run ./cmd/migrate -action=force -version=1
```

Ini memaksa version ke angka tertentu.

**⚠ Gunakan dengan hati-hati!**
