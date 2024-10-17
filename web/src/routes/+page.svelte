<script lang="ts">
	import { loggedIn, pb, user } from "$lib/client";
	import {
		users,
		groups,
		providers,
		auditlog,
		machines,
	} from "$lib/subscriptions";
	import { showToast } from "$lib/utils/Toast";
	import {
		getModalStore,
		getToastStore,
		popup,
		type ModalSettings,
	} from "@skeletonlabs/skeleton";

	const modals = getModalStore();
	const toastStore = getToastStore();
	const installAgents = async () => {
		await pb.send("/api/sync/agents", { method: "POST" });
		showToast(toastStore, "Reinstalling agents...", "warning");
	};
	const syncProviders = async () => {
		await pb.send("/api/sync/providers", { method: "POST" });
		showToast(toastStore, "Syncing providers...", "warning");
	};

	const createMachine = () => {
		const modal: ModalSettings = {
			type: "component",
			title: "New Machine",
			component: "MachineCreate",
		};
		modals.trigger(modal);
	};
	const createGroup = () => {
		const modal: ModalSettings = {
			type: "component",
			title: "New Group",
			component: "GroupCreate",
		};
		modals.trigger(modal);
	};
	const createUser = () => {
		const modal: ModalSettings = {
			type: "component",
			title: "New User",
			component: "UserCreate",
		};
		modals.trigger(modal);
	};
	const viewIssues = () => {
		const modal: ModalSettings = {
			type: "component",
			title: "Issues",
			component: "IssueModal",
		};
		modals.trigger(modal);
	};
</script>

