<script setup lang="ts">
import { reactive, ref, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import { useAuthStore } from '@/store/modules/auth'
import { updateUser } from '@/api/user'
import { uploadFile } from '@/api/file'

const props = defineProps<{ modelValue: boolean }>()
const emit = defineEmits<{ (e: 'update:modelValue', v: boolean): void }>()

const authStore = useAuthStore()

interface FormState {
  display_name: string
  phone_number: string
  email_address: string
  remark: string
}

const form = reactive<FormState>({
  display_name: '',
  phone_number: '',
  email_address: '',
  remark: '',
})

const avatarPreview = ref<string | null>(null)
const pendingAvatarFile = ref<File | null>(null)
const saving = ref(false)

watch(
  () => props.modelValue,
  (open) => {
    if (!open) return
    const u = authStore.currentUser
    form.display_name = u?.display_name ?? ''
    form.phone_number = u?.phone_number ?? ''
    form.email_address = u?.email_address ?? ''
    form.remark = u?.remark ?? ''
    avatarPreview.value = authStore.avatarBlobUrl
    pendingAvatarFile.value = null
  },
)

function handleAvatarChange(file: File) {
  pendingAvatarFile.value = file
  avatarPreview.value = URL.createObjectURL(file)
}

function triggerFileInput() {
  const input = document.createElement('input')
  input.type = 'file'
  input.accept = 'image/*'
  input.onchange = (e) => {
    const f = (e.target as HTMLInputElement).files?.[0]
    if (f) handleAvatarChange(f)
  }
  input.click()
}

async function handleSave() {
  if (!authStore.currentUser) return
  saving.value = true
  try {
    let avatarUrl = authStore.currentUser.avatar_url ?? null

    if (pendingAvatarFile.value) {
      const record = await uploadFile(pendingAvatarFile.value)
      avatarUrl = `/api/v1/files/${record.id}/raw`
    }

    await updateUser(authStore.currentUser.id, {
      display_name: form.display_name,
      phone_number: form.phone_number || null,
      email_address: form.email_address || null,
      remark: form.remark || null,
      avatar_url: avatarUrl,
    })

    await authStore.refreshProfile()
    ElMessage.success('个人信息已更新')
    emit('update:modelValue', false)
  } catch (e) {
    ElMessage.error(e instanceof Error ? e.message : '保存失败')
  } finally {
    saving.value = false
  }
}

function handleClose() {
  emit('update:modelValue', false)
}
</script>

<template>
  <el-drawer
    :model-value="modelValue"
    title="个人信息"
    size="400px"
    :append-to-body="true"
    destroy-on-close
    @update:model-value="emit('update:modelValue', $event)"
  >
    <div class="profile-body">
      <!-- Avatar -->
      <div class="avatar-section">
        <div class="avatar-wrap" @click="triggerFileInput">
          <img v-if="avatarPreview" :src="avatarPreview" class="avatar-img" alt="avatar" />
          <div v-else class="avatar-placeholder">
            <el-icon :size="28"><Plus /></el-icon>
          </div>
          <div class="avatar-overlay">更换头像</div>
        </div>
        <p class="avatar-hint">点击头像上传新图片</p>
      </div>

      <!-- Form -->
      <el-form label-position="top" @submit.prevent="handleSave">
        <el-form-item label="登录账号">
          <el-input :value="authStore.currentUser?.login_name" disabled />
        </el-form-item>

        <el-form-item label="显示姓名">
          <el-input v-model="form.display_name" placeholder="请输入显示姓名" />
        </el-form-item>

        <el-form-item label="手机号">
          <el-input v-model="form.phone_number" placeholder="请输入手机号" />
        </el-form-item>

        <el-form-item label="邮箱">
          <el-input v-model="form.email_address" placeholder="请输入邮箱" />
        </el-form-item>

        <el-form-item label="个人描述">
          <el-input
            v-model="form.remark"
            type="textarea"
            :rows="3"
            placeholder="请输入个人描述"
          />
        </el-form-item>
      </el-form>
    </div>

    <template #footer>
      <div class="drawer-footer">
        <el-button @click="handleClose">取消</el-button>
        <el-button type="primary" :loading="saving" @click="handleSave">保存</el-button>
      </div>
    </template>
  </el-drawer>
</template>

<style scoped>
.profile-body {
  padding: 0 4px;
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.avatar-section {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
}

.avatar-wrap {
  position: relative;
  width: 88px;
  height: 88px;
  border-radius: 50%;
  overflow: hidden;
  cursor: pointer;
  border: 2px dashed var(--el-border-color);
  transition: border-color 0.2s;
}

.avatar-wrap:hover {
  border-color: var(--el-color-primary);
}

.avatar-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.avatar-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--el-fill-color-light);
  color: var(--el-text-color-placeholder);
}

.avatar-overlay {
  position: absolute;
  inset: 0;
  background: rgba(0, 0, 0, 0.45);
  color: #fff;
  font-size: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  opacity: 0;
  transition: opacity 0.2s;
}

.avatar-wrap:hover .avatar-overlay {
  opacity: 1;
}

.avatar-hint {
  font-size: 12px;
  color: var(--el-text-color-placeholder);
  margin: 0;
}

.drawer-footer {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
}
</style>
