<template>
  <div
    v-if="isVisible"
    @mouseenter="$emit('show')"
    @mouseleave="$emit('hide')"
    class="cart-overlay fixed right-4 top-20 w-80 bg-white rounded-lg shadow-xl p-6 z-50 transform transition-all duration-300"
    :class="{ 'translate-x-0 opacity-100': isVisible, 'translate-x-full opacity-0': !isVisible }"
  >
    <h3 class="text-2xl font-bold text-gray-900 mb-4">Your Cart</h3>

    <div v-if="cartStore.items.length === 0" class="text-gray-600 italic text-center py-4">
      Your cart is empty.
    </div>

    <ul v-else class="divide-y divide-gray-200 max-h-60 overflow-y-auto mb-4">
      <li v-for="item in cartStore.items" :key="item.food.FoodID" class="py-3 flex items-center">
        <img
          v-if="item.food.Picture && item.food.Picture.String"
          :src="item.food.Picture.String"
          :alt="item.food.FoodName"
          class="w-16 h-16 object-cover rounded-md mr-3 shadow-sm"
        />
        <div v-else class="w-16 h-16 bg-gray-200 rounded-md mr-3 flex items-center justify-center text-gray-500 text-xs text-center">
          No Image
        </div>
        <div class="flex-grow">
          <p class="font-semibold text-gray-800 text-base truncate">{{ item.food.FoodName }}</p>
          <p class="text-sm text-gray-600">${{ item.food.Price ? item.food.Price.toFixed(2) : '0.00' }} x {{ item.quantity }}</p>
        </div>
        <div class="flex items-center ml-auto">
          <button @click="cartStore.decrementItem(item.food.FoodID)" class="text-gray-500 hover:text-gray-700 p-1 rounded-full">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 24 24" fill="currentColor">
              <path d="M19 13H5v-2h14v2z"/>
            </svg>
          </button>
          <span class="mx-2 font-medium">{{ item.quantity }}</span>
          <button @click="cartStore.incrementItem(item.food.FoodID)" class="text-gray-500 hover:text-gray-700 p-1 rounded-full">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 24 24" fill="currentColor">
              <path d="M19 13h-6v6h-2v-6H5v-2h6V5h2v6h6v2z"/>
            </svg>
          </button>
        </div>
      </li>
    </ul>

    <div class="flex justify-between items-center text-lg font-bold text-gray-900 mt-4 border-t pt-4">
      <span>Total:</span>
      <span>${{ cartStore.totalPrice.toFixed(2) }}</span>
    </div>

    <div class="mt-4">
      <label class="inline-flex items-center">
        <input
          type="checkbox"
          v-model="isDelivery"
          class="form-checkbox h-5 w-5 text-indigo-600 transition duration-150 ease-in-out"
        />
        <span class="ml-2 text-gray-700 font-medium">Order for Delivery</span>
      </label>
    </div>

    <div v-if="isDelivery" class="mt-4">
      <label class="inline-flex items-center mb-2">
        <input
          type="checkbox"
          v-model="useRegisteredAddress"
          @change="toggleRegisteredAddress"
          class="form-checkbox h-5 w-5 text-indigo-600 transition duration-150 ease-in-out"
        />
        <span class="ml-2 text-gray-700 font-medium">Deliver to my registered address</span>
      </label>
      <label for="deliveryAddress" class="block text-sm font-medium text-gray-700">Delivery Address</label>
      <input
        type="text"
        id="deliveryAddress"
        v-model="deliveryAddress"
        :disabled="useRegisteredAddress"
        class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
        placeholder="Enter delivery address"
      />
    </div>

    <button
      @click="checkout"
      :disabled="cartStore.items.length === 0 || checkoutLoading || (isDelivery && !deliveryAddress)"
      class="mt-6 w-full bg-indigo-600 text-white py-3 rounded-lg font-semibold text-lg hover:bg-indigo-700 transition duration-200 disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center"
    >
      <svg v-if="checkoutLoading" class="animate-spin -ml-1 mr-3 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
        <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
      </svg>
      {{ checkoutLoading ? 'Processing...' : 'Checkout' }}
    </button>

    <div v-if="checkoutSuccess" class="mt-4 text-green-600 text-center font-medium">
      Order placed successfully! Redirecting to payment...
    </div>
    <div v-if="checkoutError" class="mt-4 text-red-600 text-center font-medium">
      {{ checkoutError }}
    </div>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue';
import { useCartStore } from '../stores/cart';
import { useRouter } from 'vue-router';
import { useUserStore } from '../stores/user';

const props = defineProps({
  isVisible: {
    type: Boolean,
    default: false,
  },
});

const emit = defineEmits(['show', 'hide']);

const cartStore = useCartStore();
const router = useRouter();
const userStore = useUserStore();

