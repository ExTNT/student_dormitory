<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';
import { studentApi } from '@/api/student';
import { repairApi } from '@/api/repair';
import { cleaningApi } from '@/api/cleaning';
import type { MyRequest, PendingCleaning, PendingRepair } from '@/types';
import StatusTag from '@/components/StatusTag.vue';
import AttachmentList from '@/components/AttachmentList.vue';
import { formatDateTime } from '@/utils/format';

const rows = ref<MyRequest[]>([]);
const repairs = ref<PendingRepair[]>([]);
const cleanings = ref<PendingCleaning[]>([]);
const loading = ref(false);
const sortedRows = computed(() => [...rows.value].sort((a, b) => +new Date(b.created_at) - +new Date(a.created_at)));
const sortedRepairs = computed(() => [...repairs.value].sort((a, b) => +new Date(b.created_at) - +new Date(a.created_at)));
const sortedCleanings = computed(() => [...cleanings.value].sort((a, b) => +new Date(b.created_at) - +new Date(a.created_at)));

async function fetchRows() {
  loading.value = true;
  try {
    const [requestRows, repairRows, cleaningRows] = await Promise.all([
      studentApi.requests(),
      repairApi.list(),
      cleaningApi.list(),
    ]);
    rows.value = requestRows;
    repairs.value = repairRows;
    cleanings.value = cleaningRows;
  } finally {
    loading.value = false;
  }
}

onMounted(fetchRows);
</script>

<template>
  <section class="page">
    <div class="toolbar">
      <h2>我的申请总览</h2>
      <el-button @click="fetchRows">刷新</el-button>
    </div>
    <el-table v-loading="loading" :data="sortedRows">
      <el-table-column prop="request_type" label="类型" width="150" />
      <el-table-column prop="request_id" label="申请 ID" width="110" />
      <el-table-column label="状态" width="120"><template #default="{ row }"><StatusTag :status="row.status" /></template></el-table-column>
      <el-table-column label="创建时间" width="190"><template #default="{ row }">{{ formatDateTime(row.created_at) }}</template></el-table-column>
      <el-table-column prop="detail" label="详情" min-width="260" />
    </el-table>
    <div class="data-panel">
      <h2>我的维修工单详情</h2>
      <el-table :data="sortedRepairs" row-key="request_id">
        <el-table-column type="expand">
          <template #default="{ row }">
            <el-descriptions :column="2" border>
              <el-descriptions-item label="工单 ID">{{ row.request_id }}</el-descriptions-item>
              <el-descriptions-item label="状态"><StatusTag :status="row.status" /></el-descriptions-item>
              <el-descriptions-item label="房间">{{ row.room_number }}</el-descriptions-item>
              <el-descriptions-item label="维修人员">{{ row.repair_staff_name || '-' }}</el-descriptions-item>
              <el-descriptions-item label="创建时间">{{ formatDateTime(row.created_at) }}</el-descriptions-item>
              <el-descriptions-item label="完成时间">{{ formatDateTime(row.repaired_at) }}</el-descriptions-item>
              <el-descriptions-item label="报修描述" :span="2">{{ row.description }}</el-descriptions-item>
              <el-descriptions-item label="维修说明" :span="2">{{ row.repair_description || '-' }}</el-descriptions-item>
              <el-descriptions-item label="审核意见" :span="2">{{ row.review_comment || '-' }}</el-descriptions-item>
            </el-descriptions>
            <div class="attachment-section">
              <h3>维修后照片</h3>
              <AttachmentList owner-type="repair" :owner-id="row.request_id" category="after" />
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="request_id" label="工单 ID" width="110" />
        <el-table-column label="状态" width="120"><template #default="{ row }"><StatusTag :status="row.status" /></template></el-table-column>
        <el-table-column prop="room_number" label="房间" />
        <el-table-column prop="description" label="描述" min-width="220" />
        <el-table-column label="创建时间" width="190"><template #default="{ row }">{{ formatDateTime(row.created_at) }}</template></el-table-column>
      </el-table>
    </div>
    <div class="data-panel">
      <h2>我的保洁工单详情</h2>
      <el-table :data="sortedCleanings" row-key="request_id">
        <el-table-column type="expand">
          <template #default="{ row }">
            <el-descriptions :column="2" border>
              <el-descriptions-item label="工单 ID">{{ row.request_id }}</el-descriptions-item>
              <el-descriptions-item label="状态"><StatusTag :status="row.status" /></el-descriptions-item>
              <el-descriptions-item label="楼栋">{{ row.building_name }}</el-descriptions-item>
              <el-descriptions-item label="保洁人员">{{ row.cleaner_name || '-' }}</el-descriptions-item>
              <el-descriptions-item label="创建时间">{{ formatDateTime(row.created_at) }}</el-descriptions-item>
              <el-descriptions-item label="完成时间">{{ formatDateTime(row.cleaned_at) }}</el-descriptions-item>
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
          </template>
        </el-table-column>
        <el-table-column prop="request_id" label="工单 ID" width="110" />
        <el-table-column label="状态" width="120"><template #default="{ row }"><StatusTag :status="row.status" /></template></el-table-column>
        <el-table-column prop="building_name" label="楼栋" />
        <el-table-column prop="location_desc" label="位置" min-width="220" />
        <el-table-column label="创建时间" width="190"><template #default="{ row }">{{ formatDateTime(row.created_at) }}</template></el-table-column>
      </el-table>
    </div>
  </section>
</template>

<style scoped>
.attachment-section {
  margin-top: 16px;
}

.attachment-section h3 {
  margin: 0 0 10px;
  font-size: 16px;
}

.attachment-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 16px;
}

@media (max-width: 900px) {
  .attachment-grid {
    grid-template-columns: 1fr;
  }
}
</style>
