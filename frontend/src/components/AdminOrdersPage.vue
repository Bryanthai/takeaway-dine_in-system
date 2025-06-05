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
            ${{ typeof totalAverageSpending === 'number' ? totalAverageSpending.toFixed(2) : '0.00' }}
          </p>
          <p class="text-md text-indigo-600 mt-2">Average amount spent by a person on a single order.</p>
        </div>
      </section>

      <section>
        <h3 class="text-3xl font-bold text-gray-800 mb-6 text-center">Current Undone Orders</h3>

        <div class="mb-6 flex justify-end items-center">
          <label for="orderFilter" class="mr-3 text-gray-700 font-medium">Show Orders:</label>
          <select id="orderFilter" v-model="filterType"
                  class="block w-48 pl-3 pr-10 py-2 text-base border-gray-300 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm rounded-md shadow-sm">
            <option value="all">All Orders</option>
            <option value="ranged">Delivery Orders</option>
            <option value="non-ranged">Pickup Orders</option>
          </select>
        </div>

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

        <div v-else-if="filteredOrders.length === 0" class="bg-blue-100 border border-blue-400 text-blue-700 px-4 py-3 rounded relative text-center" role="alert">
          <strong class="font-bold">No Undone Orders Found!</strong>
          <span v-if="filterType === 'all'" class="block sm:inline ml-2">All orders are currently completed or no new orders have been placed.</span>
          <span v-else-if="filterType === 'ranged'" class="block sm:inline ml-2">No undone delivery orders found.</span>
          <span v-else-if="filterType === 'non-ranged'" class="block sm:inline ml-2">No undone pickup orders found.</span>
        </div>

        <div v-else class="overflow-x-auto">
          <table class="min-w-full bg-white border border-gray-200 rounded-lg shadow-md">
            <thead class="bg-gray-200">
              <tr>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-700 uppercase tracking-wider cursor-pointer hover:bg-gray-300 transition-colors duration-200"
                    :class="{ 'bg-gray-300 text-gray-900': sortColumn === 'OrderID' }"
                    @click="sortBy('OrderID')">
                  <div class="flex items-center">
                    Order ID
                    <span v-if="sortColumn === 'OrderID'" class="ml-1">
                      <span v-if="sortDirection === 'asc'">&#9650;</span>
                      <span v-else>&#9660;</span>
                    </span>
                  </div>
                </th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-700 uppercase tracking-wider">User ID</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-700 uppercase tracking-wider cursor-pointer hover:bg-gray-300 transition-colors duration-200"
                    :class="{ 'bg-gray-300 text-gray-900': sortColumn === 'TotalPrice' }"
                    @click="sortBy('TotalPrice')">
                  <div class="flex items-center">
                    Total Price
                    <span v-if="sortColumn === 'TotalPrice'" class="ml-1">
                      <span v-if="sortDirection === 'asc'">&#9650;</span>
                      <span v-else>&#9660;</span>
                    </span>
                  </div>
                </th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-700 uppercase tracking-wider cursor-pointer hover:bg-gray-300 transition-colors duration-200"
                    :class="{ 'bg-gray-300 text-gray-900': sortColumn === 'OrderTime' }"
                    @click="sortBy('OrderTime')">
                  <div class="flex items-center">
                    Ordered At
                    <span v-if="sortColumn === 'OrderTime'" class="ml-1">
                      <span v-if="sortDirection === 'asc'">&#9650;</span>
                      <span v-else>&#9660;</span>
                    </span>
                  </div>
                </th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-700 uppercase tracking-wider">Is Ranged</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-700 uppercase tracking-wider">Actions</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-200">
              <tr v-for="order in filteredOrders" :key="order.OrderID">
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{{ order.OrderID }}</td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{{ order.UserID }}</td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                  ${{ typeof order.TotalPrice === 'number' ? order.TotalPrice.toFixed(2) : '0.00' }}
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                  {{ formatDateTime(order.OrderTime) }}
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm">
                  <span :class="{'bg-blue-100 text-blue-800': order.IsRanged, 'bg-gray-100 text-gray-800': !order.IsRanged}"
                        class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full">
                    {{ order.IsRanged ? 'Delivery' : 'Pickup' }}
                  </span>
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
import { ref, onMounted, computed } from 'vue';
import { useUserStore } from '../stores/user';
import { useRouter } from 'vue-router';

