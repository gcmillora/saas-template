import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import type { JSX, PropsWithChildren } from "react";

const queryClient = new QueryClient();

export const QueryProvider = ({ children }: PropsWithChildren): JSX.Element => (
  <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
);
