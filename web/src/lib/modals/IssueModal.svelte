<script lang="ts">
	import { machines } from "$lib/subscriptions";
	import { getModalStore } from "@skeletonlabs/skeleton";
	import type { SvelteComponent } from "svelte";

	export let parent: SvelteComponent;

	const modalStore = getModalStore();
</script>

{#if $modalStore[0]}
	<div class="card p-4 w-modal shadow-xl space-y-4">
		<div class="flex justify-between items-center">
			<header class="text-2xl font-bold">
				{$modalStore[0].title ?? "(title missing)"}
			</header>
		</div>
		<ul class="space-y-2">
			{#each $machines.filter((m) => m.error !== "") as machine}
				<li class="p-2 bg-gray-700 rounded-lg shadow-inner">
					<span class="font-bold">{machine.name}:</span>
					<span class="text-red-500 text-sm">{machine.error}</span>
				</li>
			{/each}
		</ul>
		<footer class="modal-footer {parent.regionFooter}">
			<button class="btn variant-filled-surface" on:click={parent.onClose}
				>Close</button
			>
		</footer>
	</div>
{/if}
