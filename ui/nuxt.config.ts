import tailwindcss from "@tailwindcss/vite";

export default defineNuxtConfig({
  compatibilityDate: "2025-07-15",
  ssr: false,
  css: ["~/assets/main.css"],
  vite: {
    plugins: [tailwindcss()],
  },
  devtools: { enabled: true },
  imports: {
    dirs: ["./types"],
  },

  build: { analyze: { analyzerMode: "static" } },
  nitro: {
    preset: "static", // bundle as standard vuejs app
    devProxy: {
      "/api": {
        target: "http://localhost:8000/api",
        changeOrigin: true,
        prependPath: true,
      },
    },
  },
});
