# Template GO Clean Architecture


---

### **Struktur Folder**

```
project-root/
├── cmd/                # Main entry point aplikasi
│   └── main.go         # File utama untuk menjalankan aplikasi
├── config/             # Konfigurasi aplikasi (e.g., .env, config parsers)
├── internal/           # Folder utama untuk logic aplikasi
│   ├── domain/         # Entitas dan kontrak (interface)
│   │   └── product.go  # Definisi entity Product dan interface Repository
│   ├── usecase/        # Business logic
│   │   └── product.go  # Implementasi use case Product
│   ├── repository/     # Implementasi data storage (e.g., database)
│   │   └── product.go  # Repository untuk Product
│   ├── delivery/       # Delivery layer (HTTP handler, CLI, dll.)
│   │   ├── http/       # HTTP handler untuk REST API
│   │   │   └── product.go
│   │   └── middleware/ # Middleware HTTP
│   └── models/         # Struktur data yang digunakan antar layer
│       └── product.go
├── pkg/                # Helper functions, libraries, atau utilitas
│   └── logger/         # Logger custom
└── go.mod              # Dependency management file
```

---

### **Penjelasan Tiap Folder**

1. **`cmd/`**
    - Entry point untuk aplikasi.
    - `main.go` berisi inisialisasi aplikasi, seperti menghubungkan database, memulai server HTTP, dll.

2. **`config/`**
    - Konfigurasi aplikasi seperti variabel environment, pengaturan database, atau konfigurasi server.

3. **`internal/`**
    - Folder utama untuk logic aplikasi. Berisi beberapa subfolder:
        - **`domain/`**:
            - Definisi **entity** dan kontrak seperti interface repository.
            - Tidak bergantung pada implementasi spesifik.
        - **`usecase/`**:
            - Berisi implementasi business logic yang memanfaatkan entitas dari `domain/`.
        - **`repository/`**:
            - Implementasi akses data (e.g., SQL, NoSQL).
        - **`delivery/`**:
            - Layer untuk komunikasi dengan dunia luar (e.g., REST API, CLI).
        - **`models/`**:
            - Struktur data yang digunakan antar layer.

4. **`pkg/`**
    - Library atau utilitas yang reusable di luar aplikasi ini.

---

### **Contoh Implementasi**

#### **1. Entity di `internal/domain/product.go`**

```go
package domain

type Product struct {
    ID          int
    Name        string
    Description string
    Price       float64
    Stock       int
}

type ProductRepository interface {
    Save(product *Product) error
    FindByID(id int) (*Product, error)
    FindAll() ([]*Product, error)
    Delete(id int) error
}
```

---

#### **2. Use Case di `internal/usecase/product.go`**

```go
package usecase

import "project-root/internal/domain"

type ProductUsecase struct {
    repo domain.ProductRepository
}

func NewProductUsecase(repo domain.ProductRepository) *ProductUsecase {
    return &ProductUsecase{repo: repo}
}

func (u *ProductUsecase) CreateProduct(product *domain.Product) error {
    return u.repo.Save(product)
}

func (u *ProductUsecase) GetProductByID(id int) (*domain.Product, error) {
    return u.repo.FindByID(id)
}
```

---

#### **3. Repository di `internal/repository/product.go`**

```go
package repository

import (
    "database/sql"
    "project-root/internal/domain"
)

type ProductRepositoryImpl struct {
    db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepositoryImpl {
    return &ProductRepositoryImpl{db: db}
}

func (r *ProductRepositoryImpl) Save(product *domain.Product) error {
    query := "INSERT INTO products (name, description, price, stock) VALUES (?, ?, ?, ?)"
    _, err := r.db.Exec(query, product.Name, product.Description, product.Price, product.Stock)
    return err
}

func (r *ProductRepositoryImpl) FindByID(id int) (*domain.Product, error) {
    query := "SELECT id, name, description, price, stock FROM products WHERE id = ?"
    row := r.db.QueryRow(query, id)

    product := &domain.Product{}
    err := row.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Stock)
    if err != nil {
        return nil, err
    }
    return product, nil
}
```

---

#### **4. HTTP Handler di `internal/delivery/http/product.go`**

```go
package http

import (
    "net/http"
    "strconv"
    "github.com/gin-gonic/gin"
    "project-root/internal/usecase"
)

type ProductHandler struct {
    usecase *usecase.ProductUsecase
}

func NewProductHandler(router *gin.Engine, uc *usecase.ProductUsecase) {
    handler := &ProductHandler{usecase: uc}
    router.POST("/products", handler.CreateProduct)
    router.GET("/products/:id", handler.GetProductByID)
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
    var product usecase.Product
    if err := c.ShouldBindJSON(&product); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    if err := h.usecase.CreateProduct(&product); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusCreated, product)
}

func (h *ProductHandler) GetProductByID(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
        return
    }
    product, err := h.usecase.GetProductByID(id)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, product)
}
```

---

### **Menjalankan Aplikasi**

**Main File di `cmd/main.go`:**

```go
package main

import (
    "database/sql"
    "log"
    "project-root/internal/delivery/http"
    "project-root/internal/repository"
    "project-root/internal/usecase"

    "github.com/gin-gonic/gin"
    _ "github.com/go-sql-driver/mysql"
)

func main() {
    db, err := sql.Open("mysql", "user:password@tcp(localhost:3306)/krs_management")
    if err != nil {
        log.Fatalf("Could not connect to DB: %v", err)
    }
    defer db.Close()

    repo := repository.NewProductRepository(db)
    uc := usecase.NewProductUsecase(repo)

    router := gin.Default()
    http.NewProductHandler(router, uc)

    log.Println("Starting server at :8080")
    router.Run(":8080")
}
```

---

Dengan struktur ini, Anda memiliki modularitas yang jelas antara domain, use case, repository, dan delivery. Ini memudahkan pengembangan dan pemeliharaan jangka panjang!