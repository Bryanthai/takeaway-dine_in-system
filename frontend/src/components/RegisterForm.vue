<template>
  <div class="min-h-screen bg-gradient-to-br from-green-600 to-teal-700 flex items-center justify-center p-4">
    <div class="bg-white rounded-2xl shadow-2xl p-8 md:p-10 w-full max-w-md">
      <h2 class="text-4xl font-extrabold text-gray-900 mb-8 text-center">
        Register
      </h2>

      <form @submit.prevent="handleRegister" class="space-y-6">
        <div>
          <label for="username" class="block text-sm font-semibold text-gray-700 mb-1">Username</label>
          <input
            type="text"
            id="username"
            v-model="username"
            class="mt-1 block w-full px-4 py-2 border border-gray-300 rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-green-500 focus:border-green-500 transition duration-200"
            :class="{ 'border-red-500': validationErrors.username }"
            placeholder="Choose a username"
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
            class="mt-1 block w-full px-4 py-2 border border-gray-300 rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-green-500 focus:border-green-500 transition duration-200"
            :class="{ 'border-red-500': validationErrors.password }"
            placeholder="Choose a strong password"
            required
          />
          <p v-if="validationErrors.password" class="mt-1 text-sm text-red-600">{{ validationErrors.password }}</p>
        </div>

        <div>
          <label for="confirmPassword" class="block text-sm font-semibold text-gray-700 mb-1">Confirm Password</label>
          <input
            type="password"
            id="confirmPassword"
            v-model="confirmPassword"
            class="mt-1 block w-full px-4 py-2 border border-gray-300 rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-green-500 focus:border-green-500 transition duration-200"
            :class="{ 'border-red-500': validationErrors.confirmPassword }"
            placeholder="Confirm your password"
            required
          />
          <p v-if="validationErrors.confirmPassword" class="mt-1 text-sm text-red-600">{{ validationErrors.confirmPassword }}</p>
        </div>

        <button
          type="submit"
          :disabled="loading"
          class="w-full bg-gradient-to-r from-green-600 to-teal-700 text-white font-bold py-3 px-4 rounded-lg shadow-md hover:from-green-700 hover:to-teal-800 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-green-500 transition duration-300 ease-in-out transform hover:-translate-y-0.5"
          :class="{ 'opacity-50 cursor-not-allowed': loading }"
        >
          <span v-if="loading" class="flex items-center justify-center">
            <svg class="animate-spin -ml-1 mr-3 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            Registering...
          </span>
          <span v-else>Register</span>
        </button>

        <div v-if="responseMessage" :class="{
          'mt-4 p-3 rounded-lg text-center font-medium': true,
          'bg-red-100 text-red-700 border border-red-300': !registerSuccess,
          'bg-green-100 text-green-700 border border-green-300': registerSuccess
        }">
          {{ responseMessage }}
        </div>

        <p class="mt-6 text-center text-gray-700">
          Already have an account?
          <router-link to="/login" class="font-semibold text-green-600 hover:text-green-800 transition duration-200">
            Login here
          </router-link>
        </p>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue';
import { useRouter } from 'vue-router'; // Import useRouter

const router = useRouter(); // Initialize router

const username = ref('');
const password = ref('');
const confirmPassword = ref('');
const loading = ref(false);
const responseMessage = ref('');
const registerSuccess = ref(false);
const validationErrors = reactive({
  username: '',
  password: '',
  confirmPassword: ''
});

// Function to validate form inputs
const validateForm = () => {
  let valid = true;
  validationErrors.username = '';
  validationErrors.password = '';
  validationErrors.confirmPassword = '';

  if (!username.value.trim()) {
    validationErrors.username = 'Username is required.';
    valid = false;
  }
  if (!password.value.trim()) {
    validationErrors.password = 'Password is required.';
    valid = false;
  } else if (password.value.length < 6) {
    validationErrors.password = 'Password must be at least 6 characters long.';
    valid = false;
  }
  if (password.value !== confirmPassword.value) {
    validationErrors.confirmPassword = 'Passwords do not match.';
    valid = false;
  }
  return valid;
};

const handleRegister = async () => {
  responseMessage.value = '';
  registerSuccess.value = false;

  if (!validateForm()) {
    responseMessage.value = 'Please correct the errors in the form.';
    return;
  }

  loading.value = true;

  try {
    const response = await fetch('http://localhost:8080/users/register', {
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
      registerSuccess.value = true;
      responseMessage.value = data.message || 'Registration successful! You can now log in.';
      // Optionally clear form
      username.value = '';
      password.value = '';
      confirmPassword.value = '';

      // Redirect to login page after successful registration
      setTimeout(() => {
        router.push('/login'); // Navigate to the /login route
      }, 1500); // Give user a moment to see success message
    } else {
      registerSuccess.value = false;
      responseMessage.value = data.message || 'Registration failed. Please try again.';
    }
  } catch (error) {
    registerSuccess.value = false;
    responseMessage.value = `An error occurred: ${error.message}. Please check your network.`;
    console.error('Registration error:', error);
  } finally {
    loading.value = false;
  }
};
</script>

<style scoped>
/* Scoped styles specific to this component */
</style>