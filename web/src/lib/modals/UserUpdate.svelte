<script lang="ts">
	import { pb } from "$lib/client";
	import { groups, permissions } from "$lib/subscriptions";
	import { type SvelteComponent } from "svelte";
	import { getModalStore, getToastStore } from "@skeletonlabs/skeleton";
	import type { ClientResponseError, RecordModel } from "pocketbase";
	import { showToast } from "$lib/utils/Toast";
	import AutoChips from "$lib/utils/AutoChips.svelte";

	export let parent: SvelteComponent;

	const modalStore = getModalStore();
	const toastStore = getToastStore();

	let user = $modalStore[0].meta.user as RecordModel;

	const updateUser = async () => {
		try {
			await pb.collection("users").update(user.id, user);
			showToast(toastStore, `Updated user ${user.name}`, "success");
		} catch (error: ClientResponseError | any) {
			showToast(
				toastStore,
				error.data?.message || "Something went wrong.",
				"error",
			);
		}
		modalStore.close();
	};

	const deleteUser = async () => {
		try {
			await pb.collection("users").delete(user.id);
			showToast(toastStore, `Deleted user ${user.name}`, "error");
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
			updateUser();
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
			<button
				class="btn-icon btn-icon-sm variant-filled-error"
				on:click={deleteUser}
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
					bind:value={user.name}
					placeholder="Enter a name"
					autocorrect="off"
				/>
			</label>
			<label class="label">
				<span>Username</span>
				<input
					class="input variant-form-material"
					type="text"
					name="username"
					bind:value={user.username}
					placeholder="Enter the username"
					autocorrect="off"
				/>
			</label>
			<label class="label">
				<span>Email</span>
				<input
					class="input variant-form-material"
					type="email"
					name="email"
					bind:value={user.email}
					placeholder="Enter the email"
					autocorrect="off"
				/>
			</label>
			<label class="label">
				<span>Principal</span>
				<input
					class="input variant-form-material"
					type="text"
					name="principal"
					bind:value={user.principal}
					readonly
					disabled
				/>
			</label>
			<label class="label">
				<span>Permission</span>
				<select class="select" bind:value={user.permission}>
					<option value=""></option>
					{#each $permissions as permission}
						<option value={permission.id}>{permission.name}</option>
					{/each}
				</select>
			</label>
			<AutoChips name="groups" data={$groups} bind:init={user.groups} />
		</div>
		<footer class="modal-footer {parent.regionFooter}">
			<button class="btn variant-filled-success" on:click={updateUser}
				>Save</button
			>
			<button class="btn variant-filled-surface" on:click={parent.onClose}
				>Cancel</button
			>
		</footer>
	</div>
{/if}
