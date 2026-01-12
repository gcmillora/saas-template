import { Provider } from "react-redux";
import { store } from "./store";
import { RouterProvider } from "react-router";
import { router } from "./router";

function App() {
  return (
    <Provider store={store}>
      <RouterProvider router={router} />
    </Provider>
  );
}

export default App;
