# Files, folders and your workspace

_Summary:_ **Find the saved file without confusing it with another copy.**


## Why this matters

Before coding, learn to locate what you edit. A **file** is named data; a **folder** groups files; an **extension** is an ending such as .txt or .rs. Think of documents in folders. Renaming a document does not translate its contents: replacing .txt with .rs does not turn ordinary prose into valid Rust.

## The whole picture

Create a folder named rust-learning inside your home folder. The home folder is your personal location: commonly C:\Users\name on Windows, /Users/name on macOS or /home/name on Linux. Replace “name” with your account name; these are examples, not commands to paste.

Inside it, create notes.txt containing:

```text
My project is an organizer. First task: learn Rust.
```

Use a plain-text editor, not a Word document. Windows Notepad is suitable; in macOS TextEdit select plain-text mode; on Linux use your installed text editor. Lesson 5 sets up a code editor. Save as UTF-8, a way of representing text as bytes that supports different languages. We will explain bytes later; for now choose UTF-8 if the save dialog asks.

## A path is an address

An **absolute path** starts from the filesystem root or drive and gives the full location. A **relative path** starts from the current folder. Think of a full postal address versus “the next room”: the second only makes sense if your starting point is known.

notes.txt inside rust-learning is a relative path. Two files named notes.txt in different folders are different files. On Windows enable file-name extensions in File Explorer's View menu; wording varies with the Windows version. Make sure you did not create notes.txt.txt. On other systems check the full name in file properties.

**Saving** writes your edits to the file. Text visible in an editor tab may not be saved yet. Close the file, reopen it and check the line: this is our first small verification of an action.

## Recall map

File: name and contents. Folder: location. Path: address. Saving: edits reach the file. An extension does not prove that the contents are correct.

## Warm-up

1. Predict: are two files named notes.txt in different folders the same file?
2. Complete: a relative path depends on the ___ folder.
3. Repair notes.txt.txt so it has one .txt ending.

## Exercise

**Required.** Create notes.txt, save the sentence above, close it and reopen it. Record its full path separately. Expected result: the sentence persists and you can find that exact file again.

**Your own data.** Add “Менің тапсырмам” and check that the letters remain readable.

## Hint and reference answer

If text disappears, check saving and the folder before creating another copy. Reference: one notes.txt in rust-learning; both lines survive reopening. Each reader's full path differs and need not match the author's path.

Continue when you can point out the file name, its folder and its contents separately.

[Previous lesson](/read/rust-02-organizer?lang=en) · [Contents](/course/rust?lang=en) · [Next lesson](/read/rust-04-terminal?lang=en)
