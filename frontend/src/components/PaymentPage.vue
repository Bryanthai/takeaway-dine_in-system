<template>
  <div class="min-h-screen bg-gray-100 p-6 flex items-center justify-center">
    <div class="bg-white p-8 rounded-lg shadow-xl w-full max-w-lg">
      <h2 class="text-3xl font-bold text-gray-800 mb-6 text-center">Complete Your Payment</h2>

      <div v-if="orderLoading || balanceLoading" class="flex flex-col items-center justify-center h-48">
        <svg class="animate-spin h-10 w-10 text-indigo-500" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
        </svg>
        <p class="ml-3 text-lg text-gray-700 mt-4">Loading data...</p>
      </div>

      <div v-else-if="fetchError" class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative" role="alert">
        <strong class="font-bold">Error!</strong>
        <span class="block sm:inline ml-2">{{ fetchError }}</span>
      </div>

      <div v-else-if="!order" class="bg-yellow-100 border border-yellow-400 text-yellow-700 px-4 py-3 rounded relative text-center" role="alert">
        <strong class="font-bold">Order Not Found!</strong>
        <span class="block sm:inline ml-2">The requested order could not be loaded.</span>
      </div>

      <div v-else>
        <div class="mb-6 p-4 bg-gray-50 rounded-md border border-gray-200">
          <p class="text-gray-700 text-lg mb-2"><span class="font-semibold">Order ID:</span> <span class="font-bold text-indigo-600">#{{ order.OrderID }}</span></p>
          <p class="text-gray-700"><span class="font-semibold">Status:</span> {{ order.IsPaid ? 'Paid' : 'Unpaid' }}</p>
          <p class="text-gray-700"><span class="font-semibold">Type:</span> {{ order.IsRanged ? 'Delivery' : 'Pickup' }}</p>

          <p class="text-gray-800 font-semibold mt-4 mb-2">Items:</p>
          <ul class="list-disc list-inside text-gray-700 mb-4 max-h-32 overflow-y-auto">
            <li v-for="(item, index) in parseOrderInfo(order.OrderInfo)" :key="index">
              {{ item.FoodName }} (x{{ item.Quantity }}) - ${{ item.Price ? (item.Price * item.Quantity).toFixed(2) : 'N/A' }}
            </li>
          </ul>

          <p class="text-xl font-bold text-gray-900 border-t pt-4">
            Total Amount: <span class="float-right">${{ totalAmount.toFixed(2) }}</span>
          </p>
        </div>

        <div class="text-gray-800 text-lg text-center mb-6">
          Your Current Wallet Balance: <span class="font-bold text-green-600">${{ userBalance.toFixed(2) }}</span>
        </div>

        <button
          @click="processPayment"
          :disabled="paymentLoading || order.IsPaid || userBalance < totalAmount"
          class="w-full bg-green-600 text-white py-3 rounded-lg font-semibold text-lg hover:bg-green-700 transition duration-200 disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center"
        >
          <svg v-if="paymentLoading" class="animate-spin -ml-1 mr-3 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
          {{ paymentLoading ? 'Processing Payment...' : 'Pay Now' }}
        </button>

        <div v-if="paymentSuccess" class="mt-4 text-green-600 text-center font-medium">
          Payment successful! Redirecting to orders...
        </div>
        <div v-if="paymentError" class="mt-4 text-red-600 text-center font-medium">
          {{ paymentError }}
          <button v-if="paymentError.includes('balance') || paymentError.includes('funds')" @click="router.push('/recharge-wallet')" class="text-blue-600 underline ml-2">Recharge Wallet</button>
        </div>

        <button
          @click="router.back()"
          class="mt-4 w-full bg-gray-300 text-gray-800 px-6 py-2 rounded-lg font-semibold hover:bg-gray-400 transition duration-200"
        >
          Go Back
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue';
import { useRoute, useRouter } from 'vue-router';

const route = useRoute();
const router = useRouter();

const orderId = ref(null);
const order = ref(null); // Stores the fetched order details
const userBalance = ref(0); // Stores the fetched user balance

const userToken = ref(localStorage.getItem('userToken') || '');

const orderLoading = ref(true);
const balanceLoading = ref(true);
const fetchError = ref(null); // General error for initial data fetching

const paymentLoading = ref(false);
const paymentSuccess = ref(false);
const paymentError = ref(null);

// Computed property for total amount
const totalAmount = computed(() => {
  if (!order.value || !order.value.OrderInfo) return 0;
  return parseOrderInfo(order.value.OrderInfo).reduce((sum, item) => sum + (item.Price * item.Quantity || 0), 0);
});

// Helper function to parse order info from string
const parseOrderInfo = (orderInfoString) => {
  try {
    if (!orderInfoString) return [];
    return JSON.parse(orderInfoString);
  } catch (e) {
    console.error("Error parsing order info:", orderInfoString, e);
    return [{ FoodName: "Invalid/Unreadable Info", Quantity: 1, Price: 0 }];
  }
};

