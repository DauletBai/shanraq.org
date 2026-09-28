# Kubernetes: Deployment, Service, probes, and rollout

_Lead (summary):_ **Read the CloudLab manifests, build them with Kustomize, and observe a safe update and rollback.**

> This lesson starts from zero. An unfamiliar word is not treated as your deficiency: we first build an image, then give it a precise meaning, and only then use it.

## Where we are

Flight review. The previous lesson used the support **desired code → plan/diff → approve → apply → state; repeat → no surprise**. Without looking, name its four links, then compare with the written signal.

## Lesson outcome

Read the CloudLab manifests, build them with Kustomize, and observe a safe update and rollback. Do not move on until you can explain what every command verifies and which failure it can reveal.

## Begin with a familiar image

Imagine a bakery chain. A central plan requires three operating shops: the Deployment. Each shop is a Pod. One ordering number is the Service. “Can this shop take orders?” is a readiness probe. Gradual replacement of old shops is a rollout.

An analogy starts reasoning but does not replace the technology. We make every word precise below.

## The lesson's support signal

**Deployment wants Pods → Service finds ready Pods → probes steer traffic → rollout changes gradually**

Read the chain left to right. It answers not “which command should I memorise?” but “why does each step follow the previous one?”

## New words in plain language

- **Kubernetes** — A system that continually moves container applications toward declared state.
- **Pod** — The smallest runnable group of one or more tightly related containers.
- **Deployment** — A controller for desired Pod count and update strategy.
- **Service и probe** — A Service supplies a stable entry point; a probe reports which Pod is ready or alive.

## Take it apart without rushing

A Deployment manages desired Pods and updates, a Service supplies a stable address, probes control traffic and restarts, and a ConfigMap separates public configuration. A PersistentVolumeClaim preserves data, but one JSON file does not become a distributed database, so this learning Deployment keeps one replica.

Do not memorise this wording. Find its cause-and-effect link and attach it to the support signal above.

## Predict the result first

Before running anything, write what should change, what should remain unchanged, and which output would support your prediction. A wrong prediction is useful when you can explain the difference after observing it.

## Experiment in small steps

Run the commands from the `course/cloud-devops-lab` root. The commands target a learning environment. Read a command in full and check the current directory before it changes a server or data. Run one meaningful group at a time. Stop after each group and compare evidence with your prediction.

```shell
kubectl kustomize deploy/k8s >/tmp/cloudlab-k8s.yaml
grep -n 'kind: Deployment[|]readinessProbe[|]runAsNonRoot[|]resources:' /tmp/cloudlab-k8s.yaml
printf 'With a local cluster: kubectl apply -k deploy/k8s\n'
printf 'Then: kubectl -n cloudlab rollout status deployment/cloudlab\n'
```

[Project reference files](https://github.com/DauletBai/shanraq.org/tree/main/course/cloud-devops-lab).

### What this experiment actually does

`kubectl kustomize` first renders final YAML locally. A search checks key parts, and client-side dry-run asks Kubernetes to parse objects without creating them. Apply only to a learning cluster after understanding the manifest.

## What you should observe

- The rendered YAML includes probes and a security context.
- The command exits with a predictable status.
- The result can be repeated in a clean environment.

## Rebuild the whole from the support map

![Support map 23](/static/course/cloud-devops/map-23-kubernetes-en.svg)

Follow the four blocks with a finger or pencil and tell the entire lesson as one connected explanation.

## Recall without a prompt

1. Hide the map and draw its four blocks from memory.
2. Explain every transition using the word “because”.
3. Define one new term without repeating the text verbatim.
4. Name one signal that proves the result.

## Find and correct the mistake

Mistake: “Kubernetes automatically makes every application fault tolerant.” Correction: it executes declared rules, but bad probes, a single disk, or application defects remain.

## Exercise

Now rebuild the same meaning yourself. You may use the support map; the goal is a repeatable result rather than a memory exam.

**Required.** Write a POSIX sh check that renders Kustomize YAML, requires a Deployment, Service, readinessProbe, and `runAsNonRoot: true`, then performs a client-side dry run. Submit only POSIX sh commands or a script, with no real keys or tokens.

### Completion criteria

- You understand every line you run.
- A repeat gives the expected result or a safely explained difference.
- Output contains no real keys, tokens, or passwords.
- You can point to evidence that proves success.

If a check fails, record expected and actual results, revisit the one broken support link, and retry. An error does not remove your right to continue learning.

**Optional.** Write down one possible failure and the signal that would reveal it.

## Open perspective

The next lesson adds a new support: **Observability, backups, and a recovery drill**.

[Course contents](/course/cloud-devops?lang=en)
