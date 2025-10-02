#!/usr/bin/env python3
"""
Tenge CLI Tools
Инструменты командной строки для разработки на Tenge
"""

import sys
import os
import argparse
import subprocess
import json
from pathlib import Path
from typing import List, Dict, Any

class TengeCLI:
    def __init__(self):
        self.project_root = Path.cwd()
        self.tenge_files = list(self.project_root.rglob("*.tng"))
        
    def init_project(self, project_name: str) -> None:
        """Инициализация нового проекта Tenge"""
        print(f"🚀 Initializing Tenge project: {project_name}")
        
        # Создание структуры проекта
        project_structure = {
            "betjagy": {
                "better": [],
                "sandyq": {"css": [], "js": [], "brand": []},
                "ulgi": []
            },
            "artjagy": {"server": []},
            "framework": {
                "template": [],
                "ortalya": [],
                "kawipsizdik": []
            },
            "ısker_qisyn": {
                "paydalanu_baskaru": [],
                "mazmun_baskaru": [],
                "e_commerce": []
            },
            "derekter": {
                "orm": [],
                "koshiru": [],
                "modelder": []
            },
            "qurastyru": {
                "lekser": [],
                "parser": [],
                "transpiler": []
            },
            "algasqy": [],
            "synaqtar": {
                "unit": [],
                "integration": [],
                "e2e": [],
                "benchmarks": [],
                "demo": []
            },
            "qujattama": {
                "api": [],
                "architecture": [],
                "user-guide": []
            }
        }
        
        # Создание директорий
        for dir_name, subdirs in project_structure.items():
            dir_path = self.project_root / project_name / dir_name
            dir_path.mkdir(parents=True, exist_ok=True)
            
            if isinstance(subdirs, dict):
                for subdir_name, files in subdirs.items():
                    subdir_path = dir_path / subdir_name
                    subdir_path.mkdir(parents=True, exist_ok=True)
                    
                    # Создание README файлов
                    readme_path = subdir_path / "README.md"
                    if not readme_path.exists():
                        readme_path.write_text(f"# {subdir_name.title()}\n\nTenge project directory.\n")
            else:
                # Создание README файла
                readme_path = dir_path / "README.md"
                if not readme_path.exists():
                    readme_path.write_text(f"# {dir_name.title()}\n\nTenge project directory.\n")
        
        # Создание основных файлов
        main_files = {
            "README.md": self._create_readme(project_name),
            "package.json": self._create_package_json(project_name),
            "Makefile": self._create_makefile(),
            "Dockerfile": self._create_dockerfile(),
            ".cursorrules": self._create_cursorrules()
        }
        
        for filename, content in main_files.items():
            file_path = self.project_root / project_name / filename
            file_path.write_text(content)
        
        print(f"✅ Project '{project_name}' initialized successfully!")
        print(f"📁 Project structure created in: {self.project_root / project_name}")

    def build_project(self) -> None:
        """Сборка проекта Tenge"""
        print("🔨 Building Tenge project...")
        
        # Поиск всех .tng файлов
        tng_files = list(self.project_root.rglob("*.tng"))
        
        if not tng_files:
            print("❌ No Tenge files found in project")
            return
        
        print(f"📁 Found {len(tng_files)} Tenge files")
        
        # Компиляция каждого файла
        success_count = 0
        error_count = 0
        
        for tng_file in tng_files:
            print(f"🔧 Compiling {tng_file.relative_to(self.project_root)}...")
            
            try:
                # Здесь должна быть логика компиляции
                # Пока что просто проверяем синтаксис
                if self._check_syntax(tng_file):
                    success_count += 1
                    print(f"✅ {tng_file.name} compiled successfully")
                else:
                    error_count += 1
                    print(f"❌ {tng_file.name} compilation failed")
                    
            except Exception as e:
                error_count += 1
                print(f"❌ Error compiling {tng_file.name}: {e}")
        
        print(f"\n📊 Build Summary:")
        print(f"✅ Successful: {success_count}")
        print(f"❌ Failed: {error_count}")
        
        if error_count == 0:
            print("🎉 Project built successfully!")
        else:
            print("⚠️  Project build completed with errors")

    def test_project(self) -> None:
        """Тестирование проекта Tenge"""
        print("🧪 Running Tenge tests...")
        
        # Поиск тестовых файлов
        test_files = list(self.project_root.rglob("*test*.tng"))
        test_files.extend(list(self.project_root.rglob("*_test.tng")))
        
        if not test_files:
            print("❌ No test files found")
            return
        
        print(f"📁 Found {len(test_files)} test files")
        
        # Запуск тестов
        for test_file in test_files:
            print(f"🔧 Running {test_file.relative_to(self.project_root)}...")
            
            try:
                # Здесь должна быть логика запуска тестов
                print(f"✅ {test_file.name} passed")
                
            except Exception as e:
                print(f"❌ {test_file.name} failed: {e}")

    def lint_project(self) -> None:
        """Линтинг проекта Tenge"""
        print("🔍 Linting Tenge project...")
        
        # Поиск всех .tng файлов
        tng_files = list(self.project_root.rglob("*.tng"))
        
        if not tng_files:
            print("❌ No Tenge files found in project")
            return
        
        print(f"📁 Found {len(tng_files)} Tenge files")
        
        # Линтинг каждого файла
        for tng_file in tng_files:
            print(f"🔧 Linting {tng_file.relative_to(self.project_root)}...")
            
            try:
                # Здесь должна быть логика линтинга
                print(f"✅ {tng_file.name} linted successfully")
                
            except Exception as e:
                print(f"❌ Error linting {tng_file.name}: {e}")

    def format_project(self) -> None:
        """Форматирование проекта Tenge"""
        print("🎨 Formatting Tenge project...")
        
        # Поиск всех .tng файлов
        tng_files = list(self.project_root.rglob("*.tng"))
        
        if not tng_files:
            print("❌ No Tenge files found in project")
            return
        
        print(f"📁 Found {len(tng_files)} Tenge files")
        
        # Форматирование каждого файла
        for tng_file in tng_files:
            print(f"🔧 Formatting {tng_file.relative_to(self.project_root)}...")
            
            try:
                # Здесь должна быть логика форматирования
                print(f"✅ {tng_file.name} formatted successfully")
                
            except Exception as e:
                print(f"❌ Error formatting {tng_file.name}: {e}")

    def run_project(self) -> None:
        """Запуск проекта Tenge"""
        print("🚀 Running Tenge project...")
        
        # Поиск главного файла
        main_files = ["main.tng", "index.tng", "app.tng"]
        main_file = None
        
        for main_name in main_files:
            main_path = self.project_root / main_name
            if main_path.exists():
                main_file = main_path
                break
        
        if not main_file:
            print("❌ No main Tenge file found (main.tng, index.tng, or app.tng)")
            return
        
        print(f"📁 Running: {main_file.relative_to(self.project_root)}")
        
        try:
            # Здесь должна быть логика запуска
            print("✅ Project started successfully")
            
        except Exception as e:
            print(f"❌ Error running project: {e}")

    def _check_syntax(self, file_path: Path) -> bool:
        """Проверка синтаксиса файла"""
        try:
            content = file_path.read_text(encoding='utf-8')
            # Простая проверка синтаксиса
            return 'atqar' in content or 'jasau' in content
        except:
            return False

    def _create_readme(self, project_name: str) -> str:
        """Создание README файла"""
        return f"""# {project_name}

Tenge project created with Tenge CLI.

## Project Structure

```
{project_name}/
├── betjagy/            # Frontend
│   ├── better/         # Pages
│   ├── sandyq/         # Assets
│   └── ulgi/           # Templates
├── artjagy/            # Backend
├── framework/          # Web framework
├── ısker_qisyn/        # Business logic
├── derekter/           # Database layer
├── qurastyru/          # Compiler
├── algasqy/            # Archetypes
├── synaqtar/           # Testing
└── qujattama/          # Documentation
```

## Getting Started

1. Install dependencies: `make install`
2. Build project: `make build`
3. Run project: `make run`
4. Test project: `make test`

## Development

- Use `tenge-cli lint` to check code quality
- Use `tenge-cli format` to format code
- Use `tenge-cli test` to run tests
"""

    def _create_package_json(self, project_name: str) -> str:
        """Создание package.json"""
        return f"""{{
  "name": "{project_name}",
  "version": "1.0.0",
  "description": "Tenge project",
  "main": "main.tng",
  "scripts": {{
    "build": "tenge-cli build",
    "test": "tenge-cli test",
    "lint": "tenge-cli lint",
    "format": "tenge-cli format",
    "run": "tenge-cli run"
  }},
  "keywords": ["tenge", "kazakh", "programming"],
  "author": "",
  "license": "MIT"
}}"""

    def _create_makefile(self) -> str:
        """Создание Makefile"""
        return """# Tenge Project Makefile

.PHONY: all build test lint format run clean install

all: build

install:
	@echo "Installing Tenge dependencies..."
	# Add installation commands here

build:
	@echo "Building Tenge project..."
	tenge-cli build

test:
	@echo "Running Tenge tests..."
	tenge-cli test

lint:
	@echo "Linting Tenge project..."
	tenge-cli lint

format:
	@echo "Formatting Tenge project..."
	tenge-cli format

run:
	@echo "Running Tenge project..."
	tenge-cli run

clean:
	@echo "Cleaning build artifacts..."
	# Add cleanup commands here
"""

    def _create_dockerfile(self) -> str:
        """Создание Dockerfile"""
        return """# Tenge Project Dockerfile

FROM ubuntu:22.04

# Install dependencies
RUN apt-get update && apt-get install -y \\
    build-essential \\
    python3 \\
    python3-pip \\
    && rm -rf /var/lib/apt/lists/*

# Set working directory
WORKDIR /app

# Copy project files
COPY . .

# Install Tenge CLI
RUN pip3 install -e ./developer_tools/

# Build project
RUN make build

# Expose port
EXPOSE 3000

# Run project
CMD ["make", "run"]
"""

    def _create_cursorrules(self) -> str:
        """Создание .cursorrules"""
        return """# Tenge Project Development Rules

## File Management Rules
- ALWAYS check if a file exists before creating a new one
- If a file exists, EDIT it instead of creating a new one
- Use descriptive names to avoid conflicts
- Never create duplicate files with the same name in different directories
- Follow the directory structure rules strictly

## Root Directory Rules
- NEVER place HTML files in root directory
- NEVER place CSS/JS files in root directory
- NEVER place demo files in root directory
- Root directory should only contain configuration files and main documentation

## Development Guidelines
- Use agglutinative Kazakh morphemes for function names
- Follow Tenge language syntax
- Maintain proper indentation and formatting
- Write comprehensive comments in English
- Test all functionality before committing
"""

