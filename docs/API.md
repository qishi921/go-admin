# Admin System API Documentation

Base URL: `http://localhost:8080/api/v1`

Authentication: Bearer Token (JWT)

---

## Authentication

### Login
```
POST /auth/login
```

**Request Body:**
```json
{
  "username": "admin",
  "password": "password123"
}
```

**Response:**
```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "user": {
      "id": 1,
      "username": "admin",
      "email": "admin@example.com"
    }
  }
}
```

### Get Current User
```
GET /auth/me
Authorization: Bearer {token}
```

### Update Profile
```
PUT /auth/me
Authorization: Bearer {token}
```

**Request Body:**
```json
{
  "email": "new@example.com",
  "phone": "13800138000",
  "real_name": "Admin User"
}
```

### Change Password
```
PUT /auth/password
Authorization: Bearer {token}
```

**Request Body:**
```json
{
  "old_password": "oldpass",
  "new_password": "newpass123"
}
```

### Get User Permissions
```
GET /auth/permissions
Authorization: Bearer {token}
```

**Response:**
```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "permissions": ["user:manage", "role:manage"],
    "is_super_admin": false
  }
}
```

---

## Users

### List Users
```
GET /users?page=1&per_page=20&search=keyword
Authorization: Bearer {token}
```

### Get User
```
GET /users/:id
Authorization: Bearer {token}
```

### Create User
```
POST /users
Authorization: Bearer {token}
```

**Request Body:**
```json
{
  "username": "newuser",
  "password": "password123",
  "email": "user@example.com",
  "phone": "13800138000",
  "status": "active"
}
```

### Update User
```
PUT /users/:id
Authorization: Bearer {token}
```

### Delete User
```
DELETE /users/:id
Authorization: Bearer {token}
```

### Get User Roles
```
GET /users/:id/roles
Authorization: Bearer {token}
```

---

## Roles

### List Roles
```
GET /roles?page=1&per_page=20&search=keyword
Authorization: Bearer {token}
```

### Create Role
```
POST /roles
Authorization: Bearer {token}
```

**Request Body:**
```json
{
  "name": "Editor",
  "code": "editor",
  "description": "Content editor role"
}
```

### Get Role Users
```
GET /roles/:id/users
Authorization: Bearer {token}
```

### Assign User to Role
```
POST /roles/:id/users
Authorization: Bearer {token}
```

**Request Body:**
```json
{
  "user_id": 5
}
```

### Remove User from Role
```
DELETE /roles/:id/users/:userId
Authorization: Bearer {token}
```

### Get Role Permissions
```
GET /roles/:id/permissions
Authorization: Bearer {token}
```

### Assign Permission to Role
```
POST /roles/:id/permissions
Authorization: Bearer {token}
```

**Request Body:**
```json
{
  "permission_id": 3
}
```

### Remove Permission from Role
```
DELETE /roles/:id/permissions/:permissionId
Authorization: Bearer {token}
```

---

## Permissions

### List Permissions
```
GET /permissions?page=1&per_page=20&search=keyword
Authorization: Bearer {token}
```

### Create Permission
```
POST /permissions
Authorization: Bearer {token}
```

**Request Body:**
```json
{
  "name": "Article Create",
  "code": "article:create",
  "type": "action",
  "description": "Create new articles"
}
```

---

## Menus

### List Menus
```
GET /menus?page=1&per_page=20&search=keyword
Authorization: Bearer {token}
```

### Get Menu Tree
```
GET /menus/tree
Authorization: Bearer {token}
```

**Response:**
```json
{
  "code": 0,
  "message": "ok",
  "data": [
    {
      "id": 1,
      "name": "System",
      "path": "/system",
      "children": [
        {"id": 2, "name": "Users", "path": "/system/users"}
      ]
    }
  ]
}
```

### Create Menu
```
POST /menus
Authorization: Bearer {token}
```

**Request Body:**
```json
{
  "name": "Articles",
  "path": "/articles",
  "icon": "layui-icon-read",
  "sort_order": 10,
  "parent_id": null
}
```

---

## Operation Logs

### List Logs
```
GET /logs?page=1&per_page=20&search=keyword
Authorization: Bearer {token}
```

### Get Log Detail
```
GET /logs/:id
Authorization: Bearer {token}
```

---

## Error Responses

### 400 Bad Request
```json
{
  "code": 400,
  "message": "Invalid JSON"
}
```

### 401 Unauthorized
```json
{
  "code": 401,
  "message": "Unauthorized"
}
```

### 403 Forbidden
```json
{
  "code": 403,
  "message": "Permission denied: user:create"
}
```

### 422 Validation Error
```json
{
  "code": 422,
  "message": "Validation failed",
  "errors": {
    "email": ["The email field must be a valid email address"]
  }
}
```

### 429 Too Many Requests
```json
{
  "code": 429,
  "message": "登录尝试次数过多，请稍后再试"
}
```

---

## RBAC Permission Codes

| Code | Description |
|------|-------------|
| `user:manage` | View users |
| `user:create` | Create users |
| `user:edit` | Edit users |
| `user:delete` | Delete users |
| `role:manage` | View roles |
| `role:create` | Create roles |
| `role:edit` | Edit roles |
| `role:delete` | Delete roles |
| `menu:manage` | View menus |
| `permission:manage` | View permissions |
| `log:view` | View operation logs |
| `dashboard:view` | View dashboard |

Users with `super_admin` role bypass all permission checks.
