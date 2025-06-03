<template>
  <div class="min-h-screen bg-gray-100 flex items-center justify-center p-6">
    <div class="bg-white p-8 rounded-lg shadow-xl w-full max-w-md">
      <h2 class="text-3xl font-bold text-gray-800 mb-6 text-center">Recharge Your Wallet</h2>

      <div v-if="loading" class="flex flex-col items-center justify-center h-48">
        <svg class="animate-spin h-10 w-10 text-indigo-500" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
        </svg>
        <p class="ml-3 text-lg text-gray-700 mt-4">Loading balance...</p>
      </div>

      <div v-else-if="fetchError" class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative" role="alert">
        <strong class="font-bold">Error!</strong>
        <span class="block sm:inline ml-2">{{ fetchError }}</span>
        <button @click="router.push('/login')" class="underline ml-2">Please log in</button>
      </div>

      <div v-else>
        <div class="mb-6 text-center text-gray-800 text-lg">
          Your Current Wallet Balance:
          <span class="font-bold text-green-600 text-2xl">${{ currentBalance.toFixed(2) }}</span>
        </div>

        <form @submit.prevent="rechargeWallet">
          <div class="mb-4">
            <label for="amount" class="block text-gray-700 text-sm font-bold mb-2">Amount to Recharge ($):</label>
            <input
              type="number"
              id="amount"
              v-model.number="rechargeAmount"
              min="1"
              step="0.01"
              class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline focus:ring-2 focus:ring-indigo-500"
              placeholder="e.g., 50.00"
              required
            />
          </div>

          <div class="flex items-center justify-between">
            <button
              type="submit"
              :disabled="rechargeLoading || rechargeAmount <= 0"
              class="bg-blue-600 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded-lg focus:outline-none focus:shadow-outline transition duration-200 disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center"
            >
              <svg v-if="rechargeLoading" class="animate-spin -ml-1 mr-3 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              {{ rechargeLoading ? 'Processing...' : 'Recharge Wallet' }}
            </button>
          </div>
          <div v-if="rechargeError" class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative mt-4" role="alert">
            <strong class="font-bold">Error!</strong>
            <span class="block sm:inline ml-2">{{ rechargeError }}</span>
          </div>
          <div v-if="rechargeSuccess" class="bg-green-100 border border-green-400 text-green-700 px-4 py-3 rounded relative mt-4" role="alert">
            <strong class="font-bold">Success!</strong>
            <span class="block sm:inline ml-2">{{ rechargeSuccess }}</span>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { useUserStore } from '../stores/user'; // Import the user store

const router = useRouter();
const userStore = useUserStore(); // Initialize the user store

const currentBalance = ref(0);
const rechargeAmount = ref(0);

const loading = ref(true);
const fetchError = ref(null);
const rechargeLoading = ref(false);
const rechargeError = ref(null);
const rechargeSuccess = ref(null);

const fetchUserBalance = async () => {
  loading.value = true;
  fetchError.value = null;

  // Use the token from the user store
  if (!userStore.userToken) {
    fetchError.value = 'Authentication token missing. Please log in.';
    loading.value = false;
    router.push('/login');
    return;
  }

  try {
    const response = await fetch('http://localhost:8080/users', {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${userStore.userToken}`, // Use token from store
      },
    });

    // Check for network errors or non-OK status codes first
    if (!response.ok) {
      const errorData = await response.json(); // Attempt to parse error message from backend
      if (response.status === 401 || response.status === 403) {
        fetchError.value = errorData.message || 'Authentication failed. Please log in again.';
        userStore.logout(); // Use user store's logout action
        router.push('/login');
        return; // Exit function after redirect
      }
      // For other non-OK statuses, display a generic error
      throw new Error(errorData.message || `HTTP error! Status: ${response.status}`);
    }

    // If response is OK, parse the data directly as the Account struct
    const userData = await response.json();
    currentBalance.value = userData.Balance; // Directly access Balance from the Account struct

  } catch (err) {
    fetchError.value = `Error fetching user balance: ${err.message}. Please try again.`;
    console.error('Error fetching user balance:', err);
  } finally {
    loading.value = false;
  }
};

const rechargeWallet = async () => {
  rechargeLoading.value = true;
  rechargeError.value = null;
  rechargeSuccess.value = null;

  // Use the token from the user store
  if (!userStore.userToken) {
    rechargeError.value = 'Authentication token missing. Please log in.';
    rechargeLoading.value = false;
    router.push('/login');
    return;
  }

  if (rechargeAmount.value <= 0) {
    rechargeError.value = 'Please enter a valid amount to recharge.';
    rechargeLoading.value = false;
    return;
  }

  try {
    const response = await fetch('http://localhost:8080/recharge', {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${userStore.userToken}`, // Use token from store
      },
      body: JSON.stringify({
        amount: rechargeAmount.value,
      }),
    });

    const data = await response.json(); // Assuming recharge endpoint still returns {success: bool, message: string}

    if (response.ok && data.success) {
      rechargeSuccess.value = `Wallet recharged by $${rechargeAmount.value.toFixed(2)} successfully!`;
      currentBalance.value += rechargeAmount.value; // Optimistically update balance
      rechargeAmount.value = 0; // Reset input field
    } else {
      rechargeError.value = data.message || 'Failed to recharge wallet.';
      if (response.status === 401 || response.status === 403) {
        userStore.logout(); // Use user store's logout action
        router.push('/login');
      }
    }
  } catch (err) {
    rechargeError.value = `Network error: ${err.message}. Please try again.`;
    console.error('Error recharging wallet:', err);
  } finally {
    rechargeLoading.value = false;
  }
};

onMounted(() => {
  fetchUserBalance();
});
</script>

<style scoped>
/* Add your scoped styles here */
</style>