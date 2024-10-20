<script lang="ts">
    import { page } from "$app/stores";
    import autoAnimate from "@formkit/auto-animate";

    $: active = $page.url.pathname;
</script>

<main class="flex-1 flex-col gap-4 p-4 md:gap-8 md:p-10">
    <div
        class="grid grid-cols-4 items-start gap-4 md:grid-cols-[180px_1fr] lg:grid-cols-[250px_1fr]"
    >
        <nav class="grid-cols-1">
            <div class="mx-auto grid w-full max-w-6xl gap-2 mb-8">
                <h1 class="text-3xl font-semibold">Settings</h1>
            </div>
            <div class="text-muted-foreground grid gap-4 text-sm">
                <a
                    href="/settings/general"
                    class="hover:text-red-400 hover:dark:text-red-300"
                    class:font-semibold={active === "/settings/general/"}
                    class:text-primary={active === "/settings/general/"}
                >
                    General
                </a>
                <a
                    href="/settings/ssh"
                    class="hover:text-red-400 hover:dark:text-red-300"
                    class:font-semibold={active === "/settings/ssh/"}
                    class:text-primary={active === "/settings/ssh/"}
                >
                    Keys
                </a>
                <a
                    href="/settings/permissions"
                    class="hover:text-red-400 hover:dark:text-red-300"
                    class:font-semibold={active === "/settings/permissions/"}
                    class:text-primary={active === "/settings/permissions/"}
                >
                    Permissions
                </a>
                <a
                    href="/settings/providers"
                    class="hover:text-red-400 hover:dark:text-red-300"
                    class:font-semibold={active === "/settings/providers/"}
                    class:text-primary={active === "/settings/providers/"}
                >
                    Providers
                </a>
                <a
                    href="/settings/auditlog"
                    class="hover:text-red-400 hover:dark:text-red-300"
                    class:font-semibold={active === "/settings/auditlog/"}
                    class:text-primary={active === "/settings/auditlog/"}
                >
                    Auditlog
                </a>
            </div>
        </nav>

        <div class="grid-cols-3 gap-2" use:autoAnimate={{ duration: 100 }}>
            <slot />
        </div>
    </div>
</main>
