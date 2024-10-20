<script lang="ts">
    import { pb, user } from "$lib/client";
    import { Button } from "$lib/components/ui/button";
    import { Check, X } from "lucide-svelte";
    import type { ClientResponseError } from "pocketbase";
    import { toast } from "svelte-sonner";

    export let id: string;

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

    const approveAction = async () => {
        try {
            const request = await pb.collection("requests").getOne(id, {
                expand: "user, machine, group",
            });
            if (request.expand?.machine) {
                await pb
                    .collection("machines")
                    .update(request.expand.machine.id, {
                        "users+": request.expand.user.id,
                    });
            }
            if (request.expand?.group) {
                await pb.collection("users").update(request.expand.user.id, {
                    "groups+": request.expand.group.id,
                });
            }
            await pb.collection("requests").delete(id);
            toast.success(`Approved`, {
                description: `Approved request from ${request.expand?.user.username}.`,
                duration: 3000,
            });
        } catch (error: ClientResponseError | any) {
            toast.error(error.data?.message || "Something went wrong.");
        }
    };
    const denyAction = async () => {
        try {
            const request = await pb.collection("requests").getOne(id, {
                expand: "user, machine, group",
            });
            await pb.collection("requests").delete(id);
            if (canRequest()) {
                toast.success("Deleted request");
            } else {
                toast.success(`Denied`, {
                    description: `Denied request from ${request.expand?.user.username}.`,
                    duration: 3000,
                });
            }
        } catch (error: ClientResponseError | any) {
            toast.error(error.data?.message || "Something went wrong.");
        }
    };
</script>

<div class="flex items-center gap-1 dark:text-black">
    {#if !canRequest()}
        <Button
            variant="ghost"
            class="h-8 w-8 rounded-full bg-green-400"
            size="icon"
            on:click={approveAction}
        >
            <Check size="1rem" />
        </Button>
    {/if}
    <Button
        variant="ghost"
        class="h-8 w-8 rounded-full bg-red-400"
        size="icon"
        on:click={denyAction}
    >
        <X size="1rem" />
    </Button>
</div>
