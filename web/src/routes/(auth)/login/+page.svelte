<script lang="ts">
	import { goto } from "$app/navigation";
	import { loggedIn, pb } from "$lib/client";
	import { showToast } from "$lib/utils/Toast";
	import { getToastStore } from "@skeletonlabs/skeleton";
	import type { AuthProviderInfo } from "pocketbase";
	import { onMount } from "svelte";

	const toastStore = getToastStore();
	let username: string, password: string;
	const loginPassword = async () => {
		const admin = await pb.admins
			.authWithPassword(username, password)
			.catch(() => false);
		if (!admin) {
			const user = await pb
				.collection("users")
				.authWithPassword(username, password)
				.catch(() => false);
			if (!user) {
				showToast(toastStore, "Invalid username or password", "error");
				return;
			}
		}
		showToast(toastStore, "Login successful");
		loggedIn.set(true);
		goto("/");
	};
	const resetPassword = async () => {
		if (!username) return;
		try {
			await pb.collection("users").requestPasswordReset(username);
			showToast(toastStore, "Password reset email sent!");
		} catch (err) {
			showToast(toastStore, "Not a valid email address!", "error");
		}
	};

	let oauthProviders: AuthProviderInfo[] = [];
	const loginOauth = async (provider: AuthProviderInfo) => {
		try {
			await pb
				.collection("users")
				.authWithOAuth2({ provider: provider.name });
			showToast(toastStore, "Login successful");
			goto("/");
		} catch (err) {
			showToast(toastStore, "Login failed", "error");
		}
	};

	const onKeys = (e: KeyboardEvent) => {
		if (e.key === "Enter") {
			loginPassword();
		}
	};
	onMount(async () => {
		try {
			const result = await pb
				.collection("users")
				.listAuthMethods({ requestKey: null });
			oauthProviders = result.authProviders;
		} catch (err) {
			console.error("Failed to load oauth providers", err);
		}
	});
</script>

<div class="container flex flex-col gap-4 mx-auto">
	<div
		class="card lg:w-1/2 lg:self-center lg:mx-auto flex flex-col p-4 mt-8 gap-4"
		on:keydown={onKeys}
		aria-hidden
	>
		<h2 class="flex flex-col">
			<span class="text-2xl font-bold">Login</span>
			<p class="text-sm text-surface-300">
				Enter either your username or email to login
			</p>
		</h2>
		<label class="label flex flex-col gap-2">
			<span class="font-bold">Username</span>
			<input
				class="input variant-form-material"
				type="text"
				name="username"
				bind:value={username}
			/>
		</label>
		<label class="label flex flex-col gap-2">
			<span class="font-bold flex flex-row justify-between">
				Password
				<button
					class="text-sm text-surface-300 ml-2"
					on:click={resetPassword}
					tabindex="-1"
				>
					Forgot your password?
				</button>
			</span>
			<input
				class="input variant-form-material"
				type="password"
				name="password"
				bind:value={password}
			/>
		</label>
		<button class="btn variant-filled-success" on:click={loginPassword}
			>Login</button
		>
		{#if oauthProviders.length > 0}
			<div class="flex flex-col gap-2">
				{#each oauthProviders as provider}
					<button
						class="btn variant-filled-surface"
						on:click={() => loginOauth(provider)}
					>
						{#if provider.displayName}
							{provider.displayName}
						{:else}
							{provider.name}
						{/if}
					</button>
				{/each}
			</div>
		{/if}
	</div>
</div>
