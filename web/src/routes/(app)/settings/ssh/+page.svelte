<script lang="ts">
    import { pb } from "$lib/client";
    import { settings } from "$lib/subscriptions";
    import {
        clipboard,
        getModalStore,
        type ModalSettings,
    } from "@skeletonlabs/skeleton";
    import { onMount } from "svelte";

    let userKey = "";
    let hostKey = "";
    let agentToken = "";

    const rotateToken = async () => {
        await pb.send("/api/rpc/token/rotate", { method: "POST" });
        agentToken = await pb
            .send("/api/rpc/token", {})
            .then((res) => res.token);
    };
    const syncToken = async () => {
        await pb.send("/api/sync/token", { method: "POST" });
    };

    const modals = getModalStore();
    const updateKey = async () => {
        const modal: ModalSettings = {
            type: "component",
            title: "Update User CA",
            component: "UpdateUserCA",
        };
        modals.trigger(modal);
    };
    const signHostKey = () => {
        const modal: ModalSettings = {
            type: "component",
            title: "Sign Host Key",
            component: "SignHostKey",
        };
        modals.trigger(modal);
    };

    let icons = [
        "fa6-solid:clipboard",
        "fa6-solid:clipboard",
        "fa6-solid:clipboard",
        "fa6-solid:clipboard",
    ];
    const copyToClipboard = (index: number) => {
        icons[index] = "fa6-solid:thumbs-up";
        setTimeout(() => {
            icons[index] = "fa6-solid:clipboard";
        }, 2000);
    };
    const selectText = (e: any) => {
        e.target.select();
        navigator.clipboard.writeText(e.target.value);
    };

    onMount(async () => {
        if (!pb.authStore.model) return;
        userKey = await pb
            .send("/api/ssh/user/public", {})
            .then((res) => res.key);
        hostKey = await pb
            .send("/api/ssh/host/public", {})
            .then((res) => "@cert-authority * " + res.key);
        agentToken = await pb
            .send("/api/rpc/token", {})
            .then((res) => res.token);
    });
</script>

<div class="flex flex-col gap-4 w-full">
    <div class="card">
        <header class="card-header text-2xl font-bold flex items-center">
            Agent Token
            <button
                class="btn btn-sm variant-filled-surface ml-2"
                on:click={rotateToken}>Rotate Token</button
            >
            <button
                class="btn btn-sm variant-filled-surface ml-2"
                on:click={syncToken}>Sync Token</button
            >
        </header>
        <section class="px-4 py-2">
            <span class="text-sm dark:text-surface-300">
                This is current agent token, use it to register new agents.
                <br />
                <span class="font-bold underline"
                    >In case of an emergency (leaked key) you can rotate the
                    token and generate a new one. Keep in mind that this will
                    disconnect all current agents! Using the 'Sync Token' button
                    you can send the new token to all agents. Do so only after
                    removing the compromised machines!
                </span>
            </span>
        </section>
        <section class="flex flex-row items-center p-4">
            <input
                class="input"
                type="text"
                bind:value={agentToken}
                data-clipboard="agentToken"
                on:click={selectText}
                readonly
            />
            <button
                class="btn variant-filled-surface ml-2"
                use:clipboard={{ input: "agentToken" }}
                on:click={() => copyToClipboard(0)}
                value="agentToken"
            >
                <iconify-icon icon={icons[0]} />
            </button>
        </section>
    </div>

    <div class="card">
        <header class="card-header text-2xl font-bold">
            User CA Public Key
            <button
                class="btn btn-sm variant-filled-surface ml-2"
                on:click={updateKey}>Overwrite Key</button
            >
        </header>
        <section class="px-4 py-2">
            <span class="text-sm dark:text-surface-300">
                This is the public key of our user certificate authority and
                will be automatically installed on the servers.
                <br />
                <span class="font-bold underline"
                    >This key will also be used by the server to login, update
                    machines manually and for installing agents on new machines.</span
                >
            </span>
        </section>
        <section class="flex flex-row items-center p-4">
            <input
                class="input"
                type="text"
                bind:value={userKey}
                on:click={selectText}
                data-clipboard="userKey"
                readonly
            />
            <button
                class="btn variant-filled-surface ml-2"
                use:clipboard={{ input: "userKey" }}
                on:click={() => copyToClipboard(1)}
                value="serverKey"
            >
                <iconify-icon icon={icons[1]} />
            </button>
        </section>
    </div>

    <div class="card">
        <header class="card-header text-2xl font-bold">
            Host CA Public Key
            <button
                class="btn btn-sm variant-filled-surface ml-2"
                on:click={signHostKey}>Sign Host Key</button
            >
        </header>
        <section class="px-4 py-2">
            <span class="text-sm dark:text-surface-300">
                This is the public key of our host certificate authority.
                <br />
                Put this into your
                <span class="font-bold underline">known_hosts</span> file to automatically
                trust the servers.
            </span>
        </section>
        <section class="flex flex-row items-center p-4">
            <input
                class="input"
                type="text"
                bind:value={hostKey}
                on:click={selectText}
                data-clipboard="hostKey"
                readonly
            />
            <button
                class="btn variant-filled-surface ml-2"
                use:clipboard={{ input: "hostKey" }}
                on:click={() => copyToClipboard(2)}
                value="serverKey"
            >
                <iconify-icon icon={icons[2]} />
            </button>
        </section>
    </div>

    <div class="card">
        <header class="card-header text-2xl font-bold">OpenSSH Config</header>
        <section class="px-4 py-2">
            <span class="text-sm dark:text-surface-300">
                This is the custom OpenSSH server configuration file that will
                be installed on the machines.
            </span>
        </section>
        <section class="flex flex-row items-center p-4">
            {#each $settings as setting}
                {#if setting.key === "ssh_config"}
                    <textarea
                        class="textarea"
                        rows={setting.value.split("\n").length}
                        bind:value={setting.value}
                        data-clipboard="sshConfig"
                    />
                    <button
                        class="btn variant-filled-surface ml-2"
                        use:clipboard={{ input: "sshConfig" }}
                        on:click={() => copyToClipboard(3)}
                        value="sshConfig"
                        ><iconify-icon icon={icons[3]} />
                    </button>
                {/if}
            {/each}
        </section>
    </div>
</div>
