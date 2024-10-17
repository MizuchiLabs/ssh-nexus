<script lang="ts">
	import { pb, user, avatar, isAdmin } from "$lib/client";
	import { downloadConfig, generateConfig } from "$lib/utils/SSHConfig";
	import { showToast } from "$lib/utils/Toast";
	import {
		FileButton,
		getModalStore,
		getToastStore,
		popup,
		type ModalSettings,
	} from "@skeletonlabs/skeleton";
	import type { ClientResponseError } from "pocketbase";
	import { onMount } from "svelte";

	let sshKey: string = $user?.settings?.ssh_key_name || "";

	const toastStore = getToastStore();
	const updateAvatar = async (e: any) => {
		let file = e.target.files[0];
		let extension = file.name.split(".").pop();
		let filename = pb.authStore.model?.id + "_avatar" + "." + extension;

		const formData = new FormData();
		formData.append("file", file, filename);
		try {
			await pb.collection("users").update($user?.id, {
				avatar: formData.get("file"),
			});
			showToast(toastStore, `Updated avatar 👌`, "success");
		} catch (error: ClientResponseError | any) {
			showToast(
				toastStore,
				error.data?.message || "Something went wrong.",
				"error",
			);
		}
	};

	const deleteAvatar = async () => {
		try {
			await pb.collection("users").update($user?.id, { avatar: null });
			showToast(toastStore, `Deleted avatar 👌`, "error");
		} catch (error: ClientResponseError | any) {
			showToast(
				toastStore,
				error.data?.message || "Something went wrong.",
				"error",
			);
		}
	};

	const updateProfile = async (e: KeyboardEvent) => {
		if (e.key !== "Enter") return;

		let settings = { ...$user?.settings, ssh_key_name: sshKey };
		await pb.collection("users").update($user?.id, {
			name: $user?.name,
			email: $user?.email,
			settings: settings,
		});
		showToast(toastStore, `Updated profile 👌`, "success");
	};

	let sshConfig = "";
	const modals = getModalStore();
	const viewConfig = () => {
		const modal: ModalSettings = {
			type: "component",
			title: "SSH Config",
			component: "SSHConfigModal",
		};
		modals.trigger(modal);
	};
	onMount(async () => {
		sshConfig = (await generateConfig()) || "";
	});
</script>

