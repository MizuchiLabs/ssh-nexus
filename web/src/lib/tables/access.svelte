<script lang="ts">
    import { Badge } from "$lib/components/ui/badge";
    import type { RecordModel } from "pocketbase";

    export let permission: RecordModel;
</script>

{#if permission}
    <div class="flex flex-row items-center justify-start gap-2">
        <div class="flex flex-col items-start justify-start gap-0.5 text-sm">
            {#if permission.access_groups}
                <Badge variant="secondary" class="bg-blue-300 text-gray-800"
                    >Groups</Badge
                >
            {/if}
            {#if permission.access_machines}
                <Badge variant="secondary" class="bg-blue-300 text-gray-800"
                    >Machines</Badge
                >
            {/if}
            {#if permission.access_users}
                <Badge variant="secondary" class="bg-blue-300 text-gray-800"
                    >Users</Badge
                >
            {/if}
        </div>

        <div class="flex flex-col items-start justify-start gap-0.5 text-sm">
            {#if permission.can_create}
                <Badge variant="secondary" class="bg-green-300 text-gray-800"
                    >Create</Badge
                >
            {/if}
            {#if permission.can_update}
                <Badge variant="secondary" class="bg-orange-300 text-gray-800"
                    >Update</Badge
                >
            {/if}
            {#if permission.can_delete}
                <Badge variant="secondary" class="bg-red-300 text-gray-800"
                    >Delete</Badge
                >
            {/if}
        </div>
    </div>
{/if}
