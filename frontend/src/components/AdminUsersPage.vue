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
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-700 uppercase tracking-wider">User ID</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-700 uppercase tracking-wider">Username</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-700 uppercase tracking-wider">Email</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-700 uppercase tracking-wider">Wallet Balance</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-700 uppercase tracking-wider">Avg. Spending per Order</th> <th class="px-6 py-3 text-left text-xs font-medium text-gray-700 uppercase tracking-wider">Is Admin</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-700 uppercase tracking-wider">Is Active</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-700 uppercase tracking-wider">Actions</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-200">
            <tr v-for="user in users" :key="user.UserID" :class="{'bg-red-50': user.IsDeleted, 'bg-gray-50': user.IsAdmin}">
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{{ user.UserID }}</td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{{ user.Username }}</td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{{ user.Email }}</td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">${{ user.WalletBalance ? user.WalletBalance.toFixed(2) : '0.00' }}</td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                ${{ user.averageSpending !== undefined ? user.averageSpending.toFixed(2) : '0.00' }} </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm">
                <span :class="{'bg-green-100 text-green-800': user.IsAdmin, 'bg-red-100 text-red-800': !user.IsAdmin}"
                      class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full">
                  {{ user.IsAdmin ? 'Yes' : 'No' }}
                </span>
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm">
                <span :class="{'bg-green-100 text-green-800': user.IsActive, 'bg-red-100 text-red-800': !user.IsActive}"
                      class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full">
                  {{ user.IsActive ? 'Yes' : 'No' }}
                </span>
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm font-medium">
                <button
                  v-if="!user.IsDeleted"
                  @click="toggleUserStatus(user.UserID, user.IsActive)"
                  :class="user.IsActive ? 'text-yellow-600 hover:text-yellow-900' : 'text-green-600 hover:text-green-900'"
                  class="mr-4"
                  title="Toggle Active Status"
                >
                  {{ user.IsActive ? 'Deactivate' : 'Activate' }}
                </button>
                <button
                  v-if="!user.IsDeleted"
                  @click="deleteUser(user.UserID)"
                  class="text-red-600 hover:text-red-900"
                  title="Delete User"
                >
                  Delete
                </button>
                <span v-else class="text-gray-500">Deleted</span>
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

const users = ref([]);
const loading = ref(true); // Single loading state for the entire page
const error = ref(null); // Single error state for the entire page

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
    // Fetch users and average spending concurrently
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

    // Check users response
    if (!usersResponse.ok || !usersData.success) {
      error.value = usersData.message || `Failed to fetch users. HTTP Status: ${usersResponse.status}`;
      console.error("AdminUsersPage: Error fetching users:", error.value);
      if (usersResponse.status === 401 || usersResponse.status === 403) {
        userStore.logout();
        router.push('/login');
      }
      return; // Stop if user data fetch fails
    }

    // Check average spending response
    if (!averageSpendingResponse.ok || !averageSpendingData.success) {
      error.value = averageSpendingData.message || `Failed to fetch average spending. HTTP Status: ${averageSpendingResponse.status}`;
      console.error("AdminUsersPage: Error fetching average spending:", error.value);
      if (averageSpendingResponse.status === 401 || averageSpendingResponse.status === 403) {
        userStore.logout();
        router.push('/login');
      }
      return; // Stop if average spending fetch fails
    }

    // Both successful, now combine data
    const fetchedUsers = usersData.users || [];
    const averageMap = averageSpendingData.Average || {};

    // Augment each user with their average spending
    users.value = fetchedUsers.map(user => ({
      ...user,
      averageSpending: averageMap[user.Username] !== undefined ? averageMap[user.Username] : 0.00 // Default to 0 if no data
    }));

    console.log("AdminUsersPage: Users and average spending combined successfully.");

  } catch (err) {
    error.value = `Network or parsing error: ${err.message}. Please check your network connection and backend server.`;
    console.error('AdminUsersPage: Combined fetch error:', err);
  } finally {
    loading.value = false;
  }
};

const toggleUserStatus = async (userId, currentStatus) => {
  const action = currentStatus ? 'deactivate' : 'activate';
  if (!confirm(`Are you sure you want to ${action} user ID ${userId}?`)) {
    return;
  }

  try {
    const endpoint = currentStatus ? 'deactivate-user' : 'activate-user';
    const response = await fetch(`http://localhost:8080/admin/${endpoint}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${userStore.userToken}`,
      },
      body: JSON.stringify({ user_id: userId }),
    });

    const data = await response.json();

    if (response.ok && data.success) {
      alert(`User ${userId} ${action}d successfully!`);
      fetchUsersAndSpending(); // Re-fetch all data to update the list
    } else {
      alert(`Failed to ${action} user ${userId}: ${data.message || 'Unknown error.'}`);
      console.error(`Failed to ${action} user ${userId}:`, data);
    }
  } catch (err) {
    alert(`Error ${action}ing user ${userId}: ${err.message}`);
    console.error(`Network error ${action}ing user:`, err);
  }
};

const deleteUser = async (userId) => {
  if (!confirm(`Are you sure you want to DELETE user ID ${userId}? This will mark them as deleted.`)) {
    return;
  }

  try {
    const response = await fetch('http://localhost:8080/admin/delete-user', {
      method: 'DELETE',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${userStore.userToken}`,
      },
      body: JSON.stringify({ user_id: userId }),
    });

    const data = await response.json();

    if (response.ok && data.success) {
      alert(`User ${userId} marked as deleted successfully!`);
      fetchUsersAndSpending(); // Re-fetch all data to update the list
    } else {
      alert(`Failed to delete user ${userId}: ${data.message || 'Unknown error.'}`);
      console.error(`Failed to delete user ${userId}:`, data);
    }
  } catch (err) {
    alert(`Error deleting user ${userId}: ${err.message}`);
    console.error('Network error deleting user:', err);
  }
};

onMounted(() => {
  console.log("AdminUsersPage: Component mounted.");
  fetchUsersAndSpending();
});
</script>

<style scoped>
/* Highlight deleted rows */
.bg-red-50 {
  background-color: #fef2f2;
}

/* Highlight admin rows (optional, distinguish from deleted) */
.bg-gray-50 {
  background-color: #f9fafb;
}
</style>