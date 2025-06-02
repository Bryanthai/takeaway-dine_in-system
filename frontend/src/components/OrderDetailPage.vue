<template>
  <div class="min-h-screen bg-gray-100 p-6">
    <h2 class="text-4xl font-extrabold text-gray-900 mb-8 text-center">
      Order Details for #{{ orderId }}
    </h2>

    <div v-if="loading" class="flex justify-center items-center h-48">
      <svg class="animate-spin h-10 w-10 text-indigo-500" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
        <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
      </svg>
      <p class="ml-3 text-lg text-gray-700">Loading order details...</p>
    </div>

    <div v-else-if="error" class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative" role="alert">
      <strong class="font-bold">Error!</strong>
      <span class="block sm:inline ml-2">{{ error }}</span>
    </div>

    <div v-else-if="!order" class="bg-yellow-100 border border-yellow-400 text-yellow-700 px-4 py-3 rounded relative text-center" role="alert">
      <strong class="font-bold">Order Not Found!</strong>
      <span class="block sm:inline ml-2">The requested order could not be loaded.</span>
    </div>

    <div v-else class="max-w-4xl mx-auto bg-white rounded-lg shadow-xl p-8">
      <div class="mb-8 pb-8 border-b border-gray-200">
        <h3 class="text-2xl font-semibold text-gray-800 mb-4">Order Summary</h3>
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4 text-lg">
          <p class="text-gray-700"><span class="font-medium">Order ID:</span> <span class="font-bold">{{ order.OrderID }}</span></p>
          <p class="text-gray-700"><span class="font-medium">User ID:</span> {{ order.UserID }}</p>
          <p class="text-gray-700">
            <span class="font-medium">Status:</span>
            <span :class="{
              'text-green-600 font-bold': order.IsDone,
              'text-yellow-600 font-bold': !order.IsDone && !order.IsPaid,
              'text-indigo-600 font-bold': !order.IsDone && order.IsPaid
            }">
              {{ order.IsDone ? 'Completed' : (order.IsPaid ? 'Processing' : 'Pending Payment') }}
            </span>
          </p>
          <p class="text-gray-700">
              <span class="font-medium">Order Placed:</span>
              {{ formatDateTime(order.OrderTime.String) }}
          </p>
          <p class="text-gray-700">
              <span class="font-medium">Estimated Completion:</span>
              {{ formatDateTime(order.EstimatedTime) }}
          </p>
          <p class="text-gray-700">
              <span class="font-medium">Delivery:</span>
              {{ order.IsRanged ? 'Delivery' : 'Pickup' }}
          </p>
          <p class="text-gray-700">
              <span class="font-medium">Paid:</span>
              <span :class="{'text-green-600 font-bold': order.IsPaid, 'text-red-600 font-bold': !order.IsPaid}">
                  {{ order.IsPaid ? 'Yes' : 'No' }}
              </span>
          </p>
          <p class="text-gray-700">
              <span class="font-medium">Deleted:</span>
              <span :class="{'text-red-600 font-bold': order.Deleted, 'text-green-600 font-bold': !order.Deleted}">
                  {{ order.Deleted ? 'Yes' : 'No' }}
              </span>
          </p>
          <div class="md:col-span-2">
            <p class="text-gray-700 font-medium mt-2">Order Info (Summary):</p>
            <ul class="list-disc list-inside text-gray-700 bg-gray-50 p-3 rounded-md">
              <li v-for="(infoItem, index) in parseOrderInfo(order.OrderInfo)" :key="'summary-' + index">
                {{ infoItem.FoodName }} (x{{ infoItem.Quantity }}) - ${{ infoItem.Price ? (infoItem.Price * infoItem.Quantity).toFixed(2) : 'N/A' }}
              </li>
            </ul>
          </div>
        </div>
      </div>

      <div class="mb-8 pb-8 border-b border-gray-200">
        <h3 class="text-2xl font-semibold text-gray-800 mb-4">Provide Order Feedback</h3>
        <div v-if="order.Feedback.Valid" class="mb-4">
            <p class="text-gray-700 font-medium">Your current feedback:</p>
            <p class="text-gray-600 italic border border-gray-300 p-3 rounded-md bg-gray-50">"{{ order.Feedback.String }}"</p>
            <p class="text-sm text-gray-500 mt-1">You can submit new feedback below to update it.</p>
        </div>
        <textarea
          v-model="orderFeedback"
          placeholder="Share your experience with this order..."
          rows="4"
          class="w-full p-3 border border-gray-300 rounded-lg focus:ring-indigo-500 focus:border-indigo-500 mb-3"
        ></textarea>
        <div class="flex items-center space-x-2">
          <button
            @click="submitOrderFeedback"
            :disabled="feedbackLoading || !orderFeedback.trim()"
            class="bg-indigo-600 text-white px-6 py-2 rounded-lg font-semibold hover:bg-indigo-700 transition duration-200 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {{ feedbackLoading ? 'Submitting...' : 'Submit Feedback' }}
          </button>
          <span v-if="feedbackSuccess" class="text-green-600 text-sm">Feedback submitted!</span>
          <span v-if="feedbackError" class="text-red-600 text-sm">{{ feedbackError }}</span>
        </div>
      </div>

      <div>
        <h3 class="text-2xl font-semibold text-gray-800 mb-4">Rate Ordered Items</h3>
        <div v-if="itemsLoading" class="flex justify-center items-center h-24">
          <svg class="animate-spin h-8 w-8 text-blue-500" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
          <p class="ml-3 text-md text-gray-600">Loading items...</p>
        </div>
        <div v-else-if="itemsError" class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative text-center" role="alert">
          <strong class="font-bold">Item Error!</strong>
          <span class="block sm:inline ml-2">{{ itemsError }}</span>
        </div>
        <div v-else-if="orderedItems.length === 0" class="bg-yellow-100 border border-yellow-400 text-yellow-700 px-4 py-3 rounded relative text-center" role="alert">
          <strong class="font-bold">No Items Found!</strong>
          <span class="block sm:inline ml-2">No individual items found for this order.</span>
        </div>
        <div v-else>
          <ul class="divide-y divide-gray-200">
            <li v-for="item in orderedItems" :key="item.ItemID" class="py-4 flex items-center justify-between">
              <div class="flex-grow">
                <p class="text-lg font-medium text-gray-900">Food ID: {{ item.FoodID }}</p>
                <p class="text-gray-600">Quantity: {{ item.Quantity }}</p>
                <p v-if="item.Rating && item.Rating.Valid" class="text-gray-600">Current Rating: {{ item.Rating.Int32 }} / 5</p>
                <p v-else class="text-gray-500 text-sm italic">Not yet rated</p>
              </div>
              <div class="ml-4 flex items-center space-x-2">
                <select
                  v-model.number="item.newRating"
                  class="p-2 border border-gray-300 rounded-md focus:ring-indigo-500 focus:border-indigo-500"
                >
                  <option :value="null" disabled>Rate</option>
                  <option v-for="n in 5" :key="n" :value="n">{{ n }}</option>
                </select>
                <button
                  @click="submitItemRating(item)"
                  :disabled="!item.newRating || item.ratingLoading"
                  class="bg-blue-600 text-white px-4 py-2 rounded-lg font-semibold hover:bg-blue-700 transition duration-200 disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  {{ item.ratingLoading ? 'Sending...' : 'Submit Rating' }}
                </button>
              </div>
            </li>
          </ul>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';

