<template>
  <div class="min-h-screen bg-gray-100 p-6 pt-10">
    <div class="max-w-7xl mx-auto bg-white rounded-lg shadow-xl p-8">
      <h2 class="text-4xl font-extrabold text-gray-900 mb-8 text-center">
        Admin Panel: Order History
      </h2>

      <div v-if="loading" class="flex justify-center items-center h-48">
        <svg class="animate-spin h-10 w-10 text-indigo-500" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
        </svg>
        <p class="ml-3 text-lg text-gray-700">Loading order history...</p>
      </div>

      <div v-else-if="error" class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative" role="alert">
        <strong class="font-bold">Error!</strong>
        <span class="block sm:inline ml-2">{{ error }}</span>
      </div>

      <div v-else-if="orders.length === 0" class="bg-blue-100 border border-blue-400 text-blue-700 px-4 py-3 rounded relative text-center" role="alert">
        <strong class="font-bold">No Orders Found!</strong>
        <span class="block sm:inline ml-2">There is no order history to display.</span>
      </div>

      <div v-else class="overflow-x-auto">
        <table class="min-w-full bg-white border border-gray-200 rounded-lg shadow-md">
          <thead class="bg-gray-200">
            <tr>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-700 uppercase tracking-wider">Order ID</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-700 uppercase tracking-wider">User ID</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-700 uppercase tracking-wider">Order Time</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-700 uppercase tracking-wider">Order Info</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-700 uppercase tracking-wider">Is Done</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-700 uppercase tracking-wider">Is Paid</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-700 uppercase tracking-wider">Deleted</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-700 uppercase tracking-wider">Actions</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-200">
            <tr v-for="order in orders" :key="order.OrderID" :class="{'bg-red-50': order.Deleted}">
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{{ order.OrderID }}</td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{{ order.UserID }}</td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                 {{ formatDateTime(order.OrderTime.Valid ? order.OrderTime.Time : null) }}
              </td>
              <td class="px-6 py-4 text-sm text-gray-900 max-w-xs truncate" :title="order.OrderInfo">
                {{ order.OrderInfo }}
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm">
                <span :class="{'bg-green-100 text-green-800': order.IsDone, 'bg-yellow-100 text-yellow-800': !order.IsDone}"
                      class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full">
                  {{ order.IsDone ? 'Yes' : 'No' }}
                </span>
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm">
                <span :class="{'bg-green-100 text-green-800': order.IsPaid, 'bg-red-100 text-red-800': !order.IsPaid}"
                      class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full">
                  {{ order.IsPaid ? 'Yes' : 'No' }}
                </span>
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm">
                <span :class="{'bg-red-100 text-red-800': order.Deleted, 'bg-gray-100 text-gray-800': !order.Deleted}"
                      class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full">
                  {{ order.Deleted ? 'Yes' : 'No' }}
                </span>
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm font-medium">
                <router-link
                  :to="`/admin/orders/${order.OrderID}`" class="text-blue-600 hover:text-blue-900"
                  title="View Details"
                >
                  Details
                </router-link>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useUserStore } from '../stores/user';
import { useRouter } from 'vue-router';

const userStore = useUserStore();
const router = useRouter();

const orders = ref([]);
const loading = ref(true);
const error = ref(null);

const fetchAllOrders = async () => {
  console.log("AdminHistoryPage: fetchAllOrders called.");
  loading.value = true;
  error.value = null;

  if (!userStore.isLoggedIn || !userStore.isAdmin) {
    error.value = "Unauthorized access. You must be an administrator to view this page.";
    loading.value = false;
    console.warn("AdminHistoryPage: Access denied (client-side check). Redirecting to home.");
    router.push('/');
    return;
  }

  try {
    console.log("AdminHistoryPage: Attempting to fetch from:", 'http://localhost:8080/admin/orders-all');
    const response = await fetch('http://localhost:8080/admin/orders-all', {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${userStore.userToken}`,
      },
    });

    console.log("AdminHistoryPage: Fetch response received. Status:", response.status);

    const data = await response.json();
    console.log("AdminHistoryPage: Response data:", data);

    if (response.ok && data.success) {
      orders.value = data.orders || [];
      console.log("AdminHistoryPage: All orders fetched successfully. Count:", orders.value.length);
    } else {
      error.value = data.message || `Failed to fetch all orders. HTTP Status: ${response.status}`;
      console.error("AdminHistoryPage: API error:", error.value);

      if (response.status === 401 || response.status === 403) {
         console.warn("AdminHistoryPage: Unauthorized/Forbidden. Logging out and redirecting.");
         userStore.logout();
         router.push('/login');
      }
    }
  } catch (err) {
    error.value = `Error fetching all orders: ${err.message}. Please check your network connection and backend server.`;
    console.error('AdminHistoryPage: Network or parsing error:', err);
  } finally {
    loading.value = false;
    console.log("AdminHistoryPage: fetchAllOrders finished. loading set to false.");
  }
};

const formatDateTime = (dateValue) => {
  if (!dateValue) return 'N/A';
  try {
    const date = new Date(dateValue);
    return date.toLocaleString();
  } catch (e) {
    console.error("Error formatting date:", e);
    return String(dateValue);
  }
};


onMounted(() => {
  console.log("AdminHistoryPage: Component mounted.");
  fetchAllOrders();
});
</script>

<style scoped>
/* Optional: Highlight deleted rows */
.bg-red-50 {
  background-color: #fef2f2; /* Light red background for deleted orders */
}
</style>