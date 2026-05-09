<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue';
import { ElMessage } from 'element-plus';
import { useRoute } from 'vue-router';
import { repairApi } from '@/api/repair';
import type { PendingRepair } from '@/types';
import StatusTag from '@/components/StatusTag.vue';
import ReviewButtons from '@/components/ReviewButtons.vue';
import AttachmentList from '@/components/AttachmentList.vue';
import { formatDateTime } from '@/utils/format';

const rows = ref<PendingRepair[]>([]);
const route = useRoute();
const isAllOrders = computed(() => route.path.endsWith('/all'));
const title = computed(() => (isAllOrders.value ? '全部维修工单' : '待审核维修工单'));

async function fetchRows() {
  const allRows = await repairApi.list();
  rows.value = (isAllOrders.value ? allRows : allRows.filter((row) => row.status === 'repaired'))
    .sort((a, b) => +new Date(b.created_at) - +new Date(a.created_at));
}
async function review(id: number, status: 'completed' | 'rejected', comment?: string) { await repairApi.review(id, { status, comment }); ElMessage.success('审核成功'); await fetchRows(); }
onMounted(fetchRows);
watch(() => route.path, fetchRows);
</script>

<template>
  <section class="page"><div class="toolbar"><h2>{{ title }}</h2><el-button @click="fetchRows">刷新</el-button></div>
    <el-table :data="rows" row-key="request_id">
      <el-table-column type="expand">
        <template #default="{ row }">
          <div class="work-order-detail">
            <el-descriptions title="维修工单详情" :column="2" border>
              <el-descriptions-item label="工单 ID">{{ row.request_id }}</el-descriptions-item>
              <el-descriptions-item label="状态"><StatusTag :status="row.status" /></el-descriptions-item>
              <el-descriptions-item label="学生 ID">{{ row.student_id }}</el-descriptions-item>
              <el-descriptions-item label="学生">{{ row.student_name }}</el-descriptions-item>
              <el-descriptions-item label="房间 ID">{{ row.room_id }}</el-descriptions-item>
              <el-descriptions-item label="房间号">{{ row.room_number }}</el-descriptions-item>
              <el-descriptions-item label="维修人员 ID">{{ row.repair_staff_id || '-' }}</el-descriptions-item>
              <el-descriptions-item label="维修人员">{{ row.repair_staff_name || '-' }}</el-descriptions-item>
              <el-descriptions-item label="审核人 ID">{{ row.reviewer_id || '-' }}</el-descriptions-item>
              <el-descriptions-item label="审核人">{{ row.reviewer_name || '-' }}</el-descriptions-item>
              <el-descriptions-item label="创建时间">{{ formatDateTime(row.created_at) }}</el-descriptions-item>
              <el-descriptions-item label="接单时间">{{ formatDateTime(row.accepted_at) }}</el-descriptions-item>
              <el-descriptions-item label="维修完成时间">{{ formatDateTime(row.repaired_at) }}</el-descriptions-item>
              <el-descriptions-item label="审核时间">{{ formatDateTime(row.reviewed_at) }}</el-descriptions-item>
              <el-descriptions-item label="报修描述" :span="2">{{ row.description }}</el-descriptions-item>
              <el-descriptions-item label="维修说明" :span="2">{{ row.repair_description || '-' }}</el-descriptions-item>
              <el-descriptions-item label="审核意见" :span="2">{{ row.review_comment || '-' }}</el-descriptions-item>
            </el-descriptions>
            <div class="attachment-section">
              <h3>维修后照片</h3>
              <AttachmentList owner-type="repair" :owner-id="row.request_id" category="after" />
            </div>
          </div>
        </template>
      </el-table-column>
      <el-table-column prop="request_id" label="工单 ID" width="100" /><el-table-column label="状态"><template #default="{ row }"><StatusTag :status="row.status" /></template></el-table-column>
      <el-table-column prop="student_name" label="学生" /><el-table-column prop="room_number" label="房间" /><el-table-column prop="description" label="描述" min-width="220" />
      <el-table-column label="创建时间" width="190"><template #default="{ row }">{{ formatDateTime(row.created_at) }}</template></el-table-column>
      <el-table-column prop="repair_staff_name" label="维修人员" />
      <el-table-column label="操作" width="150">
        <template #default="{ row }">
          <ReviewButtons v-if="row.status === 'repaired'" with-comment approve-text="完成" @approve="(c) => review(row.request_id, 'completed', c)" @reject="(c) => review(row.request_id, 'rejected', c)" />
          <span v-else class="muted">待处理</span>
        </template>
      </el-table-column>
    </el-table>
  </section>
</template>

<style scoped>
.work-order-detail {
  display: grid;
  gap: 16px;
  padding: 12px;
}

.attachment-section h3 {
  margin: 0 0 10px;
  font-size: 16px;
}
</style>
