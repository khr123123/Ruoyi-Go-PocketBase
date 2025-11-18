<script>
    import {createEventDispatcher} from 'svelte';
    import {fade, fly} from 'svelte/transition';

    export let show = false;
    export let title = '';
    export let width = '400px';
    export let position = 'right';

    const dispatch = createEventDispatcher();

    function close() {
        dispatch('close');
    }
</script>

{#if show}
    <!-- 背景遮罩 -->
    <div class="fixed inset-0 bg-gray-300/50 z-40" transition:fade on:click={close}></div>

    <!-- 抽屉容器 -->
    <div
          class="fixed top-0 bottom-0 bg-white shadow-2xl z-50 overflow-auto"
          style="width: {width}; {position}: 0;"
          transition:fly="{{ x: position==='right' ? 300 : -300, duration: 300 }}"
    >
        <div class="flex justify-between items-center p-4 border-b">
            <h3 class="text-lg font-semibold">{title}</h3>
            <button on:click={close} class="text-gray-400 hover:text-gray-600">&times;</button>
        </div>
        <div class="p-4">
            <slot></slot>
        </div>
    </div>
{/if}
