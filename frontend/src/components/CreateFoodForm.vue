<template>
  <div class="min-h-screen bg-gradient-to-br from-purple-600 to-indigo-700 flex items-center justify-center p-4">
    <div class="bg-white rounded-2xl shadow-2xl p-8 md:p-10 w-full max-w-2xl">
      <h2 class="text-4xl font-extrabold text-gray-900 mb-8 text-center">
        Create New Food Item (Admin)
      </h2>

      <form @submit.prevent="createFood" class="space-y-6">
        <div>
          <label for="foodName" class="block text-sm font-semibold text-gray-700 mb-1">Food Name</label>
          <input
            type="text"
            id="foodName"
            v-model="formData.food_name"
            class="mt-1 block w-full px-4 py-2 border border-gray-300 rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-purple-500 transition duration-200"
            :class="{ 'border-red-500': errors.food_name }"
            placeholder="e.g., Spicy Chicken Wings"
          />
          <p v-if="errors.food_name" class="mt-1 text-sm text-red-600">{{ errors.food_name }}</p>
        </div>

        <div>
          <label for="foodTag" class="block text-sm font-semibold text-gray-700 mb-1">Food Tag</label>
          <input
            type="text"
            id="foodTag"
            v-model="formData.food_tag"
            class="mt-1 block w-full px-4 py-2 border border-gray-300 rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-purple-500 transition duration-200"
            :class="{ 'border-red-500': errors.food_tag }"
            placeholder="e.g., Appetizer, Spicy, Fast Food"
          />
          <p v-if="errors.food_tag" class="mt-1 text-sm text-red-600">{{ errors.food_tag }}</p>
        </div>

        <div>
          <label for="price" class="block text-sm font-semibold text-gray-700 mb-1">Price ($)</label>
          <input
            type="number"
            id="price"
            v-model.number="formData.price"
            step="0.01"
            min="0"
            class="mt-1 block w-full px-4 py-2 border border-gray-300 rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-purple-500 transition duration-200"
            :class="{ 'border-red-500': errors.price }"
            placeholder="e.g., 12.99"
          />
          <p v-if="errors.price" class="mt-1 text-sm text-red-600">{{ errors.price }}</p>
        </div>

        <div>
          <label for="image" class="block text-sm font-semibold text-gray-700 mb-1">Food Image</label>
          <input
            type="file"
            id="image"
            accept="image/*"
            @change="handleImageUpload"
            class="mt-1 block w-full text-sm text-gray-900 border border-gray-300 rounded-lg cursor-pointer bg-gray-50 focus:outline-none focus:border-purple-500 file:mr-4 file:py-2 file:px-4 file:rounded-md file:border-0 file:text-sm file:font-semibold file:bg-purple-50 file:text-purple-700 hover:file:bg-purple-100"
            :class="{ 'border-red-500': errors.image }"
          />
          <p v-if="errors.image" class="mt-1 text-sm text-red-600">{{ errors.image }}</p>
          <p v-if="selectedImageName" class="mt-2 text-sm text-gray-500">Image selected: {{ selectedImageName }}</p>
        </div>

        <div>
          <label for="description" class="block text-sm font-semibold text-gray-700 mb-1">Description</label>
          <textarea
            id="description"
            v-model="formData.description"
            rows="3"
            class="mt-1 block w-full px-4 py-2 border border-gray-300 rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-purple-500 transition duration-200"
            :class="{ 'border-red-500': errors.description }"
            placeholder="A brief description of the food item."
          ></textarea>
          <p v-if="errors.description" class="mt-1 text-sm text-red-600">{{ errors.description }}</p>
        </div>

        <div>
          <label for="ingredients" class="block text-sm font-semibold text-gray-700 mb-1">Ingredients</label>
          <textarea
            id="ingredients"
            v-model="formData.ingredients"
            rows="2"
            class="mt-1 block w-full px-4 py-2 border border-gray-300 rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-purple-500 transition duration-200"
            :class="{ 'border-red-500': errors.ingredients }"
            placeholder="Comma-separated list of ingredients, e.g., Chicken, Flour, Spices"
          ></textarea>
          <p v-if="errors.ingredients" class="mt-1 text-sm text-red-600">{{ errors.ingredients }}</p>
        </div>

        <div>
          <label for="info" class="block text-sm font-semibold text-gray-700 mb-1">Additional Info (Optional)</label>
          <input
            type="text"
            id="info"
            v-model="formData.info"
            class="mt-1 block w-full px-4 py-2 border border-gray-300 rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-purple-500 transition duration-200"
            placeholder="e.g., Gluten-free, Spicy Level: Medium"
          />
        </div>

        <div>
          <label for="time_needed" class="block text-sm font-semibold text-gray-700 mb-1">Preparation Time (minutes)</label>
          <input
            type="number"
            id="time_needed"
            v-model.number="formData.time_needed"
            class="mt-1 block w-full px-4 py-2 border border-gray-300 rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-green-500 focus:border-green-500 transition duration-200"
            :class="{ 'border-red-500': errors.time_needed }"
            placeholder="e.g., 30"
            min="1"
            required
          />
          <p v-if="errors.time_needed" class="mt-1 text-sm text-red-600">{{ errors.time_needed }}</p>
        </div>

        <div class="flex items-center">
          <input
            type="checkbox"
            id="long_range"
            v-model="formData.long_range"
            class="h-4 w-4 text-purple-600 focus:ring-purple-500 border-gray-300 rounded"
          />
          <label for="long_range" class="ml-2 block text-sm font-semibold text-gray-700">Available for Long Range Delivery</label>
        </div>

        <button
          type="submit"
          :disabled="loading"
          class="w-full bg-gradient-to-r from-purple-600 to-indigo-700 text-white font-bold py-3 px-4 rounded-lg shadow-md hover:from-purple-700 hover:to-indigo-800 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-purple-500 transition duration-300 ease-in-out transform hover:-translate-y-0.5"
          :class="{ 'opacity-50 cursor-not-allowed': loading }"
        >
          <span v-if="loading" class="flex items-center justify-center">
            <svg class="animate-spin -ml-1 mr-3 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            Creating Food...
          </span>
          <span v-else>Create Food Item</span>
        </button>

        <div v-if="responseMessage" :class="{
          'mt-4 p-3 rounded-lg text-center font-medium': true,
          'bg-green-100 text-green-700 border border-green-300': createSuccess,
          'bg-red-100 text-red-700 border border-red-300': !createSuccess
        }">
          {{ responseMessage }}
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue';

