# JOJJ API Документациясы
# JOJJ API Documentation

## 🎯 **JOJJ Интерфейс**

JOJJ (Jasau, Oqu, Janartu, Joiu) - бұл Shanraq.org жобасының негізгі CRUD операцияларының латин казах тіліндегі қысқартылмасы.

**JOJJ = Jasau (Жасау) + Oqu (Оқу) + Janartu (Жаңарту) + Joiu (Жою)**

## 📋 **API Эндпоинттері**

### 👥 **Пайдаланушылар (Users) - /api/v1/paydalanusylar**

#### **GET /api/v1/paydalanusylar**
Барлық пайдаланушыларды алу
```bash
curl -X GET http://localhost:8080/api/v1/paydalanusylar
```

**Жауап:**
```json
{
  "success": true,
  "users": [
    {
      "id": "user123",
      "name": "Test User",
      "email": "test@example.com",
      "role": "user",
      "status": "active",
      "created_at": 1640995200,
      "updated_at": 1640995200
    }
  ],
  "count": 1
}
```

#### **GET /api/v1/paydalanusylar/:id**
Пайдаланушыны ID бойынша алу
```bash
curl -X GET http://localhost:8080/api/v1/paydalanusylar/user123
```

#### **POST /api/v1/paydalanusylar**
Жаңа пайдаланушы жасау
```bash
curl -X POST http://localhost:8080/api/v1/paydalanusylar \
  -H "Content-Type: application/json" \
  -d '{
    "name": "New User",
    "email": "newuser@example.com",
    "password": "password123",
    "role": "user"
  }'
```

#### **PUT /api/v1/paydalanusylar/:id**
Пайдаланушыны жаңарту
```bash
curl -X PUT http://localhost:8080/api/v1/paydalanusylar/user123 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Updated User",
    "email": "updated@example.com",
    "role": "admin"
  }'
```

#### **DELETE /api/v1/paydalanusylar/:id**
Пайдаланушыны жою
```bash
curl -X DELETE http://localhost:8080/api/v1/paydalanusylar/user123
```

### 📝 **Мақалалар (Articles) - /api/v1/maqalalar**

#### **GET /api/v1/maqalalar**
Барлық мақалаларды алу
```bash
curl -X GET http://localhost:8080/api/v1/maqalalar
```

#### **GET /api/v1/maqalalar/:id**
Мақаланы ID бойынша алу
```bash
curl -X GET http://localhost:8080/api/v1/maqalalar/article123
```

#### **POST /api/v1/maqalalar**
Жаңа мақала жасау
```bash
curl -X POST http://localhost:8080/api/v1/maqalalar \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Test Article",
    "content": "This is a test article content",
    "author_id": "author123",
    "category_id": "category123",
    "status": "published"
  }'
```

#### **PUT /api/v1/maqalalar/:id**
Мақаланы жаңарту
```bash
curl -X PUT http://localhost:8080/api/v1/maqalalar/article123 \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Updated Article",
    "content": "Updated content",
    "status": "draft"
  }'
```

#### **DELETE /api/v1/maqalalar/:id**
Мақаланы жою
```bash
curl -X DELETE http://localhost:8080/api/v1/maqalalar/article123
```

### 🏷️ **Санаттар (Categories) - /api/v1/sanattar**

#### **GET /api/v1/sanattar**
Барлық санаттарды алу
```bash
curl -X GET http://localhost:8080/api/v1/sanattar
```

#### **GET /api/v1/sanattar/:id**
Санатты ID бойынша алу
```bash
curl -X GET http://localhost:8080/api/v1/sanattar/category123
```

#### **POST /api/v1/sanattar**
Жаңа санат жасау
```bash
curl -X POST http://localhost:8080/api/v1/sanattar \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Technology",
    "description": "Technology related articles",
    "parent_id": "",
    "color": "#ff0000",
    "icon": "tech-icon"
  }'
```

#### **PUT /api/v1/sanattar/:id**
Санатты жаңарту
```bash
curl -X PUT http://localhost:8080/api/v1/sanattar/category123 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Updated Category",
    "description": "Updated description",
    "color": "#00ff00"
  }'
```

#### **DELETE /api/v1/sanattar/:id**
Санатты жою
```bash
curl -X DELETE http://localhost:8080/api/v1/sanattar/category123
```

## 🔍 **Дополнительные эндпоинты**

### **GET /api/v1/maqalalar/popular**
Популярные мақалалар
```bash
curl -X GET http://localhost:8080/api/v1/maqalalar/popular
```

### **GET /api/v1/maqalalar/recent**
Соңғы мақалалар
```bash
curl -X GET http://localhost:8080/api/v1/maqalalar/recent
```

### **GET /api/v1/sanattar/hierarchy**
Санаттар иерархиясы
```bash
curl -X GET http://localhost:8080/api/v1/sanattar/hierarchy
```

### **GET /api/v1/statistics**
Жалпы статистика
```bash
curl -X GET http://localhost:8080/api/v1/statistics
```

## 📊 **Жауап форматы**

### **Сәтті жауап:**
```json
{
  "success": true,
  "data": { ... },
  "message": "Operation completed successfully"
}
```

### **Қате жауап:**
```json
{
  "success": false,
  "error": "Error message",
  "validation_errors": [ ... ]
}
```

## 🔐 **Аутентификация**

Кейбір эндпоинттер аутентификация талап етеді:

```bash
curl -X GET http://localhost:8080/api/v1/paydalanusylar \
  -H "Authorization: Bearer your-jwt-token"
```

## 📝 **Валидация**

### **Пайдаланушы валидациясы:**
- `name`: кемінде 2 символ
- `email`: жарамды email форматы
- `password`: кемінде 8 символ

### **Мақала валидациясы:**
- `title`: кемінде 3 символ
- `content`: кемінде 10 символ
- `author_id`: жарамды автор ID

### **Санат валидациясы:**
- `name`: кемінде 3 символ
- `description`: кемінде 10 символ
- `parent_id`: жарамды ата-ана санат ID

## 🚀 **Мысал пайдалану**

### **1. Пайдаланушы жасау:**
```bash
curl -X POST http://localhost:8080/api/v1/paydalanusylar \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Айдар",
    "email": "aydar@example.com",
    "password": "password123",
    "role": "user"
  }'
```

### **2. Мақала жасау:**
```bash
curl -X POST http://localhost:8080/api/v1/maqalalar \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Shanraq.org туралы",
    "content": "Shanraq.org - қазақ тілінің агглютинативтік ерекшеліктерін пайдаланатын веб-фреймворк",
    "author_id": "aydar123",
    "category_id": "tech123",
    "status": "published"
  }'
```

### **3. Санат жасау:**
```bash
curl -X POST http://localhost:8080/api/v1/sanattar \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Технология",
    "description": "Технология саласындағы мақалалар",
    "parent_id": "",
    "color": "#007bff",
    "icon": "tech"
  }'
```

## 🎯 **JOJJ Принциптері**

1. **Jasau (Жасау)** - жаңа ресурс жасау
2. **Oqu (Оқу)** - ресурсты алу немесе іздеу
3. **Janartu (Жаңарту)** - ресурсты жаңарту
4. **Joiu (Жою)** - ресурсты жою

Бұл принциптер барлық ресурстар үшін бірдей қолданылады және Shanraq.org жобасының негізгі архитектурасын құрайды.

---

**Shanraq.org JOJJ API** - қазақ тілінің агглютинативтік күшімен жасалған заманауи API! 🚀
