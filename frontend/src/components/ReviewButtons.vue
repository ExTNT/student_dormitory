<script setup lang="ts">
import { ElMessageBox } from 'element-plus';

const props = defineProps<{ approveText?: string; rejectText?: string; withComment?: boolean }>();
const emit = defineEmits<{ approve: [comment?: string]; reject: [comment?: string] }>();

async function confirmAction(action: 'approve' | 'reject') {
  const isApprove = action === 'approve';
  if (props.withComment) {
    const result = await ElMessageBox.prompt(isApprove ? '请输入审核意见' : '请输入驳回原因', '确认操作', {
      inputValue: isApprove ? '处理合格' : '',
      inputValidator: (value) => (isApprove || Boolean(value)) || '驳回原因不能为空',
      confirmButtonText: '确认',
      cancelButtonText: '取消',
      type: isApprove ? 'success' : 'warning',
    });
    if (action === 'approve') emit('approve', result.value);
    else emit('reject', result.value);
    return;
  }
  await ElMessageBox.confirm(isApprove ? '确认审批通过？' : '确认驳回该申请？', '确认操作', {
    type: isApprove ? 'success' : 'warning',
  });
  if (action === 'approve') emit('approve');
  else emit('reject');
}
</script>

<template>
  <el-button size="small" type="success" @click="confirmAction('approve')">{{ approveText || '通过' }}</el-button>
  <el-button size="small" type="danger" @click="confirmAction('reject')">{{ rejectText || '驳回' }}</el-button>
</template>
