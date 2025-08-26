
### Planner Response

I'll provide you with a comprehensive guide for integrating real backend APIs with your Vue.js file archiving system. Let me analyze the current mock implementations and create a detailed API specification.

*Semantic searched codebase*

### Planner Response

# 🚀 Backend API Integration Guide for Vue.js File Archiving System

Based on my analysis of your frontend code, here's a comprehensive guide for implementing the backend APIs. I've identified all the mock API calls and organized them by functionality.

## 📋 API Endpoints Overview

### **1. Authentication APIs**

#### **POST /auth/login**
**Location:** [src/stores/auth.js](cci:7://file:///d:/CPNS%202025/latsar/aktulisasi/frontend-file-management/web-ai/src/stores/auth.js:0:0-0:0) - [login()](cci:1://file:///d:/CPNS%202025/latsar/aktulisasi/frontend-file-management/web-ai/src/stores/auth.js:17:4-57:5) method
```javascript
// Current mock URL: https://api.example.com/auth/login
```

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "name": "John Doe",
      "email": "user@example.com",
      "username": "johndoe",
      "role": "user", // "admin" or "user"
      "createdAt": "2024-01-15T10:00:00Z",
      "lastLoginAt": "2024-01-20T15:30:00Z",
      "bio": "Optional user bio"
    }
  }
}
```

#### **POST /auth/logout** (Optional)
**Location:** [src/stores/auth.js](cci:7://file:///d:/CPNS%202025/latsar/aktulisasi/frontend-file-management/web-ai/src/stores/auth.js:0:0-0:0) - [logout()](cci:1://file:///d:/CPNS%202025/latsar/aktulisasi/frontend-file-management/web-ai/src/stores/auth.js:59:4-74:5) method
```javascript
// Currently commented out: await axios.post('https://api.example.com/auth/logout')
```

**Request:** No body required (token in Authorization header)
**Response:**
```json
{
  "success": true,
  "message": "Logged out successfully"
}
```

---

### **2. User Profile APIs**

#### **PUT /user/profile**
**Location:** [src/stores/auth.js](cci:7://file:///d:/CPNS%202025/latsar/aktulisasi/frontend-file-management/web-ai/src/stores/auth.js:0:0-0:0) - [updateProfile()](cci:1://file:///d:/CPNS%202025/latsar/aktulisasi/frontend-file-management/web-ai/src/stores/auth.js:87:4-108:5) method
```javascript
// Current mock URL: https://api.example.com/user/profile
```

**Request Body:**
```json
{
  "name": "Updated Name",
  "username": "newusername",
  "email": "newemail@example.com",
  "bio": "Updated bio text"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "user": {
      "id": 1,
      "name": "Updated Name",
      "username": "newusername",
      "email": "newemail@example.com",
      "bio": "Updated bio text",
      "role": "user",
      "createdAt": "2024-01-15T10:00:00Z",
      "lastLoginAt": "2024-01-20T15:30:00Z"
    }
  }
}
```

#### **PUT /user/password**
**Location:** [src/stores/auth.js](cci:7://file:///d:/CPNS%202025/latsar/aktulisasi/frontend-file-management/web-ai/src/stores/auth.js:0:0-0:0) - [changePassword()](cci:1://file:///d:/CPNS%202025/latsar/aktulisasi/frontend-file-management/web-ai/src/stores/auth.js:110:4-125:5) method
```javascript
// Current mock URL: https://api.example.com/user/password
```

**Request Body:**
```json
{
  "currentPassword": "oldpassword123",
  "newPassword": "newpassword456"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Password changed successfully"
}
```

---

### **3. File Management APIs**

#### **GET /files**
**Location:** [src/stores/files.js](cci:7://file:///d:/CPNS%202025/latsar/aktulisasi/frontend-file-management/web-ai/src/stores/files.js:0:0-0:0) - [fetchFiles()](cci:1://file:///d:/CPNS%202025/latsar/aktulisasi/frontend-file-management/web-ai/src/stores/files.js:82:4-102:5) method
```javascript
// Current mock URL: https://api.example.com/files?path=/folder/subfolder
```

**Query Parameters:**
- `path` (string): Current directory path (default: "/")

**Response:**
```json
{
  "success": true,
  "data": {
    "folders": [
      {
        "id": 1,
        "name": "Documents",
        "isFolder": true,
        "path": "/Documents",
        "createdAt": "2024-01-15T10:00:00Z",
        "modifiedAt": "2024-01-15T10:00:00Z",
        "owner": "John Doe",
        "fileCount": 5,
        "folderCount": 2
      }
    ],
    "files": [
      {
        "id": 2,
        "name": "document.pdf",
        "isFolder": false,
        "size": 1024000,
        "type": "file",
        "extension": "pdf",
        "mimeType": "application/pdf",
        "createdAt": "2024-01-16T09:00:00Z",
        "modifiedAt": "2024-01-16T09:00:00Z",
        "owner": "John Doe",
        "isLink": false,
        "url": null
      },
      {
        "id": 3,
        "name": "Google Docs Link",
        "isFolder": false,
        "type": "link",
        "extension": "LINK",
        "createdAt": "2024-01-16T12:00:00Z",
        "modifiedAt": "2024-01-16T12:00:00Z",
        "owner": "John Doe",
        "isLink": true,
        "url": "https://docs.google.com/document/d/example"
      }
    ]
  }
}
```

#### **POST /files/upload**
**Location:** [src/stores/files.js](cci:7://file:///d:/CPNS%202025/latsar/aktulisasi/frontend-file-management/web-ai/src/stores/files.js:0:0-0:0) - [uploadFiles()](cci:1://file:///d:/CPNS%202025/latsar/aktulisasi/frontend-file-management/web-ai/src/stores/files.js:157:4-198:5) method
```javascript
// You'll need to implement this endpoint (currently fully mocked)
```

**Request:** `multipart/form-data`
- `files[]`: Array of files
- `path`: Target directory path
- `overwrite`: Boolean (optional)

**Response:**
```json
{
  "success": true,
  "data": {
    "uploadedFiles": [
      {
        "id": 4,
        "name": "uploaded-file.jpg",
        "isFolder": false,
        "size": 2048000,
        "type": "file",
        "extension": "jpg",
        "mimeType": "image/jpeg",
        "createdAt": "2024-01-20T10:00:00Z",
        "modifiedAt": "2024-01-20T10:00:00Z",
        "owner": "Current User",
        "isLink": false
      }
    ]
  }
}
```

#### **POST /links**
**Location:** [src/stores/files.js](cci:7://file:///d:/CPNS%202025/latsar/aktulisasi/frontend-file-management/web-ai/src/stores/files.js:0:0-0:0) - [uploadLink()](cci:1://file:///d:/CPNS%202025/latsar/aktulisasi/frontend-file-management/web-ai/src/stores/files.js:200:4-224:5) method
```javascript
// You'll need to implement this endpoint (currently fully mocked)
```

**Request Body:**
```json
{
  "url": "https://docs.google.com/document/d/example",
  "name": "My Google Doc",
  "path": "/current/path"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "link": {
      "id": 5,
      "name": "My Google Doc",
      "isFolder": false,
      "type": "link",
      "extension": "LINK",
      "createdAt": "2024-01-20T10:00:00Z",
      "modifiedAt": "2024-01-20T10:00:00Z",
      "owner": "Current User",
      "isLink": true,
      "url": "https://docs.google.com/document/d/example"
    }
  }
}
```

#### **POST /folders**
**Location:** [src/stores/files.js](cci:7://file:///d:/CPNS%202025/latsar/aktulisasi/frontend-file-management/web-ai/src/stores/files.js:0:0-0:0) - [createFolder()](cci:1://file:///d:/CPNS%202025/latsar/aktulisasi/frontend-file-management/web-ai/src/stores/files.js:129:4-155:5) method
```javascript
// Current mock URL: https://api.example.com/folders
```

**Request Body:**
```json
{
  "name": "New Folder",
  "path": "/current/path"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "folder": {
      "id": 6,
      "name": "New Folder",
      "isFolder": true,
      "path": "/current/path/New Folder",
      "createdAt": "2024-01-20T10:00:00Z",
      "modifiedAt": "2024-01-20T10:00:00Z",
      "owner": "Current User"
    }
  }
}
```

#### **DELETE /items**
**Location:** [src/stores/files.js](cci:7://file:///d:/CPNS%202025/latsar/aktulisasi/frontend-file-management/web-ai/src/stores/files.js:0:0-0:0) - [deleteItems()](cci:1://file:///d:/CPNS%202025/latsar/aktulisasi/frontend-file-management/web-ai/src/stores/files.js:226:4-250:5) method
```javascript
// Current mock URL: https://api.example.com/items
```

**Request Body:**
```json
{
  "items": [
    {
      "id": 1,
      "isFolder": false
    },
    {
      "id": 2,
      "isFolder": true
    }
  ]
}
```

**Response:**
```json
{
  "success": true,
  "message": "Items deleted successfully",
  "data": {
    "deletedCount": 2
  }
}
```

---

### **4. Statistics API**

#### **GET /stats**
**Location:** [src/stores/files.js](cci:7://file:///d:/CPNS%202025/latsar/aktulisasi/frontend-file-management/web-ai/src/stores/files.js:0:0-0:0) - [fetchStats()](cci:1://file:///d:/CPNS%202025/latsar/aktulisasi/frontend-file-management/web-ai/src/stores/files.js:104:4-127:5) method
```javascript
// Current mock URL: https://api.example.com/stats
```

**Response:**
```json
{
  "success": true,
  "data": {
    "totalFolders": 15,
    "totalFiles": 127,
    "totalStorage": 2400000000,
    "usedStorage": 1800000000,
    "filesByExtension": {
      "PDF": 45,
      "DOCX": 23,
      "MP4": 12,
      "JPG": 31,
      "PNG": 16
    },
    "uploadActivity": [
      {
        "date": "2024-01-15",
        "uploads": 5,
        "downloads": 12
      },
      {
        "date": "2024-01-16",
        "uploads": 8,
        "downloads": 15
      }
    ]
  }
}
```

---

### **5. User Management APIs (Admin Only)**

#### **GET /admin/users**
**Location:** [src/views/UserManagement.vue](cci:7://file:///d:/CPNS%202025/latsar/aktulisasi/frontend-file-management/web-ai/src/views/UserManagement.vue:0:0-0:0) - [loadUsers()](cci:1://file:///d:/CPNS%202025/latsar/aktulisasi/frontend-file-management/web-ai/src/views/UserManagement.vue:241:4-275:5) method
```javascript
// You'll need to implement this endpoint (currently fully mocked)
```

**Response:**
```json
{
  "success": true,
  "data": {
    "users": [
      {
        "id": 1,
        "name": "John Doe",
        "username": "johndoe",
        "email": "john@example.com",
        "role": "user",
        "status": "active",
        "createdAt": "2024-01-15T10:00:00Z",
        "lastActiveAt": "2024-01-19T14:20:00Z"
      }
    ]
  }
}
```

#### **POST /admin/users**
**Location:** [src/views/UserManagement.vue](cci:7://file:///d:/CPNS%202025/latsar/aktulisasi/frontend-file-management/web-ai/src/views/UserManagement.vue:0:0-0:0) - [createUser()](cci:1://file:///d:/CPNS%202025/latsar/aktulisasi/frontend-file-management/web-ai/src/views/UserManagement.vue:277:4-303:5) method
```javascript
// You'll need to implement this endpoint (currently fully mocked)
```

**Request Body:**
```json
{
  "name": "New User",
  "username": "newuser",
  "email": "newuser@example.com",
  "password": "password123",
  "role": "user"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "user": {
      "id": 7,
      "name": "New User",
      "username": "newuser",
      "email": "newuser@example.com",
      "role": "user",
      "status": "active",
      "createdAt": "2024-01-20T10:00:00Z",
      "lastActiveAt": null
    }
  }
}
```

#### **PUT /admin/users/:id/role**
**Location:** [src/views/UserManagement.vue](cci:7://file:///d:/CPNS%202025/latsar/aktulisasi/frontend-file-management/web-ai/src/views/UserManagement.vue:0:0-0:0) - [promoteUser()](cci:1://file:///d:/CPNS%202025/latsar/aktulisasi/frontend-file-management/web-ai/src/views/UserManagement.vue:305:4-311:5) / [demoteUser()](cci:1://file:///d:/CPNS%202025/latsar/aktulisasi/frontend-file-management/web-ai/src/views/UserManagement.vue:313:4-319:5) methods
```javascript
// You'll need to implement this endpoint (currently fully mocked)
```

**Request Body:**
```json
{
  "role": "admin" // or "user"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "user": {
      "id": 1,
      "role": "admin"
    }
  }
}
```

---

## 🔧 Implementation Steps

### **1. Update API Base URL**
Create an environment configuration file:

```javascript
// src/config/api.js
const API_BASE_URL = process.env.VUE_APP_API_BASE_URL || 'http://localhost:3000/api'

