<script lang="ts">
    import { page } from "$app/stores";

    $: active = $page.url.pathname;
</script>

<div class="flex flex-col gap-4 min-w-52">
    <a
        href="/request/machines"
        class="btn variant-outline-surface"
        class:text-primary-300={active === "/request/machines/"}>Machines</a
    >
    <a
        href="/request/groups"
        class="btn variant-outline-surface"
        class:text-primary-300={active === "/request/groups/"}>Groups</a
    >
</div>
