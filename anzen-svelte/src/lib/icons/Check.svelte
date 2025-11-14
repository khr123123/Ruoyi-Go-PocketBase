<!-- lib/icons/Check.svelte -->
<svg
        viewBox="0 0 14 14"
        fill="none"
        width={size}
        height={size}
        class={className}
        {...$$restProps}
>
    <path
            d="M4 7L6 9L10 5"
            stroke="currentColor"
            stroke-width="1.5"
            stroke-linecap="round"
    />
</svg>

<script>
    export let size = 14;
    export let className = '';
</script>