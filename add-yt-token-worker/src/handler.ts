// Cloudflare secret -> wrangler secret put YOUTUBE_API_TOKEN
declare var YOUTUBE_API_TOKEN: string | undefined;
const API_HOST_NAME = 'youtube.googleapis.com';

// Changes the url to youtube api url and adds the key in the request
export async function handleRequest(request: Request): Promise<Response> {
  if(!YOUTUBE_API_TOKEN) {
    return Response.error();
  }

  const ytURL = new URL(request.url);
  ytURL.hostname = API_HOST_NAME;
  ytURL.searchParams.append('key', YOUTUBE_API_TOKEN);
  const ytRequest = { ...request, url: ytURL.href };

  return await fetch(ytURL.href, ytRequest);
}
