#!/usr/bin/env python3
"""Render block eight with original, accessible and localized teaching graphics."""
try:
    from .chemistry_acids_data import LESSONS
    from .chemistry_solutions_data import LESSONS as PREVIOUS
    from .build_chemistry_block import HEAD, RETURN, LANGS, LESSON_DIR, MAP_DIR
    from .build_chemistry_elements import start_svg, text
    from .build_chemistry_particles import slug
    from .build_chemistry_bonds import LABELS
    from .build_chemistry_solutions import lines
except ImportError:
    from chemistry_acids_data import LESSONS
    from chemistry_solutions_data import LESSONS as PREVIOUS
    from build_chemistry_block import HEAD, RETURN, LANGS, LESSON_DIR, MAP_DIR
    from build_chemistry_elements import start_svg, text
    from build_chemistry_particles import slug
    from build_chemistry_bonds import LABELS
    from build_chemistry_solutions import lines

MODEL={'ru':'Учебная схема · данные заданы · не инструкция для опыта',
       'kz':'Оқу сызбасы · деректер берілген · тәжірибе нұсқауы емес',
       'en':'Teaching schematic · supplied data · not an experiment procedure'}
SAFETY={
'ru':'Только бумага, калькулятор и готовые данные. Не пробуйте, не нюхайте, не смешивайте и не нагревайте химикаты; не открывайте бытовые средства и неизвестные ёмкости. Нейтральный pH не доказывает безопасность. Не применяйте учебную нейтрализацию к коже, глазам, разливам или неизвестным смесям.',
'kz':'Тек қағаз, калькулятор және дайын деректер. Химикатты татпаңыз, иіскемеңіз, араластырмаңыз және қыздырмаңыз; тұрмыстық құралды не белгісіз ыдысты ашпаңыз. Бейтарап pH қауіпсіздікті дәлелдемейді. Оқу бейтараптануын теріге, көзге, төгіндіге немесе белгісіз қоспаға қолданбаңыз.',
'en':'Use paper, a calculator and supplied data only. Do not taste, smell, mix or heat chemicals; do not open household products or unknown containers. Neutral pH does not establish safety. Do not apply teaching neutralisation to skin, eyes, spills or unknown mixtures.'}
SOURCES=[
('OpenStax: acid–base models','https://openstax.org/books/chemistry-2e/pages/14-1-bronsted-lowry-acids-and-bases'),
('OpenStax: pH and pOH','https://openstax.org/books/chemistry-2e/pages/14-2-ph-and-poh'),
('OpenStax: acid strength','https://openstax.org/books/chemistry-2e/pages/14-3-relative-strengths-of-acids-and-bases'),
('OpenStax: acid–base calculations and indicators','https://openstax.org/books/chemistry-2e/pages/14-7-acid-base-titrations'),
('OpenStax: reaction equations','https://openstax.org/books/chemistry-2e/pages/4-2-classifying-chemical-reactions'),
('CCOHS: product labels and incompatibility','https://www.ccohs.ca/oshanswers/chemicals/cleaning-products-working-safely.html')]
SOURCE_KEYS=[(0,),(3,),(1,),(0,2),(0,4),(1,3),(1,2),(5,),(0,1,3,5)]
STEPS={
'ru':[
('КИСЛОТА','Отдаёт H⁺','ОСНОВАНИЕ','Принимает H⁺','ВОДА','H₃O⁺ и OH⁻','ПРОВЕРЬ','Атомы и заряд'),
('ИНДИКАТОР','Назови вещество','ШКАЛА','Интервал цвета','НАБЛЮДЕНИЕ','Данные карточки','ВЫВОД','Среда ≠ вещество'),
('МОДЕЛЬ','Разбавленный раствор','pH','−lg(c / 1 mol/L)','ШАГ','В 10 раз по c','НЕЙТРАЛЬНО','Равенство ионов'),
('СИЛА','Равновесие','c₀','Моли на литр','α','Доля ионизации','ОПАСНОСТЬ','Отдельный вопрос'),
('ИСХОДНО','H₃O⁺ и OH⁻','РЕАКЦИЯ','2H₂O','НАБЛЮДАТЕЛИ','Na⁺ и Cl⁻','ВЫВОД','Раствор соли'),
('ОБЪЁМ','Переведи в литры','КОЛИЧЕСТВО','n = cV','БАЛАНС','Найди избыток','pH','Проверь модель'),
('СОХРАНИ','n кислоты','ОБЪЁМ','Конечный раствор','ОЦЕНИ','Вклад воды','ВЫЧИСЛИ','pH с оговорками'),
('ПРОДУКТ','Что известно?','ДОКУМЕНТ','Этикетка и паспорт','ЗАПРЕТ','Не смешивать','РЕШЕНИЕ','Ответственный взрослый'),
('ДАННЫЕ','V, c, T','БАЛАНС','Кислота и щёлочь','СОСТАВ','Соль остаётся','ПРЕДЕЛ','pH ≠ безопасность')],
'kz':[
('ҚЫШҚЫЛ','H⁺ береді','НЕГІЗ','H⁺ қабылдайды','СУ','H₃O⁺ және OH⁻','ТЕКСЕР','Атом мен заряд'),
('ИНДИКАТОР','Затты ата','ШКАЛА','Түс аралығы','БАҚЫЛАУ','Карточка дерегі','ҚОРЫТЫНДЫ','Орта ≠ зат'),
('МОДЕЛЬ','Сұйылтылған ерітінді','pH','−lg(c / 1 mol/L)','ҚАДАМ','c он есе өзгереді','БЕЙТАРАП','Иондар теңдігі'),
('КҮШ','Тепе-теңдік','c₀','Литрдегі моль','α','Ионданған үлес','ҚАУІП','Бөлек сұрақ'),
('БАСТАПҚЫ','H₃O⁺ және OH⁻','РЕАКЦИЯ','2H₂O','БАҚЫЛАУШЫ','Na⁺ және Cl⁻','ҚОРЫТЫНДЫ','Тұз ерітіндісі'),
('КӨЛЕМ','Литрге ауыстыр','МӨЛШЕР','n = cV','БАЛАНС','Артықтықты тап','pH','Модельді тексер'),
('САҚТА','Қышқылдың n-і','КӨЛЕМ','Соңғы ерітінді','БАҒАЛА','Судың үлесі','ЕСЕПТЕ','pH және шектеу'),
('ӨНІМ','Не белгілі?','ҚҰЖАТ','Заттаңба, паспорт','ТЫЙЫМ','Араластырма','ШЕШІМ','Жауапты ересек'),
('ДЕРЕК','V, c, T','БАЛАНС','Қышқыл мен сілті','ҚҰРАМ','Тұз қалады','ШЕК','pH ≠ қауіпсіздік')],
'en':[
('ACID','Donates H⁺','BASE','Accepts H⁺','WATER','H₃O⁺ and OH⁻','CHECK','Atoms and charge'),
('INDICATOR','Name the substance','SCALE','Colour interval','OBSERVATION','Supplied card','CONCLUDE','Class ≠ identity'),
('MODEL','Dilute solution','pH','−lg(c / 1 mol/L)','STEP','Tenfold change in c','NEUTRAL','Equal ion levels'),
('STRENGTH','Equilibrium','c₀','Moles per litre','α','Ionised fraction','HAZARD','Separate question'),
('START','H₃O⁺ and OH⁻','REACTION','2H₂O','SPECTATORS','Na⁺ and Cl⁻','CONCLUDE','Salt solution'),
('VOLUME','Convert to litres','AMOUNT','n = cV','BALANCE','Find the excess','pH','Check the model'),
('CONSERVE','Acid amount n','VOLUME','Final solution','ASSESS','Water contribution','CALCULATE','pH with limits'),
('PRODUCT','What is known?','DOCUMENT','Label and data sheet','PROHIBITION','Do not mix','DECISION','Responsible adult'),
('DATA','V, c, T','BALANCE','Acid and alkali','CONTENTS','Salt remains','LIMIT','pH ≠ safety')]
}
ROWS={
59:[('HCl + H₂O','→','H₃O⁺ + Cl⁻'),('NH₃ + H₂O','⇌','NH₄⁺ + OH⁻'),('HNO₃ + H₂O','→','H₃O⁺ + NO₃⁻')],
60:[('A: litmus','red','acidic'),('B: litmus','purple','neutral'),('C: litmus','blue','basic'),('D: phenolphthalein','colourless','not enough alone')],
61:[('10⁻³ mol/L','−lg(10⁻³)','pH ≈ 3'),('10⁻⁵ mol/L','−lg(10⁻⁵)','pH ≈ 5'),('10⁻⁷ mol/L; 25 °C','c(H₃O⁺) = c(OH⁻)','pH ≈ 7')],
62:[('HA','9 HA + A⁻ + H₃O⁺','α = 1/10'),('c₀ = 0.010; α = 0.10','c₀ × α','0.0010 mol/L'),('c₀ = 0.00010; α ≈ 1','c₀ × α','≈ 0.00010 mol/L')],
63:[('H₃O⁺ + OH⁻','→','2H₂O'),('Na⁺, Cl⁻','before and after','spectators'),('K⁺, NO₃⁻','before and after','spectators')],
64:[('25 mL; 0.20 mol/L','0.025 × 0.20','0.0050 mol HCl'),('40 mL; 0.10 mol/L','0.040 × 0.10','0.0040 mol NaOH'),('0.0050 − 0.0040','acid excess','0.0010 mol')],
65:[('0.010 L; 0.010 mol/L','n = 0.00010 mol','pH ≈ 2'),('V₂ = 0.100 L','c₂ = 0.0010 mol/L','pH ≈ 3'),('c(HCl) = 10⁻⁸ mol/L','water matters','not pH 8')],
66:[('A: pH 2','do not mix','follow label'),('B: pH 11','use separately','follow label'),('C: unknown','no label','responsible adult')],
67:[('20 mL + 20 mL; 0.10 mol/L','n(NaCl) = 0.0020 mol','m = 0.117 g'),('V = 0.040 L; 25 °C','c(NaCl) = 0.050 mol/L','pH ≈ 7 in model'),('20 mL HCl + 18 mL NaOH','acid excess = 0.0002 mol','pH ≈ 2.28')]
}
WORDS={
'ru':{'A: litmus':'A: лакмус','B: litmus':'B: лакмус','C: litmus':'C: лакмус','D: phenolphthalein':'D: фенолфталеин','red':'красный','purple':'фиолетовый','blue':'синий','acidic':'кислая','neutral':'нейтральная','basic':'щелочная','colourless':'бесцветный','not enough alone':'одного результата мало','before and after':'до и после','spectators':'ионы-наблюдатели','acid excess':'избыток кислоты','water matters':'нужен вклад воды','not pH 8':'не pH 8','do not mix':'не смешивать','follow label':'следовать этикетке','use separately':'использовать отдельно','C: unknown':'C: неизвестное','no label':'нет этикетки','responsible adult':'ответственный взрослый','pH ≈ 7 in model':'pH ≈ 7 в модели','acid excess = 0.0002 mol':'избыток кислоты = 0.0002 mol'},
'kz':{'A: litmus':'A: лакмус','B: litmus':'B: лакмус','C: litmus':'C: лакмус','D: phenolphthalein':'D: фенолфталеин','red':'қызыл','purple':'күлгін','blue':'көк','acidic':'қышқыл','neutral':'бейтарап','basic':'сілтілі','colourless':'түссіз','not enough alone':'бір нәтиже жеткіліксіз','before and after':'дейін және кейін','spectators':'бақылаушы иондар','acid excess':'қышқыл артықтығы','water matters':'судың үлесі маңызды','not pH 8':'pH 8 емес','do not mix':'араластырмау','follow label':'заттаңбаны орындау','use separately':'бөлек қолдану','C: unknown':'C: белгісіз','no label':'заттаңба жоқ','responsible adult':'жауапты ересек','pH ≈ 7 in model':'модельде pH ≈ 7','acid excess = 0.0002 mol':'қышқыл артықтығы = 0.0002 mol'},'en':{}}
QUIZ={
'ru':[
('В NH₃ + H₂O ⇌ NH₄⁺ + OH⁻ кто отдаёт H⁺? Назовите сопряжённую пару для NH₃.','H₂O отдаёт H⁺; NH₄⁺/NH₃. Суммарный заряд справа ноль.'),
('Бесцветный фенолфталеин доказывает pH 7 и отсутствие кислоты?','Нет. Результат совместим с широким диапазоном; нужны шкала и дополнительные данные.'),
('При pH 4 и 6 каково отношение c(H₃O⁺) в принятом приближении? Нейтральность всегда pH 7?','100:1; нет, нейтральность задаётся равенством H₃O⁺ и OH⁻, нейтральный pH зависит от температуры.'),
('c₀=0,030 моль/л и заданная α=0,020: вклад HA в H₃O⁺? Сохраняется ли α при любом разбавлении?','0,00060 моль/л при указанных допущениях; нет, требуется новое равновесие.'),
('В реакции HCl и NaOH какие ионы-наблюдатели? Исчезли ли они из раствора?','Na⁺ и Cl⁻; нет. Сокращённая запись H₃O⁺ + OH⁻ → 2H₂O не удаляет их из состава.'),
('10 мл 0,10 моль/л HCl и 10 мл 0,20 моль/л NaOH: что в избытке и сколько NaCl образуется?','Избыток NaOH 0,0010 моль; образуется 0,0010 моль NaCl. Равные объёмы не дают равных молей при разных c.'),
('10 мл 0,0010 моль/л HCl доведены в модели до 100 мл: pH при 25 °C?','При малом вкладе воды и концентрационном приближении c₂=0,00010 моль/л, pH≈4; n=0,000010 моль сохраняется.'),
('Карточки двух средств имеют pH 2 и 11. Можно ли рассчитывать пропорцию их бытового смешивания?','Нет. Состав и совместимость неизвестны; не смешивать, соблюдать этикетки и обращаться к ответственному взрослому.'),
('В исходном итоговом случае 20+20 мл остаётся ли NaCl? Назовите массу и предел вывода.','Да, 0,0020 моль и 0,117 г по заданной M. pH≈7 относится к модели при 25 °C, а не доказывает чистоту или безопасность.')],
'kz':[
('NH₃ + H₂O ⇌ NH₄⁺ + OH⁻ кезінде H⁺-ты кім береді? NH₃ үшін түйіндес жұпты атаңыз.','H₂O H⁺ береді; NH₄⁺/NH₃. Оң жақтың жалпы заряды нөл.'),
('Түссіз фенолфталеин pH 7 және қышқыл жоқтығын дәлелдей ме?','Жоқ. Нәтиже кең диапазонға сәйкес; шкала мен қосымша дерек қажет.'),
('Қабылданған жуықтауда pH 4 пен 6 үшін c(H₃O⁺) қатынасы қандай? Бейтараптық әрқашан pH 7 ме?','100:1; жоқ, бейтараптық H₃O⁺ пен OH⁻ теңдігімен анықталады, бейтарап pH температураға тәуелді.'),
('c₀=0,030 моль/л және берілген α=0,020: HA-ның H₃O⁺-қа үлесі қандай? α кез келген сұйылтуда сақтала ма?','Берілген жорамалдарда 0,00060 моль/л; жоқ, жаңа тепе-теңдік қажет.'),
('HCl мен NaOH реакциясында бақылаушы иондар қандай? Олар ерітіндіден жоғала ма?','Na⁺ пен Cl⁻; жоқ. H₃O⁺ + OH⁻ → 2H₂O қысқартылған жазбасы оларды құрамнан жоймайды.'),
('10 мл 0,10 моль/л HCl және 10 мл 0,20 моль/л NaOH: не артық, қанша NaCl түзіледі?','NaOH артықтығы 0,0010 моль; 0,0010 моль NaCl түзіледі. c әртүрлі болса, тең көлем тең моль бермейді.'),
('Модельде 10 мл 0,0010 моль/л HCl соңғы 100 мл береді: 25 °C-та pH қандай?','Су үлесі аз және концентрациялық жуықтау жарамды болса c₂=0,00010 моль/л, pH≈4; n=0,000010 моль сақталады.'),
('Екі құрал карточкасында pH 2 және 11. Оларды тұрмыста араластыру пропорциясын есептеуге бола ма?','Жоқ. Құрам мен үйлесімділік белгісіз; араластырмаңыз, заттаңбаны орындап, жауапты ересекке хабарлаңыз.'),
('Бастапқы 20+20 мл қорытынды мысалда NaCl қала ма? Массасы мен қорытынды шегін атаңыз.','Иә, берілген M бойынша 0,0020 моль және 0,117 г. pH≈7 — 25 °C моделі; ол тазалық не қауіпсіздік дәлелі емес.')],
'en':[
('In NH₃ + H₂O ⇌ NH₄⁺ + OH⁻, who donates H⁺? Name NH₃’s conjugate pair.','Water donates; NH₄⁺/NH₃. Total product charge is zero.'),
('Does colourless phenolphthalein establish pH seven and absence of acid?','No. The result fits a broad range; a scale and more evidence are needed.'),
('For pH four and six, what is the hydronium-concentration ratio in this approximation? Is neutrality always pH seven?','100:1; no. Neutrality means equal hydronium and hydroxide concentrations; neutral pH depends on temperature.'),
('For c₀=0.030 mol/L and supplied α=0.020, find HA’s hydronium contribution. Does α survive every dilution unchanged?','0.00060 mol/L under the stated assumptions; no, a new equilibrium is needed.'),
('Identify the spectator ions for HCl and NaOH. Do they disappear from solution?','Na⁺ and Cl⁻; no. H₃O⁺ + OH⁻ → 2H₂O does not remove them from the contents.'),
('10 mL of 0.10 mol/L HCl and 10 mL of 0.20 mol/L NaOH: find excess and NaCl amount.','0.0010 mol NaOH excess and 0.0010 mol NaCl formed. Equal volumes do not give equal moles at unequal concentrations.'),
('A model takes 10 mL of 0.0010 mol/L HCl to a final 100 mL. Find pH at 25 °C.','With negligible water contribution and an applicable concentration approximation, c₂=0.00010 mol/L and pH≈4; n=0.000010 mol is conserved.'),
('Two product cards show pH two and eleven. May you calculate household mixing proportions?','No. Composition and compatibility are unknown: do not mix; follow labels and contact a responsible adult.'),
('Does NaCl remain in the original 20+20 mL case? Give its mass and the conclusion’s limit.','Yes: 0.0020 mol and 0.117 g using the supplied molar mass. Approximately pH seven is a model result at 25 °C, not proof of purity or safety.')]
}
ALT={
'ru':[
'Примеры переноса протона: HCl с водой даёт H₃O⁺ и Cl⁻; NH₃ с водой образует NH₄⁺ и OH⁻. Атомы и заряд сохраняются.',
'Карточки лакмуса: красный — кислая, фиолетовый — нейтральная, синий — щелочная среда. Фенолфталеин без цвета сам по себе не определяет pH.',
'Числовая шкала pH: 3 и 5 разделены двумя единицами, отношение концентраций H₃O⁺ в разбавленной модели 100. Нейтральный pH около 7 при 25 °C.',
'Из десяти кислотных единиц девять HA и одна A⁻; рядом один H₃O⁺. Заданный ионизированный вклад 1 из 10, общий заряд ноль; вода опущена.',
'H₃O⁺ и OH⁻ дают две H₂O. Ионы Na⁺ и Cl⁻ присутствуют до и после; их сокращение в уравнении не удаляет соль из раствора.',
'Количество HCl 5,0 ммоль, NaOH 4,0 ммоль; отношение 1:1, кислотный избыток 1,0 ммоль. Длины полос пропорциональны молям, не объёмам.',
'Сохраняется 0,00010 моль HCl: конечный объём меняется от 0,010 до 0,100 л, приближённый pH от 2 до 3 при 25 °C.',
'Вымышленные этикетки A и B запрещают смешивание; неизвестная карточка C требует прекратить исследование и обратиться к ответственному взрослому.',
'При равных количествах образуется 0,0020 моль NaCl массой 0,117 г. При 18 мл щёлочи остаётся кислотный избыток; pH≈2,28 в модели.'],
'kz':[
'HCl мен су H₃O⁺ пен Cl⁻ береді; NH₃ пен су NH₄⁺ пен OH⁻ түзеді. Протон беріледі, атомдар мен заряд сақталады.',
'Лакмус карточкалары: қызыл — қышқыл, күлгін — бейтарап, көк — сілтілі орта. Түссіз фенолфталеин жеке өзі pH анықтамайды.',
'pH сандық шкаласында 3 пен 5 арасы екі бірлік; сұйылтылған модельдегі H₃O⁺ концентрацияларының қатынасы 100. 25 °C-та бейтарап pH шамамен 7.',
'Он қышқыл бірлігінің тоғызы HA, біреуі A⁻; жанында бір H₃O⁺. Берілген ионданған үлес 1/10, жалпы заряд нөл; су көрсетілмеген.',
'H₃O⁺ пен OH⁻ екі H₂O береді. Na⁺ пен Cl⁻ дейін де, кейін де бар; теңдеудегі қысқарту тұзды ерітіндіден жоймайды.',
'HCl мөлшері 5,0 ммоль, NaOH 4,0 ммоль; қатынас 1:1, қышқыл артықтығы 1,0 ммоль. Жолақ ұзындығы көлемге емес, мольге пропорционал.',
'HCl-дың 0,00010 молі сақталады: соңғы көлем 0,010-нан 0,100 литрге, жуық pH 25 °C-та 2-ден 3-ке өзгереді.',
'Ойдан алынған A және B заттаңбасы араластыруға тыйым салады; белгісіз C зерттеуді тоқтатып, жауапты ересекке хабарлауды талап етеді.',
'Тең мөлшерде массасы 0,117 г болатын 0,0020 моль NaCl түзіледі. Сілті 18 мл болса, қышқыл артық қалады; модельде pH≈2,28.'],
'en':[
'Proton-transfer examples: HCl and water give H₃O⁺ and Cl⁻; NH₃ and water give NH₄⁺ and OH⁻. Atoms and charge are conserved.',
'Litmus cards: red acidic, purple neutral, blue basic. Colourless phenolphthalein alone does not determine pH.',
'On the numerical pH scale, three and five differ by two units; the dilute-model hydronium-concentration ratio is one hundred. Neutral pH is near seven at 25 °C.',
'Ten acid units comprise nine HA and one A⁻, with one accompanying H₃O⁺. Supplied ionised fraction is one tenth and net charge zero; water is omitted.',
'H₃O⁺ and OH⁻ form two H₂O. Na⁺ and Cl⁻ are present before and after: cancelling them in the equation does not remove dissolved salt.',
'HCl amount is 5.0 mmol and NaOH 4.0 mmol: the 1:1 balance leaves 1.0 mmol acid excess. Bar lengths represent mole amounts, not volumes.',
'0.00010 mol HCl is conserved: final volume changes from 0.010 to 0.100 L and approximate pH from two to three at 25 °C.',
'Fictional A and B labels prohibit mixing; unknown C requires stopping investigation and contacting a responsible adult.',
'Equal amounts form 0.0020 mol NaCl with mass 0.117 g. With 18 mL alkali, acid excess remains and model pH is approximately 2.28.']}

