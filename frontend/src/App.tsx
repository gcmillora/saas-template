import { QueryProvider } from "./providers/QueryProvider";
import { RouterProvider } from "react-router";
import { router } from "./router";
import { ErrorBoundary } from "./components/ErrorBoundary";
import { Toaster } from "@/components/ui/sonner";

function App() {
  return (
    <ErrorBoundary>
      <QueryProvider>
        <RouterProvider router={router} />
        <Toaster richColors position="top-right" />
      </QueryProvider>
    </ErrorBoundary>
  );
}

export default App;
