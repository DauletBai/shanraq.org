# 🚀 Настройка Tenge Dev Tools
# Setup Guide for Tenge Dev Tools

## ❌ Проблема
У вас возникла ошибка:
```
/bin/sh: npm: command not found
/bin/sh: code: command not found
```

Это означает, что не установлены необходимые зависимости.

## ✅ Решение

### 1. Проверить зависимости
```bash
make check-deps
```

### 2. Установить недостающие зависимости

#### Если не установлен Node.js:
```bash
# Способ 1: Через Homebrew (рекомендуется)
brew install node

# Способ 2: Скачать с официального сайта
# Перейти на https://nodejs.org/ и скачать LTS версию
```

#### Если не установлен VS Code:
```bash
# Способ 1: Через Homebrew
brew install --cask visual-studio-code

# Способ 2: Скачать с официального сайта
# Перейти на https://code.visualstudio.com/ и скачать
```

### 3. Добавить VS Code в PATH (если нужно)
```bash
# Добавить в ~/.zshrc
echo 'export PATH="/Applications/Visual Studio Code.app/Contents/Resources/app/bin:$PATH"' >> ~/.zshrc
source ~/.zshrc
```

### 4. Проверить установку
```bash
# Проверить Node.js
node --version
npm --version

# Проверить VS Code
code --version
```

### 5. Установить Tenge расширение
```bash
# После установки всех зависимостей
make install-vscode-quick
```

## 🎯 Быстрый старт

### Если у вас есть Homebrew:
```bash
# Установить зависимости
brew install node
brew install --cask visual-studio-code

# Добавить VS Code в PATH
echo 'export PATH="/Applications/Visual Studio Code.app/Contents/Resources/app/bin:$PATH"' >> ~/.zshrc
source ~/.zshrc

# Установить Tenge расширение
make install-vscode-quick
```

### Если у вас нет Homebrew:
1. Установить Homebrew: https://brew.sh/
2. Выполнить команды выше

## 🧪 Тестирование

После установки:
```bash
# Проверить зависимости
make check-deps

# Установить расширение
make install-vscode-quick

# Открыть VS Code
code .

# Открыть тестовый файл
code test_example.tng
```

## 🚨 Если что-то не работает

### Проблема: npm не найден
```bash
# Переустановить Node.js
brew uninstall node
brew install node
```

### Проблема: code не найден
```bash
# Создать симлинк
sudo ln -s "/Applications/Visual Studio Code.app/Contents/Resources/app/bin/code" /usr/local/bin/code
```

### Проблема: Права доступа
```bash
# Дать права на выполнение
chmod +x check_dependencies.sh
chmod +x install_vscode.sh
```

## 📚 Дополнительные ресурсы

- [Node.js](https://nodejs.org/)
- [VS Code](https://code.visualstudio.com/)
- [Homebrew](https://brew.sh/)
- [INSTALL_DEPENDENCIES.md](INSTALL_DEPENDENCIES.md) - Подробная инструкция

## 🎉 После успешной установки

1. **Перезапустить терминал**
2. **Проверить:** `make check-deps`
3. **Установить:** `make install-vscode-quick`
4. **Открыть VS Code и наслаждаться поддержкой Tenge!**

---

**Удачной установки! 🚀**








