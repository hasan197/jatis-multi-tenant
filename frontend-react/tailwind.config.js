/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {},
  },
  plugins: [],
  // Penting: Nonaktifkan preflight untuk menghindari konflik dengan Ant Design
  corePlugins: {
    preflight: false,
  },
} 