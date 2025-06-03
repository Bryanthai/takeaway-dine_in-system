<template>
  <div class="min-h-screen bg-gray-100 p-6 pt-10">
    <div class="max-w-7xl mx-auto bg-white rounded-lg shadow-xl p-8">
      <h2 class="text-4xl font-extrabold text-gray-900 mb-8 text-center">
        Admin Panel: Undone Orders
      </h2>

      <section class="mb-12 p-6 bg-indigo-50 rounded-lg shadow-inner text-center">
        <h3 class="text-2xl font-semibold text-indigo-800 mb-4">Overall Average Spending Per Order</h3>
        <div v-if="totalAverageLoading" class="flex justify-center items-center">
          <svg class="animate-spin h-8 w-8 text-indigo-600" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
          <p class="ml-3 text-lg text-indigo-700">Calculating overall average...</p>
        </div>
        <div v-else-if="totalAverageError" class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative">
          <strong class="font-bold">Error!</strong>
          <span class="block sm:inline ml-2">{{ totalAverageError }}</span>
        </div>
        <div v-else>
          <p class="text-5xl font-extrabold text-indigo-900">
            ${{ totalAverageSpending !== null ? totalAverageSpending.toFixed(2) : '0.00' }}
          </p>
          <p class="text-md text-indigo-600 mt-2">Average amount spent by a person on a single order.</p>
        </div>
      </section>


      <section>
        <h3 class="text-3xl font-bold text-gray-800 mb-6 text-center">Current Undone Orders</h3>
        <div v-if="loading" class="flex justify-center items-center h-48">
          <svg class="animate-spin h-10 w-10 text-indigo-500" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
          <p class="ml-3 text-lg text-gray-700">Loading undone orders...</p>
        </div>

        <div v-else-if="error" class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative" role="alert">
          <strong class="font-bold">Error!</strong>
          <span class="block sm:inline ml-2">{{ error }}</span>
        </div>

        <div v-else-if="orders.length === 0" class="bg-blue-100 border border-blue-400 text-blue-700 px-4 py-3 rounded relative text-center" role="alert">
          <strong class="font-bold">No Undone Orders!</strong>
          <span class="block sm:inline ml-2">All orders are currently completed or no new orders have been placed.</span>
        </div>

        <div v-else class="overflow-x-auto">
          <table class="min-w-full bg-white border border-gray-200 rounded-lg shadow-md">
            <thead class="bg-gray-200">
              <tr>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-700 uppercase tracking-wider">Order ID</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-700 uppercase tracking-wider">User ID</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-700 uppercase tracking-wider">Total Price</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-700 uppercase tracking-wider">Status</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-700 uppercase tracking-wider">Ordered At</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-700 uppercase tracking-wider">Actions</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-200">
              <tr v-for="order in orders" :key="order.OrderID">
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{{ order.OrderID }}</td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{{ order.UserID }}</td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">${{ order.TotalPrice ? order.TotalPrice.toFixed(2) : '0.00' }}</td>
                <td class="px-6 py-4 whitespace-nowrap text-sm">
                  <span :class="getStatusBadgeClass(order.Status)" class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full capitalize">
                    {{ order.Status }}
                  </span>
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                  {{ formatDateTime(order.CreatedAt) }}
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm font-medium">
                  <button
                    @click="finishOrder(order.OrderID)"
                    class="text-green-600 hover:text-green-900 mr-4"
                    title="Mark as Done"
                  >
                    Finish
                  </button>
                  <button
                    @click="deleteOrder(order.OrderID)"
                    class="text-red-600 hover:text-red-900 mr-4"
                    title="Delete Order"
                  >
                    Delete
                  </button>
                  <router-link
                    :to="`/admin/orders/${order.OrderID}`"
                    class="text-blue-600 hover:text-blue-900"
                    title="View Details"
                  >
                    Details
                  </router-link>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </section>
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
const loading = ref(true); // Loading state for undone orders
const error = ref(null); // Error state for undone orders

