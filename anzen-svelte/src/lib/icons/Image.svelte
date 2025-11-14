<!-- lib/icons/Image.svelte -->
<svg
        viewBox="0 0 14 14"
        fill="none"
        width={size}
        height={size}
        class={className}
        {...$$restProps}
>
    <rect x="2" y="2" width="10" height="10" rx="1" stroke="currentColor" stroke-width="1.5"/>
</svg>

<script>
    export let size = 14;
    export let className = '';
</script>