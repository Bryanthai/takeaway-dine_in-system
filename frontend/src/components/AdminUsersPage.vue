<template>
  <div class="min-h-screen bg-gray-100 p-6 pt-10">
    <div class="max-w-7xl mx-auto bg-white rounded-lg shadow-xl p-8">
      <h2 class="text-4xl font-extrabold text-gray-900 mb-8 text-center">
        Admin Panel: Users Overview
      </h2>

      <div v-if="loading" class="flex justify-center items-center h-48">
        <svg class="animate-spin h-10 w-10 text-indigo-500" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
        </svg>
        <p class="ml-3 text-lg text-gray-700">Loading user data and spending averages...</p>
      </div>

      <div v-else-if="error" class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative" role="alert">
        <strong class="font-bold">Error!</strong>
        <span class="block sm:inline ml-2">{{ error }}</span>
      </div>

      <div v-else-if="users.length === 0" class="bg-blue-100 border border-blue-400 text-blue-700 px-4 py-3 rounded relative text-center" role="alert">
        <strong class="font-bold">No Users Found!</strong>
        <span class="block sm:inline ml-2">There are no registered users to display.</span>
      </div>

      <div v-else class="overflow-x-auto">
        <table class="min-w-full bg-white border border-gray-200 rounded-lg shadow-md">
          <thead class="bg-gray-200">
            <tr>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-700 uppercase tracking-wider cursor-pointer hover:bg-gray-300 transition-colors duration-200"
                  :class="{ 'bg-gray-300 text-gray-900': sortColumn === 'UserID' }"
                  @click="sortBy('UserID')">
                <div class="flex items-center">
                  User ID
                  <span v-if="sortColumn === 'UserID'" class="ml-1">
                    <span v-if="sortDirection === 'asc'">&#9650;</span> <span v-else>&#9660;</span> </span>
                </div>
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-700 uppercase tracking-wider">Username</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-700 uppercase tracking-wider">Email</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-700 uppercase tracking-wider cursor-pointer hover:bg-gray-300 transition-colors duration-200"
                  :class="{ 'bg-gray-300 text-gray-900': sortColumn === 'WalletBalance' }"
                  @click="sortBy('WalletBalance')">
                <div class="flex items-center">
                  Wallet Balance
                  <span v-if="sortColumn === 'WalletBalance'" class="ml-1">
                    <span v-if="sortDirection === 'asc'">&#9650;</span>
                    <span v-else>&#9660;</span>
                  </span>
                </div>
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-700 uppercase tracking-wider cursor-pointer hover:bg-gray-300 transition-colors duration-200"
                  :class="{ 'bg-gray-300 text-gray-900': sortColumn === 'averageSpending' }"
                  @click="sortBy('averageSpending')">
                <div class="flex items-center">
                  Avg. Spending per Order
                  <span v-if="sortColumn === 'averageSpending'" class="ml-1">
                    <span v-if="sortDirection === 'asc'">&#9650;</span>
                    <span v-else>&#9660;</span>
                  </span>
                </div>
              </th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-200">
            <tr v-for="user in sortedUsers" :key="user.UserID" :class="{'bg-red-50': user.IsDeleted, 'bg-gray-50': user.IsAdmin}">
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{{ user.UserID }}</td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{{ user.Username }}</td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{{ user.Email }}</td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">${{ user.WalletBalance ? user.WalletBalance.toFixed(2) : '0.00' }}</td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                ${{ user.averageSpending !== undefined ? user.averageSpending.toFixed(2) : '0.00' }}
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

const userStore = useUserStore();
const router = useRouter();

const users = ref([]);
const loading = ref(true);
const error = ref(null);

// Reactive properties for sorting
const sortColumn = ref('UserID'); // Default sort column
const sortDirection = ref('asc'); // Default sort direction

// Function to handle sorting
const sortBy = (column) => {
  if (sortColumn.value === column) {
    // If clicking the same column, reverse the sort direction
    sortDirection.value = sortDirection.value === 'asc' ? 'desc' : 'asc';
  } else {
    // If clicking a new column, set it as the sort column and default to ascending
    sortColumn.value = column;
    sortDirection.value = 'asc';
  }
};

