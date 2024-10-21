<script lang="ts">
	import type { AuthProviderInfo } from "pocketbase";
	import * as Tabs from "$lib/components/ui/tabs/index.js";
	import * as Card from "$lib/components/ui/card/index.js";
	import { Button } from "$lib/components/ui/button/index.js";
	import { Input } from "$lib/components/ui/input/index.js";
	import { Label } from "$lib/components/ui/label/index.js";
	import { goto } from "$app/navigation";
	import { loggedIn, pb } from "$lib/client";
	import { toast } from "svelte-sonner";
	import { Eye, EyeOff } from "lucide-svelte";
	import { onMount } from "svelte";

	let username: string, password: string;
	let showPassword = false;

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
				toast.error("Invalid username or password");
				return;
			}
		}
		toast.success("Login successful");
		loggedIn.set(true);
		goto("/");
	};
	const resetPassword = async () => {
		if (!username) return;
		try {
			await pb.collection("users").requestPasswordReset(username);
			toast.success("Password reset email sent!");
		} catch (err) {
			toast.error("Failed to send password reset email");
		}
	};

	let oauthProviders: AuthProviderInfo[] = [];
	const loginOauth = async (provider: AuthProviderInfo) => {
		try {
			await pb.collection("users").authWithOAuth2({ provider: provider.name });
			toast.success("Login successful");
			goto("/");
		} catch (err) {
			toast.error("Login failed");
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
			toast.error("Failed to load oauth providers");
		}
	});
</script>

{#if !$loggedIn}
	<Tabs.Root value="username" class="w-[400px]">
		<Tabs.List class="grid w-full grid-cols-2">
			<Tabs.Trigger value="username">Username</Tabs.Trigger>
			<Tabs.Trigger value="oauth">OAuth</Tabs.Trigger>
		</Tabs.List>
		<Tabs.Content value="username">
			<Card.Root class="w-[400px]">
				<Card.Header>
					<Card.Title>Login</Card.Title>
					<Card.Description>Login to your account</Card.Description>
				</Card.Header>
				<Card.Content>
					<div
						class="grid w-full items-center gap-4"
						on:keydown={onKeys}
						aria-hidden
					>
						<div class="flex flex-col space-y-1.5">
							<Label for="username">Username</Label>
							<Input id="username" bind:value={username} />
						</div>
						<div class="flex flex-col space-y-1.5">
							<Label for="password">
								<span class="flex flex-row justify-between">
									Password
									<button on:click={resetPassword} aria-hidden>
										Forgot your password?
									</button>
								</span>
							</Label>
							<div class="flex flex-row items-center justify-end gap-1">
								{#if showPassword}
									<Input id="password" type="text" bind:value={password} />
								{:else}
									<Input id="password" type="password" bind:value={password} />
								{/if}
								<Button
									variant="ghost"
									size="icon"
									class="absolute hover:bg-transparent hover:text-red-400"
									on:click={() => (showPassword = !showPassword)}
								>
									{#if showPassword}
										<Eye size="1rem" />
									{:else}
										<EyeOff size="1rem" />
									{/if}
								</Button>
							</div>
						</div>
					</div>
					<div class="mt-4 flex flex-col">
						<Button type="submit" on:click={loginPassword}>Login</Button>
					</div>
				</Card.Content>
			</Card.Root>
		</Tabs.Content>
		<Tabs.Content value="oauth">
			<Card.Root class="w-[400px]">
				<Card.Header>
					<Card.Title>Login</Card.Title>
					<Card.Description>Login using an OAuth provider</Card.Description>
				</Card.Header>
				<Card.Content class="flex flex-col items-center gap-4">
					{#if oauthProviders.length === 0}
						<p class="text-sm">No OAuth providers configured</p>
					{:else}
						{#each oauthProviders as provider}
							<Button
								variant="default"
								class="w-full"
								on:click={() => loginOauth(provider)}
							>
								{provider.name}
							</Button>
						{/each}
					{/if}
				</Card.Content>
			</Card.Root>
		</Tabs.Content>
	</Tabs.Root>
{/if}
