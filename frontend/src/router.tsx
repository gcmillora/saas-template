import { createBrowserRouter } from "react-router";
import { AppLayout } from "./components/AppLayout";
import { AuthContextProvider } from "./contexts/authContext";
import { SignIn } from "./pages/SignIn";
import { SignUp } from "./pages/SignUp";

function AuthenticatedLayout() {
  return (
    <AuthContextProvider>
      <AppLayout />
    </AuthContextProvider>
  );
}

export const router = createBrowserRouter([
  {
    element: <AuthenticatedLayout />,
    children: [
      {
        index: true,
        element: <div className="text-muted-foreground">Welcome to your SaaS app. Select a page from the sidebar.</div>,
      },
    ],
  },
  {
    path: "/signup",
    element: <SignUp />,
  },
  {
    path: "/signin",
    element: <SignIn />,
  },
]);
