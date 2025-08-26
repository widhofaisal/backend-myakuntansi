# 📋 API Checklist - Backend Project

Project Name: **File Management**  
Last Updated: **2025-07-30**  

---

## 📂 Daftar Endpoint
Base URL: **http://localhost:8080/api/v1**

<br>

### 🔐 Global

| Endpoint             | Method | Deskripsi            | Status | Auth |
|----------------------|--------|----------------------|--------|------|
| `/api/login`         | POST   | Login user           |   ✅   | 🔓   |
| `/api/logout`        | POST   | Logout user          |   ✅   | 🔒   |

---

### 📁 Users Page

| Endpoint               | Method | Deskripsi            | Status | Auth |
|------------------------|--------|----------------------|--------|------|
| `/api/users`           | GET    | Get all users        |   ✅   | 🔒   |
| `/api/users/{id}`      | GET    | Get user by ID       |   ❌   | 🔒   |
| `/api/users`           | POST   | Create user          |   ✅   | 🔒   |
| `/api/users/{id}`      | PUT    | Update user          |   ❌   | 🔒   |
| `/api/users/{id}`      | DELETE | Delete user          |   ❌   | 🔒   |

---

### 📄 Project Page

| Endpoint               | Method | Deskripsi              | Status | Auth |
|------------------------|--------|------------------------|--------|------|
| `/api/projects`        | GET    | Get list projects      |   ✅   | 🔒   |
| `/api/projects/{id}`   | GET    | Get detail project     |   ✅   | 🔒   |
| `/api/projects`        | POST   | Create project         |   ✅   | 🔒   |
| `/api/projects/{id}`   | PUT    | Update project         |   ✅   | 🔒   |
| `/api/projects/{id}`   | DELETE | Delete project         |   ✅   | 🔒   |

---

## 🧩 Catatan Tambahan

- [ ] Tambahkan validasi input untuk semua endpoint
- [ ] Tambahkan dokumentasi Swagger
- [ ] Tambahkan unit test untuk endpoint penting
- [ ] Periksa performa untuk endpoint `/api/data`

---

## 🔗 Dokumentasi API

- [ ] Swagger/OpenAPI: [Link ke dokumentasi Swagger]
- [ ] Postman Collection: [Link atau lokasi file]

