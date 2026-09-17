#!/usr/bin/env python3
"""Compile lesson snippets, check output/diagnostics and independent snapshots."""
import json
import os
from pathlib import Path
import re
import shutil
import subprocess
import tempfile

ROOT = Path(__file__).resolve().parents[2]
LESSONS = ROOT / 'course/lessons/rust'
FENCES = re.compile(r'^```([^\n]*)\n(.*?)^```\s*$', re.M | re.S)


def run(args, cwd):
    return subprocess.run(args, cwd=cwd, capture_output=True, text=True,
                          encoding='utf-8', errors='replace', timeout=90)


def require(condition, message):
    if not condition:
        raise AssertionError(message)


def check_program(code, expected, directory, fails=False, diagnostic=None):
    source = directory / 'snippet.rs'
    source.write_text(code, encoding='utf-8')
    binary = directory / ('snippet.exe' if os.name == 'nt' else 'snippet')
    built = run(['rustc', '--edition=2024', str(source), '-o', str(binary)], directory)
    if fails:
        require(built.returncode != 0, 'compile_fail unexpectedly compiled')
        require(diagnostic in built.stderr, f'missing diagnostic {diagnostic}: {built.stderr}')
        return
    require(built.returncode == 0, built.stderr)
    executed = run([str(binary)], directory)
    require(executed.returncode == 0, executed.stderr)
    require(executed.stdout == expected,
            f'output mismatch: {executed.stdout!r} != {expected!r}')


def main():
    syllabus = json.loads((ROOT / 'docs/rust-syllabus.json').read_text(encoding='utf-8'))
    require([x['number'] for x in syllabus] == list(range(1, 61)), 'syllabus numbering')
    count = 0
    with tempfile.TemporaryDirectory(prefix='rust-course-') as name:
        directory = Path(name)
        # Check every local link, including the tables of contents and prefaces.
        for page in LESSONS.glob('*.md'):
            for link in re.findall(r'\]\(([^)]+)\)', page.read_text(encoding='utf-8')):
                if not re.match(r'https?://|#|/', link):
                    require((page.parent / link.split('#')[0]).exists(), f'{page}: missing {link}')
        for entry in syllabus[:5]:
            command_sets = []
            for lang, suffix in [('ru', ''), ('kz', '-kz'), ('en', '-en')]:
                page = LESSONS / f"{entry['number']:02}-{entry['slug']}{suffix}.md"
                require(page.exists(), f'missing first-batch locale: {page}')
                text = page.read_text(encoding='utf-8')
                expected_heading = {'ru': '## Задание', 'kz': '## Тапсырма', 'en': '## Exercise'}[lang]
                require(expected_heading in text, f'{page}: missing exercise')
                command_sets.append([(m[1], m[2]) for m in FENCES.finditer(text)
                                     if m[1] in ('sh', 'powershell')])
                if lang == 'kz':
                    prose = FENCES.sub('', text)
                    suspect = re.search(r'\b(?:переменная|переменные|настройки|запустите|выполните|следующий|предыдущий|папка|папки|задание|пользователь|исходный|ошибка|ошибки)\b', prose, re.I)
                    require(not suspect, f'{page}: Russian prose fragment: {suspect}')
            require(command_sets[0] == command_sets[1] == command_sets[2],
                    f"lesson {entry['number']}: translated commands drifted")
        probe = run(['cargo', 'new', 'rust-install-check', '--edition', '2024', '--vcs', 'none'], directory)
        require(probe.returncode == 0, probe.stderr)
        probe = run(['cargo', 'run', '--quiet', '--offline'], directory / 'rust-install-check')
        require(probe.returncode == 0 and probe.stdout == 'Hello, world!\n', probe.stderr + probe.stdout)
        count += 1
        for lesson in sorted(LESSONS.glob('[0-9][0-9]-*.md')):
            text = lesson.read_text(encoding='utf-8')
            require(text.startswith('# '), f'{lesson}: missing title')
            require(len(text.splitlines()) > 4 and '**' in text.splitlines()[2],
                    f'{lesson}: missing summary')
            for link in re.findall(r'\]\(([^)]+)\)', text):
                if not re.match(r'https?://|#|/', link):
                    require((lesson.parent / link.split('#')[0]).exists(), f'{lesson}: missing {link}')
            blocks = list(FENCES.finditer(text))
            for index, match in enumerate(blocks):
                kind = match[1].strip()
                if kind not in ('rust', 'rust,compile_fail'):
                    continue
                try:
                    if kind.endswith('compile_fail'):
                        prefix = text[max(0, match.start()-100):match.start()]
                        codes = re.findall(r'error-code: (E\d+)', prefix)
                        check_program(match[2], '', directory, True,
                                      codes[-1] if codes else 'cannot find macro')
                    else:
                        require(index+1 < len(blocks) and blocks[index+1][1].strip() == 'text',
                                'runnable example must be followed by its text output')
                        check_program(match[2], blocks[index+1][2], directory)
                    count += 1
                except AssertionError as error:
                    raise AssertionError(f'{lesson.name} block {index+1}: {error}') from error
        for answer in sorted((LESSONS / 'answers').glob('*.rs')):
            check_program(answer.read_text(encoding='utf-8'),
                          answer.with_suffix('.txt').read_text(encoding='utf-8'), directory)
            count += 1
        for step in sorted((ROOT / 'course/rust-organizer').glob('step-*')):
            copied = directory / step.name
            shutil.copytree(step, copied, ignore=shutil.ignore_patterns('target'))
            # Select the Russian base file explicitly, not a translated file.
            lesson = next(p for p in LESSONS.glob(step.name.removeprefix('step-')+'-*.md')
                          if not p.stem.endswith(('-kz', '-en')))
            code = next(m[2] for m in FENCES.finditer(lesson.read_text(encoding='utf-8')) if m[1]=='rust')
            require((step/'src/main.rs').read_text(encoding='utf-8') == code, f'{step}: lesson drift')
            result = run(['cargo', 'run', '--quiet', '--locked', '--offline'], copied)
            require(result.returncode == 0, result.stderr)
            require(result.stdout == (step/'expected.txt').read_text(encoding='utf-8'), f'{step}: output drift')
            formatted = run(['cargo', 'fmt', '--check'], copied)
            require(formatted.returncode == 0, formatted.stdout + formatted.stderr)
            count += 1
        # Verify the checker rejects wrong output and a fake compile-fail fixture.
        for code, output, fails, diagnostic in [
            ('fn main() { println!("actual"); }', 'wrong\n', False, None),
            ('fn main() {}', '', True, 'E0384'),
            ('fn main() { missing(); }', '', True, 'E0384'),
        ]:
            try:
                check_program(code, output, directory, fails, diagnostic)
            except AssertionError:
                pass
            else:
                raise AssertionError('checker accepted a deliberately invalid fixture')
    print(f'PASS: {count} Rust examples, answers and Cargo snapshots; checker negative fixtures')


if __name__ == '__main__':
    main()