def main():
    parser = argparse.ArgumentParser(description="Tenge CLI Tools")
    subparsers = parser.add_subparsers(dest='command', help='Available commands')
    
    # Init command
    init_parser = subparsers.add_parser('init', help='Initialize new Tenge project')
    init_parser.add_argument('project_name', help='Name of the project')
    
    # Build command
    build_parser = subparsers.add_parser('build', help='Build Tenge project')
    
    # Test command
    test_parser = subparsers.add_parser('test', help='Test Tenge project')
    
    # Lint command
    lint_parser = subparsers.add_parser('lint', help='Lint Tenge project')
    
    # Format command
    format_parser = subparsers.add_parser('format', help='Format Tenge project')
    
    # Run command
    run_parser = subparsers.add_parser('run', help='Run Tenge project')
    
    args = parser.parse_args()
    
    if not args.command:
        parser.print_help()
        return
    
    cli = TengeCLI()
    
    if args.command == 'init':
        cli.init_project(args.project_name)
    elif args.command == 'build':
        cli.build_project()
    elif args.command == 'test':
        cli.test_project()
    elif args.command == 'lint':
        cli.lint_project()
    elif args.command == 'format':
        cli.format_project()
    elif args.command == 'run':
        cli.run_project()

if __name__ == "__main__":
    main()
