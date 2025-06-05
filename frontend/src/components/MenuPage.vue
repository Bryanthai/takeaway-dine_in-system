<template>
  <div class="bg-gray-100 p-6 min-h-screen relative">
    <div class="absolute top-4 right-4 z-50">
      <button
        @mouseenter="showCartOverlay"
        @mouseleave="hideCartOverlayDelayed"
        class="relative bg-white p-3 rounded-full shadow-lg hover:shadow-xl transition-shadow duration-200"
        aria-label="Shopping Cart"
      >
        <svg xmlns="http://www.w3.org/2000/svg" class="h-8 w-8 text-indigo-600" viewBox="0 0 24 24" fill="currentColor">
          <path d="M7 18c-1.1 0-1.99.9-1.99 2S5.9 22 7 22s2-.9 2-2-.9-2-2-2zM17 18c-1.1 0-1.99.9-1.99 2S15.9 22 17 22s2-.9 2-2-.9-2-2-2zM7.42 14.07l.79-.79h8.33c.72 0 1.35-.43 1.63-1.07L21.57 6.44c.26-.54-.08-1.16-.62-1.42-.54-.26-1.16.08-1.42.62L17.7 11H7.83L6.1 3.52c-.17-.74-.88-1.28-1.66-1.28H3c-.55 0-1 .45-1 1s.45 1 1 1h1.56c.46 0 .86.3.97.74L6.96 14.86c.21.93 1.05 1.6 2.01 1.6h7.45c.55 0 1-.45 1-1s-.45-1-1-1H9.01c-.49 0-.9-.37-1.01-.84z"/>
        </svg>
        <span v-if="cartStore.totalItems > 0" class="absolute -top-1 -right-1 bg-red-500 text-white text-xs font-bold rounded-full h-5 w-5 flex items-center justify-center">
          {{ cartStore.totalItems }}
        </span>
      </button>

      <ShoppingCart
        :isVisible="isCartVisible"
        @show="showCartOverlay"
        @hide="hideCartOverlay"
      />
    </div>

    <h2 class="text-4xl font-extrabold text-gray-900 mb-8 text-center">
      Our Delicious Menu
    </h2>

    <div v-if="initialLoading" class="flex justify-center items-center h-48">
      <svg class="animate-spin h-10 w-10 text-indigo-500" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
        <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
      </svg>
      <p class="ml-3 text-lg text-gray-700">Loading menu...</p>
    </div>

    <div v-else-if="initialError" class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative" role="alert">
      <strong class="font-bold">Error!</strong>
      <span class="block sm:inline ml-2">{{ initialError }}</span>
    </div>

    <div v-else-if="foodCategories.length === 0 && allFoods.length === 0 && recommendedFoods.length === 0" class="bg-yellow-100 border border-yellow-400 text-yellow-700 px-4 py-3 rounded relative" role="alert">
      <strong class="font-bold">No Menu Items Found!</strong>
      <span class="block sm:inline ml-2">There are no food categories, all items, or recommended items to display.</span>
    </div>

    <div v-else>
      <div class="flex flex-wrap justify-center mb-8 border-b border-gray-300">
        <button
          @click="selectTab('All Food')"
          :class="[
            'px-6 py-3 text-lg font-semibold rounded-t-lg transition-colors duration-300',
            activeTab === 'All Food'
              ? 'bg-white text-indigo-700 border-indigo-500 border-t border-l border-r -mb-px'
              : 'bg-gray-200 text-gray-700 hover:bg-gray-300 hover:text-gray-900'
          ]"
        >
          All Food
        </button>

        <button
          v-for="(category, index) in foodCategories"
          :key="category.tag"
          @click="selectTab(category.tag)"
          :class="[
            'px-6 py-3 text-lg font-semibold rounded-t-lg transition-colors duration-300',
            activeTab === category.tag
              ? 'bg-white text-indigo-700 border-indigo-500 border-t border-l border-r -mb-px'
              : 'bg-gray-200 text-gray-700 hover:bg-gray-300 hover:text-gray-900'
          ]"
        >
          {{ category.tag }}
        </button>

        <button
          @click="selectTab('Recommendations')"
          :class="[
            'px-6 py-3 text-lg font-semibold rounded-t-lg transition-colors duration-300',
            activeTab === 'Recommendations'
              ? 'bg-white text-indigo-700 border-indigo-500 border-t border-l border-r -mb-px'
              : 'bg-gray-200 text-gray-700 hover:bg-gray-300 hover:text-gray-900'
          ]"
        >
          Recommendations
        </button>
      </div>

      <div v-if="currentFoodList.length > 0" class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6">
        <div
          v-for="food in currentFoodList"
          :key="food.FoodID"
          class="bg-white rounded-xl shadow-lg overflow-hidden transform transition-transform duration-300 hover:scale-105 hover:shadow-xl"
        >
          <router-link :to="{ name: 'FoodDetail', params: { id: food.FoodID }}" class="block">
            <div class="w-full h-48 flex items-center justify-center bg-gray-200">
              <img
                v-if="food.Picture && food.Picture.String"
                :src="food.Picture.String" :alt="food.FoodName"
                class="w-full h-full object-cover"
              />
              <div v-else class="w-full h-full bg-gray-400 flex items-center justify-center text-gray-600 text-xl font-semibold">
                No Image
              </div>
            </div>

            <div class="p-4">
              <h3 class="text-xl font-bold text-gray-900 mb-2 truncate">{{ food.FoodName }}</h3>
              <p class="text-gray-600 text-sm mb-2">
                Tags:
                <span v-if="food.tags && food.tags.length > 0" class="font-semibold">{{ food.tags.join(', ') }}</span>
                <span v-else class="font-semibold text-gray-500">N/A</span>
              </p>
              <p class="text-gray-600 text-sm mb-2">Time: <span class="font-semibold">{{ food.TimeNeeded }} mins</span></p>
              <p class="text-gray-600 text-sm">
                Delivery:
                <span v-if="food.long_range" class="font-semibold text-green-600">Available</span>
                <span v-else class="font-semibold text-red-600">Not Available</span>
              </p>
            </div>
          </router-link>

          <div class="flex justify-between items-center p-4 pt-0"> <span class="text-2xl font-bold text-indigo-600">${{ food.Price.toFixed(2) }}</span>
              <button
                @click="cartStore.addItem(food)"
                class="bg-green-500 text-white px-4 py-2 rounded-lg font-semibold hover:bg-green-600 transition duration-200 flex items-center"
              >
                <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-1" viewBox="0 0 24 24" fill="currentColor">
                  <path d="M12 4c-.55 0-1 .45-1 1v6H5c-.55 0-1 .45-1 1s.45 1 1 1h6v6c0 .55.45 1 1 1s1-.45 1-1v-6h6c.55 0 1-.45 1-1s-.45-1-1-1h-6V5c0-.55-.45-1-1-1z"/>
                </svg>
                Add
              </button>
            </div>
        </div>
      </div>
      <div v-else class="bg-blue-100 border border-blue-400 text-blue-700 px-4 py-3 rounded relative text-center" role="alert">
        <span v-if="activeTab === 'Recommendations'">No recommendations available for you at the moment.</span>
        <span v-else-if="activeTab === 'All Food'">No food items found at all.</span>
        <span v-else>No food items found for this category.</span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed, watch } from 'vue';
