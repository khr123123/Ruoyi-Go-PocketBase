<!-- lib/icons/Circle.svelte -->
<svg
        viewBox="0 0 14 14"
        fill="none"
        width={size}
        height={size}
        class={className}
        {...$$restProps}
>
    <circle cx="7" cy="7" r="5" stroke="currentColor" stroke-width="1.5"/>
</svg>

<script>
    export let size = 14;
    export let className = '';
</script>