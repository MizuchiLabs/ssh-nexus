<script lang="ts">
	import { pb } from "$lib/client";
	import { groups, permissions } from "$lib/subscriptions";
	import * as Dialog from "$lib/components/ui/dialog/index.js";
	import * as Command from "$lib/components/ui/command/index.js";
	import * as Popover from "$lib/components/ui/popover/index.js";
	import { Button } from "$lib/components/ui/button/index.js";
	import { Input } from "$lib/components/ui/input/index.js";
	import { Label } from "$lib/components/ui/label/index.js";
	import type { ClientResponseError, RecordModel } from "pocketbase";
	import { Check, ChevronsUpDown } from "lucide-svelte";
	import { cn } from "$lib/utils.js";
	import { toast } from "svelte-sonner";

	export let user: RecordModel = {} as RecordModel;
	export let open = false;

	const update = async () => {
		try {
			if (!user.id) {
				await pb.collection("users").create(user);
				toast.success(`Created user ${user.name}`);
				open = false;
				return;
			}
			await pb.collection("users").update(user.id, user);
			toast.success(`Updated user ${user.name}`);
			open = false;
		} catch (error: ClientResponseError | any) {
			toast.error(error.data?.message || "Something went wrong.");
		}
	};

	const toggleGroup = (id: string) => {
		if (!user.groups) user.groups = [];
		if (!user.groups?.includes(id)) {
			user.groups.push(id);
		} else {
			user.groups = user.groups.filter((group: string) => group !== id);
		}
	};
</script>

<Dialog.Root bind:open>
	<Dialog.Content>
		<Dialog.Header>
			<Dialog.Title>
				{#if user.id}Edit{:else}Create{/if}
				User
			</Dialog.Title>
			<Dialog.Description>Provide the details of this user.</Dialog.Description>
		</Dialog.Header>

		<div class="grid gap-4 py-4">
			<div class="grid grid-cols-4 items-center gap-4">
				<Label for="name" class="text-right">Name</Label>
				<Input id="name" class="col-span-3" bind:value={user.name} />
			</div>
			<div class="grid grid-cols-4 items-center gap-4">
				<Label for="username" class="text-right">Username</Label>
				<Input id="username" class="col-span-3" bind:value={user.username} />
			</div>
			<div class="grid grid-cols-4 items-center gap-4">
				<Label for="email" class="text-right">Email</Label>
				<Input id="email" class="col-span-3" bind:value={user.email} />
			</div>

			{#if user.id}
				<div class="grid grid-cols-4 items-center gap-4">
					<Label for="principal" class="text-right">Principal</Label>
					<Input
						id="principal"
						class="col-span-3"
						bind:value={user.principal}
						placeholder="N/A"
						disabled
					/>
				</div>
			{/if}

			<!-- Permission -->
			<div class="grid grid-cols-4 items-center gap-4">
				<Popover.Root>
					<Popover.Trigger asChild let:builder>
						<Label for="permission" class="text-right">Permission</Label>
						<Button
							builders={[builder]}
							variant="outline"
							role="combobox"
							aria-expanded={open}
							class="col-span-3 justify-between"
						>
							{#if user.permission}
								{user.expand?.permission?.name}
							{:else}
								Select permission...
							{/if}
							<ChevronsUpDown class="ml-2 h-4 w-4 shrink-0 opacity-50" />
						</Button>
					</Popover.Trigger>
					<Popover.Content
						class="col-span-3 no-scrollbar p-0 max-h-[300px] overflow-y-auto"
						side="bottom"
					>
						<Command.Root>
							<Command.Input placeholder="Search permissions..." />
							<Command.Empty>No permission found.</Command.Empty>
							<Command.Group>
								{#each $permissions as permission}
									<Command.Item
										value={permission.name}
										onSelect={() => (user.permission = permission.id)}
									>
										<Check
											size="1rem"
											class={cn(
												"mr-2",
												user.permission?.includes(permission.id)
													? "opacity-100"
													: "opacity-0",
											)}
										/>
										{permission.name}
									</Command.Item>
								{/each}
							</Command.Group>
						</Command.Root>
					</Popover.Content>
				</Popover.Root>
			</div>

			<!-- Groups -->
			<div class="grid grid-cols-4 items-center gap-4">
				<Popover.Root>
					<Popover.Trigger asChild let:builder>
						<Label for="group" class="text-right">Groups</Label>
						<Button
							builders={[builder]}
							variant="outline"
							role="combobox"
							aria-expanded={open}
							class="col-span-3 justify-between"
						>
							Select groups...
							<ChevronsUpDown class="ml-2 h-4 w-4 shrink-0 opacity-50" />
						</Button>
					</Popover.Trigger>
					<Popover.Content
						class="col-span-3 no-scrollbar p-0 max-h-[300px] overflow-y-auto"
						side="bottom"
					>
						<Command.Root>
							<Command.Input placeholder="Search groups..." />
							<Command.Empty>No group found.</Command.Empty>
							<Command.Group>
								{#each $groups as group}
									<Command.Item
										value={group.name}
										onSelect={() => toggleGroup(group.id)}
									>
										<Check
											size="1rem"
											class={cn(
												"mr-2",
												user.groups?.includes(group.id)
													? "opacity-100"
													: "opacity-0",
											)}
										/>
										{group.name}
									</Command.Item>
								{/each}
							</Command.Group>
						</Command.Root>
					</Popover.Content>
				</Popover.Root>
			</div>
		</div>

		<Button class="w-full" on:click={update}>Save</Button>
	</Dialog.Content>
</Dialog.Root>
