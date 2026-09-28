# Ubuntu server: SSH keys, updates, and a firewall

_Lead (summary):_ **Connect as an unprivileged user, close unnecessary ports, and verify access from a second session.**

## Lesson outcome

Connect as an unprivileged user, close unnecessary ports, and verify access from a second session. Do not move on until you can explain what every command verifies and which failure it can reveal.

## Why this matters

Harden a server in a safe order: create a user and key, verify a second login, then disable risky access. Enabling a firewall or disabling root too early can lock you out. The provider console remains an emergency path.

## Practice

Run the commands from the `course/cloud-devops-lab` root. The commands target a learning environment. Read a command in full and check the current directory before it changes a server or data.

```shell
ssh-keygen -t ed25519 -a 64 -f ~/.ssh/cloudlab_ed25519 -C cloudlab
printf 'Copy only cloudlab_ed25519.pub to the server.\n'
printf 'On Ubuntu: allow OpenSSH before enabling ufw.\n'
printf 'Keep the current session open while testing a second login.\n'
```

[Project reference files](https://github.com/DauletBai/shanraq.org/tree/main/course/cloud-devops-lab).

## Read the result

- The private key is never copied to the server.
- The command exits with a predictable status.
- The result can be repeated in a clean environment.

![Support map 20](/static/course/cloud-devops/map-20-server-en.svg)

## Exercise

**Required.** Write a local POSIX sh script that checks for private and public key files, requires mode 600 on the private key, and prints only the public-key fingerprint. Submit only POSIX sh commands or a script, with no real keys or tokens.

**Optional.** Write down one possible failure and the signal that would reveal it.

[Course contents](/course/cloud-devops?lang=en)
