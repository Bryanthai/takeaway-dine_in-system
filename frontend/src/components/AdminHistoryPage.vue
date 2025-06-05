<template>
  <div class="min-h-screen bg-gray-100 p-6 pt-10">
    <div class="max-w-7xl mx-auto bg-white rounded-lg shadow-xl p-8">
      <h2 class="text-4xl font-extrabold text-gray-900 mb-8 text-center">
        Admin Panel: Order History
      </h2>

      <div v-if="loading" class="flex justify-center items-center h-48">
        <svg class="animate-spin h-10 w-10 text-indigo-500" xmlns="http://www.w3.org/2000/svg" fill="none"
          viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
          <path class="opacity-75" fill="currentColor"
            d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z">
          </path>
        </svg>
        <p class="ml-3 text-lg text-gray-700">Loading order history...</p>
      </div>

      <div v-else-if="error" class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative"
        role="alert">
        <strong class="font-bold">Error!</strong>
        <span class="block sm:inline ml-2">{{ error }}</span>
      </div>

      <div v-else-if="orders.length === 0"
        class="bg-blue-100 border border-blue-400 text-blue-700 px-4 py-3 rounded relative text-center" role="alert">
        <strong class="font-bold">No Orders Found!</strong>
        <span class="block sm:inline ml-2">There is no order history to display.</span>
      </div>

      <div v-else class="overflow-x-auto">
        <div class="mb-6 p-4 bg-purple-50 border border-purple-200 text-purple-800 rounded-lg text-center">
          <p class="text-lg font-semibold">Total Orders in the Past: {{ totalOrdersCount }}</p>
        </div>

        <div v-if="orderedItems.length > 0" class="mb-8 p-6 bg-white rounded-lg shadow-lg">
          <FoodIdPieChart :items="orderedItems" />
        </div>
        <div v-else
          class="bg-blue-100 border border-blue-400 text-blue-700 px-4 py-3 rounded relative text-center mb-6"
          role="alert">
          <strong class="font-bold">No Ordered Items Found!</strong>
          <span class="block sm:inline ml-2">Cannot display food ID distribution.</span>
        </div>

        <div class="mb-6 flex flex-wrap justify-end items-center gap-4">
          <label for="filterDone" class="text-gray-700 font-medium">Is Done:</label>
          <select id="filterDone" v-model="filterDone"
            class="block w-36 pl-3 pr-10 py-2 text-base border-gray-300 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm rounded-md shadow-sm">
            <option value="all">All</option>
            <option value="true">Yes</option>
            <option value="false">No</option>
          </select>

          <label for="filterPaid" class="text-gray-700 font-medium">Is Paid:</label>
          <select id="filterPaid" v-model="filterPaid"
            class="block w-36 pl-3 pr-10 py-2 text-base border-gray-300 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm rounded-md shadow-sm">
            <option value="all">All</option>
            <option value="true">Yes</option>
            <option value="false">No</option>
          </select>

          <label for="filterRanged" class="text-gray-700 font-medium">Is Delivery:</label>
          <select id="filterRanged" v-model="filterRanged"
            class="block w-36 pl-3 pr-10 py-2 text-base border-gray-300 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm rounded-md shadow-sm">
            <option value="all">All</option>
            <option value="true">Yes</option>
            <option value="false">No</option>
          </select>

          <label for="filterDeleted" class="text-gray-700 font-medium">Is Deleted:</label>
          <select id="filterDeleted" v-model="filterDeleted"
            class="block w-36 pl-3 pr-10 py-2 text-base border-gray-300 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm rounded-md shadow-sm">
            <option value="all">All</option>
            <option value="true">Yes</option>
            <option value="false">No</option>
          </select>
        </div>


        <table class="min-w-full bg-white border border-gray-200 rounded-lg shadow-md">
          <thead class="bg-gray-200">
            <tr>
              <th
                class="px-6 py-3 text-left text-xs font-medium text-gray-700 uppercase tracking-wider cursor-pointer hover:bg-gray-300 transition-colors duration-200"
                :class="{ 'bg-gray-300 text-gray-900': sortColumn === 'OrderID' }" @click="sortBy('OrderID')">
                <div class="flex items-center">
                  Order ID
                  <span v-if="sortColumn === 'OrderID'" class="ml-1">
                    <span v-if="sortDirection === 'asc'">&#9650;</span>
                    <span v-else>&#9660;</span>
                  </span>
                </div>
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-700 uppercase tracking-wider">User ID</th>
              <th
                class="px-6 py-3 text-left text-xs font-medium text-gray-700 uppercase tracking-wider cursor-pointer hover:bg-gray-300 transition-colors duration-200"
                :class="{ 'bg-gray-300 text-gray-900': sortColumn === 'OrderTime' }" @click="sortBy('OrderTime')">
                <div class="flex items-center">
                  Order Time
                  <span v-if="sortColumn === 'OrderTime'" class="ml-1">
                    <span v-if="sortDirection === 'asc'">&#9650;</span>
                    <span v-else>&#9660;</span>
                  </span>
                </div>
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-700 uppercase tracking-wider">Order Info</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-700 uppercase tracking-wider">Is Done</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-700 uppercase tracking-wider">Is Paid</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-700 uppercase tracking-wider">Is Delivery
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-700 uppercase tracking-wider">Deleted</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-700 uppercase tracking-wider">Actions</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-200">
            <tr v-for="order in filteredAndSortedOrders" :key="order.OrderID" :class="{ 'bg-red-50': order.Deleted }">
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{{ order.OrderID }}</td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{{ order.UserID }}</td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                {{ formatDateTime(order.OrderTime.Valid ? order.OrderTime.Time : null) }}
              </td>
              <td class="px-6 py-4 text-sm text-gray-900 max-w-xs truncate" :title="order.OrderInfo">
                {{ order.OrderInfo }}
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm">
                <span
                  :class="{ 'bg-green-100 text-green-800': order.IsDone, 'bg-yellow-100 text-yellow-800': !order.IsDone }"
                  class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full">
                  {{ order.IsDone ? 'Yes' : 'No' }}
                </span>
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm">
                <span
                  :class="{ 'bg-green-100 text-green-800': order.IsPaid, 'bg-red-100 text-red-800': !order.IsPaid }"
                  class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full">
                  {{ order.IsPaid ? 'Yes' : 'No' }}
                </span>
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm">
                <span
                  :class="{ 'bg-blue-100 text-blue-800': order.IsRanged, 'bg-gray-100 text-gray-800': !order.IsRanged }"
                  class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full">
                  {{ order.IsRanged ? 'Delivery' : 'Pickup' }}
                </span>
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm">
                <span
                  :class="{ 'bg-red-100 text-red-800': order.Deleted, 'bg-gray-100 text-gray-800': !order.Deleted }"
                  class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full">
                  {{ order.Deleted ? 'Yes' : 'No' }}
                </span>
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm font-medium">
                <router-link :to="`/admin/orders/${order.OrderID}`" class="text-blue-600 hover:text-blue-900"
                  title="View Details">
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
import { ref, onMounted, computed } from 'vue';
import { useUserStore } from '../stores/user';
import { useRouter } from 'vue-router';
import FoodIdPieChart from '../components/FoodIdPieChart.vue'; // Import the new component

