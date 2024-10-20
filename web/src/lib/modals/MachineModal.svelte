<script lang="ts">
	import { pb } from "$lib/client";
	import { users, groups, tags, machines } from "$lib/subscriptions";
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

	export let machine: RecordModel = {} as RecordModel;
	export let open = false;

	const update = async () => {
		try {
			if (!machine.id) {
				await pb.collection("machines").create(machine);
				toast.success(`Created machine ${machine.name}`);
				open = false;
				return;
			}
			await pb.collection("machines").update(machine.id, machine);
			toast.success(`Updated machine ${machine.name}`);
			open = false;
		} catch (error: ClientResponseError | any) {
			toast.error(error.data?.message || "Something went wrong.");
		}
	};

	const toggleGroup = (id: string) => {
		if (!machine.groups) machine.groups = [];
		if (!machine.groups?.includes(id)) {
			machine.groups.push(id);
		} else {
			machine.groups = machine.groups.filter((group: string) => group !== id);
		}
	};
	const toggleTag = (id: string) => {
		if (!machine.tags) machine.tags = [];
		if (!machine.tags?.includes(id)) {
			machine.tags.push(id);
		} else {
			machine.tags = machine.tags.filter((tag: string) => tag !== id);
		}
	};
	const toggleUser = (id: string) => {
		if (!machine.users) machine.users = [];
		if (!machine.users?.includes(id)) {
			machine.users.push(id);
		} else {
			machine.users = machine.users.filter((user: string) => user !== id);
		}
	};
</script>

<Dialog.Root bind:open>
	<Dialog.Content>
		<Dialog.Header>
			<Dialog.Title>
				{#if machine.id}Edit{:else}Create{/if}
				Machine
			</Dialog.Title>
			<Dialog.Description>
				Provide the details of this machine.
			</Dialog.Description>
		</Dialog.Header>

		<div class="grid gap-4 py-4">
			<div class="grid grid-cols-4 items-center gap-4">
				<Label for="name" class="text-right">Name</Label>
				<Input id="name" class="col-span-3" bind:value={machine.name} />
			</div>
			<div class="grid grid-cols-4 items-center gap-4">
				<Label for="host" class="text-right">Host</Label>
				<Input id="host" class="col-span-3" bind:value={machine.host} />
			</div>
			<div class="grid grid-cols-4 items-center gap-4">
				<Label for="port" class="text-right">Port</Label>
				<Input id="port" class="col-span-3" bind:value={machine.port} />
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
												machine.groups?.includes(group.id)
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

			<!-- Tags -->
			<div class="grid grid-cols-4 items-center gap-4">
				<Popover.Root>
					<Popover.Trigger asChild let:builder>
						<Label for="tag" class="text-right">Tags</Label>
						<Button
							builders={[builder]}
							variant="outline"
							role="combobox"
							aria-expanded={open}
							class="col-span-3 justify-between"
						>
							Select tags...
							<ChevronsUpDown class="ml-2 h-4 w-4 shrink-0 opacity-50" />
						</Button>
					</Popover.Trigger>
					<Popover.Content
						class="col-span-3 no-scrollbar p-0 max-h-[300px] overflow-y-auto"
						side="bottom"
					>
						<Command.Root>
							<Command.Input placeholder="Search tags..." />
							<Command.Empty>No tag found.</Command.Empty>
							<Command.Group>
								{#each $tags as tag}
									<Command.Item
										value={tag.name}
										onSelect={() => toggleTag(tag.id)}
									>
										<Check
											size="1rem"
											class={cn(
												"mr-2",
												machine.tags?.includes(tag.id)
													? "opacity-100"
													: "opacity-0",
											)}
										/>
										{tag.name}
									</Command.Item>
								{/each}
							</Command.Group>
						</Command.Root>
					</Popover.Content>
				</Popover.Root>
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
												machine.users?.includes(user.id)
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
		</div>
		<Button class="w-full" on:click={update}>Save</Button>
	</Dialog.Content>
</Dialog.Root>