// --- IMPORTANT: Get token from localStorage instead of dummy ---
const userToken = ref(localStorage.getItem('userToken') || '');

// Form data structure. `image` will now hold the File object directly.
const formData = reactive({
  food_name: '',
  food_tag: '',
  price: null,
  image: null, // Changed to null, will hold the File object
  info: '',
  description: '', // Added description to formData
  ingredients: '',
  time_needed: null, // Changed to null, will hold the number
  long_range: false, // **Added new field for long_range, default to false**
});

const errors = reactive({
  food_name: '',
  food_tag: '',
  price: '',
  image: '',
  description: '',
  ingredients: '',
  time_needed: '',
});

const loading = ref(false);
const responseMessage = ref('');
const createSuccess = ref(false);
const selectedImageName = ref(''); // To display the selected file name

// Handles file input and stores the File object
const handleImageUpload = (event) => {
  const file = event.target.files[0];
  if (!file) {
    formData.image = null;
    errors.image = '';
    selectedImageName.value = '';
    return;
  }

  // Basic file type validation
  if (!file.type.startsWith('image/')) {
    errors.image = 'Please select an image file.';
    formData.image = null;
    selectedImageName.value = '';
    return;
  }

  formData.image = file; // Store the File object directly
  selectedImageName.value = file.name; // Store the file name for display
  errors.image = ''; // Clear any previous image errors
};

