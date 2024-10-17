<script lang="ts">
    import MachineTable from "$lib/tables/MachineTable.svelte";
    import MachineGrid from "$lib/grid/MachineGrid.svelte";
    import { getModalStore, type ModalSettings } from "@skeletonlabs/skeleton";
    import { pb, user } from "$lib/client";

    let gridview = pb.authStore.model?.settings?.gridview ?? true;

    const modals = getModalStore();
    const createMachine = () => {
        const modal: ModalSettings = {
            type: "component",
            title: "New Machine",
            component: "MachineCreate",
        };
        modals.trigger(modal);
    };

    const setView = async (select: boolean) => {
        gridview = select;
        if (pb.authStore.isAdmin) return;

        await pb.collection("users").update(pb.authStore.model?.id, {
            settings: {
                ...pb.authStore.model?.settings,
                gridview: gridview,
            },
        });
    };

    const canCreateMachines = () => {
        if ($user?.expand?.permission?.is_admin || pb.authStore.isAdmin) {
            return true;
        }
        if (
            $user?.expand?.permission?.access_machines &&
            $user?.expand?.permission?.can_create
        ) {
            return true;
        }
        return false;
    };
</script>

<div class="p-4">
    <div class="flex flex-row justify-between mb-8">
        {#if canCreateMachines()}
            <button
                class="btn btn-sm bg-primary-400 font-mono"
                on:click={createMachine}
            >
                Add Machine
            </button>
        {:else}
            <div></div>
        {/if}

        <div class="btn-group variant-filled-surface">
            <button
                on:click={() => setView(true)}
                class:bg-primary-400={gridview}
            >
                <iconify-icon icon="fa6-solid:border-all" />
            </button>
            <button
                on:click={() => setView(false)}
                class:bg-primary-400={!gridview}
            >
                <iconify-icon icon="fa6-solid:table-list" />
            </button>
        </div>
    </div>

    {#if gridview}
        <MachineGrid />
    {:else}
        <MachineTable />
    {/if}
</div>
