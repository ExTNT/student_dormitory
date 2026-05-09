<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue';
import { ElMessage } from 'element-plus';
import { useRoute } from 'vue-router';
import { cleaningApi } from '@/api/cleaning';
import type { PendingCleaning } from '@/types';
import StatusTag from '@/components/StatusTag.vue';
import ReviewButtons from '@/components/ReviewButtons.vue';
import AttachmentList from '@/components/AttachmentList.vue';
import { formatDateTime } from '@/utils/format';

const rows = ref<PendingCleaning[]>([]);
const route = useRoute();
const isAllOrders = computed(() => route.path.endsWith('/all'));
const title = computed(() => (isAllOrders.value ? '全部保洁工单' : '待审核保洁工单'));

async function fetchRows() {
  const allRows = await cleaningApi.list();
  rows.value = (isAllOrders.value ? allRows : allRows.filter((row) => row.status === 'cleaned'))
    .sort((a, b) => +new Date(b.created_at) - +new Date(a.created_at));
}
async function review(id: number, status: 'completed' | 'rejected', comment?: string) { await cleaningApi.review(id, { status, comment }); ElMessage.success('审核成功'); await fetchRows(); }
onMounted(fetchRows);
watch(() => route.path, fetchRows);
</script>

<template>
  <section class="page"><div class="toolbar"><h2>{{ title }}</h2><el-button @click="fetchRows">刷新</el-button></div>
    <el-table :data="rows" row-key="request_id">
      <el-table-column type="expand">
        <template #default="{ row }">
          <div class="work-order-detail">
            <el-descriptions title="保洁工单详情" :column="2" border>
              <el-descriptions-item label="工单 ID">{{ row.request_id }}</el-descriptions-item>
              <el-descriptions-item label="状态"><StatusTag :status="row.status" /></el-descriptions-item>
              <el-descriptions-item label="学生 ID">{{ row.student_id }}</el-descriptions-item>
              <el-descriptions-item label="学生">{{ row.student_name }}</el-descriptions-item>
              <el-descriptions-item label="楼栋 ID">{{ row.building_id }}</el-descriptions-item>
              <el-descriptions-item label="楼栋">{{ row.building_name }}</el-descriptions-item>
              <el-descriptions-item label="保洁人员 ID">{{ row.cleaner_id || '-' }}</el-descriptions-item>
              <el-descriptions-item label="保洁人员">{{ row.cleaner_name || '-' }}</el-descriptions-item>
              <el-descriptions-item label="审核人 ID">{{ row.reviewer_id || '-' }}</el-descriptions-item>
              <el-descriptions-item label="审核人">{{ row.reviewer_name || '-' }}</el-descriptions-item>
              <el-descriptions-item label="创建时间">{{ formatDateTime(row.created_at) }}</el-descriptions-item>
              <el-descriptions-item label="接单时间">{{ formatDateTime(row.accepted_at) }}</el-descriptions-item>
              <el-descriptions-item label="清洁完成时间">{{ formatDateTime(row.cleaned_at) }}</el-descriptions-item>
              <el-descriptions-item label="审核时间">{{ formatDateTime(row.reviewed_at) }}</el-descriptions-item>
              <el-descriptions-item label="位置描述" :span="2">{{ row.location_desc }}</el-descriptions-item>
              <el-descriptions-item label="审核意见" :span="2">{{ row.review_comment || '-' }}</el-descriptions-item>
            </el-descriptions>
            <div class="attachment-grid">
              <div class="attachment-section">
                <h3>保洁前照片</h3>
                <AttachmentList owner-type="cleaning" :owner-id="row.request_id" category="before" />
              </div>
              <div class="attachment-section">
                <h3>保洁后照片</h3>
                <AttachmentList owner-type="cleaning" :owner-id="row.request_id" category="after" />
              </div>
            </div>
          </div>
        </template>
      </el-table-column>
      <el-table-column prop="request_id" label="工单 ID" width="100" /><el-table-column label="状态"><template #default="{ row }"><StatusTag :status="row.status" /></template></el-table-column>
      <el-table-column prop="student_name" label="学生" /><el-table-column prop="building_name" label="楼栋" /><el-table-column prop="location_desc" label="位置" min-width="220" />
      <el-table-column label="创建时间" width="190"><template #default="{ row }">{{ formatDateTime(row.created_at) }}</template></el-table-column>
      <el-table-column prop="cleaner_name" label="保洁人员" />
      <el-table-column label="操作" width="150">
        <template #default="{ row }">
          <ReviewButtons v-if="row.status === 'cleaned'" with-comment approve-text="完成" @approve="(c) => review(row.request_id, 'completed', c)" @reject="(c) => review(row.request_id, 'rejected', c)" />
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

.attachment-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 16px;
}

.attachment-section h3 {
  margin: 0 0 10px;
  font-size: 16px;
}

@media (max-width: 900px) {
  .attachment-grid {
    grid-template-columns: 1fr;
  }
}
</style>