// Function to fetch a specific order's details
const fetchOrderDetail = async () => {
  orderLoading.value = true;
  fetchError.value = null;

  if (!userToken.value) {
    fetchError.value = 'Authentication token is missing. Please log in to view order details.';
    orderLoading.value = false;
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

    // Check if the response itself indicates an error before parsing JSON
    if (!response.ok) {
      const errorData = await response.json(); // Try to parse error message if available
      fetchError.value = errorData.message || `Failed to fetch order details (HTTP ${response.status}).`;
      if (response.status === 401 || response.status === 403) {
        localStorage.removeItem('userToken');
        router.push('/login');
      }
      orderLoading.value = false;
      return;
    }

    const data = await response.json();

    // The /orders/user endpoint returns a {success: bool, orders: []} structure.
    // So `data.success` IS expected here.
    if (data.success) {
      const foundOrder = (data.orders || []).find(o => o.OrderID == orderId.value);
      if (foundOrder) {
        order.value = foundOrder;
        if (order.value.IsPaid || order.value.Deleted) {
          fetchError.value = "This order is already paid or cancelled and cannot be processed further.";
          setTimeout(() => router.push('/orders'), 2000);
          return;
        }
      } else {
        fetchError.value = 'Order not found or does not belong to this user.';
      }
    } else {
      fetchError.value = data.message || 'Failed to fetch order details.';
    }
  } catch (err) {
    fetchError.value = `Error fetching order details: ${err.message}.`;
    console.error('Error fetching order details:', err);
  } finally {
    orderLoading.value = false;
  }
};

// Function to fetch the current user's balance
const fetchUserBalance = async () => {
  balanceLoading.value = true;
  fetchError.value = null; // Clear general fetch error, specific payment errors set below

  if (!userToken.value) {
    fetchError.value = 'Authentication token is missing. Cannot fetch balance.';
    balanceLoading.value = false;
    router.push('/login');
    return;
  }

  try {
    const response = await fetch('http://localhost:8080/users', {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${userToken.value}`,
      },
    });

    // Check if the response itself indicates an error (e.g., 401, 500)
    if (!response.ok) {
      const errorData = await response.json(); // Try to parse error message
      fetchError.value = errorData.message || `Failed to fetch wallet balance (HTTP ${response.status}).`;
      if (response.status === 401 || response.status === 403) {
        localStorage.removeItem('userToken');
        router.push('/login');
      }
      balanceLoading.value = false;
      return;
    }

    // Backend returns the Account struct directly for /users
    const accountData = await response.json();

    // Now, `accountData` is the Account struct.
    userBalance.value = accountData.Balance; // Access Balance directly from the accountData object

  } catch (err) {
    fetchError.value = `Error fetching balance: ${err.message}.`;
    console.error('Error fetching user balance:', err);
  } finally {
    balanceLoading.value = false;
  }
};


// Function to process payment via backend
const processPayment = async () => {
  paymentLoading.value = true;
  paymentSuccess.value = false;
  paymentError.value = null;

  if (!userToken.value) {
    paymentError.value = 'Authentication token is missing. Please log in.';
    paymentLoading.value = false;
    router.push('/login');
    return;
  }
  if (!order.value) {
    paymentError.value = 'Order details not loaded.';
    paymentLoading.value = false;
    return;
  }
  if (order.value.IsPaid) {
    paymentError.value = 'This order has already been paid.';
    paymentLoading.value = false;
    return;
  }
  if (userBalance.value < totalAmount.value) {
    paymentError.value = 'Insufficient wallet balance. Please recharge your wallet.';
    paymentLoading.value = false;
    return;
  }

  try {
    const response = await fetch('http://localhost:8080/payment', { // PUT /payment
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${userToken.value}`,
      },
      body: JSON.stringify({
        order_id: order.value.OrderID,
      }),
    });

    const data = await response.json();

    if (response.ok && data.success) {
      paymentSuccess.value = true;
      order.value.IsPaid = true; // Optimistically update local state
      userBalance.value -= totalAmount.value; // Deduct from local balance
      setTimeout(() => {
        router.push('/orders'); // Redirect to user's orders page after successful payment
      }, 1500);
    } else {
      paymentError.value = data.message || 'Payment failed.';
      if (response.status === 401 || response.status === 403) {
        localStorage.removeItem('userToken');
        router.push('/login');
      }
    }
  } catch (err) {
    paymentError.value = `Payment network error: ${err.message}.`;
    console.error('Error processing payment:', err);
  } finally {
    paymentLoading.value = false;
  }
};

// On component mount, fetch all necessary data
onMounted(async () => {
  orderId.value = route.params.orderId;

  if (!orderId.value) {
    fetchError.value = 'No order ID provided for payment.';
    orderLoading.value = false;
    balanceLoading.value = false;
    return;
  }

  // Fetch order details and user balance concurrently
  await Promise.all([
    fetchOrderDetail(),
    fetchUserBalance()
  ]);
});
</script>

<style scoped>
/* Add your scoped styles here */
</style>