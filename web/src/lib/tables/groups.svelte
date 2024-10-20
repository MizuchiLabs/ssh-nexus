<script lang="ts">
    import { Badge } from "$lib/components/ui/badge";
    import { groups } from "$lib/subscriptions";

    export let id: string[];

    $: groupNames = $groups
        .filter((group) => id.includes(group.id))
        .map((group) => group.name);
</script>

{#if groupNames.length !== 0}
    {#each groupNames as name}
        <Badge variant="default" color="primary" class="mr-0.5 mb-0.5">
            {name}
        </Badge>
    {/each}
{/if}
