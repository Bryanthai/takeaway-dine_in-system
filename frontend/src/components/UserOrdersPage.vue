<template>
  <div class="min-h-screen bg-gray-100 p-6">
    <h2 class="text-4xl font-extrabold text-gray-900 mb-8 text-center">
      My Orders
    </h2>

    <div v-if="loading" class="flex justify-center items-center h-48">
      <svg class="animate-spin h-10 w-10 text-indigo-500" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
        <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
      </svg>
      <p class="ml-3 text-lg text-gray-700">Loading your orders...</p>
    </div>

    <div v-else-if="error" class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative" role="alert">
      <strong class="font-bold">Error!</strong>
      <span class="block sm:inline ml-2">{{ error }}</span>
    </div>

    <div v-else-if="orders.length === 0" class="bg-blue-100 border border-blue-400 text-blue-700 px-4 py-3 rounded relative text-center" role="alert">
      <strong class="font-bold">No Orders Found!</strong>
      <span class="block sm:inline ml-2">You haven't placed any orders yet.</span>
    </div>

    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <router-link
        v-for="order in orders"
        :key="order.OrderID"
        :to="{ name: 'OrderDetail', params: { id: order.OrderID }}"
        class="group block bg-white rounded-lg shadow-md p-6 border-l-4 cursor-pointer transform transition-transform duration-200 hover:scale-[1.02] hover:shadow-xl relative"
        :class="{
          'border-green-500': order.IsDone && !order.Deleted,
          'border-yellow-500': !order.IsDone && !order.IsPaid && !order.Deleted,
          'border-indigo-500': !order.IsDone && order.IsPaid && !order.Deleted,
          'border-red-500 bg-red-50 ring-2 ring-red-300': !order.IsDone && order.Deleted
        }"
      >
        <h3 class="text-xl font-semibold text-gray-800 mb-2">Order #{{ order.OrderID }}</h3>
        <p class="text-gray-600 mb-1">
          <span class="font-medium">Status:</span>
          <span :class="{
            'text-green-600 font-bold': order.IsDone && !order.Deleted,
            'text-yellow-600 font-bold': !order.IsDone && !order.IsPaid && !order.Deleted,
            'text-indigo-600 font-bold': !order.IsDone && order.IsPaid && !order.Deleted,
            'text-red-600 font-bold': !order.IsDone && order.Deleted
          }">
            {{ order.IsDone ? 'Completed' : (order.IsPaid ? 'Processing' : 'Pending Payment') }}
            <span v-if="!order.IsDone && order.Deleted" class="text-red-700 ml-1">(Cancelled by Admin)</span>
          </span>
        </p>
        <p class="text-gray-600 mb-1">
            <span class="font-medium">Order Placed:</span>
            {{ formatDateTime(order.OrderTime.String) }}
        </p>
        <p class="text-gray-600 mb-1">
            <span class="font-medium">Estimated Completion:</span>
            {{ formatDateTime(order.EstimatedTime) }}
        </p>
        <p class="text-gray-600 mb-1">
            <span class="font-medium">Delivery:</span>
            {{ order.IsRanged ? 'Delivery' : 'Pickup' }}
        </p>
        <p class="text-gray-600 mb-4">
            <span class="font-medium">Paid:</span>
            <span :class="{'text-green-600 font-bold': order.IsPaid, 'text-red-600 font-bold': !order.IsPaid}">
                {{ order.IsPaid ? 'Yes' : 'No' }}
            </span>
        </p>
        <p class="text-gray-600 mb-4">
            <span class="font-medium">Deleted (Admin):</span>
            <span :class="{'text-red-600 font-bold': order.Deleted, 'text-green-600 font-bold': !order.Deleted}">
                {{ order.Deleted ? 'Yes' : 'No' }}
            </span>
        </p>

        <div class="border-t border-gray-200 pt-4">
          <p class="text-gray-700 font-medium mb-2">Order Details (Summary):</p>
          <ul class="list-disc list-inside text-gray-700">
            <li v-for="(item, index) in parseOrderInfo(order.OrderInfo)" :key="index">
              {{ item.FoodName }} (x{{ item.Quantity }})
            </li>
          </ul>
        </div>

        <button
            v-if="!order.IsPaid && !order.Deleted"
            @click.prevent="goToPayment(order.OrderID)"
            class="mt-4 w-full bg-blue-600 text-white px-4 py-2 rounded-lg font-semibold hover:bg-blue-700 transition duration-200 flex items-center justify-center text-lg z-10"
        >
            <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 mr-2" viewBox="0 0 24 24" fill="currentColor">
                <path d="M12 2c-.55 0-1 .45-1 1v2c0 .55.45 1 1 1s1-.45 1-1V3c0-.55-.45-1-1-1zm0 20c.55 0 1-.45 1-1v-2c0-.55-.45-1-1-1s-1 .45-1 1v2c0 .55.45 1 1 1zM4 12c0-.55.45-1 1-1h2c.55 0 1 .45 1 1s-.45 1-1 1H5c-.55 0-1-.45-1-1zm16 0c0-.55-.45-1-1-1h-2c-.55 0-1 .45-1 1s.45 1 1 1h2c.55 0 1-.45 1-1zM18.36 5.64c-.39-.39-1.02-.39-1.41 0l-1.42 1.42c-.39.39-.39 1.02 0 1.41.39.39 1.02.39 1.41 0l1.42-1.42c.38-.39.38-1.02-.01-1.41zM7.05 16.95c.39.39 1.02.39 1.41 0l1.42-1.42c.39-.39.39-1.02 0-1.41-.39-.39-1.02-.39-1.41 0l-1.42 1.42c-.39.39-.39 1.02 0 1.41zM5.64 18.36c.39.39 1.02.39 1.41 0l1.42-1.42c.39-.39.39-1.02 0-1.41-.39-.39-1.02-.39-1.41 0l-1.42 1.42c-.39.39-.39 1.02 0 1.41zM16.95 7.05c.39-.39 1.02-.39 1.41 0l1.42 1.42c.39.39.39 1.02 0 1.41-.39.39-1.02.39-1.41 0l-1.42-1.42c-.39-.39-.39-1.02 0-1.41z"/>
            </svg>
            Pay Now
        </button>
      </router-link>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';

const router = useRouter();

const orders = ref([]);
const loading = ref(true);
const error = ref(null);
const userToken = ref(localStorage.getItem('userToken') || '');

const fetchUserOrders = async () => {
  loading.value = true;
  error.value = null;

  if (!userToken.value) {
    error.value = 'Authentication token is missing. Please log in to view your orders.';
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
      orders.value = data.orders || [];
    } else {
      error.value = data.message || 'Failed to fetch orders.';
      if (response.status === 401 || response.status === 403) {
        error.value = "Session expired or unauthorized. Please log in again.";
        localStorage.removeItem('userToken');
        router.push('/login');
      }
    }
  } catch (err) {
    error.value = `An error occurred: ${err.message}. Please check your network connection.`;
    console.error('Error fetching orders:', err);
  } finally {
    loading.value = false;
  }
};

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
    return [{ FoodName: "Invalid/Unreadable Order Info", Quantity: 1, Price: 0 }];
  }
};

// New method to handle payment redirection
const goToPayment = (orderId) => {
  // We will define this route and component later.
  // For now, it just navigates to a placeholder.
  router.push({ name: 'PaymentPage', params: { orderId: orderId } });
};

onMounted(() => {
  fetchUserOrders();
});
</script>

<style scoped>
/* Add any specific styles here if needed */
</style>