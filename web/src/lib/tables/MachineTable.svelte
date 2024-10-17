<script lang="ts">
    import { pb } from "$lib/client";
    import { filterCollection, machines } from "$lib/subscriptions";
    import {
        Paginator,
        type PaginationSettings,
        getModalStore,
        type ModalSettings,
    } from "@skeletonlabs/skeleton";
    import GroupSidebar from "$lib/navigation/GroupSB.svelte";
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
    $: paginationSettings.size = total;
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
                `(name ~ "${searchTerm}" || host ~ "${searchTerm}" || 
                groups.name ~ "${searchTerm}" || tags.name ~ "${searchTerm}")`,
            );
        }
        const filter = filters.join(" && ");
        total = await filterCollection(
            "machines",
            paginationSettings.page + 1,
            paginationSettings.limit,
            filter,
            sort,
        );

        if (!$machines.length && paginationSettings.page > 1) {
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

    // Selection
    const selectToggle = () => {
        const allChecked = $machines.every((machine) => machine.selected);
        $machines = $machines.map((machine) => {
            machine.selected = !allChecked;
            return machine;
        });
    };

    // Modals
    const modals = getModalStore();
    const editMachine = (machine: RecordModel) => {
        const modal: ModalSettings = {
            type: "component",
            title: "Edit Machine",
            component: "MachineUpdate",
            meta: { machine: machine },
        };
        modals.trigger(modal);
    };

    const removeGroup = async (machine: RecordModel, group: RecordModel) => {
        machine.selected = true;
        let records = $machines.filter((machine) => machine.selected);

        for (const machine of records) {
            try {
                await pb.collection("machines").update(machine.id, {
                    "groups-": group.id,
                });
            } catch (error) {}
            machine.selected = false;
        }
    };

    const removeTag = async (machine: RecordModel, tag: RecordModel) => {
        machine.selected = true;
        let records = $machines.filter((machine) => machine.selected);

        for (const machine of records) {
            try {
                await pb.collection("machines").update(machine.id, {
                    "tags-": tag.id,
                });
            } catch (error) {}
            machine.selected = false;
        }
    };

    // Drag and drop
    let draggedRecord: RecordModel[] | undefined;
    const handleDragStart = (record: RecordModel) => {
        record.selected = true;
        draggedRecord = $machines.filter((item) => item.selected);
    };

    onMount(async () => {
        total = await filterCollection(
            "machines",
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
                            on:click={() => sortTable("name")}
                            class:table-sort-asc={sortDirection === "asc" &&
                                sortColumn === "name"}
                            class:table-sort-dsc={sortDirection === "dsc" &&
                                sortColumn === "name"}>Name</th
                        >
                        <th
                            on:click={() => sortTable("host")}
                            class="hidden lg:table-cell"
                            class:table-sort-asc={sortDirection === "asc" &&
                                sortColumn === "host"}
                            class:table-sort-dsc={sortDirection === "dsc" &&
                                sortColumn === "host"}>Host</th
                        >
                        <th
                            on:click={() => sortTable("groups")}
                            class:table-sort-asc={sortDirection === "asc" &&
                                sortColumn === "groups"}
                            class:table-sort-dsc={sortDirection === "dsc" &&
                                sortColumn === "groups"}>Groups</th
                        >
                        <th
                            on:click={() => sortTable("tags")}
                            class="hidden md:table-cell"
                            class:table-sort-asc={sortDirection === "asc" &&
                                sortColumn === "tags"}
                            class:table-sort-dsc={sortDirection === "dsc" &&
                                sortColumn === "tags"}>Tags</th
                        >
                        <th
                            on:click={() => sortTable("error")}
                            class="hidden lg:table-cell"
                            class:table-sort-asc={sortDirection === "asc" &&
                                sortColumn === "error"}
                            class:table-sort-dsc={sortDirection === "dsc" &&
                                sortColumn === "error"}>Errors</th
                        >
                        <th
                            on:click={() => sortTable("agent")}
                            class="hidden sm:table-cell"
                            class:table-sort-asc={sortDirection === "asc" &&
                                sortColumn === "agent"}
                            class:table-sort-dsc={sortDirection === "dsc" &&
                                sortColumn === "agent"}>Agent</th
                        >
                    </tr>
                </thead>
                <tbody>
                    {#each $machines as row}
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
                            <td on:click={() => editMachine(row)} aria-hidden>
                                {#if row.provider !== ""}
                                    <span class="badge variant-soft-primary">
                                        <iconify-icon
                                            icon="fa6-solid:cloud"
                                            class="mr-1"
                                        />
                                        {row.name}
                                    </span>
                                {:else}
                                    <span class="badge variant-soft-secondary">
                                        {row.name}
                                    </span>
                                {/if}
                            </td>
                            <td
                                on:click={() => editMachine(row)}
                                class="hidden lg:table-cell"
                                aria-hidden
                            >
                                {row.host}
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
                            <td class="hidden md:table-cell" aria-hidden>
                                {#each row.expand?.tags || [] as tag}
                                    <span
                                        class="badge mr-0.5 mb-0.5 bg-surface-300-600-token hover:bg-transparent hover:line-through"
                                        on:click={() => removeTag(row, tag)}
                                        aria-hidden
                                    >
                                        {tag.name}
                                    </span>
                                {/each}
                            </td>
                            <td
                                on:click={() => editMachine(row)}
                                class="hidden lg:table-cell"
                                aria-hidden
                            >
                                <div class="flex flex-row gap-1 items-center">
                                    {#if row.error}
                                        <span class="badge variant-filled-error"
                                            >{row.error}</span
                                        >
                                    {:else}
                                        <iconify-icon
                                            icon="fa6-solid:check"
                                            class="text-green-500"
                                        />
                                    {/if}
                                </div>
                            </td>
                            <td
                                on:click={() => editMachine(row)}
                                class="hidden sm:table-cell"
                                aria-hidden
                            >
                                <div class="flex flex-row gap-1 items-center">
                                    {#if row.agent}
                                        <iconify-icon
                                            icon="fa6-solid:check"
                                            class="text-green-500"
                                        />
                                    {:else}
                                        <iconify-icon
                                            icon="fa6-solid:xmark"
                                            class="text-red-500"
                                        />
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
                            <code class="code"
                                >{paginationSettings.size} machines</code
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