const validateForm = () => {
  let valid = true;
  Object.keys(errors).forEach(key => errors[key] = '');

  if (!formData.food_name.trim()) {
    errors.food_name = 'Food Name is required.';
    valid = false;
  }
  if (!formData.food_tag.trim()) {
    errors.food_tag = 'Food Tag is required.';
    valid = false;
  }
  if (formData.price === null || isNaN(formData.price) || formData.price < 0) {
    errors.price = 'Price must be a non-negative number.';
    valid = false;
  } else if (String(formData.price).split('.')[1] && String(formData.price).split('.')[1].length > 2) {
    errors.price = 'Price can have at most two decimal places.';
    valid = false;
  }
  if (!formData.image) { // Validate if an image file has been selected
    errors.image = 'Food Image is required.';
    valid = false;
  }
  if (!formData.description.trim()) {
    errors.description = 'Description is required.';
    valid = false;
  }
  if (!formData.ingredients.trim()) {
    errors.ingredients = 'Ingredients are required.';
    valid = false;
  }
  if (formData.time_needed === null || isNaN(formData.time_needed) || formData.time_needed <= 0) {
    errors.time_needed = 'Preparation Time must be a positive number (in minutes).';
    valid = false;
  }

  return valid;
};

const createFood = async () => {
  responseMessage.value = '';
  createSuccess.value = false;

  if (!validateForm()) {
    responseMessage.value = 'Please correct the errors in the form.';
    createSuccess.value = false;
    return;
  }

  // Ensure we have a token before proceeding
  if (!userToken.value) {
    responseMessage.value = 'Authentication token is missing. Please log in as admin to create food items.';
    createSuccess.value = false;
    return;
  }

  loading.value = true;

  try {
    // Create a FormData object to send both text fields and the file
    const dataToSend = new FormData();
    dataToSend.append('food_name', formData.food_name);
    dataToSend.append('food_tag', formData.food_tag);
    dataToSend.append('price', parseFloat(formData.price));
    dataToSend.append('image', formData.image); // Append the actual File object
    dataToSend.append('info', formData.info);
    dataToSend.append('description', formData.description); // Ensure description is included
    dataToSend.append('ingredients', formData.ingredients);
    dataToSend.append('time_needed', parseInt(formData.time_needed, 10));
    dataToSend.append('long_range', formData.long_range); // **Append the new long_range field**

    const response = await fetch('http://localhost:8080/foods', {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${userToken.value}`,
      },
      body: dataToSend, // Send the FormData object
    });

    const data = await response.json();

    if (response.ok && data.success) {
      createSuccess.value = true;
      responseMessage.value = data.message || `Food item '${formData.food_name}' created successfully!`;
      // Clear form
      Object.assign(formData, {
        food_name: '',
        food_tag: '',
        price: null, // Reset to null
        image: null, // Reset image to null
        info: '',
        description: '', // Reset description
        ingredients: '',
        time_needed: null, // Reset to null
        long_range: false, // **Reset long_range to false**
      });
      const fileInput = document.getElementById('image');
      if (fileInput) fileInput.value = ''; // Clear the file input visually
      selectedImageName.value = ''; // Clear the displayed file name

    } else {
      createSuccess.value = false;
      responseMessage.value = data.message || 'Failed to create food item. Please try again.';
      if (response.status === 401 || response.status === 403) {
        responseMessage.value = response.status === 401 ? "Unauthorized: Please log in with admin privileges." : "Forbidden: You don't have permission to perform this action.";
        // Optionally, clear token and redirect to login if unauthorized/forbidden
        // localStorage.removeItem('userToken');
        // router.push('/login');
      }
    }
  } catch (error) {
    createSuccess.value = false;
    responseMessage.value = `An error occurred: ${error.message}. Please check your network connection or API endpoint.`;
    console.error('Create food error:', error);
  } finally {
    loading.value = false;
  }
};
</script>