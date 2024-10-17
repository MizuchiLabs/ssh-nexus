<script lang="ts">
    import { avatar, unsubscribeAuth } from "$lib/client";
    import { unsubscribeRecords } from "$lib/subscriptions";
    import { popup } from "@skeletonlabs/skeleton";

    const logout = () => {
        unsubscribeAuth();
        unsubscribeRecords();
    };
</script>

<div
    use:popup={{ event: "click", target: "popupProfile", placement: "bottom" }}
>
    {#if $avatar}
        <button class="btn-icon">
            <img
                src={$avatar}
                class="rounded-full w-10 h-10 object-cover hover:brightness-50"
                alt="avatar"
            />
        </button>
    {:else}
        <button class="btn-icon bg-gradient-to-br from-pink-500 to-violet-500">
            <iconify-icon icon="tabler:user" class="text-white" />
        </button>
    {/if}
</div>

<div class="card z-10 rounded-md dark:bg-surface-500" data-popup="popupProfile">
    <a
        href="/profile"
        class="block px-4 py-2 text-sm hover:text-primary-400 hover:dark:text-primary-300"
        tabindex="-1">Profile</a
    >
    <button
        class="block px-4 py-2 text-sm hover:text-primary-400 hover:dark:text-primary-300"
        on:click={logout}
        tabindex="-1">Logout</button
    >
</div>
