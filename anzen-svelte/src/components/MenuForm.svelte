<script>
    import { createEventDispatcher } from 'svelte';
    import { menuTypeOptions } from '../config/menuConfig.js';

    export let menu = {};
    export let parentOptions = [];
    export let isEdit = false;

    const dispatch = createEventDispatcher();

    function handleSave() {
        if (!menu.menuName || !menu.menuType) {
            dispatch('error', 'Please fill in required fields');
            return;
        }
        dispatch('save', menu);
    }

    function handleCancel() {
        dispatch('cancel');
    }
</script>

<!-- 背景遮罩 -->
<div
      class="fixed inset-0 bg-black/20 z-40 animate-fadeIn"
      on:click={handleCancel}
></div>

<!-- Drawer -->
<div
      class="fixed top-0 right-0 h-full w-[480px] bg-white shadow-2xl z-50 flex flex-col animate-slideIn"
>

    <!-- Header -->
    <div class="px-6 py-4 border-b border-gray-200 flex items-center justify-between">
        <h3 class="text-base font-semibold text-gray-900">
            {isEdit ? 'Edit menu record' : 'New menu record'}
        </h3>

        <button on:click={handleCancel} class="text-gray-400 hover:text-gray-600 transition-colors">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                      d="M6 18L18 6M6 6l12 12" />
            </svg>
        </button>
    </div>

    <!-- Content -->
    <div class="flex-1 overflow-y-auto px-6 py-5 space-y-5">

        <div>
            <label class="block text-xs font-medium text-gray-600 mb-1.5">
                menuName <span class="text-red-500">*</span>
            </label>
            <input
                  class="w-full border border-gray-300 rounded px-3 py-2 text-sm
                       bg-white text-gray-900 placeholder-gray-400
                       focus:border-gray-500 focus:ring-1 focus:ring-gray-500 outline-none
                       transition-colors"
                  bind:value={menu.menuName}
                  placeholder="e.g., Dashboard, Users"
            />
        </div>

        <div>
            <label class="block text-xs font-medium text-gray-600 mb-1.5">
                menuType <span class="text-red-500">*</span>
            </label>
            <select
                  bind:value={menu.menuType}
                  class="w-full border border-gray-300 rounded px-3 py-2 text-sm
                       bg-white text-gray-900 focus:border-gray-500
                       focus:ring-1 focus:ring-gray-500 outline-none"
            >
                <option value="">Select type...</option>
                {#each menuTypeOptions as option}
                    <option value={option.value}>{option.label}</option>
                {/each}
            </select>
        </div>

        <div>
            <label class="block text-xs font-medium text-gray-600 mb-1.5">parentId</label>
            <select
                  bind:value={menu.parentId}
                  class="w-full border border-gray-300 rounded px-3 py-2 text-sm
                       bg-white text-gray-900 focus:border-gray-500
                       focus:ring-1 focus:ring-gray-500 outline-none"
            >
                <option value="">-- Root --</option>
                {#each parentOptions as option}
                    <option value={option.id} disabled={menu.id === option.id}>
                        {option.menuName}
                    </option>
                {/each}
            </select>
        </div>

        <div>
            <label class="block text-xs font-medium text-gray-600 mb-1.5">icon</label>
            <input
                  class="w-full border border-gray-300 rounded px-3 py-2 text-sm
                       text-gray-900 bg-white placeholder-gray-400
                       focus:border-gray-500 focus:ring-1 focus:ring-gray-500 outline-none"
                  bind:value={menu.icon}
                  placeholder="e.g., fa fa-home"
            />
        </div>

        {#if menu.menuType !== 'F'}
            <div>
                <label class="block text-xs font-medium text-gray-600 mb-1.5">url</label>
                <input
                      class="w-full border border-gray-300 rounded px-3 py-2 text-sm
                           text-gray-900 bg-white placeholder-gray-400
                           focus:border-gray-500 focus:ring-1 focus:ring-gray-500 outline-none"
                      bind:value={menu.url}
                      placeholder="/path/to/page"
                />
            </div>
        {/if}

        <div>
            <label class="block text-xs font-medium text-gray-600 mb-1.5">permission</label>
            <input
                  class="w-full border border-gray-300 rounded px-3 py-2 text-sm
                       text-gray-900 bg-white placeholder-gray-400
                       focus:border-gray-500 focus:ring-1 focus:ring-gray-500 outline-none"
                  bind:value={menu.permission}
                  placeholder="e.g., sys:menu:view"
            />
        </div>

        <div>
            <label class="block text-xs font-medium text-gray-600 mb-1.5">orderNum</label>
            <input
                  class="w-full border border-gray-300 rounded px-3 py-2 text-sm
                       text-gray-900 bg-white focus:border-gray-500 focus:ring-1
                       focus:ring-gray-500 outline-none"
                  type="number"
                  bind:value={menu.orderNum}
                  placeholder="0"
            />
        </div>

    </div>

    <!-- Footer -->
    <div class="px-6 py-4 border-t border-gray-200 flex justify-end gap-2">

        <button
              on:click={handleCancel}
              class="px-4 py-2 border border-gray-300 rounded text-sm
                   text-gray-700 bg-white hover:bg-gray-50 transition-colors"
        >
            Cancel
        </button>

        <button
              on:click={handleSave}
              class="px-4 py-2 bg-gray-800 text-white rounded text-sm
                   hover:bg-gray-700 transition-colors"
        >
            Save changes
        </button>

    </div>
</div>

<style>
    /* 仅保留动画（Tailwind 无法写） */
    @keyframes fadeIn {
        from { opacity: 0; }
        to   { opacity: 1; }
    }

    @keyframes slideIn {
        from { transform: translateX(100%); }
        to   { transform: translateX(0); }
    }

    .animate-fadeIn {
        animation: fadeIn 0.2s ease;
    }

    .animate-slideIn {
        animation: slideIn .25s cubic-bezier(0.25, 0.1, 0.25, 1) forwards;
    }
</style>
