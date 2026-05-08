<script setup lang="ts">
import { reactive, ref } from 'vue';
import { ElMessage } from 'element-plus';
import { studentApi } from '@/api/student';
import type { Room } from '@/types';

const balance = ref<Room>();
const loading = ref(false);
const form = reactive({ room_id: undefined as number | undefined, amount: 0, payment_type: 'both' });

async function fetchBalance() {
  if (!form.room_id) return ElMessage.warning('请输入房间 ID');
  balance.value = await studentApi.roomBalance(form.room_id);
}

async function pay() {
  if (!form.room_id) return ElMessage.warning('请输入房间 ID');
  loading.value = true;
  try {
    await studentApi.createPayment({ room_id: form.room_id, amount: form.amount, payment_type: form.payment_type });
    ElMessage.success('缴费成功');
    await fetchBalance();
  } finally {
    loading.value = false;
  }
}
</script>

<template>
  <section class="page">
    <div class="form-card">
      <h2>水电缴费</h2>
      <el-form :model="form" label-width="100px">
        <el-form-item label="房间 ID" required>
          <el-input-number v-model="form.room_id" :min="1" />
          <el-button style="margin-left: 12px" @click="fetchBalance">查询余额</el-button>
        </el-form-item>
        <el-form-item label="金额" required><el-input-number v-model="form.amount" :min="0" :precision="2" /></el-form-item>
        <el-form-item label="类型"><el-select v-model="form.payment_type"><el-option label="水费" value="water" /><el-option label="电费" value="electricity" /><el-option label="水电" value="both" /></el-select></el-form-item>
        <el-form-item><el-button type="primary" :loading="loading" @click="pay">缴费</el-button></el-form-item>
      </el-form>
    </div>
    <el-descriptions v-if="balance" class="data-panel" :column="2" border>
      <el-descriptions-item label="房间">{{ balance.room_number }}</el-descriptions-item>
      <el-descriptions-item label="楼栋 ID">{{ balance.building_id }}</el-descriptions-item>
      <el-descriptions-item label="水费余额">{{ balance.water_balance }}</el-descriptions-item>
      <el-descriptions-item label="电费余额">{{ balance.electricity_balance }}</el-descriptions-item>
    </el-descriptions>
  </section>
</template>
