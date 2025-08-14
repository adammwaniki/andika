import { getAssetFromKV } from '@cloudflare/kv-asset-handler';

addEventListener('fetch', event => {
  event.respondWith(handleEvent(event));
});

async function handleEvent(event) {
  try {
    // Try to serve static assets from KV
    return await getAssetFromKV(event);
  } catch (e) {
    // If asset not found, return a simple 404 response
    return new Response('404 - Not Found', { status: 404 });
  }
}