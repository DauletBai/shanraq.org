# OpenTofu and Ansible: infrastructure and configuration as code

_Lead (summary):_ **Validate a resource description, produce inventory, and apply idempotent host configuration.**

## Lesson outcome

Validate a resource description, produce inventory, and apply idempotent host configuration. Do not move on until you can explain what every command verifies and which failure it can reveal.

## Why this matters

OpenTofu manages infrastructure resource lifecycles through state; Ansible brings the OS to the desired configuration over SSH. Code makes changes reviewable, but state and secrets need protection. An idempotent playbook should report no changes on a second run.

## Practice

Run the commands from the `course/cloud-devops-lab` root. The commands target a learning environment. Read a command in full and check the current directory before it changes a server or data.

```shell
tofu -chdir=infra/opentofu fmt -check 2>/dev/null || true
tofu -chdir=infra/opentofu init -backend=false 2>/dev/null || true
tofu -chdir=infra/opentofu validate 2>/dev/null || true
ansible-playbook --syntax-check -i infra/ansible/inventory.ini.example infra/ansible/site.yml 2>/dev/null || true
```

[Project reference files](https://github.com/DauletBai/shanraq.org/tree/main/course/cloud-devops-lab).

## Read the result

- Validation creates no cloud resource.
- The command exits with a predictable status.
- The result can be repeated in a clean environment.

![Support map 22](/static/course/cloud-devops/map-22-iac-en.svg)

## Exercise

**Required.** Write a POSIX sh preflight that requires `tofu` and `ansible-playbook`, runs fmt/validate and syntax-check, and never calls `apply`. Submit only POSIX sh commands or a script, with no real keys or tokens.

**Optional.** Write down one possible failure and the signal that would reveal it.

[Course contents](/course/cloud-devops?lang=en)
