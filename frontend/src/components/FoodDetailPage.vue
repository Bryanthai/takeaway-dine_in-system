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
      <router-link to="/menu" class="text-blue-600 hover:underline mt-4 block">Back to Menu</router-link>
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
        <h3 class="text-3xl font-bold text-gray-900 mb-2">
          <span v-if="!isEditing">{{ food.FoodName }}</span>
          <input v-else v-model="formFood.FoodName" type="text" class="form-input mt-1 block w-full border-gray-300 rounded-md shadow-sm" />
        </h3>
        <p class="text-indigo-600 text-2xl font-bold mb-4">
          $<span v-if="!isEditing">{{ food.Price.toFixed(2) }}</span>
          <input v-else v-model.number="formFood.Price" type="number" step="0.01" class="form-input mt-1 block w-full border-gray-300 rounded-md shadow-sm" />
        </p>

        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4 mb-6">
          <p class="text-gray-700">
            <span class="font-semibold">Category:</span>
            <span v-if="!isEditing && foodTags.length">{{ foodTags.join(', ') }}</span>
            <span v-else-if="!isEditing && !foodTags.length">N/A</span>
            <input v-else v-model="formFood.FoodTag" type="text" class="form-input mt-1 block w-full border-gray-300 rounded-md shadow-sm" />
          </p>
          <p class="text-gray-700">
            <span class="font-semibold">Preparation Time:</span>
            <span v-if="!isEditing">{{ food.TimeNeeded }} minutes</span>
            <input v-else v-model.number="formFood.TimeNeeded" type="number" class="form-input mt-1 block w-full border-gray-300 rounded-md shadow-sm" />
          </p>
          <p class="text-gray-700">
            <span class="font-semibold">Availability:</span>
            <span v-if="!isEditing">{{ food.LongRange ? 'Delivery' : 'Pickup Only' }}</span>
            <span v-else class="flex items-center mt-1">
              <input
                type="checkbox"
                id="longRangeEdit"
                v-model="formFood.LongRange"
                class="form-checkbox h-4 w-4 text-indigo-600 transition duration-150 ease-in-out mr-2"
              />
              <label for="longRangeEdit" class="text-sm text-gray-700">Available for Delivery</label>
            </span>
          </p>

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
          <p v-if="!isEditing" class="text-gray-700">{{ food.Description || 'No description available.' }}</p>
          <textarea v-else v-model="formFood.Description" rows="3" class="form-textarea mt-1 block w-full border-gray-300 rounded-md shadow-sm"></textarea>
        </div>

        <div class="mb-6">
          <p class="text-gray-800 font-semibold mb-2">Ingredients:</p>
          <p v-if="!isEditing" class="text-gray-700">{{ food.Ingredients || 'No ingredients listed.' }}</p>
          <textarea v-else v-model="formFood.Ingredients" rows="2" class="form-textarea mt-1 block w-full border-gray-300 rounded-md shadow-sm"></textarea>
        </div>

        <div class="mb-6">
          <p class="text-gray-800 font-semibold mb-2">Additional Info:</p>
          <p v-if="!isEditing" class="text-gray-700 italic">{{ food.Info && food.Info.Valid ? food.Info.String : 'No additional info.' }}</p>
          <textarea v-else v-model="formFood.Info" rows="2" class="form-textarea mt-1 block w-full border-gray-300 rounded-md shadow-sm"></textarea>
        </div>

        <div v-if="userStore.isAdmin" class="mt-8 pt-4 border-t border-gray-200">
          <h4 class="text-xl font-bold text-gray-800 mb-4">Admin Actions</h4>

          <div v-if="changeLoading || deleteLoading" class="flex justify-center items-center py-4">
            <svg class="animate-spin h-8 w-8 text-indigo-500" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            <p class="ml-3 text-lg text-gray-700">Processing action...</p>
          </div>
          <div v-if="changeError" class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative mb-4" role="alert">
            <strong class="font-bold">Change Error!</strong>
            <span class="block sm:inline ml-2">{{ changeError }}</span>
          </div>
          <div v-if="deleteError" class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative mb-4" role="alert">
            <strong class="font-bold">Delete Error!</strong>
            <span class="block sm:inline ml-2">{{ deleteError }}</span>
          </div>
          <div v-if="changeSuccess" class="bg-green-100 border border-green-400 text-green-700 px-4 py-3 rounded relative mb-4" role="alert">
            <strong class="font-bold">Success!</strong>
            <span class="block sm:inline ml-2">Food info updated!</span>
          </div>
          <div v-if="deleteSuccess" class="bg-green-100 border border-green-400 text-green-700 px-4 py-3 rounded relative mb-4" role="alert">
            <strong class="font-bold">Success!</strong>
            <span class="block sm:inline ml-2">Food deleted! Redirecting...</span>
          </div>


          <div v-if="!isEditing" class="flex flex-col sm:flex-row gap-4">
            <button
              @click="startEditing"
              class="w-full bg-indigo-600 text-white px-6 py-3 rounded-lg font-semibold text-lg hover:bg-indigo-700 transition duration-200 flex items-center justify-center"
            >
              <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z" />
              </svg>
              Edit Food Info
            </button>
            <button
              @click="deleteFood"
              class="w-full bg-red-600 text-white px-6 py-3 rounded-lg font-semibold text-lg hover:bg-red-700 transition duration-200 flex items-center justify-center"
            >
              <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
              </svg>
              Delete Food
            </button>
          </div>

          <div v-else class="flex flex-col sm:flex-row gap-4">
            <button
              @click="saveChanges"
              :disabled="changeLoading"
              class="w-full bg-green-600 text-white px-6 py-3 rounded-lg font-semibold text-lg hover:bg-green-700 transition duration-200 flex items-center justify-center"
            >
              <svg v-if="!changeLoading" xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              <span v-else class="animate-pulse">Saving...</span>
              Save Changes
            </button>
            <button
              @click="cancelEditing"
              class="w-full bg-gray-600 text-white px-6 py-3 rounded-lg font-semibold text-lg hover:bg-gray-700 transition duration-200 flex items-center justify-center"
            >
              <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
              Cancel
            </button>
          </div>
        </div>

        <button
          v-if="!userStore.isAdmin"
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
import { ref, onMounted, watch } from 'vue'; // Import watch
import { useRoute, useRouter } from 'vue-router';
import { useCartStore } from '../stores/cart';
import { useUserStore } from '../stores/user';

