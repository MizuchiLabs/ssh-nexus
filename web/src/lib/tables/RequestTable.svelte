<script lang="ts">
    import { pb } from "$lib/client";
    import { filterCollection, requests } from "$lib/subscriptions";
    import { Paginator, type PaginationSettings } from "@skeletonlabs/skeleton";
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
        const sort = sortColumn ? `${direction}${sortColumn}` : "+created";
        const filters: string[] = [];

        if (searchTerm) {
            filters.push(
                `(user.name ~ "${searchTerm}" || machine.name ~ "${searchTerm}" || 
                group.name ~ "${searchTerm}")`,
            );
        }
        const filter = filters.join(" && ");
        total = await filterCollection(
            "requests",
            paginationSettings.page + 1,
            paginationSettings.limit,
            filter,
            sort,
        );

        if (!$requests.length && paginationSettings.page > 0) {
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
    const accept = async (request: RecordModel) => {
        if (request.expand?.machine) {
            await pb.collection("machines").update(request.expand.machine.id, {
                "users+": request.expand.user.id,
            });
        }
        if (request.expand?.group) {
            await pb.collection("users").update(request.expand.user.id, {
                "groups+": request.expand.group.id,
            });
        }
        await pb.collection("requests").delete(request.id);
    };

    const deny = async (request: RecordModel) => {
        await pb.collection("requests").delete(request.id);
    };

    onMount(async () => {
        total = await filterCollection(
            "requests",
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
        <table class="table table-compact table-interactive">
            <thead>
                <tr class="cursor-pointer">
                    <th
                        on:click={() => sortTable("user")}
                        class:table-sort-asc={sortDirection === "asc" &&
                            sortColumn === "user"}
                        class:table-sort-dsc={sortDirection === "dsc" &&
                            sortColumn === "user"}>User</th
                    >
                    <th
                        on:click={() => sortTable("machine")}
                        class:table-sort-asc={sortDirection === "asc" &&
                            sortColumn === "machine"}
                        class:table-sort-dsc={sortDirection === "dsc" &&
                            sortColumn === "machine"}>Group/Machine</th
                    >
                    <th
                        on:click={() => sortTable("created")}
                        class="hidden lg:table-cell"
                        class:table-sort-asc={sortDirection === "asc" &&
                            sortColumn === "created"}
                        class:table-sort-dsc={sortDirection === "dsc" &&
                            sortColumn === "created"}>Request Time</th
                    >
                    <th>Actions</th>
                </tr>
            </thead>
            <tbody>
                {#each $requests as row}
                    <tr>
                        <td class="font-bold">
                            <span class="chip variant-soft-surface">
                                {row.expand?.user.username}
                            </span>
                        </td>
                        <td class="font-bold">
                            {#if row.expand?.machine}
                                <span class="chip variant-soft-secondary">
                                    <iconify-icon
                                        icon="fa6-solid:server"
                                        class="mr-2"
                                    />
                                    {row.expand?.machine.name}
                                </span>
                            {:else}
                                <span class="chip variant-soft-secondary">
                                    <iconify-icon
                                        icon="fa6-solid:users"
                                        class="mr-2"
                                    />
                                    {row.expand?.group.name}
                                </span>
                            {/if}
                        </td>
                        <td class="hidden md:table-cell" aria-hidden>
                            {Intl.DateTimeFormat("en", {
                                dateStyle: "full",
                                timeStyle: "short",
                            }).format(new Date(row.created))}
                        </td>
                        <td>
                            <button
                                class="btn btn-sm variant-filled-success"
                                on:click={() => accept(row)}>Accept</button
                            >
                            <button
                                class="btn btn-sm variant-filled-error"
                                on:click={() => deny(row)}>Deny</button
                            >
                        </td>
                    </tr>
                {/each}
            </tbody>
            <tfoot>
                <tr>
                    <th colspan="2">Total</th>
                    <td>
                        <code class="code"
                            >{paginationSettings.size} requests</code
                        ></td
                    >
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
</div>
