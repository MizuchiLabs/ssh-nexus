<script lang="ts">
    import { pb } from "$lib/client";
    import { filterCollection, auditlog } from "$lib/subscriptions";
    import { Paginator, type PaginationSettings } from "@skeletonlabs/skeleton";
    import type { RecordModel } from "pocketbase";
    import { onMount } from "svelte";

    // Pagination settings
    let total = 0;
    let sortGroup = "";
    let searchTerm = "";
    let sortColumn = "";
    let sortDirection = "asc";
    let paginationSettings = {
        page: 0,
        limit: pb.authStore.model?.settings?.table_limit || 30,
        size: 0,
        amounts: [30, 50, 100],
    } satisfies PaginationSettings;
    $: sortGroup, searchTerm, sortColumn, sortDirection, update();

    const update = async () => {
        const direction = sortDirection === "asc" ? "+" : "-";
        const sort = sortColumn ? `${direction}${sortColumn}` : "";
        let filter = "";
        if (searchTerm) {
            filter = `(collection ~ "${searchTerm}" || data ~ "${searchTerm}")`;
        }
        total = await filterCollection(
            "auditlog",
            paginationSettings.page + 1,
            paginationSettings.limit,
            filter,
            sort,
        );

        if (!$auditlog.length && paginationSettings.page > 0) {
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

    // Nice messages
    function getCollectionName(e: RecordModel) {
        switch (e.collection) {
            case "users":
            case "machines":
            case "groups":
            case "permissions":
                return e.data.name;
            case "settings":
                return e.data.key;
            default:
                return "Unknown"; // Handle unknown collections
        }
    }

    function getEventMessage(e: RecordModel) {
        if (!e.collection || !e.event) return null;

        let message = `${e.event}d`; // Base message with "d" suffix

        switch (e.collection) {
            case "users":
            case "machines":
            case "settings":
                return message; // No additional message for these collections
            default:
                return message;
        }
    }

    function getTriggeredBy(e: RecordModel) {
        return e.expand?.user?.name || "Admin";
    }

    onMount(async () => {
        total = await filterCollection(
            "auditlog",
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
        <table class="table table-hover table-comfortable">
            <thead>
                <tr class="cursor-pointer">
                    <th
                        on:click={() => sortTable("collection")}
                        class:table-sort-asc={sortDirection === "asc" &&
                            sortColumn === "collection"}
                        class:table-sort-dsc={sortDirection === "dsc" &&
                            sortColumn === "collection"}>Collection</th
                    >
                    <th
                        on:click={() => sortTable("change")}
                        class:table-sort-asc={sortDirection === "asc" &&
                            sortColumn === "change"}
                        class:table-sort-dsc={sortDirection === "dsc" &&
                            sortColumn === "change"}>Change</th
                    >
                    <th
                        on:click={() => sortTable("date")}
                        class="hidden md:table-cell"
                        class:table-sort-asc={sortDirection === "asc" &&
                            sortColumn === "date"}
                        class:table-sort-dsc={sortDirection === "dsc" &&
                            sortColumn === "date"}>Date</th
                    >
                </tr>
            </thead>
            <tbody>
                {#each $auditlog as row}
                    <tr>
                        <td>
                            {row.collection}
                        </td>
                        <td>
                            <span class="badge variant-filled-surface"
                                >{getCollectionName(row)}</span
                            >

                            <span
                                class="badge"
                                class:variant-filled-success={row.event ===
                                    "create"}
                                class:variant-filled-error={row.event ===
                                    "delete"}
                                class:variant-filled-warning={row.event ===
                                    "update"}
                            >
                                {getEventMessage(row)}
                            </span>
                            by
                            <span class="badge variant-filled-surface">
                                {getTriggeredBy(row)}</span
                            >
                        </td>
                        <td class="hidden md:table-cell">
                            {new Date(row.created).toUTCString()}
                        </td>
                    </tr>
                {/each}
            </tbody>
        </table>
    </div>

    <Paginator
        bind:settings={paginationSettings}
        on:page={update}
        on:amount={setLimit}
        showFirstLastButtons={false}
        showPreviousNextButtons={true}
    />
</div>
