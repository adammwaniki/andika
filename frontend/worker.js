export default {
  async fetch(request, env, ctx) {
    try {
      return await env.ASSETS.fetch(request);
    } catch (e) {
      console.error("Worker error:", e);
      return new Response('404 - Not Found', { status: 404 });
    }
  }
};
