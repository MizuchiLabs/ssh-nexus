<script lang="ts">
    import {
        users,
        machines,
        permissions,
        providers,
        groups,
        requests,
        auditlog,
    } from "$lib/subscriptions";
    import {
        createTable,
        createRender,
        Render,
        Subscribe,
    } from "svelte-headless-table";
    import {
        addPagination,
        addSortBy,
        addTableFilter,
        addColumnFilters,
        addHiddenColumns,
        addSelectedRows,
    } from "svelte-headless-table/plugins";
    import * as Table from "$lib/components/ui/table";
    import * as Select from "$lib/components/ui/select";
    import * as DropdownMenu from "$lib/components/ui/dropdown-menu";
    import { Button } from "$lib/components/ui/button";
    import { Input } from "$lib/components/ui/input";
    import { ArrowDown, ArrowUp, ChevronDown } from "lucide-svelte";
    import TablePermissions from "$lib/tables/permissions.svelte";
    import TableCheckbox from "$lib/tables/checkbox.svelte";
    import TableActions from "$lib/tables/actions.svelte";
    import TableRequestActions from "$lib/tables/reqActions.svelte";
    import TableGroups from "$lib/tables/groups.svelte";
    import TableTags from "$lib/tables/tags.svelte";
    import TableBool from "$lib/tables/bool.svelte";
    import TableAccess from "$lib/tables/access.svelte";
    import TableLog from "$lib/tables/log.svelte";
    import { onMount } from "svelte";
    import { writable } from "svelte/store";
    import type { RecordModel } from "pocketbase";

    export let collection: string;

    const plugins = {
        page: addPagination(),
        hide: addHiddenColumns(),
        select: addSelectedRows(),
        sort: addSortBy({ disableMultiSort: true }),
        filter: addTableFilter({
            includeHiddenColumns: true,
            fn: ({ filterValue, value }): boolean =>
                value.toLowerCase().includes(filterValue.toLowerCase()),
        }),
        colFilter: addColumnFilters(),
    };
    let table = createTable(writable([] as RecordModel[]), plugins);

    const selectionColumn = table.column({
        accessor: "id",
        header: (_, { pluginStates }) => {
            const { allPageRowsSelected } = pluginStates.select;
            return createRender(TableCheckbox, {
                checked: allPageRowsSelected,
            });
        },
        cell: ({ row }, { pluginStates }) => {
            const { getRowState } = pluginStates.select;
            const { isSelected } = getRowState(row);

            return createRender(TableCheckbox, {
                checked: isSelected,
            });
        },
        plugins: {
            sort: {
                disable: true,
            },
            filter: {
                exclude: true,
            },
        },
    });
    const actionsColumn = table.column({
        accessor: (item) => item,
        header: "",
        cell: ({ value }) => {
            return createRender(TableActions, {
                record: value,
                collection: collection,
            });
        },
        plugins: {
            sort: {
                disable: true,
            },
            filter: {
                exclude: true,
            },
        },
    });
    const machineColumns = [
        table.column({
            accessor: "name",
            header: "Name",
        }),
        table.column({
            accessor: "host",
            header: "Host",
        }),
        table.column({
            accessor: "groups",
            header: "Groups",
            cell: ({ value }) => {
                return createRender(TableGroups, { id: value });
            },
            plugins: {
                sort: {
                    disable: true,
                },
                colFilter: {
                    fn: ({ filterValue, value }): boolean => {
                        return value ? value.includes(filterValue) : false;
                    },
                },
            },
        }),
        table.column({
            accessor: "tags",
            header: "Tags",
            cell: ({ value }) => {
                return createRender(TableTags, { id: value });
            },
            plugins: {
                sort: {
                    disable: true,
                },
            },
        }),
        table.column({
            accessor: "error",
            header: "Errors",
            plugins: {
                sort: {
                    disable: true,
                },
                filter: {
                    exclude: true,
                },
            },
        }),
        table.column({
            accessor: "agent",
            header: "Agent",
            cell: ({ value }) => {
                return createRender(TableBool, { isTrue: value });
            },
            plugins: {
                sort: {
                    disable: true,
                },
                filter: {
                    exclude: true,
                },
            },
        }),
    ];

    const userColumns = [
        table.column({
            accessor: "username",
            header: "Username",
        }),

        table.column({
            accessor: "email",
            header: "Email",
        }),
        table.column({
            accessor: "principal",
            header: "Principal",
            plugins: {
                sort: {
                    disable: true,
                },
            },
        }),
        table.column({
            accessor: "permission",
            header: "Permission",
            cell: ({ value }) => {
                return createRender(TablePermissions, { id: value });
            },
            plugins: {
                sort: {
                    disable: true,
                },
            },
        }),
        table.column({
            accessor: "groups",
            header: "Groups",
            cell: ({ value }) => {
                return createRender(TableGroups, { id: value });
            },
            plugins: {
                sort: {
                    disable: true,
                },
                colFilter: {
                    fn: ({ filterValue, value }): boolean => {
                        return value ? value.includes(filterValue) : false;
                    },
                },
            },
        }),
    ];

    const groupColumns = [
        table.column({
            accessor: "name",
            header: "Name",
        }),
        table.column({
            accessor: "description",
            header: "Description",
            plugins: {
                sort: {
                    disable: true,
                },
                filter: {
                    exclude: true,
                },
            },
        }),
        table.column({
            accessor: "linux_username",
            header: "Linux Username",
        }),
    ];

    const requestColumns = [
        table.column({
            accessor: "user",
            header: "User",
            cell: ({ value }) => {
                return $users.find((u) => u.id === value)?.username;
            },
        }),
        table.column({
            accessor: "group",
            header: "Group",
            cell: ({ value }) => {
                return $groups.find((g) => g.id === value)?.name || "";
            },
        }),
        table.column({
            accessor: "machine",
            header: "Machine",
            cell: ({ value }) => {
                return $machines.find((g) => g.id === value)?.name || "";
            },
        }),
        table.column({
            accessor: "created",
            header: "Created",
        }),
        table.column({
            accessor: ({ id }) => id,
            header: "",
            cell: ({ value }) => {
                return createRender(TableRequestActions, {
                    id: value,
                });
            },
            plugins: {
                sort: {
                    disable: true,
                },
                filter: {
                    exclude: true,
                },
            },
        }),
    ];

    const providerColumns = [
        table.column({
            accessor: "name",
            header: "Name",
        }),
        table.column({
            accessor: "type",
            header: "Type",
        }),
        table.column({
            accessor: "error",
            header: "Error",
        }),
        table.column({
            accessor: "last_sync",
            header: "Last Sync",
        }),
    ];

    const permissionColumns = [
        table.column({
            accessor: "name",
            header: "Name",
        }),
        table.column({
            accessor: "is_admin",
            header: "Admin",
            cell: ({ value }) => {
                return createRender(TableBool, { isTrue: value });
            },
            plugins: {
                sort: {
                    disable: true,
                },
                filter: {
                    exclude: true,
                },
            },
        }),
        table.column({
            accessor: (item) => item,
            header: "Details",
            cell: ({ value }) => {
                return createRender(TableAccess, { permission: value });
            },
            plugins: {
                sort: {
                    disable: true,
                },
                filter: {
                    exclude: true,
                },
            },
        }),
    ];

    const auditlogColumns = [
        table.column({
            accessor: "collection",
            header: "Collection",
        }),
        table.column({
            accessor: (item) => item,
            header: "Action",
            cell: ({ value }) => {
                return createRender(TableLog, { log: value });
            },
        }),
        table.column({
            accessor: "created",
            header: "Date",
        }),
    ];

    let columns = table.createColumns([]);
    let hidableCols: string[] = [];

    const limits = [10, 25, 50, 100];
    let {
        headerRows,
        pageRows,
        tableAttrs,
        tableBodyAttrs,
        pluginStates,
        flatColumns,
    } = table.createViewModel(columns);
    let { pageIndex, pageCount, pageSize, hasNextPage, hasPreviousPage } =
        pluginStates.page;
    let { filterValue } = pluginStates.filter;
    let { filterValues } = pluginStates.colFilter;
    let { hiddenColumnIds } = pluginStates.hide;

    let hideForId = Object.fromEntries(
        flatColumns.map((col) => col.id).map((id) => [id, true]),
    );

    $: $hiddenColumnIds = Object.entries(hideForId)
        .filter(([, hide]) => !hide)
        .map(([id]) => id);

    onMount(() => {
        switch (collection) {
            case "machines":
                hidableCols = ["host", "groups", "tags", "error", "agent"];
                table = createTable(machines, plugins);
                columns = table.createColumns([
                    selectionColumn,
                    ...machineColumns,
                    actionsColumn,
                ]);
                break;
            case "users":
                hidableCols = [
                    "username",
                    "email",
                    "permission",
                    "groups",
                    "principal",
                ];
                table = createTable(users, plugins);
                columns = table.createColumns([
                    selectionColumn,
                    ...userColumns,
                    actionsColumn,
                ]);
                break;
            case "groups":
                hidableCols = ["name", "description", "linux_username"];
                table = createTable(groups, plugins);
                columns = table.createColumns([
                    selectionColumn,
                    ...groupColumns,
                    actionsColumn,
                ]);
                break;
            case "permissions":
                hidableCols = ["is_admin", "Details"];
                table = createTable(permissions, plugins);
                columns = table.createColumns([
                    selectionColumn,
                    ...permissionColumns,
                    actionsColumn,
                ]);
                break;
            case "providers":
                hidableCols = ["type", "error", "last_sync"];
                table = createTable(providers, plugins);
                columns = table.createColumns([
                    selectionColumn,
                    ...providerColumns,
                    actionsColumn,
                ]);
                break;
            case "requests":
                hidableCols = ["created"];
                table = createTable(requests, plugins);
                columns = table.createColumns([
                    selectionColumn,
                    ...requestColumns,
                ]);
                break;
            case "auditlog":
                hidableCols = ["collection", "Action", "created"];
                table = createTable(auditlog, plugins);
                columns = table.createColumns([
                    selectionColumn,
                    ...auditlogColumns,
                ]);
                break;
            default:
                hidableCols = [];
        }

        // Regenerate the table view model every time columns change
        ({
            headerRows,
            pageRows,
            tableAttrs,
            tableBodyAttrs,
            pluginStates,
            flatColumns,
        } = table.createViewModel(columns));

        // Extract plugin states
        ({ pageIndex, pageCount, pageSize, hasNextPage, hasPreviousPage } =
            pluginStates.page);
        filterValue = pluginStates.filter.filterValue;
        filterValues = pluginStates.colFilter.filterValues;
        hiddenColumnIds = pluginStates.hide.hiddenColumnIds;
        hideForId = Object.fromEntries(
            flatColumns.map((col) => col.id).map((id) => [id, true]),
        );
    });
