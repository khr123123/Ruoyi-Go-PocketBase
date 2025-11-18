<!-- lib/icons/ArrowRight.svelte -->
<svg
        viewBox="0 0 16 16"
        fill="none"
        width={size}
        height={size}
        class={className}
        {...$$restProps}
>
  <path
          d="M6 3 L11 8 L6 13"
          stroke="currentColor"
          stroke-width="2"
          stroke-linecap="round"
          stroke-linejoin="round"
  />
</svg>

<script>
    export let size = 16;
    export let className = '';
</script>