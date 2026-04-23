import { useSyncExternalStore } from "react";

export function useMediaQuery(query: string): boolean {
  const getSnapshot = () => {
    if (typeof window !== "undefined") {
      return window.matchMedia(query).matches;
    }
    return false;
  };

  const subscribe = (onStoreChange: () => void) => {
    const mql = window.matchMedia(query);
    mql.addEventListener("change", onStoreChange);
    return () => mql.removeEventListener("change", onStoreChange);
  };

  return useSyncExternalStore(subscribe, getSnapshot, () => false);
}
