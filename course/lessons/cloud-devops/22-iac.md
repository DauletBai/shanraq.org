# OpenTofu и Ansible: инфраструктура и настройка как код

_Лид (summary):_ **Валидируем описание ресурса, формируем inventory и применяем идемпотентную настройку хоста.**

## Результат урока

Валидируем описание ресурса, формируем inventory и применяем идемпотентную настройку хоста. Не переходите дальше, пока не можете объяснить, что проверяет каждая команда и какой отказ она обнаруживает.

## Почему это важно

OpenTofu управляет жизненным циклом инфраструктурных ресурсов через state; Ansible приводит ОС к нужной конфигурации по SSH. Код делает изменение обозримым, но state и секреты требуют защиты. Идемпотентный playbook при повторном запуске не должен каждый раз сообщать об изменении.

## Практика

Выполните команды из корня `course/cloud-devops-lab`. Команды рассчитаны на учебную среду. Перед командой, меняющей сервер или данные, прочитайте её целиком и проверьте текущий каталог.

```shell
tofu -chdir=infra/opentofu fmt -check 2>/dev/null || true
tofu -chdir=infra/opentofu init -backend=false 2>/dev/null || true
tofu -chdir=infra/opentofu validate 2>/dev/null || true
ansible-playbook --syntax-check -i infra/ansible/inventory.ini.example infra/ansible/site.yml 2>/dev/null || true
```

[Эталонные файлы проекта](https://github.com/DauletBai/shanraq.org/tree/main/course/cloud-devops-lab).

## Разберите результат

- Проверка не создаёт облачный ресурс.
- Команда завершается предсказуемым кодом.
- Результат можно повторить на чистой среде.

![Опорная схема 22](/static/course/cloud-devops/map-22-iac-ru.svg)

## Задание

**Обязательное.** Напишите POSIX sh preflight, который требует `tofu` и `ansible-playbook`, выполняет fmt/validate и syntax-check, но никогда не вызывает `apply`. Отправьте только команды или скрипт POSIX sh без реальных ключей и токенов.

**По желанию.** Запишите один возможный отказ и сигнал, по которому вы его заметите.

[Оглавление курса](/course/cloud-devops?lang=ru)
