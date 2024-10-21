export interface Route {
  name: string;
  path: string;
  icon: string;
  admin: boolean;
}

export const baseRoutes: Route[] = [
  {
    name: "Home",
    path: "/",
    icon: "fa6-solid:house",
    admin: false,
  },
  {
    name: "Machines",
    path: "/machines",
    icon: "fa6-solid:server",
    admin: false,
  },
  {
    name: "Users",
    path: "/users",
    icon: "fa6-solid:user",
    admin: true,
  },
  {
    name: "Groups",
    path: "/groups",
    icon: "fa6-solid:users",
    admin: false,
  },
  {
    name: "Settings",
    path: "/settings/general",
    icon: "fa6-solid:gear",
    admin: true,
  },
];
