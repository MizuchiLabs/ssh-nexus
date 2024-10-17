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

	let user = {} as RecordModel;

	const createUser = async () => {
		if (!passwordOk) {
			showToast(toastStore, "Passwords must match", "error");
			return;
		}
		try {
			await pb.collection("users").create(user);
			showToast(toastStore, `Created user ${user.name}`, "success");
		} catch (error: ClientResponseError | any) {
			showToast(
				toastStore,
				error.data?.message || "Something went wrong.",
				"error",
			);
		}
		modalStore.close();
	};

	let passwordOk = false;
	const passwordCheck = () => {
		if (!user.password) {
			return;
		}
		if (user.passwordConfirm !== user.password) {
			passwordOk = false;
			return;
		}
		if (user.password.length < 8) {
			passwordOk = false;
			return;
		}
		if (!/[a-z]/.test(user.password)) {
			passwordOk = false;
			return;
		}
		if (!/\d/.test(user.password)) {
			passwordOk = false;
			return;
		}
		passwordOk = true;
	};

	const onKeys = (e: KeyboardEvent) => {
		if (e.key === "Enter") {
			createUser();
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
				<span
					>Password <span class="ml-1 text-sm text-gray-500"
						>(min. 8 characters and at least 1 number)</span
					></span
				>
				<input
					class="input variant-form-material"
					class:input-success={passwordOk}
					class:input-error={!passwordOk &&
						user.password !== undefined}
					type="password"
					name="password"
					bind:value={user.password}
					placeholder="Enter the password"
					autocorrect="off"
					on:input={passwordCheck}
				/>
			</label>
			<label class="label">
				<span>Confirm Password</span>
				<input
					class="input variant-form-material"
					class:input-success={passwordOk}
					class:input-error={!passwordOk &&
						user.password !== undefined}
					type="password"
					name="password"
					bind:value={user.passwordConfirm}
					placeholder="Enter the password again"
					autocorrect="off"
					on:input={passwordCheck}
				/>
			</label>
			<label class="label">
				<span>Select permission</span>
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
			<button class="btn variant-filled-success" on:click={createUser}
				>Save</button
			>
			<button class="btn variant-filled-surface" on:click={parent.onClose}
				>Cancel</button
			>
		</footer>
	</div>
{/if}
