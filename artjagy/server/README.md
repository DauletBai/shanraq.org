# Artjagy Server - Сервер Басқарушылары
# Artjagy Server - Server Basqaru

## 📁 **Файлдар құрылымы / File Structure**

```
artjagy/server/
├── jojj_basqaru.tng              # JOJJ API басқарушылары
├── main.js                       # Негізгі сервер файлы
├── tenge_server.tng              # Tenge сервер
├── archetype_engine.tng          # Архетип қозғалтқышы
├── morpheme_engine.tng           # Морфема қозғалтқышы
├── phoneme_engine.tng            # Фонема қозғалтқышы
├── performance_optimization.tng  # Өнімділік оптимизациясы
├── simd_processor.tng            # SIMD процессор
└── README.md                     # Бұл файл
```

## 🚀 **JOJJ Basqaru (Басқарушылары)**

### **jojj_basqaru.tng**
JOJJ API басқарушылары - барлық CRUD операцияларын басқаратын негізгі файл.

**Мүмкіндіктер:**
- Пайдаланушылар басқаруы (Paydalanusylar)
- Мақалалар басқаруы (Maqalalar)  
- Санаттар басқаруы (Sanattar)
- API эндпоинттері
- Валидация және қателерді өңдеу

**Негізгі функциялар:**
```tenge
// Инициализация
jojj_basqaru_initialize() -> aqıqat

// Пайдаланушылар API
api_paydalanusylar_oqu_barlik()
api_paydalanu_oqu()
api_paydalanu_jasau()
api_paydalanu_janartu()
api_paydalanu_joiu()

// Мақалалар API
api_maqalalar_oqu_barlik()
api_maqala_oqu()
api_maqala_jasau()
api_maqala_janartu()
api_maqala_joiu()

// Санаттар API
api_sanattar_oqu_barlik()
api_sanat_oqu()
api_sanat_jasau()
api_sanat_janartu()
api_sanat_joiu()
```

## 🔧 **Негізгі сервер файлы**

### **main.js**
Express.js негізіндегі негізгі сервер файлы.

**Мүмкіндіктер:**
- HTTP сервер
- Middleware конфигурациясы
- JOJJ API маршруттары
- Статикалық файлдар
- Қателерді өңдеу

**API маршруттары:**
```javascript
// Пайдаланушылар
GET    /api/v1/paydalanusylar
GET    /api/v1/paydalanusylar/:id
POST   /api/v1/paydalanusylar
PUT    /api/v1/paydalanusylar/:id
DELETE /api/v1/paydalanusylar/:id

// Мақалалар
GET    /api/v1/maqalalar
GET    /api/v1/maqalalar/:id
POST   /api/v1/maqalalar
PUT    /api/v1/maqalalar/:id
DELETE /api/v1/maqalalar/:id

// Санаттар
GET    /api/v1/sanattar
GET    /api/v1/sanattar/:id
POST   /api/v1/sanattar
PUT    /api/v1/sanattar/:id
DELETE /api/v1/sanattar/:id
```

## ⚙️ **Сервер компоненттері**

### **tenge_server.tng**
Shanraq.org негізгі сервері - агглютинативтік ерекшеліктерді пайдаланатын жоғары өнімділікті сервер.

### **archetype_engine.tng**
Архетип қозғалтқышы - паттерн-негізделген дамыту жүйесі.

### **morpheme_engine.tng**
Морфема қозғалтқышы - қазақ тілінің морфемаларын өңдейтін қозғалтқыш.

### **phoneme_engine.tng**
Фонема қозғалтқышы - дыбыстық ерекшеліктерді өңдейтін қозғалтқыш.

### **performance_optimization.tng**
Өнімділік оптимизациясы - сервердің жылдамдығын арттыратын функциялар.

### **simd_processor.tng**
SIMD процессор - векторлық операцияларды орындайтын процессор.

## 🚀 **Серверді іске қосу**

### **Негізгі сервер:**
```bash
npm start
# немесе
node index.js
```

### **Тікелей сервер:**
```bash
npm run server
# немесе
node artjagy/server/main.js
```

### **Демо сервер:**
```bash
npm run demo
# немесе
python3 synaqtar/demo/template_server.py
```

## 📊 **API тестілеу**

### **Health Check:**
```bash
curl http://localhost:8080/api/v1/health
```

### **Status:**
```bash
curl http://localhost:8080/api/v1/status
```

### **Statistics:**
```bash
curl http://localhost:8080/api/v1/statistics
```

## 🔧 **Конфигурация**

Сервер конфигурациясы `baptaular/server_baptaular.json` файлында орналасқан:

```json
{
  "server": {
    "port": 8080,
    "host": "localhost",
    "environment": "development"
  },
  "database": {
    "type": "sqlite",
    "database": "tenge_web.db"
  },
  "security": {
    "jwt_secret": "tenge_web_secret_key",
    "bcrypt_rounds": 12
  }
}
```

## 📚 **Қосымша ақпарат**

- [JOJJ API Документациясы](../../qujattama/api/jojj_api.md)
- [Архитектура Документациясы](../../qujattama/architecture/overview.md)
- [Пайдаланушы Нұсқаулығы](../../qujattama/user-guide/getting-started.md)

---

**Shanraq.org Artjagy Server** - қазақ тілінің агглютинативтік күшімен жасалған заманауи сервер басқарушылары! 🚀
