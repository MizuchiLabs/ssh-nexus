<script lang="ts">
    import {
        AppBar,
        getModalStore,
        type ModalSettings,
    } from "@skeletonlabs/skeleton";
    import { isAdmin, loggedIn, pb } from "$lib/client";
    import logo from "$lib/assets/logo.png";
    import Drawer from "$lib/navigation/Drawer.svelte";
    import RequestDD from "$lib/dropdowns/RequestDD.svelte";
    import ProfileDD from "$lib/dropdowns/ProfileDD.svelte";
    import Verification from "$lib/utils/Verification.svelte";
    import GlobalSearch from "$lib/utils/GlobalSearch.svelte";
    import { baseRoutes } from "$lib/navigation/Routes";
    import { page } from "$app/stores";

    let drawer: Drawer;

    $: active = $page.url.pathname.split("/")[1];

    const modals = getModalStore();
    const signUserKey = () => {
        const modal: ModalSettings = {
            type: "component",
            title: "Sign User Key",
            component: "SignUserKey",
        };
        modals.trigger(modal);
    };
</script>

<Drawer bind:this={drawer} />

<Verification />

<AppBar>
    <svelte:fragment slot="lead">
        <!-- Logo & Drawer -->
        <a href="/" class="hidden md:block" tabindex="-1">
            <img
                src={logo}
                alt="logo"
                class="mx-2 hover:brightness-75"
                width="38"
                height="38"
            />
        </a>
        <button
            class="block md:hidden"
            on:click={() => drawer.open()}
            tabindex="-1"
        >
            <img
                src={logo}
                alt="logo"
                class="mx-2 hover:brightness-75"
                width="38"
                height="38"
            />
        </button>

        <!-- Navigation Left -->
        {#if $loggedIn}
            <nav
                class="hidden flex-col gap-6 ml-4 items-center font-mono md:flex md:flex-row"
            >
                {#each baseRoutes as route}
                    {#if $isAdmin && route.admin}
                        <a
                            href={route.path}
                            class="font-bold hover:text-primary-400 hover:dark:text-primary-300"
                            class:text-primary-400={active ===
                                route.path.split("/")[1]}
                            class:dark:text-primary-300={active ===
                                route.path.split("/")[1]}
                            tabindex="-1"
                        >
                            {route.name}
                        </a>
                    {:else if !route.admin}
                        <a
                            href={route.path}
                            class="font-bold hover:text-primary-400 hover:dark:text-primary-300"
                            class:text-primary-400={active ===
                                route.path.split("/")[1]}
                            class:dark:text-primary-300={active ===
                                route.path.split("/")[1]}
                            tabindex="-1"
                        >
                            {route.name}
                        </a>
                    {/if}
                {/each}
            </nav>
        {/if}
    </svelte:fragment>

    <!-- TODO: <GlobalSearch />-->

    <!-- Navigation Right -->
    <svelte:fragment slot="trail">
        <div class="flex items-center gap-2 mr-4">
            {#if $loggedIn}
                {#if pb.authStore.isAdmin}
                    <a
                        href={pb.baseUrl + "/_/"}
                        class="btn-icon btn-icon-sm variant-soft-primary hover:variant-ghost-primary"
                    >
                        <iconify-icon icon="fa6-solid:database" width="12" />
                    </a>
                {/if}
                <button
                    class="btn-icon btn-icon-sm variant-soft-primary hover:brightness-50"
                    on:click={signUserKey}
                >
                    <iconify-icon icon="fa6-solid:fingerprint" width="16" />
                </button>
                <RequestDD />
                <ProfileDD />
            {:else}
                <a href="/login" class="btn-icon" tabindex="-1">
                    <iconify-icon icon="tabler:login" />
                </a>
            {/if}
        </div>
    </svelte:fragment>
</AppBar>
