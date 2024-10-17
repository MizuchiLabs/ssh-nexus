<script lang="ts">
    import { pb } from "$lib/client";
    import {
        Autocomplete,
        InputChip,
        type AutocompleteOption,
    } from "@skeletonlabs/skeleton";
    import type { RecordModel } from "pocketbase";

    export let data = [] as RecordModel[]; // Base list, shouldn't change
    export let init = [] as string[];
    export let name = "";
    let input = "";
    let options = [] as AutocompleteOption<string>[];
    let whitelist = [""] as string[];
    let chips = [] as string[];

    $: data, mapBase();
    const mapBase = () => {
        if (data.length === 0) return;
        whitelist = data.map((u) => u.name);
        chips = data.filter((u) => init.includes(u.id)).map((u) => u.name);
        options = data.map((u) => {
            return {
                label: u.name,
                value: u.name,
                meta: {
                    id: u.id,
                },
            };
        });
    };

    // Function to add a new item to the chips array
    async function add({ detail: e }: CustomEvent) {
        input = ""; // Reset input field after adding
        chips = [...chips, e.value];
        init = data.filter((u) => chips.includes(u.name)).map((u) => u.id);
    }

    // Function to remove a item from the chips array
    function remove({ detail: e }: CustomEvent) {
        chips = chips.filter((chip) => chip !== e.chipValue);
        init = data.filter((u) => chips.includes(u.name)).map((u) => u.id);
    }

    async function onInvalidHandler() {
        if (name === "tags") {
            let record = await pb.collection("tags").create({ name: input });
            data = [...data, record];
            chips = [...chips, input];
            init = data.filter((u) => chips.includes(u.name)).map((u) => u.id);
            input = "";
        }
    }
</script>

<label class="label">
    <span class="capitalize">{name}</span>
    <input class="input hidden" type="text" {name} />
    <InputChip
        name="chips"
        class="px-4 variant-form-material"
        placeholder="Type to search/add..."
        bind:input
        bind:value={chips}
        on:remove={remove}
        on:invalid={onInvalidHandler}
        {whitelist}
    />
    {#if input}
        <div
            class="w-full max-w max-h-48 p-2 overflow-y-auto border border-surface-500"
            tabindex="-1"
        >
            <Autocomplete
                bind:input
                {options}
                denylist={chips}
                allowlist={whitelist}
                on:selection={add}
            />
        </div>
    {/if}
</label>
