import { baseApiV1Public as api } from "./client";
const injectedRtkApi = api.injectEndpoints({
  endpoints: (build) => ({
    getHealth: build.query<GetHealthApiResponse, GetHealthApiArg>({
      query: () => ({ url: `/health` }),
    }),
    postSignin: build.mutation<PostSigninApiResponse, PostSigninApiArg>({
      query: (queryArg) => ({
        url: `/signin`,
        method: "POST",
        body: queryArg.body,
      }),
    }),
    postSignup: build.mutation<PostSignupApiResponse, PostSignupApiArg>({
      query: (queryArg) => ({
        url: `/signup`,
        method: "POST",
        body: queryArg.body,
      }),
    }),
  }),
  overrideExisting: false,
});
export { injectedRtkApi as apiV1Public };
export type GetHealthApiResponse = unknown;
export type GetHealthApiArg = void;
export type PostSigninApiResponse = unknown;
export type PostSigninApiArg = {
  /** User credentials */
  body: {
    email?: string;
    password?: string;
    tenant_id?: string;
  };
};
export type PostSignupApiResponse = unknown;
export type PostSignupApiArg = {
  /** User signup credentials */
  body: {
    email?: string;
    password?: string;
    confirm_password?: string;
    tenant_id?: string;
  };
};
export const {
  useGetHealthQuery,
  useLazyGetHealthQuery,
  usePostSigninMutation,
  usePostSignupMutation,
} = injectedRtkApi;
