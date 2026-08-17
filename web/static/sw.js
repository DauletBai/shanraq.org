// Service worker for the installable site.
//
// It caches ONE thing: /static/*. Those URLs carry a content hash (AssetURL),
// so a cached copy can never be stale — a changed file is a changed URL. Pages
// are deliberately never cached: a classifieds site whose listings come back
// from yesterday's cache is worse than one that simply needs a connection.
const CACHE = 'shanraq-static-v1';

self.addEventListener('install', (e) => {
  // Take over immediately rather than waiting for every tab to close.
  self.skipWaiting();
});

self.addEventListener('activate', (e) => {
  e.waitUntil((async () => {
    const names = await caches.keys();
    await Promise.all(names.filter((n) => n !== CACHE).map((n) => caches.delete(n)));
    await self.clients.claim();
  })());
});

self.addEventListener('fetch', (e) => {
  const req = e.request;
  if (req.method !== 'GET') return;
  const url = new URL(req.url);
  if (url.origin !== self.location.origin) return;
  if (!url.pathname.startsWith('/static/')) return; // pages always go to the network

  e.respondWith((async () => {
    const hit = await caches.match(req);
    if (hit) return hit;
    const res = await fetch(req);
    if (res && res.status === 200) {
      const cache = await caches.open(CACHE);
      cache.put(req, res.clone());
    }
    return res;
  })());
});
