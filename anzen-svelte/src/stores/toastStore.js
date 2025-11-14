import {writable} from 'svelte/store';

export const toast = writable(null);

// 辅助方法，方便调用
export function showToast(message, type = 'info', duration = 3000) {
    toast.set({message, type});
    // 自动消失
    setTimeout(() => {
        toast.set(null);
    }, duration);
}
