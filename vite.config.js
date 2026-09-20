import { defineConfig } from 'vite';
import vue from '@vitejs/plugin-vue';

const backendUrl = process.env.BACKEND_URL || 'http://localhost:8080';

export default defineConfig({
  plugins: [vue()],
  server: {
    port: 5173,
    proxy: {
      '/ws': {
        target: backendUrl.replace(/^http/, 'ws'),
        ws: true
      },
      '/api': {
        target: backendUrl
      }
    }
  }
});
