/** @type {import('tailwindcss').Config} */
const config = {
  content: ["./index.html", "./src/**/*.{js,ts,jsx,tsx}"],
  theme: {
    extend: {
      width: {
        main: "70rem",
      },
      colors: {
        primary: "#1890ff",
      },
      textColor: {
        primary: "rgba(0,0,0,.85)",
      },
    },
  },
  corePlugins: {
    preflight: false,
  },
};

export default config;