const userStore = useUserStore();
const router = useRouter();

const orders = ref([]);
const loading = ref(true);
const error = ref(null);

const totalAverageSpending = ref(0);
const totalAverageLoading = ref(true);
const totalAverageError = ref(null);

// Sorting state
const sortColumn = ref('OrderID');
const sortDirection = ref('asc');

// Filtering state
const filterType = ref('all'); // 'all', 'ranged', 'non-ranged'

// Function to handle sorting
const sortBy = (column) => {
  if (sortColumn.value === column) {
    sortDirection.value = sortDirection.value === 'asc' ? 'desc' : 'asc';
  } else {
    sortColumn.value = column;
    sortDirection.value = 'asc';
  }
};

// Computed property for sorted orders
const sortedOrders = computed(() => {
  if (!orders.value || orders.value.length === 0) {
    return [];
  }

  return [...orders.value].sort((a, b) => {
    let aValue = a[sortColumn.value];
    let bValue = b[sortColumn.value];

    if (sortColumn.value === 'OrderTime') {
      const aTime = aValue && aValue.Valid ? new Date(aValue.Time).getTime() : 0;
      const bTime = bValue && bValue.Valid ? new Date(bValue.Time).getTime() : 0;

      if (aTime === 0 && bTime === 0) return 0;
      if (aTime === 0) return sortDirection.value === 'asc' ? 1 : -1;
      if (bTime === 0) return sortDirection.value === 'asc' ? -1 : 1;

      return sortDirection.value === 'asc' ? aTime - bTime : bTime - aTime;
    }

    if (aValue < bValue) {
      return sortDirection.value === 'asc' ? -1 : 1;
    }
    if (aValue > bValue) {
      return sortDirection.value === 'asc' ? 1 : -1;
    }
    return 0;
  });
});

// Computed property for filtered orders (applies filter after sorting)
const filteredOrders = computed(() => {
  const currentOrders = sortedOrders.value; // Start with the already sorted list

  if (filterType.value === 'all') {
    return currentOrders;
  } else if (filterType.value === 'ranged') {
    return currentOrders.filter(order => order.IsRanged);
  } else if (filterType.value === 'non-ranged') {
    return currentOrders.filter(order => !order.IsRanged);
  }
  return currentOrders; // Fallback to all orders if filterType is unexpected
});


