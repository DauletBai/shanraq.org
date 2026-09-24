#!/usr/bin/env python3
"""Compile lesson snippets, check output/diagnostics and independent snapshots."""
import json
import os
from pathlib import Path
import re
import shutil
import subprocess
import tempfile

import langcheck

ROOT = Path(__file__).resolve().parents[2]
LESSONS = ROOT / 'course/lessons/rust'
FENCES = re.compile(r'^```([^\n]*)\n(.*?)^```\s*$', re.M | re.S)


def run(args, cwd):
    return subprocess.run(args, cwd=cwd, stdin=subprocess.DEVNULL, capture_output=True, text=True,
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
    if '#[test]' in code:
        tested = run(['rustc', '--edition=2024', '--test', str(source), '-o', str(binary)], directory)
        require(tested.returncode == 0, tested.stderr)
        result = run([str(binary)], directory)
        require(result.returncode == 0, result.stdout + result.stderr)


def dependency_manifest(page):
    """Return a lesson snapshot manifest only when its examples need Cargo deps."""
    number = page.stem[:2]
    lang = 'kz' if page.stem.endswith('-kz') else 'en' if page.stem.endswith('-en') else ''
    step = ROOT / 'course/rust-organizer' / lang / f'step-{number}'
    manifest = step / 'Cargo.toml'
    if not manifest.exists():
        return None
    after = manifest.read_text(encoding='utf-8').split('[dependencies]', 1)
    if len(after) < 2:
        return None
    deps = after[1].split('\n[', 1)[0]
    return manifest if any(line.strip() and not line.lstrip().startswith('#')
                           for line in deps.splitlines()) else None


def check_cargo_program(code, expected, directory, manifest, fails=False, diagnostic=None):
    """Compile a lesson with the dependencies pinned by its first Cargo snapshot."""
    project = directory / 'cargo-snippet'
    (project / 'src').mkdir(parents=True, exist_ok=True)
    shutil.copy2(manifest, project / 'Cargo.toml')
    shutil.copy2(manifest.with_name('Cargo.lock'), project / 'Cargo.lock')
    (project / 'src/main.rs').write_text(code, encoding='utf-8')
    # The checker shares Cargo's target directory across lessons. Clear only
    # this package's artifacts so identical package names cannot reuse a
    # binary from a different lesson; dependencies stay cached.
    cleaned = run(['cargo', 'clean', '-p', 'organizer'], project)
    require(cleaned.returncode == 0, cleaned.stderr)
    args = ['cargo', 'check' if fails else 'run', '--quiet', '--locked']
    result = run(args, project)
    if fails:
        require(result.returncode != 0, 'compile_fail unexpectedly compiled')
        require(diagnostic in result.stderr,
                f'missing diagnostic {diagnostic}: {result.stderr}')
        return
    require(result.returncode == 0, result.stderr)
    require(result.stdout == expected,
            f'output mismatch: {result.stdout!r} != {expected!r}')
    if '#[test]' in code:
        tested = run(['cargo', 'test', '--quiet', '--locked'], project)
        require(tested.returncode == 0, tested.stdout + tested.stderr)


def main():
    syllabus = json.loads((ROOT / 'docs/rust-syllabus.json').read_text(encoding='utf-8'))
    require([x['number'] for x in syllabus] == list(range(1, 61)), 'syllabus numbering')
    count = 0
    with tempfile.TemporaryDirectory(prefix='rust-course-') as name:
        directory = Path(name)
        os.environ.setdefault('CARGO_TARGET_DIR', str(directory / 'cargo-target'))
        # Check every local link, including the tables of contents and prefaces.
        for page in LESSONS.glob('*.md'):
            for link in re.findall(r'\]\(([^)]+)\)', page.read_text(encoding='utf-8')):
                if not re.match(r'https?://|#|/', link):
                    require((page.parent / link.split('#')[0]).exists(), f'{page}: missing {link}')
        release = json.loads((ROOT / 'tools/course/rust-release.json').read_text(encoding='utf-8'))
        for entry in syllabus:
            locales = [LESSONS / f"{entry['number']:02}-{entry['slug']}{suffix}.md"
                       for suffix in ('', '-kz', '-en')]
            if entry['number'] > release['lessons'] and not all(page.exists() for page in locales):
                continue
            command_sets = []
            for lang, suffix in [('ru', ''), ('kz', '-kz'), ('en', '-en')]:
                page = LESSONS / f"{entry['number']:02}-{entry['slug']}{suffix}.md"
                require(page.exists(), f'missing released locale: {page}')
                text = page.read_text(encoding='utf-8')
                words, reason = langcheck.check(str(page))
                require(not words, f'{page}: {reason}: {words}')
                notes = langcheck.check_kazakh_prose(str(page))
                require(not notes, f'{page}: {notes}')
                expected_heading = {'ru': '## Задание', 'kz': '## Тапсырма', 'en': '## Exercise'}[lang]
                require(expected_heading in text, f'{page}: missing exercise')
                command_sets.append([(m[1], m[2]) for m in FENCES.finditer(text)
                                     if m[1] in ('sh', 'powershell')])
                if entry['number'] >= 6:
                    answer = LESSONS / 'answers' / (page.stem + '-answer.rs')
                    inline = text.split('<!-- task-answer -->', 1)
                    require(len(inline) == 2, f'{page}: missing inline answer')
                    answer_blocks = list(FENCES.finditer(inline[1]))
                    require(len(answer_blocks) >= 2 and answer_blocks[0][1] == 'rust'
                            and answer_blocks[1][1] == 'text', f'{page}: answer format')
                    require(answer_blocks[0][2] == answer.read_text(encoding='utf-8'),
                            f'{page}: inline answer drift')
                    require(answer_blocks[1][2] == answer.with_suffix('.txt').read_text(encoding='utf-8'),
                            f'{page}: inline answer output drift')
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
            answer = LESSONS / 'answers' / (lesson.stem + '-answer.rs')
            if answer.exists():
                inline = text.split('<!-- task-answer -->', 1)
                require(len(inline) == 2, f'{lesson}: missing inline answer')
                answer_blocks = list(FENCES.finditer(inline[1]))
                require(len(answer_blocks) >= 2 and answer_blocks[0][1] == 'rust'
                        and answer_blocks[1][1] == 'text', f'{lesson}: answer format')
                require(answer_blocks[0][2] == answer.read_text(encoding='utf-8'),
                        f'{lesson}: inline answer drift')
                require(answer_blocks[1][2] == answer.with_suffix('.txt').read_text(encoding='utf-8'),
                        f'{lesson}: inline answer output drift')
            for link in re.findall(r'\]\(([^)]+)\)', text):
                if not re.match(r'https?://|#|/', link):
                    require((lesson.parent / link.split('#')[0]).exists(), f'{lesson}: missing {link}')
            blocks = list(FENCES.finditer(text))
            manifest = dependency_manifest(lesson)
            for index, match in enumerate(blocks):
                kind = match[1].strip()
                if kind not in ('rust', 'rust,compile_fail'):
                    continue
                try:
                    if kind.endswith('compile_fail'):
                        prefix = text[max(0, match.start()-100):match.start()]
                        codes = re.findall(r'error-code: (E\d+)', prefix)
                        checker = check_cargo_program if manifest else check_program
                        args = (match[2], '', directory)
                        if manifest:
                            checker(*args, manifest, True,
                                    codes[-1] if codes else 'cannot find macro')
                        else:
                            checker(*args, True,
                                    codes[-1] if codes else 'cannot find macro')
                    else:
                        require(index+1 < len(blocks) and blocks[index+1][1].strip() == 'text',
                                'runnable example must be followed by its text output')
                        if manifest:
                            check_cargo_program(match[2], blocks[index+1][2], directory, manifest)
                        else:
                            check_program(match[2], blocks[index+1][2], directory)
                    count += 1
                except AssertionError as error:
                    raise AssertionError(f'{lesson.name} block {index+1}: {error}') from error
        for answer in sorted((LESSONS / 'answers').glob('*.rs')):
            page = LESSONS / (answer.stem.removesuffix('-answer') + '.md')
            manifest = dependency_manifest(page)
            if manifest:
                check_cargo_program(answer.read_text(encoding='utf-8'),
                                    answer.with_suffix('.txt').read_text(encoding='utf-8'),
                                    directory, manifest)
            else:
                check_program(answer.read_text(encoding='utf-8'),
                              answer.with_suffix('.txt').read_text(encoding='utf-8'), directory)
            count += 1
        for step in sorted((ROOT / 'course/rust-organizer').rglob('step-*')):
            copied = directory / step.relative_to(ROOT / 'course/rust-organizer')
            shutil.copytree(step, copied, ignore=shutil.ignore_patterns('target'))
            suffix = '-' + step.parent.name if step.parent.name in ('kz', 'en') else ''
            entry = syllabus[int(step.name.removeprefix('step-')) - 1]
            lesson = LESSONS / f"{entry['number']:02}-{entry['slug']}{suffix}.md"
            code = next(m[2] for m in FENCES.finditer(lesson.read_text(encoding='utf-8')) if m[1]=='rust')
            require((step/'src/main.rs').read_text(encoding='utf-8') == code, f'{step}: lesson drift')
            cleaned = run(['cargo', 'clean', '-p', 'organizer'], copied)
            require(cleaned.returncode == 0, cleaned.stderr)
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
