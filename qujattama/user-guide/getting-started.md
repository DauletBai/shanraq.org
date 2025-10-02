# Shanraq.org Getting Started Guide

## Introduction

Welcome to Shanraq.org, an agglutinative web application framework built on the Tenge programming language. This guide will help you get started with building modern web applications using Kazakh language's agglutinative features.

## What is Shanraq.org?

Shanraq.org is a revolutionary web framework that leverages the agglutinative nature of the Kazakh language to create more natural and intuitive programming syntax. It combines:

- **Agglutinative Programming**: Use Kazakh morphemes and phonemes for function composition
- **Modern Web Technologies**: Built on proven web standards
- **High Performance**: Optimized for speed and efficiency
- **Natural Syntax**: Code that reads like natural language

## Prerequisites

Before you begin, ensure you have the following installed:

- **Go 1.21+**: For the Tenge compiler and runtime
- **Node.js 18+**: For frontend development
- **Git**: For version control
- **SQLite3**: For database operations (or PostgreSQL/MySQL)

## Installation

### 1. Clone the Repository

```bash
git clone https://github.com/shanraq/shanraq.git
cd shanraq
```

### 2. Install Dependencies

```bash
# Install Go dependencies
go mod download

# Install Node.js dependencies
npm install
```

### 3. Build the Project

```bash
# Build all components
make build

# Or build individually
make build-backend
make build-frontend
```

### 4. Run Database Migrations

```bash
make migrate-up
```

### 5. Start the Development Server

```bash
# Start all services
make run-web

# Or start individually
make run-api      # API server
make run-frontend # Frontend server
```

## Your First Shanraq.org Application

Let's create a simple "Hello World" application to get familiar with Shanraq.org.

### 1. Create a Basic Server

Create a file called `main.tng`:

```tenge
// main.tng - Your first Shanraq.org application
atqar main() {
    korset("Shanraq.org серверін іске қосу...");
    
    // Create web server
    jasau server: WebServer = web_server_create(8080);
    
    // Add routes
    web_get_route_qosu(server, "/", ana_sahifa_handler);
    web_get_route_qosu(server, "/api/hello", api_hello_handler);
    
    // Start server
    korset("Сервер іске қосылды: http://localhost:8080");
    web_server_listen(server);
}

// Homepage handler
atqar ana_sahifa_handler(request: WebRequest, response: WebResponse) {
    jasau html: jol = `
<!DOCTYPE html>
<html lang="kk">
<head>
    <meta charset="UTF-8">
    <title>Shanraq.org</title>
</head>
<body>
    <h1>Қош келдіңіз, Shanraq.org!</h1>
    <p>Бұл сіздің бірінші Shanraq.org қосымшасы.</p>
</body>
</html>`;
    
    web_html_response_qaytar(response, html, 200);
}

// API handler
atqar api_hello_handler(request: WebRequest, response: WebResponse) {
    jasau data: JsonObject = json_object_create();
    json_object_set_string(data, "message", "Сәлем, Shanraq.org!");
    json_object_set_string(data, "timestamp", current_timestamp());
    
    web_json_response_qaytar(response, data, 200);
}
```

### 2. Compile and Run

```bash
# Compile the Tenge code to C
tenge compile main.tng

# Compile the generated C code
gcc -o main main.c -lm

# Run the application
./main
```

### 3. Test Your Application

Open your browser and visit:
- http://localhost:8080 - Homepage
- http://localhost:8080/api/hello - API endpoint

## Understanding Tenge Syntax

### Agglutinative Function Names

Shanraq.org uses agglutinative function names based on Kazakh morphemes:

```tenge
// Basic morphemes
jasau    // create
alu      // get/retrieve
qosu     // add
zhangartu // update
zhoyu    // delete
tekseru  // check/validate

// Function composition
web_server_jasau(port)           // create web server
user_parol_tekseru(email, pass)  // check user password
database_connection_opt()        // optimize database connection
```

### Variable Declarations

