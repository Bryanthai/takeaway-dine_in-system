<template>
  <div class="min-h-screen bg-gradient-to-br from-purple-500 to-pink-600 flex items-center justify-center p-4">
    <div class="bg-white rounded-2xl shadow-2xl p-8 md:p-10 w-full max-w-md">
      <h2 class="text-4xl font-extrabold text-gray-900 mb-8 text-center">
        Change User Information
      </h2>

      <form @submit.prevent="handleChangeInfo" class="space-y-6">
        <div>
          <label for="username" class="block text-sm font-semibold text-gray-700 mb-1">Username</label>
          <input
            type="text"
            id="username"
            v-model="username"
            class="mt-1 block w-full px-4 py-2 border border-gray-300 rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-purple-500 transition duration-200"
            :class="{ 'border-red-500': validationErrors.username }"
            placeholder="Enter new username"
          />
          <p v-if="validationErrors.username" class="mt-1 text-sm text-red-600">{{ validationErrors.username }}</p>
        </div>

        <div>
          <label for="email" class="block text-sm font-semibold text-gray-700 mb-1">Email</label>
          <input
            type="email"
            id="email"
            v-model="email"
            class="mt-1 block w-full px-4 py-2 border border-gray-300 rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-purple-500 transition duration-200"
            :class="{ 'border-red-500': validationErrors.email }"
            placeholder="Enter new email"
          />
          <p v-if="validationErrors.email" class="mt-1 text-sm text-red-600">{{ validationErrors.email }}</p>
        </div>

        <div>
          <label for="address" class="block text-sm font-semibold text-gray-700 mb-1">Address</label>
          <input
            type="text"
            id="address"
            v-model="address"
            class="mt-1 block w-full px-4 py-2 border border-gray-300 rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-purple-500 transition duration-200"
            :class="{ 'border-red-500': validationErrors.address }"
            placeholder="Enter new address"
          />
          <p v-if="validationErrors.address" class="mt-1 text-sm text-red-600">{{ validationErrors.address }}</p>
        </div>

        <div>
          <label for="phone" class="block text-sm font-semibold text-gray-700 mb-1">Phone</label>
          <input
            type="tel"
            id="phone"
            v-model="phone"
            class="mt-1 block w-full px-4 py-2 border border-gray-300 rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-purple-500 transition duration-200"
            :class="{ 'border-red-500': validationErrors.phone }"
            placeholder="Enter new phone number"
          />
          <p v-if="validationErrors.phone" class="mt-1 text-sm text-red-600">{{ validationErrors.phone }}</p>
        </div>

        <div>
          <label for="newPassword" class="block text-sm font-semibold text-gray-700 mb-1">New Password (Optional)</label>
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
          <label for="confirmNewPassword" class="block text-sm font-semibold text-gray-700 mb-1">Confirm New Password (Optional)</label>
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
import { useRouter } from 'vue-router';

const router = useRouter();

// State variables for the form fields
const username = ref('');
const email = ref('');
const address = ref('');
const phone = ref('');
const newPassword = ref('');
const confirmNewPassword = ref('');

const loading = ref(false);
const responseMessage = ref('');
const updateSuccess = ref(false);

// Reactive object for validation errors
const validationErrors = reactive({
  username: '',
  email: '',
  address: '',
  phone: '',
  newPassword: '',
  confirmNewPassword: ''
});

// Retrieve token from localStorage
const userToken = ref(localStorage.getItem('userToken') || '');

// Form validation logic
const validateForm = () => {
  let valid = true;
  // Clear all previous errors
  Object.keys(validationErrors).forEach(key => validationErrors[key] = '');

  // Validate Username (if provided)
  if (username.value.trim() && username.value.trim().length < 3) {
    validationErrors.username = 'Username must be at least 3 characters long.';
    valid = false;
  }

  // Validate Email (if provided)
  if (email.value.trim() && !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email.value.trim())) {
    validationErrors.email = 'Please enter a valid email address.';
    valid = false;
  }

  // Validate Address (if provided)
  if (address.value.trim() && address.value.trim().length < 5) {
    validationErrors.address = 'Address must be at least 5 characters long.';
    valid = false;
  }

  // Validate Phone (if provided)
  if (phone.value.trim() && !/^\d+$/.test(phone.value.trim())) {
    validationErrors.phone = 'Phone number must contain only digits.';
    valid = false;
  } else if (phone.value.trim() && phone.value.trim().length < 7) {
    validationErrors.phone = 'Phone number must be at least 7 digits long.';
    valid = false;
  }

  // Validate New Password (if provided)
  if (newPassword.value) {
    if (newPassword.value.length < 6) {
      validationErrors.newPassword = 'New password must be at least 6 characters long.';
      valid = false;
    }
    if (newPassword.value !== confirmNewPassword.value) {
      validationErrors.confirmNewPassword = 'Passwords do not match.';
      valid = false;
    }
  } else if (confirmNewPassword.value) {
    // If confirm password is filled but new password isn't
    validationErrors.newPassword = 'Please enter a new password to confirm.';
    valid = false;
  }

  return valid;
};

// Handle form submission
const handleChangeInfo = async () => {
  responseMessage.value = '';
  updateSuccess.value = false;

  if (!validateForm()) {
    responseMessage.value = 'Please correct the errors in the form.';
    return;
  }

  // Essential: Check for token
  if (!userToken.value) {
    responseMessage.value = 'Authentication token is missing. Please log in to change your information.';
    updateSuccess.value = false;
    router.push('/login'); // Redirect to login if token is missing
    return;
  }

  loading.value = true;

  try {
    // Construct the request body dynamically, only including fields that have values
    const requestBody = {};
    if (username.value.trim()) requestBody.username = username.value.trim();
    if (email.value.trim()) requestBody.email = email.value.trim();
    if (address.value.trim()) requestBody.address = address.value.trim();
    if (phone.value.trim()) requestBody.phone = parseInt(phone.value.trim(), 10);
    if (newPassword.value) requestBody.new_password = newPassword.value; // Assuming your backend takes 'new_password'

    // If no fields are provided for update
    if (Object.keys(requestBody).length === 0) {
      responseMessage.value = 'Please fill in at least one field to update.';
      loading.value = false;
      return;
    }

    const response = await fetch('http://localhost:8080/users/change-info', {
      method: 'PUT', // Assuming your backend uses PUT for updates
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
      // Clear all fields on successful update
      username.value = '';
      email.value = '';
      address.value = '';
      phone.value = '';
      newPassword.value = '';
      confirmNewPassword.value = '';
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