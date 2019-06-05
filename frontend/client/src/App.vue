<template>
  <v-app>
    <div id="nav">
      <router-link to="/">Home</router-link>
      |
      <router-link to="/about">About</router-link>
      <span v-if="isLoggedIn"> | <a @click="logout">Logout {{ user }}</a></span>
    </div>
    <notifications group="login" position="bottom right"/>
    <notifications group="crud" position="bottom right"/>
    <router-view/>
  </v-app>
</template>

<script>
  export default {
    data(){
      return {
        user : "",
      }
    },

    name: 'App',
    computed: {
      isLoggedIn: function () {
        return this.$store.getters.isLoggedIn
      },
    },
    created() {
      this.initialize()
    },

    methods: {
      initialize: function () {
        this.user = localStorage.getItem('user')
      },

      logout: function () {
        this.$store.dispatch('logout')
          .then(() => {
            this.$router.push('/login')
          })
      }
    },
  }
</script>

<style>
  #app {
    font-family: 'Avenir', Helvetica, Arial, sans-serif;
    -webkit-font-smoothing: antialiased;
    -moz-osx-font-smoothing: grayscale;
    text-align: center;
    color: #2c3e50;
  }

  #nav {
    padding: 30px;
  }

  #nav a {
    font-weight: bold;
    color: #2c3e50;
    cursor: pointer;
  }

  #nav a:hover {
    text-decoration: underline;
  }

  #nav a.router-link-exact-active {
    color: #42b983;
  }
</style>
