#!/usr/bin/env python3
"""
Tenge Language Linter
Линтер для языка программирования Tenge
"""

import sys
import re
import argparse
import json
from typing import List, Dict, Any, Tuple
from dataclasses import dataclass
from pathlib import Path

@dataclass
class LintError:
    line: int
    column: int
    message: str
    severity: str
    code: str
    suggestion: str = ""

class TengeLinter:
    def __init__(self):
        self.errors: List[LintError] = []
        self.agglutinative_suffixes = [
            '_jasau', '_alu', '_qosu', '_zhangartu', '_zhoyu', '_tekseru', '_opt',
            '_eng', '_man', '_negizgi', '_qoldanu', '_marshrut', '_baskaru'
        ]
        self.tenge_keywords = [
            'atqar', 'jasau', 'eger', 'azirshe', 'qaytar', 'end', 'endif', 'endwhile', 'endfor',
            'import', 'export', 'public', 'private', 'static', 'const', 'var', 'let',
            'true', 'false', 'null', 'NULL', 'aqıqat_ras', 'aqıqat_jin'
        ]
        self.tenge_types = [
            'san', 'jol', 'aqıqat', 'JsonObject', 'WebServer', 'Component', 'TemplateEngine',
            'ArchetypeEngine', 'MorphemeEngine', 'PhonemeEngine', 'Vector', 'Map', 'Array',
            'Function', 'Class', 'Struct', 'Interface', 'Enum'
        ]
        self.english_words = [
            'create', 'get', 'add', 'update', 'delete', 'check', 'optimize', 'engine', 'manager',
            'function', 'variable', 'class', 'method', 'property', 'return', 'if', 'else', 'while',
            'for', 'import', 'export', 'public', 'private', 'static', 'const', 'var', 'let'
        ]

    def lint_file(self, file_path: str) -> List[LintError]:
        """Линтинг файла Tenge"""
        self.errors = []
        
        try:
            with open(file_path, 'r', encoding='utf-8') as f:
                content = f.read()
            
            lines = content.split('\n')
            
            for line_num, line in enumerate(lines, 1):
                self._lint_line(line, line_num)
            
            return self.errors
            
        except Exception as e:
            return [LintError(0, 0, f"Error reading file: {e}", "error", "file.read")]

    def _lint_line(self, line: str, line_num: int) -> None:
        """Линтинг одной строки"""
        original_line = line
        line = line.strip()
        
        # Пропустить пустые строки и комментарии
        if not line or line.startswith('//') or line.startswith('/*'):
            return
        
        # Проверка агглютинативных паттернов
        self._check_agglutinative_patterns(line, line_num)
        
        # Проверка ключевых слов Tenge
        self._check_tenge_keywords(line, line_num)
        
        # Проверка типов
        self._check_types(line, line_num)
        
        # Проверка точек с запятой
        self._check_semicolons(line, line_num)
        
        # Проверка отступов
        self._check_indentation(original_line, line_num)
        
        # Проверка казахского языка
        self._check_kazakh_compliance(line, line_num)
        
        # Проверка возвращаемых значений
        self._check_return_statements(line, line_num)

    def _check_agglutinative_patterns(self, line: str, line_num: int) -> None:
        """Проверка агглютинативных паттернов"""
        # Поиск функций
        function_matches = re.finditer(r'\b[a-zA-Z_][a-zA-Z0-9_]*\s*\(', line)
        for match in function_matches:
            function_name = match.group().replace('(', '').strip()
            if not self._is_agglutinative_name(function_name):
                self.errors.append(LintError(
                    line_num, match.start(),
                    f"Function name '{function_name}' should follow agglutinative patterns",
                    "warning", "tenge.agglutinative",
                    f"Consider using suffixes like _jasau, _alu, _qosu, etc."
                ))

    def _check_tenge_keywords(self, line: str, line_num: int) -> None:
        """Проверка использования правильных ключевых слов Tenge"""
        # Проверка на английские ключевые слова
        english_keywords = ['function', 'var', 'let', 'const', 'if', 'else', 'while', 'for', 'return']
        for keyword in english_keywords:
            if re.search(rf'\b{keyword}\b', line):
                self.errors.append(LintError(
                    line_num, 0,
                    f"Use Tenge keywords instead of '{keyword}'",
                    "error", "tenge.keywords",
                    f"Use 'atqar' instead of 'function', 'jasau' instead of 'var/let'"
                ))

    def _check_types(self, line: str, line_num: int) -> None:
        """Проверка типов"""
        # Проверка объявления переменных
        var_match = re.search(r'jasau\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*:', line)
        if var_match:
            var_name = var_match.group(1)
            if not any(line.count(f': {t}') for t in self.tenge_types):
                self.errors.append(LintError(
                    line_num, var_match.start(),
                    f"Variable '{var_name}' should have a proper type annotation",
                    "warning", "tenge.type",
                    "Use types like san, jol, aqıqat, JsonObject, etc."
                ))

    def _check_semicolons(self, line: str, line_num: int) -> None:
        """Проверка точек с запятой"""
        if (line and not line.endswith(';') and not line.endswith('{') and 
            not line.endswith('}') and not line.startswith('//') and 
            not line.startswith('/*') and not re.match(r'^(atqar|jasau|eger|azirshe|qaytar|end|import|export)', line)):
            self.errors.append(LintError(
                line_num, len(line),
                "Missing semicolon at end of statement",
                "warning", "tenge.semicolon",
                "Add semicolon at the end of the statement"
            ))

    def _check_indentation(self, line: str, line_num: int) -> None:
        """Проверка отступов"""
        if line.strip():
            # Проверка на смешанные табы и пробелы
            if '\t' in line and ' ' in line:
                self.errors.append(LintError(
                    line_num, 0,
                    "Mixed tabs and spaces in indentation",
                    "warning", "tenge.indentation",
                    "Use consistent indentation (either tabs or spaces)"
                ))

    def _check_kazakh_compliance(self, line: str, line_num: int) -> None:
        """Проверка соответствия казахскому языку"""
        for word in self.english_words:
            if re.search(rf'\b{word}\b', line, re.IGNORECASE):
                self.errors.append(LintError(
                    line_num, 0,
                    f"Use Kazakh equivalents instead of '{word}'",
                    "warning", "tenge.kazakh",
                    "Replace with agglutinative Kazakh morphemes"
                ))

    def _check_return_statements(self, line: str, line_num: int) -> None:
        """Проверка возвращаемых значений"""
        if line.strip().startswith('atqar ') and 'qaytar ' not in line:
            self.errors.append(LintError(
                line_num, 0,
                "Function should have a return statement",
                "warning", "tenge.return",
                "Add 'qaytar' statement to the function"
            ))

    def _is_agglutinative_name(self, name: str) -> bool:
        """Проверка агглютинативного имени"""
        return any(name.endswith(suffix) for suffix in self.agglutinative_suffixes)

    def format_output(self, errors: List[LintError], format_type: str = "text") -> str:
        """Форматирование вывода"""
        if format_type == "json":
            return json.dumps([
                {
                    "line": e.line,
                    "column": e.column,
                    "message": e.message,
                    "severity": e.severity,
                    "code": e.code,
                    "suggestion": e.suggestion
                } for e in errors
            ], indent=2)
        
        elif format_type == "text":
            output = []
            for error in errors:
                output.append(f"{error.severity.upper()}: {error.message}")
                output.append(f"  Line {error.line}, Column {error.column}")
                output.append(f"  Code: {error.code}")
                if error.suggestion:
                    output.append(f"  Suggestion: {error.suggestion}")
                output.append("")
            return "\n".join(output)
        
        return ""

def main():
    parser = argparse.ArgumentParser(description="Tenge Language Linter")
    parser.add_argument("files", nargs="+", help="Tenge files to lint")
    parser.add_argument("-f", "--format", choices=["text", "json"], default="text", help="Output format")
    parser.add_argument("-o", "--output", help="Output file")
    parser.add_argument("--exit-code", action="store_true", help="Exit with non-zero code if errors found")
    
    args = parser.parse_args()
    
    linter = TengeLinter()
    all_errors = []
    
    for file_path in args.files:
        if not Path(file_path).exists():
            print(f"Error: File '{file_path}' not found", file=sys.stderr)
            continue
        
        errors = linter.lint_file(file_path)
        all_errors.extend(errors)
    
    # Форматирование вывода
    output = linter.format_output(all_errors, args.format)
    
    if args.output:
        with open(args.output, 'w', encoding='utf-8') as f:
            f.write(output)
    else:
        print(output)
    
    # Выход с кодом ошибки
    if args.exit_code and all_errors:
        error_count = len([e for e in all_errors if e.severity == "error"])
        warning_count = len([e for e in all_errors if e.severity == "warning"])
        sys.exit(error_count + warning_count)

if __name__ == "__main__":
    main()
