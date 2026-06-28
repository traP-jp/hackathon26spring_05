const CHECK_INTERVAL_MS = 60_000
const ASSET_SELECTOR = 'script[src^="/assets/"], link[rel="stylesheet"][href^="/assets/"]'

const getAssetPath = (element: Element) => {
  if (element instanceof HTMLScriptElement) {
    return element.getAttribute('src')
  }
  if (element instanceof HTMLLinkElement) {
    return element.getAttribute('href')
  }
  return null
}

const getAssets = (root: ParentNode) =>
  new Set(
    Array.from(root.querySelectorAll(ASSET_SELECTOR))
      .map(getAssetPath)
      .filter((path): path is string => path !== null),
  )

const areSameAssets = (currentAssets: Set<string>, nextAssets: Set<string>) => {
  if (currentAssets.size !== nextAssets.size) {
    return false
  }
  return Array.from(currentAssets).every((asset) => nextAssets.has(asset))
}

export const startUpdateChecker = () => {
  if (!import.meta.env.PROD) {
    return
  }

  const currentAssets = getAssets(document)
  if (currentAssets.size === 0) {
    return
  }

  const checkForUpdate = async () => {
    try {
      const response = await fetch(`/?__qpid_update_check=${Date.now()}`, {
        cache: 'no-store',
        headers: {
          'cache-control': 'no-store',
          pragma: 'no-cache',
        },
      })
      if (!response.ok) {
        return
      }

      const html = await response.text()
      const nextDocument = new DOMParser().parseFromString(html, 'text/html')
      const nextAssets = getAssets(nextDocument)
      if (nextAssets.size > 0 && !areSameAssets(currentAssets, nextAssets)) {
        window.location.reload()
      }
    } catch (error) {
      console.log('Update check failed:', error)
    }
  }

  document.addEventListener('visibilitychange', () => {
    if (!document.hidden) {
      void checkForUpdate()
    }
  })

  window.setInterval(() => {
    void checkForUpdate()
  }, CHECK_INTERVAL_MS)
}
