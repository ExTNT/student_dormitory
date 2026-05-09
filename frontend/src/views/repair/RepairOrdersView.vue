<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import { useRoute } from 'vue-router';
import { repairApi } from '@/api/repair';
import { attachmentApi } from '@/api/attachment';
import type { PendingRepair } from '@/types';
import StatusTag from '@/components/StatusTag.vue';
import ImageUploader from '@/components/ImageUploader.vue';
import AttachmentList from '@/components/AttachmentList.vue';
import { formatDateTime } from '@/utils/format';

const rows = ref<PendingRepair[]>([]);
const loading = ref(false);
const lists = ref<Record<number, InstanceType<typeof AttachmentList>>>({});
const route = useRoute();
const isMyOrders = computed(() => route.path.endsWith('/my-orders'));
const title = computed(() => (isMyOrders.value ? '我的维修工单' : '接收维修工单'));

async function fetchRows() {
  loading.value = true;
  try {
    if (isMyOrders.value) {
      rows.value = (await repairApi.list()).sort((a, b) => +new Date(b.created_at) - +new Date(a.created_at));
      return;
    }
    rows.value = (await repairApi.pending())
      .filter((row) => row.status === 'pending')
      .sort((a, b) => +new Date(b.created_at) - +new Date(a.created_at));
  } finally {
    loading.value = false;
  }
}

async function accept(id: number) {
  await ElMessageBox.confirm('确认接单？', '接单', { type: 'warning' });
  await repairApi.accept(id);
  ElMessage.success('接单成功');
  await fetchRows();
}

async function complete(id: number) {
  const attachments = await attachmentApi.list({ owner_type: 'repair', owner_id: id, category: 'after' });
  if (!attachments.length) return ElMessage.warning('请先上传维修后照片');
  const { value } = await ElMessageBox.prompt('请输入维修说明', '完成维修', {
    inputValidator: (v) => Boolean(v) || '维修说明不能为空',
  });
  await repairApi.complete(id, value);
  ElMessage.success('维修已完成');
  await fetchRows();
}

onMounted(fetchRows);
watch(() => route.path, fetchRows);
</script>

<template>
  <section class="page">
    <div class="toolbar"><h2>{{ title }}</h2><el-button @click="fetchRows">刷新</el-button></div>
    <el-table v-loading="loading" :data="rows" row-key="request_id">
      <el-table-column type="expand">
        <template #default="{ row }">
          <div class="data-panel">
            <div class="toolbar">
              <strong>维修后照片</strong>
              <ImageUploader
                v-if="isMyOrders && row.status === 'accepted'"
                owner-type="repair"
                :owner-id="row.request_id"
                category="after"
                button-text="上传维修后照片"
                tip="完成维修前至少上传 1 张维修后照片。"
                @success="lists[row.request_id]?.fetchList()"
              />
            </div>
            <AttachmentList :ref="(el) => { if (el) lists[row.request_id] = el as any }" owner-type="repair" :owner-id="row.request_id" category="after" />
            <el-descriptions style="margin-top: 16px" :column="2" border>
              <el-descriptions-item label="学生">{{ row.student_name }}</el-descriptions-item>
              <el-descriptions-item label="房间">{{ row.room_number }}</el-descriptions-item>
              <el-descriptions-item label="维修说明">{{ row.repair_description || '-' }}</el-descriptions-item>
              <el-descriptions-item label="审核意见">{{ row.review_comment || '-' }}</el-descriptions-item>
              <el-descriptions-item label="接单时间">{{ formatDateTime(row.accepted_at) }}</el-descriptions-item>
              <el-descriptions-item label="完成时间">{{ formatDateTime(row.repaired_at) }}</el-descriptions-item>
            </el-descriptions>
          </div>
        </template>
      </el-table-column>
      <el-table-column prop="request_id" label="工单 ID" width="100" />
      <el-table-column label="状态" width="120"><template #default="{ row }"><StatusTag :status="row.status" /></template></el-table-column>
      <el-table-column prop="student_name" label="学生" />
      <el-table-column prop="room_number" label="房间" />
      <el-table-column prop="description" label="描述" min-width="220" />
      <el-table-column label="创建时间" width="190"><template #default="{ row }">{{ formatDateTime(row.created_at) }}</template></el-table-column>
      <el-table-column prop="repair_staff_name" label="维修人员" />
      <el-table-column prop="reviewer_name" label="审核人" />
      <el-table-column label="操作" width="180" fixed="right">
        <template #default="{ row }">
          <el-button v-if="!isMyOrders && row.status === 'pending'" size="small" type="primary" @click="accept(row.request_id)">接单</el-button>
          <el-button v-if="isMyOrders && row.status === 'accepted'" size="small" type="success" @click="complete(row.request_id)">完成</el-button>
          <span v-if="(isMyOrders && row.status !== 'accepted') || (!isMyOrders && row.status !== 'pending')" class="muted">-</span>
        </template>
      </el-table-column>
    </el-table>
  </section>
</template>