<div class="container mx-auto p-4">
	{#if $loggedIn}
		<!-- Widgets -->
		<div class="grid gap-4 grid-cols-1 md:grid-cols-2 lg:grid-cols-3">
			<div class="card sm:col-span-1 col-span-2">
				<header
					class="card-header flex flex-row items-center justify-between space-y-0 pb-2"
				>
					<span class="text-lg font-medium">Machines</span>
					<iconify-icon icon="fa6-solid:server" />
				</header>
				<section
					class="flex flex-row items-center justify-between px-4 py-1 mb-2"
				>
					<span class="text-2xl font-bold">{$machines.length}</span>
					<button
						class="btn btn-sm text-sm bg-surface-300 dark:bg-surface-500"
						on:click={createMachine}
					>
						<iconify-icon icon="fa6-solid:plus" class="mr-1" />
						Add Machine
					</button>
				</section>
			</div>
			<div class="card sm:col-span-1 col-span-2">
				<header
					class="card-header flex flex-row items-center justify-between space-y-0 pb-2"
				>
					<span class="text-lg font-medium">Agents Online</span>
					<iconify-icon icon="fa6-solid:robot" />
				</header>
				<section
					class="flex flex-row items-center justify-between px-4 py-1 mb-2"
				>
					<span class="text-2xl font-bold"
						>{$machines.filter((m) => m.agent === true)
							.length}</span
					>
					<button
						class="btn btn-sm text-sm bg-surface-300 dark:bg-surface-500"
						on:click={installAgents}
					>
						<iconify-icon icon="fa6-solid:rotate" class="mr-1" />
						Install agents
					</button>
				</section>
			</div>
			<div class="card sm:col-span-1 col-span-2">
				<header
					class="card-header flex flex-row items-center justify-between space-y-0 pb-2"
				>
					<span class="text-lg font-medium">Users</span>
					<iconify-icon icon="fa6-solid:user" />
				</header>
				<section
					class="flex flex-row items-center justify-between px-4 py-1 mb-2"
				>
					<span class="text-2xl font-bold">{$users.length}</span>
					<button
						class="btn btn-sm text-sm bg-surface-300 dark:bg-surface-500"
						on:click={createUser}
					>
						<iconify-icon icon="fa6-solid:plus" class="mr-1" />
						Add User
					</button>
				</section>
			</div>
			<div class="card sm:col-span-1 col-span-2">
				<header
					class="card-header flex flex-row items-center justify-between space-y-0 pb-2"
				>
					<span class="text-lg font-medium">Groups</span>
					<iconify-icon icon="fa6-solid:users" />
				</header>
				<section
					class="flex flex-row items-center justify-between px-4 py-1 mb-2"
				>
					<span class="text-2xl font-bold">{$groups.length}</span>
					<button
						class="btn btn-sm text-sm bg-surface-300 dark:bg-surface-500"
						on:click={createGroup}
					>
						<iconify-icon icon="fa6-solid:plus" class="mr-1" />
						Add Group
					</button>
				</section>
			</div>
			<div class="card sm:col-span-1 col-span-2">
				<header
					class="card-header flex flex-row items-center justify-between space-y-0 pb-2"
				>
					<span class="text-lg font-medium">Unhealthy machines</span>
					<iconify-icon icon="fa6-solid:skull" />
				</header>
				<section
					class="flex flex-row items-center justify-between px-4 py-1 mb-2"
				>
					<span class="text-2xl font-bold"
						>{$machines.filter((m) => m.error !== "").length}</span
					>
					<button
						class="btn btn-sm text-sm bg-surface-300 dark:bg-surface-500"
						on:click={viewIssues}
					>
						<iconify-icon icon="fa6-solid:eye" class="mr-1" />
						View issues
					</button>
				</section>
			</div>
			<div class="card sm:col-span-1 col-span-2">
				<header
					class="card-header flex flex-row items-center justify-between space-y-0 pb-2"
				>
					<span class="text-lg font-medium">Memberships</span>
					<iconify-icon icon="fa6-solid:users" />
				</header>
				<section
					class="flex flex-row flex-wrap text-sm font-bold px-4 py-1 gap-2 mb-2"
				>
					{#if pb.authStore.isAdmin}
						<span class="text-surface-400 text-sm text-center">
							You are admin 🐱
						</span>
					{:else if pb.authStore.isAuthRecord}
						{#if $user?.expand?.groups}
							{#each $user?.expand?.groups as group}
								<span
									class="chip font-bold text-black bg-primary-400"
								>
									{group.name}
								</span>
							{/each}
						{:else}
							<span
								class="flex items-center justify-center text-surface-400 text-sm text-center"
							>
								You are not currently in any groups
								<iconify-icon
									icon="fa6-solid:face-frown"
									class="ml-1"
								/>
							</span>
						{/if}
					{/if}
				</section>
			</div>
			<div class="card col-span-2">
				<header
					class="card-header flex flex-row items-center justify-between space-y-0 pb-2"
				>
					<span class="text-lg font-medium">Recent Activity</span>
					<iconify-icon icon="fa6-solid:file-lines" />
				</header>
				<section class="flex flex-col text-sm px-4 py-1">
					{#each $auditlog.slice(0, 5) as log}
						<p class="text-surface-400">
							Collection
							<span class="font-bold">{log.collection}</span>
							was {log.event}d by
							<span class="font-bold">
								{#if log.expand?.user.name === undefined}
									an admin
								{:else}
									{log.expand?.user.name}
								{/if}
							</span>
							on
							{#if log.created}
								{Intl.DateTimeFormat("en", {
									dateStyle: "full",
									timeStyle: "short",
								}).format(new Date(log.created))}
							{:else}
								N/A
							{/if}
						</p>
					{/each}
				</section>
			</div>
			<div class="card lg:col-span-1 col-span-2">
				<header
					class="card-header flex flex-row items-center justify-between space-y-0 pb-2"
				>
					<span class="text-lg font-medium">Provider Status</span>
					<iconify-icon icon="fa6-solid:cloud" />
				</header>
				<section class="flex flex-col text-sm font-bold px-4 py-1">
					{#each $providers as provider, i}
						<div
							class="flex flex-row gap-4 justify-between"
							use:popup={{
								event: "hover",
								target: "providerError" + i,
								placement: "bottom",
							}}
						>
							<span>{provider.name}</span>
							{#if provider.error === null}
								<iconify-icon
									icon="fa6-solid:check"
									class="text-green-500"
								/>
							{:else}
								<iconify-icon
									icon="fa6-solid:xmark"
									class="text-red-500"
								/>

								<div
									class="card p-4 variant-filled-error"
									data-popup={"providerError" + i}
								>
									<p>{provider.error}</p>
								</div>
							{/if}
						</div>
					{/each}
					<button
						class="btn btn-sm mt-4 mb-1 bg-surface-300 dark:bg-surface-500"
						on:click={syncProviders}>Resync</button
					>
				</section>
			</div>
		</div>
	{/if}
</div>
