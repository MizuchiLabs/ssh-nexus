<script lang="ts">
	import { getModalStore, getToastStore } from "@skeletonlabs/skeleton";
	import type { ClientResponseError, RecordModel } from "pocketbase";
	import { pb } from "$lib/client";
	import { showToast } from "$lib/utils/Toast";
	import type { SvelteComponent } from "svelte";

	export let parent: SvelteComponent;

	const modalStore = getModalStore();
	const toastStore = getToastStore();

	let provider = $modalStore[0].meta.provider as RecordModel;

	const updateProvider = async () => {
		try {
			await pb.collection("providers").update(provider.id, provider);
			showToast(
				toastStore,
				`Updated provider ${provider.name}`,
				"success",
			);
		} catch (error: ClientResponseError | any) {
			showToast(
				toastStore,
				error.data?.message || "Something went wrong.",
				"error",
			);
		}
		modalStore.close();
	};

	const deleteProvider = async () => {
		try {
			await pb.collection("providers").delete(provider.id);
			showToast(toastStore, `Deleted provider ${provider.name}`, "error");
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
			updateProvider();
		}
		if (e.key === "Escape") {
			modalStore.close();
		}
	};
</script>

{#if $modalStore[0]}
	<div class="card p-4 w-modal shadow-xl space-y-4">
		<div class="flex justify-between">
			<header class="text-2xl font-bold flex items-center">
				{#if provider.type === "aws"}
					<iconify-icon
						icon="simple-icons:amazonaws"
						width="24"
						class="mr-2"
					/>
				{:else if provider.type === "azure"}
					<iconify-icon
						icon="simple-icons:azuredevops"
						width="24"
						class="mr-2"
					/>
				{:else if provider.type === "linode"}
					<iconify-icon
						icon="fa6-brands:linode"
						width="24"
						class="mr-2"
					/>
				{:else if provider.type === "google"}
					<iconify-icon
						icon="simple-icons:googlecloud"
						width="24"
						class="mr-2"
					/>
				{:else if provider.type === "hetzner"}
					<iconify-icon
						icon="simple-icons:hetzner"
						width="24"
						class="mr-2"
					/>
				{:else if provider.type === "vultr"}
					<iconify-icon
						icon="simple-icons:vultr"
						width="24"
						class="mr-2"
					/>
				{:else if provider.type === "proxmox"}
					<iconify-icon
						icon="simple-icons:proxmox"
						width="24"
						class="mr-2"
					/>
				{/if}
				<p>{$modalStore[0].title ?? "(title missing)"}</p>
			</header>
			<button
				class="btn-icon btn-icon-sm variant-filled-error"
				on:click={deleteProvider}
				tabindex="-1"
			>
				<iconify-icon icon="mdi:trash" width="20" />
			</button>
		</div>
		<div on:keydown={onKeys} aria-hidden class="flex flex-col gap-4">
			<label class="label">
				<span>Name</span>
				<input
					class="input variant-form-material"
					type="text"
					name="name"
					bind:value={provider.name}
					placeholder="Enter a name"
					autocorrect="off"
				/>
			</label>
			<label class="label">
				<span>URL</span>
				<input
					class="input variant-form-material"
					type="text"
					name="url"
					bind:value={provider.url}
					placeholder="Enter the api url"
					autocorrect="off"
				/>
			</label>
			<label class="label">
				<span>Username</span>
				<input
					class="input variant-form-material"
					type="text"
					name="username"
					bind:value={provider.username}
					placeholder="Enter the username"
					autocorrect="off"
				/>
			</label>
			<label class="label">
				<span>Password</span>
				<input
					class="input variant-form-material"
					type="password"
					name="password"
					bind:value={provider.password}
					placeholder="Enter the password"
					autocorrect="off"
				/>
			</label>
			<label class="label">
				<span>Token</span>
				<input
					class="input variant-form-material"
					type="password"
					name="token"
					bind:value={provider.token}
					placeholder="Enter the token"
					autocorrect="off"
				/>
			</label>
		</div>
		<footer class="modal-footer {parent.regionFooter}">
			<button class="btn variant-filled-success" on:click={updateProvider}
				>Save</button
			>
			<button class="btn variant-filled-surface" on:click={parent.onClose}
				>Cancel</button
			>
		</footer>
	</div>
{/if}
