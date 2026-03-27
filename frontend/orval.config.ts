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
        fetch: {
          includeHttpResponseReturnType: false,
        },
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
        fetch: {
          includeHttpResponseReturnType: false,
        },
        mutator: {
          path: "./src/services/api/axios-v1-public.ts",
          name: "customInstance",
        },
      },
    },
  },
  v1Admin: {
    input: "../backend/openapi-admin.yaml",
    output: {
      target: "./src/services/api/v1-admin.ts",
      client: "react-query",
      httpClient: "fetch",
      mode: "single",
      override: {
        fetch: {
          includeHttpResponseReturnType: false,
        },
        mutator: {
          path: "./src/services/api/axios-v1-admin.ts",
          name: "customInstance",
        },
      },
    },
  },
});
