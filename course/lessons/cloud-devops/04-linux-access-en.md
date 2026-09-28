# Linux: users, files, and least privilege

_Lead (summary):_ **Understand owners, groups, permission modes, and why a service should not run permanently as root.**

## Lesson outcome

Understand owners, groups, permission modes, and why a service should not run permanently as root. Do not move on until you can explain what every command verifies and which failure it can reveal.

## Why this matters

Least privilege limits blast radius: a compromised process receives only what its task needs. Modes `750` and `640` are not rituals; they encode read, write, and execute access for owner, group, and others.

## Practice

Run the commands from the `course/cloud-devops-lab` root. The commands target a learning environment. Read a command in full and check the current directory before it changes a server or data.

```shell
id
umask
mkdir -p /tmp/cloudlab-permissions
printf '%s\n' secret > /tmp/cloudlab-permissions/config
chmod 640 /tmp/cloudlab-permissions/config
ls -ld /tmp/cloudlab-permissions
ls -l /tmp/cloudlab-permissions/config
```

[Project reference files](https://github.com/DauletBai/shanraq.org/tree/main/course/cloud-devops-lab).

## Read the result

- The file is not executable.
- Other users have no access.
- You can decode every digit in mode 640.

![Support map 4](/static/course/cloud-devops/map-04-linux-access-en.svg)

## Exercise

**Required.** Write a POSIX sh script that creates a directory with mode 750 and a config file with mode 640, then verifies both modes with `stat`. Submit only POSIX sh commands or a script, with no real keys or tokens.

**Optional.** Write down one possible failure and the signal that would reveal it.

[Course contents](/course/cloud-devops?lang=en)
