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
    clearCart,
  };
});