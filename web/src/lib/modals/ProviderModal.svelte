<script lang="ts">
	import { pb } from "$lib/client";
	import * as Dialog from "$lib/components/ui/dialog/index.js";
	import * as Select from "$lib/components/ui/select";
	import { Button } from "$lib/components/ui/button/index.js";
	import { Input } from "$lib/components/ui/input/index.js";
	import { Label } from "$lib/components/ui/label/index.js";
	import type { ClientResponseError, RecordModel } from "pocketbase";
	import { toast } from "svelte-sonner";

	export let provider: RecordModel = {} as RecordModel;
	export let open = false;
	const providers = [
		"aws",
		"azure",
		"google",
		"hetzner",
		"linode",
		"vultr",
		"proxmox",
	];

	const update = async () => {
		try {
			if (!provider.id) {
				await pb.collection("providers").create(provider);
				toast.success(`Created provider ${provider.name}`);
				open = false;
				return;
			}
			await pb.collection("providers").update(provider.id, provider);
			toast.success(`Updated provider ${provider.name}`);
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
				{#if provider.id}Edit{:else}Create{/if}
				Provider
			</Dialog.Title>
			<Dialog.Description
				>Provide the details of this provider.</Dialog.Description
			>
		</Dialog.Header>

		<div class="grid gap-4 py-4">
			<div class="grid grid-cols-4 items-center gap-4">
				<Label for="type" class="text-right">Type</Label>
				<div class="col-span-3 flex items-center flex-row">
					<Select.Root
						selected={{
							value: provider.type,
							label: provider.type?.toString(),
						}}
						onSelectedChange={(e) => e && (provider.type = e.value)}
					>
						<Select.Trigger>
							<Select.Value placeholder="Provider" />
						</Select.Trigger>
						<Select.Content>
							{#each providers as provider}
								<Select.Item value={provider} label={provider.toString()}
									>{provider}</Select.Item
								>
							{/each}
						</Select.Content>
					</Select.Root>
				</div>
			</div>
			<div class="grid grid-cols-4 items-center gap-4">
				<Label for="name" class="text-right">Name</Label>
				<Input id="name" class="col-span-3" bind:value={provider.name} />
			</div>
			<div class="grid grid-cols-4 items-center gap-4">
				<Label for="url" class="text-right">URL</Label>
				<Input id="url" class="col-span-3" bind:value={provider.url} />
			</div>
			<div class="grid grid-cols-4 items-center gap-4">
				<Label for="username" class="text-right">Username</Label>
				<Input
					id="username"
					class="col-span-3"
					bind:value={provider.username}
				/>
			</div>
			<div class="grid grid-cols-4 items-center gap-4">
				<Label for="password" class="text-right">Password</Label>
				<Input
					id="password"
					class="col-span-3"
					bind:value={provider.password}
				/>
			</div>
			<div class="grid grid-cols-4 items-center gap-4">
				<Label for="token" class="text-right">Token</Label>
				<Input id="token" class="col-span-3" bind:value={provider.token} />
			</div>
		</div>

		<Button class="w-full" on:click={update}>Save</Button>
	</Dialog.Content>
</Dialog.Root>