// Computed property for sorted users
const sortedUsers = computed(() => {
  if (!users.value || users.value.length === 0) {
    return [];
  }

  // Create a shallow copy to avoid mutating the original array
  return [...users.value].sort((a, b) => {
    const aValue = a[sortColumn.value];
    const bValue = b[sortColumn.value];

    // Handle string comparison for non-numeric fields if needed,
    // though for UserID, WalletBalance, and averageSpending, direct comparison is fine.
    if (typeof aValue === 'string' && typeof bValue === 'string') {
      return sortDirection.value === 'asc' ? aValue.localeCompare(bValue) : bValue.localeCompare(aValue);
    } else {
      if (aValue < bValue) {
        return sortDirection.value === 'asc' ? -1 : 1;
      }
      if (aValue > bValue) {
        return sortDirection.value === 'asc' ? 1 : -1;
      }
      return 0;
    }
  });
});

const fetchUsersAndSpending = async () => {
  console.log("AdminUsersPage: fetchUsersAndSpending called.");
  loading.value = true;
  error.value = null;

  if (!userStore.isLoggedIn || !userStore.isAdmin) {
    error.value = "Unauthorized access. You must be an administrator to view this page.";
    loading.value = false;
    router.push('/');
    return;
  }

  try {
    const [usersResponse, averageSpendingResponse] = await Promise.all([
      fetch('http://localhost:8080/admin/users', {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${userStore.userToken}`,
        },
      }),
      fetch('http://localhost:8080/admin/total-average-by-user', {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${userStore.userToken}`,
        },
      })
    ]);

    const usersData = await usersResponse.json();
    const averageSpendingData = await averageSpendingResponse.json();

    if (!usersResponse.ok || !usersData.success) {
      error.value = usersData.message || `Failed to fetch users. HTTP Status: ${usersResponse.status}`;
      console.error("AdminUsersPage: Error fetching users:", error.value);
      if (usersResponse.status === 401 || usersResponse.status === 403) {
        userStore.logout();
        router.push('/login');
      }
      return;
    }

    if (!averageSpendingResponse.ok || !averageSpendingData.success) {
      error.value = averageSpendingData.message || `Failed to fetch average spending. HTTP Status: ${averageSpendingResponse.status}`;
      console.error("AdminUsersPage: Error fetching average spending:", error.value);
      if (averageSpendingResponse.status === 401 || averageSpendingResponse.status === 403) {
        userStore.logout();
        router.push('/login');
      }
      return;
    }

    const fetchedUsers = usersData.users || [];
    const averageSpendingArray = averageSpendingData.average || [];

    // Convert the array of average spending objects into a map for easier lookup
    const averageMap = averageSpendingArray.reduce((acc, item) => {
      acc[item.username] = item.average;
      return acc;
    }, {});

    users.value = fetchedUsers.map(user => {
      const avgSp = averageMap[user.Username];

      return {
        UserID: user.ID,
        Username: user.Username,
        Email: user.Email,
        WalletBalance: parseFloat(user.Balance) || 0.00,
        IsAdmin: user.IsAdmin,
        IsActive: true,
        IsDeleted: false,
        averageSpending: parseFloat(avgSp) || 0.00
      };
    });

    console.log("AdminUsersPage: Users and average spending combined successfully.");

  } catch (err) {
    error.value = `Network or parsing error: ${err.message}. Please check your network connection and backend server.`;
    console.error('AdminUsersPage: Combined fetch error:', err);
  } finally {
    loading.value = false;
  }
};

onMounted(() => {
  console.log("AdminUsersPage: Component mounted.");
  fetchUsersAndSpending();
});
</script>

<style scoped>
/* Highlight deleted rows (still kept in case you want to visually indicate deleted users based on IsDeleted data in user object) */
.bg-red-50 {
  background-color: #fef2f2;
}

/* Highlight admin rows (still kept in case you want to visually indicate admin users based on IsAdmin data in user object) */
.bg-gray-50 {
  background-color: #f9fafb;
}

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