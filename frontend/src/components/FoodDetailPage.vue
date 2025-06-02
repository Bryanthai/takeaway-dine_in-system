<template>
  <div class="min-h-screen bg-gray-100 p-6">
    <h2 class="text-4xl font-extrabold text-gray-900 mb-8 text-center">
      Food Details
    </h2>

    <div v-if="loading" class="flex justify-center items-center h-48">
      <svg class="animate-spin h-10 w-10 text-indigo-500" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
        <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
      </svg>
      <p class="ml-3 text-lg text-gray-700">Loading food details...</p>
    </div>

    <div v-else-if="error" class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative" role="alert">
      <strong class="font-bold">Error!</strong>
      <span class="block sm:inline ml-2">{{ error }}</span>
    </div>

    <div v-else-if="!food" class="bg-yellow-100 border border-yellow-400 text-yellow-700 px-4 py-3 rounded relative text-center" role="alert">
      <strong class="font-bold">Food Item Not Found!</strong>
      <span class="block sm:inline ml-2">The requested food item could not be loaded.</span>
    </div>

    <div v-else class="max-w-4xl mx-auto bg-white rounded-lg shadow-xl p-8 flex flex-col md:flex-row">
      <div class="md:w-1/3 flex items-center justify-center p-4">
        <img
          v-if="food.Picture && food.Picture.String"
          :src="food.Picture.String"
          :alt="food.FoodName"
          class="rounded-lg shadow-md max-h-64 object-cover w-full"
        />
        <div v-else class="w-full h-64 bg-gray-300 rounded-lg flex items-center justify-center text-gray-600 text-2xl font-semibold">
          No Image
        </div>
      </div>

      <div class="md:w-2/3 md:pl-8 mt-6 md:mt-0">
        <h3 class="text-3xl font-bold text-gray-900 mb-2">{{ food.FoodName }}</h3>
        <p class="text-indigo-600 text-2xl font-bold mb-4">${{ food.Price.toFixed(2) }}</p>

        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4 mb-6">
          <p class="text-gray-700"><span class="font-semibold">Category:</span> {{ food.FoodTag }}</p>
          <p class="text-gray-700"><span class="font-semibold">Type:</span> {{ food.FoodType }}</p>
          <p class="text-gray-700"><span class="font-semibold">Preparation Time:</span> {{ food.TimeNeeded }} minutes</p>
          <p class="text-gray-700"><span class="font-semibold">Availability:</span> {{ food.LongRange ? 'Delivery' : 'Pickup Only' }}</p>

          <p v-if="ratingLoading" class="text-gray-600 col-span-2">Loading rating...</p>
          <p v-else-if="ratingError" class="text-red-600 col-span-2">Rating Error: {{ ratingError }}</p>
          <p v-else-if="ratingInfo" class="text-gray-700 col-span-2">
            <span class="font-semibold">Average Rating:</span>
            <span v-if="ratingInfo.Rating > 0">{{ ratingInfo.Rating.toFixed(1) }} / 5 ({{ ratingInfo.Times }} orders rated)</span>
            <span v-else>No ratings yet ({{ ratingInfo.Times }} orders)</span>
          </p>
        </div>

        <div class="mb-6">
          <p class="text-gray-800 font-semibold mb-2">Description:</p>
          <p class="text-gray-700">{{ food.Description || 'No description available.' }}</p>
        </div>

        <div class="mb-6" v-if="food.Ingredients">
          <p class="text-gray-800 font-semibold mb-2">Ingredients:</p>
          <p class="text-gray-700">{{ food.Ingredients }}</p>
        </div>

        <div class="mb-6" v-if="food.Info && food.Info.Valid">
          <p class="text-gray-800 font-semibold mb-2">Additional Info:</p>
          <p class="text-gray-700 italic">{{ food.Info.String }}</p>
        </div>

        <button
          @click="addToCart(food)"
          class="w-full bg-green-600 text-white px-6 py-3 rounded-lg font-semibold text-lg hover:bg-green-700 transition duration-200 flex items-center justify-center mt-4"
        >
          <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 mr-2" viewBox="0 0 24 24" fill="currentColor">
            <path d="M12 4c-.55 0-1 .45-1 1v6H5c-.55 0-1 .45-1 1s.45 1 1 1h6v6c0 .55.45 1 1 1s1-.45 1-1v-6h6c.55 0 1-.45 1-1s-.45-1-1-1h-6V5c0-.55-.45-1-1-1z"/>
          </svg>
          Add to Cart
        </button>
        <p v-if="cartStore.itemAddedSuccessfully" class="text-green-500 text-center mt-2">Item added to cart!</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useCartStore } from '../stores/cart'; // Assuming you have a cart store

