import { createRouter, createWebHistory } from 'vue-router';
import { useUserStore } from '../stores/user';

// Import all your components from src/components/
import LoginForm from '../components/LoginForm.vue';
import RegisterForm from '../components/RegisterForm.vue';
import ChangeInfoForm from '../components/ChangeInfoForm.vue';
import RechargeWalletForm from '../components/RechargeWalletForm.vue';
import MenuPage from '../components/MenuPage.vue';
import CreateFoodForm from '../components/CreateFoodForm.vue';
import UserOrdersPage from '../components/UserOrdersPage.vue';
import OrderDetailPage from '../components/OrderDetailPage.vue'; // This is the USER'S Order Detail Page
import FoodDetailPage from '../components/FoodDetailPage.vue';
import PaymentPage from '../components/PaymentPage.vue';
import AdminUsersPage from '../components/AdminUsersPage.vue';
import AdminOrdersPage from '../components/AdminOrdersPage.vue'; // The undone orders page
import AdminHistoryPage from '../components/AdminHistoryPage.vue';
import AdminOrderDetailPage from '../components/AdminOrderDetailPage.vue'; // Admin-specific order detail page

const routes = [
  {
    path: '/',
    name: 'Home',
    redirect: '/menu'
  },
  {
    path: '/login',
    name: 'Login',
    component: LoginForm,
    meta: { requiresGuest: true }
  },
  {
    path: '/register',
    name: 'Register',
    component: RegisterForm,
    meta: { requiresGuest: true }
  },
  {
    path: '/menu',
    name: 'Menu',
    component: MenuPage,
  },
  {
    path: '/food/:id',
    name: 'FoodDetail',
    component: FoodDetailPage,
    props: true,
  },
  {
    path: '/orders',
    name: 'UserOrders',
    component: UserOrdersPage,
    meta: { requiresAuth: true }
  },
  {
    path: '/orders/:id', // This is the user's specific order detail route
    name: 'UserOrderDetail', // This name is correct and matches the fix in UserOrdersPage.vue
    component: OrderDetailPage,
    props: true,
    meta: { requiresAuth: true }
  },
  {
    path: '/change-info',
    name: 'ChangeInfo',
    component: ChangeInfoForm,
    meta: { requiresAuth: true }
  },
  {
    path: '/pay/:orderId',
    name: 'PaymentPage',
    component: PaymentPage,
    props: true,
    meta: { requiresAuth: true }
  },
  {
    path: '/recharge-wallet',
    name: 'RechargeWallet',
    component: RechargeWalletForm,
    meta: { requiresAuth: true }
  },
  // Admin Routes
  {
    path: '/admin', // Redirect from base admin path
    name: 'AdminDashboard',
    redirect: '/admin/orders', // Redirects to undone orders
    meta: { requiresAuth: true, requiresAdmin: true }
  },
  {
    path: '/admin/create-food',
    name: 'CreateFood',
    component: CreateFoodForm,
    meta: { requiresAuth: true, requiresAdmin: true }
  },
  {
    path: '/admin/users',
    name: 'AdminUsers',
    component: AdminUsersPage,
    meta: { requiresAuth: true, requiresAdmin: true }
  },
  {
    path: '/admin/orders', // This is the undone orders page
    name: 'AdminOrders',
    component: AdminOrdersPage,
    meta: { requiresAuth: true, requiresAdmin: true }
  },
  {
    path: '/admin/history',
    name: 'AdminHistory',
    component: AdminHistoryPage,
    meta: { requiresAuth: true, requiresAdmin: true }
  },
  {
    path: '/admin/orders/:id', // Admin-specific order detail route
    name: 'AdminOrderDetail',
    component: AdminOrderDetailPage,
    props: true,
    meta: { requiresAuth: true, requiresAdmin: true }
  }
];

const router = createRouter({
  // Using import.meta.env.BASE_URL is good practice for Vite projects
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
});

// Global Navigation Guard (remains the same)
router.beforeEach((to, from, next) => {
  const userStore = useUserStore();

  if (to.meta.requiresAuth && !userStore.isLoggedIn) {
    console.log(`Router Guard: Redirecting to login for ${to.name}. Route requires auth but user is not logged in.`);
    next('/login');
    return;
  }

  if (to.meta.requiresAdmin && !userStore.isAdmin) {
    console.log(`Router Guard: Redirecting to home for ${to.name}. Route requires admin but user is not an admin.`);
    next('/');
    return;
  }

  if (to.meta.requiresGuest && userStore.isLoggedIn) {
    console.log(`Router Guard: Redirecting from ${to.name}. User is already logged in.`);
    next('/');
    return;
  }

  next();
});

export default router;