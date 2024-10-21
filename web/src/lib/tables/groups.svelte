<script lang="ts">
    import { Badge } from "$lib/components/ui/badge";
    import { groups } from "$lib/subscriptions";
    import type { RecordModel } from "pocketbase";
    import { writable, type Writable } from "svelte/store";

    export let id: string[];
    export let selected: Writable<RecordModel[]> = writable([]);

    $: group = $groups.filter((group) => id.includes(group.id));

    const toggleGroup = (group: RecordModel) => {
        if ($selected.includes(group)) {
            $selected = $selected.filter((g) => g.id !== group.id);
        } else {
            $selected = [...$selected, group];
        }
    };
</script>

{#if group.length !== 0}
    {#each group as g}
        <button on:click={() => toggleGroup(g)}>
            <Badge variant="default" class="mr-0.5 mb-0.5">
                {g.name}
            </Badge>
        </button>
    {/each}
{/if}
