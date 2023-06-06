import colors from "vuetify/es5/util/colors";

export default {
  head: {
    titleTemplate: "%s - motor",
    title: "motor",
    htmlAttrs: {
      lang: "en",
    },
    meta: [
      { charset: "utf-8" },
      { name: "viewport", content: "width=device-width, initial-scale=1" },
      { hid: "description", name: "description", content: "" },
      { name: "format-detection", content: "telephone=no" },
    ],
    link: [{ rel: "icon", type: "image/x-icon", href: "/favicon.ico" }],
  },

  css: [],

  plugins: [],

  components: true,

  buildModules: ["@nuxtjs/vuetify"],

  vuetify: {
    customVariables: ["~/assets/variables.scss"],
    theme: {
      dark: false,
      themes: {
        light: {
          primary: colors.purple.darken4,
          background: "#E9E4ED",
          accent: colors.grey.darken3,
          secondary: colors.amber.darken3,
          info: colors.teal.lighten1,
          warning: colors.amber.base,
          error: colors.deepOrange.accent4,
          success: colors.green.accent3,
        },
      },
    },
  },

  axios: {
    baseURL: "http://localhost:8080/",
    // baseURL: "http://167.172.69.241:8080/",
  },

  build: { transpile: ["defu"] },

  modules: ["@nuxtjs/axios", "@nuxtjs/auth-next"],

  router: {
    middleware: ['auth']
  },

  auth: {
    strategies: {
      local: {
        token: {
          property: "token",
          global: true,
          required: true,
          type: 'Bearer'
        },
        user: {
          property: false,
          autoFetch: true
        },
        endpoints: {
          login: { url: "login-web", method: "post", propertyName: "token" },
          logout: false,
          user: { url: "profil", method: "get" },
        },
      },
    },
    redirect: {
      login: "/login",
      logout: "/login",
      home: "/",
    },
  },
};
