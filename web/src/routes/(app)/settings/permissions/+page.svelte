<script lang="ts">
    import { pb, user } from "$lib/client";
    import Table from "$lib/tables/Table.svelte";
    import { Button } from "$lib/components/ui/button/index";
    import PermissionModal from "$lib/modals/PermissionModal.svelte";

    let open = false;

    const canCreatePermission = () => {
        if ($user?.expand?.permission?.is_admin || pb.authStore.isAdmin) {
            return true;
        }
        return false;
    };
</script>

<PermissionModal bind:open />

<div class="p-4">
    {#if canCreatePermission()}
        <Button
            variant="default"
            class="font-mono"
            on:click={() => (open = true)}
        >
            Create Permission
        </Button>
    {/if}

    <Table collection="permissions" />
</div>
