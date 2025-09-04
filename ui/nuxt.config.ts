import tailwindcss from "@tailwindcss/vite"

export default defineNuxtConfig({
  compatibilityDate: "2025-07-15",
  ssr: false,
  css: ["~/assets/css/tailwind.css", "vue-data-ui/style.css"],
  vite: {
    plugins: [tailwindcss()]
  },
  devtools: { enabled: true },
  build: { analyze: { analyzerMode: "static" } },
  nitro: {
    preset: "static", // bundle as standard vuejs app
    devProxy: {
      "/api": {
        target: "http://localhost:8000/api",
        changeOrigin: true,
        prependPath: true
      },
      "/auth": {
        target: "http://localhost:8000/auth",
        cookieDomainRewrite: "localhost"
      }
    }
  },
  modules: ["@nuxtjs/color-mode", "@vueuse/nuxt", "@nuxt/fonts", "@vee-validate/nuxt", "reka-ui/nuxt", "motion-v/nuxt"],
  imports: {
    dirs: ["./tpes"],
    imports: [
      {
        from: "tailwind-variants",
        name: "tv"
      },
      {
        from: "tailwind-variants",
        name: "VariantProps",
        type: true
      }
    ]
  },
  colorMode: {
    storageKey: "ui-color-mode",
    classSuffix: ""
  }
})