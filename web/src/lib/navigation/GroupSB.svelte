<script lang="ts">
    import { pb } from "$lib/client";
    import { groups } from "$lib/subscriptions";
    import { getModalStore, type ModalSettings } from "@skeletonlabs/skeleton";
    import type { RecordModel } from "pocketbase";

    export let group = "";
    $: sortedGroups = [...$groups].sort((a, b) => {
        if (a.name === "default") return -1;
        if (b.name === "default") return 1;
        return a.name.localeCompare(b.name);
    });

    // Modals
    const modals = getModalStore();
    const createGroup = () => {
        const modal: ModalSettings = {
            type: "component",
            title: "New Group",
            component: "GroupCreate",
        };
        modals.trigger(modal);
    };
    const updateGroup = (group: RecordModel) => {
        const modal: ModalSettings = {
            type: "component",
            title: "Edit Group",
            component: "GroupUpdate",
            meta: { group: group },
        };
        modals.trigger(modal);
    };

    // Drag and drop objects into groups
    export let records: RecordModel[] | undefined;
    const handleAddDrop = async (group: RecordModel) => {
        if (!records) return;
        for (const record of records) {
            if (!record.selected && records.length > 1) continue;

            if (record.collectionName === "users") {
                await pb.collection("users").update(record.id, {
                    "groups+": group.id,
                });
            }

            if (record.collectionName === "machines") {
                await pb.collection("machines").update(record.id, {
                    "groups+": group.id,
                });
            }
        }
    };
</script>

<div class="flex flex-col gap-4 min-w-52 font-mono">
    <button
        class="btn hover:variant-soft-surface"
        class:bg-primary-300={group === ""}
        class:bg-surface-300={group !== ""}
        class:dark:bg-primary-400={group === ""}
        class:dark:bg-surface-400={group !== ""}
        on:click={() => (group = "")}
    >
        Overview
    </button>
    {#each sortedGroups as sg}
        <button
            class="btn hover:variant-soft-surface flex justify-between items-center"
            class:bg-primary-300={group === sg.id}
            class:bg-surface-300={group !== sg.id}
            class:dark:bg-primary-400={group === sg.id}
            class:dark:bg-surface-500={group !== sg.id}
            on:drop={() => handleAddDrop(sg)}
            on:dragover|preventDefault
            on:click={() => (group = sg.id)}
            aria-hidden
        >
            <span class="text-center font-light">{sg.name}</span>
            <button
                class="btn-icon btn-icon-sm variant-filled-primary w-7 h-7"
                on:click={() => updateGroup(sg)}
            >
                <iconify-icon icon="fa6-solid:pen" />
            </button>
        </button>
    {/each}
    <button
        class="btn variant-filled-surface hover:variant-soft-surface"
        on:click={createGroup}
    >
        <iconify-icon icon="fa6-solid:plus" />
        <span>Add Group</span>
    </button>
</div>
