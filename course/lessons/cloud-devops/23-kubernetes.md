# Kubernetes: Deployment, Service, probes и rollout

_Лид (summary):_ **Читаем манифесты CloudLab, строим их через Kustomize и наблюдаем безопасное обновление и откат.**

## Результат урока

Читаем манифесты CloudLab, строим их через Kustomize и наблюдаем безопасное обновление и откат. Не переходите дальше, пока не можете объяснить, что проверяет каждая команда и какой отказ она обнаруживает.

## Почему это важно

Deployment управляет желаемым числом Pod и обновлением, Service даёт стабильный адрес, probes управляют трафиком и перезапуском, ConfigMap отделяет открытую конфигурацию. PersistentVolumeClaim сохраняет данные, но один локальный JSON-файл не становится распределённой базой: учебный Deployment оставляет одну реплику.

## Практика

Выполните команды из корня `course/cloud-devops-lab`. Команды рассчитаны на учебную среду. Перед командой, меняющей сервер или данные, прочитайте её целиком и проверьте текущий каталог.

```shell
kubectl kustomize deploy/k8s >/tmp/cloudlab-k8s.yaml
grep -n 'kind: Deployment[|]readinessProbe[|]runAsNonRoot[|]resources:' /tmp/cloudlab-k8s.yaml
printf 'With a local cluster: kubectl apply -k deploy/k8s\n'
printf 'Then: kubectl -n cloudlab rollout status deployment/cloudlab\n'
```

[Эталонные файлы проекта](https://github.com/DauletBai/shanraq.org/tree/main/course/cloud-devops-lab).

## Разберите результат

- Собранный YAML содержит probes и securityContext.
- Команда завершается предсказуемым кодом.
- Результат можно повторить на чистой среде.

![Опорная схема 23](/static/course/cloud-devops/map-23-kubernetes-ru.svg)

## Задание

**Обязательное.** Напишите POSIX sh-проверку, которая строит Kustomize YAML, требует наличие Deployment, Service, readinessProbe и `runAsNonRoot: true`, затем запускает client-side dry-run. Отправьте только команды или скрипт POSIX sh без реальных ключей и токенов.

**По желанию.** Запишите один возможный отказ и сигнал, по которому вы его заметите.

[Оглавление курса](/course/cloud-devops?lang=ru)