const route = useRoute();
const router = useRouter();

const orderId = ref(null);
const order = ref(null);
const orderedItems = ref([]); // Now includes 'newRating' and 'ratingLoading' for each item

const loading = ref(true);
const error = ref(null);
const itemsLoading = ref(true);
const itemsError = ref(null);

const userToken = ref(localStorage.getItem('userToken') || '');

// State for Order Feedback
const orderFeedback = ref('');
const feedbackLoading = ref(false);
const feedbackSuccess = ref(false);
const feedbackError = ref(null);

// Helper functions
const formatDateTime = (isoString) => {
  if (!isoString) return 'N/A';
  try {
    const date = new Date(isoString);
    if (isNaN(date.getTime())) {
      return 'Invalid Date';
    }
    return date.toLocaleString('en-US', {
        year: 'numeric',
        month: 'short',
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit',
        hour12: true
    });
  } catch (e) {
    console.error("Error parsing date:", isoString, e);
    return 'Invalid Date';
  }
};

const parseOrderInfo = (orderInfoString) => {
  try {
    if (!orderInfoString) return [];
    return JSON.parse(orderInfoString);
  } catch (e) {
    console.error("Error parsing order info:", orderInfoString, e);
    return [{ FoodName: "Invalid/Unreadable Info", Quantity: 1, Price: 0 }];
  }
};

