<script lang="ts">
	import { onMount } from "svelte";
	import { generateConfig } from "$lib/utils/SSHConfig";
	import * as Dialog from "$lib/components/ui/dialog/index.js";
	import { Button } from "$lib/components/ui/button/index.js";
	import { Textarea } from "$lib/components/ui/textarea";

	export let open = false;

	let sshConfig = "";
	const selectText = (e: any) => {
		e.target.select();
		navigator.clipboard.writeText(e.target.value);
	};
	onMount(async () => {
		sshConfig = (await generateConfig()) || "";
	});
</script>

<Dialog.Root bind:open>
	<Dialog.Content class="no-scrollbar max-h-[80vh] overflow-y-auto">
		<Dialog.Header>
			<Dialog.Title>SSH Config</Dialog.Title>
			<Dialog.Description>This is your current ssh config.</Dialog.Description>
		</Dialog.Header>

		<Textarea
			class="textarea"
			rows={sshConfig.split("\n").length}
			bind:value={sshConfig}
			on:click={selectText}
			data-clipboard="sshConfig"
			readonly
		/>

		<Button
			class="w-full"
			on:click={() => navigator.clipboard.writeText(sshConfig)}>Copy</Button
		>
	</Dialog.Content>
</Dialog.Root>
