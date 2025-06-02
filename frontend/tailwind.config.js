/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}", // This line is crucial for scanning your Vue components
  ],
  theme: {
    extend: {},
  },
  plugins: [],
}