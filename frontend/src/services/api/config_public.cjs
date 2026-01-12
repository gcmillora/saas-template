/** @type (import("@rtk-query/codegen/opeapi").ConfigFile) */
const config = {
  schemaFile: "../../../../backend/openapi-public.yaml",
  apiFile: "./client.ts",
  apiImport: "baseApiV1Public",
  outputFile: "./v1-public.ts",
  exportName: "apiV1Public",
  hooks: { queries: true, lazyQueries: true, mutations: true },
};

module.exports = config;
