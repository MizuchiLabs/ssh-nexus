<script lang="ts">
    import { Button } from "$lib/components/ui/button/index";
    import { isAdmin, pb } from "$lib/client";
    import logo from "$lib/assets/logo.png";
    import Requests from "$lib/dropdowns/Requests.svelte";
    import Profile from "$lib/dropdowns/Profile.svelte";
    import Verification from "$lib/utils/Verification.svelte";
    import SignUserKey from "$lib/modals/SignUserKey.svelte";
    import { baseRoutes } from "$lib/navigation/Routes";
    import { page } from "$app/stores";
    import { Database, KeyRound } from "lucide-svelte";

    $: active = $page.url.pathname.split("/")[1];
    let open = false;
</script>

<Verification />

<SignUserKey bind:open />

<nav class="flex h-16 items-center justify-between border-b bg-surface-500/5">
    <div class="ml-4 flex flex-row items-center">
        <a href="/" class="hidden md:block" tabindex="-1">
            <img
                src={logo}
                alt="logo"
                class="mx-2 hover:brightness-75"
                width="38"
                height="38"
            />
        </a>

        <!-- Navigation Left -->
        <nav
            class="hidden flex-col gap-6 ml-4 items-center font-mono md:flex md:flex-row"
        >
            {#each baseRoutes as route}
                {#if $isAdmin && route.admin}
                    <a
                        href={route.path}
                        class="hover:text-red-400 hover:dark:text-red-300"
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
                        class="hover:text-red-400 hover:dark:text-red-300"
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
    </div>

    <div class="mr-2 flex flex-row items-center gap-2">
        <div class="flex items-center gap-1 mr-4">
            {#if pb.authStore.isAdmin}
                <Button
                    variant="ghost"
                    href={pb.baseUrl + "/_/"}
                    target="_blank"
                    size="icon"
                    class="h-8 w-8 rounded-full"
                >
                    <Database size="1rem" />
                </Button>
            {/if}
            <Button
                variant="ghost"
                size="icon"
                on:click={() => (open = true)}
                class="h-8 w-8 rounded-full"
            >
                <KeyRound size="1rem" />
            </Button>
            <Requests />
            <Profile />
        </div>
    </div>
</nav>
