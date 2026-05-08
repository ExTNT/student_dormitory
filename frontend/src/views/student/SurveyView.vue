<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';
import { ElMessage } from 'element-plus';
import { studentApi } from '@/api/student';
import type { LifestyleSurvey } from '@/types';

const loading = ref(false);
const form = reactive<LifestyleSurvey>({ sleep_time: '23:30:00', smoking: 0, snoring: 0, study_habit: '', remarks: '' });

async function fetchSurvey() {
  try {
    const data = await studentApi.getSurvey();
    Object.assign(form, data);
  } catch {
    // 没提交过调查时后端可能返回 404，保留空表单。
  }
}

async function submit() {
  loading.value = true;
  try {
    await studentApi.saveSurvey({
      sleep_time: form.sleep_time,
      smoking: form.smoking,
      snoring: form.snoring,
      study_habit: form.study_habit,
      remarks: form.remarks,
    });
    ElMessage.success('提交成功');
    await fetchSurvey();
  } finally {
    loading.value = false;
  }
}

onMounted(fetchSurvey);
</script>

<template>
  <section class="form-card">
    <h2>生活习惯调查</h2>
    <el-form :model="form" label-width="110px">
      <el-form-item label="就寝时间">
        <el-time-picker v-model="form.sleep_time" value-format="HH:mm:ss" format="HH:mm:ss" />
      </el-form-item>
      <el-form-item label="是否吸烟">
        <el-switch v-model="form.smoking" :active-value="1" :inactive-value="0" />
      </el-form-item>
      <el-form-item label="是否打鼾">
        <el-radio-group v-model="form.snoring">
          <el-radio :value="0">否</el-radio>
          <el-radio :value="1">是</el-radio>
        </el-radio-group>
      </el-form-item>
      <el-form-item label="学习习惯">
        <el-input v-model="form.study_habit" />
      </el-form-item>
      <el-form-item label="备注">
        <el-input v-model="form.remarks" type="textarea" />
      </el-form-item>
      <el-form-item>
        <el-button type="primary" :loading="loading" @click="submit">提交</el-button>
      </el-form-item>
    </el-form>
  </section>
</template>
