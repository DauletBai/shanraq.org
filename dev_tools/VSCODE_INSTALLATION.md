# Установка VS Code расширения для Tenge
# VS Code Extension Installation for Tenge

## 🚀 Быстрая установка / Quick Installation

### Метод 1: Автоматическая установка
```bash
# Перейти в директорию dev_tools
cd dev_tools

# Установить все зависимости и расширение
make install-vscode
```

### Метод 2: Ручная установка
```bash
# Перейти в директорию VS Code расширения
cd dev_tools/vscode

# Установить зависимости
npm install

# Скомпилировать расширение
npm run compile

# Установить расширение в VS Code
code --install-extension tenge-language-support-1.0.0.vsix
```

## 🔧 Настройка VS Code

### 1. Открыть VS Code
```bash
# Открыть VS Code в директории проекта
code .
```

### 2. Установить расширение
- Открыть панель расширений (Ctrl+Shift+X)
- Найти "Tenge Language Support"
- Нажать "Install"

### 3. Настроить параметры
Добавить в `settings.json`:

```json
{
  "tenge.linter.enabled": true,
  "tenge.formatter.enabled": true,
  "tenge.compiler.path": "tenge",
  "tenge.formatter.indentationSize": 4,
  "tenge.formatter.useSpaces": true,
  "tenge.formatter.autoFormatOnSave": true,
  "tenge.formatter.autoFormatOnType": false,
  "files.associations": {
    "*.tng": "tenge"
  }
}
```

## 🎯 Использование

### 1. Создание файла Tenge
- Создать файл с расширением `.tng`
- VS Code автоматически распознает язык Tenge

### 2. Подсветка синтаксиса
- Автоматическая подсветка ключевых слов Tenge
- Подсветка типов, функций, операторов
- Подсветка строк и комментариев

### 3. Автодополнение
- Начать печатать функцию - появится автодополнение
- Использовать Ctrl+Space для принудительного автодополнения

### 4. Сниппеты
- `func` + Tab - создать функцию
- `var` + Tab - создать переменную
- `if` + Tab - создать условие
- `while` + Tab - создать цикл
- `class` + Tab - создать класс

### 5. Линтинг
- Автоматическая проверка кода в реальном времени
- Подсветка ошибок и предупреждений
- Быстрые исправления (Ctrl+.)

### 6. Форматирование
- Ctrl+Shift+F - форматировать документ
- Автоматическое форматирование при сохранении (если включено)

### 7. Компиляция
- Ctrl+Shift+P → "Tenge: Compile Document"
- Автоматическая компиляция в C код

## 🛠️ Разработка расширения

### 1. Клонирование и установка
```bash
# Клонировать репозиторий
git clone https://github.com/shanraq/tenge-web.git
cd tenge-web/dev_tools/vscode

# Установить зависимости
npm install
```

### 2. Разработка
```bash
# Запустить в режиме разработки
npm run watch

# В другом терминале запустить VS Code
code .
```

### 3. Тестирование
```bash
# Запустить тесты
npm test

# Запустить расширение в новом окне VS Code
npm run compile
```

### 4. Сборка
```bash
# Создать пакет расширения
npm run package
```

## 🔍 Отладка

### 1. Запуск в режиме отладки
- Открыть VS Code
- Нажать F5
- Выбрать "Run Extension"

### 2. Логи отладки
- Открыть Developer Tools (Help → Toggle Developer Tools)
- Проверить консоль на ошибки

### 3. Проверка расширения
```bash
# Проверить, что расширение установлено
code --list-extensions | grep tenge
```

## 📁 Структура расширения

```
dev_tools/vscode/
├── package.json              # Манифест расширения
├── src/
│   ├── extension.ts          # Основной файл расширения
│   ├── linter.ts             # Линтер
│   ├── formatter.ts          # Форматер
│   └── compiler.ts           # Компилятор
├── syntaxes/
│   └── tenge.tmLanguage.json # Грамматика синтаксиса
├── snippets/
│   └── tenge.json            # Сниппеты
├── language-configuration.json
└── tsconfig.json
```

## 🚨 Устранение проблем

### Проблема: Расширение не работает
**Решение:**
```bash
# Переустановить расширение
code --uninstall-extension tenge-language-support
cd dev_tools/vscode
npm run compile
code --install-extension tenge-language-support-1.0.0.vsix
```

### Проблема: Нет подсветки синтаксиса
**Решение:**
1. Проверить, что файл имеет расширение `.tng`
2. Проверить настройки VS Code:
```json
{
  "files.associations": {
    "*.tng": "tenge"
  }
}
```

### Проблема: Линтер не работает
**Решение:**
1. Проверить настройки:
```json
{
  "tenge.linter.enabled": true
}
```
2. Перезапустить VS Code

### Проблема: Форматер не работает
**Решение:**
1. Проверить настройки:
```json
{
  "tenge.formatter.enabled": true
}
```
2. Использовать Ctrl+Shift+F для принудительного форматирования

## 🎯 Горячие клавиши

| Действие | Горячая клавиша |
|----------|------------------|
| Форматирование документа | Ctrl+Shift+F |
| Линтинг документа | Ctrl+Shift+P → "Tenge: Lint Document" |
| Компиляция документа | Ctrl+Shift+P → "Tenge: Compile Document" |
| Автодополнение | Ctrl+Space |
| Быстрые исправления | Ctrl+. |

## 📚 Дополнительные ресурсы

- [VS Code Extension API](https://code.visualstudio.com/api)
- [Language Server Protocol](https://microsoft.github.io/language-server-protocol/)
- [TextMate Grammar](https://macromates.com/manual/en/language_grammars)
- [Tenge Language Documentation](https://shanraq.org/docs)

---

**Tenge VS Code Extension** - Полная поддержка языка Tenge в VS Code! 🚀




