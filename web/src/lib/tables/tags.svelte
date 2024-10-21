<script lang="ts">
    import { Badge } from "$lib/components/ui/badge";
    import { tags } from "$lib/subscriptions";
    import type { RecordModel } from "pocketbase";
    import { writable, type Writable } from "svelte/store";

    export let id: string[];
    export let selected: Writable<RecordModel[]> = writable([]);

    $: tag = $tags.filter((tag) => id.includes(tag.id));

    const toggleTag = (tag: RecordModel) => {
        if ($selected.includes(tag)) {
            $selected = $selected.filter((t) => t.id !== tag.id);
        } else {
            $selected = [...$selected, tag];
        }
    };
</script>

{#if tag.length !== 0}
    {#each tag as t}
        <button on:click={() => toggleTag(t)}>
            <Badge variant="outline" class="mr-0.5 mb-0.5">
                {t.name}
            </Badge>
        </button>
    {/each}
{/if}
