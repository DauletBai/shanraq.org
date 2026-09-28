# Kubernetes: Deployment, Service, probe және rollout

_Лид (summary):_ **CloudLab манифесттерін оқып, Kustomize арқылы құрып, қауіпсіз жаңарту мен кері қайтаруды бақылаймыз.**

## Сабақ нәтижесі

CloudLab манифесттерін оқып, Kustomize арқылы құрып, қауіпсіз жаңарту мен кері қайтаруды бақылаймыз. Әр пәрмен нені тексеретінін және қандай ақауды табатынын түсіндіре алмайынша келесі сабаққа өтпеңіз.

## Бұл не үшін маңызды

Deployment қажетті Pod санын және жаңартуды, Service тұрақты мекенжайды, probe трафик пен қайта іске қосуды, ConfigMap ашық конфигурацияны басқарады. PersistentVolumeClaim деректі сақтайды, бірақ бір JSON файл таратылған база болмайды: оқу Deployment бір replica ұстайды.

## Тәжірибе

Пәрмендерді `course/cloud-devops-lab` түбірінен орындаңыз. Пәрмендер оқу ортасына арналған. Серверді не деректі өзгертетін пәрменді іске қоспас бұрын оны толық оқып, ағымдағы қалтаны тексеріңіз.

```shell
kubectl kustomize deploy/k8s >/tmp/cloudlab-k8s.yaml
grep -n 'kind: Deployment[|]readinessProbe[|]runAsNonRoot[|]resources:' /tmp/cloudlab-k8s.yaml
printf 'With a local cluster: kubectl apply -k deploy/k8s\n'
printf 'Then: kubectl -n cloudlab rollout status deployment/cloudlab\n'
```

[Жобаның эталон файлдары](https://github.com/DauletBai/shanraq.org/tree/main/course/cloud-devops-lab).

## Нәтижені талдаңыз

- Жиналған YAML ішінде probe және securityContext бар.
- Пәрмен болжамды кодпен аяқталады.
- Нәтижені таза ортада қайталауға болады.

![Тірек сызба 23](/static/course/cloud-devops/map-23-kubernetes-kz.svg)

## Тапсырма

**Міндетті.** Kustomize YAML құрып, Deployment, Service, readinessProbe және `runAsNonRoot: true` болуын талап етіп, client-side dry-run жасайтын POSIX sh тексеруін жазыңыз. Тек POSIX sh пәрменін не скриптін, нақты кілт пен токенсіз жіберіңіз.

**Қалауыңызша.** Бір ықтимал ақауды және оны байқайтын сигналды жазыңыз.

[Курс мазмұны](/course/cloud-devops?lang=kz)
