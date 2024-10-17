<script lang="ts">
    import { pb, user } from "$lib/client";
    import { machines, requests } from "$lib/subscriptions";
    import { showToast } from "$lib/utils/Toast";
    import { getToastStore } from "@skeletonlabs/skeleton";

    let toastStore = getToastStore();

    // Filter machines and groups which the user has already access too
    let fMachines = $machines;

    $: fMachines = $machines.filter(
        (machine) =>
            !machine.users?.includes($user?.id) &&
            !$requests.some(
                (request) =>
                    request.machine === machine.id &&
                    request.user === $user?.id,
            ),
    );

    const requestMachines = async () => {
        const selected = fMachines.filter((m) => m.selected);
        fMachines = fMachines.map((m) => ({ ...m, selected: false }));
        for (const machine of selected) {
            try {
                await pb.collection("requests").create({
                    user: pb.authStore.model?.id,
                    machine: machine.id,
                });
            } catch (error) {
                console.log(error);
            }
        }
        showToast(
            toastStore,
            `Requested access to ${selected.length} machines 👌`,
            "success",
        );
    };
</script>

<div class="flex flex-col w-full">
    <div class="input-group input-group-divider grid-cols-[auto_1fr_auto]">
        <div class="input-group-shim">
            <iconify-icon icon="fa6-solid:magnifying-glass" />
        </div>
        <input type="search" placeholder="Search..." />
    </div>
    <div class="flex flex-row justify-start gap-2 my-4">
        {#if fMachines.filter((m) => m.selected).length > 0}
            <button
                class="btn btn-sm variant-filled-success"
                on:click={requestMachines}
            >
                <iconify-icon icon="fa6-solid:plus" class="mr-1" />
                Request Access
            </button>
            <button
                class="btn btn-sm variant-filled-error"
                on:click={() =>
                    (fMachines = fMachines.map((m) => ({
                        ...m,
                        selected: false,
                    })))}
            >
                <iconify-icon icon="fa6-solid:minus" class="mr-1" />
                Clear
            </button>
        {/if}
    </div>
    <div
        class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 items-center justify-center gap-4"
    >
        {#each fMachines as machine}
            <div
                class="card flex flex-col gap-2 divide-y divide-surface-500 p-1"
            >
                <section
                    class="flex flex-row items-center justify-start gap-4 px-4 py-2"
                    on:click={() => (machine.selected = !machine.selected)}
                    aria-hidden
                >
                    <iconify-icon
                        icon="fa6-solid:server"
                        class="rounded-full p-3"
                        class:bg-surface-500={!machine.selected}
                        class:bg-green-600={machine.selected}
                    />
                    <div class="flex flex-col">
                        <span class="text-lg font-bold">{machine.name}</span>
                        <span class="text-sm text-surface-300"
                            >{machine.host}</span
                        >
                    </div>
                </section>
            </div>
        {/each}
    </div>
</div>
