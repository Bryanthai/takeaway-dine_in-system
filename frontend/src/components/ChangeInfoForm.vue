<template>
  <div class="min-h-screen bg-gradient-to-br from-purple-500 to-pink-600 flex items-center justify-center p-4">
    <div class="bg-white rounded-2xl shadow-2xl p-8 md:p-10 w-full max-w-md">
      <h2 class="text-4xl font-extrabold text-gray-900 mb-8 text-center">
        Change User Information
      </h2>

      <form @submit.prevent="handleChangeInfo" class="space-y-6">
        <div>
          <label for="newPassword" class="block text-sm font-semibold text-gray-700 mb-1">New Password</label>
          <input
            type="password"
            id="newPassword"
            v-model="newPassword"
            class="mt-1 block w-full px-4 py-2 border border-gray-300 rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-purple-500 transition duration-200"
            :class="{ 'border-red-500': validationErrors.newPassword }"
            placeholder="Enter new password"
          />
          <p v-if="validationErrors.newPassword" class="mt-1 text-sm text-red-600">{{ validationErrors.newPassword }}</p>
        </div>

        <div>
          <label for="confirmNewPassword" class="block text-sm font-semibold text-gray-700 mb-1">Confirm New Password</label>
          <input
            type="password"
            id="confirmNewPassword"
            v-model="confirmNewPassword"
            class="mt-1 block w-full px-4 py-2 border border-gray-300 rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-purple-500 transition duration-200"
            :class="{ 'border-red-500': validationErrors.confirmNewPassword }"
            placeholder="Confirm new password"
          />
          <p v-if="validationErrors.confirmNewPassword" class="mt-1 text-sm text-red-600">{{ validationErrors.confirmNewPassword }}</p>
        </div>

        <button
          type="submit"
          :disabled="loading"
          class="w-full bg-gradient-to-r from-purple-600 to-pink-700 text-white font-bold py-3 px-4 rounded-lg shadow-md hover:from-purple-700 hover:to-pink-800 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-purple-500 transition duration-300 ease-in-out transform hover:-translate-y-0.5"
          :class="{ 'opacity-50 cursor-not-allowed': loading }"
        >
          <span v-if="loading" class="flex items-center justify-center">
            <svg class="animate-spin -ml-1 mr-3 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            Updating Info...
          </span>
          <span v-else>Update Information</span>
        </button>

        <div v-if="responseMessage" :class="{
          'mt-4 p-3 rounded-lg text-center font-medium': true,
          'bg-green-100 text-green-700 border border-green-300': updateSuccess,
          'bg-red-100 text-red-700 border border-red-300': !updateSuccess
        }">
          {{ responseMessage }}
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue';
import { useRouter } from 'vue-router'; // Import useRouter for potential redirect

const router = useRouter();

const newPassword = ref('');
const confirmNewPassword = ref('');
// Add other fields you want to change, e.g.
// const newEmail = ref('');

const loading = ref(false);
const responseMessage = ref('');
const updateSuccess = ref(false);
const validationErrors = reactive({
  newPassword: '',
  confirmNewPassword: '',
  // Add errors for other fields
});

// Retrieve token from localStorage
const userToken = ref(localStorage.getItem('userToken') || '');

const validateForm = () => {
  let valid = true;
  Object.keys(validationErrors).forEach(key => validationErrors[key] = ''); // Clear previous errors

  if (newPassword.value) { // If new password is provided, validate it
    if (newPassword.value.length < 6) {
      validationErrors.newPassword = 'New password must be at least 6 characters long.';
      valid = false;
    }
    if (newPassword.value !== confirmNewPassword.value) {
      validationErrors.confirmNewPassword = 'Passwords do not match.';
      valid = false;
    }
  }

  return valid;
};

const handleChangeInfo = async () => {
  responseMessage.value = '';
  updateSuccess.value = false;

  if (!validateForm()) {
    return;
  }

  // Essential: Check for token
  if (!userToken.value) {
    responseMessage.value = 'Authentication token is missing. Please log in to change your information.';
    updateSuccess.value = false;
    // Optionally redirect to login if token is missing
    // router.push('/login');
    return;
  }

  loading.value = true;

  try {
    const requestBody = {
      // Only include fields that are actually being changed
      // For example, if only password is being changed:
      new_password: newPassword.value,
      // new_email: newEmail.value, // If you add an email field
    };

    const response = await fetch('http://localhost:8080/users/change-info', { // Adjust endpoint as per your backend
      method: 'PUT', // or POST, depending on your backend
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${userToken.value}`, // Send the token!
      },
      body: JSON.stringify(requestBody),
    });

    const data = await response.json();

    if (response.ok && data.success) {
      updateSuccess.value = true;
      responseMessage.value = data.message || 'Information updated successfully!';
      // Clear form fields if desired
      newPassword.value = '';
      confirmNewPassword.value = '';
      // newEmail.value = '';
    } else {
      updateSuccess.value = false;
      responseMessage.value = data.message || 'Failed to update information. Please try again.';
      if (response.status === 401 || response.status === 403) {
         responseMessage.value = "Session expired or unauthorized. Please log in again.";
         localStorage.removeItem('userToken'); // Clear invalid token
         router.push('/login'); // Redirect to login
      }
    }
  } catch (error) {
    updateSuccess.value = false;
    responseMessage.value = `An error occurred: ${error.message}. Please check your network or API endpoint.`;
    console.error('Change info error:', error);
  } finally {
    loading.value = false;
  }
};
</script>

<style scoped>
/* Add your scoped styles here */
</style>