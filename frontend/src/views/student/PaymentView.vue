<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue';
import { ElMessage } from 'element-plus';
import { studentApi } from '@/api/student';
import { useAuthStore } from '@/stores/auth';
import type { Room } from '@/types';

const auth = useAuthStore();
const balance = ref<Room>();
const loading = ref(false);
const balanceLoading = ref(false);
const form = reactive({ amount: 0, payment_type: 'both' });
const roomId = computed(() => auth.user?.room_id);
const hasDormitory = computed(() => Boolean(auth.user?.has_bed && roomId.value));
const dormitoryLabel = computed(() => {
  const building = auth.user?.building_name || '宿舍楼';
  const room = auth.user?.room_number || roomId.value;
  return `${building} ${room}`;
});

async function fetchBalance() {
  if (!roomId.value) return;
  balanceLoading.value = true;
  try {
    balance.value = await studentApi.roomBalance(roomId.value);
  } finally {
    balanceLoading.value = false;
  }
}

async function pay() {
  if (!roomId.value) return ElMessage.warning('当前账号未分配宿舍，无法缴费');
  if (form.amount <= 0) return ElMessage.warning('请输入大于 0 的缴费金额');
  loading.value = true;
  try {
    await studentApi.createPayment({ room_id: roomId.value, amount: form.amount, payment_type: form.payment_type });
    ElMessage.success('缴费成功');
    await fetchBalance();
  } finally {
    loading.value = false;
  }
}

onMounted(fetchBalance);
</script>

<template>
  <section class="page">
    <el-empty v-if="!hasDormitory" class="empty-state" description="当前账号未分配宿舍，暂无水电缴费入口" />
    <div class="form-card">
      <h2>水电缴费</h2>
      <el-form v-if="hasDormitory" :model="form" label-width="100px">
        <el-form-item label="当前宿舍">
          <el-tag type="primary">{{ dormitoryLabel }}</el-tag>
          <el-button style="margin-left: 12px" :loading="balanceLoading" @click="fetchBalance">刷新余额</el-button>
        </el-form-item>
        <el-form-item label="金额" required><el-input-number v-model="form.amount" :min="0" :precision="2" /></el-form-item>
        <el-form-item label="类型"><el-select v-model="form.payment_type"><el-option label="水费" value="water" /><el-option label="电费" value="electricity" /><el-option label="水电" value="both" /></el-select></el-form-item>
        <el-form-item><el-button type="primary" :loading="loading" @click="pay">缴费</el-button></el-form-item>
      </el-form>
    </div>
    <el-descriptions v-if="hasDormitory && balance" class="data-panel" :column="2" border>
      <el-descriptions-item label="房间">{{ balance.room_number }}</el-descriptions-item>
      <el-descriptions-item label="楼栋">{{ auth.user?.building_name || balance.building_id }}</el-descriptions-item>
      <el-descriptions-item label="水费余额">{{ balance.water_balance }}</el-descriptions-item>
      <el-descriptions-item label="电费余额">{{ balance.electricity_balance }}</el-descriptions-item>
    </el-descriptions>
  </section>
</template>

<style scoped>
.empty-state {
  min-height: 220px;
  border: 1px dashed rgba(148, 163, 184, 0.45);
  border-radius: 12px;
  background: rgba(248, 250, 252, 0.72);
}
</style>
