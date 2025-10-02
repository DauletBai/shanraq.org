#!/usr/bin/env python3
"""
Tenge Language Formatter
Форматер для языка программирования Tenge
"""

import sys
import re
import argparse
from typing import List, Tuple
from pathlib import Path

class TengeFormatter:
    def __init__(self, indent_size: int = 4, use_spaces: bool = True):
        self.indent_size = indent_size
        self.use_spaces = use_spaces
        self.indent_char = ' ' if use_spaces else '\t'
        self.indent_multiplier = indent_size if use_spaces else 1

    def format_file(self, file_path: str) -> str:
        """Форматирование файла Tenge"""
        try:
            with open(file_path, 'r', encoding='utf-8') as f:
                content = f.read()
            
            return self.format_text(content)
            
        except Exception as e:
            print(f"Error reading file {file_path}: {e}", file=sys.stderr)
            return ""

    def format_text(self, text: str) -> str:
        """Форматирование текста"""
        lines = text.split('\n')
        formatted_lines = []
        indent_level = 0
        in_comment = False
        in_string = False
        string_char = ''

        for i, line in enumerate(lines):
            original_line = line
            
            # Пропустить пустые строки
            if line.strip() == '':
                formatted_lines.append('')
                continue

            # Обработка комментариев
            if line.strip().startswith('//'):
                formatted_lines.append(self._indent_line(line.strip(), indent_level))
                continue

            # Обработка блочных комментариев
            if '/*' in line:
                in_comment = True
            if in_comment:
                formatted_lines.append(self._indent_line(line.strip(), indent_level))
                if '*/' in line:
                    in_comment = False
                continue

            # Обработка строк
            if not in_string:
                if '"' in line or "'" in line or '`' in line:
                    in_string = True
                    string_char = '"' if '"' in line else ("'" if "'" in line else '`')
            else:
                if string_char in line:
                    in_string = False

            if in_string:
                formatted_lines.append(self._indent_line(line.strip(), indent_level))
                continue

            # Удалить существующие отступы
            line = line.strip()

            # Обработка закрывающих скобок
            if (line.startswith('}') or line.startswith('end') or 
                line.startswith('endif') or line.startswith('endwhile') or 
                line.startswith('endfor')):
                indent_level = max(0, indent_level - 1)

            # Форматирование строки
            formatted_line = self._format_line(line)
            formatted_lines.append(self._indent_line(formatted_line, indent_level))

            # Обработка открывающих скобок и управляющих структур
            if ('{' in line or 
                re.match(r'^(atqar|jasau|eger|azirshe|while|for|if|else|function|class|struct)\b', line)):
                indent_level += 1

        return '\n'.join(formatted_lines)

    def _format_line(self, line: str) -> str:
        """Форматирование одной строки"""
        # Добавить пробелы вокруг операторов
        line = re.sub(r'([=+\-*/%<>!&|^])(?!=)', r' \1 ', line)
        line = re.sub(r'([=+\-*/%<>!&|^])(?!=)', r' \1 ', line)
        
        # Удалить множественные пробелы
        line = re.sub(r'\s+', ' ', line)
        
        # Добавить пробелы после запятых
        line = line.replace(',', ', ')
        
        # Добавить пробелы вокруг скобок
        line = line.replace('(', ' (')
        line = line.replace(')', ') ')
        
        # Удалить множественные пробелы снова
        line = re.sub(r'\s+', ' ', line)
        
        # Обрезать строку
        line = line.strip()
        
        # Убедиться в наличии точки с запятой в конце операторов
        if (line and not line.endswith(';') and not line.endswith('{') and 
            not line.endswith('}') and not line.startswith('//') and 
            not line.startswith('/*') and 
            not re.match(r'^(atqar|jasau|eger|azirshe|qaytar|end|import|export)', line)):
            line += ';'

        return line

    def _indent_line(self, line: str, indent_level: int) -> str:
        """Добавление отступов к строке"""
        if line.strip() == '':
            return ''

        indent = self.indent_char * (indent_level * self.indent_multiplier)
        return indent + line

    def format_agglutinative_names(self, text: str) -> str:
        """Форматирование агглютинативных имен"""
        # Автоматическое исправление английских слов на казахские эквиваленты
        replacements = {
            'create': 'jasau',
            'get': 'alu',
            'add': 'qosu',
            'update': 'zhangartu',
            'delete': 'zhoyu',
            'check': 'tekseru',
            'optimize': 'opt',
            'engine': 'eng',
            'manager': 'man',
            'function': 'atqar',
            'variable': 'jasau',
            'if': 'eger',
            'while': 'azirshe',
            'return': 'qaytar'
        }
        
        for english, kazakh in replacements.items():
            # Заменить только в контексте функций и переменных
            text = re.sub(rf'\b{english}\b', kazakh, text)
        
        return text

    def format_imports(self, text: str) -> str:
        """Форматирование импортов"""
        lines = text.split('\n')
        import_lines = []
        other_lines = []
        
        for line in lines:
            if line.strip().startswith('import '):
                import_lines.append(line.strip())
            else:
                other_lines.append(line)
        
        # Сортировка импортов
        import_lines.sort()
        
        # Объединение
        if import_lines:
            formatted_lines = import_lines + [''] + other_lines
        else:
            formatted_lines = other_lines
        
        return '\n'.join(formatted_lines)

    def format_functions(self, text: str) -> str:
        """Форматирование функций"""
        # Добавить пустые строки между функциями
        text = re.sub(r'}\s*atqar', '}\n\natqar', text)
        text = re.sub(r'}\s*jasau', '}\n\njasau', text)
        
        return text

def main():
    parser = argparse.ArgumentParser(description="Tenge Language Formatter")
    parser.add_argument("files", nargs="+", help="Tenge files to format")
    parser.add_argument("-i", "--in-place", action="store_true", help="Format files in place")
    parser.add_argument("-s", "--indent-size", type=int, default=4, help="Indentation size")
    parser.add_argument("-t", "--use-tabs", action="store_true", help="Use tabs instead of spaces")
    parser.add_argument("--agglutinative", action="store_true", help="Format agglutinative names")
    parser.add_argument("--imports", action="store_true", help="Format imports")
    parser.add_argument("--functions", action="store_true", help="Format functions")
    
    args = parser.parse_args()
    
    formatter = TengeFormatter(
        indent_size=args.indent_size,
        use_spaces=not args.use_tabs
    )
    
    for file_path in args.files:
        if not Path(file_path).exists():
            print(f"Error: File '{file_path}' not found", file=sys.stderr)
            continue
        
        try:
            formatted_content = formatter.format_file(file_path)
            
            if args.agglutinative:
                formatted_content = formatter.format_agglutinative_names(formatted_content)
            
            if args.imports:
                formatted_content = formatter.format_imports(formatted_content)
            
            if args.functions:
                formatted_content = formatter.format_functions(formatted_content)
            
            if args.in_place:
                with open(file_path, 'w', encoding='utf-8') as f:
                    f.write(formatted_content)
                print(f"Formatted: {file_path}")
            else:
                print(formatted_content)
                
        except Exception as e:
            print(f"Error formatting {file_path}: {e}", file=sys.stderr)

if __name__ == "__main__":
    main()
