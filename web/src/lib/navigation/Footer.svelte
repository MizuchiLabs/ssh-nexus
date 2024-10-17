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

<div
    class="flex flex-row items-center justify-start ml-2 mb-1 gap-2 text-sm font-light"
>
    <!-- <a -->
    <!--     href="https://github.com/MizuchiLabs/ssh-nexus" -->
    <!--     class="btn-icon btn-icon-sm variant-soft-surface hover:text-primary-300" -->
    <!-- > -->
    <!--     <iconify-icon icon="mingcute:github-line" width="12" /> -->
    <!-- </a> -->
    {#if version}
        <a
            href="https://github.com/MizuchiLabs/ssh-nexus/releases"
            class="flex flex-row items-center gap-1 hover:text-primary-300"
            target="_blank"
        >
            {version}
        </a>
    {/if}
</div>
