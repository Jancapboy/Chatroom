import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import tailwindcss from '@tailwindcss/vite'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react(), tailwindcss()],
  server: {
    port: 3000,
    host: true,
    proxy: {
      '/api': 'http://localhost:4001',
      '/ws': {
        target: 'ws://localhost:4001',
        ws: true,
      },
    },
  },
})
