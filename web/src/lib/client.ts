import PocketBase, {
  ClientResponseError,
  type AuthModel,
  type RecordModel,
  type RecordSubscription,
} from "pocketbase";
import { writable } from "svelte/store";
import { goto } from "$app/navigation";

export let pb: PocketBase;
if (import.meta.env.PROD) {
  pb = new PocketBase();
} else {
  pb = new PocketBase("http://127.0.0.1:8090");
}

export const user = writable({} as AuthModel);
export const loggedIn = writable(false);
export const isAdmin = writable(false);
export const avatar: any = writable(undefined);

export const subscribeAuth = async () => {
  if (!pb.authStore.model) return;

  if (pb.authStore.isAdmin) {
    user.set(pb.authStore.model);
    loggedIn.set(true);
    avatar.set(undefined);
    isAdmin.set(true);
  }

  if (pb.authStore.isAuthRecord) {
    try {
      // Fetch initial data
      user.set(
        await pb
          .collection("users")
          .getOne(pb.authStore.model.id, { expand: "permission,groups" }),
      );
      loggedIn.set(true);

      // Check if admin
      if (pb.authStore.model.permission) {
        const adminPermission = await pb
          .collection("permissions")
          .getOne(pb.authStore.model.permission);
        isAdmin.set(adminPermission.is_admin);
      }

      const avatarUrl = pb.files.getUrl(
        pb.authStore.model as RecordModel,
        pb.authStore.model?.avatar,
      );
      if (avatarUrl) {
        avatar.set(avatarUrl);
      }

      // Subscribe to SSE
      pb.collection("users").subscribe(
        pb.authStore.model.id,
        async (event: RecordSubscription<AuthModel>) => {
          if (event.action == "update") {
            user.set(event.record);
            avatar.set(
              pb.files.getUrl(
                event.record as RecordModel,
                event.record?.avatar,
              ),
            );
          }
          if (event.action == "delete") {
            unsubscribeAuth();
          }
        },
        { expand: "permission,groups" },
      );
    } catch (error: ClientResponseError | any) {
      console.error(`Failed to subscribe to user, error: `, error.message);
    }
  }
};

export const unsubscribeAuth = () => {
  try {
    user.set(null);
    loggedIn.set(false);
    isAdmin.set(false);
    avatar.set(undefined);
    pb.realtime.unsubscribe();
    pb.authStore.clear();
    goto("/login");
  } catch (error: ClientResponseError | any) {
    console.error(`Failed to unsubscribe from user, error: `, error.message);
  }
};
