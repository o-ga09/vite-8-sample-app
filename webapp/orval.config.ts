import { defineConfig } from "orval";

export default defineConfig({
  budgetApp: {
    input: {
      target: "../api/openapi.yaml",
    },
    output: {
      mode: "tags-split",
      target: "./src/lib/api/generated",
      schemas: "./src/lib/api/model",
      client: "react-query",
      useTypeImports: true,
      override: {
        mutator: {
          path: "./src/lib/api/mutator.ts",
          name: "customInstance",
        },
        query: {
          useQuery: true,
          useSuspenseQuery: false,
          useInfiniteQuery: false,
        },
      },
    },
  },
});
