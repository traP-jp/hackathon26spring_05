export const apiFetch = (input: RequestInfo | URL, init: RequestInit = {}) => {
  const headers = new Headers(init.headers)

  if (!headers.has('cache-control')) {
    headers.set('cache-control', 'no-store')
  }
  if (!headers.has('pragma')) {
    headers.set('pragma', 'no-cache')
  }

  return fetch(input, {
    ...init,
    cache: init.cache ?? 'no-store',
    headers,
  })
}
