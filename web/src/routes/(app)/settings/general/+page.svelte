<script lang="ts">
    import { pb } from "$lib/client";
    import { settings } from "$lib/subscriptions";
    import { showToast } from "$lib/utils/Toast";
    import { SlideToggle, getToastStore } from "@skeletonlabs/skeleton";
    import type { RecordModel, ClientResponseError } from "pocketbase";

    const toastStore = getToastStore();
    const updateSettings = async (setting: RecordModel) => {
        try {
            await pb.collection("settings").update(setting.id, setting);
            showToast(
                toastStore,
                `Updated ${setting.key} to ${setting.value}`,
                "success",
            );
        } catch (error: ClientResponseError | any) {
            showToast(
                toastStore,
                error.data?.message || "Something went wrong.",
                "error",
            );
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

<div class="card flex flex-col gap-4 w-full">
    <header class="card-header text-2xl font-bold">Settings</header>
    <section class="p-4">
        <ul class="flex flex-col gap-6 list">
            {#each $settings as setting}
                {#if setting.key === "default_retention"}
                    <li on:keydown={(e) => onKeys(e, setting)} aria-hidden>
                        <iconify-icon icon="fa6-solid:clock-rotate-left" />
                        <span class="flex-auto">
                            <dt class="font-bold">
                                Default Audit Log Retention
                                <span class="text-surface-400 ml-2"
                                    >{formatDuration(setting.value)}</span
                                >
                            </dt>
                            <dd class="text-sm dark:text-surface-300">
                                Set the default audit log retention (in seconds)
                            </dd>
                        </span>
                        <input
                            class="input variant-form-material w-1/3"
                            type="text"
                            bind:value={setting.value}
                            placeholder={setting.value}
                        />
                    </li>
                {/if}
                {#if setting.key === "user_lease"}
                    <li on:keydown={(e) => onKeys(e, setting)} aria-hidden>
                        <iconify-icon icon="fa6-solid:hourglass-half" />
                        <span class="flex-auto">
                            <dt class="font-bold">
                                Default TTL of new user ssh certificates
                                <span class="text-surface-400 ml-2"
                                    >{formatDuration(setting.value)}</span
                                >
                            </dt>
                            <dd class="text-sm dark:text-surface-300">
                                Set the default TTL of new user ssh certificates
                                (in seconds)
                            </dd>
                        </span>
                        <input
                            class="input variant-form-material w-1/3"
                            type="text"
                            bind:value={setting.value}
                            placeholder={setting.value}
                        />
                    </li>
                {/if}
                {#if setting.key === "max_lease"}
                    <li on:keydown={(e) => onKeys(e, setting)} aria-hidden>
                        <iconify-icon icon="fa6-solid:hourglass-end" />
                        <span class="flex-auto">
                            <dt class="font-bold">
                                Max TTL of new user ssh certificates
                                <span class="text-surface-400 ml-2"
                                    >{formatDuration(setting.value)}</span
                                >
                            </dt>
                            <dd class="text-sm dark:text-surface-300">
                                Set the maximum TTL of new user ssh certificates
                                (in seconds)
                            </dd>
                        </span>
                        <input
                            class="input variant-form-material w-1/3"
                            type="text"
                            bind:value={setting.value}
                            placeholder={setting.value}
                        />
                    </li>
                {/if}
                {#if setting.key === "host_lease"}
                    <li on:keydown={(e) => onKeys(e, setting)} aria-hidden>
                        <iconify-icon icon="fa6-solid:hourglass-start" />
                        <span class="flex-auto">
                            <dt class="font-bold">
                                Max TTL of new host ssh certificates
                                <span class="text-surface-400 ml-2"
                                    >{formatDuration(setting.value)}</span
                                >
                            </dt>
                            <dd class="text-sm dark:text-surface-300">
                                Set the maximum TTL of new host ssh certificates
                                (in seconds)
                            </dd>
                        </span>
                        <input
                            class="input variant-form-material w-1/3"
                            type="text"
                            bind:value={setting.value}
                            placeholder={setting.value}
                        />
                    </li>
                {/if}
                {#if setting.key === "install_agent"}
                    <li>
                        <iconify-icon icon="fa6-solid:robot" />
                        <span class="flex-auto">
                            <dt class="font-bold">
                                Install Agents on new machines
                            </dt>
                            <dd class="text-sm dark:text-surface-300">
                                Check if agents should be installed per default
                                on newly created machines
                            </dd>
                        </span>
                        <SlideToggle
                            name="toggleAgents"
                            checked={setting.value === "true"}
                            active="bg-primary-600 dark:bg-primary-400"
                            size="sm"
                            on:change={() => changeBool(setting)}
                        />
                    </li>
                {/if}
            {/each}
        </ul>
    </section>
</div>
