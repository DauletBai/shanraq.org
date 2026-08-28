// Play the recorded reading of an article, and show where it has got to.
//
// This used to drive the browser's own speech synthesis, and that is gone. It
// had no Kazakh voice on any system, so a Kazakh article was read in Russian;
// it fell silent when a phone locked; and when it failed there was nothing
// useful to tell the reader, because the cause lived in their operating system
// rather than on this page. Every one of those was a dead end for the person
// who just wanted the article read to them.
//
// What is left is a file and a cue map. The file plays the way any audio plays.
// The cue map says which paragraph is sounding, so the page can highlight it
// and scroll to follow -- the point being that you can be doing something else
// and still see where the reading is.
//
// The controls appear only when a recording exists. There is nothing to offer
// otherwise, and a button that produces silence is worse than no button.
(function () {
  'use strict';

  var box = document.getElementById('listen');
  var prose = document.querySelector('.prose');
  var audio = document.getElementById('listen-audio');
  if (!box || !prose || !audio) return;

  var playBtn = document.getElementById('listen-play');
  var stopBtn = document.getElementById('listen-stop');
  var label = document.getElementById('listen-label');
  var note = document.getElementById('listen-note');
  var playTxt = box.getAttribute('data-play') || '';
  var pauseTxt = box.getAttribute('data-pause') || '';

  // The same block rule the recording was cut along, and the same one in
  // scripts/narrate.py: outermost blocks only, nothing inside code, nothing
  // hidden, nothing shorter than a word. Cue N and block N must be the same
  // paragraph. Change it in one place and the highlight drifts from the sound.
  var BLOCK_SEL = 'p, h2, h3, h4, li, blockquote, figcaption, th, td';
  var els = [];
  prose.querySelectorAll(BLOCK_SEL).forEach(function (el) {
    if (el.parentElement && el.parentElement.closest(BLOCK_SEL)) return;
    if (el.closest('pre, code, [aria-hidden="true"]')) return;
    if ((el.textContent || '').trim().length < 2) return;
    els.push(el);
  });

  var cues = [];
  try { cues = JSON.parse(audio.getAttribute('data-cues') || '[]') || []; } catch (e) {}

  box.hidden = false;

  // ---- saying what went wrong ---------------------------------------------
  //
  // Silence with no explanation is what sent readers hunting through settings
  // the last time, and there was nothing there for them to find. Now the causes
  // are ours -- a missing file, a format the browser will not decode, a network
  // that dropped -- so each of them gets said out loud.
  function say(msg) {
    if (!note) return;
    note.textContent = msg || '';
    note.hidden = !msg;
  }

  function mediaError() {
    var e = audio.error;
    var code = e && e.code;
    if (code === 2) return note && note.getAttribute('data-network');
    if (code === 3) return note && note.getAttribute('data-decode');
    if (code === 4) return note && note.getAttribute('data-missing');
    return note && note.getAttribute('data-failed');
  }

  audio.addEventListener('error', function () { say(mediaError()); });
  audio.addEventListener('stalled', function () { say(note && note.getAttribute('data-slow')); });
  audio.addEventListener('playing', function () { say(''); });

  // The recording no longer matches the text: it still plays, and saying so is
  // better than letting it quietly read something the page does not say.
  if (audio.getAttribute('data-stale') === '1') {
    say(note && note.getAttribute('data-stale-note'));
  }

  // ---- following along ----------------------------------------------------

  var shown = -1;

  function clearMark() {
    prose.querySelectorAll('.listen-now').forEach(function (n) {
      n.classList.remove('listen-now');
    });
  }

  function follow() {
    if (!cues.length) return;
    var t = audio.currentTime, i = -1;
    for (var k = 0; k < cues.length; k++) {
      if (t >= cues[k].a && t < cues[k].b) { i = cues[k].i; break; }
    }
    if (i < 0 || i === shown) return;
    shown = i;
    clearMark();
    var el = els[i];
    if (!el) return;
    el.classList.add('listen-now');
    // Only scroll when the line has left the comfortable middle of the screen.
    // Scrolling at every paragraph fights a reader who is following with their
    // own eyes, and the highlight exists so that they do not have to.
    var r = el.getBoundingClientRect();
    var h = window.innerHeight || 0;
    if (r.top < h * 0.15 || r.bottom > h * 0.85) {
      el.scrollIntoView({ behavior: 'smooth', block: 'center' });
    }
  }

  // ---- how far they listened ---------------------------------------------

  var slug = prose.getAttribute('data-read-progress') || '';
  var memo = 'shanraq-listen:' + slug;
  var heard = {};
  try { heard = JSON.parse(sessionStorage.getItem(memo) || '{}'); } catch (e) {}

  function progress(pct) {
    if (!slug || !navigator.sendBeacon) return;
    [25, 50, 75, 100].forEach(function (m) {
      if (pct < m || heard[m]) return;
      heard[m] = true;
      try { sessionStorage.setItem(memo, JSON.stringify(heard)); } catch (e) {}
      try {
        navigator.sendBeacon('/read/' + encodeURIComponent(slug) + '/progress?d=' + m + '&m=listen');
      } catch (e) {}
    });
  }

  // ---- controls -----------------------------------------------------------

  audio.addEventListener('timeupdate', function () {
    follow();
    if (audio.duration) progress((audio.currentTime / audio.duration) * 100);
  });

  audio.addEventListener('play', function () {
    if (label) label.textContent = pauseTxt;
    if (stopBtn) stopBtn.hidden = false;
    if (playBtn) playBtn.setAttribute('aria-pressed', 'true');
  });

  audio.addEventListener('pause', function () {
    if (label) label.textContent = playTxt;
    if (playBtn) playBtn.setAttribute('aria-pressed', 'false');
  });

  audio.addEventListener('ended', function () {
    if (label) label.textContent = playTxt;
    if (stopBtn) stopBtn.hidden = true;
    clearMark();
    shown = -1;
    progress(100);
  });

  if (playBtn) {
    playBtn.addEventListener('click', function () {
      if (!audio.paused) { audio.pause(); return; }
      var p = audio.play();
      // A rejected play() is the one failure with no media error behind it:
      // the browser refused, usually because it does not count this as a real
      // gesture. Nothing in the audio events will fire, so it is reported here
      // or not at all.
      if (p && typeof p.catch === 'function') {
        p.catch(function () { say(note && note.getAttribute('data-blocked')); });
      }
    });
  }

  if (stopBtn) {
    stopBtn.addEventListener('click', function () {
      audio.pause();
      try { audio.currentTime = 0; } catch (e) {}
      stopBtn.hidden = true;
      clearMark();
      shown = -1;
    });
  }
})();
