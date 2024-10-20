<script lang="ts">
    import { pb } from "$lib/client";
    import { onMount } from "svelte";

    let version = "";

    onMount(async () => {
        await pb.send("/api/version", { method: "GET" }).then((res) => {
            if (res.version !== "unknown") {
                version = res.version ? res.version : "";
            }
        });
    });
</script>

<footer
    class="bottom-0 right-0 left-0 flex flex-row items-center justify-between gap-2 bg-background px-4 py-2"
>
    <div class="flex flex-row items-center gap-2 text-xs text-gray-700/50">
        {#if version && version !== "unknown"}
            <a
                href="https://github.com/MizuchiLabs/ssh-nexus/releases"
                class="flex flex-row items-center gap-1 hover:text-primary-300"
                target="_blank"
            >
                SSH Nexus {version}
            </a>
        {/if}
    </div>

    <div class="flex flex-row items-center gap-2">
        <!-- <a href="/" target="_blank" rel="noreferrer"> -->
        <!-- 	<iconify-icon icon="solar:book-bold-duotone" /> -->
        <!-- </a> -->
        <a
            href="https://github.com/mizuchilabs/ssh-nexus"
            target="_blank"
            rel="noreferrer"
        >
            <iconify-icon icon="line-md:github-loop" />
        </a>
    </div>
</footer>
