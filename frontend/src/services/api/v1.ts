import { baseApiV1 as api } from "./client";
const injectedRtkApi = api.injectEndpoints({
  endpoints: (build) => ({
    getUser: build.query<GetUserApiResponse, GetUserApiArg>({
      query: () => ({ url: `/user` }),
    }),
  }),
  overrideExisting: false,
});
export { injectedRtkApi as apiV1 };
export type GetUserApiResponse =
  /** status 200 Current user information */ BaseUser;
export type GetUserApiArg = void;
export type BaseUser = {
  email: string;
  id: string;
  firstName?: string | null;
  lastName?: string | null;
};
export type Error = {
  message?: string;
};
export const { useGetUserQuery, useLazyGetUserQuery } = injectedRtkApi;
