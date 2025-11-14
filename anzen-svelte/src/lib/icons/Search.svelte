<!-- lib/icons/Search.svelte -->
<svg
        viewBox="0 0 16 16"
        fill="none"
        width={size}
        height={size}
        class={className}
        {...$$restProps}
>
    <circle cx="7" cy="7" r="5" stroke="currentColor" stroke-width="1.5"/>
    <path d="M11 11L14 14" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
</svg>

<script>
    export let size = 16;
    export let className = '';
</script>