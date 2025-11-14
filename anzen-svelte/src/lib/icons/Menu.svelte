<!-- lib/icons/Menu.svelte -->
<svg
        viewBox="0 0 14 14"
        fill="none"
        width={size}
        height={size}
        class={className}
        {...$$restProps}
>
    <path
            d="M2 3H12M2 7H12M2 11H8"
            stroke="currentColor"
            stroke-width="1.5"
            stroke-linecap="round"
    />
</svg>

<script>
    export let size = 14;
    export let className = '';
</script>