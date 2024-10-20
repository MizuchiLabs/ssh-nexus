<script lang="ts">
	import { pb, user, avatar, isAdmin } from "$lib/client";
	import { toast } from "svelte-sonner";
	import * as Card from "$lib/components/ui/card";
	import * as Avatar from "$lib/components/ui/avatar/index.js";
	import { Button } from "$lib/components/ui/button/index.js";
	import { Input } from "$lib/components/ui/input/index.js";
	import { Label } from "$lib/components/ui/label/index.js";
	import type { ClientResponseError } from "pocketbase";
	import SSHConfigModal from "$lib/modals/SSHConfigModal.svelte";
	import { Bomb } from "lucide-svelte";

	let open = false;
	let sshKey: string = $user?.settings?.ssh_key_name || "";

	const updateAvatar = async (e: any) => {
		let file = e.target.files[0];
		let extension = file.name.split(".").pop();
		let filename = pb.authStore.model?.id + "_avatar" + "." + extension;

		const formData = new FormData();
		formData.append("file", file, filename);
		try {
			await pb.collection("users").update($user?.id, {
				avatar: formData.get("file"),
			});
			toast.success("Updated avatar");
		} catch (error: ClientResponseError | any) {
			toast.error(error.data?.message || "Something went wrong.");
		}
	};

	const deleteAvatar = async () => {
		try {
			await pb.collection("users").update($user?.id, { avatar: null });
			toast.success("Cleared avatar");
		} catch (error: ClientResponseError | any) {
			toast.error(error.data?.message || "Something went wrong.");
		}
	};

	const updateProfile = async () => {
		let settings = { ...$user?.settings, ssh_key_name: sshKey };
		await pb.collection("users").update($user?.id, {
			name: $user?.name,
			email: $user?.email,
			settings: settings,
		});
		toast.success("Updated profile");
	};
</script>

<SSHConfigModal bind:open />

{#if $user}
	<Card.Root class="max-w-md mx-auto p-2 shadow-md rounded-lg mt-8">
		<Card.Header class="flex items-center space-x-4">
			<Avatar.Root class=" w-16 h-16 rounded-full">
				<Avatar.Image src={$avatar} alt={"@" + $user.username} />
				<Avatar.Fallback class="text-xl">
					{$user.username?.slice(0, 2).toUpperCase()}
				</Avatar.Fallback>
			</Avatar.Root>
			<div class="flex flex-col items-center">
				<Card.Title class="text-xl font-semibold">{$user.name}</Card.Title>
				<Card.Description class="text-gray-500">
					@{$user.username}
				</Card.Description>
			</div>
		</Card.Header>

		<Card.Content class="mt-4 flex flex-col gap-4">
			<div class="flex flex-col items-start gap-1">
				<Label for="name">Name</Label>
				<Input id="name" type="text" bind:value={$user.name} />
			</div>

			<div class="flex flex-col items-start gap-1">
				<Label for="username">Username</Label>
				<Input id="username" type="text" bind:value={$user.username} />
			</div>

			<div class="flex flex-col items-start gap-1">
				<Label for="email">Email</Label>
				<Input id="email" type="email" bind:value={$user.email} />
			</div>

			{#if !$isAdmin}
				<div class="flex flex-col items-start gap-1">
					<Label for="avatar">Upload Avatar</Label>
					<div class="flex flex-row items-center justify-end gap-1">
						<Input id="avatar" type="file" on:change={updateAvatar} />
						<Button
							variant="ghost"
							class="bg-red-300 rounded-md hover:bg-red-500 text-black"
							on:click={deleteAvatar}
						>
							<Bomb size="1rem" />
						</Button>
					</div>
				</div>
			{/if}

			<div class="flex flex-col items-start gap-1">
				<Label for="principal">Principal ID</Label>
				<Input
					id="principal"
					type="text"
					bind:value={$user.principal}
					readonly
					disabled
				/>
			</div>
		</Card.Content>

		<Card.Footer class="grid grid-cols-2 gap-2">
			<Button
				variant="ghost"
				class="bg-purple-300 rounded-md hover:bg-purple-500 text-black"
				on:click={() => (open = true)}
			>
				Show SSH Config
			</Button>
			<Button
				variant="ghost"
				class="bg-red-300 rounded-md hover:bg-red-500 text-black"
				on:click={updateProfile}
			>
				Save Changes
			</Button>
		</Card.Footer>
	</Card.Root>
{/if}
