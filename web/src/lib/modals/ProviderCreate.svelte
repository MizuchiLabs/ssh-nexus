<script lang="ts">
	import { getModalStore, getToastStore } from "@skeletonlabs/skeleton";
	import type { ClientResponseError, RecordModel } from "pocketbase";
	import { pb } from "$lib/client";
	import { showToast } from "$lib/utils/Toast";
	import type { SvelteComponent } from "svelte";

	export let parent: SvelteComponent;

	const modalStore = getModalStore();
	const toastStore = getToastStore();

	let provider = {} as RecordModel;

	const create = async () => {
		try {
			await pb.collection("providers").create(provider);
			showToast(
				toastStore,
				`Created provider ${provider.name}`,
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

	const onKeys = (e: KeyboardEvent) => {
		if (e.key === "Enter") {
			create();
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
					bind:value={provider.name}
					placeholder="Enter a name"
					autocorrect="off"
				/>
			</label>
			<label class="label flex flex-col">
				<span>Type</span>
				<div
					class="btn-group variant-filled-surface [&>*+*]:border-gray-500 self-center"
				>
					<button
						value="aws"
						on:click={() => (provider.type = "aws")}
					>
						<iconify-icon
							icon="simple-icons:amazonaws"
							width="24"
						/>
					</button>
					<button
						value="azure"
						on:click={() => (provider.type = "azure")}
					>
						<iconify-icon
							icon="simple-icons:azuredevops"
							width="24"
						/>
					</button>
					<button
						value="linode"
						on:click={() => (provider.type = "linode")}
					>
						<iconify-icon icon="fa6-brands:linode" width="24" />
					</button>
					<button
						value="google"
						on:click={() => (provider.type = "google")}
					>
						<iconify-icon
							icon="simple-icons:googlecloud"
							width="24"
						/>
					</button>
					<button
						value="hetzner"
						on:click={() => (provider.type = "hetzner")}
					>
						<iconify-icon icon="simple-icons:hetzner" width="24" />
					</button>
					<button
						value="vultr"
						on:click={() => (provider.type = "vultr")}
					>
						<iconify-icon icon="simple-icons:vultr" width="24" />
					</button>
					<button
						value="proxmox"
						on:click={() => (provider.type = "proxmox")}
					>
						<iconify-icon icon="simple-icons:proxmox" width="24" />
					</button>
				</div>
			</label>
			{#if provider.type === "proxmox"}
				<label class="label">
					<span>URL</span>
					<input
						class="input variant-form-material"
						type="text"
						name="url"
						bind:value={provider.url}
						placeholder="Enter the API URL"
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
						placeholder="Enter the Username/Token ID"
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
			{/if}
			{#if provider.type === "aws"}
				<label class="label">
					<span>ID</span>
					<input
						class="input variant-form-material"
						type="text"
						name="id"
						bind:value={provider.username}
						placeholder="Enter an ID"
						autocorrect="off"
					/>
				</label>
				<label class="label">
					<span>Secret</span>
					<input
						class="input variant-form-material"
						type="password"
						name="secret"
						bind:value={provider.password}
						placeholder="Enter the secret"
						autocorrect="off"
					/>
				</label>
			{/if}
			{#if provider.type != undefined}
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
			{/if}
		</div>
		<footer class="modal-footer {parent.regionFooter}">
			<button class="btn variant-filled-success" on:click={create}
				>Save</button
			>
			<button class="btn variant-filled-surface" on:click={parent.onClose}
				>Cancel</button
			>
		</footer>
	</div>
{/if}