import { useCartStore } from '../stores/cart';
import ShoppingCart from './ShoppingCart.vue';
import { useRouter } from 'vue-router';
import { useUserStore } from '../stores/user';

const router = useRouter();
const cartStore = useCartStore();
const userStore = useUserStore();

const foodCategories = ref([]);
const recommendedFoods = ref([]);
const allFoods = ref([]);

const activeTab = ref(null);
const initialLoading = ref(true);
const initialError = ref(null);
const recommendationsLoading = ref(false);
const recommendationsError = ref(null);
const allFoodsLoading = ref(false);
const allFoodsError = ref(null);

const isCartVisible = ref(false);
let cartHideTimeout = null;

const showCartOverlay = () => {
  clearTimeout(cartHideTimeout);
  isCartVisible.value = true;
};
const hideCartOverlay = () => {
  clearTimeout(cartHideTimeout);
  isCartVisible.value = false;
};
const hideCartOverlayDelayed = () => {
  clearTimeout(cartHideTimeout);
  cartHideTimeout = setTimeout(() => {
    isCartVisible.value = false;
  }, 500);
};

const currentFoodList = computed(() => {
  if (activeTab.value === 'All Food') {
    return allFoods.value;
  } else if (activeTab.value === 'Recommendations') {
    return recommendedFoods.value;
  }
  const selectedCategory = foodCategories.value.find(cat => cat.tag === activeTab.value);
  return selectedCategory ? selectedCategory.foods : [];
});

// CORRECTED FUNCTION: fetchFoodTags to handle direct array response
const fetchFoodTags = async (food) => {
  if (!food || !food.FoodName) return;

  try {
    const response = await fetch(`http://localhost:8080/foods/tags?food_name=${encodeURIComponent(food.FoodName)}`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
    });

    if (!response.ok) {
      console.error(`Failed to fetch tags for ${food.FoodName}: HTTP error! Status: ${response.status}`);
      food.tags = []; // Ensure it's an empty array on HTTP error
      return;
    }

    // *** THIS IS THE CRUCIAL CHANGE: Directly parse as an array ***
    const tagsArray = await response.json();

    if (Array.isArray(tagsArray)) {
      food.tags = tagsArray;
    } else {
      // Fallback if the API returns something unexpected (not an array)
      console.warn(`Unexpected data format for tags for ${food.FoodName}: Expected an array, but got`, tagsArray);
      food.tags = [];
    }
  } catch (err) {
    console.error(`Error fetching tags for ${food.FoodName}:`, err);
    food.tags = []; // Set to empty array on network or parsing error
  }
};

