import { createBrowserRouter } from "react-router";
import { AppLayout } from "./components/AppLayout";
import { SignIn } from "./pages/SignIn";
import { SignUp } from "./pages/SignUp";

export const router = createBrowserRouter([
  {
    element: <AppLayout />,
    children: [
      // routes
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
