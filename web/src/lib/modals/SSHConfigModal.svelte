<script lang="ts">
	import { clipboard, getModalStore } from "@skeletonlabs/skeleton";
	import { onMount, type SvelteComponent } from "svelte";
	import { generateConfig } from "$lib/utils/SSHConfig";

	export let parent: SvelteComponent;

	const modalStore = getModalStore();

	let sshConfig = "";
	const onKeys = (e: KeyboardEvent) => {
		if (e.key === "Escape") {
			modalStore.close();
		}
	};
	onMount(async () => {
		sshConfig = (await generateConfig()) || "";
	});
</script>

{#if $modalStore[0]}
	<div class="card p-4 w-modal shadow-xl space-y-4">
		<div on:keydown={onKeys} aria-hidden class="flex flex-col gap-4">
			<textarea
				class="textarea"
				rows={sshConfig.split("\n").length}
				value={sshConfig}
				on:click={(event) => event.target.select()}
				data-clipboard="sshConfig"
				readonly
			/>
		</div>
		<footer class="modal-footer {parent.regionFooter}">
			<button
				class="btn variant-filled-success"
				use:clipboard={{ input: "sshConfig" }}
				on:click={parent.onClose}
				value="sshConfig"
				>Copy
			</button>
			<button class="btn variant-filled-surface" on:click={parent.onClose}
				>Close</button
			>
		</footer>
	</div>
{/if}
