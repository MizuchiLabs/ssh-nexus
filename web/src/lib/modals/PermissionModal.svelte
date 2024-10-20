<script lang="ts">
	import { pb } from "$lib/client";
	import { groups, machines, users } from "$lib/subscriptions";
	import * as Dialog from "$lib/components/ui/dialog/index.js";
	import * as Command from "$lib/components/ui/command/index.js";
	import * as Popover from "$lib/components/ui/popover/index.js";
	import { Button } from "$lib/components/ui/button/index.js";
	import { Input } from "$lib/components/ui/input/index.js";
	import { Label } from "$lib/components/ui/label/index.js";
	import { Toggle } from "$lib/components/ui/toggle";
	import { Switch } from "$lib/components/ui/switch/index.js";
	import type { ClientResponseError, RecordModel } from "pocketbase";
	import { Check, ChevronsUpDown } from "lucide-svelte";
	import { cn } from "$lib/utils.js";
	import { toast } from "svelte-sonner";

	export let permission: RecordModel = {} as RecordModel;
	export let open = false;

	const update = async () => {
		try {
			if (!permission.id) {
				await pb.collection("permissions").create(permission);
				toast.success(`Created permission ${permission.name}`);
				open = false;
				return;
			}
			await pb.collection("permissions").update(permission.id, permission);
			toast.success(`Updated permission ${permission.name}`);
			open = false;
		} catch (error: ClientResponseError | any) {
			toast.error(error.data?.message || "Something went wrong.");
		}
	};

	const toggleUser = (id: string) => {
		if (!permission.users) permission.users = [];
		if (!permission.users?.includes(id)) {
			permission.users.push(id);
		} else {
			permission.users = permission.users.filter((user: string) => user !== id);
		}
	};
	const toggleGroup = (id: string) => {
		if (!permission.groups) permission.groups = [];
		if (!permission.groups?.includes(id)) {
			permission.groups.push(id);
		} else {
			permission.groups = permission.groups.filter(
				(group: string) => group !== id,
			);
		}
	};
	const toggleMachine = (id: string) => {
		if (!permission.machines) permission.machines = [];
		if (!permission.machines?.includes(id)) {
			permission.machines.push(id);
		} else {
			permission.machines = permission.machines.filter(
				(machine: string) => machine !== id,
			);
		}
	};
</script>

<Dialog.Root bind:open>
	<Dialog.Content>
		<Dialog.Header>
			<Dialog.Title>
				{#if permission.id}Edit{:else}Create{/if}
				Permission
			</Dialog.Title>
			<Dialog.Description
				>Provide the details of this permission.</Dialog.Description
			>
		</Dialog.Header>

		<div class="grid gap-4 py-4">
			<div class="grid grid-cols-4 items-center gap-4">
				<Label for="name" class="text-right">Name</Label>
				<Input id="name" class="col-span-3" bind:value={permission.name} />
			</div>
			<div class="grid grid-cols-4 items-center gap-4">
				<Label for="description" class="text-right">Description</Label>
				<Input
					id="description"
					class="col-span-3"
					bind:value={permission.description}
				/>
			</div>
			<div class="grid grid-cols-4 items-center gap-4">
				<Label for="is_admin" class="text-right">Admin</Label>
				<Switch
					id="is_admin"
					class="col-span-3"
					bind:checked={permission.is_admin}
				/>
			</div>

			<div class="grid grid-cols-4 items-center gap-4">
				<Label for="access" class="text-right">Access</Label>
				<div class="col-span-3">
					<Toggle
						variant="outline"
						pressed={permission.access_users}
						onPressedChange={() =>
							(permission.access_users = !permission.access_users)}
					>
						Users
					</Toggle>
					<Toggle
						variant="outline"
						pressed={permission.access_machines}
						onPressedChange={() =>
							(permission.access_machines = !permission.access_machines)}
					>
						Machines
					</Toggle>
					<Toggle
						variant="outline"
						pressed={permission.access_groups}
						onPressedChange={() =>
							(permission.access_groups = !permission.access_groups)}
					>
						Groups
					</Toggle>
				</div>
			</div>

			<div class="grid grid-cols-4 items-center gap-4">
				<Label for="access" class="text-right">Type</Label>
				<div class="col-span-3">
					<Toggle
						variant="outline"
						pressed={permission.can_create}
						onPressedChange={() =>
							(permission.can_create = !permission.can_create)}
					>
						Create
					</Toggle>
					<Toggle
						variant="outline"
						pressed={permission.can_update}
						onPressedChange={() =>
							(permission.can_update = !permission.can_update)}
					>
						Update
					</Toggle>
					<Toggle
						variant="outline"
						pressed={permission.can_delete}
						onPressedChange={() =>
							(permission.can_delete = !permission.can_delete)}
					>
						Delete
					</Toggle>
				</div>
			</div>

			<!-- Users -->
			<div class="grid grid-cols-4 items-center gap-4">
				<Popover.Root>
					<Popover.Trigger asChild let:builder>
						<Label for="user" class="text-right">Users</Label>
						<Button
							builders={[builder]}
							variant="outline"
							role="combobox"
							aria-expanded={open}
							class="col-span-3 justify-between"
						>
							Select users...
							<ChevronsUpDown class="ml-2 h-4 w-4 shrink-0 opacity-50" />
						</Button>
					</Popover.Trigger>
					<Popover.Content
						class="col-span-3 no-scrollbar p-0 max-h-[300px] overflow-y-auto"
						side="bottom"
					>
						<Command.Root>
							<Command.Input placeholder="Search users..." />
							<Command.Empty>No user found.</Command.Empty>
							<Command.Group>
								{#each $users as user}
									<Command.Item
										value={user.name}
										onSelect={() => toggleUser(user.id)}
									>
										<Check
											size="1rem"
											class={cn(
												"mr-2",
												permission.users?.includes(user.id)
													? "opacity-100"
													: "opacity-0",
											)}
										/>
										{user.name}
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
												permission.groups?.includes(group.id)
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

			<!-- Machines -->
			<div class="grid grid-cols-4 items-center gap-4">
				<Popover.Root>
					<Popover.Trigger asChild let:builder>
						<Label for="user" class="text-right">Machines</Label>
						<Button
							builders={[builder]}
							variant="outline"
							role="combobox"
							aria-expanded={open}
							class="col-span-3 justify-between"
						>
							Select machines...
							<ChevronsUpDown class="ml-2 h-4 w-4 shrink-0 opacity-50" />
						</Button>
					</Popover.Trigger>
					<Popover.Content
						class="col-span-3 no-scrollbar p-0 max-h-[300px] overflow-y-auto"
						side="bottom"
					>
						<Command.Root>
							<Command.Input placeholder="Search machines..." />
							<Command.Empty>No user found.</Command.Empty>
							<Command.Group>
								{#each $machines as machine}
									<Command.Item
										value={machine.name}
										onSelect={() => toggleMachine(machine.id)}
									>
										<Check
											size="1rem"
											class={cn(
												"mr-2",
												permission.machines?.includes(machine.id)
													? "opacity-100"
													: "opacity-0",
											)}
										/>
										{machine.name}
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
