-- +goose Up
-- Add a Ludwig Erhard quote and an accountability paragraph to the inflation
-- column (all three languages).
-- +goose StatementBegin
UPDATE article_translations SET body_md = $md$Инфляция — это когда за те же деньги завтра можно купить меньше, чем сегодня. Звучит абстрактно, но чувствует это каждый у кассы. Откуда она берётся? Причин обычно две группы: либо денег в экономике больше, чем товаров, либо сами товары дорожают из-за издержек — топлива, логистики, урожая.

Важно понимать, что инфляция несправедлива по своей природе. Сильнее всего она бьёт по тем, у кого нет «подушки» и чьи доходы фиксированы, — по пенсионерам и небогатым семьям. Дорогие активы вроде недвижимости часто растут в цене вместе с инфляцией, а зарплата и сбережения «на книжке» обесцениваются. Так инфляция тихо перераспределяет — от тех, кто беднее, к тем, кто владеет активами.

Внешние шоки — подорожавшая нефть, сбои логистики, неурожай — действительно случаются. Но устойчивая, длящаяся годами инфляция почти всегда имеет внутреннюю, денежную природу: это следствие того, что в стране печатают больше денег, чем производят товаров. Поэтому, когда руководители списывают рост цен только на «внешние факторы», это часто способ снять с себя ответственность за собственное, нередко бездарное управление. Людвиг Эрхард, архитектор немецкого экономического чуда, был в этом резок:

> «Инфляция — непростительный грех» (нем. *Die Inflation — eine unverzeihliche Sünde*).
> — Людвиг Эрхард

Что помогает? На уровне государства — трезвая денежная политика и предсказуемость. На личном уровне — не держать всё в дешевеющих деньгах, вкладываться в навыки (они дорожают вместе с ценами) и не поддаваться панике, которая сама разгоняет цены. Универсального рецепта нет; есть здравый смысл.

Моё мнение: инфляция — это налог, который никто не голосовал вводить. Поэтому честность власти измеряется не обещаниями «заморозить цены», а способностью объяснять причины и не решать проблему способами, которые завтра станут новой проблемой.$md$
WHERE article_id = 'c2000000-0000-0000-0000-000000000006' AND lang = 'ru';

UPDATE article_translations SET body_md = $md$Инфляция — бұл ертең сол ақшаға бүгінгіден азырақ сатып алуға болатын кез. Дерексіз естіледі, бірақ мұны касса алдында әркім сезеді. Ол қайдан шығады? Себептің әдетте екі тобы бар: не экономикада ақша тауардан көп, не тауардың өзі шығындан — жанармай, логистика, өнім — қымбаттайды.

Инфляцияның табиғаты бойынша әділетсіз екенін түсіну маңызды. Ол «жастығы» жоқ, табысы тұрақты адамдарға — зейнеткерлер мен ауқатты емес отбасыларға — қатты соғады. Жылжымайтын мүлік сияқты қымбат активтер инфляциямен бірге өседі, ал жалақы мен «кітапшадағы» жинақ құнсызданады. Осылай инфляция үнсіз қайта бөледі — кедейлерден активі барларға.

Сыртқы соққылар — қымбаттаған мұнай, логистика ақаулары, өнім тапшылығы — шынымен болады. Бірақ жылдар бойы созылатын тұрақты инфляцияның табиғаты әрдайым дерлік ішкі, ақшалай: бұл елде тауардан көп ақша басып шығарудың салдары. Сондықтан басшылар баға өсуін тек «сыртқы факторларға» сілтегенде, бұл көбіне өз басқаруы, көбінесе дәрменсіз басқаруы үшін жауапкершіліктен қашу тәсілі. Неміс экономикалық кереметінің сәулетшісі Людвиг Эрхард бұған қатты болған:

> «Инфляция — кешірілмес күнә» (нем. *Die Inflation — eine unverzeihliche Sünde*).
> — Людвиг Эрхард

Не көмектеседі? Мемлекет деңгейінде — байсалды ақша саясаты мен болжамдылық. Жеке деңгейде — бәрін құнсызданатын ақшада ұстамау, дағдыға салыну (олар бағамен бірге қымбаттайды) және бағаны өзі үдететін дүрбелеңге берілмеу. Әмбебап рецепт жоқ; парасат бар.

Пікірім: инфляция — енгізуге ешкім дауыс бермеген салық. Сондықтан биліктің адалдығы «бағаны тоқтату» уәделерімен емес, себептерді түсіндіре білуімен және мәселені ертең жаңа мәселеге айналатын жолмен шешпеуімен өлшенеді.$md$
WHERE article_id = 'c2000000-0000-0000-0000-000000000006' AND lang = 'kz';

UPDATE article_translations SET body_md = $md$Inflation is when the same money buys less tomorrow than today. It sounds abstract, but everyone feels it at the checkout. Where does it come from? There are usually two groups of causes: either there is more money in the economy than goods, or the goods themselves grow costlier because of expenses — fuel, logistics, the harvest.

It is important to understand that inflation is unjust by nature. It hits hardest those with no cushion and fixed incomes — pensioners and less wealthy families. Expensive assets like real estate often rise in price along with inflation, while wages and savings in the bank lose value. So inflation quietly redistributes — from the poorer to those who own assets.

External shocks — pricier oil, logistics failures, poor harvests — do happen. But sustained inflation that lasts for years is almost always monetary in nature: the result of a country printing more money than it produces goods. So when leaders blame rising prices only on external factors, it is often a way to dodge responsibility for their own, often inept, management. Ludwig Erhard, the architect of the German economic miracle, was blunt about it:

> "Inflation is an unforgivable sin" (German: *Die Inflation — eine unverzeihliche Sünde*).
> — Ludwig Erhard

What helps? At the state level — sober monetary policy and predictability. At the personal level — not keeping everything in depreciating money, investing in skills (which grow costlier along with prices), and not giving in to panic, which itself drives prices up. There is no universal recipe; there is common sense.

My opinion: inflation is a tax no one voted to introduce. So a government's honesty is measured not by promises to freeze prices but by the ability to explain the causes and not to solve the problem in ways that become a new problem tomorrow.$md$
WHERE article_id = 'c2000000-0000-0000-0000-000000000006' AND lang = 'en';
-- +goose StatementEnd

-- +goose Down
-- No-op: the enriched text supersedes the original seed.
SELECT 1;
