<script lang="ts">
    import { CircleCheck, CircleX } from "lucide-svelte";

    export let isTrue: boolean;
</script>

{#if isTrue}
    <CircleCheck size="1rem" class="text-green-500" />
{:else}
    <CircleX size="1rem" class="text-red-500" />
{/if}
