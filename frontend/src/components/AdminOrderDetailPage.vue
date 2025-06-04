<template>
  <div class="min-h-screen bg-gray-100 p-6 pt-10">
    <div class="max-w-4xl mx-auto bg-white rounded-lg shadow-xl p-8">
      <h2 class="text-4xl font-extrabold text-gray-900 mb-8 text-center">
        Admin Order Details
      </h2>

      <div v-if="loading" class="flex justify-center items-center h-48">
        <svg class="animate-spin h-10 w-10 text-indigo-500" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
        </svg>
        <p class="ml-3 text-lg text-gray-700">Loading admin order details...</p>
      </div>

      <div v-else-if="error" class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative mb-6" role="alert">
        <strong class="font-bold">Error!</strong>
        <span class="block sm:inline ml-2">{{ error }}</span>
      </div>

      <div v-else-if="!order.OrderID" class="bg-blue-100 border border-blue-400 text-blue-700 px-4 py-3 rounded relative text-center mb-6" role="alert">
        <strong class="font-bold">Order Not Found.</strong>
        <span class="block sm:inline ml-2">The requested order could not be loaded.</span>
      </div>

      <div v-else>
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4 mb-8 p-6 bg-gray-50 rounded-lg shadow-inner">
          <div>
            <p class="text-gray-700 text-lg mb-2"><span class="font-semibold">Order ID:</span> {{ order.OrderID }}</p>
            <p class="text-gray-700 text-lg mb-2"><span class="font-semibold">User ID:</span> {{ order.UserID }}</p>
            <p class="text-gray-700 text-lg mb-2"><span class="font-semibold">Order Time:</span> {{ formatDateTime(order.OrderTime.Valid ? order.OrderTime.Time : null) }}</p>
            <p class="text-gray-700 text-lg mb-2"><span class="font-semibold">Estimated Time:</span> {{ formatDateTime(order.EstimatedTime) }}</p>
            <p class="text-gray-700 text-lg mb-2"><span class="font-semibold">Order Info:</span> {{ order.OrderInfo }}</p>
            <p v-if="order.Feedback.Valid" class="text-gray-700 text-lg mb-2"><span class="font-semibold">Feedback:</span> {{ order.Feedback.String }}</p>
          </div>
          <div>
            <p class="text-gray-700 text-lg mb-2"><span class="font-semibold">Is Done:</span>
              <span :class="{'bg-green-100 text-green-800': order.IsDone, 'bg-yellow-100 text-yellow-800': !order.IsDone}" class="px-2 inline-flex text-sm leading-5 font-semibold rounded-full">
                {{ order.IsDone ? 'Yes' : 'No' }}
              </span>
            </p>
            <p class="text-gray-700 text-lg mb-2"><span class="font-semibold">Is Paid:</span>
              <span :class="{'bg-green-100 text-green-800': order.IsPaid, 'bg-red-100 text-red-800': !order.IsPaid}" class="px-2 inline-flex text-sm leading-5 font-semibold rounded-full">
                {{ order.IsPaid ? 'Yes' : 'No' }}
              </span>
            </p>
            <p class="text-gray-700 text-lg mb-2"><span class="font-semibold">Is Ranged:</span>
              <span :class="{'bg-green-100 text-green-800': order.IsRanged, 'bg-gray-100 text-gray-800': !order.IsRanged}" class="px-2 inline-flex text-sm leading-5 font-semibold rounded-full">
                {{ order.IsRanged ? 'Yes' : 'No' }}
              </span>
            </p>
            <p v-if="order.IsRanged && order.DeliveryAddress && order.DeliveryAddress.Valid" class="text-gray-700 text-lg mb-2">
                <span class="font-semibold">Delivery Address:</span> {{ order.DeliveryAddress.String }}
            </p>
            <p class="text-gray-700 text-lg mb-2"><span class="font-semibold">Deleted:</span>
              <span :class="{'bg-red-100 text-red-800': order.Deleted, 'bg-gray-100 text-gray-800': !order.Deleted}" class="px-2 inline-flex text-sm leading-5 font-semibold rounded-full">
                {{ order.Deleted ? 'Yes' : 'No' }}
              </span>
            </p>
          </div>
        </div>

        <h3 class="text-3xl font-bold text-gray-800 mb-6 text-center">Ordered Items</h3>
        <div v-if="itemsLoading" class="flex justify-center items-center h-24">
          <svg class="animate-spin h-8 w-8 text-indigo-500" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
          <p class="ml-2 text-md text-gray-600">Loading items...</p>
        </div>
        <div v-else-if="itemsError" class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative mb-6" role="alert">
          <strong class="font-bold">Error loading items!</strong>
          <span class="block sm:inline ml-2">{{ itemsError }}</span>
        </div>
        <div v-else-if="items.length === 0" class="bg-yellow-100 border border-yellow-400 text-yellow-700 px-4 py-3 rounded relative text-center">
          <strong class="font-bold">No Items Found!</strong>
          <span class="block sm:inline ml-2">This order contains no items.</span>
        </div>
        <div v-else>
          <div v-for="item in items" :key="item.ItemID" class="bg-gray-50 p-4 rounded-lg shadow-sm mb-3">
            <div class="flex justify-between items-center">
              <p class="text-gray-800 font-semibold">Item ID: {{ item.ItemID }}</p>
              <p class="text-gray-800 font-semibold">Food ID: {{ item.FoodID }}</p>
              <p class="text-gray-800">Quantity: {{ item.Quantity }}</p>
            </div>
            <p v-if="item.Rating.Valid" class="text-gray-600 text-sm mt-1">Rating: {{ item.Rating.Int32 }} stars</p>
            <p v-else class="text-gray-600 text-sm mt-1">Rating: Not yet rated</p>
          </div>
        </div>

        <div class="mt-8 text-center">
          <button @click="router.back()" class="bg-gray-300 hover:bg-gray-400 text-gray-800 font-bold py-2 px-6 rounded-lg transition duration-200">
            Go Back
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useUserStore } from '../stores/user';