// Function to fetch a single order's details
const fetchOrderDetail = async () => {
  error.value = null;
  if (!userToken.value) {
    error.value = 'Authentication token is missing. Please log in to view order details.';
    loading.value = false;
    router.push('/login');
    return;
  }

  try {
    const response = await fetch('http://localhost:8080/orders/user', {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${userToken.value}`,
      },
    });

    const data = await response.json();

    if (response.ok && data.success) {
      const foundOrder = data.orders.find(o => o.OrderID == orderId.value);
      if (foundOrder) {
        order.value = foundOrder;
        // Initialize feedback field with existing feedback if any
        orderFeedback.value = foundOrder.Feedback.Valid ? foundOrder.Feedback.String : '';
      } else {
        error.value = 'Order not found for this user.';
      }
    } else {
      error.value = data.message || 'Failed to fetch order details.';
      if (response.status === 401 || response.status === 403) {
        localStorage.removeItem('userToken');
        router.push('/login');
      }
    }
  } catch (err) {
    error.value = `Error fetching order details: ${err.message}.`;
    console.error('Error fetching order details:', err);
  } finally {
    loading.value = false;
  }
};

// Function to fetch ordered items for the current orderId
const fetchOrderedItems = async () => {
  itemsLoading.value = true;
  itemsError.value = null;

  if (!userToken.value) {
    itemsError.value = 'Authentication token is missing. Cannot fetch items.';
    itemsLoading.value = false;
    return;
  }
  if (!orderId.value) {
    itemsError.value = 'No Order ID provided to fetch items.';
    itemsLoading.value = false;
    return;
  }

  try {
    const response = await fetch(`http://localhost:8080/orders/items?order_id=${orderId.value}`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${userToken.value}`,
      },
    });

    const data = await response.json();

    if (response.ok && data.success) {
      // Add reactive properties for newRating and loading state to each item
      orderedItems.value = (data.items || []).map(item => ({
        ...item,
        newRating: item.Rating.Valid ? item.Rating.Int32 : null, // Pre-fill if already rated
        ratingLoading: false,
        ratingError: null,
        ratingSuccess: false
      }));
    } else {
      itemsError.value = data.message || 'Failed to fetch ordered items.';
      if (response.status === 401 || response.status === 403) {
        if (!error.value) {
            localStorage.removeItem('userToken');
            router.push('/login');
        }
      }
    }
  } catch (err) {
    itemsError.value = `Error fetching items: ${err.message}.`;
    console.error('Error fetching ordered items:', err);
  } finally {
    itemsLoading.value = false;
  }
};

// Function to submit overall order feedback
const submitOrderFeedback = async () => {
  feedbackLoading.value = true;
  feedbackSuccess.value = false;
  feedbackError.value = null;

  if (!userToken.value) {
    feedbackError.value = 'Authentication token is missing. Please log in.';
    feedbackLoading.value = false;
    return;
  }
  if (!orderId.value) {
    feedbackError.value = 'No Order ID available.';
    feedbackLoading.value = false;
    return;
  }
  if (!orderFeedback.value.trim()) {
    feedbackError.value = 'Feedback cannot be empty.';
    feedbackLoading.value = false;
    return;
  }

  try {
    const response = await fetch('http://localhost:8080/orders/feedback', {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${userToken.value}`,
      },
      body: JSON.stringify({
        order_id: parseInt(orderId.value), // Ensure it's an integer
        feedback: orderFeedback.value.trim(),
      }),
    });

    const data = await response.json();

    if (response.ok && data.success) {
      feedbackSuccess.value = true;
      // Optionally, update the order object's feedback locally
      if (order.value) {
        order.value.Feedback = { String: orderFeedback.value.trim(), Valid: true };
      }
      setTimeout(() => feedbackSuccess.value = false, 3000); // Hide success message after 3 seconds
    } else {
      feedbackError.value = data.message || 'Failed to submit feedback.';
      if (response.status === 401 || response.status === 403) {
        localStorage.removeItem('userToken');
        router.push('/login');
      }
    }
  } catch (err) {
    feedbackError.value = `Error submitting feedback: ${err.message}.`;
    console.error('Error submitting feedback:', err);
  } finally {
    feedbackLoading.value = false;
  }
};

// Function to submit rating for an individual item
const submitItemRating = async (item) => {
  item.ratingLoading = true;
  item.ratingSuccess = false;
  item.ratingError = null;

  if (!userToken.value) {
    item.ratingError = 'Authentication token is missing. Please log in.';
    item.ratingLoading = false;
    return;
  }
  if (!item.newRating || item.newRating < 1 || item.newRating > 5) {
    item.ratingError = 'Please select a rating between 1 and 5.';
    item.ratingLoading = false;
    return;
  }

  try {
    const response = await fetch('http://localhost:8080/orders/rate', {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${userToken.value}`,
      },
      body: JSON.stringify({
        order_id: parseInt(orderId.value),
        food_id: item.FoodID,
        rating: item.newRating,
      }),
    });

    const data = await response.json();

    if (response.ok && data.success) {
      item.ratingSuccess = true;
      item.Rating = { Int32: item.newRating, Valid: true }; // Update locally
      setTimeout(() => item.ratingSuccess = false, 3000);
    } else {
      item.ratingError = data.message || 'Failed to submit rating.';
      if (response.status === 401 || response.status === 403) {
        localStorage.removeItem('userToken');
        router.push('/login');
      }
    }
  } catch (err) {
    item.ratingError = `Error submitting rating: ${err.message}.`;
    console.error('Error submitting item rating:', err);
  } finally {
    item.ratingLoading = false;
  }
};

// On component mount, get orderId from route and fetch data
onMounted(async () => {
  orderId.value = route.params.id;

  if (orderId.value) {
    await fetchOrderDetail();
    await fetchOrderedItems();
  } else {
    loading.value = false;
    itemsLoading.value = false;
    error.value = 'No order ID provided in the URL.';
  }
});
</script>

<style scoped>
/* Optional: Add some star rating styles if you prefer stars over a select dropdown later */
/* Example for disabled select to hint it's not changeable */
select:disabled {
  background-color: #f0f0f0;
  cursor: not-allowed;
}
</style>