const userStore = useUserStore();
const router = useRouter();

const orders = ref([]);
const orderedItems = ref([]); // Ref for ordered items
const loading = ref(true);
const error = ref(null);

// Sorting state
const sortColumn = ref('OrderID'); // Default sort column
const sortDirection = ref('asc'); // Default sort direction

// Filtering state
const filterDone = ref('all');     // 'all', 'true', 'false'
const filterPaid = ref('all');     // 'all', 'true', 'false'
const filterRanged = ref('all');   // 'all', 'true', 'false'
const filterDeleted = ref('all');  // 'all', 'true', 'false'


// New computed property to count total orders
const totalOrdersCount = computed(() => {
  return orders.value.length;
});

// Function to handle sorting
const sortBy = (column) => {
  if (sortColumn.value === column) {
    sortDirection.value = sortDirection.value === 'asc' ? 'desc' : 'asc';
  } else {
    sortColumn.value = column;
    sortDirection.value = 'asc';
  }
};

// Computed property for filtered and then sorted orders
const filteredAndSortedOrders = computed(() => {
  let currentOrders = [...orders.value]; // Start with a shallow copy of the original data

  // 1. Apply Filtering
  if (filterDone.value !== 'all') {
    const isDoneBool = filterDone.value === 'true';
    currentOrders = currentOrders.filter(order => order.IsDone === isDoneBool);
  }
  if (filterPaid.value !== 'all') {
    const isPaidBool = filterPaid.value === 'true';
    currentOrders = currentOrders.filter(order => order.IsPaid === isPaidBool);
  }
  if (filterRanged.value !== 'all') {
    const isRangedBool = filterRanged.value === 'true';
    currentOrders = currentOrders.filter(order => order.IsRanged === isRangedBool);
  }
  if (filterDeleted.value !== 'all') {
    const isDeletedBool = filterDeleted.value === 'true';
    currentOrders = currentOrders.filter(order => order.Deleted === isDeletedBool);
  }

  // 2. Apply Sorting
  return currentOrders.sort((a, b) => {
    let aValue = a[sortColumn.value];
    let bValue = b[sortColumn.value];

    // Special handling for OrderTime (sql.NullTime)
    if (sortColumn.value === 'OrderTime') {
      const aTime = aValue && aValue.Valid ? new Date(aValue.Time).getTime() : 0;
      const bTime = bValue && bValue.Valid ? new Date(bValue.Time).getTime() : 0;

      // Handle null/invalid dates by treating them as smaller/larger based on direction
      if (aTime === 0 && bTime === 0) return 0;
      if (aTime === 0) return sortDirection.value === 'asc' ? 1 : -1; // Null comes last in asc
      if (bTime === 0) return sortDirection.value === 'asc' ? -1 : 1; // Null comes first in asc

      return sortDirection.value === 'asc' ? aTime - bTime : bTime - aTime;
    }

    // Default string/number comparison
    if (aValue < bValue) {
      return sortDirection.value === 'asc' ? -1 : 1;
    }
    if (aValue > bValue) {
      return sortDirection.value === 'asc' ? 1 : -1;
    }
    return 0;
  });
});


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

