"""Verify meaningful boundary variations promised in Rust lessons 11–35."""
import os,subprocess,sys,tempfile,re
from pathlib import Path
sys.path.insert(0, str(Path(__file__).resolve().parent))
from rustcheck import check_program,FENCES
R=Path(__file__).resolve().parents[2] / 'course/lessons/rust'; count=0
with tempfile.TemporaryDirectory() as name:
 for lang,suffix,left,done,day,total,short,invalid in [
 ('ru','','Осталось','Все задачи выполнены','День','Всего минут','Коротких задач','Ошибка: выполнено больше общего числа'),
 ('kz','-kz','Қалғаны','Барлық тапсырма аяқталды','Күн','Барлық минут','Қысқа тапсырмалар','Қате: аяқталғаны жалпы саннан көп'),
 ('en','-en','Remaining','All tasks completed','Day','Total minutes','Short tasks','Error: completed exceeds total')]:
  def check(stem,old,new,expected):
   global count
   s=(R/'answers'/f'{stem}{suffix}-answer.rs').read_text(encoding='utf-8');assert s.count(old)==1,(stem,old)
   check_program(s.replace(old,new),expected,Path(name));count+=1
  for old,new in [('let done: bool = false','let done: bool = true'),('let archived: bool = false','let archived: bool = true'),('minutes: u32 = 15','minutes: u32 = 16')]:
   label={'ru':'Показать','kz':'Көрсету','en':'Show'}[lang];check('11-boolean',old,new,label+': false\n')
  check('12-branches','done: u32 = 6','done: u32 = 4',done+'\n')
  check('12-branches','done: u32 = 6','done: u32 = 1',left+': 3\n')
  check('13-loops','1..=5','1..1',total+': 0\n')
  check('14-functions','is_short(0)','is_short(1)','true\ntrue\nfalse\n')
  check('15-arrays','[0, 10, 25, 15]','[0, 16, 25, 30]',short+': 0\n'+total+': 0\n')
  check('15-arrays','[0, 10, 25, 15]','[1, 15, 2, 3]',short+': 4\n'+total+': 21\n')
  text=(R/f'15-arrays{suffix}.md').read_text(encoding='utf-8')
  s=next(m[2] for m in FENCES.finditer(text) if m[1]=='rust' and 'fn print_at' in m[2])
  check_program(s.replace('print_at(minutes, 3);','print_at(minutes, 1);'),'25\n',Path(name));count+=1
  first,last,elements={
   'ru':('Первые','Последние','Всего элементов'),
   'kz':('Алғашқылары','Соңғылары','Барлық элемент'),
   'en':('First pair','Last pair','Element count'),
  }[lang]
  check('23-slices','&minutes[0..2]','&minutes[0..0]',
        first+': 0\n'+last+': 35\n'+elements+': 4\n')
  removed,remaining,count_label={
   'ru':('Чтение',('Прогулка','Отдых'),'Осталось'),
   'kz':('Оқу',('Серуен','Демалыс'),'Қалды'),
   'en':('Reading',('Walk','Rest'),'Remaining'),
  }[lang]
  removed_label={'ru':'Удалено','kz':'Өшірілді','en':'Removed'}[lang]
  task_label={'ru':'Задача','kz':'Тапсырма','en':'Task'}[lang]
  check('24-vectors','titles.remove(1)','titles.remove(0)',
        removed_label+': '+removed+'\n'+task_label+': '+remaining[0]+'\n'
        +task_label+': '+remaining[1]+'\n'+count_label+': 2\n')
  tuple_code=(R/'answers'/f'26-tuples{suffix}-answer.rs').read_text(encoding='utf-8')
  tuple_title={'ru':'Прогулка','kz':'Серуен','en':'Walk'}[lang]
  assert tuple_code.count(f'(String::from("{tuple_title}"), 30)')==1
  check_program(tuple_code.replace(f'(String::from("{tuple_title}"), 30)',
                                   f'(30, String::from("{tuple_title}"))'),
                '',Path(name),True,'E0308');count+=1
  status_done={'ru':'Завершено','kz':'Аяқталған','en':'Finished'}[lang]
  check('29-enums','Status::Doing(25)','Status::Done',status_done+'\n')
  active_label={'ru':('Работа','После завершения'),
                'kz':('Жұмыс','Аяқталғаннан кейін'),
                'en':('Active','After finishing')}[lang]
  check('30-patterns','Status::Doing(25)','Status::Planned',
        active_label[0]+': 0 min\n'+active_label[1]+': 0 min\n'
        if lang=='en' else active_label[0]+': 0 мин\n'+active_label[1]+': 0 мин\n')
  def check_live(stem, input_text, arguments, expected):
   global count
   source=Path(name)/'boundary.rs'
   source.write_text((R/'answers'/f'{stem}{suffix}-answer.rs').read_text(encoding='utf-8'),encoding='utf-8')
   binary=Path(name)/('boundary.exe' if os.name=='nt' else 'boundary')
   built=subprocess.run(['rustc','--edition=2024',str(source),'-o',str(binary)],
                        capture_output=True,text=True,encoding='utf-8',errors='replace',timeout=90)
   assert built.returncode==0,(stem,lang,built.stderr)
   result=subprocess.run([str(binary),*arguments],input=input_text,capture_output=True,
                         text=True,encoding='utf-8',errors='replace',timeout=90)
   assert result.returncode==0 and result.stdout==expected,(stem,lang,result.stdout,result.stderr,expected)
   count+=1
  title={'ru':'Мой план','kz':'Менің жоспарым','en':'My plan'}[lang]
  input_labels={'ru':('Добавить: ','Нужно название'),
                'kz':('Қосу: ','Атау қажет'),
                'en':('Add: ','A title is needed')}[lang]
  check_live('34-input',title+'\n',[],input_labels[0]+title+'\n')
  check_live('34-input','\n',[],input_labels[1]+'\n')
  argument_labels={'ru':('Запланировано: ','Нужно: add "Название"','Команда не найдена'),
                   'kz':('Жоспарланды: ','Қажет: add "Атау"','Пәрмен табылмады'),
                   'en':('Planned: ','Need: add "Title"','Command not found')}[lang]
  check_live('35-arguments','',['add',title],argument_labels[0]+title+'\n')
  check_live('35-arguments','',['add',*title.split(' ')],argument_labels[1]+'\n')
  check_live('35-arguments','',['unknown',title],argument_labels[2]+'\n')
 print(f'PASS: {count} additional exercise boundary cases across three locales')
