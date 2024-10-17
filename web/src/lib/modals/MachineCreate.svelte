<script lang="ts">
	import { pb } from "$lib/client";
	import { users, groups, tags } from "$lib/subscriptions";
	import { type SvelteComponent } from "svelte";
	import { getModalStore, getToastStore } from "@skeletonlabs/skeleton";
	import type { ClientResponseError, RecordModel } from "pocketbase";
	import { showToast } from "$lib/utils/Toast";
	import AutoChips from "$lib/utils/AutoChips.svelte";

	export let parent: SvelteComponent;

	const modalStore = getModalStore();
	const toastStore = getToastStore();

	let machine = {} as RecordModel;

	// Function to generate a random string of specified length
	function generateRandomString(length = 8) {
		const characters =
			"ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz";
		let result = "";
		for (let i = 0; i < length; i++) {
			result += characters.charAt(
				Math.floor(Math.random() * characters.length),
			);
		}
		return result;
	}

	const createMachine = async () => {
		try {
			if (machine.name == null) {
				machine.name = "machine-" + generateRandomString(8);
			}
			if (machine.port == null) {
				machine.port = 22;
			}
			await pb.collection("machines").create(machine);
			showToast(toastStore, `Added machine ${machine.name}`, "success");
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
			createMachine();
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
					bind:value={machine.name}
					placeholder="Enter name..."
				/>
			</label>
			<label class="label">
				<span>Host</span>
				<input
					class="input variant-form-material"
					type="text"
					name="host"
					bind:value={machine.host}
					placeholder="Enter host..."
				/>
			</label>
			<label class="label">
				<span>Port</span>
				<input
					class="input variant-form-material"
					type="number"
					name="port"
					bind:value={machine.port}
					placeholder="Enter port..."
				/>
			</label>
		</div>
		<AutoChips name="users" data={$users} bind:init={machine.users} />
		<AutoChips name="groups" data={$groups} bind:init={machine.groups} />
		<AutoChips name="tags" data={$tags} bind:init={machine.tags} />
		<footer class="modal-footer {parent.regionFooter}">
			<button class="btn variant-filled-success" on:click={createMachine}
				>Save</button
			>
			<button class="btn variant-filled-surface" on:click={parent.onClose}
				>Cancel</button
			>
		</footer>
	</div>
{/if}
