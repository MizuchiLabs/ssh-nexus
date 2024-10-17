<script lang="ts">
	import { pb } from "$lib/client";
	import { clipboard, getModalStore } from "@skeletonlabs/skeleton";
	import type { SvelteComponent } from "svelte";

	export let parent: SvelteComponent;

	const modalStore = getModalStore();

	let hostname = "";
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
				let response = await pb.send("/api/ssh/host/sign", {
					method: "POST",
					body: {
						publickey: publicKey,
						hostname: hostname ?? "unknown",
					},
				});
				signedCertificate = response.certificate;
				expiryDate = Intl.DateTimeFormat("en", {
					dateStyle: "full",
					timeStyle: "short",
				}).format(new Date(response.expiry * 1000));
			} catch (error) {
				signFail = true;
				publicKey = "";
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
					Signed Host Certificate
				</header>
				<section class="flex flex-col p-4 gap-4">
					<code class="code text-wrap">
						Warning! You will not be able to access this information
						later, so please copy the information below!
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
					Get a new host certificate
				</header>
				<section class="flex flex-col gap-4 p-4">
					<label class="label" on:keydown={onKeys} aria-hidden>
						<span>Public Key</span>
						<input
							class="input variant-form-material"
							type="text"
							bind:value={publicKey}
							placeholder="ssh-ed25519 ..."
							class:input-success={validKey}
							class:input-error={!validKey && publicKey !== ""}
						/>
					</label>
					<label class="label" on:keydown={onKeys} aria-hidden>
						<span> Hostname </span>
						<span class="text-sm text-surface-300 ml-1"
							>(Optional)</span
						>
						<input
							class="input variant-form-material"
							type="text"
							bind:value={hostname}
							placeholder="localhost"
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
						<code>ssh-ed25519 AAAAB3Nza... username@hostname</code>
					</p>
				</section>
			{/if}
		</div>
	</div>
{/if}
