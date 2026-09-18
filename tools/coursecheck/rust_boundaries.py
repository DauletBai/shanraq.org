"""Verify boundary variations promised in Rust lessons 11–15 in each locale."""
import sys,tempfile,re
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
 print(f'PASS: {count} additional exercise boundary cases across three locales')