```tenge
// Basic variable declaration
jasau name: jol = "Shanraq.org";
jasau port: san = 8080;
jasau is_running: aqıqat = jan;

// Complex types
jasau user: JsonObject = json_object_create();
jasau server: WebServer = web_server_create(port);
```

### Control Structures

```tenge
// Conditional statements
eger (condition) {
    // do something
} aitpese {
    // do something else
}

// Loops
azirshe (i < 10) {
    korset("Iteration: " + int_to_string(i));
    i = i + 1;
}
```

### Function Definitions

```tenge
// Function with parameters and return type
atqar user_jasau(name: jol, email: jol) -> JsonObject {
    jasau user: JsonObject = json_object_create();
    json_object_set_string(user, "name", name);
    json_object_set_string(user, "email", email);
    qaytar user;
}
```

## Building a Complete Application

### 1. User Management System

Let's build a simple user management system:

```tenge
// user_management.tng
atqar user_management_jasau() -> UserManagement {
    jasau mgmt: UserManagement = user_management_create();
    qaytar mgmt;
}

atqar user_tirkelu_jasau(name: jol, email: jol, password: jol) -> JsonObject {
    // Validate input
    eger (name == "" || email == "" || password == "") {
        jasau error: JsonObject = json_object_create();
        json_object_set_string(error, "error", "All fields are required");
        qaytar error;
    }
    
    // Create user
    jasau user: JsonObject = json_object_create();
    json_object_set_string(user, "id", uuid_generate());
    json_object_set_string(user, "name", name);
    json_object_set_string(user, "email", email);
    json_object_set_string(user, "password", hash_password(password));
    json_object_set_string(user, "status", "active");
    json_object_set_number(user, "created_at", current_timestamp());
    
    // Save to database
    jasau save_result: aqıqat = user_derekterna_saktau(user);
    
    eger (save_result) {
        qaytar user;
    } aitpese {
        qaytar NULL;
    }
}

atqar user_kimdik_tekseru(email: jol, password: jol) -> JsonObject {
    jasau user: JsonObject = user_email_boyynsha_tabu(email);
    
    eger (user == NULL) {
        qaytar NULL;
    }
    
    jasau stored_password: jol = json_object_get_string(user, "password");
    jasau password_valid: aqıqat = parol_tekseru(password, stored_password);
    
    eger (password_valid) {
        qaytar user;
    } aitpese {
        qaytar NULL;
    }
}
```

### 2. Web Routes

```tenge
// routes.tng
atqar web_route_qosu(server: WebServer) {
    // User routes
    web_get_route_qosu(server, "/users", user_listesi_handler);
    web_post_route_qosu(server, "/users", user_jasau_handler);
    web_get_route_qosu(server, "/users/:id", user_alu_handler);
    web_put_route_qosu(server, "/users/:id", user_zhangartu_handler);
    web_delete_route_qosu(server, "/users/:id", user_zhoyu_handler);
    
    // Authentication routes
    web_post_route_qosu(server, "/login", login_handler);
    web_post_route_qosu(server, "/register", register_handler);
    web_post_route_qosu(server, "/logout", logout_handler);
}

atqar user_listesi_handler(request: WebRequest, response: WebResponse) {
    jasau limit: san = 10;
    jasau offset: san = 0;
    
    jasau users: JsonObject[] = user_tizimi(limit, offset);
    jasau total: san = user_sany();
    
    jasau result: JsonObject = json_object_create();
    json_object_set_array(result, "users", users);
    json_object_set_number(result, "total", total);
    
    web_json_response_qaytar(response, result, 200);
}
```

### 3. Database Models

```tenge
// models.tng
atqar user_model_jasau() -> Model {
    jasau model: Model = model_create("users");
    
    jasau id_field: Field = model_field_qosu(model, "id", "VARCHAR(36)", [constraint_primary_key()]);
    jasau name_field: Field = model_field_qosu(model, "name", "VARCHAR(100)", [constraint_not_null()]);
    jasau email_field: Field = model_field_qosu(model, "email", "VARCHAR(255)", [constraint_not_null(), constraint_unique()]);
    jasau password_field: Field = model_field_qosu(model, "password", "VARCHAR(255)", [constraint_not_null()]);
    jasau status_field: Field = model_field_qosu(model, "status", "VARCHAR(20)", [constraint_not_null()]);
    jasau created_at_field: Field = model_field_qosu(model, "created_at", "TIMESTAMP", [constraint_not_null()]);
    
    qaytar model;
}
```

