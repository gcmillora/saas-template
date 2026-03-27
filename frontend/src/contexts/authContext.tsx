import { createContext, type JSX, type PropsWithChildren } from "react";
import { Navigate } from "react-router";
import { useGetApiV1User, type BaseUser } from "../services/api/v1";

export type AuthContextInterface = {
  user: BaseUser | null;
};

export const AuthContext = createContext<AuthContextInterface>({ user: null });

export const AuthContextProvider = (props: PropsWithChildren): JSX.Element => {
  const { data, error, isLoading } = useGetApiV1User({ query: { retry: false } });

  if (isLoading) return <></>;

  if (error) {
    return <Navigate to="/signin" replace />;
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
