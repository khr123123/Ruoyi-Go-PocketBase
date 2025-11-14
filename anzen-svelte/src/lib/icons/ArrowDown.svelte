<!-- lib/icons/ArrowDown.svelte -->
<svg
        viewBox="0 0 14 14"
        fill="none"
        width={size}
        height={size}
        class={className}
        {...$$restProps}
>
    <path
            d="M7 2V12M4 9L7 12L10 9"
            stroke="currentColor"
            stroke-width="1.5"
            stroke-linecap="round"
    />
</svg>

<script>
    export let size = 14;
    export let className = '';
</script>