## Frontend Development

### 1. Component System

Shanraq.org includes a component-based frontend system:

```tenge
// components/header.tng
atqar header_component_jasau() -> Component {
    jasau component: Component = component_create("header");
    
    jasau html: jol = `
<header class="header">
    <div class="container">
        <div class="logo">
            <a href="/">Shanraq.org</a>
        </div>
        <nav class="navigation">
            <a href="/">Басты бет</a>
            <a href="/about">Біз жөнінде</a>
            <a href="/contact">Байланыс</a>
        </nav>
    </div>
</header>`;
    
    jasau css: jol = `
.header {
    background: #fff;
    box-shadow: 0 2px 4px rgba(0,0,0,0.1);
    padding: 1rem 0;
}

.container {
    max-width: 1200px;
    margin: 0 auto;
    display: flex;
    justify-content: space-between;
    align-items: center;
}

.logo a {
    font-size: 1.5rem;
    font-weight: bold;
    text-decoration: none;
    color: #333;
}

.navigation a {
    margin-left: 2rem;
    text-decoration: none;
    color: #666;
}
`;
    
    component_set_html(component, html);
    component_set_css(component, css);
    
    qaytar component;
}
```

### 2. Page Components

```tenge
// pages/home.tng
atqar home_page_jasau() -> Component {
    jasau component: Component = component_create("home_page");
    
    jasau html: jol = `
<!DOCTYPE html>
<html lang="kk">
<head>
    <meta charset="UTF-8">
    <title>Shanraq.org - Агглютинативтік Веб-Қосымша</title>
    <link rel="stylesheet" href="/static/css/main.css">
</head>
<body>
    <div id="header-container"></div>
    
    <main class="main-content">
        <section class="hero">
            <h1>Shanraq.org</h1>
            <p>Қазақ тілінің агглютинативтік ерекшеліктерін пайдаланатын заманауи веб-қосымша</p>
            <a href="/docs" class="btn btn-primary">Бастау</a>
        </section>
        
        <section class="features">
            <div class="feature">
                <h3>Морфемалар</h3>
                <p>Қазақ тілінің морфемаларын пайдалана отырып, табиғи функция атауларын жасау</p>
            </div>
            <div class="feature">
                <h3>Фонемалар</h3>
                <p>Фонемаларды пайдалана отырып, дыбыстық ерекшеліктерді кодта көрсету</p>
            </div>
            <div class="feature">
                <h3>Архетиптер</h3>
                <p>Архетиптер жүйесімен күрделі архитектураларды қарапайым түрде құру</p>
            </div>
        </section>
    </main>
    
    <div id="footer-container"></div>
    
    <script src="/static/js/main.js"></script>
</body>
</html>`;
    
    component_set_html(component, html);
    qaytar component;
}
```

## Database Operations

### 1. Using the ORM

```tenge
// database_operations.tng
atqar user_derekterna_saktau(user: JsonObject) -> aqıqat {
    jasau model: Model = user_model_jasau();
    jasau result: JsonObject = model_create_record(model, user);
    qaytar result != NULL;
}

atqar user_email_boyynsha_tabu(email: jol) -> JsonObject {
    jasau model: Model = user_model_jasau();
    jasau conditions: JsonObject = json_object_create();
    json_object_set_string(conditions, "email", email);
    
    jasau users: JsonObject[] = model_find_records(model, conditions, 1, 0);
    
    eger (users.length > 0) {
        qaytar users[0];
    } aitpese {
        qaytar NULL;
    }
}

atqar user_tizimi(limit: san, offset: san) -> JsonObject[] {
    jasau model: Model = user_model_jasau();
    jasau conditions: JsonObject = json_object_create();
    qaytar model_find_records(model, conditions, limit, offset);
}
```

