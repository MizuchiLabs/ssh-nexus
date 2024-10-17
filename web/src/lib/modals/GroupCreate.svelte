<script lang="ts">
	import { type SvelteComponent } from "svelte";
	import { getModalStore, getToastStore } from "@skeletonlabs/skeleton";
	import type { ClientResponseError, RecordModel } from "pocketbase";
	import { pb } from "$lib/client";
	import { showToast } from "$lib/utils/Toast";

	export let parent: SvelteComponent;

	const modalStore = getModalStore();
	const toastStore = getToastStore();

	let group = {} as RecordModel;

	const createGroup = async () => {
		try {
			if (!group.linux_username) {
				group.linux_username = "root";
			}
			await pb.collection("groups").create(group);
			showToast(toastStore, `Created group ${group.name}`);
		} catch (error: ClientResponseError | any) {
			showToast(
				toastStore,
				error.data?.message || "Something went wrong.",
				"error",
			);
		}
		modalStore.close();
	};

	const onKeys = (e: KeyboardEvent) => {
		if (e.key === "Enter") {
			createGroup();
		}
		if (e.key === "Escape") {
			modalStore.close();
		}
	};
</script>

{#if $modalStore[0]}
	<div class="card p-4 w-modal shadow-xl space-y-4">
		<div class="flex justify-between items-center">
			<header class="text-2xl font-bold">
				{$modalStore[0].title ?? "(title missing)"}
			</header>
		</div>
		<div on:keydown={onKeys} aria-hidden class="flex flex-col gap-4">
			<label class="label">
				<span>Name</span>
				<input
					class="input variant-form-material"
					type="text"
					name="name"
					bind:value={group.name}
					placeholder="Enter the group name"
				/>
			</label>
			<label class="label">
				<span>Description</span>
				<input
					class="input variant-form-material"
					type="text"
					name="host"
					bind:value={group.description}
					placeholder="Enter a description (optional)"
				/>
			</label>
			<label class="label">
				<span>Linux Username</span>
				<input
					class="input variant-form-material"
					type="text"
					name="linux_username"
					bind:value={group.linux_username}
					placeholder="Enter the username which will be used to login on machines"
				/>
			</label>
		</div>
		<footer class="modal-footer {parent.regionFooter}">
			<button class="btn variant-filled-success" on:click={createGroup}
				>Save</button
			>
			<button class="btn variant-filled-surface" on:click={parent.onClose}
				>Cancel</button
			>
		</footer>
	</div>
{/if}
