<script>
    import {createEventDispatcher} from 'svelte';
    import {menuTypeOptions} from '../config/menuConfig.js';

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
<div class="drawer-mask" on:click={handleCancel}></div>

<!-- Drawer -->
<div class="drawer-panel">
    <!-- Header -->
    <div class="drawer-header">
        <h3 class="title">{isEdit ? 'Edit Menu' : 'Add Menu'}</h3>
        <button on:click={handleCancel} class="close-btn">✕</button>
    </div>

    <!-- Content -->
    <div class="drawer-body">

        <div>
            <label class="pb-label">Menu Name *</label>
            <input class="pb-input" bind:value={menu.menuName}/>
        </div>

        <div>
            <label class="pb-label">Menu Type *</label>
            <select bind:value={menu.menuType} class="pb-input">
                {#each menuTypeOptions as option}
                    <option value={option.value}>{option.label}</option>
                {/each}
            </select>
        </div>

        <div>
            <label class="pb-label">Parent Menu</label>
            <select bind:value={menu.parentId} class="pb-input">
                {#each parentOptions as option}
                    <option value={option.id} disabled={menu.id === option.id}>
                        {option.menuName}
                    </option>
                {/each}
            </select>
        </div>

        <div>
            <label class="pb-label">Icon</label>
            <input class="pb-input" bind:value={menu.icon}/>
        </div>

        {#if menu.menuType !== 'F'}
            <div>
                <label class="pb-label">Route Path</label>
                <input class="pb-input" bind:value={menu.url}/>
            </div>
        {/if}

        <div>
            <label class="pb-label">Permission</label>
            <input class="pb-input" bind:value={menu.permission}/>
        </div>

        <div>
            <label class="pb-label">Sort Order</label>
            <input class="pb-input" type="number" bind:value={menu.orderNum}/>
        </div>

    </div>

    <!-- Footer -->
    <div class="drawer-footer">
        <button on:click={handleCancel} class="btn-cancel">Cancel</button>
        <button on:click={handleSave} class="btn-primary">Save</button>
    </div>
</div>

<style lang="postcss">
    @reference "tailwindcss";

    /* ------------------------------
       遮罩：淡入淡出
    ------------------------------ */
    .drawer-mask {
        @apply fixed inset-0 bg-black/30 z-40;
        animation: fadeIn 0.25s ease forwards;
    }

    @keyframes fadeIn {
        from {
            opacity: 0;
        }
        to {
            opacity: 1;
        }
    }

    /* ------------------------------
       Drawer 右侧滑入
       PocketBase 风格（直角 + 阴影）
    ------------------------------ */
    .drawer-panel {
        @apply fixed top-0 right-0 h-full w-[420px] bg-white
        border-l border-gray-300 shadow-xl z-50
        flex flex-col;

        transform: translateX(100%);
        animation: slideIn 0.28s cubic-bezier(0.25, 0.1, 0.25, 1) forwards;
    }

    @keyframes slideIn {
        from {
            transform: translateX(100%);
        }
        to {
            transform: translateX(0%);
        }
    }

    /* ------------------------------ */
    .drawer-header {
        @apply px-5 py-4 border-b border-gray-300 flex items-center justify-between;
    }

    .title {
        @apply text-base font-semibold text-gray-800;
    }

    .close-btn {
        @apply text-gray-500 hover:text-gray-700 transition;
    }

    /* ------------------------------ */
    .drawer-body {
        @apply flex-1 overflow-y-auto px-5 py-5 space-y-5;
    }

    /* label */
    .pb-label {
        @apply text-sm font-medium text-gray-700;
    }

    /* PocketBase 风输入框：更硬朗 */
    .pb-input {
        @apply w-full border border-gray-300 rounded-md px-3 py-2.5
        bg-white text-gray-800
        focus:border-gray-700 focus:ring-0
        transition;
    }

    /* ------------------------------ */
    .drawer-footer {
        @apply px-5 py-4 border-t border-gray-300 flex justify-end gap-2;
    }

    /* Buttons — 冷淡风 PocketBase */
    .btn-cancel {
        @apply px-4 py-2 border border-gray-300 rounded-md
        text-gray-700 bg-white hover:bg-gray-100 transition;
    }

    .btn-primary {
        @apply px-4 py-2 bg-gray-800 text-white rounded-md
        hover:bg-gray-700 transition;
    }
</style>
