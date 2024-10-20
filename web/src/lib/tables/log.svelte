<script lang="ts">
    import { Badge } from "$lib/components/ui/badge";
    import type { RecordModel } from "pocketbase";

    export let log: RecordModel;

    // Nice messages
    function getCollectionName(e: RecordModel) {
        switch (e.collection) {
            case "users":
            case "machines":
            case "groups":
            case "permissions":
                return e.data.name;
            case "settings":
                return e.data.key;
            default:
                return "Unknown"; // Handle unknown collections
        }
    }

    function getEventMessage(e: RecordModel) {
        if (!e.collection || !e.event) return null;

        let message = `${e.event}d`; // Base message with "d" suffix

        switch (e.collection) {
            case "users":
            case "machines":
            case "settings":
                return message; // No additional message for these collections
            default:
                return message;
        }
    }

    function getTriggeredBy(e: RecordModel) {
        return e.expand?.user?.name || "Admin";
    }
</script>

<div>
    <Badge variant="secondary">
        {getCollectionName(log)}
    </Badge>
    {#if log.event === "create"}
        <Badge variant="secondary" class="bg-green-300 text-gray-800">
            {getEventMessage(log)}
        </Badge>
    {/if}
    {#if log.event === "delete"}
        <Badge variant="secondary" class="bg-red-300 text-gray-800">
            {getEventMessage(log)}
        </Badge>
    {/if}
    {#if log.event === "update"}
        <Badge variant="secondary" class="bg-orange-300 text-gray-800">
            {getEventMessage(log)}
        </Badge>
    {/if}
    by
    <Badge variant="secondary">
        {getTriggeredBy(log)}
    </Badge>
</div>