def word(value,lang):return WORDS[lang].get(value,value)

def map_svg(number,lang):
    p=start_svg(LESSONS[number-59][lang][0],MODEL[lang],height=600)
    lines(p,40,55,LESSONS[number-59][lang][0],58,29)
    p.append(text(40,155,MODEL[lang],20,'#bedde4'))
    for i in range(4):
        x=40+i*310
        p.append(f'<rect x="{x}" y="210" width="275" height="300" rx="18" fill="#e6f4ef"/>')
        lines(p,x+18,265,STEPS[lang][number-59][i*2],18,19,'#143242')
        lines(p,x+18,345,STEPS[lang][number-59][i*2+1],17,23,'#143242')
        if i<3:p.append(text(x+282,385,'→',27))
    return '\n'.join(p+['</svg>'])+'\n'

def visual(number,lang):
    title=LESSONS[number-59][lang][0];p=start_svg(title,ALT[lang][number-59],height=820)
    lines(p,40,55,f'{number} · '+title,60,29)
    p.append(text(40,155,MODEL[lang],20,'#bedde4'))
    if number==61:
        # Equal steps in pH represent multiplicative, not additive, c changes.
        for value in range(15):
            x=70+value*80
            p.append(f'<rect data-ph="{value}" x="{x}" y="290" width="74" height="110" rx="10" fill="{"#ac685f" if value<7 else "#77aaa2" if value==7 else "#5e85b3"}"/>')
            p.append(text(x+37,355,value,28,anchor='middle'))
        for value,label in [(3,'10⁻³ mol/L'),(5,'10⁻⁵ mol/L')]:
            x=107+value*80
            p.append(f'<path d="M{x} 410 V465" stroke="#efbc7e" stroke-width="4"/>')
            p.append(text(x,510,label,24,anchor='middle'))
        p.append(text(830,525,'10⁻³ / 10⁻⁵ = 100',32,anchor='middle'))
        p.append(text(640,610,'25 °C · pH ≈ 7 ↔ c(H₃O⁺) = c(OH⁻)',28,anchor='middle'))
    elif number==62:
        for i in range(10):
            x=120+(i%5)*200;y=270+(i//5)*150
            species='A⁻' if i==9 else 'HA'
            p.append(f'<circle data-acid-unit="{species}" cx="{x}" cy="{y}" r="48" fill="{"#e9bb80" if i==9 else "#a1d9dd"}"/>')
            p.append(text(x,y+9,species,27,'#102b3b','middle'))
        p.append('<circle data-hydronium="1" cx="1140" cy="420" r="58" fill="#e9bb80"/>')
        p.append(text(1140,429,'H₃O⁺',27,'#102b3b','middle'))
        p.append(text(640,590,'α = 1/10 = 0.10     q = 0',38,anchor='middle'))
    elif number==64:
        for i,(label,value) in enumerate([('HCl',5),('NaOH',4),('HCl − NaOH',1)]):
            y=225+i*130
            p.append(text(50,y+47,label,29))
            p.append(f'<rect data-amount-mmol="{value}" x="350" y="{y}" width="{value*140}" height="75" rx="10" fill="{"#e9bb80" if i==2 else "#75b6bf"}"/>')
            p.append(text(370+value*140,y+48,f'{value}.0 mmol',27))
        p.append(text(640,650,'1 : 1     5.0 − 4.0 = 1.0 mmol',34,anchor='middle'))
    elif number in (59,63):
        equations=ROWS[number] if number==59 else [('H₃O⁺ + OH⁻','→','2H₂O'),('Na⁺ + Cl⁻','→','Na⁺ + Cl⁻')]
        for i,parts in enumerate(equations):
            y=230+i*150
            p.append(f'<rect x="40" y="{y}" width="1200" height="120" rx="18" fill="#173e4d"/>')
            p.append(text(640,y+70,'   '.join(parts),37,anchor='middle'))
    elif number==60:
        for i,(label,colour,meaning) in enumerate(ROWS[60]):
            x=40+i*310
            p.append(f'<rect x="{x}" y="220" width="275" height="370" rx="18" fill="#173e4d"/>')
            lines(p,x+18,265,word(label,lang),17,22)
            p.append(f'<circle cx="{x+137}" cy="365" r="42" fill="{["#e78b89","#a297ca","#85b4e8","#f1f4ef"][i]}"/>')
            lines(p,x+18,465,word(colour,lang),17,24)
            lines(p,x+18,525,word(meaning,lang),17,23)
    else:
        for i,row in enumerate(ROWS[number]):
            y=205+i*145
            p.append(f'<rect x="40" y="{y}" width="1200" height="130" rx="18" fill="#173e4d"/>')
            lines(p,65,y+42,word(row[0],lang)+'  →  '+word(row[1],lang),67,25)
            lines(p,65,y+98,word(row[2],lang),65,29)
    lines(p,40,725,ALT[lang][number-59],92,21)
    return '\n'.join(p+['</svg>'])+'\n'

def table(number,lang):
    head={'ru':('Объект / шаг','Дано / расчёт','Вывод'), 'kz':('Нысан / қадам','Берілгені / есеп','Қорытынды'), 'en':('Object / step','Given / calculation','Conclusion')}[lang]
    return '| '+' | '.join(head)+' |\n|---|---|---|\n'+'\n'.join('| '+' | '.join(word(s,lang) for s in row)+' |' for row in ROWS[number])

def render(number,lang):
    row=LESSONS[number-59];title,summary,*sections=row[lang];labels=LABELS[lang]
    body=[f'# {title}',f'_Лид (summary):_ **{summary}**',f'## {HEAD[lang][0]}',f'![{title}](/static/course/chemistry/map-{number:02d}-{lang}.svg)',f'## {labels[0]}',f'![{ALT[lang][number-59]}](/static/course/chemistry/visuals/{number:02d}-acids-{lang}.svg)',SAFETY[lang]]
    for i,(heading,paragraph) in enumerate(zip(HEAD[lang][1:9],sections,strict=True)):
        body.extend([f'## {heading}',paragraph])
        if i==3:body.append(table(number,lang))
    q,a=QUIZ[lang][number-59]
    body.extend([f'## {labels[1]}',q,f'## {labels[2]}',a,f'## {HEAD[lang][10]}',RETURN[lang],f'## {labels[4]}','\n'.join(f'- [{SOURCES[k][0]}]({SOURCES[k][1]})' for k in SOURCE_KEYS[number-59])])
    prev=PREVIOUS[-1] if number==59 else LESSONS[number-60]
    body.append(f'[{prev[lang][0]}](/read/{slug(number-1,prev["stem"])}?lang={lang})')
    if number<67:
        nxt=LESSONS[number-58];body.append(f'[{nxt[lang][0]}](/read/{slug(number+1,nxt["stem"])}?lang={lang})')
    return '\n\n'.join(body)+'\n'

def bridge_body(lang):
    try:from .build_chemistry_solutions import render as previous_render
    except ImportError:from build_chemistry_solutions import render as previous_render
    row=LESSONS[0]
    return previous_render(58,lang).split('\n\n',2)[2].rstrip()+f'\n\n[{row[lang][0]}](/read/{slug(59,row["stem"])}?lang={lang})\n'

def main():
    (MAP_DIR/'visuals').mkdir(parents=True,exist_ok=True)
    for number,row in enumerate(LESSONS,59):
        for lang in LANGS:
            suffix='' if lang=='ru' else '-'+lang
            (LESSON_DIR/f'{row["stem"]}{suffix}.md').write_text(render(number,lang),encoding='utf8')
            (MAP_DIR/f'map-{number:02d}-{lang}.svg').write_text(map_svg(number,lang),encoding='utf8')
            (MAP_DIR/f'visuals/{number:02d}-acids-{lang}.svg').write_text(visual(number,lang),encoding='utf8')

if __name__=='__main__':main()
