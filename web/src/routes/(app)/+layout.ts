import { goto } from "$app/navigation";
import { pb } from "$lib/client";

export async function load() {
  return new Promise((resolve) => {
    if (!pb.authStore.isValid) {
      goto("/login");
    }
    resolve({});
  });
}
