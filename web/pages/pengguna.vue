<template>
  <v-card class="ma-4">
    <v-card-title>
      <v-text-field
        v-model="search"
        append-icon="mdi-magnify"
        label="Search"
        single-line
        hide-details
      ></v-text-field>
    </v-card-title>
    <v-data-table
      :headers="headers"
      :items="desserts"
      :search="search"
      :items-per-page="20"
    ></v-data-table>
  </v-card>
</template>

<script>
export default {
  name: "PenggunaPage",
  data() {
    return {
      search: "",
      
    };
  },
};
</script>