const route = useRoute(); // Access route parameters
const router = useRouter();
const userStore = useUserStore();

const order = ref({}); // Use an empty object for a single order
const items = ref([]); // Use an empty array for order items
const loading = ref(true);
const error = ref(null);
const itemsLoading = ref(true);
const itemsError = ref(null);

const orderId = route.params.id; // Get the order ID from the route parameters

const fetchOrderDetails = async () => {
  console.log(`AdminOrderDetailPage: Fetching details for Order ID: ${orderId}`);
  loading.value = true;
  error.value = null;

  // Admin page, so check admin access
  if (!userStore.isLoggedIn || !userStore.isAdmin) {
    error.value = "Unauthorized access. You must be an administrator to view this page.";
    loading.value = false;
    router.push('/'); // Redirect to home or login page if not authorized
    return;
  }

  try {
    const response = await fetch(`http://localhost:8080/orders/info?order_id=${orderId}`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${userStore.userToken}`,
      },
    });

    const data = await response.json();

    if (response.ok && data.success) {
      order.value = data.order || {}; // Ensure it's an object, even if null
      console.log('AdminOrderDetailPage: Order details fetched:', order.value);
    } else {
      error.value = data.message || `Failed to fetch order details. Status: ${response.status}`;
      console.error('AdminOrderDetailPage: Error fetching order details:', error.value);
      order.value = {}; // Clear order if fetch failed
      // Handle 401/403 for redirection
      if (response.status === 401 || response.status === 403) {
        userStore.logout();
        router.push('/login');
      }
    }
  } catch (err) {
    error.value = `Network or parsing error fetching order details: ${err.message}`;
    console.error('AdminOrderDetailPage: Network error:', err);
    order.value = {};
  } finally {
    loading.value = false;
  }
};

const fetchOrderItems = async () => {
  console.log(`AdminOrderDetailPage: Fetching items for Order ID: ${orderId}`);
  itemsLoading.value = true;
  itemsError.value = null;

  if (!userStore.isLoggedIn || !userStore.isAdmin) {
    itemsError.value = "Unauthorized access to fetch items.";
    itemsLoading.value = false;
    return;
  }

  try {
    const response = await fetch(`http://localhost:8080/orders/items?order_id=${orderId}`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${userStore.userToken}`,
      },
    });

    const data = await response.json();

    if (response.ok && data.success) {
      items.value = data.items || []; // Ensure it's an array
      console.log('AdminOrderDetailPage: Order items fetched:', items.value);
    } else {
      itemsError.value = data.message || `Failed to fetch order items. Status: ${response.status}`;
      console.error('AdminOrderDetailPage: Error fetching order items:', itemsError.value);
      items.value = [];
       // Handle 401/403 for redirection
      if (response.status === 401 || response.status === 403) {
        userStore.logout();
        router.push('/login');
      }
    }
  } catch (err) {
    itemsError.value = `Network or parsing error fetching order items: ${err.message}`;
    console.error('AdminOrderDetailPage: Network error items:', err);
    items.value = [];
  } finally {
    itemsLoading.value = false;
  }
};


// Helper function to format date/time, handling sql.NullTime
const formatDateTime = (dateValue) => {
  if (!dateValue) return 'N/A';
  try {
    const date = new Date(dateValue);
    // Use toLocaleString for a user-friendly format, e.g., "6/2/2025, 9:13:18 PM"
    // Setting timezone to Singapore for consistency as per previous instructions
    return date.toLocaleString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit',
      hour12: true,
      timeZone: 'Asia/Singapore'
    });
  } catch (e) {
    console.error("Error formatting date:", e);
    return String(dateValue); // Return original string/value if parsing fails
  }
};

onMounted(() => {
  if (orderId) {
    fetchOrderDetails();
    fetchOrderItems();
  } else {
    error.value = "No Order ID provided in the URL.";
    loading.value = false;
    itemsLoading.value = false;
  }
});
</script>

<style scoped>
/* Add any specific styles for this page if needed */
</style>