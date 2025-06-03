import { defineStore } from 'pinia';
import { ref, computed } from 'vue';

export const useUserStore = defineStore('user', () => {
  const userToken = ref(localStorage.getItem('userToken'));
  const isLoggedIn = computed(() => !!userToken.value); // Derived state
  const isAdmin = ref(localStorage.getItem('isAdmin') === 'true'); // Store and retrieve isAdmin status

  // This function would be called after a successful login API call
  const setLoginState = (token, is_admin_status) => {
    userToken.value = token;
    localStorage.setItem('userToken', token);
    isAdmin.value = is_admin_status;
    localStorage.setItem('isAdmin', is_admin_status); // Store isAdmin status
  };

  const logout = () => {
    userToken.value = null;
    localStorage.removeItem('userToken');
    isAdmin.value = false;
    localStorage.removeItem('isAdmin'); // Clear isAdmin status on logout
  };

  return { userToken, isLoggedIn, isAdmin, setLoginState, logout };
});