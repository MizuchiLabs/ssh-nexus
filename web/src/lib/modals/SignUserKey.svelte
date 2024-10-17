<script lang="ts">
	import { pb } from "$lib/client";
	import { clipboard, getModalStore } from "@skeletonlabs/skeleton";
	import type { SvelteComponent } from "svelte";

	export let parent: SvelteComponent;

	const modalStore = getModalStore();

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
		if (e.key === "Escape") {
			modalStore.close();
		}
	};

	const selectText = (e: any) => {
		e.target.select();
		navigator.clipboard.writeText(e.target.value);
	};
</script>

{#if $modalStore[0]}
	<div class="container flex flex-col mx-auto p-4 gap-4 lg:w-2/3 xl:w-1/2">
		<div class="card">
			{#if signedCertificate}
				<header class="card-header text-2xl font-bold">
					Signed Certificate
				</header>
				<section class="flex flex-col p-4 gap-4">
					<code class="code text-wrap">
						You will not be able to access this information later,
						so please copy the information below!
					</code>
					<textarea
						class="textarea"
						rows="11"
						on:click={selectText}
						bind:value={signedCertificate}
						data-clipboard="signedCertificate"
						readonly
					/>
					<span class="text-sm ml-2">Expiry: {expiryDate}</span>
					<button
						class="btn variant-filled-success ml-2"
						use:clipboard={{ input: "signedCertificate" }}
						on:click={parent.onClose()}
						value="signedCertificate"
						>Copy
					</button>
				</section>
			{:else}
				<header class="card-header text-2xl font-bold">
					Get a new user certificate
				</header>
				<section class="p-4">
					<p class="text-sm mb-4">
						Enter your public SSH key below to get it signed.
						Signing your SSH key allows you to securely authenticate
						with our servers. Make sure your public key starts with
						<span class="font-mono"
							>ssh-ed25519, ssh-rsa, ssh-ecdsa,...</span
						> or any other supported key type.
					</p>
					<label class="label" on:keydown={onKeys} aria-hidden>
						<input
							class="input variant-form-material"
							type="text"
							bind:value={publicKey}
							placeholder="ssh-ed25519 ..."
							class:input-success={validKey}
							class:input-error={!validKey && publicKey !== ""}
						/>
					</label>
					<p
						class="text-sm text-error-500 mt-2"
						class:hidden={!signFail}
					>
						Failed to sign key, make sure you have entered a valid
						public key!
					</p>
					<p class="text-sm text-surface-400">
						Example of a valid public key format: <br />
						<code
							>ssh-ed25519 AAAAB3Nza... your_email@example.com</code
						>
					</p>
				</section>
			{/if}
		</div>
	</div>
{/if}
