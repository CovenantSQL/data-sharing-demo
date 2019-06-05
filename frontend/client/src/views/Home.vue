<template>
  <div class="home">
    <Crud/>
  </div>
</template>

<script>
  // @ is an alias to /src
  import Crud from '@/components/Crud.vue'
  export default {
    name: 'home',
    components: {
      Crud
    }
  }
</script>
