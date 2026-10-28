/* First-party aggregate measurements for the published online Go excerpt. */
(function () {
  if (navigator.webdriver || !navigator.sendBeacon) return;
  var sent = false, fired = {}, article = document.querySelector('main');
  function view() {
    if (!sent && document.visibilityState === 'visible') {
      sent = navigator.sendBeacon('/api/view');
    }
  }
  function progress() {
    if (!article || document.visibilityState !== 'visible') return;
    var r = article.getBoundingClientRect();
    if (r.top >= innerHeight) return;
    var pct = (innerHeight - r.top) / (r.height || 1) * 100;
    [25, 50, 75, 100].forEach(function (n) {
      if (pct >= n && !fired[n]) {
        fired[n] = true;
        navigator.sendBeacon('/api/commercial?e=depth&d=' + n);
      }
    });
  }
  if (document.visibilityState === 'visible') setTimeout(view, 250);
  document.addEventListener('visibilitychange', function () { if (document.visibilityState === 'visible') { view(); progress(); } });
  window.addEventListener('scroll', progress, {passive: true});
  window.addEventListener('resize', progress, {passive: true});
  progress();
})();
