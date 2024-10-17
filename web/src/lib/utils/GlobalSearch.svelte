<script lang="ts">
    import { pb } from "$lib/client";
    import { popup } from "@skeletonlabs/skeleton";
    import type { RecordModel } from "pocketbase";
    import { onDestroy, onMount } from "svelte";

    let input = "";
    let searchInput: HTMLInputElement;
    let currentFocusIndex = -1;
    let results: Record<string, RecordModel[]> = {};
    $: input, search();

    const collectionFields: Record<string, string[]> = {
        users: ["id", "username", "name", "email", "principal"],
        groups: ["id", "name"],
        machines: ["id", "name", "host", "uuid", "error"],
        permissions: ["id", "name"],
        providers: ["id", "name", "type"],
        tags: ["id", "name"],
    };

    const search = async () => {
        results = {};
        if (!input) return;
        try {
            for (const collection of Object.keys(collectionFields)) {
                const filterString = collectionFields[collection]
                    .map((field) => `${field} ~ "${input}"`)
                    .join(" || ");
                const records = await pb.collection(collection).getList(1, 30, {
                    filter: filterString,
                });

                if (records.items.length > 0) {
                    results[collection] = records.items;
                }
            }
        } catch (error) {}
    };

    const handleKeyDown = (event: KeyboardEvent) => {
        const totalResults = Object.values(results).flat().length;

        if (event.key === "Tab") {
            event.preventDefault(); // Prevent the default tab behavior
            currentFocusIndex = (currentFocusIndex + 1) % (totalResults + 1); // Loop through the results

            if (currentFocusIndex === 0) {
                searchInput.focus(); // Focus back on the input
            } else {
                const resultElements =
                    document.querySelectorAll(".result-item");
                if (resultElements[currentFocusIndex - 1]) {
                    (
                        resultElements[currentFocusIndex - 1] as HTMLElement
                    ).focus();
                }
            }
        }
    };

    onMount(() => {
        window.addEventListener("keydown", handleKeyDown);
    });

    onDestroy(() => {
        window.removeEventListener("keydown", handleKeyDown);
    });
</script>

<div
    class="flex items-center justify-between w-full max-w-xl mx-auto bg-white dark:bg-surface-500 rounded-full"
>
    <input
        class="rounded-full w-full h-12 px-4 font-mono dark:bg-surface-500 outline-none hover:outline-none border-none focus:outline-none focus:ring-0 focus:border-none focus:shadow-none"
        type="text"
        name="search"
        id="search"
        placeholder="Search..."
        bind:value={input}
        bind:this={searchInput}
        on:focus={search}
        on:blur={() => (results = {})}
        use:popup={{
            event: "focus-blur",
            target: "popupResults",
            placement: "bottom",
        }}
    />
    <button
        type="submit"
        class="btn bg-primary-400 mr-2 text-white"
        on:click={search}
        tabindex="-1"
    >
        <iconify-icon icon="fa6-solid:magnifying-glass" />
    </button>
</div>

<div
    class="z-10 mt-2 bg-surface-100 dark:bg-surface-700 rounded-lg"
    data-popup="popupResults"
    class:hidden={!Object.keys(results).length}
>
    {#each Object.entries(results) as [collection, records]}
        <h3 class="px-2 py-1 capitalize bg-primary-400 rounded-lg">
            {collection}
        </h3>
        {#each records as record}
            <a
                href="/{collection}/{record.id}"
                class="block p-2 hover:bg-surface-200 dark:hover:bg-surface-600 focus:bg-surface-100 rounded-lg result-item"
            >
                {#if collection === "users"}
                    <p>{record.username} - {record.email}</p>
                {:else if collection === "machines"}
                    <p>{record.name} - {record.host}</p>
                {:else}
                    <p>{record.name}</p>
                {/if}
            </a>
        {/each}
    {/each}
</div>
