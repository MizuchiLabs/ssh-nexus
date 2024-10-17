import type { AuthModel, RecordModel } from "pocketbase";
import { pb } from "./client";

export interface Permission {
  id: string;
  name: string;
  description: string;
  is_admin: boolean;
  can_create: boolean;
  can_update: boolean;
  can_delete: boolean;
  access_users: boolean;
  access_groups: boolean;
  access_machines: boolean;
  users: string[];
  groups: string[];
  machines: string[];
}

export interface AccessLevel {
  can_create_machines: boolean;
  can_update_machines: boolean;
  can_delete_machines: boolean;
  can_update_users: boolean;
  can_delete_users: boolean;
  can_create_groups: boolean;
  can_update_groups: boolean;
  can_delete_groups: boolean;
}

export async function getPermission() {
  if (!pb.authStore.model) return;

  if (pb.authStore.isAuthRecord) {
    const permission: Permission = await pb
      .collection("permissions")
      .getOne(pb.authStore.model?.permission);
    return permission;
  }
  if (pb.authStore.isAdmin) {
    let permission: Permission = {
      id: "",
      name: "admin",
      description: "Default admin permission",
      is_admin: true,
      can_create: true,
      can_update: true,
      can_delete: true,
      access_users: true,
      access_groups: true,
      access_machines: true,
      users: [],
      groups: [],
      machines: [],
    };
    return permission;
  }
  return undefined;
}

export async function getAccessLevel(record?: RecordModel) {
  if (!pb.authStore.model) return;
  const permission = await getPermission();

  let accessLevel: AccessLevel = {
    can_create_machines: false,
    can_update_machines: false,
    can_delete_machines: false,
    can_update_users: false,
    can_delete_users: false,
    can_create_groups: false,
    can_update_groups: false,
    can_delete_groups: false,
  };

  if (permission) {
    if (permission.is_admin) {
      accessLevel.can_create_machines = true;
      accessLevel.can_update_machines = true;
      accessLevel.can_delete_machines = true;
      accessLevel.can_update_users = true;
      accessLevel.can_delete_users = true;
      accessLevel.can_create_groups = true;
      accessLevel.can_update_groups = true;
      accessLevel.can_delete_groups = true;
      return accessLevel;
    }
    switch (record?.collection) {
      case "users":
        if (permission.access_users) {
          accessLevel.can_update_users = true;
          accessLevel.can_delete_users = true;
        }
        break;
      case "machines":
        if (permission.access_machines && record) {
          let inList = permission.machines.includes(record.id);
          accessLevel.can_create_machines = permission.can_create;
          accessLevel.can_update_machines = permission.can_update && inList;
          accessLevel.can_delete_machines = permission.can_delete && inList;
        }
        break;
      case "groups":
        if (permission.access_groups && record) {
          let inList = permission.groups.includes(record.id);
          accessLevel.can_create_groups = permission.can_create;
          accessLevel.can_update_groups = permission.can_update && inList;
          accessLevel.can_delete_groups = permission.can_delete && inList;
        }
        break;
      default:
        break;
    }
  }
  return accessLevel;
}
