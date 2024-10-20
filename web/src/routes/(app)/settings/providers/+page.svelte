<script lang="ts">
    import { pb, user } from "$lib/client";
    import Table from "$lib/tables/Table.svelte";
    import { Button } from "$lib/components/ui/button/index";
    import ProviderModal from "$lib/modals/ProviderModal.svelte";

    let open = false;

    const canCreateProvider = () => {
        if ($user?.expand?.permission?.is_admin || pb.authStore.isAdmin) {
            return true;
        }
        return false;
    };
</script>

<ProviderModal bind:open />

<div class="p-4">
    {#if canCreateProvider()}
        <Button
            variant="default"
            class="font-mono"
            on:click={() => (open = true)}
        >
            Create Provider
        </Button>
    {/if}

    <Table collection="providers" />
</div>
