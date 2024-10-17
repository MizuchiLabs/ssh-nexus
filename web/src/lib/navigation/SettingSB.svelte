<script lang="ts">
    import { page } from "$app/stores";

    $: active = $page.url.pathname;
</script>

<div class="flex flex-col gap-4 min-w-52 font-mono font-light">
    <a
        href="/settings/general"
        class="btn variant-outline-surface hover:text-primary-300"
        class:text-primary-400={active === "/settings/general/"}
        class:dark:text-primary-300={active === "/settings/general/"}
    >
        General
    </a>
    <a
        href="/settings/ssh"
        class="btn variant-outline-surface hover:text-primary-300"
        class:text-primary-400={active === "/settings/ssh/"}
        class:dark:text-primary-300={active === "/settings/ssh/"}>Keys</a
    >
    <a
        href="/settings/permissions"
        class="btn variant-outline-surface hover:text-primary-300"
        class:text-primary-400={active === "/settings/permissions/"}
        class:dark:text-primary-300={active === "/settings/permissions/"}
        >Permissions</a
    >
    <a
        href="/settings/providers"
        class="btn variant-outline-surface hover:text-primary-300"
        class:text-primary-400={active === "/settings/providers/"}
        class:dark:text-primary-300={active === "/settings/providers/"}
        >Providers</a
    >
    <a
        href="/settings/audit"
        class="btn variant-outline-surface hover:text-primary-300"
        class:text-primary-400={active === "/settings/audit/"}
        class:dark:text-primary-300={active === "/settings/audit/"}>Logs</a
    >
</div>
