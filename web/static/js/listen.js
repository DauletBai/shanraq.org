// Read the article aloud with the browser's own speech synthesis.
//
// The capability is already in every reader's browser and almost nobody knows
// it is there; this is a button that points at it, not a new service. Nothing
// leaves the machine, there is no key and no per-character bill.
//
// Three things about speechSynthesis decide the shape of this file.
//
// Voices are per-device, not per-browser. Russian and English are on nearly
// everything; Kazakh is on almost nothing. So the control stays hidden until a
// voice for this page's language is actually found: a button that produces
// silence -- or Kazakh spoken with a Russian voice -- is worse than no button.
//
// The voice list arrives asynchronously. Ask too early and it is empty, which
// is why `voiceschanged` is waited for.
//
// Chrome stops speaking after roughly fifteen seconds of a single utterance.
// The text is therefore cut into sentences and queued, which is also what makes
// pause, stop and the moving highlight possible.
(function () {
  'use strict';

  var box = document.getElementById('listen');
  var prose = document.querySelector('.prose');
  var synth = window.speechSynthesis;
  if (!box || !prose || !synth || typeof SpeechSynthesisUtterance === 'undefined') return;

  var lang = (box.getAttribute('data-lang') || 'ru').toLowerCase();
  // The page calls Kazakh "kz"; BCP 47, and therefore every voice, calls it "kk".
  var tag = { kz: 'kk', ru: 'ru', en: 'en' }[lang] || lang;

  var playBtn = document.getElementById('listen-play');
  var stopBtn = document.getElementById('listen-stop');
  var label = document.getElementById('listen-label');
  var speedBox = document.getElementById('listen-speedbox');
  var speed = document.getElementById('listen-speed');
  var note = document.getElementById('listen-note');
  var offer = document.getElementById('listen-offer');
  var offerText = document.getElementById('listen-offer-text');

  var chunks = [];   // {text, el}
  var at = 0;
  var voice = null;
  var playing = false;

  // ---- what is worth saying out loud -------------------------------------

  // Characters that carry meaning to the eye and none to the ear. A backslash,
  // a pipe or a hash read aloud is just noise in the middle of a sentence.
  var MUTE = /[\\|*#~`^<>{}\[\]_]+/g;

  // A bare address is unlistenable: "h t t p s colon slash slash..." for twenty
  // characters. Links keep their visible label, which is prose; the address
  // itself is dropped.
  var URLS = /\bhttps?:\/\/\S+|\bwww\.\S+/gi;

  // The separator between sources reads as "middle dot". It is a pause.
  var DOTS = /\s*[·•]\s*/g;

  // A thousands separator is a space to the eye and a full stop to the voice:
  // "16 700" comes out as "sixteen, seven hundred". Closed up, it is read as
  // the one number it is. Only groups of exactly three digits are joined, so an
  // ordinary sentence is left alone.
  var GROUPS = /(\d)[    ](?=\d{3}\b)/g;

  function clean(s) {
    return s
      .replace(URLS, ' ')
      .replace(MUTE, ' ')
      .replace(GROUPS, '$1')
      .replace(GROUPS, '$1')   // twice: 55 800 000 has two separators
      .replace(DOTS, ', ')
      .replace(/\s+/g, ' ')
      .trim();
  }

  // Blocks that are read, in the order a reader meets them. Headings included:
  // skipping them loses the structure the writer built.
  var BLOCKS = 'p, h2, h3, h4, li, blockquote, figcaption, th, td';

  function collect() {
    var out = [];
    prose.querySelectorAll(BLOCKS).forEach(function (el) {
      // A cell inside a paragraph, a list item inside a quote: read once, at
      // the outermost block that holds it.
      if (el.parentElement && el.parentElement.closest(BLOCKS)) return;
      if (el.closest('pre, code, [aria-hidden="true"]')) return;
      var text = clean(el.textContent || '');
      if (text.length < 2) return;
      // Sentence by sentence: short utterances survive Chrome's timeout and let
      // pause take effect at a natural place rather than mid-word.
      var parts = text.match(/[^.!?…]+[.!?…]+["»)]?\s*|[^.!?…]+$/g) || [text];
      var buf = '';
      parts.forEach(function (part) {
        buf += part;
        if (buf.length >= 160) { out.push({ text: buf.trim(), el: el }); buf = ''; }
      });
      if (buf.trim()) out.push({ text: buf.trim(), el: el });
    });
    return out;
  }

  // ---- voices -------------------------------------------------------------

  function best(list) {
    // A voice the device calls local speaks without the network and without a
    // pause before every sentence.
    var local = list.filter(function (v) { return v.localService; });
    return local[0] || list[0];
  }

  function pickVoice() {
    var all = synth.getVoices() || [];
    if (!all.length) return { voice: null, exact: false };
    var exact = all.filter(function (v) {
      return (v.lang || '').toLowerCase().replace('_', '-').indexOf(tag) === 0;
    });
    if (exact.length) return { voice: best(exact), exact: true };
    // Nothing for this language. Rather than hide the control and leave the
    // reader guessing, offer what the machine does have and let them decide:
    // a Kazakh article read by a Russian voice is a poor substitute, but it is
    // the reader's call whether a poor substitute beats nothing.
    var ru = all.filter(function (v) { return /^ru/i.test(v.lang || ''); });
    var en = all.filter(function (v) { return /^en/i.test(v.lang || ''); });
    var alt = ru.length ? best(ru) : (en.length ? best(en) : all[0]);
    return { voice: alt || null, exact: false };
  }

  var substitute = false;

  function ready() {
    var pick = pickVoice();
    if (!pick.voice) return;       // the machine speaks nothing at all
    voice = pick.voice;
    substitute = !pick.exact;
    chunks = collect();
    if (!chunks.length) return;
    box.hidden = false;
    if (substitute && offer) {
      var tpl = offer.getAttribute('data-substitute') || '%s';
      offerText.textContent = tpl.replace('%s', voice.name || voice.lang);
      offer.hidden = false;
      playBtn.hidden = true;       // the substitute is accepted, not assumed
    }
  }

  // The list is empty on first call in most browsers and arrives later.
  ready();
  if (!voice && typeof synth.onvoiceschanged !== 'undefined') {
    synth.addEventListener('voiceschanged', function once() {
      synth.removeEventListener('voiceschanged', once);
      ready();
    });
  }

  // ---- how far they listened ---------------------------------------------
  //
  // The same milestones the scroll tracker reports, down the same endpoint,
  // marked as listening. Kept apart in the figures because the two mean
  // different things: scrolling away at half is giving up, stopping a recording
  // at half may be arriving at work.
  //
  // Milestones are remembered for the session, so replaying an article cannot
  // report the same listener twice -- the same guard the reading beacon uses,
  // and for the same reason.
  var slug = prose.getAttribute('data-read-progress') || '';
  var heard = {};
  var memo = 'shanraq-listen:' + slug;
  try { heard = JSON.parse(sessionStorage.getItem(memo) || '{}'); } catch (e) {}

  function reached(i) {
    if (!slug || !navigator.sendBeacon || !chunks.length) return;
    var pct = ((i + 1) / chunks.length) * 100;
    [25, 50, 75, 100].forEach(function (mark) {
      if (pct < mark || heard[mark]) return;
      heard[mark] = true;
      try { sessionStorage.setItem(memo, JSON.stringify(heard)); } catch (e) {}
      try {
        navigator.sendBeacon('/read/' + encodeURIComponent(slug) + '/progress?d=' + mark + '&m=listen');
      } catch (e) {}
    });
  }

  // ---- speaking -----------------------------------------------------------

  function mark(el) {
    prose.querySelectorAll('.listen-now').forEach(function (n) { n.classList.remove('listen-now'); });
    if (el) el.classList.add('listen-now');
  }

  function speak(i) {
    if (i >= chunks.length) { finish(); return; }
    at = i;
    var u = new SpeechSynthesisUtterance(chunks[i].text);
    // The language is always stated; the voice object only when it is not the
    // engine's own choice for that language. A pinned voice can go stale
    // between page loads, and a stale one is spoken by nobody.
    u.lang = voice.lang;
    if (substitute || pinned) u.voice = voice;
    u.rate = parseFloat(speed.value) || 1;
    u.onstart = function () { started = true; mark(chunks[i].el); say(''); };
    u.onend = function () { reached(i); if (playing) speak(i + 1); };
    // One failure stops the reading. It used to step to the next sentence,
    // which meant a browser refusing to speak at all raced through the whole
    // article in silence: the reader pressed the button and nothing happened,
    // with nothing said about why.
    u.onerror = function (ev) { fail(ev && ev.error); };
    started = false;
    // Chrome can leave the queue paused after a navigation, and everything
    // spoken afterwards waits behind it for ever. resume() unsticks that.
    // cancel() is deliberately NOT called here: it is asynchronous, and an
    // utterance handed over in the same breath is swallowed by it -- which is
    // its own way of producing silence.
    try { synth.resume(); } catch (e) {}
    synth.speak(u);
    watch();
  }

  // Speech that neither starts nor errors is the worst case: silence with
  // nothing to explain it. If nothing has begun within a few seconds, say so.
  var started = false, timer = null, pinned = false;
  var retried = false;
  function watch() {
    clearTimeout(timer);
    timer = setTimeout(function () {
      if (!playing || started) return;
      // Some browsers register the gesture only on the second attempt. One
      // retry is worth more than a message; a second would just be noise.
      if (!retried) {
        retried = true;
        try { synth.cancel(); } catch (e) {}
        setTimeout(function () { if (playing && !started) speak(at); }, 120);
        return;
      }
      fail('silent');
    }, 2500);
  }

  function say(msg) {
    if (!note) return;
    note.textContent = msg || '';
    note.hidden = !msg;
  }

  // Where the voices actually live, for the machine in front of the reader.
  // "Check your settings" is useless advice: they have already looked, and
  // found nothing, because voices belong to the operating system and not to
  // the browser. This says which panel, by name.
  function whereVoicesLive() {
    var ua = navigator.userAgent;
    if (/iPhone|iPad|iPod/.test(ua)) return note.getAttribute('data-ios');
    if (/Android/.test(ua)) return note.getAttribute('data-android');
    if (/Windows/.test(ua)) return note.getAttribute('data-win');
    if (/Mac OS X|Macintosh/.test(ua)) return note.getAttribute('data-mac');
    return note.getAttribute('data-android');
  }

  function fail(reason) {
    clearTimeout(timer);
    playing = false;
    try { synth.cancel(); } catch (e) {}
    setPlaying(false);
    mark(null);
    if (!note) return;

    var parts = [];
    parts.push(note.getAttribute(reason === 'not-allowed' ? 'data-blocked' : 'data-failed'));
    // Naming the voice we found settles the commonest wrong guess: the reader
    // stops hunting for a missing voice that is in fact installed.
    if (voice) {
      var found = note.getAttribute('data-found') || '';
      if (found) parts.push(found.replace('%s', voice.name || voice.lang));
    }
    // Chrome alone can mute a single site, and that switch is the one nobody
    // finds. Chromium-based browsers report themselves as Chrome; Safari does
    // not, and has no such switch.
    if (/Chrome\//.test(navigator.userAgent) && !/Edg\//.test(navigator.userAgent)) {
      parts.push(note.getAttribute('data-sound'));
    }
    parts.push(whereVoicesLive());
    say(parts.filter(Boolean).join(' '));
  }

  // Chrome pauses a long reading of its own accord; a nudge every few seconds
  // keeps it going. Harmless when nothing is speaking.
  setInterval(function () {
    if (playing && synth.paused) { try { synth.resume(); } catch (e) {} }
  }, 5000);

  function setPlaying(on) {
    playing = on;
    label.textContent = box.getAttribute(on ? 'data-pause' : 'data-play');
    stopBtn.hidden = !on && at === 0;
    speedBox.hidden = !on && at === 0;
    playBtn.setAttribute('aria-pressed', String(on));
  }

  function finish() {
    clearTimeout(timer);
    say('');
    synth.cancel();
    setPlaying(false);
    at = 0;
    mark(null);
    stopBtn.hidden = true;
    speedBox.hidden = true;
  }

  playBtn.addEventListener('click', function () {
    if (playing) {
      // Chrome's pause() is unreliable on long queues; cancelling and resuming
      // from the current sentence is what actually stops and restarts cleanly.
      playing = false;
      synth.cancel();
      setPlaying(false);
      return;
    }
    say('');
    retried = false;
    setPlaying(true);
    speak(at);
  });

  if (document.getElementById('listen-accept')) {
    document.getElementById('listen-accept').addEventListener('click', function () {
      // Consent given: from here the substitute voice is used deliberately.
      pinned = true;
      offer.hidden = true;
      playBtn.hidden = false;
      playBtn.click();
    });
  }

  stopBtn.addEventListener('click', finish);

  speed.addEventListener('change', function () {
    if (!playing) return;
    synth.cancel();
    speak(at);           // the new rate takes effect from the current sentence
  });

  // Leaving the page with a voice still talking would follow the reader to the
  // next one.
  window.addEventListener('beforeunload', function () { synth.cancel(); });
  window.addEventListener('pagehide', function () { synth.cancel(); });
})();
