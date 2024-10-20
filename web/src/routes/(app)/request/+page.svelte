<script lang="ts">
    import { pb, user } from "$lib/client";
    import Table from "$lib/tables/Table.svelte";
    import { Button } from "$lib/components/ui/button/index";
    import RequestModal from "$lib/modals/RequestModal.svelte";

    let open = false;

    // User is admin and doesn't need to request permissions
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

<RequestModal bind:open />

<div class="p-4">
    <div class="flex flex-row justify-between mb-4">
        {#if canRequest()}
            <Button
                variant="default"
                class="font-mono"
                on:click={() => (open = true)}
            >
                Create Request
            </Button>
        {/if}
    </div>

    <Table collection="requests" />
</div>
