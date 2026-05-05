<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  CircleCheckFilled,
  CircleCloseFilled,
  Delete,
  EditPen,
  Key,
  Plus,
  Refresh,
  Search,
  User as UserIcon,
} from '@element-plus/icons-vue'
import {
  assignUserRoles,
  createUser,
  deleteUser,
  listRoleOptions,
  listUsers,
  resetUserPassword,
  updateUser,
  updateUserStatus,
} from '@/api/user'
import type {
  CreateUserPayload,
  UpdateUserPayload,
  UserItem,
  UserListQuery,
  UserRole,
} from '@/types/user'

const STATUS_ACTIVE = 1
const STATUS_FROZEN = 2

const loading = ref(false)
const items = ref<UserItem[]>([])
const total = ref(0)
const roleOptions = ref<UserRole[]>([])

const filters = reactive<UserListQuery>({
  page: 1,
  page_size: 10,
  keyword: '',
  user_status: 0,
  role_id: 0,
})

// ── Modal state ───────────────────────────────────────────────────────
const formDialogVisible = ref(false)
const formMode = ref<'create' | 'edit'>('create')
const formSubmitting = ref(false)
const editingId = ref<number | null>(null)
const formState = reactive({
  login_name: '',
  password: '',
  display_name: '',
  email_address: '',
  phone_number: '',
  primary_role_id: undefined as number | undefined,
  user_status: STATUS_ACTIVE,
  remark: '',
  role_ids: [] as number[],
})

const passwordDialogVisible = ref(false)
const passwordSubmitting = ref(false)
const passwordTarget = ref<UserItem | null>(null)
const passwordForm = reactive({
  password: '',
  must_change_password: true,
})

const rolesDialogVisible = ref(false)
const rolesSubmitting = ref(false)
const rolesTarget = ref<UserItem | null>(null)
const rolesForm = reactive({ role_ids: [] as number[] })

// ── Derived ───────────────────────────────────────────────────────────
const isAdminTarget = (row?: UserItem | null) =>
  Boolean(row && (row.id === 1 || row.login_name === 'admin'))

const formTitle = computed(() => (formMode.value === 'create' ? '新增用户' : '编辑用户'))

// ── Loaders ───────────────────────────────────────────────────────────
async function loadData() {
  loading.value = true
  try {
    const result = await listUsers({
      page: filters.page,
      page_size: filters.page_size,
      keyword: filters.keyword?.trim() || undefined,
      user_status: filters.user_status || undefined,
      role_id: filters.role_id || undefined,
    })
    items.value = result.items ?? []
    total.value = result.total
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : '加载用户列表失败')
  } finally {
    loading.value = false
  }
}

async function loadRoleOptions() {
  try {
    roleOptions.value = await listRoleOptions()
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : '加载角色选项失败')
  }
}

onMounted(async () => {
  await Promise.all([loadRoleOptions(), loadData()])
})

// ── Filters ───────────────────────────────────────────────────────────
function applyFilters() {
  filters.page = 1
  void loadData()
}

function resetFilters() {
  filters.keyword = ''
  filters.user_status = 0
  filters.role_id = 0
  filters.page = 1
  void loadData()
}

function onPageChange(page: number) {
  filters.page = page
  void loadData()
}

function onPageSizeChange(size: number) {
  filters.page = 1
  filters.page_size = size
  void loadData()
}

// ── Create / Edit ─────────────────────────────────────────────────────
function resetFormState() {
  formState.login_name = ''
  formState.password = ''
  formState.display_name = ''
  formState.email_address = ''
  formState.phone_number = ''
  formState.primary_role_id = undefined
  formState.user_status = STATUS_ACTIVE
  formState.remark = ''
  formState.role_ids = []
}

function openCreate() {
  resetFormState()
  formMode.value = 'create'
  editingId.value = null
  formDialogVisible.value = true
}

function openEdit(row: UserItem) {
  resetFormState()
  formMode.value = 'edit'
  editingId.value = row.id
  formState.login_name = row.login_name
  formState.display_name = row.display_name
  formState.email_address = row.email_address ?? ''
  formState.phone_number = row.phone_number ?? ''
  formState.primary_role_id = row.primary_role_id ?? undefined
  formState.user_status = row.user_status
  formState.remark = row.remark ?? ''
  formState.role_ids = row.roles?.map((r) => r.id) ?? []
  formDialogVisible.value = true
}

function trimmedOrNull(value: string): string | null {
  const v = value.trim()
  return v ? v : null
}

