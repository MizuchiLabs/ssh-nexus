<script lang="ts">
	import { pb } from "$lib/client";
	import * as Dialog from "$lib/components/ui/dialog/index.js";
	import { Button } from "$lib/components/ui/button/index.js";
	import { Input } from "$lib/components/ui/input/index.js";
	import { Label } from "$lib/components/ui/label/index.js";
	import { Textarea } from "$lib/components/ui/textarea";
	import { Copy } from "lucide-svelte";

	export let open = false;

	let publicKey = "";
	let expiryDate = "";
	let signedCertificate = "";
	let validKey = false;
	let signFail = false;
	const onKeys = async (e: KeyboardEvent) => {
		const sshKeyPattern =
			/ssh-(ed25519|rsa|dss|ecdsa) AAAA(?:[A-Za-z0-9+\/]{4})*(?:[A-Za-z0-9+\/]{2}==|[A-Za-z0-9+\/]{3}=|[A-Za-z0-9+\/]{4})( [^@]+@[^@]+)?/;
		validKey = sshKeyPattern.test(publicKey);

		if (e.key === "Enter" && publicKey && validKey) {
			try {
				let response = await pb.send("/api/ssh/user/sign", {
					method: "POST",
					body: { publickey: publicKey },
				});
				signedCertificate = response.certificate;
				expiryDate = Intl.DateTimeFormat("en", {
					dateStyle: "full",
					timeStyle: "short",
				}).format(new Date(response.expiry * 1000));
			} catch (error) {
				signFail = true;
			}
		}
	};

	const selectText = (e: any) => {
		e.target.select();
		navigator.clipboard.writeText(e.target.value);
	};
</script>

<Dialog.Root bind:open>
	<Dialog.Content>
		<Dialog.Header>
			<Dialog.Title>
				{#if signedCertificate}
					Signed Certificate
				{:else}
					Sign Key
				{/if}
			</Dialog.Title>
			<Dialog.Description>
				{#if signedCertificate}
					You will not be able to access this information later, so please copy
					the information below!
				{:else}
					Get a new user certificate by signing your public key.
				{/if}
			</Dialog.Description>
		</Dialog.Header>
		{#if signedCertificate}
			<Textarea
				rows={10}
				bind:value={signedCertificate}
				on:click={selectText}
				readonly
			/>

			<div class="flex flex-row items-center justify-between gap-4">
				<span class="text-sm ml-2">Expiry: {expiryDate}</span>

				<Button
					variant="ghost"
					size="icon"
					class="h-8 hover:bg-transparent hover:text-red-400"
					on:click={() => navigator.clipboard.writeText(signedCertificate)}
				>
					<Copy size="1rem" />
				</Button>
			</div>
		{:else}
			<div class="flex flex-col gap-4" on:keydown={onKeys} aria-hidden>
				<div class="flex flex-row items-center gap-4">
					<Label for="publickey" class="text-right min-w-[80px]"
						>Public Key</Label
					>
					<Input
						id="publickey"
						class="col-span-3"
						bind:value={publicKey}
						placeholder="ssh-ed25519 ..."
					/>
				</div>
				{#if signFail}
					<p class="text-xs text-red-400 text-right">
						Failed to sign key, make sure you have entered a valid public key!
					</p>
				{/if}
			</div>
		{/if}
	</Dialog.Content>
</Dialog.Root>