const isDelivery = ref(false);
const deliveryAddress = ref('');
const useRegisteredAddress = ref(false); // New ref for the checkbox
const checkoutLoading = ref(false);
const checkoutSuccess = ref(false);
const checkoutError = ref(null);

// Watch for changes in `isDelivery` to reset `useRegisteredAddress`
watch(isDelivery, (newValue) => {
  if (!newValue) {
    useRegisteredAddress.value = false;
    deliveryAddress.value = ''; // Clear delivery address if delivery is unchecked
  }
});

// Function to fetch user data
const fetchUserData = async () => {
  if (!userStore.userToken) {
    return null;
  }
  try {
    const response = await fetch('http://localhost:8080/users', {
      headers: {
        'Authorization': `Bearer ${userStore.userToken}`,
      },
    });
    if (!response.ok) {
      if (response.status === 401 || response.status === 403) {
        userStore.logout();
        router.push('/login');
      }
      throw new Error('Failed to fetch user data');
    }
    const data = await response.json();
    return data; // Assuming the API returns the Account struct directly
  } catch (error) {
    console.error('Error fetching user data:', error);
    checkoutError.value = 'Failed to load user address.';
    return null;
  }
};

// Toggle registered address functionality
const toggleRegisteredAddress = async () => {
  if (useRegisteredAddress.value) {
    const userData = await fetchUserData();
    if (userData && userData.Address) {
      deliveryAddress.value = userData.Address;
    } else {
      deliveryAddress.value = '';
      useRegisteredAddress.value = false; // Uncheck if address not found
      checkoutError.value = 'No registered address found or failed to fetch.';
    }
  } else {
    deliveryAddress.value = ''; // Clear address if checkbox is unchecked
  }
};

const checkout = async () => {
  checkoutLoading.value = true;
  checkoutSuccess.value = false;
  checkoutError.value = null;

  // Validate delivery address if delivery is selected
  if (isDelivery.value && !deliveryAddress.value.trim()) {
    checkoutError.value = 'Please enter a delivery address.';
    checkoutLoading.value = false;
    return;
  }

  // Ensure user is logged in (token exists)
  if (!userStore.userToken) {
    checkoutError.value = 'Please log in to place an order.';
    checkoutLoading.value = false;
    userStore.logout(); // Clear any stale token and redirect
    router.push('/login');
    return;
  }
  const userToken = userStore.userToken; // Get token from store

  // Prepare OrderInfo string (summary of items for the order overall)
  const orderInfoArray = cartStore.items.map(item => ({
    FoodID: item.food.FoodID,
    FoodName: item.food.FoodName,
    Quantity: item.quantity,
    Price: item.food.Price,
  }));
  const orderInfoString = JSON.stringify(orderInfoArray);

  // Prepare OrderItems for the backend (ONLY FoodID and Quantity, with CORRECT keys)
  const orderItems = cartStore.items.map(item => ({
    food_id: item.food.FoodID,
    quantity: item.quantity,
  }));

  // Construct the request body including delivery_address
  const requestBody = {
    order_info: orderInfoString,
    is_ranged: isDelivery.value,
    order_items: orderItems,
    delivery_address: isDelivery.value ? deliveryAddress.value.trim() : "", // Send address only if delivery is true
  };

  console.log('Sending order request:', JSON.stringify(requestBody, null, 2));

  try {
    const response = await fetch('http://localhost:8080/orders', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${userToken}`, // User identification is handled by the token
      },
      body: JSON.stringify(requestBody),
    });

    const data = await response.json();

    if (response.ok && data.success) {
      checkoutSuccess.value = true;
      cartStore.clearCart(); // Clear the cart on successful order
      emit('hide'); // Hide the cart overlay

      const newOrderId = data.order_id; // Assuming your backend returns the new order_id
      if (newOrderId) {
        setTimeout(() => {
          router.push({ name: 'PaymentPage', params: { orderId: newOrderId } }); // Redirect to payment page for the new order
        }, 1500); // Give user a moment to see success message
      } else {
        checkoutError.value = 'Order created but order ID not returned. Redirecting to orders page.';
        setTimeout(() => {
          router.push('/orders'); // Fallback to orders page if ID is missing
        }, 1500);
      }
    } else {
      checkoutError.value = data.message || 'Failed to create order.';
      if (response.status === 401 || response.status === 403) {
        userStore.logout(); // User is unauthorized/forbidden, log them out
        router.push('/login');
      }
    }
  } catch (err) {
    checkoutError.value = `Network error: ${err.message}. Please try again.`;
    console.error('Error creating order:', err);
  } finally {
    checkoutLoading.value = false;
  }
};
</script>

<style scoped>
/* Scoped styles for the cart overlay */
.cart-overlay {
  transition: transform 0.3s ease-out, opacity 0.3s ease-out;
}
</style>