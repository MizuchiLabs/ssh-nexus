<script lang="ts">
    import * as Avatar from "$lib/components/ui/avatar/index.js";
    import * as DropdownMenu from "$lib/components/ui/dropdown-menu/index.js";
    import { avatar, user, isAdmin, unsubscribeAuth } from "$lib/client";
    import { unsubscribeRecords } from "$lib/subscriptions";
    import { LogOut, User, Moon, Sun } from "lucide-svelte";
    import { onMount } from "svelte";

    const logout = () => {
        unsubscribeAuth();
        unsubscribeRecords();
    };

    let darkMode = false;
    function handleSwitchDarkMode() {
        darkMode = !darkMode;
        localStorage.setItem("mode", darkMode ? "dark" : "light");

        darkMode
            ? document.documentElement.classList.add("dark")
            : document.documentElement.classList.remove("dark");
    }

    onMount(async () => {
        if (
            localStorage.mode === "dark" ||
            (!("mode" in localStorage) &&
                window.matchMedia("(prefers-color-scheme: dark)").matches)
        ) {
            darkMode = true;
        } else {
            darkMode = false;
        }
    });
</script>

<DropdownMenu.Root>
    <DropdownMenu.Trigger class="ml-4">
        <Avatar.Root>
            <Avatar.Image src={$avatar} />
            <Avatar.Fallback
                class="bg-gradient-to-br from-pink-500 to-violet-500"
            >
                {#if $user?.username}
                    {$user?.username?.slice(0, 2).toUpperCase()}
                {:else if $isAdmin}
                    AD
                {/if}
            </Avatar.Fallback>
        </Avatar.Root>
    </DropdownMenu.Trigger>
    <DropdownMenu.Content>
        <DropdownMenu.Group>
            <DropdownMenu.Item
                href="/profile"
                class="flex flex-row items-center gap-2"
            >
                <User size="1rem" />
                <span>Profile</span>
            </DropdownMenu.Item>
            <DropdownMenu.Item
                class="flex flex-row items-center gap-2"
                on:click={handleSwitchDarkMode}
            >
                {#if darkMode}
                    <Sun size="1rem" />
                    Light
                {:else}
                    <Moon size="1rem" />
                    Dark
                {/if}
            </DropdownMenu.Item>
            <DropdownMenu.Item
                class="flex flex-row items-center gap-2"
                on:click={logout}
            >
                <LogOut size="1rem" />
                <span>Logout</span>
            </DropdownMenu.Item>
        </DropdownMenu.Group>
    </DropdownMenu.Content>
</DropdownMenu.Root>