</script>

<div class="flex items-center py-4">
    <Input
        class="max-w-sm"
        placeholder="Search..."
        type="text"
        bind:value={$filterValue}
    />
    <DropdownMenu.Root closeOnItemClick={false}>
        <DropdownMenu.Trigger asChild let:builder>
            <Button variant="outline" class="ml-auto" builders={[builder]}>
                Columns
                <ChevronDown size="1rem" class="ml-1" />
            </Button>
        </DropdownMenu.Trigger>
        <DropdownMenu.Content>
            {#each flatColumns as col}
                {#if hidableCols.includes(col.id)}
                    <DropdownMenu.CheckboxItem bind:checked={hideForId[col.id]}>
                        {col.header}
                    </DropdownMenu.CheckboxItem>
                {/if}
            {/each}
        </DropdownMenu.Content>
    </DropdownMenu.Root>
</div>

<div class="rounded-md border">
    <Table.Root {...$tableAttrs}>
        <Table.Header>
            {#each $headerRows as headerRow}
                <Subscribe rowAttrs={headerRow.attrs()}>
                    <Table.Row>
                        {#each headerRow.cells as cell (cell.id)}
                            <Subscribe
                                attrs={cell.attrs()}
                                let:attrs
                                props={cell.props()}
                                let:props
                            >
                                <Table.Head {...attrs}>
                                    {#if !props.sort.disabled}
                                        <span
                                            on:click={props.sort.toggle}
                                            class="flex items-center gap-1 cursor-pointer"
                                            aria-hidden
                                        >
                                            <Render of={cell.render()} />
                                            {#if props.sort.order === "asc"}
                                                <ArrowUp class="h-4 w-4" />
                                            {:else}
                                                <ArrowDown class="h-4 w-4" />
                                            {/if}
                                        </span>
                                    {:else}
                                        <Render of={cell.render()} />
                                    {/if}
                                </Table.Head>
                            </Subscribe>
                        {/each}
                    </Table.Row>
                </Subscribe>
            {/each}
        </Table.Header>
        <Table.Body {...$tableBodyAttrs}>
            {#each $pageRows as row (row.id)}
                <Subscribe rowAttrs={row.attrs()} let:rowAttrs>
                    <Table.Row {...rowAttrs}>
                        {#each row.cells as cell (cell.id)}
                            <Subscribe attrs={cell.attrs()} let:attrs>
                                <Table.Cell {...attrs}>
                                    <Render of={cell.render()} />
                                </Table.Cell>
                            </Subscribe>
                        {/each}
                    </Table.Row>
                </Subscribe>
            {/each}
        </Table.Body>
    </Table.Root>
</div>

<div class="flex items-center justify-between py-4">
    <Select.Root
        selected={{ value: $pageSize, label: $pageSize.toString() }}
        onSelectedChange={(e) => e && ($pageSize = e.value)}
    >
        <Select.Trigger class="w-[180px]">
            <Select.Value placeholder="Limit" />
        </Select.Trigger>
        <Select.Content>
            {#each limits as limit}
                <Select.Item value={limit} label={limit.toString()}
                    >{limit}</Select.Item
                >
            {/each}
        </Select.Content>
    </Select.Root>

    <div class="flex items-center gap-2">
        <Button
            variant="outline"
            size="sm"
            on:click={() => ($pageIndex = $pageIndex - 1)}
            disabled={!$hasPreviousPage}>Previous</Button
        >
        <span class="text-sm">
            {$pageIndex + 1}-{$pageSize} of {$pageCount}
        </span>
        <Button
            variant="outline"
            size="sm"
            disabled={!$hasNextPage}
            on:click={() => ($pageIndex = $pageIndex + 1)}>Next</Button
        >
    </div>
</div>
