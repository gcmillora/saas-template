import { createContext, type JSX, type PropsWithChildren } from "react";
import { useGetUserQuery, type BaseUser } from "../services/api/v1";
import { UnauthorizedRoutes } from "../Routes";

export type AuthContextInterface = {
  user: BaseUser | null;
};

export const AuthContext = createContext<AuthContextInterface>({ user: null });

export const AuthContextProvider = (props: PropsWithChildren): JSX.Element => {
  const { data, error, isLoading } = useGetUserQuery();

  if (
    Object.values(UnauthorizedRoutes).includes(
      window.location.pathname as UnauthorizedRoutes,
    )
  ) {
    return (
      <AuthContext.Provider
        value={{
          user: null,
        }}
      >
        {props.children}
      </AuthContext.Provider>
    );
  }

  // biome-ignore lint/complexity/noUselessFragments: ignore
  if (isLoading) return <></>;

  if (error) {
    console.log(error);
    console.log("Not authenticated");
    return <></>;
  }
  const authCtx: AuthContextInterface = {
    user: data ?? null,
  };

  return (
    <AuthContext.Provider value={authCtx}>
      {props.children}
    </AuthContext.Provider>
  );
};
