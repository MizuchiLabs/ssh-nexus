<script lang="ts">
    import { pb, user } from "$lib/client";
    import { toast } from "svelte-sonner";

    const verifyEmail = async () => {
        await pb
            .collection("users")
            .requestVerification(pb.authStore.model?.email);
        toast.success("Verification email sent");
    };
</script>

{#if $user?.verified === false}
    <div
        class="flex items-center justify-center w-full p-2 bg-red-300 dark:bg-red-400"
    >
        <span class="flex items-center text-sm font-normal dark:text-black"
            >You are not verified yet,
            <p
                on:click={verifyEmail}
                class="font-medium text-purple-500 dark:text-purple-700 cursor-pointer px-1"
                aria-hidden
            >
                click here
            </p>
            to verify your email.
        </span>
    </div>
{/if}