const fetchOrderTotalPrice = async (orderId) => {
    try {
        const response = await fetch(`http://localhost:8080/orders/price?order_id=${orderId}`, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${userStore.userToken}`,
            },
        });

        if (!response.ok) {
            const errorData = await response.json().catch(() => ({}));
            console.error(`AdminOrdersPage: Failed to fetch total price for OrderID ${orderId}:`, errorData.message || response.status);
            return 0.00;
        }

        const data = await response.json();
        if (data.success) {
            return parseFloat(data.total) || 0.00;
        } else {
            console.error(`AdminOrdersPage: Backend reported failure for total price for OrderID ${orderId}:`, data.message);
            return 0.00;
        }
    } catch (err) {
        console.error(`AdminOrdersPage: Network or parsing error fetching total price for OrderID ${orderId}:`, err);
        return 0.00;
    }
};

const fetchUndoneOrders = async () => {
  console.log("AdminOrdersPage: fetchUndoneOrders called.");
  loading.value = true;
  error.value = null;

  if (!userStore.isLoggedIn || !userStore.isAdmin) {
    error.value = "Unauthorized access. You must be an administrator to view this page.";
    loading.value = false;
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

    if (!response.ok) {
      const errorData = await response.json().catch(() => ({}));
      error.value = errorData.message || `Failed to fetch undone orders. HTTP Status: ${response.status}`;
      console.error("AdminOrdersPage: API error fetching undone orders:", error.value, errorData);

      if (response.status === 401 || response.status === 403) {
        console.warn("AdminOrdersPage: Unauthorized/Forbidden for undone orders. Logging out and redirecting.");
        userStore.logout();
        router.push('/login');
      }
      return;
    }

    const data = await response.json();
    console.log("AdminOrdersPage: Raw response data for undone orders:", JSON.stringify(data, null, 2));

    if (data.success) {
      const fetchedOrders = await Promise.all((data.orders || []).map(async order => {
        const totalPrice = await fetchOrderTotalPrice(order.OrderID);
        
        console.log(`OrderID ${order.OrderID}: Raw OrderTime data received:`, order.OrderTime);

        return {
          OrderID: order.OrderID,
          UserID: order.UserID,
          TotalPrice: totalPrice,
          OrderTime: order.OrderTime,
          IsRanged: order.IsRanged
        };
      }));
      orders.value = fetchedOrders;

      console.log("AdminOrdersPage: Processed undone orders:", orders.value);
      if (orders.value.length === 0) {
        console.log("AdminOrdersPage: No undone orders found after processing.");
      }
    } else {
      error.value = data.message || 'Failed to fetch undone orders (backend reported failure).';
      console.error('AdminOrdersPage: Backend reported failure for undone orders:', data.message);
    }
  } catch (err) {
    error.value = `Error fetching undone orders: ${err.message}. Please check your network connection and backend server.`;
    console.error('AdminOrdersPage: Network or parsing error fetching undone orders:', err);
  } finally {
    loading.value = false;
    console.log("AdminOrdersPage: fetchUndoneOrders finished. loading set to false.");
  }
};

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

    if (!response.ok) {
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
      totalAverageSpending.value = parseFloat(data.average) || 0.00;
      console.log("AdminOrdersPage: Total average spending fetched and parsed:", totalAverageSpending.value);
    } else {
      totalAverageError.value = data.message || 'Failed to fetch total average spending (backend reported failure).';
      console.error("AdminOrdersPage: Backend reported failure for total average spending:", data.message);
      totalAverageSpending.value = 0;
    }
  } catch (err) {
    totalAverageError.value = `Error fetching total average spending: ${err.message}.`;
    console.error('AdminOrdersPage: Network error fetching total average spending:', err);
    totalAverageSpending.value = 0;
  } finally {
    totalAverageLoading.value = false;
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

const formatDateTime = (sqlNullTimeObject) => {
  if (sqlNullTimeObject && typeof sqlNullTimeObject.Valid === 'boolean') {
    if (sqlNullTimeObject.Valid) {
      try {
        const date = new Date(sqlNullTimeObject.Time);
        if (!isNaN(date.getTime())) {
          return date.toLocaleString('en-SG', {
            year: 'numeric',
            month: 'numeric',
            day: 'numeric',
            hour: '2-digit',
            minute: '2-digit',
            second: '2-digit',
            hour12: false
          });
        } else {
          console.warn("formatDateTime: new Date() created an invalid date from:", sqlNullTimeObject.Time);
        }
      } catch (e) {
        console.error("formatDateTime: Error parsing date from sql.NullTime.Time:", e, sqlNullTimeObject.Time);
      }
    } else {
      console.log("formatDateTime: OrderTime is NULL in the database (Valid: false).");
    }
  } else if (typeof sqlNullTimeObject === 'string' && sqlNullTimeObject) {
      try {
          const date = new Date(sqlNullTimeObject);
          if (!isNaN(date.getTime())) {
              return date.toLocaleString('en-SG', {
                year: 'numeric',
                month: 'numeric',
                day: 'numeric',
                hour: '2-digit',
                minute: '2-digit',
                second: '2-digit',
                hour12: false
              });
          } else {
            console.warn("formatDateTime: new Date() created an invalid date from direct string:", sqlNullTimeObject);
          }
      } catch (e) {
          console.error("formatDateTime: Error parsing date from direct string:", e, sqlNullTimeObject);
      }
  } else if (sqlNullTimeObject === null) {
      console.log("formatDateTime: OrderTime received as explicit null.");
  } else {
      console.warn("formatDateTime: Unexpected format for OrderTime:", sqlNullTimeObject);
  }

  return 'Not recorded yet'; 
};

onMounted(() => {
  console.log("AdminOrdersPage: Component mounted.");
  fetchUndoneOrders();
  fetchTotalAverageSpending();
});
</script>

<style scoped>
/* Cursor pointer for sortable columns */
th.cursor-pointer {
  cursor: pointer;
}

/* Style for active sort column header */
th.bg-gray-300 {
  background-color: #d1d5db; /* A slightly darker gray */
  color: #1f2937; /* Darker text */
}

/* Flexbox for header content to align text and icon */
.flex.items-center {
  display: flex;
  align-items: center;
}
</style>