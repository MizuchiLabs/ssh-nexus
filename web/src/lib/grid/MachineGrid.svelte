<script lang="ts">
    import { pb } from "$lib/client";
    import { machines } from "$lib/subscriptions";
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
    const editMachine = (machine: RecordModel) => {
        const modal: ModalSettings = {
            type: "component",
            title: "Edit Machine",
            component: "MachineUpdate",
            meta: { machine: machine },
        };
        modals.trigger(modal);
    };

    const removeUser = async (machine: RecordModel, user: RecordModel) => {
        await pb.collection("machines").update(machine.id, {
            "users-": user.id,
        });
    };

    const removeGroup = async (machine: RecordModel, group: RecordModel) => {
        await pb.collection("machines").update(machine.id, {
            "groups-": group.id,
        });
    };

    const removeTag = async (machine: RecordModel, tag: RecordModel) => {
        await pb.collection("machines").update(machine.id, {
            "tags-": tag.id,
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
    {#each $machines as machine}
        <div
            class="card min-w-[300px] flex flex-col gap-2 divide-y dark:divide-surface-500 p-1"
        >
            <header
                class="card-header flex flex-row items-center justify-between py-2"
            >
                {#if machine.agent}
                    <iconify-icon
                        icon="fa6-solid:robot"
                        class="text-success-500 animate-pulse"
                        on:mouseover={() => (helpText = "Agent online")}
                        on:focus={() => (helpText = "Agent online")}
                        use:popup={popupHelp}
                        aria-hidden
                    />
                {:else}
                    <iconify-icon
                        icon="fa6-solid:skull"
                        on:mouseover={() => (helpText = "Agent offline")}
                        on:focus={() => (helpText = "Agent offline")}
                        use:popup={popupHelp}
                        class="text-primary-400"
                        aria-hidden
                    />
                {/if}
                <iconify-icon icon="fa6-solid:server" />
            </header>
            <section
                class="flex flex-row items-center justify-between px-4 py-2"
                on:click={() => (showDetails = !showDetails)}
                aria-hidden
            >
                <div class="flex flex-col">
                    <span class="text-lg font-bold">{machine.name}</span>
                    <span class="text-sm text-surface-400 dark:text-surface-300"
                        >{machine.host}</span
                    >
                </div>
                {#if showDetails}
                    <iconify-icon icon="fa6-regular:lightbulb" />
                {:else}
                    <iconify-icon icon="fa6-solid:lightbulb" />
                {/if}
            </section>
            {#if showDetails}
                {#if machine.expand?.users || machine.expand?.groups || machine.expand?.tags}
                    <section class="flex flex-col px-4 py-4 gap-4">
                        {#if machine.expand?.users}
                            <div
                                class="flex flex-row items-center justify-start gap-1"
                            >
                                <iconify-icon
                                    icon="fa6-solid:user"
                                    class="mr-2"
                                    on:mouseover={() => (helpText = "Users")}
                                    on:focus={() => (helpText = "Users")}
                                    use:popup={popupHelp}
                                    aria-hidden
                                />
                                {#each machine.expand.users as user}
                                    <span
                                        class="badge bg-surface-300-600-token hover:bg-primary-400"
                                        on:click={() =>
                                            removeUser(machine, user)}
                                        aria-hidden
                                    >
                                        {user.username}
                                    </span>
                                {/each}
                            </div>
                        {/if}
                        {#if machine.expand?.groups}
                            <div
                                class="flex flex-row items-center justify-start gap-1"
                            >
                                <iconify-icon
                                    icon="fa6-solid:users"
                                    class="mr-2"
                                    on:mouseover={() => (helpText = "Groups")}
                                    on:focus={() => (helpText = "Groups")}
                                    use:popup={popupHelp}
                                    aria-hidden
                                />
                                {#each machine.expand.groups as group}
                                    <span
                                        class="badge bg-surface-300-600-token hover:bg-primary-400"
                                        on:click={() =>
                                            removeGroup(machine, group)}
                                        aria-hidden
                                    >
                                        {group.name}
                                    </span>
                                {/each}
                            </div>
                        {/if}
                        {#if machine.expand?.tags}
                            <div
                                class="flex flex-row items-center justify-start gap-1"
                            >
                                <iconify-icon
                                    icon="fa6-solid:tag"
                                    class="mr-2"
                                    on:mouseover={() => (helpText = "Tags")}
                                    on:focus={() => (helpText = "Tags")}
                                    use:popup={popupHelp}
                                    aria-hidden
                                />
                                {#each machine.expand.tags as tag}
                                    <span
                                        class="badge bg-surface-300-600-token hover:bg-primary-400"
                                        on:click={() => removeTag(machine, tag)}
                                        aria-hidden
                                    >
                                        {tag.name}
                                    </span>
                                {/each}
                            </div>
                        {/if}
                    </section>
                {/if}
            {/if}
            <footer
                class="card-footer flex flex-row items-center justify-between pt-4 pb-2"
            >
                {#if machine.error}
                    <button class="btn btn-sm text-sm bg-primary-400">
                        <iconify-icon icon="fa6-solid:heart" />
                        <span>Unhealthy</span>
                    </button>
                {:else}
                    <button class="btn btn-sm text-sm variant-filled-success">
                        <iconify-icon icon="fa6-regular:heart" />
                        <span>Healthy</span>
                    </button>
                {/if}
                <button
                    class="btn btn-sm text-sm bg-surface-300-600-token"
                    on:click={() => editMachine(machine)}
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