const route = useRoute();
const router = useRouter();
const cartStore = useCartStore();
const userStore = useUserStore();

const foodId = ref(null);
const food = ref(null);
const ratingInfo = ref(null);
const foodTags = ref([]); // New ref to hold the fetched food tags

const loading = ref(true);
const error = ref(null);
const ratingLoading = ref(true);
const ratingError = ref(null);

const isEditing = ref(false);
const formFood = ref({});

const changeLoading = ref(false);
const changeError = ref(null);
const changeSuccess = ref(false);

const deleteLoading = ref(false);
const deleteError = ref(null);
const deleteSuccess = ref(false);

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
        'Authorization': userStore.userToken ? `Bearer ${userStore.userToken}` : '',
      },
    });

    if (!response.ok) {
      const errorData = await response.json();
      throw new Error(errorData.message || `HTTP error! Status: ${response.status}`);
    }

    const data = await response.json();
    food.value = data;

    if (!food.value || !food.value.FoodID) {
      error.value = 'Food item not found or invalid data received.';
      food.value = null;
    } else {
        // Initialize formFood with current food data (excluding FoodTag for now)
        formFood.value = {
            FoodName: food.value.FoodName,
            // FoodTag will be populated from the /foods/tags endpoint results for editing
            Price: food.value.Price,
            Info: food.value.Info.Valid ? food.value.Info.String : '',
            Ingredients: food.value.Ingredients,
            TimeNeeded: food.value.TimeNeeded,
            LongRange: food.value.LongRange,
            Description: food.value.Description,
        };
        // After food details are loaded, fetch its tags
        await fetchFoodTags(food.value.FoodName);
    }

  } catch (err) {
    error.value = `Error fetching food details: ${err.message}.`;
    console.error('Error fetching food details:', err);
  } finally {
    loading.value = false;
  }
};

