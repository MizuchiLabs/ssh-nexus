<script lang="ts">
    import { getToastStore, popup } from "@skeletonlabs/skeleton";
    import { pb, user } from "$lib/client";
    import { requests } from "$lib/subscriptions";
    import type { RecordModel } from "pocketbase";
    import { showToast } from "$lib/utils/Toast";

    const toastStore = getToastStore();
    const deleteRequest = async (request: RecordModel) => {
        await pb.collection("requests").delete(request.id);
        showToast(toastStore, `Deleted request 👌`, "success");
    };

    const clearRequests = async () => {
        let requests = await pb.collection("requests").getFullList({
            filter: pb.filter("user.id = {:user_id}", {
                user_id: pb.authStore.model?.id,
            }),
        });

        for (let request of requests) {
            await pb.collection("requests").delete(request.id);
        }
        showToast(toastStore, `Cleared all requests 👌`, "success");
    };

    const canRequest = () => {
        if ($user?.expand?.permission?.is_admin || pb.authStore.isAdmin) {
            return false;
        }
        if (
            $user?.expand?.permission?.access_machines &&
            $user?.expand?.permission?.can_update
        ) {
            return false;
        }
        return true;
    };
</script>

<div
    use:popup={{
        event: "focus-blur",
        target: "popupRequests",
        placement: "bottom",
    }}
>
    {#if canRequest()}
        <button
            class="btn-icon btn-icon-sm relative variant-soft-primary hover:brightness-50 mr-2"
        >
            <iconify-icon icon="fa6-solid:bell" class="hover:brightness-50" />
            <span class="absolute badge-icon -top-1 -right-1 shadow-none"
                >{$requests.length}
            </span>
        </button>
        <div
            class="card z-10 rounded-md dark:bg-surface-500 divide-y divide-gray-300"
            data-popup="popupRequests"
            tabindex="-1"
        >
            <div class="block px-4 py-2 text-sm font-semibold">
                Pending Requests
            </div>
            <ul
                class="flex flex-col p-2 gap-2 text-sm text-gray-700 dark:text-gray-200"
            >
                {#each $requests as request}
                    <li
                        class="flex flex-row justify-between items-center hover:bg-primary-500/20 hover:line-through rounded-full p-2"
                        on:click={() => deleteRequest(request)}
                        aria-hidden
                    >
                        {#if request.machine}
                            <button
                                class="flex flex-row items-center justify-between gap-2"
                            >
                                <iconify-icon icon="fa6-solid:server" />

                                {request.expand?.machine.name}
                            </button>
                        {/if}
                        {#if request.group}
                            <button class="flex flex-row items-center gap-2">
                                <iconify-icon icon="fa6-solid:users" />
                                {request.expand?.group.name}
                            </button>
                        {/if}
                    </li>
                {/each}
                <button
                    on:click={clearRequests}
                    class="btn btn-sm bg-surface-500 hover:bg-surface-400 mt-1"
                >
                    Clear All
                </button>
            </ul>
        </div>
    {:else}
        <a
            href="/request/machines/"
            class="btn-icon btn-icon-sm relative variant-soft-primary hover:brightness-50 mr-2"
        >
            <iconify-icon icon="fa6-solid:bell" class="hover:brightness-50" />
            <span class="absolute badge-icon -top-1 -right-1 shadow-none"
                >{$requests.length}
            </span>
        </a>
    {/if}
</div>
