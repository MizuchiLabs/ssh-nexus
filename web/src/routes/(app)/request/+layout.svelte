<script lang="ts">
    import { pb, user } from "$lib/client";
    import RequestTable from "$lib/tables/RequestTable.svelte";
    import RequestSidebar from "$lib/navigation/RequestSB.svelte";

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

<div class="flex flex-row mx-auto p-4 gap-8">
    {#if canRequest()}
        <RequestSidebar />
        <slot />
    {:else}
        <RequestTable />
    {/if}
</div>
