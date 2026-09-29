"""Reviewed Kazakh and English copy for mathematics lessons 52–59."""


def E(title, summary, where, alt, image, meaning, signal, worked, mistake, task, next_):
    return {
        "title": title, "summary": summary, "where": where, "alt": alt,
        "image": image, "meaning": meaning, "signal": signal,
        "worked": worked, "mistake": mistake, "task": task, "next": next_,
    }


DATA_EXTRA = {
    "52-combinatorics": {
        "kz": E(
            "Комбинаторика: нұсқаларды қолмен қайта санамай есептеу",
            "Таңдау қадамдарын көреміз, реттелген таңдауды жиыннан ажыратамыз, қосу мен көбейту ережелерін, орналастырулар мен терулерді шағын тізім арқылы тексереміз.",
            "Көбейту, бөлшек және айнымалы енді нақты бір нәтижені емес, мүмкін шешімдердің санын табуға көмектеседі.",
            "Бес нысаннан 20 реттелген жұп және 10 реттелмеген жұп шығады",
            "A, B, C, D, E деген бес кітап бар делік. Бірінші және екінші орын үшін AB мен BA — екі бөлек нәтиже. Ал бір сөреге екі кітапты таңдау кезінде екеуі бір жұпты білдіреді.",
            "Қосу ережесі өзара үйлеспейтін жағдайлардың сандарын қосады. Көбейту ережесі тізбекті қадамдардағы мүмкіндіктерді көбейтеді. `P(n,k)=n!/(n−k)!` ретті ескереді; `C(n,k)=n!/(k!(n−k)!)` ретті ескермейді. Формуладан бұрын қайталауға бола ма және рет маңызды ма деген екі сұраққа жауап береміз.",
            "таңдауды сипаттау → қадамдарға бөлу → рет? → қайталау? → қосу не көбейту → формула → тізіммен тексеру",
            "Бес кітаптан жеңімпаз бен екінші орынды таңдау саны `P(5,2)=5·4=20`. Орынсыз екі кітап таңдағанда әр жұп екі рет саналған, сондықтан `C(5,2)=20/2=10`. AB, AC, AD, AE, BC, BD, BE, CD, CE, DE тізімі он жұпты дәлелдейді.",
            "Бес адамның арасынан төраға мен хатшыны таңдауға `C(5,2)=10` қолдану қате. Рөлдер әртүрлі, сондықтан әр жұп екі тағайындау береді; дұрыс жауап `P(5,2)=20`.",
            "8 оқушыдан: а) сынып жетекшісі мен орынбасарын; ә) екі адамдық комиссияны; б) белгілі бір оқушы міндетті кіретін үш адамдық комиссияны неше тәсілмен таңдауға болатынын табыңыз. Әр жауапты төрт оқушыдан тұратын шағын мысалмен тексеріңіз.",
            "Келесі сабақ тең мүмкін нәтижелердің санын белгісіздік өлшеміне — ықтималдыққа айналдырады."
        ),
        "en": E(
            "Combinatorics: count possibilities without listing them all",
            "We identify stages of choice, distinguish an ordered selection from a set, and use sum, product, permutation, and combination rules with a small enumeration check.",
            "Multiplication, fractions, and variables now help count all possible solutions before one outcome is chosen.",
            "Five objects produce 20 ordered pairs and 10 unordered pairs",
            "Take five books A, B, C, D, and E. For first and second place, AB and BA are different results. For choosing two books for one shelf, both orders describe the same pair.",
            "The sum rule adds mutually exclusive cases. The product rule multiplies possibilities across successive stages. `P(n,k)=n!/(n−k)!` keeps order; `C(n,k)=n!/(k!(n−k)!)` ignores order. Before using a formula, ask whether repetition is allowed and whether order matters.",
            "describe the choice → split into stages → order? → repetition? → sum or product → formula → enumerate to check",
            "Choosing a winner and runner-up from five books gives `P(5,2)=5·4=20`. Choosing two books without ranks counts every pair twice, so `C(5,2)=20/2=10`. The list AB, AC, AD, AE, BC, BD, BE, CD, CE, DE confirms ten pairs.",
            "Using `C(5,2)=10` to choose a chair and secretary is wrong. The roles differ, so each pair supports two assignments; the correct count is `P(5,2)=20`.",
            "From 8 students, count ways to choose: (a) a class representative and deputy; (b) a two-person committee; (c) a three-person committee containing one specified student. Check every answer on a smaller four-student case.",
            "The next lesson turns counts of equally possible outcomes into a measure of uncertainty: probability."
        ),
    },
    "53-probability": {
        "kz": E(
            "Ықтималдық: нөл мен бір арасындағы белгісіздік өлшемі",
            "Нәтижелер кеңістігін құрамыз, оқиғаны жеке нәтижеден ажыратамыз, толықтауыш пен шартты ықтималдықты қолданамыз және затты қайтармай алу ағашын оқимыз.",
            "Комбинаторика нұсқаларды санады. Ықтималдық оларға салмақ беріп, оқиға бақылауға дейін қаншалықты күтілетінін көрсетеді.",
            "Үш қызыл және екі көк шары бар дорбадан екі рет алу ағашы екі қызылдың ықтималдығын 3/10 деп көрсетеді",
            "Мөлдір емес дорбада үш қызыл және екі көк шар бар. Келесі түсті білмейміз, бірақ құрамын білеміз. Шарды қайтармасақ, екінші қадамның ықтималдығы өзгереді.",
            "Оқиға ықтималдығы 0 мен 1 арасында. Тең мүмкін нәтижелерде `P(A)=m/n`. Толықтауыш `P(not A)=1−P(A)`. Тәуелді қадамдарда `P(A∩B)=P(A)P(B|A)`. Тәуелсіздік — A туралы білім B ықтималдығын өзгертпейді; оны жай ғана ыңғай үшін болжауға болмайды.",
            "нәтижелер → оқиға → бұтақ салмақтары → шарт өзгерді ме? → жол бойымен көбейту → қажетті жолдарды қосу → барлық нәтиже қосындысы 1",
            "Қайтармай екі шар алғанда бірінші қызылдың ықтималдығы `3/5`. Одан кейін төрт шардың екеуі қызыл: `P(RR)=3/5·2/4=3/10`. Терулермен тексеру: барлық жұп `C(5,2)=10`, қызыл жұп `C(3,2)=3`.",
            "Бірінші қызылдан кейін екіншісінің ықтималдығын `3/5` күйінде қалдыру қате. Дорбада төрт шардың тек екеуі қызыл қалды, сондықтан шартты ықтималдық `2/4`.",
            "Дорбада 4 ақ, 3 қара және 1 жасыл шар бар. Екі шар қайтармай алынады. Екеуі ақ, түстері әртүрлі және кемінде біреуі жасыл болу ықтималдықтарын табыңыз. Соңғысын толықтауышпен шешіп, жағдайлар қосындысын тексеріңіз.",
            "Келесі сабақ кездейсоқ модельден нақты жиналған деректерді адал сипаттауға өтеді."
        ),
        "en": E(
            "Probability: a measure of uncertainty from zero to one",
            "We build a sample space, distinguish an event from one outcome, and use complements, conditional probability, and a tree for sampling without replacement.",
            "Combinatorics counted possibilities. Probability assigns weights and says how expected an event is before observing the result.",
            "A two-draw tree for a bag with three red and two blue balls gives probability 3/10 for two red balls",
            "An opaque bag contains three red and two blue balls. We do not know the next colour, but we know the composition. Without replacement, that composition changes after the first draw.",
            "An event probability lies from 0 to 1. For equally likely outcomes, `P(A)=m/n`. A complement has `P(not A)=1−P(A)`. Dependent steps use `P(A∩B)=P(A)P(B|A)`. Independence means learning A does not change the probability of B; it cannot be assumed for convenience.",
            "outcomes → event → branch weights → did the condition change? → multiply along paths → add matching paths → total probability 1",
            "Without replacement, the first red has probability `3/5`. Then two of four balls are red, so `P(RR)=3/5·2/4=3/10`. Combinations check it: `C(5,2)=10` total pairs and `C(3,2)=3` red pairs.",
            "Keeping the second-red probability at `3/5` after drawing a red ball is wrong. Only two red balls remain among four, so the conditional probability is `2/4`.",
            "A bag contains 4 white, 3 black, and 1 green ball. Draw two without replacement. Find probabilities that both are white, their colours differ, and at least one is green. Solve the last by a complement and check the complete set of cases.",
            "The next lesson moves from a random mechanism to an honest description of observed data."
        ),
    },
    "54-descriptive-statistics": {
        "kz": E(
            "Сипаттамалық статистика: модельге дейін дерек не айтады",
            "Деректер қатарын үлестірім, орташа, медиана, ауқым және ауытқыма арқылы оқимыз, жиын пішінін бір ыңғайлы санмен жасырмаймыз.",
            "Ықтималдық тәжірибеге дейінгі нәтижелерді сипаттады. Енді бақылауларды өлшенген күйінде анық көрсетеміз.",
            "2, 3, 3, 4, 8 деректерінің нүктелік диаграммасы орташа 4, медиана 3 және ауқым 6 екенін көрсетеді",
            "Бес оқушы есепке 2, 3, 3, 4 және 8 минут жұмсады. 8 қалғандардан алыс. Орташа оған медианадан күштірек жауап береді.",
            "Үлестірім мәндер мен жиіліктерді көрсетеді. Орташа `x̄=Σxᵢ/n`, медиана реттелген қатарды екіге бөледі, мода — ең жиі мән. Ауқым `max−min` шашырауды көрсетеді. Ауытқыманы автоматты өшірмейміз: алдымен өлшеу қатесін және мағынасын тексереміз. Графикте ось, бірлік және таңдау көлемі болуы керек.",
            "сұрақ → бақылау бірлігі → реттеу → үлестірімді көрсету → орталық → шашырау → ерекше нүкте → шектеулі қорытынды",
            "`2,3,3,4,8` қосындысы 20, сондықтан `x̄=20/5=4`. Ортадағы мән медиананы 3 етеді, мода 3, ауқым `8−2=6`. Шын 8 мәні орташа мен медиана айырмасын түсіндіреді.",
            "Үлестірімге қарамай, типтік жалақыны орташамен атау қате. Бірнеше өте үлкен мән орташаға әсер етеді; медиана мен график контекст береді.",
            "`5,7,7,8,9,9,10,15` қатары үшін орташа, медиана, модалар, ауқым және квартильдерді өз ережеңізді атап табыңыз. 15-ті 10-ға ауыстырып, әр өлшемнің сезімталдығын салыстырыңыз.",
            "Таңдауды сипаттау үлкен топ туралы дәлел емес. Келесі сабақ осындай ауысуды сақ орындайды."
        ),
        "en": E(
            "Descriptive statistics: what data say before a model",
            "We read a dataset through its distribution, mean, median, range, and outliers without hiding its shape behind one convenient number.",
            "Probability described possible outcomes before an experiment. We now display what was actually measured.",
            "A dot plot of 2, 3, 3, 4, 8 shows mean 4, median 3, and range 6",
            "Five learners spent 2, 3, 3, 4, and 8 minutes on a problem. Eight lies apart from the others. The mean reacts to it more strongly than the median.",
            "A distribution shows values and frequencies. The mean is `x̄=Σxᵢ/n`, the median splits ordered data, and the mode is most frequent. Range `max−min` describes spread. Do not automatically delete an outlier: first check measurement and context. A graph needs axes, units, and sample size.",
            "question → observational unit → order data → show distribution → centre → spread → unusual points → bounded conclusion",
            "For `2,3,3,4,8`, the sum is 20, so `x̄=20/5=4`. The middle value gives median 3, the mode is 3, and range is `8−2=6`. A genuine 8 explains why the mean exceeds the median.",
            "Calling the mean salary typical without viewing the distribution is wrong. A few very large values can raise the mean; the median and distribution provide context.",
            "For `5,7,7,8,9,9,10,15`, find mean, median, modes, range, and quartiles under a stated convention. Replace 15 with 10, recalculate, and compare each measure's sensitivity.",
            "Describing a sample does not establish properties of a population. The next lesson makes that transition cautiously."
        ),
    },
    "55-inference": {
        "kz": E(
            "Статистикалық қорытынды: таңдау жиынтық туралы не айта алады",
            "Жиынтық пен таңдауды, параметр мен бағалауды, кездейсоқ қатені ығысудан ажыратып, сенімділік аралығын дерек бермейтін уәдесіз оқимыз.",
            "Сипаттамалық статистика жиналған бақылауларды көрсетті. Енді нәтижені кеңірек жиынтыққа көшіруге бола ма және белгісіздік қандай екенін сұраймыз.",
            "100 бақылаудың кездейсоқ таңдауында 0,62 үлес және шамамен 0,52-ден 0,72-ге дейінгі 95 пайыздық аралық алынды",
            "Қала тұрғындарының бәрін сұрамай, таңдау алады. Күндіз тек стационар телефонға қоңырау шалу қамту ығысуын туғызады; үлкен n оны түземейді.",
            "Бас жиынтық — қорытынды жасалатын топ, таңдау — зерттелген бөлік. Параметр белгісіз, таңдау статистикасы оны бағалайды. Стандарттық қате бағалаудың кездейсоқ өзгерісін сипаттайды. Сенімділік аралығы қайталанған тәжірибелерде параметрді берілген жиілікпен қамтитын рәсімге жатады. Бақылаулық байланыс себептілікті дәлелдемейді.",
            "сұрақ → жиынтық → іріктеу тәсілі → бағалау → кездейсоқ қате → аралық → ығысулар → рұқсат етілген қорытынды",
            "`n=100` таңдауда 62 жетістік бар: `p̂=0,62`. `SE≈√(0,62·0,38/100)≈0,049`; шамамен `0,62±1,96·0,049=[0,52;0,72]`. Бұл кездейсоқ қатені көрсетеді, бірақ қисық іріктеуді түземейді.",
            "Құрылған бір аралық параметрді 95% ықтималдықпен қамтиды деу жиіліктік мағынаны бұрмалайды. Параметр тұрақты, ал 95% — қайталанатын рәсімнің қамту жиілігі.",
            "225 кездейсоқ бақылаудың 126-сында белгі бар. `p̂`, жуық стандарттық қате және 95% аралықты табыңыз. Аралық есептемейтін екі ығысуды және бір рұқсат етілген, бір рұқсат етілмеген қорытындыны жазыңыз.",
            "Келесі сабақ есептеуден қандай қорытынды шығып, қайсысы шықпайтынын логика арқылы анықтайды."
        ),
        "en": E(
            "Statistical inference: what a sample can say about a population",
            "We separate population from sample, parameter from estimate, and random error from bias, then read a confidence interval without promises the data cannot support.",
            "Descriptive statistics reported the observations. We now ask whether they support a wider population claim and with what uncertainty.",
            "A random sample of 100 gives proportion 0.62 and an approximate 95 percent interval from 0.52 to 0.72",
            "A city survey usually samples residents. Calling only landlines during working hours creates coverage bias; a large sample does not repair that design.",
            "A population is the group of interest; a sample is the observed part. An unknown population parameter is estimated by a sample statistic. Standard error describes random variation. A confidence interval belongs to a repeated procedure that covers the true parameter at its stated rate. Association in observational data does not prove causation.",
            "question → population → sampling design → estimate → random error → interval → biases → supported conclusion",
            "With 62 successes among `n=100`, `p̂=0.62`. `SE≈√(0.62·0.38/100)≈0.049`, so `0.62±1.96·0.049` is about `[0.52,0.72]`. It reflects random error, but does not fix biased selection.",
            "Saying a completed interval has a 95% probability of containing the parameter misstates the frequentist meaning. The parameter is fixed; 95% describes coverage of the repeated procedure.",
            "In a random sample of 225 observations, 126 have a feature. Find `p̂`, an approximate standard error, and a 95% interval. Name two biases it misses and write one supported and one unsupported conclusion.",
            "The next lesson uses logic to make explicit which conclusions follow from evidence and which do not."
        ),
    },
    "56-logic-proofs": {
        "kz": E(
            "Логика және дәлелдеу: қорытынды шарттан неге шығады",
            "Пайымдарды, терістеуді, логикалық амалдарды, импликация мен кванторларды талдап, тікелей дәлел, қарсы мысал және қайшылықтан дәлел құрамыз.",
            "Статистикалық қорытынды дерек қолдайтын және одан шықпайтын тұжырымды ажыратуды талап етті. Логика осы шекараға дәл тіл береді.",
            "Импликацияның ақиқат кестесі және жұп сан квадратының жұп екенін дәлелдеу тізбегі",
            "«Жаңбыр жауса, жол суланады» деген сөйлем жаңбырдың жалғыз себеп екенін айтпайды. Импликацияны кері бұру — жиі қате.",
            "`p→q` тек p ақиқат, q жалған кезде жалған. Кері `q→p` — бөлек тұжырым. Жалпы тұжырымды жоққа шығаруға бір қарсы мысал жеткілікті. Тікелей дәлел анықтамадан қорытындыға өтеді; қайшылықтан дәлел мақсаттың терістеуін қосып, қайшылық алады. «Барлығы үшін» және «бар» кванторлары орын ауыстырмайды.",
            "шарттар → анықтамалар → рұқсат етілген қадам → аралық қорытынды → мақсат; тексеру: кері ме? қарсы мысал? жасырын болжам?",
            "n жұп болса, `n=2k`. Сонда `n²=(2k)²=4k²=2(2k²)`. `2k²` бүтін, демек n² екі еселенген бүтін сан және жұп.",
            "`n²` жұптығынан n жұп деп дәлелденген тура импликацияға ғана сүйену қате. Бұл — кері тұжырым және оған жеке дәлел, мысалы контрапозиция керек.",
            "Төрт тапсырма орындаңыз: 3-ке еселі сандар қосындысын дәлелдеңіз; «барлық жай сан тақ» пікірін жоққа шығарыңыз; «әр x үшін y>x бар» сөйлемінің терістеуін жазыңыз; `√2` иррационалдығын қайшылықтан дәлелдеңіз.",
            "Келесі сабақ жеке логикалық қадамдарды түйіндер мен байланыстардан тұратын алгоритмге жинайды."
        ),
        "en": E(
            "Logic and proof: why a conclusion follows from assumptions",
            "We analyze statements, negation, connectives, implication, and quantifiers, then build direct proofs, counterexamples, and contradiction arguments.",
            "Statistical inference required us to distinguish what evidence supports from what does not follow. Logic gives that boundary a precise language.",
            "A truth table for implication and a proof chain showing that the square of an even integer is even",
            "“If it rains, the road is wet” does not say rain is the only cause of a wet road. Reversing an implication is a common error.",
            "`p→q` is false only when p is true and q false. Its converse `q→p` is separate. One counterexample refutes a universal claim. A direct proof moves from assumptions by definitions; contradiction adds the negation of the goal and derives an impossibility. “For every” and “there exists” cannot be interchanged.",
            "assumptions → definitions → valid step → intermediate result → goal; check: converse? counterexample? hidden assumption?",
            "If n is even, `n=2k`. Then `n²=(2k)²=4k²=2(2k²)`. Since `2k²` is an integer, n² is twice an integer and is even.",
            "Using the proved forward implication alone to conclude n is even from n² even is wrong. That is the converse and needs its own proof, for example by contraposition.",
            "Complete four tasks: prove sums of multiples of 3 are multiples of 3; refute “all primes are odd”; negate “for every x there exists y>x”; and prove `√2` irrational by contradiction.",
            "The next lesson assembles valid local steps into algorithms over nodes and connections."
        ),
    },
    "57-graphs-algorithms": {
        "kz": E(
            "Графтар және алгоритмдер: түйіндер, байланыстар және ең жақсы жол",
            "Желіні төбелер мен қырлар арқылы модельдеп, жол мен циклді ажыратамыз, Дейкстра алгоритмін қолмен орындаймыз және маршруттың ең қысқа екенін тексереміз.",
            "Логика ауысу ережелерін берді. Граф оларды нысандар арасындағы қырларға айналдырады, алгоритм әрекеттердің қайталанатын ретін белгілейді.",
            "Алты төбелі салмақталған граф A-дан F-ке дейін ұзындығы 13 ең қысқа жолды көрсетеді",
            "Жол картасы, достар желісі, тапсырма тәуелділігі және веб-сілтемелер бір құрылымды бөліседі: нысандар мен байланыстар. Қыр салмағының бірлігі ортақ болуы тиіс.",
            "Граф төбелер мен қырлардан тұрады. Қыр бағытталған не бағытталмаған, салмақталған не салмақсыз болады. Жол — көршілес төбелер тізбегі. Теріс емес салмақта Дейкстра ең жақын бекітілмеген төбені таңдап, көршілер бағасын жаңартады. Теріс қырлар басқа алгоритмді талап етеді.",
            "түйіндер мен қырлар → салмақ бірлігі → бастапқы арақашықтық → ең жақын төбе → көршілерді жаңарту → қайталау → жолды қалпына келтіру",
            "A-дан C=2, B=4. C арқылы B=3 болып жақсарады. B арқылы D=8, D арқылы E=10, E арқылы F=13. Алдыңғы төбелер `A→C→B→D→E→F`, қосынды `2+1+5+2+3=13` жолын береді.",
            "Әр қиылыстағы ең қысқа шығыс қырды таңдап, оны жалпы ең қысқа жол деу қате. Жергілікті арзан қадам қымбат жалғасуға апаруы мүмкін; толық қашықтық салыстырылады.",
            "Сызбада A-дан барлық төбеге дейінгі қысқа қашықтықтарды Дейкстра кестесімен табыңыз. D–E қырын алып тастап, F-ке есепті қайталаңыз. Салмақсыз ен бойынша іздеу керек болатын бір мысал келтіріңіз.",
            "Келесі сабақ дерек, ықтималдық, функция, матрица және графты бір тексерілетін модельдеу цикліне біріктіреді."
        ),
        "en": E(
            "Graphs and algorithms: nodes, connections, and the best route",
            "We model a network with vertices and edges, distinguish paths and cycles, execute Dijkstra's algorithm by hand, and verify that the route is truly shortest.",
            "Logic supplied rules for valid transitions. A graph turns them into edges between objects, while an algorithm gives a reproducible order of actions.",
            "A weighted graph with six vertices shows a shortest A-to-F path of length 13",
            "Road maps, social networks, task dependencies, and web links share one structure: objects and connections. Edge weights may mean time, distance, or cost, but their units must be consistent.",
            "A graph has vertices and edges. Edges may be directed or undirected, weighted or unweighted. A path follows adjacent vertices. With nonnegative weights, Dijkstra repeatedly finalizes the nearest unsettled vertex and relaxes its neighbours. Negative edges require another algorithm.",
            "model nodes and edges → weight unit → initial distances → nearest vertex → update neighbours → repeat → reconstruct path",
            "From A we get C=2 and B=4. Through C, B improves to 3. Through B, D=8; through D, E=10; through E, F=13. Predecessors give `A→C→B→D→E→F`, with `2+1+5+2+3=13`.",
            "Picking the shortest outgoing edge at every intersection and calling the result globally shortest is wrong. A cheap local step can lead to an expensive continuation; the algorithm compares full distances from the start.",
            "Use a Dijkstra table to find shortest distances from A to every vertex. Remove edge D–E and recompute the route to F. Give one example where unweighted breadth-first search is the right tool instead.",
            "The next lesson joins data, probability, functions, matrices, and graphs in one testable modeling cycle."
        ),
    },
    "58-capstone-modeling": {
        "kz": E(
            "Қорытынды математикалық модельдеу: нақты сұрақтан тексерілетін шешімге",
            "Курсты бір циклге жинаймыз: сұрақ қоямыз, шамалар мен болжамдарды таңдаймыз, модель құрамыз, болжам есептейміз, дерекпен салыстырып, қолдану шегін нақтылаймыз.",
            "Барлық алдыңғы блок құрал берді. Енді сан, функция, ықтималдық, матрица немесе графтың сұраққа сай келуін негіздеу керек.",
            "Модельдеу циклі дерек пен болжамнан модельге, болжамға, тексеруге және автобус саны мысалында нақтылауға апарады",
            "Қала артық көлік шығармай, автобус күтуді азайтқысы келеді. Формуладан бұрын маршрут, бақылау уақыты, жолаушы ағыны, сыйымдылық, құн және рұқсат етілген күту анықталады.",
            "Математикалық модель — нысанның әдейі ықшамдалған бейнесі. Онда сұрақ, кіріс, шығыс, бірлік, болжам, теңдеу не алгоритм, параметр дерегі, тексеру және қолдану шегі болуы тиіс. Күрделі модель міндетті емес; ол түсінікті бастапқы тәсілден жақсы әрі тексерілетін болуы керек. Қате мен сезімталдық нәтижемен бірге беріледі.",
            "сұрақ → дерек пен бірлік → болжам → модель → есептеу → болжам → жаңа дерекпен тексеру → нақтылау не бас тарту",
            "`C(b)=b²−12b+52` болса, `C′(b)=2b−12`, сондықтан үздіксіз минимум b=6; `C″(b)=2>0`. Автобус саны бүтін болғандықтан 5,6,7 мәндерін салыстырамыз: 17,16,17. Бірақ сыйымдылық шектеуі мен коэффициенттердің жаңа күндердегі тексеруі қажет.",
            "Формула минимумын тауып, алты автобусты нақты өмірдің оптимумы деп жариялау қате. Сағаттық шың, ақау, маршрут байланысы және таңдау сапасы модельде болмауы мүмкін.",
            "Мини-жоба дайындаңыз: сұрақ, кемінде 20 бақылау, шамалар сөздігі, график, модель, есеп, кейінге қалдырылған дерекпен тексеру, қате мен шектеу талдауы. Бір балама тәсілді салыстырыңыз.",
            "Соңғы бақылау таныс есепті ғана емес, құралды таңдау, болжамды түсіндіру және шешімді тексеру қабілетін бағалайды."
        ),
        "en": E(
            "Mathematical modeling capstone: from a real question to a testable solution",
            "We assemble the course into one cycle: state a question, choose quantities and assumptions, build a model, calculate a prediction, compare with data, and refine its domain.",
            "Earlier blocks supplied tools. Choosing whether the question needs a number, function, probability, matrix, or graph is now part of the mathematics.",
            "A modeling cycle moves from data and assumptions to model, prediction, validation, and revision in an optimal-bus example",
            "A city wants shorter bus waits without unnecessary vehicles. Before a formula, define the route, observation hours, passenger flow, capacity, cost, and acceptable wait.",
            "A mathematical model is an intentionally simplified representation. It needs a question, inputs, outputs, units, assumptions, equations or an algorithm, parameter data, validation, and a domain of use. Complexity is optional; a testable improvement over a clear baseline is essential. Report prediction error and sensitivity with the answer.",
            "question → data and units → assumptions → model → calculation → prediction → test on new data → revise or reject",
            "For `C(b)=b²−12b+52`, `C′(b)=2b−12`, so the continuous minimum is b=6; `C″(b)=2>0`. Buses are whole, so compare 5, 6, 7: costs are 17, 16, 17. Capacity constraints and new-day validation still matter.",
            "Finding the formula's minimum and declaring six buses the real optimum is wrong. Rush hours, failures, route connections, and sample quality may be absent from the model.",
            "Prepare a mini-project with a question, at least 20 observations, a data dictionary, graph, model, calculation, held-out validation, error and limitation analysis, and one alternative approach.",
            "The final checkpoint tests tool choice, explicit assumptions, and verification rather than repetition of a familiar calculation."
        ),
    },
    "59-data-mastery": {
        "kz": E(
            "Қорытынды бақылау: ықтималдық, дерек, логика, графтар және модельдеу",
            "Он тапсырма нұсқа санауды, шартты ықтималдықты, дерек сипаттауды, статистикалық қорытындыны, дәлелдеуді, қысқа жолды және жеті күннен кейінгі қайталауды тексереді.",
            "Бұл соңғы блок пен бүкіл бағыттың қорытындысы: құралды таңдап, есептеп, болжамды атап, тәуелсіз тексеру қажет.",
            "Соңғы блоктың жеті тірегі қазір 8/10 және жеті күннен кейін 7/10 нәтижесіне апарады",
            "Нақты мәселеде нұсқалар, кездейсоқтық, дерек, логика және желі бірге кездеседі. Сенімді шешім әр санның қайдан келгенін және әр қорытындының шегін сақтайды.",
            "Меңгеру жаңа жағдайда құрылымды танудан көрінеді: математикалық тілді таңдайсыз, нәтижені бағалайсыз және тексеру модельге қайшы болса, қорытындыны өзгертесіз.",
            "нұсқалар → ықтималдық → дерек сипаттамасы → қорытынды → дәлел → граф алгоритмі → модель → тәуелсіз тексеру",
            "Бір жоба маршрутты графпен, сұранысты таңдаумен, белгісіздікті аралықпен және функция параметрін минимуммен біріктіре алады. Байланыс ортақ тақырыппен емес, бірлік пен болжам арқылы түсіндіріледі.",
            "Бір сан шыққан соң есеп шешілді деу қате. Бірліксіз, тексерусіз, белгісіздіксіз және қолдану шартынсыз санды шешімге қауіпсіз қолдануға болмайды.",
            "52–58 сабақтардан он аралас тапсырманы орындаңыз: орналастыру мен теруді салыстырыңыз; қайтармай алу ықтималдығын табыңыз; орташа мен медиананы түсіндіріңіз; аралық құрыңыз; пайымды терістеңіз; дәлел жазыңыз; A–F жолын қалпына келтіріңіз; Дейкстра шегін атаңыз; бүтін минимумды табыңыз; өз моделіңізді тексеріңіз. Кемінде 8/10, жеті күннен кейін 7/10 алыңыз.",
            "Курстан кейін жаңа математикалық бөлімді осы тәуелділік картасындағы тексерілген тіректерге қосыңыз."
        ),
        "en": E(
            "Final checkpoint: probability, data, logic, graphs, and modeling",
            "Ten tasks test counting, conditional probability, data description, inference, proof, shortest paths, and honest modeling with delayed recall after seven days.",
            "This closes the last block and the full route: choose the tool, calculate, state assumptions, and show an independent check.",
            "Seven supports of the final block lead to 8/10 now and 7/10 after seven days",
            "In a real problem, choices, randomness, data, logic, and networks appear together. A reliable solution preserves the source of every number and the boundary of every claim.",
            "Mastery appears in transfer: recognize structure in an unfamiliar setting, choose a mathematical language, evaluate the result, and change the conclusion when validation disagrees.",
            "possibilities → probability → data description → inference → proof → graph algorithm → model → independent check",
            "One project may count routes with a graph, estimate demand from a sample, express uncertainty with an interval, and optimize a function. Units and assumptions must connect the parts.",
            "Calling a task solved as soon as a number appears is wrong. Without units, validation, uncertainty, and conditions of use, a number cannot safely guide a decision.",
            "Complete ten mixed tasks from lessons 52–58: compare permutations and combinations; compute sampling-without-replacement probability; explain mean and median; build an interval; negate a claim; write a proof; reconstruct A–F; name Dijkstra's limit; optimize over integers; and validate your own model. Reach 8/10 now and 7/10 after seven days.",
            "After the course, attach any new area of mathematics to the verified supports on this dependency map."
        ),
    },
}
