import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";
import { fileURLToPath, URL } from "node:url";

// @ts-expect-error process is a nodejs global
const host = process.env.TAURI_DEV_HOST;

// https://vite.dev/config/
export default defineConfig(async () => ({
  plugins: [vue()],

  // 路径别名配置
  resolve: {
    alias: {
      "@": fileURLToPath(new URL("./src", import.meta.url))
    }
  },

  // Vite options tailored for Tauri development and only applied in `tauri dev` or `tauri build`
  //
  // 1. prevent Vite from obscuring rust errors
  clearScreen: false,
  // 2. tauri expects a fixed port, fail if that port is not available
  server: {
    port: 1420,
    strictPort: true,
    host: host || false,
    hmr: host
      ? {
          protocol: "ws",
          host,
          port: 1421,
        }
      : undefined,
    watch: {
      // 3. tell Vite to ignore watching `src-tauri`
      ignored: ["**/src-tauri/**"],
    },
  },

  // 构建优化
  build: {
    // 代码分割配置
    rollupOptions: {
      output: {
        manualChunks: {
          // Monaco Editor 相关代码
          'monaco-editor': ['monaco-editor'],
          // 第三方依赖
          'naive-ui': ['naive-ui'],
          'xterm': ['@xterm/xterm', '@xterm/addon-fit', '@xterm/addon-web-links'],
          'icons': ['@tabler/icons-vue', 'lucide-vue-next', '@vicons/carbon', '@vicons/ionicons5'],
          'vue-i18n': ['vue-i18n'],
          'pinia': ['pinia']
        }
      }
    },
    // 增加chunk大小警告限制
    chunkSizeWarningLimit: 1000
  },
}));