async function submitForm() {
  if (!formState.login_name.trim() || !formState.display_name.trim()) {
    ElMessage.warning('请填写账号和显示名称')
    return
  }
  if (formMode.value === 'create' && formState.password.length < 6) {
    ElMessage.warning('密码至少 6 位')
    return
  }
  formSubmitting.value = true
  try {
    if (formMode.value === 'create') {
      const payload: CreateUserPayload = {
        login_name: formState.login_name.trim(),
        password: formState.password,
        display_name: formState.display_name.trim(),
        email_address: trimmedOrNull(formState.email_address),
        phone_number: trimmedOrNull(formState.phone_number),
        primary_role_id: formState.primary_role_id ?? null,
        user_status: formState.user_status,
        remark: trimmedOrNull(formState.remark),
        role_ids: formState.role_ids,
      }
      await createUser(payload)
      ElMessage.success('已创建用户')
    } else if (editingId.value) {
      const payload: UpdateUserPayload = {
        display_name: formState.display_name.trim(),
        email_address: trimmedOrNull(formState.email_address),
        phone_number: trimmedOrNull(formState.phone_number),
        primary_role_id: formState.primary_role_id ?? null,
        remark: trimmedOrNull(formState.remark),
      }
      await updateUser(editingId.value, payload)
      ElMessage.success('已更新用户')
    }
    formDialogVisible.value = false
    await loadData()
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : '保存失败')
  } finally {
    formSubmitting.value = false
  }
}

// ── Reset password ────────────────────────────────────────────────────
function openResetPassword(row: UserItem) {
  passwordTarget.value = row
  passwordForm.password = ''
  passwordForm.must_change_password = true
  passwordDialogVisible.value = true
}

async function submitResetPassword() {
  if (!passwordTarget.value) return
  if (passwordForm.password.length < 6) {
    ElMessage.warning('密码至少 6 位')
    return
  }
  passwordSubmitting.value = true
  try {
    await resetUserPassword(passwordTarget.value.id, {
      password: passwordForm.password,
      must_change_password: passwordForm.must_change_password,
    })
    ElMessage.success('已重置密码')
    passwordDialogVisible.value = false
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : '重置失败')
  } finally {
    passwordSubmitting.value = false
  }
}

// ── Assign roles ──────────────────────────────────────────────────────
function openAssignRoles(row: UserItem) {
  rolesTarget.value = row
  rolesForm.role_ids = row.roles?.map((r) => r.id) ?? []
  rolesDialogVisible.value = true
}

async function submitAssignRoles() {
  if (!rolesTarget.value) return
  rolesSubmitting.value = true
  try {
    await assignUserRoles(rolesTarget.value.id, { role_ids: rolesForm.role_ids })
    ElMessage.success('已更新角色')
    rolesDialogVisible.value = false
    await loadData()
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : '保存失败')
  } finally {
    rolesSubmitting.value = false
  }
}

// ── Toggle status / delete ────────────────────────────────────────────
async function toggleStatus(row: UserItem) {
  if (isAdminTarget(row) && row.user_status === STATUS_ACTIVE) {
    ElMessage.warning('内置 admin 用户不可冻结')
    return
  }
  const next = row.user_status === STATUS_ACTIVE ? STATUS_FROZEN : STATUS_ACTIVE
  const action = next === STATUS_FROZEN ? '冻结' : '解冻'
  try {
    await ElMessageBox.confirm(`确定${action}用户「${row.display_name}」吗？`, '操作确认', {
      type: 'warning',
    })
  } catch {
    return
  }
  try {
    await updateUserStatus(row.id, { user_status: next })
    ElMessage.success(`已${action}用户`)
    await loadData()
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : `${action}失败`)
  }
}

async function removeUser(row: UserItem) {
  if (isAdminTarget(row)) {
    ElMessage.warning('内置 admin 用户不可删除')
    return
  }
  try {
    await ElMessageBox.confirm(`确定删除用户「${row.display_name}」吗？该操作为软删除。`, '危险操作', {
      type: 'warning',
      confirmButtonText: '删除',
      confirmButtonClass: 'el-button--danger',
    })
  } catch {
    return
  }
  try {
    await deleteUser(row.id)
    ElMessage.success('已删除用户')
    await loadData()
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : '删除失败')
  }
}