// FIX: Initialize totalAverageSpending to a number (0) to prevent toFixed errors
const totalAverageSpending = ref(0);
const totalAverageLoading = ref(true); // NEW: Loading state for total average
const totalAverageError = ref(null); // NEW: Error state for total average

const fetchUndoneOrders = async () => {
  console.log("AdminOrdersPage: fetchUndoneOrders called.");
  loading.value = true; // Ensure loading is true at the start of fetch
  error.value = null;

  if (!userStore.isLoggedIn || !userStore.isAdmin) {
    error.value = "Unauthorized access. You must be an administrator to view this page.";
    loading.value = false; // Set loading to false on immediate redirect
    console.warn("AdminOrdersPage: Access denied (client-side check). Redirecting to home.");
    router.push('/');
    return;
  }

  try {
    console.log("AdminOrdersPage: Attempting to fetch undone orders from:", 'http://localhost:8080/admin/undone-orders');
    const response = await fetch('http://localhost:8080/admin/undone-orders', {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${userStore.userToken}`,
      },
    });

    console.log("AdminOrdersPage: Fetch undone orders response received. Status:", response.status);

    if (!response.ok) { // Handle non-2xx responses
      const errorData = await response.json().catch(() => ({})); // Try to parse error data
      error.value = errorData.message || `Failed to fetch undone orders. HTTP Status: ${response.status}`;
      console.error("AdminOrdersPage: API error fetching undone orders:", error.value, errorData);

      if (response.status === 401 || response.status === 403) {
           console.warn("AdminOrdersPage: Unauthorized/Forbidden for undone orders. Logging out and redirecting.");
           userStore.logout();
           router.push('/login');
       }
       return; // Exit here if response was not ok
    }

    const data = await response.json();
    console.log("AdminOrdersPage: Raw response data for undone orders:", data);

    if (data.success) {
      orders.value = data.orders || [];
      console.log("AdminOrdersPage: Undone orders fetched successfully. Count:", orders.value.length);
      if (orders.value.length === 0) {
        console.log("AdminOrdersPage: No undone orders found.");
      }
    } else {
      error.value = data.message || 'Failed to fetch undone orders (backend reported failure).';
      console.error('AdminOrdersPage: Backend reported failure for undone orders:', data.message);
    }
  } catch (err) {
    error.value = `Error fetching undone orders: ${err.message}. Please check your network connection and backend server.`;
    console.error('AdminOrdersPage: Network or parsing error fetching undone orders:', err);
  } finally {
    loading.value = false; // Always set loading to false in finally
    console.log("AdminOrdersPage: fetchUndoneOrders finished. loading set to false.");
  }
};

// Function to fetch total average spending
const fetchTotalAverageSpending = async () => {
  console.log("AdminOrdersPage: fetchTotalAverageSpending called.");
  totalAverageLoading.value = true;
  totalAverageError.value = null;

  if (!userStore.isLoggedIn || !userStore.isAdmin) {
    totalAverageError.value = "Unauthorized to fetch total average spending.";
    totalAverageLoading.value = false;
    return;
  }

  try {
    console.log("AdminOrdersPage: Attempting to fetch total average from:", 'http://localhost:8080/admin/total-average');
    const response = await fetch('http://localhost:8080/admin/total-average', {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${userStore.userToken}`,
      },
    });

    console.log("AdminOrdersPage: Fetch total average response received. Status:", response.status);

    if (!response.ok) { // Handle non-2xx responses
      const errorData = await response.json().catch(() => ({}));
      totalAverageError.value = errorData.message || `Failed to fetch total average spending. Status: ${response.status}`;
      console.error("AdminOrdersPage: API error fetching total average spending:", totalAverageError.value, errorData);
      if (response.status === 401 || response.status === 403) {
        console.warn("AdminOrdersPage: Unauthorized/Forbidden for total average. Logging out and redirecting.");
        userStore.logout();
        router.push('/login');
      }
      return;
    }

    const data = await response.json();
    console.log("AdminOrdersPage: Raw response data for total average:", data);

    if (data.success) {
      // FIX: Ensure data.Average is a number, default to 0 if not
      totalAverageSpending.value = typeof data.Average === 'number' ? data.Average : 0;
      console.log("AdminOrdersPage: Total average spending fetched:", totalAverageSpending.value);
    } else {
      totalAverageError.value = data.message || 'Failed to fetch total average spending (backend reported failure).';
      console.error("AdminOrdersPage: Backend reported failure for total average spending:", data.message);
      // FIX: Set totalAverageSpending to 0 in case of backend failure
      totalAverageSpending.value = 0;
    }
  } catch (err) {
    totalAverageError.value = `Error fetching total average spending: ${err.message}.`;
    console.error('AdminOrdersPage: Network error fetching total average spending:', err);
    // FIX: Set totalAverageSpending to 0 in case of network/parsing error
    totalAverageSpending.value = 0;
  } finally {
    totalAverageLoading.value = false; // Always set loading to false in finally
    console.log("AdminOrdersPage: fetchTotalAverageSpending finished. loading set to false.");
  }
};


