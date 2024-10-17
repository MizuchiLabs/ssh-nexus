<script lang="ts">
    import { pb } from "$lib/client";
    import { users } from "$lib/subscriptions";
    import {
        getModalStore,
        popup,
        type ModalSettings,
        type PopupSettings,
    } from "@skeletonlabs/skeleton";
    import type { RecordModel } from "pocketbase";

    let showDetails = false;

    // Modals
    const modals = getModalStore();
    const editUser = (user: RecordModel) => {
        const modal: ModalSettings = {
            type: "component",
            title: "Edit User",
            component: "UserUpdate",
            meta: { user: user },
        };
        modals.trigger(modal);
    };

    const removeGroup = async (user: RecordModel, group: RecordModel) => {
        await pb.collection("users").update(user.id, {
            "groups-": group.id,
        });
    };

    // Helper popup
    let helpText = "";
    const popupHelp: PopupSettings = {
        event: "hover",
        target: "popupHelp",
        placement: "bottom",
    };
</script>

<div
    class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 items-center justify-center gap-4"
>
    {#each $users as user}
        <div
            class="card min-w-[300px] flex flex-col gap-2 divide-y divide-surface-500 p-1"
        >
            <header
                class="card-header flex flex-row items-center justify-between py-2"
            >
                <span class="text-sm text-surface-600-300-token"
                    >{user.principal.slice(0, 8)}...</span
                >
                {#if user.expand?.permission?.is_admin}
                    <iconify-icon
                        icon="fa6-solid:crown"
                        on:mouseover={() => (helpText = "Admin")}
                        on:focus={() => (helpText = "Admin")}
                        use:popup={popupHelp}
                        aria-hidden
                    />
                {:else}
                    <iconify-icon
                        icon="fa6-solid:user"
                        on:mouseover={() => (helpText = "Member")}
                        on:focus={() => (helpText = "Member")}
                        use:popup={popupHelp}
                        aria-hidden
                    />
                {/if}
            </header>
            <section
                class="flex flex-row items-center justify-between px-4 py-2"
                on:click={() => (showDetails = !showDetails)}
                aria-hidden
            >
                <div class="flex flex-col">
                    <span class="text-lg font-bold">{user.name}</span>
                    <span class="text-sm text-surface-600-300-token"
                        >{user.email}</span
                    >
                </div>
                {#if showDetails}
                    <iconify-icon icon="fa6-regular:lightbulb" />
                {:else}
                    <iconify-icon icon="fa6-solid:lightbulb" />
                {/if}
            </section>
            {#if showDetails && user.expand?.groups}
                <section class="flex flex-col px-4 py-4 gap-4">
                    <div class="flex flex-row items-center justify-start gap-1">
                        <iconify-icon
                            icon="fa6-solid:users"
                            class="mr-2"
                            on:mouseover={() => (helpText = "Groups")}
                            on:focus={() => (helpText = "Groups")}
                            use:popup={popupHelp}
                            aria-hidden
                        />
                        {#each user.expand.groups as group}
                            <span
                                class="badge bg-surface-300-600-token hover:bg-primary-400"
                                on:click={() => removeGroup(user, group)}
                                aria-hidden
                            >
                                {group.name}
                            </span>
                        {/each}
                    </div>
                </section>
            {/if}
            <footer
                class="card-footer flex flex-row items-center justify-between pt-4 pb-2"
            >
                {#if user.verified}
                    <button class="btn btn-sm text-sm variant-filled-success">
                        <iconify-icon icon="fa6-solid:circle-check" />
                        <span>Verified</span>
                    </button>
                {:else}
                    <button class="btn btn-sm text-sm variant-filled-error">
                        <iconify-icon icon="fa6-solid:xmark" />
                        <span>Unverified</span>
                    </button>
                {/if}
                <button
                    class="btn btn-sm text-sm bg-surface-300-600-token"
                    on:click={() => editUser(user)}
                >
                    <iconify-icon
                        icon="fa6-solid:pen"
                        class="mr-2"
                        width="12"
                    />
                    Edit
                </button>
            </footer>
        </div>
    {/each}
    <div class="card p-4 variant-filled-surface" data-popup="popupHelp">
        <div><p>{helpText}</p></div>
        <div class="arrow variant-filled-surface" />
    </div>
</div>
