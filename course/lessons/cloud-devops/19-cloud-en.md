# Cloud: IaaS, cost, IAM, and shared responsibility

_Lead (summary):_ **Design minimal virtual infrastructure, a budget, and access before creating a billable resource.**

> This lesson starts from zero. An unfamiliar word is not treated as your deficiency: we first build an image, then give it a precise meaning, and only then use it.

## Where we are

Flight review. The previous lesson used the support **deploy X → smoke test → observe → keep X OR rollback Y → verify**. Without looking, name its four links, then compare with the written signal.

## Lesson outcome

Design minimal virtual infrastructure, a budget, and access before creating a billable resource. Do not move on until you can explain what every command verifies and which failure it can reveal.

## Begin with a familiar image

Cloud resembles renting a furnished office: it is available quickly, but time, size, and extras cost money. The landlord protects the building; you manage keys, people, documents, and anything left exposed. This is shared responsibility.

An analogy starts reasoning but does not replace the technology. We make every word precise below.

## The lesson's support signal

**need → region/size/access → price alert → create → observe → delete**

Read the chain left to right. It answers not “which command should I memorise?” but “why does each step follow the previous one?”

## New words in plain language

- **IaaS** — Rental of virtual compute, networking, and disks while you manage the OS and application.
- **Region** — A geographic area in which cloud resources run.
- **IAM** — Rules defining who may perform which action on a resource.
- **Budget alert** — A spending warning; it does not necessarily stop resources automatically.

## Take it apart without rushing

In IaaS, the provider protects the facility and virtualization layer, while the customer owns the OS, access, firewall, data, and application. Free credit does not cancel later billing. Before creating anything, record region, size, disk, traffic, backups, and the deletion condition.

Do not memorise this wording. Find its cause-and-effect link and attach it to the support signal above.

## Predict the result first

Before running anything, write what should change, what should remain unchanged, and which output would support your prediction. A wrong prediction is useful when you can explain the difference after observing it.

## Experiment in small steps

Run the commands from the `course/cloud-devops-lab` root. The commands target a learning environment. Read a command in full and check the current directory before it changes a server or data. Run one meaningful group at a time. Stop after each group and compare evidence with your prediction.

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

### What this experiment actually does

The lesson first creates a text plan rather than a billable resource. It records an owner, deletion date, open ports, and a budget alert. Only a understood plan becomes later infrastructure.

## What you should observe

- The plan has an owner and deletion date.
- The command exits with a predictable status.
- The result can be repeated in a clean environment.

## Rebuild the whole from the support map

![Support map 19](/static/course/cloud-devops/map-19-cloud-en.svg)

Follow the four blocks with a finger or pencil and tell the entire lesson as one connected explanation.

## Recall without a prompt

1. Hide the map and draw its four blocks from memory.
2. Explain every transition using the word “because”.
3. Define one new term without repeating the text verbatim.
4. Name one signal that proves the result.

## Find and correct the mistake

Mistake: “Free credit means a bill is impossible.” Correction: limits, dates, and terms vary; set an alert and deletion condition before creation.

## Exercise

Now rebuild the same meaning yourself. You may use the support map; the goal is a repeatable result rather than a memory exam.

**Required.** Write a POSIX sh check for a plan file that requires `owner=`, `expires=`, and `budget_alert=` lines and fails if any is missing. Submit only POSIX sh commands or a script, with no real keys or tokens.

### Completion criteria

- You understand every line you run.
- A repeat gives the expected result or a safely explained difference.
- Output contains no real keys, tokens, or passwords.
- You can point to evidence that proves success.

If a check fails, record expected and actual results, revisit the one broken support link, and retry. An error does not remove your right to continue learning.

**Optional.** Write down one possible failure and the signal that would reveal it.

## Open perspective

The next lesson adds a new support: **Ubuntu server: SSH keys, updates, and a firewall**.

[Course contents](/course/cloud-devops?lang=en)
