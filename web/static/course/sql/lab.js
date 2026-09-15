import sqlite3InitModule from "/static/vendor/sqlite-wasm/index.mjs";

const copy = {
  ru:{title:"SQL-лаборатория семейного бюджета",lead:"Запускайте запросы прямо на телефоне, планшете или компьютере. База хранится в этом браузере и не отправляется на Shanraq.",privacyTitle:"Важно:",privacy:"учебные данные вымышлены. Не вводите реальные финансовые сведения на общем устройстве. Очистка данных браузера может удалить базу — сохраняйте резервную копию.",example:"Готовый пример",exampleMonth:"Итог месяца",exampleLargest:"Крупные расходы",exampleCategories:"Расходы по категориям",exampleTables:"Список таблиц",run:"Выполнить",reset:"Вернуть учебную базу",backup:"Скачать backup",restore:"Восстановить из .sqlite",loading:"Запускаем SQLite…",ready:"Готово",saved:"База сохранена на устройстве",temporary:"Временная база: скачайте backup",empty:"Запрос выполнен. Строк результата нет.",helpTitle:"Как пользоваться на телефоне",help1:"Выберите пример или измените запрос в поле.",help2:"Нажмите «Выполнить» — результат появится под редактором.",help3:"Чтобы продолжить позже, скачайте backup и восстановите его в следующем сеансе.",confirm:"Вернуть исходные учебные данные? Ваши изменения будут удалены.",restored:"Резервная копия восстановлена.",resetDone:"Учебная база восстановлена."},
  kz:{title:"Отбасы бюджетіне арналған SQL зертханасы",lead:"Сұрауларды телефонда, планшетте немесе компьютерде орындаңыз. Дерекқор осы браузерде сақталады және Shanraq серверіне жіберілмейді.",privacyTitle:"Маңызды:",privacy:"оқу деректері ойдан алынған. Ортақ құрылғыға шынайы қаржы мәліметтерін енгізбеңіз. Браузер деректерін тазаласаңыз, дерекқор жойылуы мүмкін — сақтық көшірме жасаңыз.",example:"Дайын мысал",exampleMonth:"Ай қорытындысы",exampleLargest:"Ірі шығындар",exampleCategories:"Санаттар бойынша шығын",exampleTables:"Кестелер тізімі",run:"Орындау",reset:"Оқу дерекқорын қалпына келтіру",backup:"Сақтық көшірмені жүктеу",restore:".sqlite файлынан қалпына келтіру",loading:"SQLite іске қосылып жатыр…",ready:"Дайын",saved:"Дерекқор құрылғыда сақталады",temporary:"Уақытша дерекқор: сақтық көшірмені жүктеңіз",empty:"Сұрау орындалды. Нәтиже жолдары жоқ.",helpTitle:"Телефонда қалай пайдалану керек",help1:"Дайын мысалды таңдаңыз немесе өрістегі сұрауды өзгертіңіз.",help2:"«Орындау» батырмасын басыңыз — нәтиже редактордың астында шығады.",help3:"Кейін жалғастыру үшін сақтық көшірмені жүктеп, келесі сеанста қалпына келтіріңіз.",confirm:"Бастапқы оқу деректерін қайтару керек пе? Өзгерістеріңіз жойылады.",restored:"Сақтық көшірме қалпына келтірілді.",resetDone:"Оқу дерекқоры қалпына келтірілді."},
  en:{title:"Family-budget SQL lab",lead:"Run queries on a phone, tablet, or computer. The database stays in this browser and is not sent to Shanraq.",privacyTitle:"Important:",privacy:"the sample data is fictional. Do not enter real financial details on a shared device. Clearing browser data may erase the database, so download backups.",example:"Ready-made example",exampleMonth:"Monthly summary",exampleLargest:"Largest expenses",exampleCategories:"Spending by category",exampleTables:"List tables",run:"Run",reset:"Reset sample database",backup:"Download backup",restore:"Restore from .sqlite",loading:"Starting SQLite…",ready:"Ready",saved:"Database saved on this device",temporary:"Temporary database: download backups",empty:"Query completed. It returned no rows.",helpTitle:"Using the lab on a phone",help1:"Choose an example or edit the query in the field.",help2:"Tap Run; the result appears below the editor.",help3:"To continue later, download a backup and restore it in the next session.",confirm:"Reset to the original sample data? Your changes will be deleted.",restored:"Backup restored.",resetDone:"Sample database restored."}
};
const examples={month:`SELECT c.kind, SUM(t.amount_tiyn) / 100.0 AS amount_kzt
FROM transactions AS t JOIN categories AS c ON c.id=t.category_id
WHERE t.happened_on >= '2026-09-01' AND t.happened_on < '2026-10-01'
GROUP BY c.kind;`,largest:`SELECT happened_on, note, amount_tiyn / 100.0 AS amount_kzt
FROM transactions
WHERE category_id <> 1
ORDER BY amount_tiyn DESC LIMIT 5;`,categories:`SELECT c.name, SUM(t.amount_tiyn) / 100.0 AS amount_kzt
FROM transactions AS t JOIN categories AS c ON c.id=t.category_id
WHERE c.kind='expense' GROUP BY c.id, c.name ORDER BY amount_kzt DESC;`,tables:`SELECT name FROM sqlite_schema WHERE type='table' ORDER BY name;`};
const schema=`PRAGMA foreign_keys=ON;
CREATE TABLE accounts(id INTEGER PRIMARY KEY,name TEXT NOT NULL UNIQUE,opening_balance_tiyn INTEGER NOT NULL DEFAULT 0);
CREATE TABLE categories(id INTEGER PRIMARY KEY,name TEXT NOT NULL UNIQUE,kind TEXT NOT NULL CHECK(kind IN('income','expense')),monthly_limit_tiyn INTEGER CHECK(monthly_limit_tiyn IS NULL OR monthly_limit_tiyn>=0));
CREATE TABLE transactions(id INTEGER PRIMARY KEY,account_id INTEGER NOT NULL REFERENCES accounts(id),category_id INTEGER NOT NULL REFERENCES categories(id),happened_on TEXT NOT NULL CHECK(date(happened_on) IS NOT NULL AND happened_on=date(happened_on)),amount_tiyn INTEGER NOT NULL CHECK(amount_tiyn>0),note TEXT,created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP);
CREATE INDEX transactions_month_category ON transactions(happened_on,category_id);
INSERT INTO accounts VALUES(1,'Card',18000000),(2,'Cash',2500000);
INSERT INTO categories VALUES(1,'Salary','income',NULL),(2,'Food','expense',9000000),(3,'Home','expense',6500000),(4,'Transport','expense',3500000),(5,'Learning','expense',2500000);
INSERT INTO transactions(id,account_id,category_id,happened_on,amount_tiyn,note) VALUES(1,1,1,'2026-09-01',42000000,'September salary'),(2,1,3,'2026-09-02',5500000,'Utilities and rent'),(3,1,2,'2026-09-03',1865000,'Groceries'),(4,2,4,'2026-09-04',420000,'Bus card'),(5,1,5,'2026-09-06',1200000,'Books'),(6,2,2,'2026-09-08',975000,'Market'),(7,1,4,'2026-09-10',680000,'Fuel'),(8,1,2,'2026-09-12',2140000,'Groceries');`;
const params=new URLSearchParams(location.search);const lang=["kz","ru","en"].includes(params.get("lang"))?params.get("lang"):"ru";const t=copy[lang];document.documentElement.lang=lang;document.title=`${t.title} · Shanraq`;document.querySelectorAll("[data-t]").forEach(el=>el.textContent=t[el.dataset.t]);document.querySelectorAll("[data-lang]").forEach(el=>el.classList.toggle("active",el.dataset.lang===lang));document.querySelector(".lab-brand").href=`/course/sql?lang=${lang}`;
let sqlite3;
let db;
const importedFilename = "/shanraq-budget.sqlite3";

