<template>
  <div class="bg-gray-100 p-6 min-h-screen">
    <h2 class="text-4xl font-extrabold text-gray-900 mb-8 text-center">
      All Food Items
    </h2>

    <div v-if="loading" class="flex justify-center items-center h-48">
      <svg class="animate-spin h-10 w-10 text-indigo-500" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
        <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
      </svg>
      <p class="ml-3 text-lg text-gray-700">Loading food items...</p>
    </div>

    <div v-else-if="error" class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative" role="alert">
      <strong class="font-bold">Error!</strong>
      <span class="block sm:inline ml-2">{{ error }}</span>
    </div>

    <div v-else-if="foods.length === 0" class="bg-yellow-100 border border-yellow-400 text-yellow-700 px-4 py-3 rounded relative" role="alert">
      <strong class="font-bold">No Food Found!</strong>
      <span class="block sm:inline ml-2">There are no food items to display.</span>
    </div>

    <div v-else class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6">
      <div
        v-for="food in foods"
        :key="food.FoodID"
        class="bg-white rounded-xl shadow-lg overflow-hidden transform transition-transform duration-300 hover:scale-105 hover:shadow-xl"
      >
        <div class="w-full h-48 flex items-center justify-center bg-gray-200">
          <img
            v-if="food.Picture && food.Picture.String"
            :src="food.Picture.String"
            :alt="food.FoodName"
            class="w-full h-full object-cover"
          />
          <div v-else class="w-full h-full bg-gray-400 flex items-center justify-center text-gray-600 text-xl font-semibold">
            No Image
          </div>
        </div>

        <div class="p-4">
          <h3 class="text-xl font-bold text-gray-900 mb-2 truncate">{{ food.FoodName }}</h3>
          <p class="text-gray-600 text-sm mb-2">Tag: <span class="font-semibold">{{ food.FoodTag }}</span></p>
          <div class="flex justify-between items-center">
            <span class="text-2xl font-bold text-indigo-600">${{ food.Price.toFixed(2) }}</span>
            <button class="bg-indigo-500 text-white px-4 py-2 rounded-lg hover:bg-indigo-600 transition duration-200">
              Details
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';

// Reactive state for food items, loading, and errors
const foods = ref([]);
const loading = ref(true);
const error = ref(null);

// Function to fetch all food items
const fetchAllFood = async () => {
  loading.value = true;
  error.value = null;

  try {
    const response = await fetch('http://localhost:8080/menu', {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
        // 'Authorization': `Bearer YOUR_AUTH_TOKEN_HERE`, // If this endpoint requires authentication
      },
    });

    const data = await response.json();

    if (response.ok && data.success) {
      // Assuming 'foods' in the response directly matches the Food struct array
      foods.value = data.foods;
    } else {
      error.value = data.message || 'Failed to retrieve food items.';
    }
  } catch (err) {
    error.value = `An error occurred: ${err.message}. Please check your network or API endpoint.`;
    console.error('Error fetching food:', err);
  } finally {
    loading.value = false;
  }
};

// Fetch food items when the component is mounted
onMounted(() => {
  fetchAllFood();
});
</script>

<style scoped>
/* Scoped styles specific to this component */
</style>