### 2. Migrations

```tenge
// migrations/001_create_users.tng
atqar migration_001_create_users() -> Migration {
    jasau migration: Migration = migration_jasau("create_users_table");
    
    jasau id_field: Field = field_qosu("id", "VARCHAR(36)", [constraint_primary_key()]);
    jasau name_field: Field = field_qosu("name", "VARCHAR(100)", [constraint_not_null()]);
    jasau email_field: Field = field_qosu("email", "VARCHAR(255)", [constraint_not_null(), constraint_unique()]);
    jasau password_field: Field = field_qosu("password", "VARCHAR(255)", [constraint_not_null()]);
    jasau status_field: Field = field_qosu("status", "VARCHAR(20)", [constraint_not_null()]);
    jasau created_at_field: Field = field_qosu("created_at", "TIMESTAMP", [constraint_not_null()]);
    
    jasau fields: Field[] = [id_field, name_field, email_field, password_field, status_field, created_at_field];
    migration_table_jasau(migration, "users", fields);
    
    qaytar migration;
}
```

## Testing Your Application

### 1. Unit Tests

```tenge
// synaqtar/unit/user_test.tng
atqar user_test_jasau() {
    korset("=== User Tests ===");
    
    // Test user creation
    jasau user: JsonObject = user_tirkelu_jasau("Test User", "test@example.com", "password123");
    
    eger (user != NULL) {
        korset("✅ User creation test passed");
    } aitpese {
        korset("❌ User creation test failed");
    }
    
    // Test user authentication
    jasau auth_user: JsonObject = user_kimdik_tekseru("test@example.com", "password123");
    
    eger (auth_user != NULL) {
        korset("✅ User authentication test passed");
    } aitpese {
        korset("❌ User authentication test failed");
    }
}
```

### 2. Integration Tests

```tenge
// synaqtar/integration/api_test.tng
atqar api_integration_test() {
    korset("=== API Integration Tests ===");
    
    // Test user registration endpoint
    jasau register_data: JsonObject = json_object_create();
    json_object_set_string(register_data, "name", "Test User");
    json_object_set_string(register_data, "email", "test@example.com");
    json_object_set_string(register_data, "password", "password123");
    
    jasau response: WebResponse = web_request_simulate("POST", "/api/v1/users", register_data);
    
    eger (response.status == 201) {
        korset("✅ User registration API test passed");
    } aitpese {
        korset("❌ User registration API test failed");
    }
}
```

## Deployment

### 1. Production Build

```bash
# Build for production
make build

# Run database migrations
make migrate-up

# Start production server
./bin/shanraq-server
```

### 2. Docker Deployment

```dockerfile
# Dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY . .
RUN go mod download
RUN make build

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/bin/shanraq-server .
COPY --from=builder /app/frontend/dist ./static

EXPOSE 8080
CMD ["./shanraq-server"]
```

### 3. Environment Configuration

```bash
# .env
DATABASE_URL=postgres://user:password@localhost/tenge_web
REDIS_URL=redis://localhost:6379
JWT_SECRET=your-secret-key
PORT=8080
```

## Next Steps

Now that you have a basic understanding of Shanraq.org, you can:

1. **Explore the API Documentation**: Learn about all available endpoints
2. **Build Complex Applications**: Use the full framework features
3. **Contribute to the Project**: Help improve Shanraq.org
4. **Join the Community**: Connect with other developers

## Resources

- **Documentation**: https://docs.shanraq.org
- **GitHub Repository**: https://github.com/shanraq/shanraq
- **Discord Community**: https://discord.gg/shanraq
- **Examples**: https://github.com/shanraq/examples

## Support

If you need help:

1. Check the documentation
2. Search existing issues on GitHub
3. Ask questions on Discord
4. Create a new issue for bugs

Welcome to the Shanraq.org community! 🚀

