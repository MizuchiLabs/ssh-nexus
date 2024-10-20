<script lang="ts">
    import * as Card from "$lib/components/ui/card/index.js";
    import { Switch } from "$lib/components/ui/switch/index.js";
    import { Input } from "$lib/components/ui/input/index.js";
    import { pb } from "$lib/client";
    import { settings } from "$lib/subscriptions";
    import type { RecordModel, ClientResponseError } from "pocketbase";
    import { toast } from "svelte-sonner";

    const updateSettings = async (setting: RecordModel) => {
        try {
            await pb.collection("settings").update(setting.id, setting);
            toast.success(`Updated ${setting.key} to ${setting.value}`);
        } catch (error: ClientResponseError | any) {
            toast.error(error.data?.message || "Something went wrong.");
        }
    };

    const onKeys = (e: KeyboardEvent, setting: RecordModel) => {
        if (e.key === "Enter") {
            updateSettings(setting);
        }
    };

    const changeBool = (setting: RecordModel) => {
        setting.value = setting.value === "true" ? "false" : "true";
        updateSettings(setting);
    };

    function formatDuration(seconds: number) {
        const secondsInMinute = 60;
        const secondsInHour = secondsInMinute * 60;
        const secondsInDay = secondsInHour * 24;

        const days = Math.floor(seconds / secondsInDay);
        seconds -= days * secondsInDay;

        const hours = Math.floor(seconds / secondsInHour);
        seconds -= hours * secondsInHour;

        const minutes = Math.floor(seconds / secondsInMinute);
        seconds -= minutes * secondsInMinute;

        const remainingSeconds = seconds;

        const parts = [];

        if (days > 0) {
            if (days === 1) {
                parts.push("1 day");
            } else {
                parts.push(`${days} days`);
            }
        }
        if (hours > 0) {
            if (hours === 1) {
                parts.push("1 hour");
            } else {
                parts.push(`${hours} hours`);
            }
        }
        if (minutes > 0) {
            if (minutes === 1) {
                parts.push("1 minute");
            } else {
                parts.push(`${minutes} minutes`);
            }
        }
        if (remainingSeconds > 0) {
            if (remainingSeconds === 1) {
                parts.push("1 second");
            } else {
                parts.push(`${remainingSeconds} seconds`);
            }
        }

        return parts.join(" ");
    }
</script>

<div class="flex flex-col gap-2">
    {#each $settings as setting}
        {#if setting.key === "default_retention"}
            <Card.Root>
                <Card.Header>
                    <Card.Title>
                        Default Audit Log Retention
                        <span class="text-gray-400/75 ml-2 text-sm">
                            {formatDuration(setting.value)}
                        </span>
                    </Card.Title>
                    <Card.Description>
                        Set the default audit log retention (in seconds)
                    </Card.Description>
                </Card.Header>
                <Card.Content>
                    <Input
                        type="text"
                        bind:value={setting.value}
                        placeholder={setting.value}
                        on:keydown={(e) => onKeys(e, setting)}
                    />
                </Card.Content>
            </Card.Root>
        {/if}
        {#if setting.key === "user_lease"}
            <Card.Root>
                <Card.Header>
                    <Card.Title>
                        Default TTL of new user ssh certificates
                        <span class="text-gray-400/75 ml-2 text-sm">
                            {formatDuration(setting.value)}
                        </span>
                    </Card.Title>
                    <Card.Description>
                        Set the default TTL of new user ssh certificates (in
                        seconds)
                    </Card.Description>
                </Card.Header>
                <Card.Content>
                    <Input
                        type="text"
                        bind:value={setting.value}
                        placeholder={setting.value}
                        on:keydown={(e) => onKeys(e, setting)}
                    />
                </Card.Content>
            </Card.Root>
        {/if}
        {#if setting.key === "max_lease"}
            <Card.Root>
                <Card.Header>
                    <Card.Title>
                        Max TTL of new user ssh certificates
                        <span class="text-gray-400/75 ml-2 text-sm">
                            {formatDuration(setting.value)}
                        </span>
                    </Card.Title>
                    <Card.Description>
                        Set the maximum TTL of new user ssh certificates (in
                        seconds)
                    </Card.Description>
                </Card.Header>
                <Card.Content>
                    <Input
                        type="text"
                        bind:value={setting.value}
                        placeholder={setting.value}
                        on:keydown={(e) => onKeys(e, setting)}
                    />
                </Card.Content>
            </Card.Root>
        {/if}
        {#if setting.key === "host_lease"}
            <Card.Root>
                <Card.Header>
                    <Card.Title>
                        Max TTL of new host ssh certificates
                        <span class="text-gray-400/75 ml-2 text-sm">
                            {formatDuration(setting.value)}
                        </span>
                    </Card.Title>
                    <Card.Description>
                        Set the maximum TTL of new host ssh certificates (in
                        seconds)
                    </Card.Description>
                </Card.Header>
                <Card.Content>
                    <Input
                        type="text"
                        bind:value={setting.value}
                        placeholder={setting.value}
                        on:keydown={(e) => onKeys(e, setting)}
                    />
                </Card.Content>
            </Card.Root>
        {/if}
        {#if setting.key === "install_agent"}
            <Card.Root>
                <Card.Content
                    class="flex flex-row justify-between items-center"
                >
                    <div>
                        <h2 class="text-lg font-semibold">
                            Install Agents on new machines
                        </h2>
                        <p class="text-gray-400/75 text-sm">
                            Check if agents should be installed per default on
                            newly created machines
                        </p>
                    </div>
                    <Switch
                        id="toggleAgents"
                        checked={setting.value === "true"}
                        onCheckedChange={(_) => changeBool(setting)}
                    />
                </Card.Content>
            </Card.Root>
        {/if}
    {/each}
</div>
