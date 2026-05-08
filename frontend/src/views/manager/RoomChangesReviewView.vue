<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { ElMessage } from 'element-plus';
import { managerApi } from '@/api/manager';
import type { PendingRoomChange } from '@/types';
import StatusTag from '@/components/StatusTag.vue';
import ReviewButtons from '@/components/ReviewButtons.vue';
import { formatDateTime } from '@/utils/format';

const rows = ref<PendingRoomChange[]>([]);
async function fetchRows() { rows.value = await managerApi.roomChanges(); }
async function review(id: number, status: 'approved' | 'rejected') { await managerApi.reviewRoomChange(id, status); ElMessage.success('审批成功'); await fetchRows(); }
function targetText(row: PendingRoomChange) {
  if (row.target_bed_id) return `${row.target_building_name} ${row.target_room_number} ${row.target_bed_label}`;
  return `${row.recommended_building_name || '-'} ${row.recommended_room_number || ''} ${row.recommended_bed_label || ''}`;
}
onMounted(fetchRows);
</script>

<template>
  <section class="page"><div class="toolbar"><h2>待审批换寝申请</h2><el-button @click="fetchRows">刷新</el-button></div>
    <el-table :data="rows">
      <el-table-column prop="id" label="ID" width="80" /><el-table-column prop="student_name" label="学生" />
      <el-table-column label="原床位" min-width="180"><template #default="{ row }">{{ row.from_building_name }} {{ row.from_room_number }} {{ row.from_bed_label }}</template></el-table-column>
      <el-table-column label="目标/推荐床位" min-width="190"><template #default="{ row }">{{ targetText(row) }}</template></el-table-column>
      <el-table-column prop="reason" label="原因" min-width="180" /><el-table-column label="状态"><template #default="{ row }"><StatusTag :status="row.status" /></template></el-table-column>
      <el-table-column label="创建时间" width="190"><template #default="{ row }">{{ formatDateTime(row.created_at) }}</template></el-table-column>
      <el-table-column label="操作" width="150"><template #default="{ row }"><ReviewButtons @approve="review(row.id, 'approved')" @reject="review(row.id, 'rejected')" /></template></el-table-column>
    </el-table>
  </section>
</template>
