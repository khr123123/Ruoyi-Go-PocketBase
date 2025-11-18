<script>
    import {confirmStore} from '../stores/confirmStore.js';
    import {onDestroy} from 'svelte';

    let current = null;

    const unsubscribe = confirmStore.subscribe(value => {
        current = value;
    });

    function handleConfirm() {
        current?.resolve(true);
        confirmStore.set(null);
    }

    function handleCancel() {
        current?.resolve(false);
        confirmStore.set(null);
    }

    onDestroy(() => unsubscribe());
</script>

{#if current}
    <div class="fixed inset-x-0 top-20 z-50 flex justify-center pointer-events-none">
        <div class="bg-white rounded-lg p-6 transform transition-all duration-200 scale-95 opacity-0 animate-scale-in pointer-events-auto
                shadow-xl shadow-gray-400/30 border border-gray-200 max-w-sm w-full">
            <!-- 图标 + 文本 -->
            <div class="flex items-center gap-3">
                <!-- 红色感叹号 -->
                <svg class="w-6 h-6 text-red-600 flex-shrink-0" fill="currentColor" viewBox="0 0 20 20">
                    <path fill-rule="evenodd"
                          d="M18 10c0 4.418-3.582 8-8 8s-8-3.582-8-8 3.582-8 8-8 8 3.582 8 8zm-8-4a1 1 0 00-1 1v4a1 1 0 002 0V7a1 1 0 00-1-1zm0 8a1.5 1.5 0 100-3 1.5 1.5 0 000 3z"
                          clip-rule="evenodd"/>
                </svg>
                <p class="text-gray-800 text-sm flex-1">{current.message}</p>
            </div>

            <!-- 按钮 -->
            <div class="mt-6 flex justify-end gap-3">
                <button
                      on:click={handleCancel}
                      class="px-4 py-2 rounded bg-gray-200 hover:bg-gray-300 transition text-gray-700">
                    Cancel
                </button>
                <button
                      on:click={handleConfirm}
                      class="px-4 py-2 rounded bg-red-600 hover:bg-red-700 transition text-white">
                    Delete
                </button>
            </div>
        </div>
    </div>
{/if}

<style>
    @keyframes scale-in {
        0% {
            transform: scale(0.95);
            opacity: 0;
        }
        100% {
            transform: scale(1);
            opacity: 1;
        }
    }

    .animate-scale-in {
        animation: scale-in 0.15s ease-out forwards;
    }
</style>
