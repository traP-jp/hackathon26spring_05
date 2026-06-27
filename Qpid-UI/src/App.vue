<script setup lang="ts">
import Matching from './components/Matching.vue'
import Header from './components/Header.vue';
import Myself from './components/template/Myself.vue';
</script>

<template>
  <div class="app-layout">
    <Header/>
    <main>
      <router-view />
    </main>
  </div>
</template>

<style scoped>
/* header {
  line-height: 1.5;
}

.logo {
  display: block;
  margin: 0 auto 2rem;
} */

/* @media (min-width: 1024px) {
  header {
    display: flex;
    place-items: center;
    padding-right: calc(var(--section-gap) / 2);
  }

  .logo {
    margin: 0 2rem 0 0;
  }

  header .wrapper {
    display: flex;
    place-items: flex-start;
    flex-wrap: wrap;
  }
} */

body {
  margin: 0;
  font-family: 'Helvetica Neue', Arial, sans-serif;
  background-color: #f8f9fa;
}

.app-layout {
  min-height: 100vh;
}

/* ヘッダーの高さ分だけ上部に余白を作る */
.main-content {
  padding-top: 70px; 
  max-width: 1200px;
  margin: 0 auto;
  padding-left: 20px;
  padding-right: 20px;
}

</style>
