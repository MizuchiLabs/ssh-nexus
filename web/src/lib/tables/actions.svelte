<script lang="ts">
    import { pb } from "$lib/client";
    import { Button } from "$lib/components/ui/button";
    import { Eye, X } from "lucide-svelte";
    import type { ClientResponseError } from "pocketbase";
    import { toast } from "svelte-sonner";
    import type { RecordModel } from "pocketbase";
    import MachineModal from "../modals/MachineModal.svelte";
    import UserModal from "../modals/UserModal.svelte";
    import GroupModal from "../modals/GroupModal.svelte";
    import PermissionModal from "../modals/PermissionModal.svelte";
    import ProviderModal from "../modals/ProviderModal.svelte";

    export let record: RecordModel;
    export let collection: string;

    let open = false;
    const deleteAction = async () => {
        try {
            await pb.collection(collection).delete(record.id);
            toast.success(`Deleted ${record.name}`, {
                description: `Successfully deleted from ${record.collectionName}.`,
                duration: 3000,
            });
        } catch (error: ClientResponseError | any) {
            toast.error(error.data?.message || "Something went wrong.");
        }
    };
</script>

<div class="flex items-center gap-1 dark:text-black">
    <Button
        variant="ghost"
        class="h-8 w-8 rounded-full bg-green-400"
        size="icon"
        on:click={() => (open = true)}
    >
        <Eye size="1rem" />
    </Button>
    <Button
        variant="ghost"
        class="h-8 w-8 rounded-full bg-red-400"
        size="icon"
        on:click={deleteAction}
    >
        <X size="1rem" />
    </Button>
</div>

<!-- Modals -->
{#if collection === "machines"}
    <MachineModal bind:open machine={record} />
{/if}

{#if collection === "users"}
    <UserModal bind:open user={record} />
{/if}

{#if collection === "groups"}
    <GroupModal bind:open group={record} />
{/if}

{#if collection === "permissions"}
    <PermissionModal bind:open permission={record} />
{/if}

{#if collection === "providers"}
    <ProviderModal bind:open provider={record} />
{/if}