<div class="flex flex-col mx-auto p-4 gap-8">
	<h1 class="text-2xl font-bold my-2 flex justify-center items-center">
		<span class="mr-2">Profile</span>
	</h1>
	<div
		class="card self-center relative mx-auto flex flex-col w-full md:w-1/2"
	>
		<span
			class="badge-icon variant-filled-primary absolute -top-2 -right-2 z-10 w-8 h-8"
		>
			{#if pb.authStore.isAdmin}
				<iconify-icon
					icon="fa6-solid:crown"
					use:popup={{
						event: "hover",
						target: "popupSuperAdmin",
						placement: "bottom",
					}}
				/>
				<div
					class="variant-filled-secondary card p-4 w-40"
					data-popup="popupSuperAdmin"
				>
					<p>You are super admin!</p>
					<div class="arrow variant-filled-secondary" />
				</div>
			{/if}
			{#if pb.authStore.isAuthRecord}
				{#if $isAdmin}
					<iconify-icon
						icon="fa6-solid:crown"
						use:popup={{
							event: "hover",
							target: "popupAdmin",
							placement: "bottom",
						}}
					/>
					<div
						class="variant-filled-secondary card p-4 w-32"
						data-popup="popupAdmin"
					>
						<p>You are admin!</p>
						<div class="arrow variant-filled-secondary" />
					</div>
				{:else}
					<iconify-icon
						icon="fa6-solid:user"
						use:popup={{
							event: "hover",
							target: "popupUser",
							placement: "bottom",
						}}
					/>
					<div
						class="variant-filled-secondary card p-4 min-w-24"
						data-popup="popupUser"
					>
						{#if $user?.expand?.permission?.name === undefined}
							<p>You don't have any permissions</p>
						{:else}
							<p>You are {$user?.expand?.permission.name}</p>
						{/if}
						<div class="arrow variant-filled-secondary" />
					</div>
				{/if}
			{/if}
		</span>
		<header class="card-header text-center self-center">
			{#if $avatar}
				<FileButton
					name="file"
					button="bg-transparent rounded-full"
					on:change={updateAvatar}
				>
					<img
						src={$avatar}
						class="object-cover rounded-full w-40 h-40 hover:brightness-50"
						alt="avatar"
					/>
				</FileButton>
				<button
					type="button"
					class="btn-icon btn-icon-sm variant-filled-primary mt-2"
					on:click={deleteAvatar}
				>
					<iconify-icon icon="fa6-solid:skull" />
				</button>
			{:else}
				<FileButton
					name="file"
					button="bg-transparent rounded-full"
					on:change={updateAvatar}
				>
					<button
						class="btn-icon w-40 h-40 hover:brightness-50 bg-gradient-to-br from-pink-500 to-violet-500"
					>
						<iconify-icon
							icon="tabler:user"
							width="60"
							class="text-white"
						/>
					</button>
				</FileButton>
			{/if}
		</header>
		<section
			class="flex flex-col gap-2 p-4"
			on:keydown={updateProfile}
			aria-hidden
		>
			{#if !pb.authStore.isAdmin}
				<span>Name</span>
				<input
					class="input"
					type="text"
					placeholder="Name"
					bind:value={$user.name}
				/>
				<span>Email</span>
				<input
					class="input"
					type="email"
					placeholder="Email"
					bind:value={$user.email}
				/>
				<span class="flex items-center gap-1">
					SSH Key Name
					<iconify-icon
						icon="fa6-solid:circle-question"
						use:popup={{
							event: "hover",
							target: "popupKey",
							placement: "right",
						}}
					/>
				</span>
				<div
					class="card p-4 variant-filled-secondary w-64"
					data-popup="popupKey"
				>
					<p>
						Set a key name that will be used when generating your
						ssh config. Use the local name of your ssh key e.g.
						"id_rsa".
					</p>
					<div class="arrow variant-filled-secondary" />
				</div>
				<input
					class="input"
					type="text"
					placeholder="Key Name"
					bind:value={sshKey}
				/>
				<span class="flex items-center gap-1">
					Principal ID
					<iconify-icon
						icon="fa6-solid:circle-question"
						use:popup={{
							event: "hover",
							target: "popupPrincipal",
							placement: "right",
						}}
					/>
				</span>
				<div
					class="card p-4 variant-filled-secondary w-64"
					data-popup="popupPrincipal"
				>
					<p>
						This is a random generated principal id. It is fixed and
						not changeable.
					</p>
					<div class="arrow variant-filled-secondary" />
				</div>
				<input
					class="input text-gray-500"
					type="text"
					placeholder="Principal Name"
					value={$user.principal}
					tabindex="-1"
					readonly
				/>
			{/if}
		</section>
	</div>

	<div
		class="card self-center relative mx-auto flex flex-col w-full md:w-1/2"
	>
		<header
			class="card-header flex flex-row items-center justify-between space-y-0 pb-2"
		>
			<span class="text-lg font-medium">SSH Config</span>
			<iconify-icon icon="fa6-solid:key" />
		</header>
		<section class="flex flex-col text-sm font-bold px-4 gap-2 mb-2">
			<span class="text-surface-400 text-sm"
				>Your personal SSH configuration</span
			>
			<button
				class="btn btn-sm variant-filled-surface"
				on:click={viewConfig}
			>
				View</button
			>
			<button
				class="btn btn-sm variant-filled-surface"
				on:click={() => downloadConfig(sshConfig)}
			>
				Download
			</button>
		</section>
	</div>
</div>
