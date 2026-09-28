# Cloud: IaaS, cost, IAM, and shared responsibility

_Lead (summary):_ **Design minimal virtual infrastructure, a budget, and access before creating a billable resource.**

## Lesson outcome

Design minimal virtual infrastructure, a budget, and access before creating a billable resource. Do not move on until you can explain what every command verifies and which failure it can reveal.

## Why this matters

In IaaS, the provider protects the facility and virtualization layer, while the customer owns the OS, access, firewall, data, and application. Free credit does not cancel later billing. Before creating anything, record region, size, disk, traffic, backups, and the deletion condition.

## Practice

Run the commands from the `course/cloud-devops-lab` root. The commands target a learning environment. Read a command in full and check the current directory before it changes a server or data.

```shell
cat > /tmp/cloud-budget.txt <<'EOF'
resource=one small Linux VM
owner=student
expires=after lesson 22
public_ports=22,80,443
budget_alert=set before creation
EOF
cat /tmp/cloud-budget.txt
printf 'No cloud resource is created by this lesson.\n'
```

[Project reference files](https://github.com/DauletBai/shanraq.org/tree/main/course/cloud-devops-lab).

## Read the result

- The plan has an owner and deletion date.
- The command exits with a predictable status.
- The result can be repeated in a clean environment.

![Support map 19](/static/course/cloud-devops/map-19-cloud-en.svg)

## Exercise

**Required.** Write a POSIX sh check for a plan file that requires `owner=`, `expires=`, and `budget_alert=` lines and fails if any is missing. Submit only POSIX sh commands or a script, with no real keys or tokens.

**Optional.** Write down one possible failure and the signal that would reveal it.

[Course contents](/course/cloud-devops?lang=en)
