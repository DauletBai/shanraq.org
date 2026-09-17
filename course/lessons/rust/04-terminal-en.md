# Understanding the terminal

_Summary:_ **Enter a folder whose name contains spaces.**


## Why this matters

A **terminal** is a window for text interaction. A **shell** interprets the commands entered there: PowerShell on Windows, usually zsh on macOS and often bash on Linux. Think of a service window and the dispatcher behind it. The terminal does not compile Rust source code for you.

## The whole picture

Open PowerShell from the Windows menu, or Terminal on macOS/Linux. The **prompt** shows that the shell is waiting for a command. Do not copy prompt decorations such as $ or PS C:\...>. Our command blocks omit them.

Find your current folder. In PowerShell:

```powershell
Get-Location
```

On macOS/Linux:

```sh
pwd
```

The current folder is the starting point for relative paths. Enter the folder from lesson 3 and list its contents. PowerShell:

```powershell
cd "$HOME/rust-learning"
Get-ChildItem
```

macOS/Linux:

```sh
cd "$HOME/rust-learning"
ls
```

cd changes directory. The following value is an **argument**, extra information supplied to the command. $HOME refers to your home folder in these shells; do not assign a new value to it. Double quotes keep the path one argument even when it contains spaces. Get-ChildItem and ls list files; notes.txt should be present.

## Commands and code are different

Entering fn main() into the shell does not create a Rust program. Source goes in .rs files; tool commands go in the terminal. Enter submits the command. Do not append the printed output to the command itself.

For a normally running terminal program, Ctrl+C requests interruption. The exit command leaves the shell itself and may close its window or tab. Those are different actions.

If cd reports a missing folder, check its name and location from lesson 3. Running as administrator does not fix the wrong address.

## Recall map

Terminal → shell → command and arguments → action → prompt. The current folder determines where relative paths start.

## Warm-up

1. Predict: does opening a file in an editor automatically change the shell's folder?
2. Complete: cd receives a path as an ___.
3. Repair a path containing spaces by quoting the whole path.

## Exercise

**Required.** Enter rust-learning and list its files. Run cd .., then show your current folder again. The two dots mean the parent folder. Expect rust-learning with notes.txt first, then your home folder.

**Your own data.** Create “My tasks” inside rust-learning using the file manager, enter it with a quoted path and return using cd .. .

## Hint and reference answer

Check location after each change using Get-Location or pwd. The path ends in rust-learning before leaving, and no longer does afterwards. Account and drive names vary. Continue when you can follow just your operating system's branch and distinguish commands from output.

[Previous lesson](/read/rust-03-workspace?lang=en) · [Contents](/course/rust?lang=en) · [Next lesson](/read/rust-05-installation?lang=en)
