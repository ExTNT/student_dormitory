<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { studentApi } from '@/api/student';
import type { Roommate } from '@/types';
import AttachmentImage from '@/components/AttachmentImage.vue';

const rows = ref<Roommate[]>([]);
onMounted(async () => {
  rows.value = await studentApi.roommates();
});
</script>

<template>
  <section class="page">
    <h2>舍友信息</h2>
    <el-table :data="rows">
      <el-table-column label="头像" width="110">
        <template #default="{ row }"><AttachmentImage v-if="row.avatar_attachment_id" :id="row.avatar_attachment_id" /></template>
      </el-table-column>
      <el-table-column prop="roommate_name" label="姓名" />
      <el-table-column prop="roommate_phone" label="电话" />
      <el-table-column prop="bed_label" label="床位" />
    </el-table>
  </section>
</template>
