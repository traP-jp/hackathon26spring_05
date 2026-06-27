<template>
  <Matching v-if="isPC" />
  <Matching_Mobile v-else />
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue';
import Matching from './Matching.vue';
import Matching_Mobile from './Matching_Mobile.vue';

const isPC = ref(window.innerWidth >= 768);

const handleResize = () => {
  isPC.value = window.innerWidth >= 768;
};

onMounted(() => window.addEventListener('resize', handleResize));
onUnmounted(() => window.removeEventListener('resize', handleResize));
</script>