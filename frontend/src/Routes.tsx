export const Routes = {
  home: "/",
  forgotPassword: "/forgot-password",
  resetPassword: "/reset-password",
} as const;

export const UnauthorizedRoutes = {
  signin: "/signin",
  signup: "/signup",
  forgotPassword: "/forgot-password",
  resetPassword: "/reset-password",
} as const;

export type UnauthorizedRoutes =
  (typeof UnauthorizedRoutes)[keyof typeof UnauthorizedRoutes];

export default Routes;
