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
      <li v-for="item in cartStore.items" :key="item.FoodID" class="py-3 flex items-center">
        <img
          v-if="item.Picture && item.Picture.String"
          :src="item.Picture.String"
          :alt="item.FoodName"
          class="w-16 h-16 object-cover rounded-md mr-3 shadow-sm"
        />
        <div v-else class="w-16 h-16 bg-gray-200 rounded-md mr-3 flex items-center justify-center text-gray-500 text-xs text-center">
          No Image
        </div>
        <div class="flex-grow">
          <p class="font-semibold text-gray-800 text-base truncate">{{ item.FoodName }}</p>
          <p class="text-sm text-gray-600">${{ item.Price ? item.Price.toFixed(2) : '0.00' }} x {{ item.quantity }}</p>
        </div>
        <div class="flex items-center ml-auto">
          <button @click="cartStore.decrementItem(item.FoodID)" class="text-gray-500 hover:text-gray-700 p-1 rounded-full">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 24 24" fill="currentColor">
              <path d="M12 2c-.55 0-1 .45-1 1v18c0 .55.45 1 1 1s1-.45 1-1V3c0-.55-.45-1-1-1z"/>
            </svg>
          </button>
          <span class="mx-2 font-medium">{{ item.quantity }}</span>
          <button @click="cartStore.incrementItem(item.FoodID)" class="text-gray-500 hover:text-gray-700 p-1 rounded-full">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 24 24" fill="currentColor">
              <path d="M12 4c-.55 0-1 .45-1 1v6H5c-.55 0-1 .45-1 1s.45 1 1 1h6v6c0 .55.45 1 1 1s1-.45 1-1v-6h6c.55 0 1-.45 1-1s-.45-1-1-1h-6V5c0-.55-.45-1-1-1z"/>
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

    <button
      @click="checkout"
      :disabled="cartStore.items.length === 0 || checkoutLoading"
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
import { ref } from 'vue';
import { useCartStore } from '../stores/cart';
import { useRouter } from 'vue-router'; // Import useRouter

const props = defineProps({
  isVisible: {
    type: Boolean,
    default: false,
  },
});

const emit = defineEmits(['show', 'hide']);

const cartStore = useCartStore();
const router = useRouter(); // Initialize router

const isDelivery = ref(false); // New reactive state for delivery option
const checkoutLoading = ref(false);
const checkoutSuccess = ref(false);
const checkoutError = ref(null);

const checkout = async () => {
  checkoutLoading.value = true;
  checkoutSuccess.value = false;
  checkoutError.value = null;

  const userToken = localStorage.getItem('userToken');
  if (!userToken) {
    checkoutError.value = 'Please log in to place an order.';
    checkoutLoading.value = false;
    router.push('/login'); // Redirect to login if no token
    return;
  }

  // Decode token to get UserID (assuming user_id is in the JWT payload)
  let userID = null;
  try {
    const base64Url = userToken.split('.')[1];
    const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
    const payload = JSON.parse(atob(base64));
    userID = payload.user_id; // Adjust if your JWT claims differ (e.g., 'sub', 'id')
    if (!userID) {
      throw new Error('User ID not found in token.');
    }
  } catch (e) {
    console.error('Error decoding JWT for UserID:', e);
    checkoutError.value = 'Authentication error. Please log in again.';
    checkoutLoading.value = false;
    localStorage.removeItem('userToken');
    router.push('/login');
    return;
  }

  // Prepare OrderInfo string (summary of items for the order overall)
  const orderInfoArray = cartStore.items.map(item => ({
    FoodID: item.FoodID,
    FoodName: item.FoodName,
    Quantity: item.quantity,
    Price: item.Price, // Include price for accurate summary
  }));
  const orderInfoString = JSON.stringify(orderInfoArray);

  // Prepare OrderItems for the backend (only FoodID and Quantity)
  const orderItems = cartStore.items.map(item => ({
    FoodID: item.FoodID,
    Quantity: item.quantity,
  }));

  const requestBody = {
    user_id: userID,
    order_info: orderInfoString,
    is_ranged: isDelivery.value, // True for delivery, false for pickup
    order_items: orderItems,
  };

  console.log('Sending order request:', requestBody); // For debugging

  try {
    const response = await fetch('http://localhost:8080/orders', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${userToken}`,
      },
      body: JSON.stringify(requestBody),
    });

    const data = await response.json();

    if (response.ok && data.success) {
      checkoutSuccess.value = true;
      cartStore.clearCart(); // Clear the cart on successful order
      emit('hide'); // Hide the cart overlay
      // --- CHANGE IS HERE ---
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
      // --- END CHANGE ---
    } else {
      checkoutError.value = data.message || 'Failed to create order.';
      if (response.status === 401 || response.status === 403) {
        localStorage.removeItem('userToken');
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