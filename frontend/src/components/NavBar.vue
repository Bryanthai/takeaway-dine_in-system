<template>
  <nav class="bg-indigo-600 p-4 shadow-md sticky top-0 z-50">
    <div class="container mx-auto flex justify-between items-center">
      <router-link to="/" class="text-white text-2xl font-bold hover:text-indigo-100 transition duration-300">
        Food App
      </router-link>

      <div class="md:hidden">
        <button @click="toggleMobileMenu" class="text-white focus:outline-none">
          <svg class="h-8 w-8" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16m-7 6h7"></path>
          </svg>
        </button>
      </div>

      <div class="hidden md:flex space-x-6 items-center">
        <router-link to="/menu" class="text-white hover:text-indigo-100 transition duration-300 text-lg font-medium">Menu</router-link>

        <template v-if="userStore.isAdmin">
          <router-link to="/admin/users" class="text-white hover:text-indigo-100 transition duration-300 text-lg font-medium">Admin Users</router-link>
          <router-link to="/admin/orders" class="text-white hover:text-indigo-100 transition duration-300 text-lg font-medium">Admin Undone Orders</router-link>
          <router-link to="/admin/history" class="text-white hover:text-indigo-100 transition duration-300 text-lg font-medium">Order History</router-link>
          <router-link to="/admin/create-food" class="text-white hover:text-indigo-100 transition duration-300 text-lg font-medium">Create Food</router-link>
        </template>

        <template v-if="userStore.isLoggedIn && !userStore.isAdmin">
          <router-link to="/orders" class="text-white hover:text-indigo-100 transition duration-300 text-lg font-medium">My Orders</router-link>
          <router-link to="/change-info" class="text-white hover:text-indigo-100 transition duration-300 text-lg font-medium">Change Info</router-link>
          <router-link to="/recharge-wallet" class="text-white hover:text-indigo-100 transition duration-300 text-lg font-medium">Recharge Wallet</router-link>
        </template>

        <template v-if="!userStore.isLoggedIn">
          <router-link to="/login" class="bg-indigo-500 text-white px-5 py-2 rounded-full hover:bg-indigo-700 transition duration-300 text-lg font-medium">Login</router-link>
        </template>
        <template v-else>
          <button @click="handleLogout" class="bg-red-500 text-white px-5 py-2 rounded-full hover:bg-red-700 transition duration-300 text-lg font-medium">Logout</button>
        </template>
      </div>
    </div>

    <div v-if="isMobileMenuOpen" class="md:hidden mt-4 space-y-3 px-2 pb-3 pt-2 bg-indigo-700 rounded-md">
      <router-link @click="toggleMobileMenu" to="/menu" class="block text-white hover:bg-indigo-600 rounded-md px-3 py-2 text-base font-medium">Menu</router-link>

      <template v-if="userStore.isAdmin">
        <router-link @click="toggleMobileMenu" to="/admin/users" class="block text-white hover:bg-indigo-600 rounded-md px-3 py-2 text-base font-medium">Admin Users</router-link>
        <router-link @click="toggleMobileMenu" to="/admin/orders" class="block text-white hover:bg-indigo-600 rounded-md px-3 py-2 text-base font-medium">Admin Undone Orders</router-link>
        <router-link @click="toggleMobileMenu" to="/admin/history" class="block text-white hover:bg-indigo-600 rounded-md px-3 py-2 text-base font-medium">Order History</router-link>
        <router-link @click="toggleMobileMenu" to="/admin/create-food" class="block text-white hover:bg-indigo-600 rounded-md px-3 py-2 text-base font-medium">Create Food</router-link>
      </template>

      <template v-if="userStore.isLoggedIn && !userStore.isAdmin">
        <router-link @click="toggleMobileMenu" to="/orders" class="block text-white hover:bg-indigo-600 rounded-md px-3 py-2 text-base font-medium">My Orders</router-link>
        <router-link @click="toggleMobileMenu" to="/change-info" class="block text-white hover:bg-indigo-600 rounded-md px-3 py-2 text-base font-medium">Change Info</router-link>
        <router-link @click="toggleMobileMenu" to="/recharge-wallet" class="block text-white hover:bg-indigo-600 rounded-md px-3 py-2 text-base font-medium">Recharge Wallet</router-link>
      </template>

      <template v-if="!userStore.isLoggedIn">
        <router-link @click="toggleMobileMenu" to="/login" class="block bg-indigo-500 text-white rounded-md px-3 py-2 text-base font-medium text-center">Login</router-link>
      </template>
      <template v-else>
        <button @click="handleLogoutAndCloseMenu" class="w-full bg-red-500 text-white rounded-md px-3 py-2 text-base font-medium text-center">Logout</button>
      </template>
    </div>
  </nav>
</template>

<script setup>
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { useUserStore } from '../stores/user';

const router = useRouter();
const userStore = useUserStore();

const isMobileMenuOpen = ref(false);

const toggleMobileMenu = () => {
  isMobileMenuOpen.value = !isMobileMenuOpen.value;
};

const handleLogout = () => {
  userStore.logout();
  router.push('/login');
};

const handleLogoutAndCloseMenu = () => {
  handleLogout();
  toggleMobileMenu();
};
</script>

<style scoped>
/* No specific scoped styles needed if using Tailwind CSS */
</style>