<template>
  <v-container fluid>
    <v-row>
      <v-col cols="2">
        <v-select
          v-model="limit"
          :items="filterOptions"
          :disabled="loading"
          outlined
          placeholder="Show (Default 100)"
          @change="setLimit(limit)"
        ></v-select>
      </v-col>
      <v-col cols="10">
        <v-text-field
          v-model="search"
          placeholder="Cari berdasarkan leasing, cabang atau nama debitur"
          outlined
          prepend-inner-icon="mdi-magnify"
          @input="debouncedFetchData"
        ></v-text-field>
      </v-col>
    </v-row>

    <v-data-table
      :headers="headers"
      :items="numberedItems"
      :loading="loading"
      :items-per-page="limit"
      :page.sync="currentPage"
      hide-default-footer
      class="elevation-1"
    >
      <template v-slot:item.sisa_hutang="{ item }">
        {{ formatCurrency(item.sisa_hutang) }}
      </template>
    </v-data-table>

    <v-pagination
      v-model="currentPage"
      :length="totalPages"
      @input="fetchData"
      color="primary"
      circle
      class="my-5 custom-pagination"
      :max="10"
    ></v-pagination>
  </v-container>
</template>

<script>
import { debounce } from "lodash";

export default {
  data() {
    return {
      items: [],
      loading: false,
      pagination: {
        page: 1,
        itemsPerPage: 100,
      },
      totalPages: 10,
      search: "",
      currentPage: 1,
      limit: 100,
      filterOptions: [
        { text: "10", value: "10" },
        { text: "25", value: "25" },
        { text: "50", value: "50" },
        { text: "100", value: "100" },
        { text: "200", value: "200" },
      ],
      debouncedFetchData: debounce(this.fetchData, 300),
    };
  },
  computed: {
    headers() {
      return [
        { text: "No", value: "id" },
        { text: "Leasing", value: "leasing" },
        { text: "Cabang", value: "cabang" },
        { text: "Nama Debitur", value: "nama_debitur" },
        { text: "No Polisi", value: "nomor_polisi" },
        { text: "Sisa Hutang", value: "sisa_hutang" },
      ];
    },
    numberedItems() {
      const startIndex = (this.currentPage - 1) * this.limit;
      return this.items.map((item, index) => ({
        ...item,
        id: startIndex + index + 1,
      }));
    },
  },
  mounted() {
    this.fetchData();
  },
  methods: {
    setLimit(limit) {
      this.limit = limit;
      this.fetchData();
    },
    fetchData() {
      this.loading = true;
      this.$axios
        .get("leasing", {
          params: {
            page: this.currentPage,
            limit: this.limit,
            search: this.search,
          },
        })
        .then((response) => {
          console.log(response)
          this.items = response.data.data.leasing;
          this.totalPages = Math.ceil(response.data.data.total / this.limit);
        })
        .catch((error) => {
          console.error(error);
        })
        .finally(() => {
          this.loading = false;
        });
    },
    formatCurrency(value) {
      const formatter = new Intl.NumberFormat("id-ID", {
        style: "currency",
        currency: "IDR",
      });
      return formatter.format(value);
    },
  },
};
</script>
