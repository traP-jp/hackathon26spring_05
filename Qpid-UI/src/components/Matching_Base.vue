<template>
  <Matching v-if="isPC" />
  <Matching_Mobile v-else />
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue';
import Matching from './Matching.vue';
import Matching_Mobile from './Matching_Mobile.vue';


// 768 から 576 に変更（576px以上ならすべてPC・タブレット扱いにする）
const isPC = ref(window.innerWidth >= 576);

const handleResize = () => {
  isPC.value = window.innerWidth >= 576;
};

onMounted(() => {
  window.addEventListener('resize', handleResize)
  console.log("matching base loaded...")
});
onUnmounted(() => window.removeEventListener('resize', handleResize));
</script>