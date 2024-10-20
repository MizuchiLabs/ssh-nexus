<script lang="ts">
	import { pb, user } from "$lib/client";
	import { groups, machines } from "$lib/subscriptions";
	import * as Dialog from "$lib/components/ui/dialog/index.js";
	import * as Command from "$lib/components/ui/command/index.js";
	import * as Popover from "$lib/components/ui/popover/index.js";
	import { Button } from "$lib/components/ui/button/index.js";
	import { Label } from "$lib/components/ui/label/index.js";
	import type { ClientResponseError, RecordModel } from "pocketbase";
	import { Check, ChevronsUpDown } from "lucide-svelte";
	import { cn } from "$lib/utils.js";
	import { toast } from "svelte-sonner";

	export let request: RecordModel = {} as RecordModel;
	export let open = false;

	const create = async () => {
		if (!$user) return;

		try {
			for (const group of request.groups) {
				const groupRequest = {
					user: $user.id,
					description: request.description,
					group: group,
				};
				await pb.collection("requests").create(groupRequest);
			}

			for (const machine of request.machines) {
				const machineRequest = {
					user: $user.id,
					description: request.description,
					machine: machine,
				};
				await pb.collection("requests").create(machineRequest);
			}
			toast.success(`Sent request!`, {
				description: "An admin will get back to you shortly",
			});
			open = false;
		} catch (error: ClientResponseError | any) {
			toast.error(error.data?.message || "Something went wrong.");
		}
	};

	const toggleGroup = (id: string) => {
		if (!request.groups) request.groups = [];
		if (!request.groups?.includes(id)) {
			request.groups.push(id);
		} else {
			request.groups = request.groups.filter(
				(machine: string) => machine !== id,
			);
		}
	};

	const toggleMachine = (id: string) => {
		if (!request.machines) request.machines = [];
		if (!request.machines?.includes(id)) {
			request.machines.push(id);
		} else {
			request.machines = request.machines.filter(
				(machine: string) => machine !== id,
			);
		}
	};
</script>

<Dialog.Root bind:open>
	<Dialog.Content>
		<Dialog.Header>
			<Dialog.Title>
				{#if request.id}Edit{:else}Create{/if}
				Request
			</Dialog.Title>
			<Dialog.Description>
				Request access to one or more machines/groups.
			</Dialog.Description>
		</Dialog.Header>

		<div class="grid gap-4 py-4">
			<!-- Groups -->
			<div class="grid grid-cols-4 items-center gap-4">
				<Popover.Root>
					<Popover.Trigger asChild let:builder>
						<Label for="machine" class="text-right">Groups</Label>
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
							<Command.Empty>No machine found.</Command.Empty>
							<Command.Group>
								{#each $groups as machine}
									<Command.Item
										value={machine.name}
										onSelect={() => toggleGroup(machine.id)}
									>
										<Check
											size="1rem"
											class={cn(
												"mr-2",
												request.groups?.includes(machine.id)
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

			<!-- Machines -->
			<div class="grid grid-cols-4 items-center gap-4">
				<Popover.Root>
					<Popover.Trigger asChild let:builder>
						<Label for="machine" class="text-right">Machines</Label>
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
							<Command.Empty>No machine found.</Command.Empty>
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
												request.machines?.includes(machine.id)
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

		<Button class="w-full" on:click={create}>Send</Button>
	</Dialog.Content>
</Dialog.Root>
