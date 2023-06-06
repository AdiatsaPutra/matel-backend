<template>
  <v-app>
    <v-main>
      <v-container fill-height>
        <v-layout align-center justify-center>
          <v-flex xs12 sm8 md4>
            <h1 class="text-center">Login form</h1>
            <div class="mt-6"></div>
            <v-form>
              <v-text-field
                class="rounded-pill"
                placeholder="Masukan Alamat Email"
                name="email"
                type="text"
                v-model="login.email"
                outlined
              ></v-text-field>
              <v-text-field
                class="rounded-pill"
                v-model="login.password"
                :append-icon="show1 ? 'mdi-eye' : 'mdi-eye-off'"
                :type="show1 ? 'text' : 'password'"
                name="input-10-1"
                placeholder="Kata Sandi / Password"
                outlined
                @click:append="show1 = !show1"
              ></v-text-field>
            </v-form>
            <v-spacer></v-spacer>
            <v-btn
              class="rounded-pill"
              height="50"
              :disabled="loading"
              :loading="loading"
              block
              color="primary"
              @click="userLogin"
              >Login</v-btn
            >
            <div class="mt-6"></div>
            <v-alert v-show="isError" type="error">
              <div>{{ errorMessage }}</div>
            </v-alert>
            <div class="mt-6"></div>
            <v-row class="align-center justify-center">
              <p>Belum memiliki akun?&nbsp;</p>
              <p class="purple--text text--darken-2">Daftar</p>
            </v-row>
          </v-flex>
        </v-layout>
      </v-container>
    </v-main>
  </v-app>
</template>

<script>
export default {
  name: "Login",
  data() {
    return {
      show1: false,
      show2: true,
      isError: false,
      errorMessage: "",
      loading: false,
      rules: {
        min: (v) => v.length >= 8 || "Min 8 characters",
        emailMatch: () => `The email and password you entered don't match`,
      },
      login: {
        email: "",
        password: "",
      },
    };
  },
  methods: {
    async userLogin() {
      this.isError = false;
      this.loading = true;
      try {
        let response = await this.$auth.loginWith("local", {
          data: this.login,
        });
        this.loading = false;
        if (this.$auth.loggedIn) {
          if (this.$auth.user.data.is_admin === 1) {
            this.$router.push("/");
          } else {
            this.loading = false;
            this.isError = true;
            this.errorMessage = "Anda bukan admin";
          }
        }
      } catch (error) {
        console.log(error)
        this.loading = false;
        this.isError = true;
        this.errorMessage = error;
      }
    },
  },
};
</script>
