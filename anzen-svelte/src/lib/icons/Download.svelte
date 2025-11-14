<!-- lib/icons/Download.svelte -->
<svg
        viewBox="0 0 16 16"
        fill="none"
        width={size}
        height={size}
        class={className}
        {...$$restProps}
>
    <path d="M8 2V10M5 7L8 10L11 7" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
    <path d="M2 12V14H14V12" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
</svg>

<script>
    export let size = 16;
    export let className = '';
</script>