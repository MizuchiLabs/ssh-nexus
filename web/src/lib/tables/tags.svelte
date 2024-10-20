<script lang="ts">
    import { Badge } from "$lib/components/ui/badge";
    import { tags } from "$lib/subscriptions";

    export let id: string[];

    $: tagNames = $tags
        .filter((tag) => id.includes(tag.id))
        .map((tag) => tag.name);
</script>

{#if tagNames.length !== 0}
    {#each tagNames as name}
        <Badge variant="outline" color="primary" class="mr-0.5 mb-0.5">
            {name}
        </Badge>
    {/each}
{/if}
