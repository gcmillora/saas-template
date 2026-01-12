/** @type (import("@rtk-query/codegen/opeapi").ConfigFile) */
const config = {
  schemaFile: "../../../../backend/openapi.yaml",
  apiFile: "./client.ts",
  apiImport: "baseApiV1",
  outputFile: "./v1.ts",
  exportName: "apiV1",
  hooks: { queries: true, lazyQueries: true, mutations: true },
};

module.exports = config;
