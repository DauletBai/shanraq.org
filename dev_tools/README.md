# Tenge Dev Tools
# Инструменты разработчика для языка Tenge

Полный набор инструментов для разработки на языке программирования Tenge, включая подсветку синтаксиса, линтер, форматер и CLI инструменты.

## 🚀 Быстрый старт / Quick Start

### Установка / Installation

```bash
# Клонирование репозитория
git clone https://github.com/shanraq/tenge-web.git
cd tenge-web/developer_tools

# Установка всех инструментов
make install

# Или установка по отдельности
make install-vscode    # VS Code extension
make install-sublime   # Sublime Text package
```

### Использование / Usage

```bash
# Создание нового проекта
make new-project

# Сборка проекта
make build-project

# Тестирование проекта
make test-project

# Запуск проекта
make run-project

# Линтинг проекта
make lint-project

# Форматирование проекта
make format-project
```

## 🛠️ Инструменты / Tools

### 1. VS Code Extension
**Расширение для Visual Studio Code**

- ✅ Подсветка синтаксиса Tenge
- ✅ Автодополнение кода
- ✅ Линтинг в реальном времени
- ✅ Форматирование кода
- ✅ Сниппеты для быстрой разработки
- ✅ Интеграция с компилятором

**Установка:**
```bash
cd vscode
npm install
npm run compile
code --install-extension tenge-language-support-1.0.0.vsix
```

### 2. Sublime Text Package
**Пакет для Sublime Text**

- ✅ Подсветка синтаксиса Tenge
- ✅ Сниппеты для быстрой разработки
- ✅ Автодополнение кода

**Установка:**
```bash
# Автоматическая установка
make install-sublime

# Ручная установка
cp sublime/tenge.sublime-syntax ~/.config/sublime-text/Packages/User/
cp sublime/tenge.sublime-snippets ~/.config/sublime-text/Packages/User/
```

### 3. Linter
**Линтер для проверки качества кода**

- ✅ Проверка агглютинативных паттернов
- ✅ Проверка использования ключевых слов Tenge
- ✅ Проверка типов
- ✅ Проверка точек с запятой
- ✅ Проверка отступов
- ✅ Проверка соответствия казахскому языку

**Использование:**
```bash
# Линтинг одного файла
python3 linter/tenge-linter.py file.tng

# Линтинг всех файлов в проекте
python3 linter/tenge-linter.py *.tng

# Линтинг с выводом в JSON
python3 linter/tenge-linter.py -f json file.tng

# Линтинг с выходным кодом ошибки
python3 linter/tenge-linter.py --exit-code file.tng
```

### 4. Formatter
**Форматер для автоматического форматирования кода**

- ✅ Автоматическое форматирование отступов
- ✅ Форматирование операторов
- ✅ Форматирование строк
- ✅ Форматирование агглютинативных имен
- ✅ Форматирование импортов
- ✅ Форматирование функций

**Использование:**
```bash
# Форматирование файла
python3 formatter/tenge-formatter.py file.tng

# Форматирование с заменой файла
python3 formatter/tenge-formatter.py -i file.tng

# Форматирование с настройками
python3 formatter/tenge-formatter.py -s 2 -t file.tng

# Форматирование агглютинативных имен
python3 formatter/tenge-formatter.py --agglutinative file.tng
```

### 5. CLI Tools
**Инструменты командной строки**

- ✅ Инициализация проектов
- ✅ Сборка проектов
- ✅ Тестирование проектов
- ✅ Линтинг проектов
- ✅ Форматирование проектов
- ✅ Запуск проектов

**Использование:**
```bash
# Инициализация нового проекта
python3 cli/tenge-cli.py init my-project

# Сборка проекта
python3 cli/tenge-cli.py build

# Тестирование проекта
python3 cli/tenge-cli.py test

# Линтинг проекта
python3 cli/tenge-cli.py lint

# Форматирование проекта
python3 cli/tenge-cli.py format

# Запуск проекта
python3 cli/tenge-cli.py run
```

