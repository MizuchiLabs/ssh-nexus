<script lang="ts">
    import { pb, user } from "$lib/client";
    import {
        getModalStore,
        popup,
        type ModalSettings,
        type PopupSettings,
    } from "@skeletonlabs/skeleton";

    let open = false;

    const modals = getModalStore();
    const signUserKey = () => {
        const modal: ModalSettings = {
            type: "component",
            title: "Sign User Key",
            component: "SignUserKey",
        };
        modals.trigger(modal);
    };

    function clickOutside(node: any) {
        function onClick(event: MouseEvent) {
            if (!node.contains(event.target)) {
                open = false;
            }
        }

        document.body.addEventListener("click", onClick);
        return {
            destroy() {
                document.body.removeEventListener("click", onClick);
            },
        };
    }

    // User is admin and doesn't need to request permissions
    const canRequest = () => {
        if ($user?.expand?.permission?.is_admin || pb.authStore.isAdmin) {
            return false;
        }
        if (
            $user?.expand?.permission?.access_machines &&
            $user?.expand?.permission?.can_update
        ) {
            return false;
        }
        return true;
    };

    let helpText = "";
    const popupHelp: PopupSettings = {
        event: "hover",
        target: "popupHelp",
        placement: "top",
    };
</script>

<div class="fixed flex flex-row gap-1 bottom-4 right-4" use:clickOutside>
    <div class="card p-4 text-sm variant-soft-surface" data-popup="popupHelp">
        <p>{helpText}</p>
        <div class="arrow variant-soft-surface" />
    </div>

    {#if open}
        <div class="flex items-center gap-1" use:popup={popupHelp}>
            {#if pb.authStore.isAdmin}
                <a
                    href={pb.baseUrl + "/_/"}
                    class="btn btn-icon-sm variant-soft-primary hover:variant-ghost-primary"
                    on:mouseover={() => (helpText = "Backend")}
                    on:focus={() => (helpText = "Backend")}
                >
                    <iconify-icon icon="fa6-solid:database" width="16" />
                </a>
            {/if}
            {#if canRequest()}
                <a
                    href="/request/machines"
                    class="btn btn-icon-sm variant-soft-primary"
                    on:mouseover={() => (helpText = "Request Machines")}
                    on:focus={() => (helpText = "Request Machines")}
                >
                    <iconify-icon icon="fa6-solid:circle-plus" width="16" />
                </a>
            {/if}
            <button
                class="btn btn-icon-sm variant-soft-primary"
                on:click={signUserKey}
                on:mouseover={() => (helpText = "Sign User Key")}
                on:focus={() => (helpText = "Sign User Key")}
            >
                <iconify-icon icon="fa6-solid:key" width="16" />
            </button>
        </div>
    {/if}

    <button
        class="btn btn-icon-sm bg-primary-500/40 hover:bg-primary-500/60"
        on:click={() => (open = !open)}
    >
        <iconify-icon icon="fa6-solid:plus" width="16" />
    </button>
</div>
