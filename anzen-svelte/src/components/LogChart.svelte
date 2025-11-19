<!-- src/components/LogChart.svelte -->
<script>
    import {onMount, afterUpdate} from 'svelte';

    export let data = [];

    let canvas;
    let ctx;

    function drawChart() {
        if (!canvas || !ctx || data.length === 0) return;

        const width = canvas.width;
        const height = canvas.height;
        const padding = {top: 20, right: 20, bottom: 30, left: 40};

        // 清空画布
        ctx.clearRect(0, 0, width, height);

        // 计算数据范围
        const maxCount = Math.max(...data.map(d => d.count), 1);
        const minCount = 0;

        // 绘制网格线
        ctx.strokeStyle = '#f0f0f0';
        ctx.lineWidth = 1;
        for (let i = 0; i <= 5; i++) {
            const y = padding.top + (height - padding.top - padding.bottom) * i / 5;
            ctx.beginPath();
            ctx.moveTo(padding.left, y);
            ctx.lineTo(width - padding.right, y);
            ctx.stroke();
        }

        // 绘制折线
        if (data.length > 0) {
            ctx.strokeStyle = '#ef4444';
            ctx.lineWidth = 2;
            ctx.beginPath();

            data.forEach((point, index) => {
                const x = padding.left + (width - padding.left - padding.right) * index / (data.length - 1 || 1);
                const y = height - padding.bottom - ((point.count - minCount) / (maxCount - minCount)) * (height - padding.top - padding.bottom);

                if (index === 0) {
                    ctx.moveTo(x, y);
                } else {
                    ctx.lineTo(x, y);
                }
            });

            ctx.stroke();
        }

        // 绘制X轴标签
        ctx.fillStyle = '#666';
        ctx.font = '11px sans-serif';
        ctx.textAlign = 'center';
        data.forEach((point, index) => {
            if (index % Math.ceil(data.length / 8) === 0) {
                const x = padding.left + (width - padding.left - padding.right) * index / (data.length - 1 || 1);
                ctx.fillText(point.date, x, height - 10);
            }
        });

        // 绘制Y轴标签
        ctx.textAlign = 'right';
        for (let i = 0; i <= 5; i++) {
            const value = Math.round(maxCount * (5 - i) / 5);
            const y = padding.top + (height - padding.top - padding.bottom) * i / 5;
            ctx.fillText(value.toString(), padding.left - 10, y + 4);
        }
    }

    onMount(() => {
        ctx = canvas.getContext('2d');
        canvas.width = canvas.offsetWidth * 2;
        canvas.height = 300;
        drawChart();
    });

    afterUpdate(() => {
        drawChart();
    });
</script>

<div class="bg-white border border-gray-300 rounded p-4">
    <canvas
          bind:this={canvas}
          class="w-full"
          style="height: 150px;"
    ></canvas>
</div>