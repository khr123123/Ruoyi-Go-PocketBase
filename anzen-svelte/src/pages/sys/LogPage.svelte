<script>
    import { onMount } from 'svelte';
    import { Refresh, Search, Settings } from '../../lib/icons/index.js';
    import LogChart from "../../components/LogChart.svelte";
    import LogTable from "../../components/LogTable.svelte";

    const LOG_LEVELS = [
        { value: -4, label: 'DEBUG' },
        { value: 0, label: 'INFO' },
        { value: 4, label: 'WARN' },
        { value: 8, label: 'ERROR' }
    ];

    let logs = [];
    let chartData = [];
    let totalLogs = 0;
    let currentPage = 1;
    let perPage = 50;
    let searchQuery = "";
    let selectedLevels = LOG_LEVELS.map(l => l.value);
    let includeSuperusers = false;
    let isLoading = false;

    /** -------------------------------
     * 生成虚拟日志
     * -------------------------------- */
    function generateFakeLogs(total = 500) {
        const messages = [
            "User login successful",
            "User logout",
            "Created new post",
            "Failed to fetch data",
            "Password updated",
            "Access denied",
            "Server error occurred",
            "Settings changed"
        ];
        const authUsers = ["alice", "bob", "charlie", "superuser"];
        const result = [];

        for (let i = 0; i < total; i++) {
            const level = LOG_LEVELS[Math.floor(Math.random() * LOG_LEVELS.length)].value;
            const message = messages[Math.floor(Math.random() * messages.length)];
            const auth = authUsers[Math.floor(Math.random() * authUsers.length)];
            const daysAgo = Math.floor(Math.random() * 30); // 最近 30 天
            const created = new Date();
            created.setDate(created.getDate() - daysAgo);
            created.setHours(Math.floor(Math.random() * 24), Math.floor(Math.random() * 60));

            result.push({
                id: i + 1,
                level,
                message,
                auth,
                created: created.toISOString()
            });
        }

        return result.sort((a, b) => new Date(b.created) - new Date(a.created));
    }

    let allFakeLogs = generateFakeLogs(500); // 全量虚拟日志

    /** -------------------------------
     * 构建过滤
     * -------------------------------- */
    function buildFilter(logList) {
        return logList.filter(log => {
            // 1. 过滤 level
            if (!selectedLevels.includes(log.level)) return false;

            // 2. 搜索 message
            if (searchQuery.trim() && !log.message.toLowerCase().includes(searchQuery.toLowerCase())) return false;

            // 3. 过滤 superuser
            if (!includeSuperusers && log.auth === "superuser") return false;

            return true;
        });
    }

    /** -------------------------------
     * 加载日志
     * -------------------------------- */
    function fetchLogs() {
        isLoading = true;
        setTimeout(() => { // 模拟异步
            const filtered = buildFilter(allFakeLogs);
            totalLogs = filtered.length;
            const start = (currentPage - 1) * perPage;
            const end = start + perPage;
            logs = filtered.slice(start, end);
            isLoading = false;
        }, 200);
    }

    /** -------------------------------
     * 加载图表数据
     * -------------------------------- */
    function fetchChartData() {
        const filtered = buildFilter(allFakeLogs);
        const grouped = {};
        filtered.forEach(log => {
            const date = new Date(log.created).toLocaleDateString("zh-CN", { month: "2-digit", day: "2-digit" });
            grouped[date] = (grouped[date] || 0) + 1;
        });
        chartData = Object.entries(grouped).map(([date, count]) => ({ date, count }));
    }

    /** -------------------------------
     * 事件处理
     * -------------------------------- */
    function toggleLevel(level) {
        if (selectedLevels.includes(level)) {
            selectedLevels = selectedLevels.filter(l => l !== level);
        } else {
            selectedLevels = [...selectedLevels, level];
        }
        fetchLogs();
        fetchChartData();
    }

    function refresh() {
        fetchLogs();
        fetchChartData();
    }

    function handleSearch() {
        currentPage = 1;
        fetchLogs();
        fetchChartData();
    }

    function handlePageChange(event) {
        currentPage = event.detail.page;
        fetchLogs();
    }

    onMount(() => {
        fetchLogs();
        fetchChartData();
    });
</script>


<!-- Page -->
<div class="min-h-screen bg-gray-50">
    <!-- Main -->
    <main class="max-w-[1600px] mx-auto">
        <!-- Log Levels -->
        <div class="mb-4 flex items-center gap-3">
            <span class="text-sm text-gray-600">Default log levels:</span>

            <div class="flex gap-2">
                {#each LOG_LEVELS as level}
                    <button
                          on:click={() => toggleLevel(level.value)}
                          class="px-3 py-1 text-xs rounded border transition
                               {selectedLevels.includes(level.value)
                                   ? 'bg-gray-200 border-gray-400 text-gray-800'
                                   : 'bg-gray-50 border-gray-300 text-gray-500'}"
                    >
                        {level.label}
                    </button>
                {/each}
            </div>

            <span class="text-sm text-gray-500 ml-auto">
                Found <strong class="text-gray-900">{totalLogs}</strong> logs
            </span>
        </div>

        <!-- Chart -->
        <div class="mb-6">
            <LogChart data={chartData}/>
        </div>

        <!-- Table -->
        <LogTable
              {logs}
              {totalLogs}
              {currentPage}
              {isLoading}
              on:pageChange={handlePageChange}
        />
    </main>
</div>
