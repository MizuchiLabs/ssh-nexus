<script lang="ts">
	import { pb } from "$lib/client";
	import { toast } from "svelte-sonner";
	import * as Dialog from "$lib/components/ui/dialog/index.js";
	import { Button } from "$lib/components/ui/button/index.js";
	import { Textarea } from "$lib/components/ui/textarea";

	export let open = false;

	let key = "";
	const update = async () => {
		try {
			await pb.send("/api/ssh/user/set", {
				method: "POST",
				body: { key: key },
			});
			toast.success("Updated user CA");
			open = false;
		} catch (error: any) {
			toast.error(
				error.message ||
					"Failed to update key, make sure you are using the correct openssh format!",
			);
		}
	};
</script>

<Dialog.Root bind:open>
	<Dialog.Content>
		<Dialog.Header>
			<Dialog.Title>Set a custom user CA</Dialog.Title>
		</Dialog.Header>
		<Textarea
			rows={10}
			bind:value={key}
			placeholder="-----BEGIN OPENSSH PRIVATE KEY-----
<private key>
-----END OPENSSH PRIVATE KEY-----"
		/>
		<Button class="w-full" on:click={update}>Overwrite Key</Button>
	</Dialog.Content>
</Dialog.Root>
