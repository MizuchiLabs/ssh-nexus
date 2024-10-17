<script lang="ts">
    import { isAdmin, loggedIn } from "$lib/client";
    import {
        getDrawerStore,
        type DrawerSettings,
        Drawer,
    } from "@skeletonlabs/skeleton";
    import { onMount } from "svelte";
    import { baseRoutes, type Route } from "$lib/navigation/Routes";

    const drawerStore = getDrawerStore();
    const drawerSettings: DrawerSettings = {
        id: "example-3",
        bgDrawer: "bg-slate-800 text-white",
        bgBackdrop: "bg-gradient-to-tr from-purple-600/50 to-slate-800/50",
        width: "w-56",
        padding: "p-2",
        rounded: "rounded-xl",
    };
    let routes = [] as Route[];

    export const open = () => {
        drawerStore.open(drawerSettings);
    };

    export const close = () => {
        drawerStore.close();
    };

    onMount(async () => {
        routes = baseRoutes.filter((route) => !route.admin || isAdmin);
    });
</script>

{#if $loggedIn}
    <Drawer>
        <div
            class="lg:flex flex-col justify-between h-full bg-surface-500/5 w-56 p-4"
        >
            <nav class="list-nav flex flex-col gap-1">
                {#each routes as route}
                    {#if $isAdmin && route.admin}
                        <a href={route.path} on:click={close} tabindex="-1">
                            <iconify-icon icon={route.icon} class="mr-2 w-5" />
                            {route.name}
                        </a>
                    {:else if !route.admin}
                        <a href={route.path} on:click={close} tabindex="-1">
                            <iconify-icon icon={route.icon} class="mr-2 w-5" />
                            {route.name}
                        </a>
                    {/if}
                {/each}
            </nav>
        </div>
    </Drawer>
{/if}
