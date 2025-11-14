<!-- Layout.svelte -->
<script>
    import Router, {location, push} from "svelte-spa-router";
    import {wrap} from 'svelte-spa-router/wrap';
    import Home from './pages/Home.svelte';
    import About from './pages/About.svelte';

    // 导入图标
    import {Chart, IdCard, Logo, Refresh, Settings, SettingsAlt, Users} from './lib/icons';

    let sidebarOpen = true;
    let searchQuery = '';

    const layoutRoutes = {
        '/': wrap({component: Home}),
        '/about': wrap({component: About})
    }

    function handleLogout() {
        push('/login');
    }

    function isActive(path) {
        return $location === path;
    }
</script>

<div class="layout">
    <!-- Sidebar -->
    <aside class="sidebar" class:collapsed={!sidebarOpen}>
        <div class="sidebar-header">
            <div class="logo-container">
                <Logo size={32}/>
                <span class="tracking-wide">RUOYI-PB</span>
            </div>
            <input
                    type="text"
                    placeholder="Search collections..."
                    class="search-input"
                    bind:value={searchQuery}
            />
        </div>

        <nav class="sidebar-nav">
            <div class="nav-section">
                <a href="#/" class="nav-item" class:active={isActive('/')}>
                    <span class="nav-icon">
                        <Users size={20}/>
                    </span>
                    <span class="nav-text">users</span>
                </a>

                <a href="#/about" class="nav-item" class:active={isActive('/about')}>
                    <span class="nav-icon">
                        <Chart size={20}/>
                    </span>
                    <span class="nav-text">sys_menu</span>
                </a>

                <a href="#/services" class="nav-item" class:active={isActive('/services')}>
                    <span class="nav-icon">
                        <IdCard size={20}/>
                    </span>
                    <span class="nav-text">sys_role</span>
                </a>
            </div>

            <div class="nav-section">
                <div class="section-label">System</div>
                <a href="#/contact" class="nav-item" class:active={isActive('/contact')}>
                    <span class="nav-icon">
                        <Settings size={20}/>
                    </span>
                    <span class="nav-text">Settings</span>
                </a>
            </div>
        </nav>

        <div class="sidebar-footer">
            <button class="new-collection-btn">
                <span>+</span>
                <span>New collection</span>
            </button>
        </div>
    </aside>

    <!-- Main Content -->
    <main class="main-content">
        <header class="content-header">
            <div class="breadcrumb">
                <span class="breadcrumb-item">Collections</span>
                <span class="breadcrumb-separator">/</span>
                <span class="breadcrumb-item active">
                    {$location === '/' ? 'users' : $location.slice(1)}
                </span>
            </div>

            <div class="header-actions">
                <button class="header-btn" title="Settings">
                    <SettingsAlt size={20}/>
                </button>
                <button class="header-btn" title="Refresh">
                    <Refresh size={20}/>
                </button>
                <button class="logout-btn" on:click={handleLogout}>
                    Logout
                </button>
            </div>
        </header>

        <div class="content-body">
            <Router routes={layoutRoutes}/>
        </div>
    </main>
</div>


<style>
    .layout {
        display: flex;
        height: 100vh;
        background-color: #f7f7f7;
    }

    /* Sidebar Styles */
    .sidebar {
        width: 200px;
        background: white;
        border-right: 1px solid #e5e5e5;
        display: flex;
        flex-direction: column;
        transition: width 0.3s;
    }

    .sidebar.collapsed {
        width: 60px;
    }

    .sidebar-header {
        padding: 20px 16px;
        border-bottom: 1px solid #e5e5e5;
    }

    .logo-container {
        display: flex;
        align-items: center;
        justify-content: center;
        margin-bottom: 16px;
        color: #2b3034;
    }

    .search-input {
        width: 100%;
        padding: 8px 12px;
        border: 1px solid #e5e5e5;
        border-radius: 4px;
        font-size: 13px;
        background-color: #f7f7f7;
    }

    .search-input:focus {
        outline: none;
        border-color: #5a9fd4;
        background-color: white;
    }

    .sidebar-nav {
        flex: 1;
        overflow-y: auto;
        padding: 12px 8px;
    }

    .nav-section {
        margin-bottom: 20px;
    }

    .section-label {
        font-size: 11px;
        color: #999;
        text-transform: uppercase;
        letter-spacing: 0.5px;
        padding: 8px 12px;
        font-weight: 600;
    }

    .nav-item {
        display: flex;
        align-items: center;
        gap: 12px;
        padding: 10px 12px;
        border-radius: 6px;
        color: #2b3034;
        text-decoration: none;
        font-size: 14px;
        transition: background-color 0.2s;
        margin-bottom: 2px;
    }

    .nav-item:hover {
        background-color: #f0f4f8;
    }

    .nav-item.active {
        background-color: #e8f0f7;
        color: #2b3034;
        font-weight: 500;
    }

    .nav-icon {
        font-size: 16px;
        width: 20px;
        display: flex;
        align-items: center;
        justify-content: center;
    }

    .sidebar-footer {
        padding: 16px;
        border-top: 1px solid #e5e5e5;
    }

    .new-collection-btn {
        width: 100%;
        padding: 10px 16px;
        border: 1px solid #e5e5e5;
        border-radius: 6px;
        background: white;
        color: #2b3034;
        font-size: 14px;
        font-weight: 500;
        cursor: pointer;
        display: flex;
        align-items: center;
        justify-content: center;
        gap: 8px;
        transition: all 0.2s;
    }

    .new-collection-btn:hover {
        border-color: #2b3034;
        background-color: #f7f7f7;
    }

    /* Main Content Styles */
    .main-content {
        flex: 1;
        display: flex;
        flex-direction: column;
        overflow: hidden;
    }

    .content-header {
        background: white;
        border-bottom: 1px solid #e5e5e5;
        padding: 16px 24px;
        display: flex;
        align-items: center;
        justify-content: space-between;
    }

    .breadcrumb {
        display: flex;
        align-items: center;
        gap: 8px;
        font-size: 14px;
    }

    .breadcrumb-item {
        color: #666;
    }

    .breadcrumb-item.active {
        color: #2b3034;
        font-weight: 500;
    }

    .breadcrumb-separator {
        color: #ccc;
    }

    .header-actions {
        display: flex;
        align-items: center;
        gap: 12px;
    }

    .header-btn {
        padding: 8px;
        border: 1px solid #e5e5e5;
        border-radius: 4px;
        background: white;
        color: #666;
        cursor: pointer;
        display: flex;
        align-items: center;
        justify-content: center;
        transition: all 0.2s;
    }

    .header-btn:hover {
        border-color: #2b3034;
        color: #2b3034;
    }

    .logout-btn {
        padding: 8px 16px;
        border: 1px solid #e5e5e5;
        border-radius: 4px;
        background: white;
        color: #2b3034;
        font-size: 13px;
        font-weight: 500;
        cursor: pointer;
        transition: all 0.2s;
    }

    .logout-btn:hover {
        background-color: #2b3034;
        color: white;
    }

    .content-body {
        flex: 1;
        overflow-y: auto;
        padding: 24px;
    }

</style>