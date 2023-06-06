<template>
  <div>
    <v-navigation-drawer v-model="drawer" fixed app>
      <div class="text-h4 py-3 px-6 font-weight-700">Matel</div>
      <v-list>
        <div
          class="text-h6 py-3 px-6 font-weight-700 purple--text text--darken-4"
        >
          DASHBOARD
        </div>
        <v-list-item
          v-for="(item, i) in dashboardItem"
          :key="i"
          :to="item.to"
          router
          exact
        >
          <v-list-item-action>
            <v-icon>{{ item.icon }}</v-icon>
          </v-list-item-action>
          <v-list-item-content>
            <v-list-item-title>{{ item.title }}</v-list-item-title>
          </v-list-item-content>
        </v-list-item>
      </v-list>
      <v-list>
        <div
          class="text-h6 py-3 px-6 font-weight-700 purple--text text--darken-4"
        >
          CONTENT MANAGEMENT
        </div>
        <v-list-item
          v-for="(item, i) in contentManagementItem"
          :key="i"
          :to="item.to"
          router
          exact
        >
          <v-list-item-action>
            <v-icon>{{ item.icon }}</v-icon>
          </v-list-item-action>
          <v-list-item-content>
            <v-list-item-title>{{ item.title }}</v-list-item-title>
          </v-list-item-content>
        </v-list-item>
      </v-list>
    </v-navigation-drawer>
    <v-app-bar fixed app>
      <v-app-bar-nav-icon @click.stop="drawer = !drawer" />
      <v-btn @click="logout">
        Logout
      </v-btn>
    </v-app-bar>
  </div>
</template>

<script>
export default {
  name: "NavigationDrawer",
  data() {
    return {
      drawer: true,
      dashboardItem: [
        {
          icon: "mdi-home-outline",
          title: "Home",
          to: "/",
        },
        {
          icon: "mdi-account-outline",
          title: "Member",
          to: "/member",
        },
        {
          icon: "mdi-file-upload-outline",
          title: "Upload Data",
          to: "/upload-data",
        },
        {
          icon: "mdi-account-multiple-outline",
          title: "Leasing",
          to: "/leasing",
        },

        {
          icon: "mdi-map-marker-outline",
          title: "Wilayah",
          to: "/wilayah",
        },
      ],
      contentManagementItem: [
        {
          icon: "mdi-bank-outline",
          title: "Info Pembayaran",
          to: "/info-pembayaran",
        },
        {
          icon: "mdi-credit-card-outline",
          title: "Nomor Rekening",
          to: "/nomor-rekening",
        },
        {
          icon: "mdi-face-agent",
          title: "Customer Support",
          to: "/customer-support",
        },
      ],
      right: true,
      rightDrawer: false,
    };
  },
  methods:{
    async logout() {
      await this.$auth.logout()
      this.$router.push("/login");
    }
  }
};
</script>
