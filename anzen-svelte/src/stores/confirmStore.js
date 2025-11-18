import { writable } from 'svelte/store';

export const confirmStore = writable(null);

/**
 * JS 调用确认框
 * @param {string} message - 提示文字
 * @returns {Promise<boolean>} - 用户是否点击确认
 */
export function confirmDialog(message) {
    return new Promise((resolve) => {
        confirmStore.set({ message, resolve });
    });
}
