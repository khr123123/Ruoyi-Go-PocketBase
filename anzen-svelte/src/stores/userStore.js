import {writable} from "svelte/store";

// key 名称
const STORAGE_KEY = "ruoyi-pb-user";

// 从 localStorage 读取初始值
const storedUser = localStorage.getItem(STORAGE_KEY);
const defaultUser = storedUser ? JSON.parse(storedUser) : {
    id: null, name: "", email: "", avatar: "", role: [], token: ""
};

// 创建全局 store
export const user = writable(defaultUser);

// 订阅 store，变化时同步到 localStorage
user.subscribe((value) => {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(value));
});

// 登录/更新用户信息
export function setUser(data) {
    if (!data) return;
    const {record, token} = data;
    user.set({
        id: record.id,
        name: record.name,
        email: record.email,
        avatar: record.avatar,
        role: record.role || [],
        token: token || "",
    });
}

// 清空用户信息（登出）
export function clearUser() {
    user.set({id: null, name: "", email: "", avatar: "", role: [], token: ""});
    localStorage.removeItem(STORAGE_KEY); // 清除 localStorage
}
