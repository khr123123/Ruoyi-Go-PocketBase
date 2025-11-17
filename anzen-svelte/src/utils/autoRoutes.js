const modules = import.meta.glob('../../pages/**/*.svelte');

/**
 * 根据 URL 获取对应的组件异步加载函数
 * 示例: /system/user -> ../../pages/system/user.svelte
 */
export function getComponentAsync(path) {
    // 规范化路径
    const normalizedPath = path.startsWith('/') ? path : '/' + path;
    // 构造相对路径
    const fullPath = `/src/pages${normalizedPath}.svelte`;
    // 查找匹配模块
    const importer = modules[fullPath];
    if (!importer) {
        console.warn(`No page found for route: ${path}`);
        return null;
    }
    return async () => {
        const {default: component} = await importer();
        return component;
    };
}