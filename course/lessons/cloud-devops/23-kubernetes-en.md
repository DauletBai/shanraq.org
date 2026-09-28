# Kubernetes: Deployment, Service, probes, and rollout

_Lead (summary):_ **Read the CloudLab manifests, build them with Kustomize, and observe a safe update and rollback.**

## Lesson outcome

Read the CloudLab manifests, build them with Kustomize, and observe a safe update and rollback. Do not move on until you can explain what every command verifies and which failure it can reveal.

## Why this matters

A Deployment manages desired Pods and updates, a Service supplies a stable address, probes control traffic and restarts, and a ConfigMap separates public configuration. A PersistentVolumeClaim preserves data, but one JSON file does not become a distributed database, so this learning Deployment keeps one replica.

## Practice

Run the commands from the `course/cloud-devops-lab` root. The commands target a learning environment. Read a command in full and check the current directory before it changes a server or data.

```shell
kubectl kustomize deploy/k8s >/tmp/cloudlab-k8s.yaml
grep -n 'kind: Deployment[|]readinessProbe[|]runAsNonRoot[|]resources:' /tmp/cloudlab-k8s.yaml
printf 'With a local cluster: kubectl apply -k deploy/k8s\n'
printf 'Then: kubectl -n cloudlab rollout status deployment/cloudlab\n'
```

[Project reference files](https://github.com/DauletBai/shanraq.org/tree/main/course/cloud-devops-lab).

## Read the result

- The rendered YAML includes probes and a security context.
- The command exits with a predictable status.
- The result can be repeated in a clean environment.

![Support map 23](/static/course/cloud-devops/map-23-kubernetes-en.svg)

## Exercise

**Required.** Write a POSIX sh check that renders Kustomize YAML, requires a Deployment, Service, readinessProbe, and `runAsNonRoot: true`, then performs a client-side dry run. Submit only POSIX sh commands or a script, with no real keys or tokens.

**Optional.** Write down one possible failure and the signal that would reveal it.

[Course contents](/course/cloud-devops?lang=en)
