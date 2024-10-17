<script lang="ts">
    import { pb } from "$lib/client";
    import { filterCollection, permissions } from "$lib/subscriptions";
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
    let sortColumn = "name";
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
            filter = `(name ~ "${searchTerm}")`;
        }
        total = await filterCollection(
            "permissions",
            paginationSettings.page + 1,
            paginationSettings.limit,
            filter,
            sort,
        );

        if (!$permissions.length && paginationSettings.page > 0) {
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
    const editPermission = (permission: RecordModel) => {
        const modal: ModalSettings = {
            type: "component",
            title: "Edit Permission",
            component: "PermissionUpdate",
            meta: { permission: permission },
        };
        modals.trigger(modal);
    };
    const createPermission = () => {
        const modal: ModalSettings = {
            type: "component",
            title: "New Permission",
            component: "PermissionCreate",
        };
        modals.trigger(modal);
    };

    onMount(async () => {
        total = await filterCollection(
            "permissions",
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
        <table class="table table-interactive table-fixed">
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
                        on:click={() => sortTable("is_admin")}
                        class:table-sort-asc={sortDirection === "asc" &&
                            sortColumn === "is_admin"}
                        class:table-sort-dsc={sortDirection === "dsc" &&
                            sortColumn === "is_admin"}>Admin</th
                    >
                    <th class="hidden md:table-cell">Access</th>
                    <th class="hidden md:table-cell">Type</th>
                </tr>
            </thead>
            <tbody>
                {#each $permissions as row}
                    <tr>
                        <td on:click={() => editPermission(row)} aria-hidden>
                            {row.name}
                        </td>
                        <td on:click={() => editPermission(row)} aria-hidden>
                            {#if row.is_admin}
                                <iconify-icon icon="fa6-solid:check" />
                            {:else}
                                <iconify-icon icon="fa6-solid:xmark" />
                            {/if}
                        </td>
                        <td class="hidden md:table-cell">
                            <div class="flex flex-col gap-1">
                                {#if row.access_users}
                                    <span class="badge variant-filled-primary"
                                        >Users
                                    </span>
                                {/if}
                                {#if row.access_machines}
                                    <span class="badge variant-filled-primary"
                                        >Machines
                                    </span>
                                {/if}
                                {#if row.access_groups}
                                    <span class="badge variant-filled-primary"
                                        >Groups
                                    </span>
                                {/if}
                            </div>
                        </td>
                        <td class="hidden md:table-cell">
                            <div class="flex flex-col gap-1">
                                {#if row.can_create}
                                    <span class="badge variant-filled-success"
                                        >Create
                                    </span>
                                {/if}
                                {#if row.can_update}
                                    <span class="badge variant-filled-warning"
                                        >Update
                                    </span>
                                {/if}
                                {#if row.can_delete}
                                    <span class="badge variant-filled-error"
                                        >Delete
                                    </span>
                                {/if}
                            </div>
                        </td>
                    </tr>
                {/each}
            </tbody>
            <tfoot>
                <tr>
                    <th colspan="2">Total</th>
                    <td>
                        <code class="code">
                            {paginationSettings.size} permissions
                        </code>
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
    <button
        class="flex-auto btn variant-filled-success mt-4"
        on:click={createPermission}>Add permission rule</button
    >
</div>
