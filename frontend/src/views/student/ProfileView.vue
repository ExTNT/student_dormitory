<script setup lang="ts">
import { computed, reactive, ref, watch } from 'vue';
import { ElMessage } from 'element-plus';
import { useAuthStore } from '@/stores/auth';
import { attachmentApi } from '@/api/attachment';
import AttachmentImage from '@/components/AttachmentImage.vue';
import ImageUploader from '@/components/ImageUploader.vue';
import { formatDateTime, roleLabels } from '@/utils/format';
import type { AttachmentMeta } from '@/types';

const auth = useAuthStore();
const avatarId = ref<number>();
const savingPhone = ref(false);
const phoneForm = reactive({ phone: '' });
const canEditProfile = computed(() => ['student', 'repair_staff', 'cleaning_staff'].includes(auth.user?.role || ''));
const showDormitory = computed(() => auth.user?.role === 'student');
const dormitoryText = computed(() => {
  const user = auth.user;
  if (!user?.has_bed) return '未分配宿舍';
  return [user.building_name, user.room_number, user.bed_label].filter(Boolean).join(' / ');
});

async function fetchAvatar() {
  if (!auth.user?.id) {
    avatarId.value = undefined;
    return;
  }
  avatarId.value = auth.user.avatar_attachment_id;
  if (avatarId.value) return;
  try {
    const metas = await attachmentApi.list({ owner_type: 'user_avatar', owner_id: auth.user.id, category: 'avatar' });
    avatarId.value = metas[0]?.id;
  } catch {
    avatarId.value = undefined;
  }
}

async function handleAvatarSuccess(meta: AttachmentMeta) {
  avatarId.value = meta.id;
  await auth.fetchMe();
}

async function savePhone() {
  savingPhone.value = true;
  try {
    await auth.updateMe({ phone: phoneForm.phone.trim() || undefined });
    ElMessage.success('电话已更新');
  } finally {
    savingPhone.value = false;
  }
}

watch(
  () => auth.user,
  (user) => {
    phoneForm.phone = user?.phone || '';
    fetchAvatar();
  },
  { immediate: true },
);
</script>

<template>
  <section class="page profile-page">
    <div class="profile-card">
      <div class="profile-avatar">
        <AttachmentImage v-if="avatarId" :id="avatarId" width="112px" height="112px" />
        <div v-else class="default-avatar">{{ auth.user?.name?.slice(0, 1) || auth.user?.username?.slice(0, 1) || '用' }}</div>
        <ImageUploader
          v-if="canEditProfile && auth.user?.id"
          owner-type="user_avatar"
          :owner-id="auth.user.id"
          category="avatar"
          button-text="上传头像"
          tip="仅支持 jpg/png，最大 5MB"
          @success="handleAvatarSuccess"
        />
      </div>
      <div class="profile-info">
        <div class="toolbar">
          <h2>个人信息</h2>
        </div>
        <el-descriptions :column="1" border>
          <el-descriptions-item label="用户名">{{ auth.user?.username }}</el-descriptions-item>
          <el-descriptions-item label="姓名">{{ auth.user?.name }}</el-descriptions-item>
          <el-descriptions-item label="角色">{{ auth.user ? roleLabels[auth.user.role] : '' }}</el-descriptions-item>
          <el-descriptions-item label="电话">{{ auth.user?.phone || '-' }}</el-descriptions-item>
          <el-descriptions-item v-if="showDormitory" label="宿舍号">{{ dormitoryText }}</el-descriptions-item>
          <el-descriptions-item label="创建时间">{{ formatDateTime(auth.user?.created_at) }}</el-descriptions-item>
        </el-descriptions>
        <el-form v-if="canEditProfile" class="phone-form" :model="phoneForm" label-width="80px">
          <el-form-item label="修改电话">
            <el-input v-model="phoneForm.phone" placeholder="请输入联系电话" maxlength="20" clearable />
            <el-button type="primary" :loading="savingPhone" @click="savePhone">保存</el-button>
          </el-form-item>
        </el-form>
      </div>
    </div>
  </section>
</template>

<style scoped>
.profile-card {
  display: grid;
  grid-template-columns: 220px minmax(0, 1fr);
  gap: 20px;
  align-items: start;
}

.profile-avatar,
.profile-info {
  border: 1px solid var(--border);
  border-radius: 12px;
  background: rgba(255, 255, 255, 0.86);
  box-shadow: var(--shadow-sm);
}

.profile-avatar {
  display: grid;
  justify-items: center;
  gap: 14px;
  padding: 22px 16px;
}

.profile-info {
  padding: 20px;
}

.default-avatar {
  display: grid;
  place-items: center;
  width: 112px;
  height: 112px;
  border: 1px solid rgba(11, 93, 102, 0.18);
  border-radius: 16px;
  background: linear-gradient(135deg, rgba(11, 93, 102, 0.14), rgba(245, 158, 11, 0.16));
  color: var(--brand-strong);
  font-size: 42px;
  font-weight: 900;
}

.phone-form {
  margin-top: 18px;
  max-width: 520px;
}

.phone-form :deep(.el-form-item__content) {
  display: flex;
  flex-wrap: nowrap;
  gap: 10px;
}

@media (max-width: 900px) {
  .profile-card {
    grid-template-columns: 1fr;
  }
}
</style>
