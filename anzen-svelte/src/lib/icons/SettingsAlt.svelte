<!-- lib/icons/SettingsAlt.svelte -->
<svg
        viewBox="0 0 20 20"
        fill="none"
        width={size}
        height={size}
        class={className}
        {...$$restProps}
>
    <path
            d="M10 12.5C11.3807 12.5 12.5 11.3807 12.5 10C12.5 8.61929 11.3807 7.5 10 7.5C8.61929 7.5 7.5 8.61929 7.5 10C7.5 11.3807 8.61929 12.5 10 12.5Z"
            stroke="currentColor"
            stroke-width="1.5"
    />
    <path
            d="M10 2.5V4.5M10 15.5V17.5M17.5 10H15.5M4.5 10H2.5M15.3033 4.69672L13.8891 6.11094M6.11094 13.8891L4.69672 15.3033M15.3033 15.3033L13.8891 13.8891M6.11094 6.11094L4.69672 4.69672"
            stroke="currentColor"
            stroke-width="1.5"
            stroke-linecap="round"
    />
</svg>

<script>
    export let size = 20;
    export let className = '';
</script>