<script lang="ts">
    import { Button } from "$lib/components/ui/button/index";
    import Table from "$lib/tables/Table.svelte";
    import { pb, user } from "$lib/client";
    import MachineModal from "$lib/modals/MachineModal.svelte";

    let open = false;

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

<MachineModal bind:open />

<div class="p-4">
    {#if canCreateMachines()}
        <Button
            variant="default"
            class="font-mono mb-4"
            on:click={() => (open = true)}
        >
            Create Machine
        </Button>
    {/if}

    <Table collection="machines" />
</div>
