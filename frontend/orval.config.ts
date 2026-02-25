import { defineConfig } from "orval";

export default defineConfig({
  v1: {
    input: "../backend/openapi.yaml",
    output: {
      target: "./src/services/api/v1.ts",
      client: "react-query",
      httpClient: "fetch",
      mode: "single",
      override: {
        mutator: {
          path: "./src/services/api/axios-v1.ts",
          name: "customInstance",
        },
      },
    },
  },
  v1Public: {
    input: "../backend/openapi-public.yaml",
    output: {
      target: "./src/services/api/v1-public.ts",
      client: "react-query",
      httpClient: "fetch",
      mode: "single",
      override: {
        mutator: {
          path: "./src/services/api/axios-v1-public.ts",
          name: "customInstance",
        },
      },
    },
  },
});
