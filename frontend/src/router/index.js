import { createRouter, createWebHistory } from 'vue-router';

// Import your components
import LoginForm from '../components/LoginForm.vue';
import RegisterForm from '../components/RegisterForm.vue';
import ChangeInfoForm from '../components/ChangeInfoForm.vue';
import RechargeWalletForm from '../components/RechargeWalletForm.vue';
import MenuPage from '../components/MenuPage.vue';
import CreateFoodForm from '../components/CreateFoodForm.vue';
import UserOrdersPage from '../components/UserOrdersPage.vue';
import OrderDetailPage from '../components/OrderDetailPage.vue';
import FoodDetailPage from '../components/FoodDetailPage.vue';
import PaymentPage from '../components/PaymentPage.vue';

// Define your routes with meta fields for authentication/authorization
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
  },
  {
    path: '/register',
    name: 'Register',
    component: RegisterForm,
  },
  {
    path: '/menu',
    name: 'Menu',
    component: MenuPage,
  },
  {
    path: '/food/:id', // <-- NEW DYNAMIC ROUTE FOR FOOD DETAILS
    name: 'FoodDetail',
    component: FoodDetailPage,
    props: true,
  },
  {
    path: '/orders',
    name: 'UserOrders',
    component: UserOrdersPage,
  },
  {
    path: '/orders/:id',
    name: 'OrderDetail',
    component: OrderDetailPage,
    props: true,
  },
  {
    path: '/change-info',
    name: 'ChangeInfo',
    component: ChangeInfoForm,
  },
  {
    path: '/pay/:orderId', 
    name: 'PaymentPage',
    component: PaymentPage,
    props: true, 
  },
  {
    path: '/recharge-wallet',
    name: 'RechargeWallet',
    component: RechargeWalletForm,
  },
  {
    path: '/admin/create-food',
    name: 'CreateFood',
    component: CreateFoodForm,
  },
  {
    path: '/:catchAll(.*)',
    name: 'NotFound',
    component: { template: '<div class="text-center p-8 text-2xl font-bold text-gray-700">404 - Page Not Found</div>' }
  }
];

// Create the router instance
const router = createRouter({
  history: createWebHistory(),
  routes,
});

export default router;