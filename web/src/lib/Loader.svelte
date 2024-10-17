<script lang="ts">
    type SpinnerTypes = {
        size: string | number;
        unit: string;
        duration: string;
        pause: boolean;
    };
    const durationUnitRegex = /[a-zA-Z]/;
    export const range = (size: number, startAt = 0) =>
        [...Array(size).keys()].map((i) => i + startAt);

    export let unit: SpinnerTypes["unit"] = "px";
    export let duration: SpinnerTypes["duration"] = "1s";
    export let size: SpinnerTypes["size"] = "60";
    export let pause: SpinnerTypes["pause"] = false;
    let durationUnit: string = duration.match(durationUnitRegex)?.[0] ?? "s";
    let durationNum: string = duration.replace(durationUnitRegex, "");
</script>

<div class="wrapper" style="--size: {size}{unit}; --duration: {duration};">
    {#each range(3, 1) as version}
        <div
            class="circle bg-primary-500 top-1/3 left-1/2"
            class:pause-animation={pause}
            style="animation-delay: {(+durationNum / 3) * (version - 1) +
                durationUnit};"
        />
    {/each}
</div>

<style>
    .wrapper {
        width: var(--size);
        height: var(--size);
    }
    .circle {
        border-radius: 100%;
        animation-fill-mode: both;
        position: absolute;
        opacity: 0;
        width: var(--size);
        height: var(--size);
        animation: bounce var(--duration) linear infinite;
    }
    .pause-animation {
        animation-play-state: paused;
    }
    @keyframes bounce {
        0% {
            opacity: 0;
            transform: scale(0);
        }
        5% {
            opacity: 1;
        }
        100% {
            opacity: 0;
            transform: scale(1);
        }
    }
</style>
