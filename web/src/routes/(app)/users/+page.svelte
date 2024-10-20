<script lang="ts">
    import { Button } from "$lib/components/ui/button/index";
    import Table from "$lib/tables/Table.svelte";
    import { pb, user } from "$lib/client";
    import UserModal from "$lib/modals/UserModal.svelte";

    let open = false;

    const canCreateUsers = () => {
        if ($user?.expand?.permission?.is_admin || pb.authStore.isAdmin) {
            return true;
        }
        if (
            $user?.expand?.permission?.access_users &&
            $user?.expand?.permission?.can_create
        ) {
            return true;
        }
        return false;
    };
</script>

<UserModal bind:open />

<div class="p-4">
    {#if canCreateUsers()}
        <Button
            variant="default"
            class="font-mono mb-4"
            on:click={() => (open = true)}
        >
            Create User
        </Button>
    {/if}

    <Table collection="users" />
</div>
