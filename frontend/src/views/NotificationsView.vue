<script setup lang="ts">
import { onMounted } from 'vue';
import { useNotificationStore } from '@/stores/notification';
import { formatDateTime } from '@/utils/format';

const store = useNotificationStore();
function rowClassName({ row }: { row: { is_read: number } }) {
  return row.is_read ? '' : 'unread';
}
onMounted(store.fetchNotifications);
</script>

<template>
  <section class="page">
    <div class="toolbar">
      <h2>通知</h2>
      <el-button @click="store.fetchNotifications">刷新</el-button>
    </div>
    <el-table :data="store.notifications" row-key="id" :row-class-name="rowClassName">
      <el-table-column prop="type" label="类型" width="140" />
      <el-table-column prop="message" label="内容" min-width="260" />
      <el-table-column label="时间" width="190">
        <template #default="{ row }">{{ formatDateTime(row.created_at) }}</template>
      </el-table-column>
      <el-table-column label="操作" width="120">
        <template #default="{ row }">
          <el-button v-if="!row.is_read" size="small" type="primary" @click="store.markRead(row.id)">标记已读</el-button>
          <el-tag v-else type="success">已读</el-tag>
        </template>
      </el-table-column>
    </el-table>
  </section>
</template>