const route = useRoute();
const router = useRouter();
const cartStore = useCartStore();

const foodId = ref(null);
const food = ref(null); // Will hold the Food object
const ratingInfo = ref(null); // Will hold { Rating: float64, Times: int64 }

const loading = ref(true); // Loading state for main food details
const error = ref(null);    // Error for main food details
const ratingLoading = ref(true); // Loading state for rating info
const ratingError = ref(null);    // Error for rating info

const userToken = ref(localStorage.getItem('userToken') || ''); // Token might be needed for some info

// Function to fetch food details
const fetchFoodDetails = async () => {
  loading.value = true;
  error.value = null;

  if (!foodId.value) {
    error.value = 'No Food ID provided.';
    loading.value = false;
    return;
  }

  try {
    const response = await fetch(`http://localhost:8080/foods?food_id=${foodId.value}`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
        // 'Authorization': `Bearer ${userToken.value}`, // Might be needed if this endpoint is protected
      },
    });

    const data = await response.json();

    if (response.ok && data.success) {
      food.value = data.data; // Assuming your backend returns the single Food object under 'data'
      if (!food.value) {
        error.value = 'Food item not found.';
      }
    } else {
      error.value = data.message || 'Failed to fetch food details.';
      // Handle 401/403 if this endpoint becomes protected in the future
    }
  } catch (err) {
    error.value = `Error fetching food details: ${err.message}.`;
    console.error('Error fetching food details:', err);
  } finally {
    loading.value = false;
  }
};

// Function to fetch food rating and times ordered
const fetchFoodRating = async () => {
  ratingLoading.value = true;
  ratingError.value = null;

  if (!foodId.value) {
    ratingError.value = 'No Food ID provided to fetch rating.';
    ratingLoading.value = false;
    return;
  }

  try {
    const response = await fetch(`http://localhost:8080/menu/rating-times-info?food_id=${foodId.value}`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
        // This endpoint probably doesn't need auth, but add if it becomes protected:
        // 'Authorization': `Bearer ${userToken.value}`,
      },
    });

    const data = await response.json();

    if (response.ok && data.success) {
      ratingInfo.value = {
        Rating: data.rating || 0, // Default to 0 if not present
        Times: data.times_ordered || 0 // Default to 0 if not present
      };
    } else {
      ratingError.value = data.message || 'Failed to fetch rating information.';
    }
  } catch (err) {
    ratingError.value = `Error fetching rating: ${err.message}.`;
    console.error('Error fetching food rating:', err);
  } finally {
    ratingLoading.value = false;
  }
};

// Add to cart functionality (assuming you have a cart store)
const addToCart = (foodItem) => {
  cartStore.addItem(foodItem);
  // Optional: show a temporary success message
  cartStore.itemAddedSuccessfully = true; // Assuming cart store has this reactive prop
  setTimeout(() => {
    cartStore.itemAddedSuccessfully = false;
  }, 2000);
};

// On component mount
onMounted(async () => {
  foodId.value = route.params.id; // Get the ID from the URL parameter

  if (foodId.value) {
    await Promise.all([
      fetchFoodDetails(),
      fetchFoodRating()
    ]);
  } else {
    loading.value = false;
    ratingLoading.value = false;
    error.value = 'No food ID provided in the URL.';
  }
});
</script>

<style scoped>
/* Add your scoped styles here */
</style>