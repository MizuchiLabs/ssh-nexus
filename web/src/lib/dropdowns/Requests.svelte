<script lang="ts">
    import { pb, user } from "$lib/client";
    import { requests } from "$lib/subscriptions";
    import * as DropdownMenu from "$lib/components/ui/dropdown-menu/index.js";
    import { Button } from "$lib/components/ui/button/index.js";
    import type { RecordModel } from "pocketbase";
    import { Bell, Bomb, Server, Users, X, SquarePlus } from "lucide-svelte";
    import { toast } from "svelte-sonner";

    const deleteRequest = async (request: RecordModel) => {
        await pb.collection("requests").delete(request.id);
        toast.success("Deleted request");
    };

    const clearRequests = async () => {
        let requests = await pb.collection("requests").getFullList({
            filter: pb.filter("user.id = {:user_id}", {
                user_id: pb.authStore.model?.id,
            }),
        });

        for (let request of requests) {
            await pb.collection("requests").delete(request.id);
        }
        toast.success("Cleared all requests");
    };

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

{#if canRequest()}
    <DropdownMenu.Root>
        <DropdownMenu.Trigger asChild let:builder>
            <Button
                builders={[builder]}
                variant="ghost"
                size="icon"
                class="h-8 w-8 rounded-full relative"
            >
                <Bell size="1rem" />
                <span class="absolute -top-1 -right-1 text-xs text-gray-500">
                    {$requests.length}
                </span>
            </Button>
        </DropdownMenu.Trigger>
        <DropdownMenu.Content
            class="no-scrollbar max-w-md max-h-96 overflow-y-auto"
        >
            <DropdownMenu.Label
                class="flex flex-row items-center justify-between gap-2"
            >
                Requests
                <Button
                    variant="ghost"
                    href="/request/"
                    size="icon"
                    class="h-8 w-8 rounded-full"
                >
                    <SquarePlus size="1rem" />
                </Button>
            </DropdownMenu.Label>
            {#if $requests.length !== 0}
                <DropdownMenu.Separator />
                <DropdownMenu.Group>
                    {#each $requests as request}
                        <div class="flex flex-row justify-between gap-1">
                            <div class="flex flex-col gap-1">
                                {#if request.machine}
                                    <DropdownMenu.Item
                                        class="flex flex-row items-center gap-2"
                                    >
                                        <Server size="1rem" />
                                        <span
                                            >{request.expand?.machine
                                                .name}</span
                                        >
                                    </DropdownMenu.Item>
                                {/if}
                                {#if request.group}
                                    <DropdownMenu.Item
                                        class="flex flex-row items-center gap-2"
                                    >
                                        <Users size="1rem" />
                                        <span>{request.expand?.group.name}</span
                                        >
                                    </DropdownMenu.Item>
                                {/if}
                            </div>
                            <DropdownMenu.Item
                                on:click={() => deleteRequest(request)}
                            >
                                <X size="1rem" />
                            </DropdownMenu.Item>
                        </div>
                    {/each}
                </DropdownMenu.Group>
                <DropdownMenu.Separator />
                <DropdownMenu.Item
                    on:click={clearRequests}
                    class="flex flex-row items-center gap-2"
                >
                    <Bomb size="1rem" />
                    <span>Clear requests</span>
                </DropdownMenu.Item>
            {/if}
        </DropdownMenu.Content>
    </DropdownMenu.Root>
{:else}
    <Button
        href="/request/"
        variant="ghost"
        size="icon"
        class="h-8 w-8 rounded-full relative"
    >
        <Bell size="1rem" />
        <span class="absolute -top-1 -right-1 text-xs text-gray-500">
            {$requests.length}
        </span>
    </Button>
{/if}
