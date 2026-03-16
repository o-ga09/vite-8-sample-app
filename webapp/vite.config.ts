import { defineConfig } from "vite-plus";
import react from "@vitejs/plugin-react";
import tailwindcss from "@tailwindcss/vite";

export default defineConfig({
  lazy: async () => ({
    plugins: [...react(), ...tailwindcss()],
  }),
  staged: {
    "*": "vp check --fix",
  },
});
