<!-- pages/Home.svelte -->
<script>
    import {ArrowDown, Check, Circle, Code, Download, Image, Menu, MoreVertical, Search} from '../lib/icons';

    const users = [
        {
            id: 'user0000000002',
            email: '1076731326@yahu.com',
            emailVisibility: false,
            verified: false,
            name: 'zhang',
            avatar: '🧑'
        },
        {
            id: 'user0000000001',
            email: '1076731326@qq.com',
            emailVisibility: false,
            verified: false,
            name: 'admin',
            avatar: '👨‍💼'
        }
    ];
</script>

<div class="page-container">
    <div class="page-header">
        <div class="search-bar">
            <Search size={16}/>
            <input
                    type="text"
                    placeholder='Search term or filter like created > "2022-01-01"...'
                    class="search-input"
            />
        </div>

        <div class="page-actions">
            <button class="action-btn secondary">
                <Code size={16}/>
                API Preview
            </button>
            <button class="action-btn primary">
                <span>+</span>
                New record
            </button>
        </div>
    </div>

    <div class="table-container">
        <table class="data-table">
            <thead>
            <tr>
                <th class="checkbox-col">
                    <input type="checkbox"/>
                </th>
                <th>
                    <div class="th-content">
                        <ArrowDown size={14}/>
                        id
                    </div>
                </th>
                <th>
                    <div class="th-content">
                        <Menu size={14}/>
                        email
                    </div>
                </th>
                <th>
                    <div class="th-content">
                        <Circle size={14}/>
                        emailVisibility
                    </div>
                </th>
                <th>
                    <div class="th-content">
                        <Check size={14}/>
                        verified
                    </div>
                </th>
                <th>
                    <div class="th-content">
                        <span>T</span>
                        name
                    </div>
                </th>
                <th>
                    <div class="th-content">
                        <Image size={14}/>
                        avatar
                    </div>
                </th>
                <th class="actions-col"></th>
            </tr>
            </thead>
            <tbody>
            {#each users as user}
                <tr>
                    <td class="checkbox-col">
                        <input type="checkbox"/>
                    </td>
                    <td>
                        <div class="id-cell">
                            <ArrowDown size={12}/>
                            {user.id}
                        </div>
                    </td>
                    <td>{user.email}</td>
                    <td>
                        <span class="badge {user.emailVisibility ? 'badge-success' : 'badge-default'}">
                            {user.emailVisibility ? 'true' : 'false'}
                        </span>
                    </td>
                    <td>
                        <span class="badge {user.verified ? 'badge-success' : 'badge-default'}">
                            {user.verified ? 'true' : 'false'}
                        </span>
                    </td>
                    <td>{user.name}</td>
                    <td>
                        <div class="avatar">{user.avatar}</div>
                    </td>
                    <td class="actions-col">
                        <button class="icon-btn">
                            <MoreVertical size={16}/>
                        </button>
                        <button class="icon-btn">→</button>
                    </td>
                </tr>
            {/each}
            </tbody>
        </table>
    </div>

    <div class="table-footer">
        <div class="footer-info">
            Total found: <strong>2</strong>
        </div>
        <div class="footer-actions">
            <button class="footer-btn" title="Export">
                <Download size={16}/>
            </button>
            <button class="footer-btn" title="Settings">⚙️</button>
            <button class="footer-btn" title="Help">❓</button>
            <button class="footer-btn" title="Info">ℹ️</button>
        </div>
    </div>
</div>

<style>
    .page-container {
        background: white;
        border-radius: 8px;
        box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
        overflow: hidden;
    }

    .page-header {
        padding: 16px 20px;
        border-bottom: 1px solid #e5e5e5;
        display: flex;
        align-items: center;
        justify-content: space-between;
        gap: 16px;
    }

    .search-bar {
        flex: 1;
        max-width: 600px;
        position: relative;
        display: flex;
        align-items: center;
        gap: 8px;
        padding: 8px 12px;
        border: 1px solid #e5e5e5;
        border-radius: 4px;
        background-color: #f7f7f7;
    }

    .search-bar svg {
        color: #999;
    }

    .search-bar .search-input {
        flex: 1;
        border: none;
        background: transparent;
        font-size: 13px;
        outline: none;
        color: #2b3034;
    }

    .page-actions {
        display: flex;
        gap: 8px;
    }

    .action-btn {
        padding: 8px 16px;
        border-radius: 4px;
        font-size: 13px;
        font-weight: 500;
        cursor: pointer;
        display: flex;
        align-items: center;
        gap: 6px;
        transition: all 0.2s;
    }

    .action-btn.secondary {
        border: 1px solid #e5e5e5;
        background: white;
        color: #2b3034;
    }

    .action-btn.secondary:hover {
        border-color: #2b3034;
        background: #f7f7f7;
    }

    .action-btn.primary {
        border: none;
        background: #2b3034;
        color: white;
    }

    .action-btn.primary:hover {
        background: #1a1d1f;
    }

    .table-container {
        overflow-x: auto;
    }

    .data-table {
        width: 100%;
        border-collapse: collapse;
        font-size: 13px;
    }

    .data-table th {
        background: #f7f7f7;
        padding: 12px 16px;
        text-align: left;
        font-weight: 500;
        color: #666;
        border-bottom: 1px solid #e5e5e5;
    }

    .th-content {
        display: flex;
        align-items: center;
        gap: 6px;
    }

    .data-table td {
        padding: 12px 16px;
        border-bottom: 1px solid #f0f0f0;
        color: #2b3034;
    }

    .data-table tr:hover {
        background-color: #fafafa;
    }

    .checkbox-col {
        width: 40px;
    }

    .actions-col {
        width: 80px;
        text-align: right;
    }

    .id-cell {
        display: flex;
        align-items: center;
        gap: 6px;
        color: #666;
        font-family: monospace;
        font-size: 12px;
    }

    .badge {
        display: inline-block;
        padding: 3px 8px;
        border-radius: 12px;
        font-size: 11px;
        font-weight: 500;
    }

    .badge-default {
        background: #f0f0f0;
        color: #666;
    }

    .badge-success {
        background: #e8f5e9;
        color: #2e7d32;
    }

    .avatar {
        width: 32px;
        height: 32px;
        border-radius: 50%;
        background: #f0f4f8;
        display: flex;
        align-items: center;
        justify-content: center;
        font-size: 18px;
    }

    .icon-btn {
        padding: 4px 8px;
        border: none;
        background: transparent;
        color: #999;
        cursor: pointer;
        border-radius: 4px;
        transition: all 0.2s;
    }

    .icon-btn:hover {
        background: #f0f0f0;
        color: #2b3034;
    }

    .table-footer {
        padding: 12px 20px;
        border-top: 1px solid #e5e5e5;
        display: flex;
        align-items: center;
        justify-content: space-between;
        background: #fafafa;
    }

    .footer-info {
        font-size: 13px;
        color: #666;
    }

    .footer-actions {
        display: flex;
        gap: 8px;
    }

    .footer-btn {
        padding: 6px;
        border: none;
        background: transparent;
        color: #666;
        cursor: pointer;
        border-radius: 4px;
        transition: all 0.2s;
        display: flex;
        align-items: center;
    }

    .footer-btn:hover {
        background: white;
        color: #2b3034;
    }
</style>