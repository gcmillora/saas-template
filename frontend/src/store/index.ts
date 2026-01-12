import { configureStore } from "@reduxjs/toolkit";
import { apiV1Public } from "../services/api/v1-public";
import { apiV1 } from "../services/api/v1";
import {
  useDispatch,
  useSelector,
  type TypedUseSelectorHook,
} from "react-redux";

export const store = configureStore({
  reducer: {
    [apiV1.reducerPath]: apiV1.reducer,
    [apiV1Public.reducerPath]: apiV1Public.reducer,
  },
  middleware: (getDefaultMiddleware) =>
    getDefaultMiddleware().concat(apiV1.middleware, apiV1Public.middleware),
});

export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;
export const useAppDispatch: () => AppDispatch = useDispatch;
export const useAppSelector: TypedUseSelectorHook<RootState> = useSelector;
