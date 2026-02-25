import { createContext, type JSX, type PropsWithChildren } from "react";
import { useGetApiV1User, type BaseUser } from "../services/api/v1";
import { UnauthorizedRoutes } from "../Routes";

export type AuthContextInterface = {
  user: BaseUser | null;
};

export const AuthContext = createContext<AuthContextInterface>({ user: null });

export const AuthContextProvider = (props: PropsWithChildren): JSX.Element => {
  const { data, error, isLoading } = useGetApiV1User();

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
    window.location.replace("/signin");
    return <></>;
  }

  const user = data?.status === 200 ? data.data : null;
  const authCtx: AuthContextInterface = {
    user: user ?? null,
  };

  return (
    <AuthContext.Provider value={authCtx}>
      {props.children}
    </AuthContext.Provider>
  );
};