function execute(statement) {
  const resultSets = [];
  db.exec({
    sql: statement,
    rowMode: "array",
    callback: (row, stmt) => {
      let current = resultSets.at(-1);
      const columns = stmt.getColumnNames();
      if (!current) {
        current = {columns, rows: []};
        resultSets.push(current);
      }
      current.rows.push(row.map(value => typeof value === "bigint" ? value.toString() : value));
    }
  });
  return resultSets;
}

async function call(type, data={}) {
  if (type === "open") {
    sqlite3 = await sqlite3InitModule();
    db = new sqlite3.oo1.DB(":memory:", "c");
    db.exec("PRAGMA foreign_keys=ON");
    return {persistent: false, version: sqlite3.version.libVersion};
  }
  if (type === "exec") return execute(data.sql);
  if (type === "reset") {
    db.exec("PRAGMA foreign_keys=OFF; DROP TABLE IF EXISTS transactions; DROP TABLE IF EXISTS categories; DROP TABLE IF EXISTS accounts; PRAGMA foreign_keys=ON;");
    return null;
  }
  if (type === "export") return sqlite3.capi.sqlite3_js_db_export(db.pointer);
  if (type === "import") {
    db.close();
    sqlite3.capi.sqlite3_js_posix_create_file(importedFilename, new Uint8Array(data.bytes));
    db = new sqlite3.oo1.DB(importedFilename, "c");
    db.exec("PRAGMA foreign_keys=ON");
    return null;
  }
  throw new Error(`Unknown lab operation: ${type}`);
}
const status=document.querySelector("#status"),engine=document.querySelector("#engine"),result=document.querySelector("#result"),sql=document.querySelector("#sql");function render(sets){result.replaceChildren();if(!sets.length){const p=document.createElement("p");p.className="empty";p.textContent=t.empty;result.append(p);return;}for(const set of sets){const table=document.createElement("table"),head=document.createElement("thead"),hr=document.createElement("tr");set.columns.forEach(name=>{const th=document.createElement("th");th.textContent=name;hr.append(th);});head.append(hr);table.append(head);const body=document.createElement("tbody");set.rows.forEach(row=>{const tr=document.createElement("tr");row.forEach(value=>{const td=document.createElement("td");td.textContent=value===null?"NULL":String(value);tr.append(td);});body.append(tr);});table.append(body);result.append(table);}}function fail(error){result.innerHTML="";const p=document.createElement("pre");p.className="error";p.textContent=error.message;result.append(p);status.textContent="SQL error";}
async function reset(first=false){if(!first&&!confirm(t.confirm))return;await call("reset");await call("exec",{sql:schema});status.textContent=first?t.ready:t.resetDone;render(await call("exec",{sql:examples.month}));}
document.querySelector("#example").onchange=e=>{sql.value=examples[e.target.value];};document.querySelector("#run").onclick=async()=>{try{status.textContent="…";render(await call("exec",{sql:sql.value}));status.textContent=t.ready;}catch(e){fail(e);}};document.querySelector("#reset").onclick=()=>reset();document.querySelector("#backup").onclick=async()=>{try{const bytes=await call("export"),url=URL.createObjectURL(new Blob([bytes],{type:"application/vnd.sqlite3"})),a=document.createElement("a");a.href=url;a.download="shanraq-budget.sqlite";a.click();URL.revokeObjectURL(url);}catch(e){fail(e);}};document.querySelector("#restore").onchange=async e=>{try{const bytes=await e.target.files[0].arrayBuffer();await call("import",{bytes});status.textContent=t.restored;render(await call("exec",{sql:examples.month}));}catch(err){fail(err);}finally{e.target.value="";}};
try{const info=await call("open");engine.textContent=`SQLite ${info.version}`;status.textContent=info.persistent?t.saved:t.temporary;let tables=await call("exec",{sql:"SELECT count(*) AS n FROM sqlite_schema WHERE type='table' AND name='transactions';"});if(!tables[0]?.rows[0]?.[0])await reset(true);else render(await call("exec",{sql:examples.month}));document.querySelector(".lab-result-panel").setAttribute("aria-busy","false");}catch(e){fail(e);}
