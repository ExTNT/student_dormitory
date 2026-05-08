<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { dashboardApi } from '@/api/dashboard';
import type { LowBalanceRoom } from '@/types';

const rows = ref<LowBalanceRoom[]>([]);
async function fetchRows() { rows.value = await dashboardApi.lowBalance(); }
onMounted(fetchRows);
</script>

<template>
  <section class="page">
    <div class="toolbar"><h2>低余额房间</h2><el-button @click="fetchRows">刷新</el-button></div>
    <el-table :data="rows">
      <el-table-column prop="room_id" label="房间 ID" />
      <el-table-column prop="building_id" label="楼栋 ID" />
      <el-table-column prop="room_number" label="房号" />
      <el-table-column prop="water_balance" label="水费余额" />
      <el-table-column prop="electricity_balance" label="电费余额" />
    </el-table>
  </section>
</template>
