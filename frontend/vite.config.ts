import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";
import wails from "@wailsio/runtime/plugins/vite";

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue(), wails("./bindings")],
  base: './',  // 使用相对路径，适配嵌入式文件系统
  build: {
    rollupOptions: {
      external: ["@wailsio/runtime"]
    }
  }
});
