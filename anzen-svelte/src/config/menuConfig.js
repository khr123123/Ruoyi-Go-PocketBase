import {ArrowDown, Check, Circle, Image, Menu} from "../lib/icons/index.js";

/**
 * 菜单列配置
 */
export const menuColumns = [{
    key: 'menuName', label: 'Menu Name', sortable: true, visible: true, icon: ArrowDown,class:"font-mono"
}, {
    key: 'icon', label: 'Icon', icon: Circle, sortable: false, class: 'text-xs font-mono text-gray-600',
}, {
    key: 'menuType', label: 'Type', icon: Image, sortable: true, render: (v) => {
        const types = {M: 'Directory', C: 'Menu', F: 'Button'};
        const colors = {M: 'blue', C: 'green', F: 'yellow'};
        return `<span class="px-2 py-1 text-xs rounded bg-${colors[v] || 'gray'}-100 text-${colors[v] || 'gray'}-700">${types[v] || v || '-'}</span>`;
    }
}, {
    key: 'orderNum',
    label: 'Sort',
    sortable: true,
    render: (v) => `<span class="px-2 py-1 text-xs rounded bg-gray-100 text-gray-700">${v ?? '-'}</span>`
}, {
    key: 'permission',
    label: 'Permission',
    sortable: false,
    class: 'text-xs font-mono',
    icon: Check,
    render: (v) => v ? `<code class="text-blue-600">${v}</code>` : '<span class="text-gray-400">-</span>'
}, {
    key: 'url',
    label: 'Route',
    sortable: false,
    class: 'text-xs',
    icon: Menu,
    render: (v) => v ? `<code class="text-green-600">${v}</code>` : '<span class="text-gray-400">-</span>'
}];

/**
 * 菜单类型选项
 */
export const menuTypeOptions = [{value: 'M', label: 'Directory (目录)'}, {
    value: 'C', label: 'Menu (菜单)'
}, {value: 'F', label: 'Button (按钮)'}];

/**
 * 默认菜单数据
 */
export const defaultMenuData = {
    menuName: '', menuType: 'M', parentId: '0', orderNum: 0, url: '', icon: '', permission: '', status: 'show'
};
