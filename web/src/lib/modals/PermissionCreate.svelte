<script lang="ts">
	import { pb } from "$lib/client";
	import { users, groups, machines } from "$lib/subscriptions";
	import { type SvelteComponent } from "svelte";
	import { getModalStore, getToastStore } from "@skeletonlabs/skeleton";
	import type { ClientResponseError, RecordModel } from "pocketbase";
	import { showToast } from "$lib/utils/Toast";
	import AutoChips from "$lib/utils/AutoChips.svelte";

	export let parent: SvelteComponent;

	const modalStore = getModalStore();
	const toastStore = getToastStore();

	let permission = {} as RecordModel;

	let access: Record<string, boolean> = {
		Users: false,
		Groups: false,
		Machines: false,
	};
	let type: Record<string, boolean> = {
		Create: false,
		Update: false,
		Delete: false,
	};

	const toggleAccess = async (a: string) => {
		access[a] = !access[a];
		permission = {
			...permission,
			access_users: access.Users,
			access_groups: access.Groups,
			access_machines: access.Machines,
		};
	};
	const toggleType = async (t: string) => {
		type[t] = !type[t];
		permission = {
			...permission,
			can_create: type.Create,
			can_update: type.Update,
			can_delete: type.Delete,
		};
	};
	const toggleAdmin = async () => {
		permission = {
			...permission,
			is_admin: !permission.is_admin,
		};
	};

	const createPermission = async () => {
		try {
			await pb.collection("permissions").create(permission);
			showToast(
				toastStore,
				`Created permission ${permission.name}`,
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
					bind:value={permission.name}
					placeholder="Enter a name"
					autocorrect="off"
				/>
			</label>
			<label class="label">
				<span>Description</span>
				<input
					class="input variant-form-material"
					type="text"
					name="description"
					bind:value={permission.description}
					placeholder="Enter a description"
					autocorrect="off"
				/>
			</label>
			<label class="label">
				<span>Admin</span>
				<div class="flex flex-row gap-2">
					<button
						class="chip {permission.is_admin
							? 'variant-filled-primary'
							: 'variant-filled-surface'}"
						on:click={() => {
							toggleAdmin();
						}}
						on:keypress
					>
						{#if permission.is_admin}
							<iconify-icon icon="fa6-solid:check" />
						{/if}
						<span>Set as admin</span>
					</button>
				</div>
			</label>
			<label class="label">
				<span>Access</span>
				<div class="flex flex-row gap-2">
					{#each Object.keys(access) as f}
						<button
							class="chip {access[f]
								? 'variant-filled-primary'
								: 'variant-filled-surface'}"
							on:click={() => {
								toggleAccess(f);
							}}
							on:keypress
						>
							{#if access[f]}
								<iconify-icon icon="fa6-solid:check" />
							{/if}
							<span class="capitalize">{f}</span>
						</button>
					{/each}
				</div>
			</label>
			<label class="label">
				<span>Type</span>
				<div class="flex flex-row gap-2">
					{#each Object.keys(type) as f}
						<button
							class="chip {type[f]
								? 'variant-filled-primary'
								: 'variant-filled-surface'}"
							on:click={() => {
								toggleType(f);
							}}
							on:keypress
						>
							{#if type[f]}
								<iconify-icon icon="fa6-solid:check" />
							{/if}
							<span class="capitalize">{f}</span>
						</button>
					{/each}
				</div>
			</label>
			<AutoChips
				name="users"
				data={$users}
				bind:init={permission.users}
			/>
			<AutoChips
				name="groups"
				data={$groups}
				bind:init={permission.groups}
			/>
			<AutoChips
				name="machines"
				data={$machines}
				bind:init={permission.machines}
			/>
		</div>
		<footer class="modal-footer {parent.regionFooter}">
			<button
				class="btn variant-filled-success"
				on:click={createPermission}>Save</button
			>
			<button class="btn variant-filled-surface" on:click={parent.onClose}
				>Cancel</button
			>
		</footer>
	</div>
{/if}
