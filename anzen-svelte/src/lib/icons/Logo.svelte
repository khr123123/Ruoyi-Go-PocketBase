<!-- lib/icons/Logo.svelte -->
<svg
        width={size}
        height={size}
        viewBox="0 0 48 48"
        fill="none"
        class={className}
        {...$$restProps}
>
    <rect x="8" y="8" width="32" height="32" rx="4" stroke="currentColor" stroke-width="2"/>
    <path d="M16 20 H24 M16 24 H24 M16 28 H24" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
    <rect x="28" y="20" width="4" height="8" fill="currentColor"/>
</svg>

<script>
    export let size = 32;
    export let className = '';
</script>