/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly VITE_WS_HOST: string
  readonly VITE_WS_PORT: string
  readonly VITE_WS_PROTOCOL: string
  readonly VITE_API_BASE_URL: string
  readonly VITE_APP_NAME: string
  readonly VITE_APP_VERSION: string
  readonly VITE_DEBUG: string
}

interface ImportMeta {
  readonly env: ImportMetaEnv
}
