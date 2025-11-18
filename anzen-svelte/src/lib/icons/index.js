export {default as Logo} from './Logo.svelte';
export {default as Users} from './Users.svelte';
export {default as Chart} from './Chart.svelte';
export {default as IdCard} from './IdCard.svelte';
export {default as Settings} from './Settings.svelte';
export {default as SettingsAlt} from './SettingsAlt.svelte';
export {default as Refresh} from './Refresh.svelte';
export {default as Search} from './Search.svelte';
export {default as ArrowRight} from './ArrowRight.svelte';
export {default as ArrowDown} from './ArrowDown.svelte';
export {default as Menu} from './Menu.svelte';
export {default as Circle} from './Circle.svelte';
export {default as Check} from './Check.svelte';
export {default as Image} from './Image.svelte';
export {default as MoreVertical} from './MoreVertical.svelte';
export {default as Download} from './Download.svelte';
export {default as Code} from './Code.svelte';
export {default as Bars4} from './Bars4.svelte';
export {default as ChevronDown} from './ChevronDown.svelte';
export {default as ChevronUpDown} from './ChevronUpDown.svelte';
export {default as ArrowUp} from './ArrowUp.svelte';
export {default as Home} from './Home.svelte';
export {default as Identifi} from './Identifi.svelte';
export {default as ShieldCheck} from './ShieldCheck.svelte';
export {default as Github} from './Github.svelte';
export {default as Bug} from './Bug.svelte';
export {default as Expanded} from './Expanded.svelte';
export {default as NoExpanded} from './NoExpanded.svelte';

// src/lib/icons/index.js

// ... 其他图标 ...

export function Upload(props) {
    return `<svg ${props.class ? `class="${props.class}"` : ''} fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12"/>
    </svg>`;
}

export function X(props) {
    return `<svg ${props.class ? `class="${props.class}"` : ''} fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
    </svg>`;
}
