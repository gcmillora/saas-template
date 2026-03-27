import { useContext } from "react";
import { Home, Settings, LogOut, ChevronsUpDown } from "lucide-react";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { useNavigate, useLocation } from "react-router";

import {
  Sidebar,
  SidebarContent,
  SidebarGroup,
  SidebarGroupContent,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarHeader,
  SidebarFooter,
  SidebarRail,
  SidebarSeparator,
} from "@/components/ui/sidebar";
import { usePostApiV1Signout } from "@/services/api/v1-public";
import { AuthContext } from "@/contexts/authContext";
import { queryClient } from "@/providers/QueryProvider";

export function AppSidebar() {
  const navigate = useNavigate();
  const location = useLocation();
  const signout = usePostApiV1Signout();
  const { user } = useContext(AuthContext);

  const userInitials = user
    ? [user.firstName, user.lastName]
        .filter(Boolean)
        .map((n) => n![0].toUpperCase())
        .join("") || user.email[0].toUpperCase()
    : "?";

  const userDisplayName = user
    ? [user.firstName, user.lastName].filter(Boolean).join(" ") || user.email
    : "";

  const handleSignout = async () => {
    try {
      await signout.mutateAsync();
    } catch {
      // continue even if request fails
    }
    queryClient.clear();
    navigate("/signin");
  };

  return (
    <Sidebar collapsible="icon">
      <SidebarHeader className="pb-0">
        <SidebarMenu>
          <SidebarMenuItem>
            <SidebarMenuButton size="lg" asChild>
              <a href="/" className="group/brand">
                <div className="flex aspect-square size-9 shrink-0 items-center justify-center rounded-xl bg-sidebar-primary text-sidebar-primary-foreground transition-transform duration-200 group-hover/brand:scale-105">
                  <Home className="size-4" />
                </div>
                <div className="flex flex-col gap-0.5 leading-none">
                  <span className="text-sm font-semibold tracking-tight">
                    SaaS Template
                  </span>
                  <span className="text-[10px] text-sidebar-foreground/50">
                    v1.0.0
                  </span>
                </div>
              </a>
            </SidebarMenuButton>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarHeader>

      <SidebarSeparator className="my-3" />

      <SidebarContent>
        <SidebarGroup className="pt-0">
          <SidebarGroupContent>
            <SidebarMenu>
              <SidebarMenuItem>
                <SidebarMenuButton asChild isActive={location.pathname === "/"}>
                  <a href="/">
                    <Home className="size-4" />
                    <span>Home</span>
                  </a>
                </SidebarMenuButton>
              </SidebarMenuItem>
              <SidebarMenuItem>
                <SidebarMenuButton
                  asChild
                  isActive={location.pathname === "/settings"}
                >
                  <a href="/settings">
                    <Settings className="size-4" />
                    <span>Settings</span>
                  </a>
                </SidebarMenuButton>
              </SidebarMenuItem>
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>
      </SidebarContent>

      <SidebarFooter>
        <SidebarSeparator className="mb-2" />
        <SidebarMenu>
          <SidebarMenuItem>
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <SidebarMenuButton
                  size="lg"
                  className="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
                >
                  <Avatar className="size-8 rounded-lg">
                    <AvatarFallback className="rounded-lg bg-sidebar-accent text-sidebar-accent-foreground text-xs font-medium">
                      {userInitials}
                    </AvatarFallback>
                  </Avatar>
                  <div className="grid flex-1 text-left text-sm leading-tight">
                    <span className="truncate font-semibold">
                      {userDisplayName}
                    </span>
                    <span className="truncate text-xs text-sidebar-foreground/50">
                      {user?.email}
                    </span>
                  </div>
                  <ChevronsUpDown className="ml-auto size-4" />
                </SidebarMenuButton>
              </DropdownMenuTrigger>
              <DropdownMenuContent
                className="w-(--radix-dropdown-menu-trigger-width) min-w-56 rounded-lg"
                side="top"
                align="end"
                sideOffset={4}
              >
                <DropdownMenuItem onClick={handleSignout}>
                  <LogOut className="mr-2 size-4" />
                  Sign out
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarFooter>
      <SidebarRail />
    </Sidebar>
  );
}
