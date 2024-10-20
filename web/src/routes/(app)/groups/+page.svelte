<script lang="ts">
  import { Button } from "$lib/components/ui/button/index";
  import Table from "$lib/tables/Table.svelte";
  import { pb, user } from "$lib/client";
  import GroupModal from "$lib/modals/GroupModal.svelte";

  let open = false;

  const canCreateGroups = () => {
    if ($user?.expand?.permission?.is_admin || pb.authStore.isAdmin) {
      return true;
    }
    if (
      $user?.expand?.permission?.access_groups &&
      $user?.expand?.permission?.can_create
    ) {
      return true;
    }
    return false;
  };
</script>

<GroupModal bind:open />

<div class="p-4">
  {#if canCreateGroups()}
    <Button
      variant="default"
      class="font-mono mb-4"
      on:click={() => (open = true)}
    >
      Create Group
    </Button>
  {/if}

  <Table collection="groups" />
</div>
