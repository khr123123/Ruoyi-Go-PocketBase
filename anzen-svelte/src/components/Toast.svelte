<script>
    import {onMount} from "svelte";

    export let message = "";
    export let type = "success"; // success | error | warning
    export let onClose = () => {
    };

    let toastEl; // Toast DOM 元素
    let left = 0; // 动态计算 left

    const typeConfig = {
        success: {bg: "bg-green-50 border-green-300 text-green-800", icon: "✓"},
        error: {bg: "bg-red-50 border-red-300 text-red-800", icon: "✕"},
        warning: {bg: "bg-yellow-50 border-yellow-300 text-yellow-800", icon: "!"},
    };

    const config = typeConfig[type] || typeConfig.success;

    // 自动消失
    const timer = setTimeout(() => onClose(), 3000);

    onMount(() => {
        if (toastEl) {
            // 水平居中计算
            left = window.innerWidth / 2 - toastEl.offsetWidth / 2;
        }
    });
</script>

<div
        bind:this={toastEl}
        class="fixed top-1/16 z-50 animate-drop-in"
        style="left: {left}px;"
>
    <div class={`flex items-center gap-3 px-4 py-3 rounded-lg border shadow ${config.bg}`}>
        <span class="text-xl font-bold">{config.icon}</span>
        <span class="font-medium">{message}</span>
        <button class="ml-2 opacity-60 hover:opacity-100" on:click={onClose}>×</button>
    </div>
</div>

<style>
    @keyframes drop-in {
        from {
            opacity: 0;
            transform: translateY(-30px); /* 起始位置稍微更高一点 */
        }
        to {
            opacity: 1;
            transform: translateY(0);
        }
    }

    /* 缓慢顺滑的动画 */
    .animate-drop-in {
        animation: drop-in 0.6s cubic-bezier(0.25, 0.8, 0.5, 1) forwards;
    }
</style>
