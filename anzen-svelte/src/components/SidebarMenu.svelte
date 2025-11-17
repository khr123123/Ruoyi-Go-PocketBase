<script>
    export let nodes = [];
    let openMap = new Map();

    const toggle = (id) => {
        openMap.set(id, !openMap.get(id));
        openMap = new Map(openMap);
    };

    const hasChildren = (node) =>
        node.children && node.children.length > 0;
</script>

<style>
    .menu-item {
        padding: 10px 12px;
        cursor: pointer;
        display: flex;
        justify-content: space-between;
    }
    .submenu {
        margin-left: 16px;
        border-left: 1px solid #eee;
    }
</style>

{#each nodes as node}
    <div class="menu-item" on:click={() => hasChildren(node) && toggle(node.id)}>
        <span>{node.menuName}</span>
        {#if hasChildren(node)}
            <span>{openMap.get(node.id) ? "▾" : "▸"}</span>
        {/if}
    </div>

    {#if hasChildren(node) && openMap.get(node.id)}
        <div class="submenu">
            <SidebarMenu nodes={node.children}/>
        </div>
    {/if}
{/each}
