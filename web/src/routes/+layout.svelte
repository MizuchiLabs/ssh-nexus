<script lang="ts">
	import "../app.postcss";
	import { loggedIn, pb, subscribeAuth, unsubscribeAuth } from "$lib/client";
	import { Toaster } from "$lib/components/ui/sonner";
	import Navbar from "$lib/navigation/Navbar.svelte";
	import Footer from "$lib/navigation/Footer.svelte";
	import autoAnimate from "@formkit/auto-animate";
	import { subscribeRecords, unsubscribeRecords } from "$lib/subscriptions";
	import { goto } from "$app/navigation";
	import { onMount } from "svelte";

	pb.authStore.onChange(() => {
		if (pb.authStore.isValid) {
			subscribeAuth();
			subscribeRecords();
		}
	});

	onMount(async () => {
		if (
			localStorage.mode === "dark" ||
			(!("mode" in localStorage) &&
				window.matchMedia("(prefers-color-scheme: dark)").matches)
		) {
			document.documentElement.classList.add("dark");
		} else {
			document.documentElement.classList.remove("dark");
		}

		try {
			// get an up-to-date auth store state by verifying and refreshing the loaded auth model (if any)
			if (pb.authStore.isValid) {
				if (pb.authStore.isAdmin) {
					await pb.admins.authRefresh();
				} else if (pb.authStore.isAuthRecord) {
					await pb.collection("users").authRefresh();
				}
			} else {
				goto("/login");
			}
		} catch (_) {
			// clear the auth store on failed refresh
			unsubscribeAuth();
			unsubscribeRecords();
		}
	});
</script>

<Toaster />

<div class="app flex min-h-screen flex-col">
	{#if $loggedIn}
		<Navbar />

		<main class="mb-auto px-6 py-4" use:autoAnimate={{ duration: 100 }}>
			<slot />
		</main>

		<Footer />
	{:else}
		<div
			class="flex h-screen flex-col items-center justify-center"
			use:autoAnimate={{ duration: 100 }}
		>
			<slot />
		</div>
	{/if}
</div>
