import { defineStore } from 'pinia';
import { ref, computed } from 'vue';

export const useCartStore = defineStore('cart', () => {
  // State
  const items = ref([]); // Array of { food: FoodObject, quantity: number }

  // Getters
  const totalItems = computed(() => items.value.reduce((sum, item) => sum + item.quantity, 0));
  const totalPrice = computed(() =>
    items.value.reduce((sum, item) => sum + item.food.Price * item.quantity, 0)
  );

  // Actions
  function addItem(foodItem) {
    const existingItem = items.value.find(item => item.food.FoodID === foodItem.FoodID);
    if (existingItem) {
      existingItem.quantity++;
    } else {
      items.value.push({ food: foodItem, quantity: 1 });
    }
  }

  function removeItem(foodId) {
    items.value = items.value.filter(item => item.food.FoodID !== foodId);
  }

  function updateQuantity(foodId, newQuantity) {
    const item = items.value.find(item => item.food.FoodID === foodId);
    if (item) {
      if (newQuantity <= 0) {
        removeItem(foodId);
      } else {
        item.quantity = newQuantity;
      }
    }
  }

  // NEW ACTIONS: incrementItem and decrementItem
  function incrementItem(foodId) {
    const item = items.value.find(item => item.food.FoodID === foodId);
    if (item) {
      item.quantity++;
    }
  }

  function decrementItem(foodId) {
    const item = items.value.find(item => item.food.FoodID === foodId);
    if (item) {
      item.quantity--;
      if (item.quantity <= 0) {
        removeItem(foodId); // Remove if quantity drops to 0 or less
      }
    }
  }

  function clearCart() {
    items.value = [];
  }

  return {
    items,
    totalItems,
    totalPrice,
    addItem,
    removeItem,
    updateQuantity,
    incrementItem, // EXPOSE NEW ACTION
    decrementItem, // EXPOSE NEW ACTION
    clearCart,
  };
});