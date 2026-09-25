// Follow the section being read. The prose and TOC use the same sec-N IDs.
(function () {
  var nav = document.querySelector('.article-aside .toc-box:not(.course-box)');
  if (!nav) return;
  var list = nav.querySelector('.toc-box__list');
  if (!list) return;

  var sections = Array.from(list.querySelectorAll('a[href^="#"]')).map(function (link) {
    var heading = document.getElementById(link.getAttribute('href').slice(1));
    return heading && heading.closest('.prose') ? { link: link, heading: heading } : null;
  }).filter(Boolean);
  if (!sections.length) return;

  var current = null;
  var pending = false;
  function update() {
    pending = false;
    // A heading becomes current after it crosses the upper reading area.
    var line = Math.min(180, window.innerHeight * 0.3);
    var next = null;
    sections.forEach(function (section) {
      if (section.heading.getBoundingClientRect().top <= line) next = section;
    });
    if (next === current) return;
    if (current) current.link.removeAttribute('aria-current');
    current = next;
    if (!current) return;
    current.link.setAttribute('aria-current', 'location');

    // Move only the TOC list. scrollIntoView would also move the article page.
    var box = list.getBoundingClientRect();
    var row = current.link.getBoundingClientRect();
    if (row.top < box.top) list.scrollTop += row.top - box.top;
    else if (row.bottom > box.bottom) list.scrollTop += row.bottom - box.bottom;
  }
  function schedule() {
    if (pending) return;
    pending = true;
    requestAnimationFrame(update);
  }
  window.addEventListener('scroll', schedule, { passive: true });
  window.addEventListener('resize', schedule, { passive: true });
  update();
})();
