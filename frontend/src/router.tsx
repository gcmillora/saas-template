import { createBrowserRouter } from "react-router";
import { AppLayout } from "./components/AppLayout";
import { AuthContextProvider } from "./contexts/authContext";
import { GuestRoute } from "./components/GuestRoute";
import { SignIn } from "./pages/SignIn";
import { SignUp } from "./pages/SignUp";
import { ForgotPassword } from "./pages/ForgotPassword";
import { ResetPassword } from "./pages/ResetPassword";
import { ErrorPage } from "./pages/ErrorPage";

function AuthenticatedLayout() {
  return (
    <AuthContextProvider>
      <AppLayout />
    </AuthContextProvider>
  );
}

export const router = createBrowserRouter([
  {
    errorElement: <ErrorPage />,
    children: [
      {
        element: <AuthenticatedLayout />,
        children: [
          {
            index: true,
            element: (
              <div className="text-muted-foreground">
                Welcome to your SaaS app. Select a page from the sidebar.
              </div>
            ),
          },
        ],
      },
      {
        path: "/signin",
        element: (
          <GuestRoute>
            <SignIn />
          </GuestRoute>
        ),
      },
      {
        path: "/signup",
        element: (
          <GuestRoute>
            <SignUp />
          </GuestRoute>
        ),
      },
      {
        path: "/forgot-password",
        element: (
          <GuestRoute>
            <ForgotPassword />
          </GuestRoute>
        ),
      },
      {
        path: "/reset-password",
        element: (
          <GuestRoute>
            <ResetPassword />
          </GuestRoute>
        ),
      },
      {
        path: "*",
        element: <ErrorPage />,
      },
    ],
  },
]);
