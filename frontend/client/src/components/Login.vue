<template>
  <v-content>
    <v-container fluid fill-height>
      <v-layout align-center justify-center>
        <v-flex xs12 sm8 md4>
          <v-card class="elevation-12">
            <v-toolbar dark color="primary">
              <v-toolbar-title>Login form</v-toolbar-title>
            </v-toolbar>
            <v-card-text>
              <v-form @submit.prevent="login">
                <v-text-field prepend-icon="person" name="user" label="user" v-model="user" type="text"></v-text-field>
                <v-text-field prepend-icon="lock" name="password" label="password" v-model="password" id="password"
                              type="password"></v-text-field>
              </v-form>
            </v-card-text>
            <v-card-actions>
              <v-spacer></v-spacer>
              <v-btn color="primary" @click="login">Login</v-btn>
            </v-card-actions>
          </v-card>
        </v-flex>
      </v-layout>
    </v-container>
  </v-content>
</template>
<script>
  export default {
    data() {
      return {
        user: "",
        password: "",
      }
    },
    methods: {
      login: function () {
        let user = this.user;
        let password = this.password;
        this.$store.dispatch('login', {user, password})
          .then(
            () => this.$router.push({path: '/'})
          )
          .catch(err => {
            this.$notify({
              group: 'login',
              type: 'error',
              title: "Login failed",
              text: err.toString(),
            });
            return this.$log.error(err);
          })
      }
    }
  }
</script>
