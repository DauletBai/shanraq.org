// Loan calculator.
//
// The reference implementation is LoanCalc in pkg/modules/articles/loans.go and
// that is the one under test; this mirrors it so the totals follow every
// keystroke without a round trip. Change one, change the other -- the worked
// example on the page is rendered by the Go side, so a divergence shows up as
// the same inputs giving two different answers on one page.
//
// It lives in a file rather than inline so the browser can cache it and so no
// translated string ever has to survive html/template's JavaScript escaping.
// Everything it needs -- programmes, market rates, labels -- arrives in data
// attributes, read once at start-up.
(function () {
  'use strict';

  var form = document.getElementById('calc');
  if (!form) return;

  function attr(name, fallback) {
    try { return JSON.parse(form.getAttribute(name)) || fallback; } catch (e) { return fallback; }
  }
  var progs = attr('data-programs', []);
  var market = attr('data-market', {});
  var L = attr('data-labels', {});

  var $ = function (id) { return document.getElementById(id); };
  var nf = new Intl.NumberFormat('ru-RU');
  var money = function (v) { return nf.format(Math.round(v)) + ' ₸'; };
  // Money fields are text so they can carry thousands separators while the
  // reader types: a bare 20000000 is unreadable, and type=number forbids the
  // spaces that make it legible.
  var num = function (id) {
    var el = $(id);
    if (!el) return 0;
    var v = parseFloat(String(el.value).replace(/[^0-9.,-]/g, '').replace(',', '.'));
    return isFinite(v) ? v : 0;
  };

  function groupMoneyFields() {
    form.querySelectorAll('.calc__money').forEach(function (el) {
      el.addEventListener('input', function () {
        var caretFromEnd = el.value.length - (el.selectionStart || el.value.length);
        var digits = el.value.replace(/[^0-9]/g, '');
        el.value = digits ? nf.format(parseInt(digits, 10)) : '';
        var pos = Math.max(0, el.value.length - caretFromEnd);
        try { el.setSelectionRange(pos, pos); } catch (e) {}
      });
    });
  }
  var kind = 'mortgage';

  // Kinds that buy a thing take a price and a down payment; a cash loan takes
  // the amount straight.
  function hasPrice(k) { return k === 'mortgage' || k === 'auto' || k === 'installment'; }

  // Interest a schedule of this shape carries at a given rate. Used twice: for
  // the loan itself, and for the market rate a subsidy is measured against.
  function interestAt(loan, rate, months, diff) {
    var i = rate / 100 / 12;
    var flat = loan / months;
    var base = 0;
    if (!diff) base = i <= 0 ? flat : loan * i / (1 - Math.pow(1 + i, -months));
    var bal = loan, total = 0;
    for (var k = 1; k <= months; k++) {
      var int_ = Math.round(bal * i);
      var pr = diff ? Math.round(flat) : Math.round(base) - int_;
      if (k === months || pr > bal) pr = bal;
      if (pr < 0) pr = 0;
      bal -= pr;
      total += int_;
    }
    return total;
  }

  function calc() {
    var months = Math.min(600, Math.max(0, Math.round(num('c-months'))));
    var rate = Math.max(0, num('c-rate'));
    var price = 0, down = 0, loan = 0;
    if (hasPrice(kind)) {
      price = Math.max(0, num('c-price'));
      down = Math.min(Math.max(0, num('c-down')), price);
      loan = price - down;
    } else {
      loan = Math.max(0, num('c-amount'));
    }
    if (loan <= 0 || months <= 0) return null;

    var i = rate / 100 / 12;
    var feeM = Math.round(loan * num('c-feem') / 100);
    var insM = Math.round(loan * num('c-ins') / 100 / 12);
    var per = feeM + insM;
    var feeOnce = Math.max(0, num('c-fee1f') + Math.round(loan * num('c-fee1p') / 100));
    var scheme = form.querySelector('input[name=scheme]:checked');
    var diff = scheme && scheme.value === 'differentiated';

    var flat = loan / months;
    var base = 0;
    if (!diff) base = i <= 0 ? flat : loan * i / (1 - Math.pow(1 + i, -months));

    var bal = loan, rows = [], interest = 0, paid = 0;
    for (var k = 1; k <= months; k++) {
      var int_ = Math.round(bal * i);
      var pr = diff ? Math.round(flat) : Math.round(base) - int_;
      if (k === months || pr > bal) pr = bal;
      if (pr < 0) pr = 0;
      bal -= pr;
      var pay = pr + int_ + per;
      rows.push({ n: k, pay: pay, pr: pr, int: int_, fee: per, bal: bal });
      interest += int_;
      paid += pay;
    }
    var bankFees = feeM * months + feeOnce;
    var insurer = insM * months;
    paid += feeOnce;

    // The effective rate: the rate at which what is repaid discounts back to
    // what was actually received. Bisection, because it is monotone on a
    // bracket that cannot be escaped and converges from any input a form makes.
    var net = loan - feeOnce, eff = 0;
    if (net > 0) {
      var npv = function (m) {
        var sum = 0;
        for (var j = 0; j < rows.length; j++) sum += rows[j].pay / Math.pow(1 + m, rows[j].n);
        return sum - net;
      };
      if (npv(0) > 0) {
        var lo = 0, hi = 1;
        while (npv(hi) > 0 && hi < 10) hi *= 2;
        for (var s = 0; s < 200; s++) {
          var mid = (lo + hi) / 2;
          if (npv(mid) > 0) lo = mid; else hi = mid;
        }
        eff = (Math.pow(1 + (lo + hi) / 2, 12) - 1) * 100;
      }
    }

    // A subsidised rate is the same loan with part of the interest paid from
    // the budget. Measured against the dearest rate on offer for this kind.
    var subsidy = 0, mkt = market[kind] || 0;
    if (mkt > rate) {
      var full = interestAt(loan, mkt, months, diff);
      if (full > interest) subsidy = full - interest;
    }

    return {
      loan: loan, rows: rows, interest: interest, paid: paid,
      overpay: paid - loan, total: down + paid,
      mul: price > 0 ? (down + paid) / price : 0,
      eff: eff, diff: diff,
      seller: price > 0 ? price : loan,
      bankFees: bankFees, insurer: insurer, subsidy: subsidy
    };
  }

  function clear() {
    ['o-total', 'o-eff', 'o-pay', 'o-loan', 'o-paid'].forEach(function (id) {
      if ($(id)) $(id).textContent = '—';
    });
    ['o-mul', 'o-over', 'o-paylast'].forEach(function (id) { if ($(id)) $(id).textContent = ''; });
    if ($('o-rows')) $('o-rows').innerHTML = '';
    if ($('o-who')) $('o-who').innerHTML = '';
  }

  function render() {
    var r = calc();
    if (!r) { clear(); return; }

    $('o-total').textContent = money(r.total);
    $('o-mul').textContent = r.mul > 0 ? (L.multiple || '%s').replace('%s', r.mul.toFixed(1)) : '';
    $('o-eff').textContent = r.eff.toFixed(2) + ' %';
    $('o-paylab').textContent = r.diff ? (L.first || '') : (L.monthly || '');
    $('o-pay').textContent = money(r.rows[0].pay);
    $('o-paylast').textContent = r.diff
      ? (L.last || '') + ': ' + money(r.rows[r.rows.length - 1].pay)
      : '';
    $('o-loan').textContent = money(r.loan);
    $('o-paid').textContent = money(r.paid);
    // Only the parts that exist: "комиссии и страховка 0 ₸" is noise, and
    // repeating the same sum twice reads like an error.
    var extra = r.bankFees + r.insurer;
    var over = (L.overpay || '') + ' ' + money(r.overpay);
    if (extra > 0) {
      over += ': ' + (L.interest || '') + ' ' + money(r.interest) +
              ', ' + (L.fees || '') + ' ' + money(extra);
    }
    $('o-over').textContent = over;

    // Who ends up with the money. Only the seller's line buys something that
    // exists; the budget's line is money nobody in this room paid.
    var who = [
      [L.seller, r.seller, 'seller'],
      [L.wint, r.interest, ''],
      [L.wfees, r.bankFees, ''],
      [L.wins, r.insurer, ''],
      [L.subsidy, r.subsidy, 'budget']
    ];
    var out = '';
    for (var w = 0; w < who.length; w++) {
      if (!(who[w][1] > 0)) continue;
      var pct = r.total > 0 ? Math.round(who[w][1] / r.total * 100) : 0;
      out += '<li class="calc__who-row' + (who[w][2] ? ' calc__who-row--' + who[w][2] : '') + '">' +
        '<span class="calc__who-lab"></span>' +
        '<span class="calc__who-bar"><i style="width:' + Math.min(100, pct) + '%"></i></span>' +
        '<span class="calc__who-sum"></span></li>';
    }
    out += '<li class="calc__who-row calc__who-row--you">' +
      '<span class="calc__who-lab"></span><span class="calc__who-bar"></span>' +
      '<span class="calc__who-sum"></span></li>';
    $('o-who').innerHTML = out;

    // Labels and sums are written as text, never interpolated into the markup:
    // they come from a data attribute and must not be able to carry tags.
    var items = $('o-who').querySelectorAll('.calc__who-row');
    var shown = [];
    for (var x = 0; x < who.length; x++) if (who[x][1] > 0) shown.push(who[x]);
    shown.push([L.you, r.total, 'you']);
    for (var y = 0; y < items.length && y < shown.length; y++) {
      items[y].querySelector('.calc__who-lab').textContent = shown[y][0] || '';
      items[y].querySelector('.calc__who-sum').textContent = money(shown[y][1]);
    }
    $('o-subnote').hidden = !(r.subsidy > 0);

    var rowsOut = '';
    for (var j = 0; j < r.rows.length; j++) {
      var q = r.rows[j];
      rowsOut += '<tr><td>' + q.n + '</td><td>' + money(q.pay) + '</td><td>' + money(q.pr) +
        '</td><td>' + money(q.int) + '</td><td>' + money(q.fee) + '</td><td>' + money(q.bal) + '</td></tr>';
    }
    $('o-rows').innerHTML = rowsOut;
  }

  function fillPrograms() {
    var sel = $('c-prog');
    var own = sel.options[0];
    sel.innerHTML = '';
    sel.appendChild(own);
    progs.filter(function (p) { return p.kind === kind; }).forEach(function (p) {
      var o = document.createElement('option');
      o.value = p.code;
      o.textContent = p.lender ? p.name + ' · ' + p.lender : p.name;
      sel.appendChild(o);
    });
    form.querySelectorAll('[data-when=price]').forEach(function (el) { el.hidden = !hasPrice(kind); });
    form.querySelectorAll('[data-when=amount]').forEach(function (el) { el.hidden = hasPrice(kind); });
  }

  form.querySelectorAll('.calc__kind').forEach(function (b) {
    b.addEventListener('click', function () {
      kind = b.getAttribute('data-kind');
      form.querySelectorAll('.calc__kind').forEach(function (x) {
        x.setAttribute('aria-pressed', String(x === b));
      });
      fillPrograms();
      $('c-prognote').hidden = true;
      render();
    });
  });

  $('c-prog').addEventListener('change', function () {
    var code = $('c-prog').value;
    var p = progs.filter(function (x) { return x.code === code; })[0];
    var note = $('c-prognote');
    if (!p) { note.hidden = true; render(); return; }
    $('c-rate').value = p.rate;
    $('c-months').value = p.months;
    if (hasPrice(kind) && p.down > 0) {
      $('c-down').value = nf.format(Math.round(Math.max(0, num('c-price')) * p.down / 100));
    }
    note.textContent = p.note || '';
    note.hidden = !p.note;
    render();
  });

  form.addEventListener('input', render);
  form.addEventListener('change', render);

  var tog = $('o-toggle');
  if (tog) {
    tog.addEventListener('click', function () {
      var box = $('o-sched');
      var open = box.hidden;
      box.hidden = !open;
      tog.textContent = tog.getAttribute(open ? 'data-hide' : 'data-show');
    });
  }

  groupMoneyFields();
  var firstKind = form.querySelector('.calc__kind');
  if (firstKind) firstKind.setAttribute('aria-pressed', 'true');
  fillPrograms();
  render();
})();
