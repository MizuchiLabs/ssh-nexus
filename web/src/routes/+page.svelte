<script lang="ts">
	import { pb } from "$lib/client";
	import {
		users,
		machines,
		groups,
		permissions,
		tags,
		providers,
	} from "$lib/subscriptions";
	import * as Card from "$lib/components/ui/card/index.js";
	import { Button } from "$lib/components/ui/button/index.js";
	import { User, Server, Users, Lock, Tag, Cloud } from "lucide-svelte";
	import { toast } from "svelte-sonner";

	$: stats = [
		{
			title: "Users",
			icon: User,
			color: "text-blue-500",
			count: $users.length,
		},
		{
			title: "Machines",
			icon: Server,
			color: "text-purple-500",
			count: $machines.length,
		},
		{
			title: "Groups",
			icon: Users,
			color: "text-green-500",
			count: $groups.length,
		},
		{
			title: "Permissions",
			icon: Lock,
			color: "text-red-500",
			count: $permissions.length,
		},
		{
			title: "Tags",
			icon: Tag,
			color: "text-yellow-500",
			count: $tags.length,
		},
		{
			title: "Providers",
			icon: Cloud,
			color: "text-indigo-500",
			count: $providers.length,
		},
	];

	const installAgents = async () => {
		await pb.send("/api/sync/agents", { method: "POST" });
		toast.warning("Installing agents...");
	};
	const syncProviders = async () => {
		await pb.send("/api/sync/providers", { method: "POST" });
		toast.warning("Syncing providers...");
	};
</script>

<div class="px-4 py-6">
	<h1 class="text-3xl font-bold mb-6">Overview</h1>

	<div class="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3">
		{#each stats as stat}
			<Card.Root>
				<Card.Header
					class="flex flex-row items-center justify-between space-y-0 pb-2"
				>
					<Card.Title class="text-sm font-medium">
						{stat.title}
					</Card.Title>
					<svelte:component
						this={stat.icon}
						class={`${stat.color}`}
						size="1.1rem"
					/>
				</Card.Header>
				<Card.Content class="flex flex-row items-start justify-between">
					<div class="text-2xl font-bold">{stat.count}</div>
					{#if stat.title === "Machines"}
						<Button
							variant="secondary"
							class="h-8 rounded-full hover:bg-transparent"
							on:click={installAgents}
						>
							Sync Agents
						</Button>
					{/if}
					{#if stat.title === "Providers"}
						<Button
							variant="secondary"
							class="h-8 rounded-full hover:bg-transparent"
							on:click={syncProviders}
						>
							Sync Providers
						</Button>
					{/if}
				</Card.Content>
			</Card.Root>
		{/each}
	</div>
</div>
