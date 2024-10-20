<script lang="ts">
    import * as HoverCard from "$lib/components/ui/hover-card/index.js";
    import { Badge } from "$lib/components/ui/badge";
    import { permissions } from "$lib/subscriptions";

    export let id: string;

    $: details = $permissions.filter((p) => p.id === id)[0];
</script>

{#if details}
    <HoverCard.Root>
        <HoverCard.Trigger>
            <Badge variant="secondary" color="primary" class="mr-0.5 mb-0.5">
                {details.name}
            </Badge>
        </HoverCard.Trigger>
        <HoverCard.Content class="w-100 shadow-lg rounded-lg">
            <div class="space-y-2">
                <!-- Permission Name -->
                <h4 class="text-sm font-semibold">{details.name}</h4>

                <!-- Permission Access Info -->
                <div class="flex flex-col justify-start gap-1 py-2">
                    {#if details.is_admin}
                        <div class="text-sm">
                            <Badge
                                variant="secondary"
                                class="bg-red-300 text-gray-800">Admin</Badge
                            >
                        </div>
                    {/if}
                    <div class="text-sm">
                        {#if details.access_groups}
                            <Badge
                                variant="secondary"
                                class="bg-green-300 text-gray-800"
                                >Access Groups</Badge
                            >
                        {/if}
                        {#if details.access_machines}
                            <Badge
                                variant="secondary"
                                class="bg-green-300 text-gray-800"
                                >Access Machines</Badge
                            >
                        {/if}
                        {#if details.access_users}
                            <Badge
                                variant="secondary"
                                class="bg-green-300 text-gray-800"
                                >Access Users</Badge
                            >
                        {/if}
                    </div>
                    <div class="text-sm">
                        {#if details.can_create}
                            <Badge
                                variant="secondary"
                                class="bg-orange-300 text-gray-800"
                            >
                                Create
                            </Badge>
                        {/if}
                        {#if details.can_update}
                            <Badge
                                variant="secondary"
                                class="bg-orange-300 text-gray-800"
                            >
                                Update
                            </Badge>
                        {/if}
                        {#if details.can_delete}
                            <Badge
                                variant="secondary"
                                class="bg-orange-300 text-gray-800"
                            >
                                Delete
                            </Badge>
                        {/if}
                    </div>
                </div>

                <!-- Associated Groups -->
                {#if details.expand && details.expand?.groups.length > 0}
                    <div class="text-sm">
                        <p class="font-medium">Groups:</p>
                        <ul class="list-disc list-inside">
                            {#each details.expand?.groups as group}
                                <li>{group.name}</li>
                            {/each}
                        </ul>
                    </div>
                {/if}

                <!-- Associated Machines -->
                {#if details.expand && details.expand?.machines.length > 0}
                    <div class="text-sm">
                        <p class="font-medium">Machines:</p>
                        <ul class="list-disc list-inside">
                            {#each details.expand?.machines as machine}
                                <li>
                                    {machine.name} ({machine.host}:{machine.port})
                                </li>
                            {/each}
                        </ul>
                    </div>
                {/if}

                <!-- Associated Users -->
                {#if details.expand && details.expand?.users.length > 0}
                    <div class="text-sm">
                        <p class="font-medium">Users:</p>
                        <ul class="list-disc list-inside">
                            {#each details.expand?.users as user}
                                <li>{user.username}</li>
                            {/each}
                        </ul>
                    </div>
                {/if}

                <!-- Created and Updated Timestamps -->
                <div class="text-xs text-gray-500 pt-2">
                    <p>
                        Created: {new Date(
                            details.created,
                        ).toLocaleDateString()}
                    </p>
                    <p>
                        Updated: {new Date(
                            details.updated,
                        ).toLocaleDateString()}
                    </p>
                </div>
            </div>
        </HoverCard.Content>
    </HoverCard.Root>
{/if}
