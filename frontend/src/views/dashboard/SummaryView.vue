<script setup lang="ts">
import { nextTick, onBeforeUnmount, onMounted, ref } from 'vue';
import * as echarts from 'echarts';
import { dashboardApi } from '@/api/dashboard';
import type { DormitorySummary } from '@/types';

const rows = ref<DormitorySummary[]>([]);
const chartEl = ref<HTMLDivElement>();
let chart: echarts.ECharts | null = null;

function renderChart() {
  if (!chartEl.value) return;
  chart ||= echarts.init(chartEl.value);
  chart.setOption({
    tooltip: { trigger: 'axis' },
    legend: { data: ['总床位', '已入住', '空床位'] },
    xAxis: { type: 'category', data: rows.value.map((r) => r.building_name) },
    yAxis: { type: 'value' },
    series: [
      { name: '总床位', type: 'bar', data: rows.value.map((r) => r.total_beds) },
      { name: '已入住', type: 'bar', data: rows.value.map((r) => r.occupied_beds) },
      { name: '空床位', type: 'bar', data: rows.value.map((r) => r.free_beds) },
    ],
  });
}

async function fetchRows() {
  rows.value = await dashboardApi.summary();
  await nextTick();
  renderChart();
}

onMounted(fetchRows);
onBeforeUnmount(() => chart?.dispose());
</script>

<template>
  <section class="page">
    <div class="toolbar"><h2>楼栋统计</h2><el-button @click="fetchRows">刷新</el-button></div>
    <div ref="chartEl" class="data-panel chart"></div>
    <el-table :data="rows">
      <el-table-column prop="building_name" label="楼栋" />
      <el-table-column prop="total_rooms" label="房间数" />
      <el-table-column prop="total_beds" label="总床位" />
      <el-table-column prop="occupied_beds" label="已入住" />
      <el-table-column prop="free_beds" label="空床位" />
    </el-table>
  </section>
</template>
