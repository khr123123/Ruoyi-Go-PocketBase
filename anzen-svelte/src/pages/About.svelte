<!-- pages/About.svelte -->
<div class="page-content">
    <div class="content-card">
        <h2>sys_menu 集合</h2>
        <p>这是系统菜单管理页面，用于配置应用程序的菜单结构。</p>

        <div class="info-grid">
            <div class="info-item">
                <div class="info-label">记录总数</div>
                <div class="info-value">12</div>
            </div>
            <div class="info-item">
                <div class="info-label">最后更新</div>
                <div class="info-value">2024-01-15</div>
            </div>
            <div class="info-item">
                <div class="info-label">状态</div>
                <div class="info-value">
                    <span class="status-badge active">活跃</span>
                </div>
            </div>
        </div>
    </div>
</div>

<style>
    .page-content {
        max-width: 1200px;
    }

    .content-card {
        background: white;
        border-radius: 8px;
        padding: 32px;
        box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
    }

    .content-card h2 {
        margin: 0 0 16px 0;
        color: #2b3034;
        font-size: 24px;
    }

    .content-card p {
        color: #666;
        line-height: 1.6;
        margin-bottom: 32px;
    }

    .info-grid {
        display: grid;
        grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
        gap: 24px;
    }

    .info-item {
        padding: 20px;
        background: #f7f7f7;
        border-radius: 6px;
    }

    .info-label {
        font-size: 12px;
        color: #999;
        text-transform: uppercase;
        letter-spacing: 0.5px;
        margin-bottom: 8px;
    }

    .info-value {
        font-size: 24px;
        font-weight: 600;
        color: #2b3034;
    }

    .status-badge {
        display: inline-block;
        padding: 6px 12px;
        border-radius: 12px;
        font-size: 13px;
        font-weight: 500;
    }

    .status-badge.active {
        background: #e8f5e9;
        color: #2e7d32;
    }

    .page-content {
        max-width: 1200px;
    }

    .content-card {
        background: white;
        border-radius: 8px;
        padding: 32px;
        box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
    }

    .content-card h2 {
        margin: 0 0 16px 0;
        color: #2b3034;
        font-size: 24px;
    }

    .content-card p {
        color: #666;
        line-height: 1.6;
        margin-bottom: 32px;
    }

    .roles-list {
        display: grid;
        grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
        gap: 20px;
    }

    .role-card {
        padding: 24px;
        border: 1px solid #e5e5e5;
        border-radius: 8px;
        transition: all 0.2s;
    }

    .role-card:hover {
        border-color: #2b3034;
        box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
    }

    .role-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 12px;
    }

    .role-header h3 {
        margin: 0;
        color: #2b3034;
        font-size: 18px;
    }

    .role-badge {
        padding: 4px 10px;
        border-radius: 12px;
        font-size: 11px;
        font-weight: 600;
        text-transform: uppercase;
        letter-spacing: 0.5px;
    }

    .role-badge.admin {
        background: #fce4ec;
        color: #c2185b;
    }

    .role-badge.editor {
        background: #e3f2fd;
        color: #1976d2;
    }

    .role-badge.viewer {
        background: #f3e5f5;
        color: #7b1fa2;
    }

    .role-description {
        color: #666;
        font-size: 14px;
        margin-bottom: 16px;
    }

    .role-stats {
        display: flex;
        gap: 16px;
        font-size: 13px;
        color: #999;
    }
</style>
<!-- pages/Services.svelte -->
<div class="page-content">
    <div class="content-card">
        <h2>sys_role 集合</h2>
        <p>系统角色权限管理，控制用户访问权限和操作范围。</p>

        <div class="roles-list">
            <div class="role-card">
                <div class="role-header">
                    <h3>管理员</h3>
                    <span class="role-badge admin">Admin</span>
                </div>
                <p class="role-description">拥有系统的完全访问权限</p>
                <div class="role-stats">
                    <span>👥 5 用户</span>
                    <span>🔒 全部权限</span>
                </div>
            </div>

            <div class="role-card">
                <div class="role-header">
                    <h3>编辑者</h3>
                    <span class="role-badge editor">Editor</span>
                </div>
                <p class="role-description">可以创建和编辑内容</p>
                <div class="role-stats">
                    <span>👥 12 用户</span>
                    <span>📝 读写权限</span>
                </div>
            </div>

            <div class="role-card">
                <div class="role-header">
                    <h3>访客</h3>
                    <span class="role-badge viewer">Viewer</span>
                </div>
                <p class="role-description">只能查看内容</p>
                <div class="role-stats">
                    <span>👥 28 用户</span>
                    <span>👁️ 只读权限</span>
                </div>
            </div>
        </div>
    </div>
</div>

