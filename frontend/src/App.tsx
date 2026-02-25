import { QueryProvider } from "./providers/QueryProvider";
import { RouterProvider } from "react-router";
import { router } from "./router";
import { ErrorBoundary } from "./components/ErrorBoundary";

function App() {
  return (
    <ErrorBoundary>
      <QueryProvider>
        <RouterProvider router={router} />
      </QueryProvider>
    </ErrorBoundary>
  );
}

export default App;
