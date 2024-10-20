<script lang="ts">
	import { pb } from "$lib/client";
	import * as Dialog from "$lib/components/ui/dialog/index.js";
	import { Button } from "$lib/components/ui/button/index.js";
	import { Input } from "$lib/components/ui/input/index.js";
	import { Label } from "$lib/components/ui/label/index.js";
	import type { ClientResponseError, RecordModel } from "pocketbase";
	import { toast } from "svelte-sonner";

	export let group: RecordModel = {} as RecordModel;
	export let open = false;

	const update = async () => {
		try {
			if (!group.id) {
				await pb.collection("groups").create(group);
				toast.success(`Created group ${group.name}`);
				open = false;
				return;
			}
			await pb.collection("groups").update(group.id, group);
			toast.success(`Updated group ${group.name}`);
			open = false;
		} catch (error: ClientResponseError | any) {
			toast.error(error.data?.message || "Something went wrong.");
		}
	};
</script>

<Dialog.Root bind:open>
	<Dialog.Content>
		<Dialog.Header>
			<Dialog.Title>
				{#if group.id}Edit{:else}Create{/if}
				Group
			</Dialog.Title>
			<Dialog.Description>
				Provide the details of this group.
			</Dialog.Description>
		</Dialog.Header>

		<div class="grid gap-4 py-4">
			<div class="grid grid-cols-4 items-center gap-4">
				<Label for="name" class="text-right">Name</Label>
				<Input id="name" class="col-span-3" bind:value={group.name} />
			</div>
			<div class="grid grid-cols-4 items-center gap-4">
				<Label for="description" class="text-right">Description</Label>
				<Input
					id="description"
					class="col-span-3"
					bind:value={group.description}
				/>
			</div>
			<div class="grid grid-cols-4 items-center gap-4">
				<Label for="username" class="text-right">Linux username</Label>
				<Input
					id="username"
					class="col-span-3"
					bind:value={group.linux_username}
				/>
			</div>
		</div>
		<Button class="w-full" on:click={update}>Save</Button>
	</Dialog.Content>
</Dialog.Root>
