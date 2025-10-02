# Shanraq.org REST API Documentation

## Overview

The Shanraq.org REST API provides a comprehensive set of endpoints for building agglutinative web applications. All API endpoints follow RESTful conventions and use JSON for data exchange.

## Base URL

```
http://localhost:8080/api/v1
```

## Authentication

Most API endpoints require authentication using JWT tokens. Include the token in the Authorization header:

```
Authorization: Bearer <your-jwt-token>
```

## Response Format

All API responses follow this format:

```json
{
  "success": true,
  "data": { ... },
  "message": "Operation completed successfully",
  "timestamp": "2024-01-01T00:00:00Z"
}
```

Error responses:

```json
{
  "success": false,
  "error": "Error message",
  "code": "ERROR_CODE",
  "timestamp": "2024-01-01T00:00:00Z"
}
```

## Endpoints

### Health & Status

#### GET /health
Check API health status.

**Response:**
```json
{
  "status": "ok",
  "timestamp": "2024-01-01T00:00:00Z",
  "version": "1.0.0"
}
```

#### GET /status
Get detailed system status.

**Response:**
```json
{
  "server": "running",
  "database": "connected",
  "cache": "active",
  "uptime": 3600
}
```

### User Management

#### GET /users
Get list of users.

**Query Parameters:**
- `limit` (optional): Number of users to return (default: 10)
- `offset` (optional): Number of users to skip (default: 0)
- `role` (optional): Filter by user role
- `status` (optional): Filter by user status

**Response:**
```json
{
  "users": [
    {
      "id": "user-123",
      "name": "John Doe",
      "email": "john@example.com",
      "role": "user",
      "status": "active",
      "created_at": "2024-01-01T00:00:00Z"
    }
  ],
  "total": 100,
  "limit": 10,
  "offset": 0
}
```

#### POST /users
Create a new user.

**Request Body:**
```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "securepassword123",
  "role": "user"
}
```

**Response:**
```json
{
  "id": "user-123",
  "name": "John Doe",
  "email": "john@example.com",
  "role": "user",
  "status": "active",
  "created_at": "2024-01-01T00:00:00Z"
}
```

#### GET /users/{id}
Get user by ID.

**Response:**
```json
{
  "id": "user-123",
  "name": "John Doe",
  "email": "john@example.com",
  "role": "user",
  "status": "active",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

#### PUT /users/{id}
Update user.

**Request Body:**
```json
{
  "name": "John Smith",
  "email": "johnsmith@example.com",
  "role": "admin"
}
```

#### DELETE /users/{id}
Delete user.

**Response:**
```json
{
  "message": "User deleted successfully"
}
```

### Authentication

#### POST /login
Authenticate user.

**Request Body:**
```json
{
  "email": "john@example.com",
  "password": "securepassword123"
}
```

**Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "user-123",
    "name": "John Doe",
    "email": "john@example.com",
    "role": "user"
  }
}
```

#### POST /register
Register new user.

**Request Body:**
```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "securepassword123"
}
```

#### POST /logout
Logout user.

**Response:**
```json
{
  "message": "Logged out successfully"
}
```

#### GET /profile
Get current user profile.

**Headers:**
- `Authorization: Bearer <token>`

**Response:**
```json
{
  "id": "user-123",
  "name": "John Doe",
  "email": "john@example.com",
  "role": "user",
  "created_at": "2024-01-01T00:00:00Z"
}
```

### Content Management

#### GET /content
Get list of content.

**Query Parameters:**
- `limit` (optional): Number of items to return (default: 10)
- `offset` (optional): Number of items to skip (default: 0)
- `category` (optional): Filter by category
- `status` (optional): Filter by status
- `author_id` (optional): Filter by author

**Response:**
```json
{
  "contents": [
    {
      "id": "content-123",
      "title": "Sample Article",
      "slug": "sample-article",
      "content": "Article content...",
      "category": "news",
      "status": "published",
      "author_id": "user-123",
      "views": 150,
      "likes": 25,
      "created_at": "2024-01-01T00:00:00Z"
    }
  ],
  "total": 50,
  "limit": 10,
  "offset": 0
}
```

#### POST /content
Create new content.

**Request Body:**
```json
{
  "title": "New Article",
  "content": "Article content...",
  "category": "news",
  "status": "draft"
}
```

#### GET /content/{id}
Get content by ID.

#### PUT /content/{id}
Update content.

#### DELETE /content/{id}
Delete content.

#### POST /content/{id}/like
Like content.

#### DELETE /content/{id}/like
Unlike content.

### E-Commerce

#### GET /products
Get list of products.

**Query Parameters:**
- `limit` (optional): Number of products to return (default: 10)
- `offset` (optional): Number of products to skip (default: 0)
- `category` (optional): Filter by category
- `min_price` (optional): Minimum price filter
- `max_price` (optional): Maximum price filter
- `sort_by` (optional): Sort by field (price, name, created_at)

**Response:**
```json
{
  "products": [
    {
      "id": "product-123",
      "name": "Sample Product",
      "description": "Product description...",
      "price": 2999,
      "category": "electronics",
      "stock": 50,
      "seller_id": "user-123",
      "views": 200,
      "sales": 10,
      "created_at": "2024-01-01T00:00:00Z"
    }
  ],
  "total": 100,
  "limit": 10,
  "offset": 0
}
```