const fetchFoodByTypes = async () => {
  initialError.value = null;
  try {
    const response = await fetch('http://localhost:8080/menu/sort-type', {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
    });

    if (!response.ok) {
      const errorData = await response.json();
      throw new Error(errorData.message || `HTTP error! Status: ${response.status}`);
    }

    const data = await response.json();
    if (data.success) {
      foodCategories.value = data.data || [];
      // Initialize tags property for reactivity
      foodCategories.value.forEach(category => {
        if (category.foods) {
          category.foods.forEach(food => {
            food.tags = []; // Initialize as empty array
          });
        }
      });
    } else {
      throw new Error(data.message || 'Failed to fetch food categories.');
    }
  } catch (err) {
    initialError.value = `Failed to load menu categories: ${err.message}`;
    console.error('Fetch food by types error:', err);
  }
};

const fetchRecommendedFood = async () => {
  recommendationsLoading.value = true;
  recommendationsError.value = null;

  if (!userStore.userToken) {
    recommendationsError.value = 'Authentication token is missing. Please log in to see recommendations.';
    recommendedFoods.value = [];
    recommendationsLoading.value = false;
    return;
  }

  try {
    const response = await fetch('http://localhost:8080/menu/sort-by-usertag', {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${userStore.userToken}`,
      },
    });

    const data = await response.json();

    if (response.ok && data.success) {
      recommendedFoods.value = data.foods || [];
      // Initialize tags property for reactivity
      recommendedFoods.value.forEach(food => {
        food.tags = []; // Initialize as empty array
      });
    } else {
      recommendedFoods.value = [];
      recommendationsError.value = data.message || 'Failed to retrieve recommendations.';
      if (response.status === 401 || response.status === 403) {
          userStore.logout();
          router.push('/login');
      }
    }
  } catch (err) {
    recommendationsError.value = `Error fetching recommendations: ${err.message}`;
    console.error('Fetch recommended food error:', err);
  } finally {
    recommendationsLoading.value = false;
  }
};

const fetchAllFood = async () => {
  allFoodsLoading.value = true;
  allFoodsError.value = null;

  try {
    const response = await fetch('http://localhost:8080/menu', {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
    });

    if (!response.ok) {
      const errorData = await response.json();
      throw new Error(errorData.message || `HTTP error! Status: ${response.status}`);
    }

    const data = await response.json();
    if (data.success) {
      allFoods.value = data.foods || [];
      // Initialize tags property for reactivity
      allFoods.value.forEach(food => {
        food.tags = []; // Initialize as empty array
      });
    } else {
      throw new Error(data.message || 'Failed to retrieve all food items.');
    }
  } catch (err) {
    allFoodsError.value = `Error fetching all food items: ${err.message}`;
    console.error('Fetch all food error:', err);
  } finally {
    allFoodsLoading.value = false;
  }
};

const selectTab = async (tag) => {
  activeTab.value = tag;
  if (tag === 'Recommendations' && recommendedFoods.value.length === 0 && !recommendationsLoading.value) {
      if (!recommendationsError.value || recommendationsError.value === 'Authentication token is missing. Please log in to see recommendations.') {
          await fetchRecommendedFood(); // Await to ensure recommendedFoods is populated before fetching tags
      }
  } else if (tag === 'All Food' && allFoods.value.length === 0 && !allFoodsLoading.value) {
      await fetchAllFood(); // Await to ensure allFoods is populated before fetching tags
  }
  // Trigger fetching tags for the newly selected list
  currentFoodList.value.forEach(food => {
    // Only fetch if tags haven't been fetched OR if the existing tags array is empty
    if (!food.tags || food.tags.length === 0) {
      fetchFoodTags(food);
    }
  });
};

onMounted(async () => {
  initialLoading.value = true;
  await Promise.all([
    fetchFoodByTypes(),
    fetchAllFood()
  ]);

  if (foodCategories.value.length > 0) {
    activeTab.value = foodCategories.value[0].tag;
  } else if (allFoods.value.length > 0) {
    activeTab.value = 'All Food';
  } else {
    activeTab.value = 'Recommendations';
    await fetchRecommendedFood(); // Await to ensure recommendedFoods is populated
  }
  initialLoading.value = false;

  // Fetch tags for the initial active tab's food items
  currentFoodList.value.forEach(food => {
    if (!food.tags || food.tags.length === 0) {
      fetchFoodTags(food);
    }
  });
});

watch(currentFoodList, (newList) => {
  // When the currentFoodList changes (due to tab switch or initial load),
  // fetch tags for all items in the new list.
  newList.forEach(food => {
    if (!food.tags || food.tags.length === 0) {
      fetchFoodTags(food);
    }
  });
}, { immediate: true }); // immediate: true runs the watcher once immediately on component mount
</script>

<style scoped>
/* Scoped styles specific to this component */
</style>