const finishOrder = async (orderId) => {
  if (!confirm(`Are you sure you want to FINISH Order ID ${orderId}? This cannot be undone.`)) {
    return;
  }

  try {
    const response = await fetch('http://localhost:8080/orders/finish', {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${userStore.userToken}`,
      },
      body: JSON.stringify({ order_id: orderId }),
    });

    const data = await response.json();

    if (response.ok && data.success) {
      alert(`Order ${orderId} finished successfully!`);
      fetchUndoneOrders(); // Re-fetch undone orders
      fetchTotalAverageSpending(); // Re-fetch total average as it might change
    } else {
      alert(`Failed to finish order ${orderId}: ${data.message || 'Unknown error.'}`);
      console.error(`Failed to finish order ${orderId}:`, data);
    }
  } catch (err) {
    alert(`Error finishing order ${orderId}: ${err.message}`);
    console.error('Network error finishing order:', err);
  }
};

const deleteOrder = async (orderId) => {
  if (!confirm(`Are you sure you want to DELETE Order ID ${orderId}? This will mark it as deleted.`)) {
    return;
  }

  try {
    const response = await fetch('http://localhost:8080/orders', {
      method: 'DELETE',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${userStore.userToken}`,
      },
      body: JSON.stringify({ order_id: orderId }),
    });

    const data = await response.json();

    if (response.ok && data.success) {
      alert(`Order ${orderId} marked as deleted successfully!`);
      fetchUndoneOrders(); // Re-fetch undone orders
      fetchTotalAverageSpending(); // Re-fetch total average as it might change
    } else {
      alert(`Failed to delete order ${orderId}: ${data.message || 'Unknown error.'}`);
      console.error(`Failed to delete order ${orderId}:`, data);
    }
  } catch (err) {
    alert(`Error deleting order ${orderId}: ${err.message}`);
    console.error('Network error deleting order:', err);
  }
};

// FIX: Added a null/undefined check for 'status'
const getStatusBadgeClass = (status) => {
  // Default to an empty string if status is null or undefined to prevent .toLowerCase() error
  const safeStatus = (status || '').toLowerCase(); 
  switch (safeStatus) {
    case 'pending': return 'bg-yellow-100 text-yellow-800';
    case 'processing': return 'bg-blue-100 text-blue-800';
    case 'ready': return 'bg-purple-100 text-purple-800';
    case 'shipped': return 'bg-indigo-100 text-indigo-800';
    case 'delivered': return 'bg-green-100 text-green-800';
    case 'cancelled': return 'bg-red-100 text-red-800';
    default: return 'bg-gray-100 text-gray-800'; // Fallback for 'unknown' or other unexpected values
  }
};

const formatDateTime = (isoString) => {
  if (!isoString) return 'N/A';
  try {
    const date = new Date(isoString);
    return date.toLocaleString();
  } catch (e) {
    console.error("Error formatting date:", e);
    return isoString;
  }
};


onMounted(() => {
  console.log("AdminOrdersPage: Component mounted.");
  fetchUndoneOrders();
  fetchTotalAverageSpending();
});
</script>

<style scoped>
/* Scoped styles specific to this component */
</style>