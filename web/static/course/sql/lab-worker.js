import sqlite3InitModule from "/static/vendor/sqlite-wasm/index.mjs";

let sqlite3;
let db;
const filename = "/shanraq-budget.sqlite3";

function run(sql) {
  const resultSets = [];
  db.exec({
    sql,
    rowMode: "array",
    callback: (row, statement) => {
      let current = resultSets.at(-1);
      const columns = statement.getColumnNames();
      if (!current || current.sql !== statement.getSQL()) {
        current = { sql: statement.getSQL(), columns, rows: [] };
        resultSets.push(current);
      }
      current.rows.push(row.map(value => typeof value === "bigint" ? value.toString() : value));
    }
  });
  return resultSets;
}

async function open() {
  sqlite3 = await sqlite3InitModule();
  const persistent = false;
  db = new sqlite3.oo1.DB(":memory:", "c");
  db.exec("PRAGMA foreign_keys = ON");
  return { persistent, version: sqlite3.version.libVersion };
}

self.onmessage = async event => {
  const { id, type, sql, bytes } = event.data;
  try {
    if (type === "open") self.postMessage({ id, ok: true, data: await open() });
    else if (type === "exec") self.postMessage({ id, ok: true, data: run(sql) });
    else if (type === "reset") {
      db.exec("PRAGMA foreign_keys=OFF; DROP TABLE IF EXISTS transactions; DROP TABLE IF EXISTS categories; DROP TABLE IF EXISTS accounts; PRAGMA foreign_keys=ON;");
      self.postMessage({ id, ok: true, data: null });
    } else if (type === "export") {
      const data = sqlite3.capi.sqlite3_js_db_export(db.pointer);
      self.postMessage({ id, ok: true, data }, [data.buffer]);
    } else if (type === "import") {
      db.close();
      sqlite3.capi.sqlite3_js_posix_create_file(filename, bytes);
      db = new sqlite3.oo1.DB(filename, "c");
      db.exec("PRAGMA foreign_keys = ON");
      self.postMessage({ id, ok: true, data: null });
    }
  } catch (error) {
    self.postMessage({ id, ok: false, error: String(error?.message || error) });
  }
};
