import { Outlet } from "react-router";
import { AppSidebar } from "./app-sidebar";
import {
  SidebarInset,
  SidebarProvider,
  SidebarTrigger,
} from "@/components/ui/sidebar";
import { useMediaQuery } from "@/hooks/use-media-query";

export function AppLayout() {
  const isDesktop = useMediaQuery("(min-width: 1024px)");

  return (
    <SidebarProvider defaultOpen={isDesktop}>
      <AppSidebar />
      <SidebarInset>
        <header className="flex h-12 shrink-0 items-center gap-2 border-b px-4">
          <SidebarTrigger className="-ml-1" />
        </header>
        <main className="flex flex-1 flex-col gap-4 p-4">
          <Outlet />
        </main>
      </SidebarInset>
    </SidebarProvider>
  );
}
