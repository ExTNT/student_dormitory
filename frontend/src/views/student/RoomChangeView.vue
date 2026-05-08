<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';
import { ElMessage } from 'element-plus';
import { studentApi } from '@/api/student';
import type { AvailableBed } from '@/types';

const beds = ref<AvailableBed[]>([]);
const selected = ref<AvailableBed>();
const specify = ref(false);
const loading = ref(false);
const form = reactive({ building_id: undefined as number | undefined, floor: undefined as number | undefined, reason: '' });

async function fetchBeds() {
  beds.value = await studentApi.availableBeds({ building_id: form.building_id, floor: form.floor });
}

async function submit() {
  const data: Record<string, unknown> = { reason: form.reason };
  if (specify.value) {
    if (!selected.value) {
      ElMessage.warning('请选择目标空床位');
      return;
    }
    data.target_room_id = selected.value.room_id;
    data.target_bed_id = selected.value.bed_id;
  }
  loading.value = true;
  try {
    await studentApi.createRoomChange(data);
    ElMessage.success('申请已提交');
  } finally {
    loading.value = false;
  }
}

onMounted(fetchBeds);
</script>

<template>
  <section class="page">
    <div class="form-card">
      <h2>换寝申请</h2>
      <el-form :model="form" label-width="120px">
        <el-form-item label="指定目标">
          <el-switch v-model="specify" active-text="指定空床位" inactive-text="系统推荐" />
        </el-form-item>
        <el-form-item label="原因" required><el-input v-model="form.reason" type="textarea" /></el-form-item>
        <el-form-item><el-button type="primary" :loading="loading" @click="submit">提交申请</el-button></el-form-item>
      </el-form>
    </div>
    <div v-if="specify" class="data-panel">
      <div class="toolbar">
        <strong>空床位</strong>
        <div>
          <el-input-number v-model="form.building_id" :min="1" placeholder="楼栋" />
          <el-input-number v-model="form.floor" :min="1" placeholder="楼层" />
          <el-button @click="fetchBeds">查询</el-button>
        </div>
      </div>
      <el-table :data="beds" highlight-current-row @current-change="(row: AvailableBed) => (selected = row)">
        <el-table-column prop="building_name" label="楼栋" />
        <el-table-column prop="floor" label="楼层" />
        <el-table-column prop="room_number" label="房间" />
        <el-table-column prop="bed_label" label="床位" />
      </el-table>
    </div>
  </section>
</template>