export default {
  baseURL: API_BASE_URL,
  timeout: 10000
}
```

### **2. Update Axios Configuration**
```javascript
// src/main.js or src/config/axios.js
import axios from 'axios'
import apiConfig from './config/api'

axios.defaults.baseURL = apiConfig.baseURL
axios.defaults.timeout = apiConfig.timeout

// Add request interceptor for auth token
axios.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// Add response interceptor for error handling
axios.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      // Handle unauthorized access
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)
```

### **3. Replace Mock URLs**
Update all the mock URLs in your stores:

**In [src/stores/auth.js](cci:7://file:///d:/CPNS%202025/latsar/aktulisasi/frontend-file-management/web-ai/src/stores/auth.js:0:0-0:0):**
```javascript
// Replace:
const response = await axios.post('https://api.example.com/auth/login', credentials)
// With:
const response = await axios.post('/auth/login', credentials)
```

**In [src/stores/files.js](cci:7://file:///d:/CPNS%202025/latsar/aktulisasi/frontend-file-management/web-ai/src/stores/files.js:0:0-0:0):**
```javascript
// Replace all mock URLs with relative paths:
'/files'
'/files/upload'
'/links'
'/folders'
'/items'
'/stats'
```

### **4. Error Handling**
Ensure all API calls handle errors properly:

```javascript
try {
  const response = await axios.post('/api/endpoint', data)
  return { success: true, data: response.data }
} catch (error) {
  console.error('API Error:', error)
  return { 
    success: false, 
    message: error.response?.data?.message || 'Operation failed' 
  }
}
```

### **5. Authentication Headers**
Make sure all protected endpoints receive the JWT token:

```javascript
// This is already implemented in your auth store
axios.defaults.headers.common['Authorization'] = `Bearer ${token}`
```

---

## 🎯 Priority Implementation Order

1. **Authentication APIs** - Essential for login/logout
2. **File Management APIs** - Core functionality
3. **Statistics API** - For dashboard
4. **User Management APIs** - Admin features
5. **Profile APIs** - User settings

This guide provides everything you need to implement a fully functional backend for your Vue.js file archiving system. Each endpoint includes the exact request/response format expected by your frontend code.