<script lang="ts">
	import { getModalStore, getToastStore } from "@skeletonlabs/skeleton";
	import { type SvelteComponent } from "svelte";
	import { pb } from "$lib/client";
	import { showToast } from "$lib/utils/Toast";

	export let parent: SvelteComponent;

	const modalStore = getModalStore();
	const toastStore = getToastStore();

	let key = "";
	const update = async () => {
		try {
			await pb.send("/api/ssh/user/set", {
				method: "POST",
				body: { key: key },
			});
			showToast(toastStore, "Updated user CA", "success");
		} catch (error: any) {
			showToast(
				toastStore,
				"Failed to update key, make sure you are using the correct openssh format!",
				"error",
			);
		}
		modalStore.close();
	};

	const onKeys = (e: KeyboardEvent) => {
		if (e.key === "Escape") {
			modalStore.close();
		}
	};
</script>

{#if $modalStore[0]}
	<div class="card p-4 w-modal shadow-xl space-y-4">
		<div on:keydown={onKeys} aria-hidden class="flex flex-col gap-4">
			<textarea
				class="textarea"
				rows="10"
				bind:value={key}
				placeholder="-----BEGIN OPENSSH PRIVATE KEY-----
<private key>
-----END OPENSSH PRIVATE KEY-----"
			/>
		</div>
		<footer class="modal-footer {parent.regionFooter}">
			<button class="btn variant-filled-error" on:click={update}>
				Overwrite Key
			</button>
			<button class="btn variant-filled-surface" on:click={parent.onClose}
				>Close</button
			>
		</footer>
	</div>
{/if}