// New function to fetch food tags
const fetchFoodTags = async (foodName) => {
  foodTags.value = []; // Clear previous tags
  if (!foodName) return;

  try {
    const response = await fetch(`http://localhost:8080/foods/tags?food_name=${encodeURIComponent(foodName)}`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': userStore.userToken ? `Bearer ${userStore.userToken}` : '',
      },
    });

    if (!response.ok) {
        // Handle non-200 responses, e.g., if food_name is not found or no tags exist
        if (response.status === 404) {
            console.warn(`No tags found for food: ${foodName}`);
        } else {
            const errorData = await response.json();
            throw new Error(errorData.message || `HTTP error! Status: ${response.status}`);
        }
    } else {
        const data = await response.json();
        // Assuming the backend returns an array of strings directly, e.g., ["tag1", "tag2"]
        if (Array.isArray(data)) {
            foodTags.value = data;
        } else {
            console.warn('Unexpected data format for food tags:', data);
            foodTags.value = [];
        }
    }
  } catch (err) {
    console.error(`Error fetching food tags for ${foodName}:`, err);
    foodTags.value = [];
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
      },
    });

    const data = await response.json();

    if (response.ok && data.success) {
      ratingInfo.value = {
        Rating: data.rating || 0,
        Times: data.times_ordered || 0
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

// Admin Action: Start Editing
const startEditing = () => {
  isEditing.value = true;
  // Deep copy food data to formFood to allow cancellation without affecting original
  formFood.value = {
    FoodName: food.value.FoodName,
    // When starting edit, pre-fill the FoodTag input with the current tags (joined by comma for editing)
    FoodTag: foodTags.value.join(', '),
    Price: food.value.Price,
    Info: food.value.Info.Valid ? food.value.Info.String : '',
    Ingredients: food.value.Ingredients,
    TimeNeeded: food.value.TimeNeeded,
    LongRange: food.value.LongRange,
    Description: food.value.Description,
  };
  clearAdminMessages();
};

// Admin Action: Cancel Editing
const cancelEditing = () => {
  isEditing.value = false;
  clearAdminMessages();
};

// Admin Action: Save Changes (PUT /foods/change-info)
const saveChanges = async () => {
  changeLoading.value = true;
  clearAdminMessages();

  if (!userStore.isAdmin) {
    changeError.value = "Unauthorized: Not an admin.";
    changeLoading.value = false;
    return;
  }

  try {
    const response = await fetch('http://localhost:8080/foods/change-info', {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${userStore.userToken}`,
      },
      body: JSON.stringify({
        food_id: foodId.value,
        food_name: formFood.value.FoodName,
        // Send the updated FoodTag from the form.
        // If the backend expects an array, you might need to split formFood.FoodTag by comma:
        food_tag: formFood.value.FoodTag, // Assuming backend can handle comma-separated string or expects single tag
        price: formFood.value.Price,
        info: formFood.value.Info,
        ingredients: formFood.value.Ingredients,
        time_needed: formFood.value.TimeNeeded,
        long_range: formFood.value.LongRange,
        description: formFood.value.Description,
      }),
    });

    const data = await response.json();

    if (response.ok && data.success) {
      changeSuccess.value = true;
      isEditing.value = false;
      await fetchFoodDetails(); // Re-fetch all data, including tags, to update displayed data
      setTimeout(() => {
        changeSuccess.value = false;
      }, 3000);
    } else {
      changeError.value = data.message || 'Failed to update food information.';
      if (response.status === 401 || response.status === 403) {
          userStore.logout();
          router.push('/login');
      }
    }
  } catch (err) {
    changeError.value = `Network or unexpected error: ${err.message}`;
    console.error('Error changing food info:', err);
  } finally {
    changeLoading.value = false;
  }
};

// Admin Action: Delete Food (DELETE /foods)
const deleteFood = async () => {
  if (!confirm(`Are you sure you want to DELETE "${food.value.FoodName}"? This action cannot be undone.`)) {
    return;
  }

  deleteLoading.value = true;
  clearAdminMessages();

  if (!userStore.isAdmin) {
    deleteError.value = "Unauthorized: Not an admin.";
    deleteLoading.value = false;
    return;
  }

  try {
    const response = await fetch('http://localhost:8080/foods', {
      method: 'DELETE',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${userStore.userToken}`,
      },
      body: JSON.stringify({
        food_name: food.value.FoodName,
      }),
    });

    const data = await response.json();

    if (response.ok && data.success) {
      deleteSuccess.value = true;
      alert(`Food "${food.value.FoodName}" deleted successfully!`);
      router.push('/admin/create-food');
    } else {
      deleteError.value = data.message || 'Failed to delete food.';
      if (response.status === 401 || response.status === 403) {
          userStore.logout();
          router.push('/login');
      }
    }
  } catch (err) {
    deleteError.value = `Network or unexpected error: ${err.message}`;
    console.error('Error deleting food:', err);
  } finally {
    deleteLoading.value = false;
  }
};

const clearAdminMessages = () => {
    changeError.value = null;
    changeSuccess.value = false;
    deleteError.value = null;
    deleteSuccess.value = false;
};

// Add to cart functionality
const addToCart = (foodItem) => {
  cartStore.addItem(foodItem);
  cartStore.itemAddedSuccessfully = true;
  setTimeout(() => {
    cartStore.itemAddedSuccessfully = false;
  }, 2000);
};

// Watch for changes in food.value.FoodName to re-fetch tags if the food name itself changes
// (though in this context, the food name typically won't change while on the same details page)
watch(() => food.value?.FoodName, (newFoodName) => {
  if (newFoodName) {
    fetchFoodTags(newFoodName);
  }
}, { immediate: false }); // immediate: false as food details are fetched on mount already.

// On component mount
onMounted(async () => {
  foodId.value = route.params.id;

  if (foodId.value) {
    // Fetch food details and rating concurrently
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
/* Basic form styling for inputs/textareas */
.form-input, .form-textarea {
  @apply px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500;
}
</style>