export const Routes = {
  home: "/home",
} as const;

export const UnauthorizedRoutes = {
  signin: "/signin",
  signup: "/signup",
} as const;

export type UnauthorizedRoutes =
  (typeof UnauthorizedRoutes)[keyof typeof UnauthorizedRoutes];

export default Routes;
