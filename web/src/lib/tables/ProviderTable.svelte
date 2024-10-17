<script lang="ts">
    import { pb } from "$lib/client";
    import { filterCollection, providers } from "$lib/subscriptions";
    import {
        Paginator,
        type PaginationSettings,
        getModalStore,
        type ModalSettings,
    } from "@skeletonlabs/skeleton";
    import type { RecordModel } from "pocketbase";
    import { onMount } from "svelte";

    // Pagination settings
    let total = 0;
    let searchTerm = "";
    let sortColumn = "";
    let sortDirection = "asc";
    let paginationSettings = {
        page: 0,
        limit: pb.authStore.model?.settings?.table_limit || 30,
        size: 0,
        amounts: [30, 50, 100],
    } satisfies PaginationSettings;
    $: searchTerm, sortColumn, sortDirection, update();

    const update = async () => {
        const direction = sortDirection === "asc" ? "+" : "-";
        const sort = sortColumn ? `${direction}${sortColumn}` : "";
        let filter = "";
        if (searchTerm) {
            filter = `(name ~ "${searchTerm}" || type ~ "${searchTerm}")`;
        }
        total = await filterCollection(
            "providers",
            paginationSettings.page + 1,
            paginationSettings.limit,
            filter,
            sort,
        );

        if (!$providers.length && paginationSettings.page > 0) {
            paginationSettings.page -= 1;
            update();
        }
    };

    const setLimit = async () => {
        if (pb.authStore.isAdmin) return;

        await pb.collection("users").update(pb.authStore.model?.id, {
            settings: {
                ...pb.authStore.model?.settings,
                table_limit: paginationSettings.limit,
            },
        });
    };

    const sortTable = (column: string) => {
        sortDirection =
            sortColumn === column && sortDirection === "asc" ? "dsc" : "asc";
        sortColumn = column;
    };

    const modals = getModalStore();
    const edit = (data: RecordModel) => {
        const modal: ModalSettings = {
            type: "component",
            title: "Edit Provider",
            component: "ProviderUpdate",
            meta: { provider: data },
        };
        modals.trigger(modal);
    };
    const create = () => {
        const modal: ModalSettings = {
            type: "component",
            title: "Add Provider",
            component: "ProviderCreate",
        };
        modals.trigger(modal);
    };

    onMount(async () => {
        total = await filterCollection(
            "providers",
            paginationSettings.page + 1,
            paginationSettings.limit,
        );
    });
</script>

<div class="flex flex-col gap-4 w-full">
    <div
        class="input-group input-group-divider grid-cols-[auto_1fr_auto] md:max-w-xs"
    >
        <div class="input-group-shim">
            <iconify-icon icon="fa6-solid:magnifying-glass" />
        </div>
        <input type="search" placeholder="Filter..." bind:value={searchTerm} />
    </div>
    <div class="table-container">
        <table class="table table-interactive">
            <thead>
                <tr class="cursor-pointer">
                    <th
                        on:click={() => sortTable("name")}
                        class:table-sort-asc={sortDirection === "asc" &&
                            sortColumn === "name"}
                        class:table-sort-dsc={sortDirection === "dsc" &&
                            sortColumn === "name"}>Name</th
                    >
                    <th
                        on:click={() => sortTable("type")}
                        class="hidden md:table-cell"
                        class:table-sort-asc={sortDirection === "asc" &&
                            sortColumn === "type"}
                        class:table-sort-dsc={sortDirection === "dsc" &&
                            sortColumn === "type"}>Type</th
                    >
                    <th
                        on:click={() => sortTable("error")}
                        class:table-sort-asc={sortDirection === "asc" &&
                            sortColumn === "error"}
                        class:table-sort-dsc={sortDirection === "dsc" &&
                            sortColumn === "error"}>Error</th
                    >
                    <th
                        on:click={() => sortTable("last_sync")}
                        class:table-sort-asc={sortDirection === "asc" &&
                            sortColumn === "last_sync"}
                        class:table-sort-dsc={sortDirection === "dsc" &&
                            sortColumn === "last_sync"}>Last Sync</th
                    >
                </tr>
            </thead>
            <tbody>
                {#each $providers as row}
                    <tr>
                        <td on:click={() => edit(row)} aria-hidden>
                            {row.name}
                        </td>
                        <td
                            on:click={() => edit(row)}
                            class="hidden md:table-cell"
                            aria-hidden
                        >
                            {row.type}
                        </td>
                        <td
                            on:click={() => edit(row)}
                            class="hidden md:table-cell"
                            aria-hidden
                        >
                            {row.error}
                        </td>
                        <td
                            on:click={() => edit(row)}
                            class="hidden md:table-cell"
                            aria-hidden
                        >
                            {#if row.last_sync}
                                {Intl.DateTimeFormat("en", {
                                    dateStyle: "full",
                                    timeStyle: "short",
                                }).format(new Date(row.last_sync))}
                            {:else}
                                Never
                            {/if}
                        </td>
                    </tr>
                {/each}
            </tbody>
            <tfoot>
                <tr>
                    <th colspan="2">Total</th>
                    <td>
                        <code class="code"
                            >{paginationSettings.size} providers</code
                        >
                    </td>
                </tr>
            </tfoot>
        </table>
    </div>

    <Paginator
        bind:settings={paginationSettings}
        on:page={update}
        on:amount={setLimit}
        showFirstLastButtons={false}
        showPreviousNextButtons={true}
    />

    <button class="flex-auto btn variant-filled-success mt-4" on:click={create}
        >Add External Provider</button
    >
</div>
