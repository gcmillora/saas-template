import { createApi, fetchBaseQuery } from "@reduxjs/toolkit/query/react";

export const baseApiV1 = createApi({
  reducerPath: "apiV1",
  baseQuery: fetchBaseQuery({
    baseUrl: "/api/v1",
    credentials: "include",
    prepareHeaders: (headers) => {
      headers.set("Accpet", "application/json");
      headers.set("Content-Type", "application/json");
      return headers;
    },
  }),
  endpoints: () => ({}),
});

export const baseApiV1Public = createApi({
  reducerPath: "apiV1Public",
  baseQuery: fetchBaseQuery({
    baseUrl: "/api/public/v1",
    credentials: "include",
    prepareHeaders: (headers) => {
      headers.set("Accpet", "application/json");
      headers.set("Content-Type", "application/json");
      return headers;
    },
  }),
  endpoints: () => ({}),
});
