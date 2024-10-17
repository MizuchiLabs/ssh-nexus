<script lang="ts">
	import "../app.postcss";
	import {
		Toast,
		initializeStores,
		type ModalComponent,
		Modal,
		autoModeWatcher,
	} from "@skeletonlabs/skeleton";
	import { loggedIn, pb, subscribeAuth, unsubscribeAuth } from "$lib/client";
	import MachineCreate from "$lib/modals/MachineCreate.svelte";
	import MachineUpdate from "$lib/modals/MachineUpdate.svelte";
	import GroupCreate from "$lib/modals/GroupCreate.svelte";
	import GroupUpdate from "$lib/modals/GroupUpdate.svelte";
	import UserUpdate from "$lib/modals/UserUpdate.svelte";
	import UserCreate from "$lib/modals/UserCreate.svelte";
	import PermissionUpdate from "$lib/modals/PermissionUpdate.svelte";
	import PermissionCreate from "$lib/modals/PermissionCreate.svelte";
	import ProviderUpdate from "$lib/modals/ProviderUpdate.svelte";
	import ProviderCreate from "$lib/modals/ProviderCreate.svelte";
	import SSHConfigModal from "$lib/modals/SSHConfigModal.svelte";
	import UpdateUserCA from "$lib/modals/UpdateUserCA.svelte";
	import SignUserKey from "$lib/modals/SignUserKey.svelte";
	import SignHostKey from "$lib/modals/SignHostKey.svelte";
	import IssueModal from "$lib/modals/IssueModal.svelte";

	const modalRegistry: Record<string, ModalComponent> = {
		MachineCreate: { ref: MachineCreate },
		MachineUpdate: { ref: MachineUpdate },
		GroupCreate: { ref: GroupCreate },
		GroupUpdate: { ref: GroupUpdate },
		UserUpdate: { ref: UserUpdate },
		UserCreate: { ref: UserCreate },
		PermissionUpdate: { ref: PermissionUpdate },
		PermissionCreate: { ref: PermissionCreate },
		ProviderUpdate: { ref: ProviderUpdate },
		ProviderCreate: { ref: ProviderCreate },
		SSHConfigModal: { ref: SSHConfigModal },
		UpdateUserCA: { ref: UpdateUserCA },
		SignUserKey: { ref: SignUserKey },
		SignHostKey: { ref: SignHostKey },
		IssueModal: { ref: IssueModal },
	};

	// Highlight JS
	import hljs from "highlight.js/lib/core";
	import "highlight.js/styles/github-dark.css";
	import { storeHighlightJs } from "@skeletonlabs/skeleton";
	import xml from "highlight.js/lib/languages/xml"; // for HTML
	import css from "highlight.js/lib/languages/css";
	import javascript from "highlight.js/lib/languages/javascript";
	import typescript from "highlight.js/lib/languages/typescript";

	hljs.registerLanguage("xml", xml); // for HTML
	hljs.registerLanguage("css", css);
	hljs.registerLanguage("javascript", javascript);
	hljs.registerLanguage("typescript", typescript);
	storeHighlightJs.set(hljs);

	// Floating UI for Popups
	import {
		computePosition,
		autoUpdate,
		flip,
		shift,
		offset,
		arrow,
	} from "@floating-ui/dom";
	import { storePopup } from "@skeletonlabs/skeleton";
	import { onMount } from "svelte";
	import { fade } from "svelte/transition";
	import { page } from "$app/stores";
	import { cubicOut } from "svelte/easing";
	import Navbar from "$lib/navigation/Navbar.svelte";
	import Footer from "$lib/navigation/Footer.svelte";
	import Fab from "$lib/utils/FAB.svelte";
	import { subscribeRecords, unsubscribeRecords } from "$lib/subscriptions";
	import { goto } from "$app/navigation";
	storePopup.set({ computePosition, autoUpdate, flip, shift, offset, arrow });
	initializeStores();

	pb.authStore.onChange(() => {
		if (pb.authStore.isValid) {
			subscribeAuth();
			subscribeRecords();
		}
	});

	onMount(async () => {
		autoModeWatcher(); // auto switch between light and dark mode based on OS setting
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

<Toast />
<Modal components={modalRegistry} />

<div class="flex flex-col h-screen overflow-y-auto">
	<Navbar />

	{#key $page.url.pathname}
		<main
			in:fade={{ easing: cubicOut, duration: 300, delay: 100 }}
			class="mb-auto"
		>
			<slot />
		</main>
	{/key}

	{#if $loggedIn}
		<Footer />
	{/if}
</div>
