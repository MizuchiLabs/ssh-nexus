<script lang="ts">
    import { pb } from "$lib/client";
    import { settings } from "$lib/subscriptions";
    import * as Card from "$lib/components/ui/card/index.js";
    import { Input } from "$lib/components/ui/input/index.js";
    import { Button } from "$lib/components/ui/button/index.js";
    import { Textarea } from "$lib/components/ui/textarea/index.js";
    import SignHostKey from "$lib/modals/SignHostKey.svelte";
    import UpdateUserCa from "$lib/modals/UpdateUserCA.svelte";
    import { onMount } from "svelte";
    import { Copy } from "lucide-svelte";

    let userKey = "";
    let hostKey = "";
    let agentToken = "";
    let openHostKey = false;
    let openUserCA = false;

    const rotateToken = async () => {
        await pb.send("/api/rpc/token/rotate", { method: "POST" });
        agentToken = await pb
            .send("/api/rpc/token", {})
            .then((res) => res.token);
    };
    const syncToken = async () => {
        await pb.send("/api/sync/token", { method: "POST" });
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

<SignHostKey bind:open={openHostKey} />
<UpdateUserCa bind:open={openUserCA} />

<div class="flex flex-col gap-4 w-full">
    <Card.Root>
        <Card.Header>
            <Card.Title class="flex items-center gap-4">
                Agent Token
                <div>
                    <Button
                        variant="default"
                        class="h-8 rounded-full"
                        on:click={rotateToken}
                    >
                        Rotate Token
                    </Button>
                    <Button
                        variant="default"
                        class="h-8 rounded-full"
                        on:click={syncToken}
                    >
                        Sync Token
                    </Button>
                </div>
            </Card.Title>
            <Card.Description>
                <span class="text-sm dark:text-surface-300">
                    This is current agent token, use it to register new agents.
                    <br />
                    <span class="font-bold underline"
                        >In case of an emergency you can rotate the token and
                        generate a new one. Keep in mind that this will
                        disconnect all current agents! Using the 'Sync Token'
                        button you can send the new token to all agents. Do so
                        only after removing the compromised machines!
                    </span>
                </span>
            </Card.Description>
        </Card.Header>
        <Card.Content class="flex flex-row items-center justify-end gap-1">
            <Input
                type="text"
                bind:value={agentToken}
                on:click={selectText}
                class="pr-10"
                readonly
            />
            <Button
                variant="ghost"
                size="icon"
                class="h-8 rounded-full absolute hover:bg-transparent hover:text-red-400"
                on:click={() => navigator.clipboard.writeText(agentToken)}
            >
                <Copy size="1rem" />
            </Button>
        </Card.Content>
    </Card.Root>

    <Card.Root>
        <Card.Header>
            <Card.Title class="flex items-center gap-4">
                User CA Public Key

                <Button
                    variant="default"
                    class="h-8 rounded-full"
                    on:click={() => (openUserCA = true)}
                >
                    Overwrite Key
                </Button>
            </Card.Title>
            <Card.Description>
                <span class="text-sm dark:text-surface-300">
                    This is the public key of our user certificate authority and
                    will be automatically installed on the servers.
                    <br />
                    <span class="font-bold underline"
                        >This key will also be used by the server to login,
                        update machines manually and for installing agents on
                        new machines.</span
                    >
                </span>
            </Card.Description>
        </Card.Header>
        <Card.Content class="flex flex-row items-center justify-end gap-1">
            <Input
                type="text"
                bind:value={userKey}
                on:click={selectText}
                class="pr-10"
                readonly
            />
            <Button
                variant="ghost"
                size="icon"
                class="h-8 rounded-full absolute hover:bg-transparent hover:text-red-400"
                on:click={() => navigator.clipboard.writeText(userKey)}
            >
                <Copy size="1rem" />
            </Button>
        </Card.Content>
    </Card.Root>

    <Card.Root>
        <Card.Header>
            <Card.Title class="flex items-center gap-4">
                Host CA Public Key

                <Button
                    variant="default"
                    class="h-8 rounded-full"
                    on:click={() => (openHostKey = true)}
                >
                    Sign Host Key
                </Button>
            </Card.Title>
            <Card.Description>
                <span class="text-sm dark:text-surface-300">
                    This is the public key of our host certificate authority.
                    <br />
                    Put this into your
                    <span class="font-bold underline">known_hosts</span> file to
                    automatically trust the servers.
                </span>
            </Card.Description>
        </Card.Header>
        <Card.Content class="flex flex-row items-center justify-end gap-1">
            <Input
                type="text"
                bind:value={hostKey}
                on:click={selectText}
                class="pr-10"
                readonly
            />
            <Button
                variant="ghost"
                size="icon"
                class="h-8 rounded-full absolute hover:bg-transparent hover:text-red-400"
                on:click={() => navigator.clipboard.writeText(hostKey)}
            >
                <Copy size="1rem" />
            </Button>
        </Card.Content>
    </Card.Root>

    <Card.Root>
        <Card.Header>
            <Card.Title class="flex items-center gap-4">
                OpenSSH Config
            </Card.Title>
            <Card.Description>
                <span class="text-sm dark:text-surface-300">
                    This is the custom OpenSSH server configuration file that
                    will be installed on the machines.
                </span>
            </Card.Description>
        </Card.Header>
        <Card.Content class="flex flex-row items-center justify-end gap-1">
            {#each $settings as setting}
                {#if setting.key === "ssh_config"}
                    <Textarea
                        class="textarea"
                        rows={setting.value.split("\n").length}
                        bind:value={setting.value}
                        data-clipboard="sshConfig"
                    />
                    <Button
                        variant="ghost"
                        size="icon"
                        class="h-8 rounded-full absolute hover:bg-transparent hover:text-red-400"
                        on:click={() =>
                            navigator.clipboard.writeText(setting.value)}
                        value="sshConfig"
                    >
                        <Copy size="1rem" />
                    </Button>
                {/if}
            {/each}
        </Card.Content>
    </Card.Root>
</div>
