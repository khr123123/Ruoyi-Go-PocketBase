<!-- lib/icons/Refresh.svelte -->
<svg
        viewBox="0 0 20 20"
        fill="none"
        width={size}
        height={size}
        class={className}
        {...$$restProps}
>
    <path
            d="M17.5 7.5C16.5 4.5 13.5 2.5 10 2.5C5.86 2.5 2.5 5.86 2.5 10C2.5 14.14 5.86 17.5 10 17.5C13.5 17.5 16.5 15.5 17.5 12.5"
            stroke="currentColor"
            stroke-width="1.5"
            stroke-linecap="round"
    />
    <path
            d="M14 7.5H17.5V4"
            stroke="currentColor"
            stroke-width="1.5"
            stroke-linecap="round"
            stroke-linejoin="round"
    />
</svg>

<script>
    export let size = 20;
    export let className = '';
</script>