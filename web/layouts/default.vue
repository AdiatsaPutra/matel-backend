<template>
  <v-app v-bind:style="{ background: $vuetify.theme.themes.light.background }">
    <div v-if="$nuxt.$route.name === 'login'">
      <Nuxt />
    </div>
    <div v-else>
      <NavigationDrawer />
      <v-main>
        <v-container fluid>
          <Nuxt />
        </v-container>
      </v-main>
      <MainFooter />
    </div>
  </v-app>
</template>

<script>
import NavigationDrawer from "../components/Main/NavigationDrawer.vue";

export default {
  name: "DefaultLayout",
  data() {
    return {
      clipped: false,
      drawer: false,
      items: [
        {
          icon: "mdi-view-dashboard",
          title: "Dashboard",
          to: "/",
        },
        {
          icon: "mdi-motorbike",
          title: "Leasing",
          to: "/leasing",
        },
        {
          icon: "mdi-account-multiple",
          title: "Pengguna",
          to: "/pengguna",
        },
        {
          icon: "mdi-information-outline",
          title: "CMS",
          to: "/cms",
        },
      ],
      miniVariant: false,
      right: true,
      rightDrawer: false,
      title: "Vuetify.js",
    };
  },
  components: { NavigationDrawer },
};
</script>
