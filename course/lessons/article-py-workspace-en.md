# A workplace: Python, the project's environment and the first run

_Лид (summary):_ **The second lesson of the Python course. We install the language, create the project's environment and move yesterday's program into it. Measured: a package installed into the environment is invisible outside it — `ModuleNotFoundError` — which is exactly what an environment is for. The `.venv` folder weighs 17 MB and never goes into the repository.**

## Why this is needed

Yesterday's program ran on bare Python: nothing had to be installed. But within a few lessons we will need other people's libraries — for requests, for tables, for charts. And this is where a beginner meets the first real unpleasantness.

A library is installed into the system. Then another, for another task. Then they are upgraded. Then something that worked for half a year stops starting, and it is not clear what exactly did it. On other people's machines this story ends with reinstalling the system.

One habit cures it: **every project gets an environment of its own**. Today we make one and make sure it really is separate.

## The whole thing at once

Everything below is the terminal. Type it line by line; the `$` is the shell's prompt and is not typed.

```
$ python3 --version
Python 3.14.5

$ python3 -m venv .venv          # create the project's environment
$ source .venv/bin/activate
$ python -c "import sys, pathlib; print(pathlib.Path(sys.executable).relative_to(pathlib.Path.cwd()))"
.venv/bin/python

$ pip install requests
Successfully installed certifi-2026.7.22 charset_normalizer-3.5.1 idna-3.19 requests-2.34.2 urllib3-2.7.0

$ python -c "import requests; print(requests.__version__)"
2.34.2

$ pip freeze > requirements.txt
$ cat requirements.txt
certifi==2026.7.22
charset-normalizer==3.5.1
idna==3.19
requests==2.34.2
urllib3==2.7.0

$ deactivate                     # leave the environment
$ python3 -c "import requests"
Traceback (most recent call last):
  File "<string>", line 1, in <module>
    import requests
ModuleNotFoundError: No module named 'requests'

$ du -sh .venv
 17M    .venv
```

## The walk-through

### First check that the language is there at all

`python3 --version` is the first command of any day. Answered with a number — the language is installed. Answered "command not found" — install it:

- **Windows:** open python.org, download the installer and **be sure** to tick "Add python.exe to PATH" on the first screen. Forgot it? Reinstall; that is quicker than repairing.
- **macOS:** the simplest way is Homebrew: `brew install python`. Better to leave the Python already in the system alone: the system itself uses it.
- **Linux:** it is almost certainly there already; if not, `sudo apt install python3 python3-venv` or the same through your package manager.

Any version from 3.12 will do. The course is written and checked on 3.14.

### What an environment is and why it lives in the project

`python3 -m venv .venv` creates a `.venv` folder inside the project. Inside it are a `python` of its own and a place of its own for libraries.

The `.venv/bin/python` line in the output confirms it: after `activate`, the word `python` means not the system language but the one lying in your project. Nothing was "configured globally" — only a path changed.

> **Picture it.** A first-aid kit of your own in the car rather than the shared one in the yard. In the shared one the bandages run out, or somebody leaves something expired; yours is always the way you packed it.

### The proof this was all for

Look at two lines of the output in a row.

Inside the environment `import requests` works and reports version `2.34.2`. Leave the environment with `deactivate` and the same import gives `ModuleNotFoundError: No module named 'requests'`.

That is not a failure. That is the isolation: the library went **into the project**, not into the system. There is nothing left to break a neighbouring project with.

### `requirements.txt`, the list an environment is rebuilt from

`pip freeze` prints everything installed in the environment with exact versions:

```
requests==2.34.2
urllib3==2.7.0
```

We saved that into `requirements.txt`. A year from now, on another machine, `pip install -r requirements.txt` assembles the same thing. Notice: you asked for `requests` alone and the list has five lines — the rest came with it, and that is normal.

The `.venv` itself does **not** go into the repository: 17 megabytes, different on every system, and rebuilt from the list with one command. That is why the project now has a two-line `.gitignore`.

### What the project looks like now

```
$ ls -a
.gitignore .venv requirements.txt tsena.py 

$ source .venv/bin/activate
$ python tsena.py | head -5
== how many times prices grew in Kazakhstan
price index in 2010: 100.0
price index in 2025: 348.1
prices grew 3.48 times
1000 tenge of 2010 is worth 287 tenge today
```

Four names, each in its place: the code, the list of dependencies, the environment, and the rule about what stays out of the repository. Yesterday's program runs from the project and prints the same thing.

This folder will grow — modules, a database and a report will appear in it. But the structure was set today and will not have to change.

### The editor

Any will do, but if you are choosing — **VS Code** with Microsoft's **Python** extension. One setting matters: in the bottom right corner select the interpreter, the one inside `.venv`. Then the editor suggests from the libraries you actually have rather than from ones it imagined.

Checking is easy: open `tsena.py` and press "Run" — the program should behave exactly as it does from the terminal.

### One warning about `pip`

Better never to type `pip install` **outside an environment**. If you are unsure where you are, look at the prompt: an activated environment puts `(.venv)` on the left. No `(.venv)`? Then `source .venv/bin/activate` first (on Windows: `.venv\Scripts\activate`).

## The map of the lesson

![The map of the lesson: the system, the environment and the project](/static/course/py/map-workspace-en.svg)

## Say it in your own words

Without looking, answer aloud or on paper. The answers are at the end of the lesson.

1. What exactly does `activate` change, if nothing is installed anew?
2. Why is `.venv` kept out of the repository while `requirements.txt` is kept in it?
3. How can you tell in one second that you are outside the environment?

## The exercise

**Required.** Make a `digest` folder, create an environment in it, move yesterday's program there and run it from inside the environment. Then repeat the check: install `requests`, leave the environment, and make sure the import no longer works.

**If you want more.**

- Look inside `.venv`: `ls .venv/lib/python3*/site-packages`. Exactly those five packages from the list are there.
- Delete the whole `.venv` folder and rebuild it: `python3 -m venv .venv && source .venv/bin/activate && pip install -r requirements.txt`. That is how you check the list is complete.
- Make a second project with a different version of `requests` and see that the first one is unchanged.

## Where this goes in the project

The digest has a home. From now on everything we write lives in the `digest` folder: today one program and a list of dependencies, by the end of the course modules, a database, charts and a report.

The debts. The program still dies without the internet and fetches the data afresh on every run. We save nothing to disk yet — and data that was fetched but not saved has to be asked for again. Those are the next lessons.

## The answers

1. Only the paths. `activate` puts `.venv/bin` first in `PATH`, so the word `python` starts meaning the interpreter from the project and `pip` starts installing packages there. Nothing is installed in the process.
2. Because `.venv` is a result: rebuilt with one command, 17 megabytes, and system-dependent. `requirements.txt` is the cause: five lines from which the result follows on any machine.
3. Look at the terminal prompt: an activated environment shows `(.venv)` on the left. If in doubt, `python -c "import sys; print(sys.executable)"` shows whose python this is.

## Sources

- [Python: the venv module](https://docs.python.org/3/library/venv.html)
- [Python: installing packages](https://packaging.python.org/en/latest/tutorials/installing-packages/)
- [VS Code: working with Python environments](https://code.visualstudio.com/docs/python/environments)