const fetchOrderedItems = async () => {
  console.log("AdminHistoryPage: fetchOrderedItems called.");

  if (!userStore.isLoggedIn || !userStore.isAdmin) {
    console.warn("AdminHistoryPage: Access denied for fetching ordered items (client-side check).");
    return;
  }

  try {
    console.log("AdminHistoryPage: Attempting to fetch ordered items from:", 'http://localhost:8080/orders/all-items');
    const response = await fetch('http://localhost:8080/orders/all-items', {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${userStore.userToken}`,
      },
    });

    console.log("AdminHistoryPage: Ordered items fetch response received. Status:", response.status);

    const data = await response.json();
    console.log("AdminHistoryPage: Ordered items response data:", data);

    if (response.ok && data.success) {
      orderedItems.value = data.items || [];
      console.log("AdminHistoryPage: All ordered items fetched successfully. Count:", orderedItems.value.length);
    } else {
      console.error("AdminHistoryPage: API error fetching ordered items:", data.message || `Failed to fetch ordered items. HTTP Status: ${response.status}`);
    }
  } catch (err) {
    console.error('AdminHistoryPage: Network or parsing error fetching ordered items:', err);
  }
};

const formatDateTime = (dateValue) => {
  if (!dateValue) return 'N/A';
  try {
    const date = new Date(dateValue);
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
    return String(dateValue);
  }
};

onMounted(() => {
  console.log("AdminHistoryPage: Component mounted.");
  fetchAllOrders();
  fetchOrderedItems();
});
</script>

<style scoped>
/* Cursor pointer for sortable columns */
th.cursor-pointer {
  cursor: pointer;
}

/* Style for active sort column header */
th.bg-gray-300 {
  background-color: #d1d5db;
  /* A slightly darker gray */
  color: #1f2937;
  /* Darker text */
}

/* Flexbox for header content to align text and icon */
.flex.items-center {
  display: flex;
  align-items: center;
}

/* Optional: Highlight deleted rows */
.bg-red-50 {
  background-color: #fef2f2;
  /* Light red background for deleted orders */
}
</style>