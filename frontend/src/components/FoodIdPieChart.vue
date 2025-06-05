<script setup>
import { computed } from 'vue';
import { Pie } from 'vue-chartjs';
import { Chart as ChartJS, Title, Tooltip, Legend, ArcElement, CategoryScale } from 'chart.js';

ChartJS.register(Title, Tooltip, Legend, ArcElement, CategoryScale);

const props = defineProps({
  items: {
    type: Array,
    required: true,
  },
});

const chartData = computed(() => {
  const foodIdCounts = {};
  props.items.forEach(item => {
    foodIdCounts[item.FoodID] = (foodIdCounts[item.FoodID] || 0) + item.Quantity;
  });

  const labels = Object.keys(foodIdCounts).map(id => `Food ID: ${id}`);
  const data = Object.values(foodIdCounts);

  // Generate a distinct color for each slice
  const backgroundColors = data.map((_, index) => `hsl(${index * 60}, 70%, 70%)`);

  return {
    labels: labels,
    datasets: [
      {
        backgroundColor: backgroundColors,
        data: data,
      },
    ],
  };
});

const chartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    title: {
      display: true,
      text: 'Distribution of Food Items Ordered',
      font: {
        size: 18
      }
    },
    tooltip: {
        callbacks: {
            label: function(context) {
                let label = context.label || '';
                if (label) {
                    label += ': ';
                }
                if (context.parsed !== null) {
                    label += context.parsed + ' items';
                }
                return label;
            }
        }
    }
  }
};
</script>

<template>
  <div class="h-96">
    <Pie :data="chartData" :options="chartOptions" />
  </div>
</template>

<style scoped>
/* You can add specific styles for the chart container here if needed */
</style>