#### POST /products
Create new product.

**Request Body:**
```json
{
  "name": "New Product",
  "description": "Product description...",
  "price": 2999,
  "category": "electronics",
  "stock": 50
}
```

#### GET /products/{id}
Get product by ID.

#### PUT /products/{id}
Update product.

#### DELETE /products/{id}
Delete product.

### Shopping Cart

#### GET /cart
Get user's shopping cart.

**Headers:**
- `Authorization: Bearer <token>`

**Response:**
```json
{
  "id": "cart-123",
  "user_id": "user-123",
  "items": [
    {
      "product_id": "product-123",
      "quantity": 2,
      "price": 2999,
      "name": "Sample Product"
    }
  ],
  "total": 5998,
  "created_at": "2024-01-01T00:00:00Z"
}
```

#### POST /cart/items
Add item to cart.

**Request Body:**
```json
{
  "product_id": "product-123",
  "quantity": 2
}
```

#### DELETE /cart/items/{product_id}
Remove item from cart.

### Orders

#### POST /orders
Create new order.

**Request Body:**
```json
{
  "cart_id": "cart-123",
  "shipping_address": {
    "street": "123 Main St",
    "city": "Almaty",
    "postal_code": "050000",
    "country": "Kazakhstan"
  },
  "payment_method": "credit_card"
}
```

#### GET /orders
Get user's orders.

#### GET /orders/{id}
Get order by ID.

#### PUT /orders/{id}/status
Update order status.

**Request Body:**
```json
{
  "status": "shipped"
}
```

## Error Codes

| Code | Description |
|------|-------------|
| `VALIDATION_ERROR` | Request validation failed |
| `AUTHENTICATION_REQUIRED` | Authentication required |
| `INSUFFICIENT_PERMISSIONS` | User lacks required permissions |
| `RESOURCE_NOT_FOUND` | Requested resource not found |
| `DUPLICATE_RESOURCE` | Resource already exists |
| `RATE_LIMIT_EXCEEDED` | Rate limit exceeded |
| `INTERNAL_SERVER_ERROR` | Internal server error |

## Rate Limiting

API requests are rate limited to prevent abuse:

- **Authenticated users**: 1000 requests per hour
- **Unauthenticated users**: 100 requests per hour
- **Specific endpoints**: May have additional limits

Rate limit headers are included in responses:

```
X-RateLimit-Limit: 1000
X-RateLimit-Remaining: 999
X-RateLimit-Reset: 1640995200
```

## Pagination

List endpoints support pagination using `limit` and `offset` parameters:

```
GET /api/v1/users?limit=20&offset=40
```

Response includes pagination metadata:

```json
{
  "data": [...],
  "total": 100,
  "limit": 20,
  "offset": 40,
  "has_next": true,
  "has_prev": true
}
```

## Filtering and Sorting

Many endpoints support filtering and sorting:

```
GET /api/v1/products?category=electronics&min_price=1000&sort_by=price
```

## Search

Search endpoints are available for content and products:

```
GET /api/v1/content/search?q=keyword&category=news
GET /api/v1/products/search?q=keyword&category=electronics
```

## Webhooks

The API supports webhooks for real-time notifications:

- **User Registration**: `POST /webhooks/user-registered`
- **Order Created**: `POST /webhooks/order-created`
- **Content Published**: `POST /webhooks/content-published`

## SDKs and Libraries

Official SDKs are available for:

- **JavaScript/Node.js**: `npm install @shanraq/sdk`
- **Python**: `pip install shanraq-sdk`
- **Go**: `go get github.com/shanraq/sdk-go`

## Examples

### Complete User Registration Flow

```javascript
// 1. Register user
const registerResponse = await fetch('/api/v1/register', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    name: 'John Doe',
    email: 'john@example.com',
    password: 'securepassword123'
  })
});

const { token, user } = await registerResponse.json();

// 2. Use token for authenticated requests
const profileResponse = await fetch('/api/v1/profile', {
  headers: {
    'Authorization': `Bearer ${token}`
  }
});

const profile = await profileResponse.json();
```

### Content Management

```javascript
// Create content
const contentResponse = await fetch('/api/v1/content', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'Authorization': `Bearer ${token}`
  },
  body: JSON.stringify({
    title: 'My Article',
    content: 'Article content...',
    category: 'news',
    status: 'published'
  })
});

// Get content list
const contentsResponse = await fetch('/api/v1/content?limit=10&offset=0');
const { contents } = await contentsResponse.json();
```

### E-Commerce Operations

```javascript
// Add product to cart
const cartResponse = await fetch('/api/v1/cart/items', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'Authorization': `Bearer ${token}`
  },
  body: JSON.stringify({
    product_id: 'product-123',
    quantity: 2
  })
});

// Create order
const orderResponse = await fetch('/api/v1/orders', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'Authorization': `Bearer ${token}`
  },
  body: JSON.stringify({
    cart_id: 'cart-123',
    shipping_address: {
      street: '123 Main St',
      city: 'Almaty',
      postal_code: '050000',
      country: 'Kazakhstan'
    },
    payment_method: 'credit_card'
  })
});
```

## Support

For API support and questions:

- **Documentation**: https://docs.shanraq.org
- **GitHub**: https://github.com/shanraq/shanraq
- **Discord**: https://discord.gg/shanraq
- **Email**: api-support@shanraq.org

