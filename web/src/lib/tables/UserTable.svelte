<script lang="ts">
    import { pb } from "$lib/client";
    import {
        Paginator,
        type PaginationSettings,
        getModalStore,
        type ModalSettings,
    } from "@skeletonlabs/skeleton";
    import GroupSidebar from "$lib/navigation/GroupSB.svelte";
    import type { RecordModel } from "pocketbase";
    import { onMount } from "svelte";
    import { filterCollection, users } from "$lib/subscriptions";

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
        const filters: string[] = [];
        if (sortGroup) {
            filters.push(`groups.id ?= "${sortGroup}"`);
        }
        if (searchTerm) {
            filters.push(
                `(username ~ "${searchTerm}" || name ~ "${searchTerm}")`,
            );
        }
        const filter = filters.join(" && ");
        total = await filterCollection(
            "users",
            paginationSettings.page + 1,
            paginationSettings.limit,
            filter,
            sort,
        );
        if (!$users.length && paginationSettings.page > 0) {
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

    const selectToggle = () => {
        const allChecked = $users.every((user) => user.selected);
        $users = $users.map((user) => {
            user.selected = !allChecked;
            return user;
        });
    };

    const modals = getModalStore();
    const editUser = (user: RecordModel) => {
        const modal: ModalSettings = {
            type: "component",
            title: "Edit User",
            component: "UserUpdate",
            meta: { user: user },
        };
        modals.trigger(modal);
    };

    const removeGroup = async (user: RecordModel, group: RecordModel) => {
        user.selected = true;
        let records = $users.filter((user) => user.selected);

        for (const user of records) {
            try {
                await pb.collection("users").update(user.id, {
                    "groups-": group.id,
                });
            } catch (error) {}
            user.selected = false;
        }
    };

    // Drag and drop users into groups
    let draggedRecord: RecordModel[] | undefined;
    const handleDragStart = (record: RecordModel) => {
        record.selected = true;
        draggedRecord = $users.filter((item) => item.selected);
    };

    onMount(async () => {
        total = await filterCollection(
            "users",
            paginationSettings.page + 1,
            paginationSettings.limit,
        );
    });
</script>

<div class="flex flex-col md:flex-row gap-4">
    <GroupSidebar bind:group={sortGroup} bind:records={draggedRecord} />

    <div class="flex flex-col gap-4 w-full">
        <div
            class="input-group input-group-divider grid-cols-[auto_1fr_auto] md:max-w-xs"
        >
            <div class="input-group-shim">
                <iconify-icon icon="fa6-solid:magnifying-glass" />
            </div>
            <input
                type="search"
                placeholder="Filter..."
                bind:value={searchTerm}
            />
        </div>
        <div class="table-container">
            <table class="table table-compact table-interactive">
                <thead>
                    <tr class="cursor-pointer">
                        <th on:click={selectToggle} class="table-cell-fit">*</th
                        >
                        <th
                            on:click={() => sortTable("username")}
                            class:table-sort-asc={sortDirection === "asc" &&
                                sortColumn === "username"}
                            class:table-sort-dsc={sortDirection === "dsc" &&
                                sortColumn === "username"}>Username</th
                        >
                        <th
                            on:click={() => sortTable("email")}
                            class="hidden md:table-cell"
                            class:table-sort-asc={sortDirection === "asc" &&
                                sortColumn === "email"}
                            class:table-sort-dsc={sortDirection === "dsc" &&
                                sortColumn === "email"}>Email</th
                        >
                        <th
                            on:click={() => sortTable("permission")}
                            class="hidden md:table-cell"
                            class:table-sort-asc={sortDirection === "asc" &&
                                sortColumn === "permission"}
                            class:table-sort-dsc={sortDirection === "dsc" &&
                                sortColumn === "permission"}>Permission</th
                        >
                        <th
                            on:click={() => sortTable("principal")}
                            class="hidden lg:table-cell"
                            class:table-sort-asc={sortDirection === "asc" &&
                                sortColumn === "principal"}
                            class:table-sort-dsc={sortDirection === "dsc" &&
                                sortColumn === "principal"}>Principal</th
                        >
                        <th
                            on:click={() => sortTable("groups")}
                            class:table-sort-asc={sortDirection === "asc" &&
                                sortColumn === "groups"}
                            class:table-sort-dsc={sortDirection === "dsc" &&
                                sortColumn === "groups"}>Groups</th
                        >
                    </tr>
                </thead>
                <tbody>
                    {#each $users as row}
                        <tr
                            draggable="true"
                            on:dragstart={() => handleDragStart(row)}
                            class:table-row-checked={row.selected}
                        >
                            <td>
                                <label
                                    class="flex items-center space-x-2 table-cell-fit"
                                >
                                    <input
                                        class="checkbox"
                                        type="checkbox"
                                        bind:checked={row.selected}
                                        on:click|stopPropagation
                                    />
                                </label>
                            </td>
                            <td on:click={() => editUser(row)} aria-hidden>
                                {#if row.username}
                                    {row.username}
                                {/if}
                            </td>
                            <td
                                on:click={() => editUser(row)}
                                class="hidden md:table-cell"
                                aria-hidden
                            >
                                {#if row.email}
                                    {row.email}
                                {/if}
                            </td>
                            <td
                                on:click={() => editUser(row)}
                                class="hidden md:table-cell"
                                aria-hidden
                            >
                                <div
                                    class="text-sm text-center px-2 py-1 rounded-xl bg-surface-300 dark:bg-surface-500"
                                >
                                    {#if row.permission}
                                        {row.expand?.permission?.name}
                                    {:else}
                                        None
                                    {/if}
                                </div>
                            </td>
                            <td
                                on:click={() => editUser(row)}
                                class="hidden lg:table-cell"
                                aria-hidden
                            >
                                {row.principal}
                            </td>
                            <td aria-hidden>
                                {#each row.expand?.groups || [] as group}
                                    <span
                                        class="badge variant-soft-primary mr-0.5 mb-0.5 hover:variant-ringed-primary hover:line-through"
                                        on:click={() => removeGroup(row, group)}
                                        aria-hidden
                                    >
                                        {group.name}
                                    </span>
                                {/each}
                            </td>
                        </tr>
                    {/each}
                </tbody>
                <tfoot>
                    <tr>
                        <th colspan="2">Total</th>
                        <td>
                            <code class="code"
                                >{paginationSettings.size} users</code
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
</div>