## 📁 Структура проекта / Project Structure

```
developer_tools/
├── vscode/                 # VS Code extension
│   ├── package.json       # Extension manifest
│   ├── src/               # Source code
│   ├── syntaxes/          # Syntax highlighting
│   └── snippets/          # Code snippets
├── sublime/               # Sublime Text package
│   ├── tenge.sublime-syntax
│   └── tenge.sublime-snippets
├── linter/                # Linter
│   └── tenge-linter.py
├── formatter/             # Formatter
│   └── tenge-formatter.py
├── cli/                   # CLI tools
│   └── tenge-cli.py
├── Makefile               # Build automation
├── requirements.txt       # Python dependencies
└── README.md              # This file
```

## 🔧 Настройка / Configuration

### VS Code Settings
```json
{
  "tenge.linter.enabled": true,
  "tenge.formatter.enabled": true,
  "tenge.compiler.path": "tenge",
  "tenge.formatter.indentationSize": 4,
  "tenge.formatter.useSpaces": true,
  "tenge.formatter.autoFormatOnSave": true,
  "tenge.formatter.autoFormatOnType": false
}
```

### Sublime Text Settings
```json
{
  "tenge_linter_enabled": true,
  "tenge_formatter_enabled": true,
  "tenge_indent_size": 4,
  "tenge_use_spaces": true
}
```

## 🧪 Тестирование / Testing

```bash
# Тестирование всех инструментов
make test

# Тестирование линтера
python3 linter/tenge-linter.py --help

# Тестирование форматера
python3 formatter/tenge-formatter.py --help

# Тестирование CLI
python3 cli/tenge-cli.py --help
```

## 🚀 Разработка / Development

### Установка среды разработки
```bash
# Установка зависимостей
make dev-env

# Установка Python зависимостей
pip3 install -r requirements.txt

# Установка Node.js зависимостей
cd vscode && npm install
```

### Сборка
```bash
# Сборка всех инструментов
make build

# Сборка VS Code extension
make build-vscode
```

### Очистка
```bash
# Очистка артефактов сборки
make clean
```

## 📚 Документация / Documentation

### Линтер
- Проверяет агглютинативные паттерны именования
- Проверяет использование ключевых слов Tenge
- Проверяет типы переменных
- Проверяет точки с запятой
- Проверяет отступы
- Проверяет соответствие казахскому языку

### Форматер
- Автоматическое форматирование отступов
- Форматирование операторов и скобок
- Форматирование строк и комментариев
- Форматирование агглютинативных имен
- Форматирование импортов и функций

### CLI Tools
- Инициализация проектов с правильной структурой
- Сборка и компиляция проектов
- Тестирование проектов
- Линтинг и форматирование проектов
- Запуск проектов

## 🤝 Вклад в проект / Contributing

1. Форкните репозиторий
2. Создайте ветку для новой функции
3. Внесите изменения
4. Добавьте тесты
5. Создайте Pull Request

## 📄 Лицензия / License

MIT License - см. файл LICENSE для деталей.

## 🆘 Поддержка / Support

- 📧 Email: support@shanraq.org
- 🐛 Issues: https://github.com/shanraq/tenge-web/issues
- 📖 Documentation: https://shanraq.org/docs
- 💬 Community: https://shanraq.org/community

## 🎯 Roadmap

- [ ] Поддержка других редакторов (Vim, Emacs, IntelliJ)
- [ ] Интеграция с CI/CD
- [ ] Автоматическое тестирование
- [ ] Профилирование производительности
- [ ] Отладчик для Tenge
- [ ] Интеграция с Git hooks
- [ ] Поддержка плагинов
- [ ] Веб-интерфейс для инструментов

---

**Tenge Developer Tools** - Полный набор инструментов для разработки на языке Tenge! 🚀
