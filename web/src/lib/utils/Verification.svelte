<script lang="ts">
    import { pb, user } from "$lib/client";
    import { showToast } from "$lib/utils/Toast";
    import { getToastStore } from "@skeletonlabs/skeleton";

    const toastStore = getToastStore();
    const verifyEmail = async () => {
        await pb
            .collection("users")
            .requestVerification(pb.authStore.model?.email);
        showToast(toastStore, `Email sent 👌`, "success");
    };
</script>

{#if $user?.verified === false}
    <div
        class="flex items-center justify-center w-full p-2 bg-primary-400 dark:bg-surface-700"
    >
        <span class="flex items-center text-sm font-normal dark:text-gray-400"
            >You are not verified yet,
            <p
                on:click={verifyEmail}
                class="font-medium text-success-500 dark:text-primary-300 hover:text-success-300 cursor-pointer px-1"
                aria-hidden
            >
                click here
            </p>
            to verify your email.
        </span>
    </div>
{/if}
