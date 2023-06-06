<template>
  <v-row class="mb-6 text-center" no-gutters>
    <v-col xs="12" sm="12" md="6" lg="4" xl="2">
      <v-card color="primary pa-6 ma-6">
        <div v-if="loading">
          <h1 class="text-medium white--text">Loading...</h1>
        </div>
        <div v-else>
          <h1 class="text-medium white--text">{{ total }}</h1>
        </div>
        <p class="text-body-1 white--text">Total leasing</p>
      </v-card>
    </v-col>
    <v-col xs="12" sm="12" md="6" lg="4" xl="2">
      <v-card color="primary pa-6 ma-6">
        <h1 class="text-medium white--text">12</h1>
        <p class="text-body-1 white--text">Total Data Kendaraan</p>
      </v-card>
    </v-col>
    <v-col xs="12" sm="12" md="6" lg="4" xl="2">
      <v-card color="primary pa-6 ma-6">
        <h1 class="text-medium white--text">12</h1>
        <p class="text-body-1 white--text">Total Data Kendaraan Per Leasing</p>
      </v-card>
    </v-col>
    <v-col xs="12" sm="12" md="6" lg="4" xl="2">
      <v-card color="primary pa-6 ma-6">
        <h1 class="text-medium white--text">12</h1>
        <p class="text-body-1 white--text">Total Member Trial</p>
      </v-card>
    </v-col>
    <v-col xs="12" sm="12" md="6" lg="4" xl="2">
      <v-card color="primary pa-6 ma-6">
        <h1 class="text-medium white--text">12</h1>
        <p class="text-body-1 white--text">Total member Premium</p>
      </v-card>
    </v-col>
    <v-col xs="12" sm="12" md="6" lg="4" xl="2">
      <v-card color="primary pa-6 ma-6">
        <h1 class="text-medium white--text">12</h1>
        <p class="text-body-1 white--text">Total Payment Confirmation</p>
      </v-card>
    </v-col>
    <v-col xs="12" sm="12" md="6" lg="4" xl="2">
      <v-card color="primary pa-6 ma-6">
        <h1 class="text-medium white--text">12</h1>
        <p class="text-body-1 white--text">Total member expired</p>
      </v-card>
    </v-col>
  </v-row>
</template>

<script>
export default {
  name: "IndexPage",
  middleware: 'auth',
  data() {
    return {
      loading: false,
      total: null,
    };
  },
  mounted() {
    this.fetchData();
  },
  methods: {
    fetchData() {
      this.loading = true;
      this.$axios
        .get("home")
        .then((response) => {
          this.total = response.data.data.leasing;
          this.loading = false;
        })
        .catch((error) => {
          console.error(error);
          this.loading = false;
        })
        .finally(() => {
          this.loading = false;
        });
    },
    
  },
};
</script>
