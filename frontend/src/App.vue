<template>
  <div id="app">
    <NavBar />

    <main class="pt-16"> <router-view></router-view>
    </main>
  </div>
</template>

<script setup>
import NavBar from './components/NavBar.vue'; // Import your NavBar component
// No other script logic needed here typically, unless for global app setup
</script>

<style>
/* You can keep global styles here if needed */
body {
  margin: 0;
  font-family: sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}
/* Ensure your Tailwind CSS is properly imported in your main.js or index.css */
</style>