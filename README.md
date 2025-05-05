# URL Shortener

![Go](https://img.shields.io/badge/Go-1.18+-00ADD8?logo=go)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-13+-4169E1?logo=postgresql)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

A high-performance URL shortening service with authentication, built with Go and PostgreSQL.

## Features

- 🚀 **RESTful API** with JSON support
- 🔐 JWT Authentication
- 🗄️ PostgreSQL storage with soft deletes
- 📦 Bulk URL operations
- ⚡ Gzip compression middleware
- 🧪 Comprehensive unit tests
- ⚙️ Dual configuration (env vars + CLI flags)

## Quick Start

### Prerequisites
- Go 1.18+
- PostgreSQL 13+

### Installation
```bash
git clone https://github.com/allensuvorov/url-shortener.git
cd url-shortener

# Set up environment
export DATABASE_DSN="postgres://user:password@localhost:5432/dbname?sslmode=disable"
export SERVER_ADDRESS=":8080"

# Run
go run cmd/shortener/main.go
