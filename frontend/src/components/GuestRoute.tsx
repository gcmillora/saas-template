import { type PropsWithChildren } from "react";
import { Navigate } from "react-router";
import { useGetApiV1User } from "@/services/api/v1";

export function GuestRoute({ children }: PropsWithChildren) {
  const { data, isLoading } = useGetApiV1User({ query: { retry: false } });

  if (isLoading) return null;

  if (data) {
    return <Navigate to="/" replace />;
  }

  return <>{children}</>;
}
