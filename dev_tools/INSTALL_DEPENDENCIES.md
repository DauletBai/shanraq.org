# Установка зависимостей для Tenge Dev Tools
# Installing Dependencies for Tenge Dev Tools

## 🚨 Проблема
У вас не установлены необходимые зависимости:
- Node.js (для npm)
- VS Code (для code команды)

## 🔧 Решение

### 1. Установка Node.js

#### Способ 1: Через официальный сайт (рекомендуется)
1. Перейти на https://nodejs.org/
2. Скачать LTS версию для macOS
3. Установить скачанный файл

#### Способ 2: Через Homebrew
```bash
# Установить Homebrew (если не установлен)
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

# Установить Node.js
brew install node
```

#### Способ 3: Через nvm (Node Version Manager)
```bash
# Установить nvm
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.0/install.sh | bash

# Перезапустить терминал или выполнить
source ~/.bashrc

# Установить последнюю LTS версию Node.js
nvm install --lts
nvm use --lts
```

### 2. Установка VS Code

#### Способ 1: Через официальный сайт (рекомендуется)
1. Перейти на https://code.visualstudio.com/
2. Скачать VS Code для macOS
3. Установить скачанный файл

#### Способ 2: Через Homebrew
```bash
# Установить VS Code
brew install --cask visual-studio-code
```

#### Способ 3: Через Mac App Store
1. Открыть Mac App Store
2. Найти "Visual Studio Code"
3. Установить

### 3. Проверка установки

```bash
# Проверить Node.js
node --version
npm --version

# Проверить VS Code
code --version
```

## 🚀 После установки зависимостей

### 1. Установка Tenge VS Code расширения
```bash
# Перейти в директорию dev_tools
cd dev_tools

# Установить расширение
make install-vscode-quick
```

### 2. Альтернативная установка
```bash
# Если make не работает, использовать прямой скрипт
./install_vscode.sh
```

### 3. Ручная установка
```bash
# Перейти в директорию vscode
cd dev_tools/vscode

# Установить зависимости
npm install

# Скомпилировать расширение
npm run compile

# Создать пакет расширения
npm run package

# Установить расширение
code --install-extension tenge-language-support-1.0.0.vsix
```

## 🧪 Тестирование

### 1. Проверить установку
```bash
# Проверить, что расширение установлено
code --list-extensions | grep tenge
```

### 2. Открыть тестовый файл
```bash
# Открыть VS Code
code .

# Открыть тестовый файл
code test_example.tng
```

### 3. Проверить функциональность
- Подсветка синтаксиса
- Автодополнение
- Сниппеты
- Линтинг
- Форматирование

## 🚨 Устранение проблем

### Проблема: npm не найден
**Решение:**
```bash
# Переустановить Node.js
brew uninstall node
brew install node

# Или использовать nvm
nvm install --lts
nvm use --lts
```

### Проблема: code не найден
**Решение:**
```bash
# Добавить VS Code в PATH
echo 'export PATH="/Applications/Visual Studio Code.app/Contents/Resources/app/bin:$PATH"' >> ~/.zshrc
source ~/.zshrc

# Или создать симлинк
sudo ln -s "/Applications/Visual Studio Code.app/Contents/Resources/app/bin/code" /usr/local/bin/code
```

### Проблема: Права доступа
**Решение:**
```bash
# Дать права на выполнение скрипта
chmod +x install_vscode.sh

# Запустить скрипт
./install_vscode.sh
```

## 📚 Дополнительные ресурсы

- [Node.js Official Website](https://nodejs.org/)
- [VS Code Official Website](https://code.visualstudio.com/)
- [Homebrew](https://brew.sh/)
- [nvm](https://github.com/nvm-sh/nvm)

## 🎯 После успешной установки

1. **Перезапустить терминал**
2. **Проверить установку:**
   ```bash
   node --version
   npm --version
   code --version
   ```
3. **Установить Tenge расширение:**
   ```bash
   make install-vscode-quick
   ```
4. **Открыть VS Code и наслаждаться поддержкой Tenge!**

---

**Удачной установки! 🚀**




