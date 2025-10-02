# 🚀 Быстрый старт с Tenge Dev Tools
# Quick Start with Tenge Dev Tools

## ⚡ Установка за 3 шага / Installation in 3 steps

### 1. Установка VS Code расширения
```bash
# Перейти в директорию dev_tools
cd dev_tools

# Быстрая установка VS Code расширения
make install-vscode-quick
```

### 2. Проверка установки
```bash
# Открыть VS Code
code .

# Открыть тестовый файл
code test_example.tng
```

### 3. Настройка VS Code
Добавить в настройки VS Code (`Ctrl+,` → Settings → Open Settings JSON):

```json
{
  "tenge.linter.enabled": true,
  "tenge.formatter.enabled": true,
  "tenge.formatter.autoFormatOnSave": true,
  "files.associations": {
    "*.tng": "tenge"
  }
}
```

## 🎯 Что вы получите

### ✅ Подсветка синтаксиса
- Ключевые слова Tenge (atqar, jasau, eger, etc.)
- Типы данных (san, jol, aqıqat, JsonObject, etc.)
- Функции и переменные
- Строки и комментарии

### ✅ Автодополнение
- Начать печатать функцию → автодополнение
- Ctrl+Space для принудительного автодополнения
- Подсказки для параметров функций

### ✅ Сниппеты
- `func` + Tab → создать функцию
- `var` + Tab → создать переменную
- `if` + Tab → создать условие
- `while` + Tab → создать цикл
- `class` + Tab → создать класс

### ✅ Линтинг
- Проверка агглютинативных паттернов
- Проверка ключевых слов Tenge
- Проверка типов переменных
- Проверка точек с запятой
- Проверка соответствия казахскому языку

### ✅ Форматирование
- Ctrl+Shift+F → форматировать документ
- Автоматическое форматирование при сохранении
- Правильные отступы и пробелы

### ✅ Компиляция
- Ctrl+Shift+P → "Tenge: Compile Document"
- Автоматическая компиляция в C код

## 🧪 Тестирование

### 1. Открыть тестовый файл
```bash
# В VS Code открыть
code test_example.tng
```

### 2. Проверить подсветку синтаксиса
- Ключевые слова должны быть подсвечены
- Типы должны быть подсвечены
- Строки должны быть подсвечены

### 3. Проверить автодополнение
- Начать печатать `atqar` → должно появиться автодополнение
- Начать печатать `jasau` → должно появиться автодополнение

### 4. Проверить сниппеты
- Напечатать `func` и нажать Tab
- Должен появиться шаблон функции

### 5. Проверить линтинг
- Должны появиться подсказки об ошибках
- Подсветка проблемных мест

### 6. Проверить форматирование
- Ctrl+Shift+F
- Код должен отформатироваться

## 🔧 Дополнительные настройки

### Настройки VS Code
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

### Горячие клавиши
| Действие | Горячая клавиша |
|----------|------------------|
| Форматирование | Ctrl+Shift+F |
| Линтинг | Ctrl+Shift+P → "Tenge: Lint Document" |
| Компиляция | Ctrl+Shift+P → "Tenge: Compile Document" |
| Автодополнение | Ctrl+Space |
| Быстрые исправления | Ctrl+. |

## 🚨 Устранение проблем

### Проблема: Нет подсветки синтаксиса
**Решение:**
1. Проверить, что файл имеет расширение `.tng`
2. Проверить настройки VS Code
3. Перезапустить VS Code

### Проблема: Линтер не работает
**Решение:**
1. Проверить настройки: `"tenge.linter.enabled": true`
2. Перезапустить VS Code

### Проблема: Форматер не работает
**Решение:**
1. Проверить настройки: `"tenge.formatter.enabled": true`
2. Использовать Ctrl+Shift+F

### Проблема: Расширение не установлено
**Решение:**
```bash
# Переустановить расширение
make install-vscode-quick
```

## 📚 Дополнительные ресурсы

- [VSCODE_INSTALLATION.md](VSCODE_INSTALLATION.md) - Подробная инструкция по установке
- [README.md](README.md) - Полная документация
- [test_example.tng](test_example.tng) - Пример файла Tenge

## 🎉 Готово!

Теперь у вас есть полная поддержка языка Tenge в VS Code! 

**Следующие шаги:**
1. Создайте свой первый файл `.tng`
2. Начните писать код на Tenge
3. Используйте все возможности расширения

**Удачной разработки! 🚀**




