<!-- lib/icons/Code.svelte -->
<svg
        viewBox="0 0 16 16"
        fill="none"
        width={size}
        height={size}
        class={className}
        {...$$restProps}
>
    <path
            d="M2 2L6 6M6 6V2M6 6H2"
            stroke="currentColor"
            stroke-width="1.5"
            stroke-linecap="round"
    />
</svg>

<script>
    export let size = 16;
    export let className = '';
</script>