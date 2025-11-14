<!-- lib/icons/MoreVertical.svelte -->
<svg
        viewBox="0 0 16 16"
        fill="none"
        width={size}
        height={size}
        class={className}
        {...$$restProps}
>
    <circle cx="8" cy="4" r="1" fill="currentColor"/>
    <circle cx="8" cy="8" r="1" fill="currentColor"/>
    <circle cx="8" cy="12" r="1" fill="currentColor"/>
</svg>

<script>
    export let size = 16;
    export let className = '';
</script>