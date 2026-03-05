import { defineConfig, loadEnv } from 'vite'
import react from '@vitejs/plugin-react-swc'
import { resolve } from 'node:path'

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd(), '')
  const apiProxyTarget = env.VITE_DEV_API_PROXY_TARGET || 'http://localhost:8080'

  return {
    plugins: [react()],
    resolve: {
      alias: { '@': resolve(__dirname, './src') },
    },
    server: {
      proxy: {
        '/api': { target: apiProxyTarget, changeOrigin: true },
      },
    },
  }
})
