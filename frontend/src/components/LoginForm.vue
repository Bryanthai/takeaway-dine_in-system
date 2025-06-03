<template>
  <div class="min-h-screen bg-gradient-to-br from-blue-600 to-indigo-700 flex items-center justify-center p-4">
    <div class="bg-white rounded-2xl shadow-2xl p-8 md:p-10 w-full max-w-md">
      <h2 class="text-4xl font-extrabold text-gray-900 mb-8 text-center">
        Login
      </h2>

      <form @submit.prevent="handleLogin" class="space-y-6">
        <div>
          <label for="username" class="block text-sm font-semibold text-gray-700 mb-1">Username</label>
          <input
            type="text"
            id="username"
            v-model="username"
            class="mt-1 block w-full px-4 py-2 border border-gray-300 rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition duration-200"
            :class="{ 'border-red-500': validationErrors.username }"
            placeholder="Enter your username"
            required
          />
          <p v-if="validationErrors.username" class="mt-1 text-sm text-red-600">{{ validationErrors.username }}</p>
        </div>

        <div>
          <label for="password" class="block text-sm font-semibold text-gray-700 mb-1">Password</label>
          <input
            type="password"
            id="password"
            v-model="password"
            class="mt-1 block w-full px-4 py-2 border border-gray-300 rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition duration-200"
            :class="{ 'border-red-500': validationErrors.password }"
            placeholder="Enter your password"
            required
          />
          <p v-if="validationErrors.password" class="mt-1 text-sm text-red-600">{{ validationErrors.password }}</p>
        </div>

        <button
          type="submit"
          :disabled="loading"
          class="w-full bg-gradient-to-r from-blue-600 to-indigo-700 text-white font-bold py-3 px-4 rounded-lg shadow-md hover:from-blue-700 hover:to-indigo-800 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 transition duration-300 ease-in-out transform hover:-translate-y-0.5"
          :class="{ 'opacity-50 cursor-not-allowed': loading }"
        >
          <span v-if="loading" class="flex items-center justify-center">
            <svg class="animate-spin -ml-1 mr-3 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            Logging in...
          </span>
          <span v-else>Login</span>
        </button>

        <div v-if="responseMessage" :class="{
          'mt-4 p-3 rounded-lg text-center font-medium': true,
          'bg-red-100 text-red-700 border border-red-300': !loginSuccess,
          'bg-green-100 text-green-700 border border-green-300': loginSuccess
        }">
          {{ responseMessage }}
        </div>

        <p class="mt-6 text-center text-gray-700">
          Don't have an account?
          <router-link to="/register" class="font-semibold text-blue-600 hover:text-blue-800 transition duration-200">
            Register here
          </router-link>
        </p>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue';
import { useRouter } from 'vue-router';
import { useUserStore } from '../stores/user'; // Import the user store

const router = useRouter();
const userStore = useUserStore(); // Initialize the user store

const username = ref('');
const password = ref('');
const loading = ref(false);
const responseMessage = ref('');
const loginSuccess = ref(false);
const validationErrors = reactive({
  username: '',
  password: ''
});

// Function to validate form inputs
const validateForm = () => {
  let valid = true;
  validationErrors.username = '';
  validationErrors.password = '';

  if (!username.value.trim()) {
    validationErrors.username = 'Username is required.';
    valid = false;
  }
  if (!password.value.trim()) {
    validationErrors.password = 'Password is required.';
    valid = false;
  }
  return valid;
};

const handleLogin = async () => {
  responseMessage.value = '';
  loginSuccess.value = false;

  if (!validateForm()) {
    responseMessage.value = 'Please fill in all required fields.';
    return;
  }

  loading.value = true;

  try {
    const response = await fetch('http://localhost:8080/users/login', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        username: username.value,
        password: password.value,
      }),
    });

    const data = await response.json();

    if (response.ok && data.success) {
      loginSuccess.value = true;
      responseMessage.value = data.message || 'Login successful!';

      // --- IMPORTANT CHANGE HERE ---
      // Use the setLoginState action from your userStore
      // The backend MUST return 'token' and 'is_admin' in the response
      userStore.setLoginState(data.token, data.is_admin);
      // --- END IMPORTANT CHANGE ---

      console.log('Login successful! User Token:', data.token);
      console.log('Is Admin:', data.is_admin);


      // Redirect after a short delay
      setTimeout(() => {
        router.push('/menu'); // Navigate to the /menu route
      }, 1000); // Give user a moment to see success message
    } else {
      loginSuccess.value = false;
      responseMessage.value = data.message || 'Login failed. Please check your credentials.';
    }
  } catch (error) {
    loginSuccess.value = false;
    responseMessage.value = `An error occurred: ${error.message}. Please check your network.`;
    console.error('Login error:', error);
  } finally {
    loading.value = false;
  }
};
</script>

<style scoped>
/* Scoped styles specific to this component */
</style>