<template>
  <v-container fluid>
    <v-card width="760">
      <v-responsive max-width="1000">
        <div class="pa-6">
          <!--  -->
          <div class="text-h6 purple--text text--darken-4">UPLOAD DATA</div>
          <div class="py-1"></div>

          <v-alert v-show="success" type="success">
            <div>
              <div class="text-subtitle-1 text--black">
                Data Berhasil di Upload! Dalam {{ this.time }} Detik
              </div>
            </div>
          </v-alert>

          <v-alert v-show="error" type="error">
            <div>
              <div class="text-subtitle-1 text--black">
                {{ error }}
              </div>
            </div>
          </v-alert>

          <!--  -->
          <v-divider
            :thickness="4"
            class="border-opacity-100"
            color="purple"
          ></v-divider>
          <div class="py-1"></div>

          <!--  -->
          <div class="text-body-1">Upload Data Per Leasing</div>
          <div class="py-2"></div>

          <!--  -->
          <v-row class="px-3">
            <v-responsive max-width="500">
              <v-file-input
                v-model="file"
                multiple
                outlined
                dense
                placeholder="Pilih File"
                prepend-icon
                @change="handleFileChange"
              ></v-file-input>
            </v-responsive>
            <!-- <div class="px-2"></div>
            <v-btn color="purple darken-4" outlined height="40">
              Pilih File
            </v-btn> -->
            <div class="px-2"></div>
            <v-btn
              color="purple darken-4 white--text"
              height="40"
              :disabled="isLoading || success || file === null"
              :loading="isLoading"
              @click="uploadFile"
            >
              Simpan
            </v-btn>
          </v-row>
          <div class="py-2"></div>

          <!--  -->
          <div class="text-body-1">Upload Data Per Semua</div>
          <div class="py-2"></div>

          <!--  -->
          <v-row class="px-3">
            <v-responsive max-width="500">
              <v-file-input
                multiple
                outlined
                dense
                placeholder="Pilih File"
                prepend-icon
              ></v-file-input>
            </v-responsive>
            <div class="px-2"></div>
            <v-btn color="purple darken-4 white--text" height="40">
              Simpan
            </v-btn>
          </v-row>
        </div>
      </v-responsive>
    </v-card>
  </v-container>
</template>

<script>
export default {
  data() {
    return {
      isLoading: false,
      success: false,
      error: null,
      file: null,
      formData: null,
      time: null,
    };
  },
  methods: {
    handleFileChange() {
      this.formData = new FormData();
      if (this.file) {
        this.formData.append("file", this.file[0]);
      }
    },
    uploadFile() {
      if (this.formData) {
        this.success = false;
        this.error = false;
        this.isLoading = true;
        this.$axios
          .post("upload-leasing", this.formData, {
            headers: {
              "Content-Type": "multipart/form-data",
            },
          })
          .then((response) => {
            this.formData = null;
            this.formData = null;
            this.success = true;
            this.isLoading = false;
            this.time = response.data.data
            // Handle successful upload
            console.log(response.data);
          })
          .catch((error) => {
            this.isLoading = false;
            this.error = error;
            // Handle upload error
            console.error(error);
          });
      }
    },
  },
};
</script>
