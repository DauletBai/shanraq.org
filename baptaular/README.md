# Baptaular - Конфигурация / Баптаулар

Бұл директория Shanraq.org жобасының конфигурация файлдарын қамтиды.

## Файлдар / Files

- `server_baptaular.json` - Сервер конфигурациясы
- `database_baptaular.json` - Деректер базасы конфигурациясы

## Қолдану / Usage

```bash
# Серверді іске қосу
npm start

# Деректер базасын конфигурациялау
node baptaular/database_baptaular.json
```

## Конфигурация параметрлері / Configuration Parameters

### Сервер конфигурациясы / Server Configuration
- `port` - Сервер порты
- `host` - Сервер хості
- `environment` - Орта (development/production)

### Деректер базасы конфигурациясы / Database Configuration
- `type` - Деректер базасы түрі (sqlite/postgres)
- `host` - Деректер базасы хості
- `port` - Деректер базасы порты
- `database` - Деректер базасы атауы
- `username` - Пайдаланушы атауы
- `password` - Пароль

