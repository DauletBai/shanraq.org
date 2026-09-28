# OpenTofu және Ansible: инфрақұрылым мен баптау код ретінде

_Лид (summary):_ **Ресурс сипаттамасын тексеріп, inventory құрып, хостқа идемпотентті баптау қолданамыз.**

## Сабақ нәтижесі

Ресурс сипаттамасын тексеріп, inventory құрып, хостқа идемпотентті баптау қолданамыз. Әр пәрмен нені тексеретінін және қандай ақауды табатынын түсіндіре алмайынша келесі сабаққа өтпеңіз.

## Бұл не үшін маңызды

OpenTofu state арқылы инфрақұрылым ресурсының өмір циклін басқарады; Ansible SSH арқылы ОС-ты қажетті күйге әкеледі. Код өзгерісті көрінетін етеді, бірақ state пен құпияны қорғау керек. Идемпотентті playbook қайта іске қосылған сайын өзгеріс жасамауы тиіс.

## Тәжірибе

Пәрмендерді `course/cloud-devops-lab` түбірінен орындаңыз. Пәрмендер оқу ортасына арналған. Серверді не деректі өзгертетін пәрменді іске қоспас бұрын оны толық оқып, ағымдағы қалтаны тексеріңіз.

```shell
tofu -chdir=infra/opentofu fmt -check 2>/dev/null || true
tofu -chdir=infra/opentofu init -backend=false 2>/dev/null || true
tofu -chdir=infra/opentofu validate 2>/dev/null || true
ansible-playbook --syntax-check -i infra/ansible/inventory.ini.example infra/ansible/site.yml 2>/dev/null || true
```

[Жобаның эталон файлдары](https://github.com/DauletBai/shanraq.org/tree/main/course/cloud-devops-lab).

## Нәтижені талдаңыз

- Тексеру бұлт ресурсын жасамайды.
- Пәрмен болжамды кодпен аяқталады.
- Нәтижені таза ортада қайталауға болады.

![Тірек сызба 22](/static/course/cloud-devops/map-22-iac-kz.svg)

## Тапсырма

**Міндетті.** `tofu` және `ansible-playbook` талап етіп, fmt/validate және syntax-check орындайтын, бірақ `apply` ешқашан шақырмайтын POSIX sh preflight жазыңыз. Тек POSIX sh пәрменін не скриптін, нақты кілт пен токенсіз жіберіңіз.

**Қалауыңызша.** Бір ықтимал ақауды және оны байқайтын сигналды жазыңыз.

[Курс мазмұны](/course/cloud-devops?lang=kz)