function formatDateTime(value?: string | null) {
  if (!value) return '—'
  const d = new Date(value)
  if (Number.isNaN(d.getTime())) return value
  const pad = (n: number) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}`
}
</script>

<template>
  <div class="user-page">
    <!-- Filter bar -->
    <section class="bp-placeholder filter-bar">
      <el-form inline @submit.prevent>
        <el-form-item label="关键词">
          <el-input
            v-model="filters.keyword"
            placeholder="账号 / 姓名 / 邮箱"
            clearable
            style="width: 220px"
            @keyup.enter="applyFilters"
          />
        </el-form-item>
        <el-form-item label="状态">
          <el-select
            v-model="filters.user_status"
            placeholder="全部"
            clearable
            style="width: 140px"
            @change="applyFilters"
          >
            <el-option :value="0" label="全部" />
            <el-option :value="STATUS_ACTIVE" label="正常" />
            <el-option :value="STATUS_FROZEN" label="冻结" />
          </el-select>
        </el-form-item>
        <el-form-item label="角色">
          <el-select
            v-model="filters.role_id"
            placeholder="全部"
            clearable
            style="width: 200px"
            @change="applyFilters"
          >
            <el-option :value="0" label="全部" />
            <el-option
              v-for="role in roleOptions"
              :key="role.id"
              :value="role.id"
              :label="role.role_name"
            />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button :icon="Search" type="primary" @click="applyFilters">查询</el-button>
          <el-button :icon="Refresh" @click="resetFilters">重置</el-button>
        </el-form-item>
        <el-form-item style="margin-left: auto">
          <el-button :icon="Plus" type="primary" @click="openCreate">新增用户</el-button>
        </el-form-item>
      </el-form>
    </section>

    <!-- Table -->
    <section class="bp-placeholder table-card" v-loading="loading">
      <el-table :data="items" stripe row-key="id" empty-text="暂无数据">
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column prop="login_name" label="账号" min-width="120" />
        <el-table-column prop="display_name" label="显示名称" min-width="140" />
        <el-table-column label="角色" min-width="180">
          <template #default="{ row }">
            <template v-if="row.roles?.length">
              <el-tag
                v-for="role in row.roles"
                :key="role.id"
                size="small"
                style="margin-right: 4px"
              >
                {{ role.role_name }}
              </el-tag>
            </template>
            <span v-else style="color: var(--text-tertiary)">—</span>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag
              :type="row.user_status === STATUS_ACTIVE ? 'success' : 'danger'"
              size="small"
            >
              {{ row.user_status === STATUS_ACTIVE ? '正常' : '冻结' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="email_address" label="邮箱" min-width="180">
          <template #default="{ row }">
            <span>{{ row.email_address || '—' }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="phone_number" label="手机" min-width="120">
          <template #default="{ row }">
            <span>{{ row.phone_number || '—' }}</span>
          </template>
        </el-table-column>
        <el-table-column label="最近登录" min-width="160">
          <template #default="{ row }">
            <span>{{ formatDateTime(row.last_login_at) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="创建时间" min-width="160">
          <template #default="{ row }">
            <span>{{ formatDateTime(row.created_at) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="操作" fixed="right" width="280">
          <template #default="{ row }">
            <el-button :icon="EditPen" link size="small" @click="openEdit(row)">编辑</el-button>
            <el-button :icon="UserIcon" link size="small" @click="openAssignRoles(row)">角色</el-button>
            <el-button :icon="Key" link size="small" @click="openResetPassword(row)">重置密码</el-button>
            <el-button
              :icon="row.user_status === STATUS_ACTIVE ? CircleCloseFilled : CircleCheckFilled"
              link
              size="small"
              :type="row.user_status === STATUS_ACTIVE ? 'warning' : 'success'"
              :disabled="isAdminTarget(row) && row.user_status === STATUS_ACTIVE"
              @click="toggleStatus(row)"
            >
              {{ row.user_status === STATUS_ACTIVE ? '冻结' : '解冻' }}
            </el-button>
            <el-button
              :icon="Delete"
              link
              size="small"
              type="danger"
              :disabled="isAdminTarget(row)"
              @click="removeUser(row)"
            >删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-row">
        <el-pagination
          background
          layout="total, sizes, prev, pager, next, jumper"
          :total="total"
          :page-size="filters.page_size"
          :current-page="filters.page"
          :page-sizes="[10, 20, 50, 100]"
          @current-change="onPageChange"
          @size-change="onPageSizeChange"
        />
      </div>
    </section>

    <!-- Create / Edit dialog -->
    <el-dialog v-model="formDialogVisible" :title="formTitle" width="560px">
      <el-form label-position="top">
        <el-form-item label="账号" required>
          <el-input
            v-model="formState.login_name"
            placeholder="3-32 位字母/数字/_-."
            :disabled="formMode === 'edit'"
            autocomplete="off"
          />
        </el-form-item>
        <el-form-item v-if="formMode === 'create'" label="初始密码" required>
          <el-input
            v-model="formState.password"
            type="password"
            placeholder="至少 6 位"
            show-password
            autocomplete="new-password"
          />
        </el-form-item>
        <el-form-item label="显示名称" required>
          <el-input v-model="formState.display_name" placeholder="例如：张三" />
        </el-form-item>
        <div class="form-row">
          <el-form-item label="邮箱">
            <el-input v-model="formState.email_address" placeholder="user@example.com" />
          </el-form-item>
          <el-form-item label="手机">
            <el-input v-model="formState.phone_number" placeholder="13800000000" />
          </el-form-item>
        </div>
        <div class="form-row">
          <el-form-item label="主角色">
            <el-select
              v-model="formState.primary_role_id"
              placeholder="选择主角色"
              clearable
              style="width: 100%"
            >
              <el-option
                v-for="role in roleOptions"
                :key="role.id"
                :value="role.id"
                :label="role.role_name"
              />
            </el-select>
          </el-form-item>
          <el-form-item v-if="formMode === 'create'" label="状态">
            <el-select v-model="formState.user_status" style="width: 100%">
              <el-option :value="STATUS_ACTIVE" label="正常" />
              <el-option :value="STATUS_FROZEN" label="冻结" />
            </el-select>
          </el-form-item>
        </div>
        <el-form-item v-if="formMode === 'create'" label="授权角色">
          <el-select
            v-model="formState.role_ids"
            multiple
            collapse-tags
            collapse-tags-tooltip
            placeholder="可多选；主角色会自动包含"
            style="width: 100%"
          >
            <el-option
              v-for="role in roleOptions"
              :key="role.id"
              :value="role.id"
              :label="role.role_name"
            />
          </el-select>
          <div class="form-hint">用户的菜单与按钮权限来自这里勾选的角色；主角色会自动并入。</div>
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="formState.remark" type="textarea" :rows="2" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="formDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="formSubmitting" @click="submitForm">保存</el-button>
      </template>
    </el-dialog>

    <!-- Reset password dialog -->
    <el-dialog v-model="passwordDialogVisible" title="重置密码" width="440px">
      <p class="dialog-hint">为「{{ passwordTarget?.display_name }}」（{{ passwordTarget?.login_name }}）设置新密码。</p>
      <el-form label-position="top">
        <el-form-item label="新密码" required>
          <el-input
            v-model="passwordForm.password"
            type="password"
            show-password
            placeholder="至少 6 位"
            autocomplete="new-password"
          />
        </el-form-item>
        <el-form-item>
          <el-checkbox v-model="passwordForm.must_change_password">用户下次登录必须修改密码</el-checkbox>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="passwordDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="passwordSubmitting" @click="submitResetPassword">确认重置</el-button>
      </template>
    </el-dialog>

    <!-- Assign roles dialog -->
    <el-dialog v-model="rolesDialogVisible" title="分配角色" width="440px">
      <p class="dialog-hint">为「{{ rolesTarget?.display_name }}」分配角色。</p>
      <el-checkbox-group v-model="rolesForm.role_ids" class="roles-group">
        <el-checkbox v-for="role in roleOptions" :key="role.id" :value="role.id">
          {{ role.role_name }}
          <span class="role-code">({{ role.role_code }})</span>
        </el-checkbox>
      </el-checkbox-group>
      <template #footer>
        <el-button @click="rolesDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="rolesSubmitting" @click="submitAssignRoles">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.user-page {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.filter-bar :deep(.el-form) {
  display: flex;
  flex-wrap: wrap;
  align-items: flex-end;
  gap: 8px 16px;
  margin: 0;
}

.filter-bar :deep(.el-form-item) {
  margin: 0;
}

.table-card {
  padding: 16px;
}

.pagination-row {
  margin-top: 14px;
  display: flex;
  justify-content: flex-end;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
}

.dialog-hint {
  margin-bottom: 14px;
  color: var(--text-secondary);
  font-size: 13px;
}

.roles-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.role-code {
  margin-left: 4px;
  color: var(--text-tertiary);
  font-size: 12px;
}

.form-hint {
  margin-top: 4px;
  color: var(--text-tertiary);
  font-size: 12px;
  line-height: 1.5;
}
</style>
