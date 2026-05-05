<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { Delete, EditPen, Plus, Refresh, Search } from '@element-plus/icons-vue'
import {
  createAPI,
  createSkipRule,
  deleteAPI,
  deleteSkipRule,
  getAPIGroups,
  listAPIs,
  listSkipRules,
  updateAPI,
} from '@/api/api'
import { APIStatusActive, APIStatusDisabled, HTTP_METHODS } from '@/types/api'
import type {
  APIListQuery,
  APIResource,
  SaveAPIPayload,
  SaveSkipRulePayload,
  SkipRule,
  SkipRuleListQuery,
} from '@/types/api'
import { formatDateTime } from '@/utils/datetime'

// ── API Resources ─────────────────────────────────────────────────────────

const apiLoading = ref(false)
const apiItems = ref<APIResource[]>([])
const apiTotal = ref(0)
const apiGroups = ref<string[]>([])

const apiFilters = reactive<APIListQuery>({
  page: 1,
  page_size: 10,
  keyword: '',
  api_group: '',
  api_method: '',
  api_status: 0,
})

const apiDialogVisible = ref(false)
const apiMode = ref<'create' | 'edit'>('create')
const apiSubmitting = ref(false)
const editingApiId = ref<number | null>(null)
const apiFormRef = ref<FormInstance>()
const apiForm = reactive<SaveAPIPayload>({
  api_path: '',
  api_method: 'GET',
  api_group: '',
  api_desc: '',
  api_status: APIStatusActive,
})
const apiFormRules: FormRules = {
  api_path: [
    { required: true, message: '请填写接口路径', trigger: 'blur' },
    { pattern: /^\//, message: '路径必须以 / 开头', trigger: 'blur' },
  ],
  api_method: [{ required: true, message: '请选择请求方法', trigger: 'change' }],
}

async function loadAPIs() {
  apiLoading.value = true
  try {
    const res = await listAPIs(apiFilters)
    apiItems.value = res.items
    apiTotal.value = res.total
  } catch {
    ElMessage.error('加载失败')
  } finally {
    apiLoading.value = false
  }
}

async function loadGroups() {
  try {
    apiGroups.value = await getAPIGroups()
  } catch {
    // non-critical
  }
}

function handleAPISearch() {
  apiFilters.page = 1
  loadAPIs()
}

function handleAPIReset() {
  apiFilters.keyword = ''
  apiFilters.api_group = ''
  apiFilters.api_method = ''
  apiFilters.api_status = 0
  apiFilters.page = 1
  loadAPIs()
}

function openCreateAPI() {
  apiMode.value = 'create'
  editingApiId.value = null
  Object.assign(apiForm, { api_path: '', api_method: 'GET', api_group: '', api_desc: '', api_status: APIStatusActive })
  apiDialogVisible.value = true
}

function openEditAPI(row: APIResource) {
  apiMode.value = 'edit'
  editingApiId.value = row.id
  Object.assign(apiForm, {
    api_path: row.api_path,
    api_method: row.api_method,
    api_group: row.api_group,
    api_desc: row.api_desc,
    api_status: row.api_status,
  })
  apiDialogVisible.value = true
}

async function submitAPI() {
  const valid = await apiFormRef.value?.validate().catch(() => false)
  if (!valid) return
  apiSubmitting.value = true
  try {
    if (apiMode.value === 'create') {
      await createAPI({ ...apiForm })
      ElMessage.success('创建成功')
    } else {
      await updateAPI(editingApiId.value!, { ...apiForm })
      ElMessage.success('更新成功')
    }
    apiDialogVisible.value = false
    loadAPIs()
    loadGroups()
  } catch (e: any) {
    ElMessage.error(e?.message || '操作失败')
  } finally {
    apiSubmitting.value = false
  }
}

async function handleDeleteAPI(row: APIResource) {
  try {
    await ElMessageBox.confirm(`确认删除 [${row.api_method} ${row.api_path}]？`, '删除确认', {
      confirmButtonText: '删除',
      cancelButtonText: '取消',
      type: 'warning',
    })
    await deleteAPI(row.id)
    ElMessage.success('已删除')
    loadAPIs()
  } catch {
    // cancelled
  }
}

// ── Skip Rules ────────────────────────────────────────────────────────────

const skipLoading = ref(false)
const skipItems = ref<SkipRule[]>([])
const skipTotal = ref(0)

const skipFilters = reactive<SkipRuleListQuery>({
  page: 1,
  page_size: 10,
  keyword: '',
  api_method: '',
})

const skipDialogVisible = ref(false)
const skipSubmitting = ref(false)
const skipFormRef = ref<FormInstance>()
const skipForm = reactive<SaveSkipRulePayload>({
  api_path: '',
  api_method: 'GET',
  skip_reason: '',
})
const skipFormRules: FormRules = {
  api_path: [
    { required: true, message: '请填写接口路径', trigger: 'blur' },
    { pattern: /^\//, message: '路径必须以 / 开头', trigger: 'blur' },
  ],
  api_method: [{ required: true, message: '请选择请求方法', trigger: 'change' }],
}

async function loadSkipRules() {
  skipLoading.value = true
  try {
    const res = await listSkipRules(skipFilters)
    skipItems.value = res.items
    skipTotal.value = res.total
  } catch {
    ElMessage.error('加载失败')
  } finally {
    skipLoading.value = false
  }
}

function handleSkipSearch() {
  skipFilters.page = 1
  loadSkipRules()
}

function handleSkipReset() {
  skipFilters.keyword = ''
  skipFilters.api_method = ''
  skipFilters.page = 1
  loadSkipRules()
}

function openCreateSkip() {
  Object.assign(skipForm, { api_path: '', api_method: 'GET', skip_reason: '' })
  skipDialogVisible.value = true
}

async function submitSkipRule() {
  const valid = await skipFormRef.value?.validate().catch(() => false)
  if (!valid) return
  skipSubmitting.value = true
  try {
    await createSkipRule({ ...skipForm })
    ElMessage.success('创建成功')
    skipDialogVisible.value = false
    loadSkipRules()
  } catch (e: any) {
    ElMessage.error(e?.message || '操作失败')
  } finally {
    skipSubmitting.value = false
  }
}

async function handleDeleteSkip(row: SkipRule) {
  try {
    await ElMessageBox.confirm(`确认删除白名单 [${row.api_method} ${row.api_path}]？`, '删除确认', {
      confirmButtonText: '删除',
      cancelButtonText: '取消',
      type: 'warning',
    })
    await deleteSkipRule(row.id)
    ElMessage.success('已删除')
    loadSkipRules()
  } catch {
    // cancelled
  }
}

// ── Helpers ───────────────────────────────────────────────────────────────

function methodTag(method: string) {
  const map: Record<string, string> = { GET: 'success', POST: 'primary', PUT: 'warning', DELETE: 'danger', PATCH: 'warning' }
  return map[method] ?? 'info'
}

function fetchGroupSuggestions(q: string, cb: (arr: { value: string }[]) => void) {
  cb(apiGroups.value.filter(g => g.includes(q)).map(g => ({ value: g })))
}

onMounted(() => {
  loadAPIs()
  loadGroups()
  loadSkipRules()
})
</script>

<template>
  <div class="page-container">
    <el-tabs type="border-card">
      <!-- ── API 资源 ── -->
      <el-tab-pane label="API 资源">
        <el-card shadow="never" class="filter-card">
          <el-form inline @submit.prevent="handleAPISearch">
            <el-form-item label="关键词">
              <el-input v-model="apiFilters.keyword" placeholder="路径 / 描述" clearable style="width:180px" />
            </el-form-item>
            <el-form-item label="分组">
              <el-select v-model="apiFilters.api_group" clearable style="width:130px">
                <el-option label="全部" value="" />
                <el-option v-for="g in apiGroups" :key="g" :label="g" :value="g" />
              </el-select>
            </el-form-item>
            <el-form-item label="方法">
              <el-select v-model="apiFilters.api_method" clearable style="width:110px">
                <el-option label="全部" value="" />
                <el-option v-for="m in HTTP_METHODS" :key="m" :label="m" :value="m" />
              </el-select>
            </el-form-item>
            <el-form-item label="状态">
              <el-select v-model="apiFilters.api_status" style="width:100px">
                <el-option label="全部" :value="0" />
                <el-option label="启用" :value="APIStatusActive" />
                <el-option label="禁用" :value="APIStatusDisabled" />
              </el-select>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" :icon="Search" @click="handleAPISearch">查询</el-button>
              <el-button :icon="Refresh" @click="handleAPIReset">重置</el-button>
            </el-form-item>
            <el-form-item style="margin-left: auto">
              <el-button type="primary" :icon="Plus" @click="openCreateAPI">新增</el-button>
            </el-form-item>
          </el-form>
        </el-card>

        <el-card shadow="never" class="table-card">
          <el-table v-loading="apiLoading" :data="apiItems" stripe>
            <el-table-column prop="id" label="ID" width="70" />
            <el-table-column label="方法" width="90" align="center">
              <template #default="{ row }">
                <el-tag :type="methodTag(row.api_method)" size="small">{{ row.api_method }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="api_path" label="路径" min-width="240" show-overflow-tooltip />
            <el-table-column prop="api_group" label="分组" width="110" />
            <el-table-column prop="api_desc" label="描述" min-width="160" show-overflow-tooltip />
            <el-table-column label="状态" width="80" align="center">
              <template #default="{ row }">
                <el-tag :type="row.api_status === APIStatusActive ? 'success' : 'info'" size="small">
                  {{ row.api_status === APIStatusActive ? '启用' : '禁用' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="140" fixed="right">
              <template #default="{ row }">
                <el-button link size="small" :icon="EditPen" @click="openEditAPI(row)">编辑</el-button>
                <el-button link size="small" type="danger" :icon="Delete" @click="handleDeleteAPI(row)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>

          <div class="pagination-wrapper">
            <el-pagination
              v-model:current-page="apiFilters.page"
              v-model:page-size="apiFilters.page_size"
              :total="apiTotal"
              :page-sizes="[10, 20, 50, 100]"
              layout="total, sizes, prev, pager, next"
              @current-change="loadAPIs"
              @size-change="() => { apiFilters.page = 1; loadAPIs() }"
            />
          </div>
        </el-card>
      </el-tab-pane>

      <!-- ── 白名单 ── -->
      <el-tab-pane label="免鉴权白名单">
        <el-card shadow="never" class="filter-card">
          <el-form inline @submit.prevent="handleSkipSearch">
            <el-form-item label="关键词">
              <el-input v-model="skipFilters.keyword" placeholder="路径 / 原因" clearable style="width:200px" />
            </el-form-item>
            <el-form-item label="方法">
              <el-select v-model="skipFilters.api_method" clearable style="width:110px">
                <el-option label="全部" value="" />
                <el-option v-for="m in HTTP_METHODS" :key="m" :label="m" :value="m" />
              </el-select>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" :icon="Search" @click="handleSkipSearch">查询</el-button>
              <el-button :icon="Refresh" @click="handleSkipReset">重置</el-button>
            </el-form-item>
            <el-form-item style="margin-left: auto">
              <el-button type="primary" :icon="Plus" @click="openCreateSkip">新增</el-button>
            </el-form-item>
          </el-form>
        </el-card>

        <el-card shadow="never" class="table-card">
          <el-table v-loading="skipLoading" :data="skipItems" stripe>
            <el-table-column prop="id" label="ID" width="70" />
            <el-table-column label="方法" width="90" align="center">
              <template #default="{ row }">
                <el-tag :type="methodTag(row.api_method)" size="small">{{ row.api_method }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="api_path" label="路径" min-width="260" show-overflow-tooltip />
            <el-table-column prop="skip_reason" label="原因" min-width="200" show-overflow-tooltip />
            <el-table-column prop="created_at" label="创建时间" width="180">
              <template #default="{ row }">{{ formatDateTime(row.created_at) }}</template>
            </el-table-column>
            <el-table-column label="操作" width="80" fixed="right">
              <template #default="{ row }">
                <el-button link size="small" type="danger" :icon="Delete" @click="handleDeleteSkip(row)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>

          <div class="pagination-wrapper">
            <el-pagination
              v-model:current-page="skipFilters.page"
              v-model:page-size="skipFilters.page_size"
              :total="skipTotal"
              :page-sizes="[10, 20, 50, 100]"
              layout="total, sizes, prev, pager, next"
              @current-change="loadSkipRules"
              @size-change="() => { skipFilters.page = 1; loadSkipRules() }"
            />
          </div>
        </el-card>
      </el-tab-pane>
    </el-tabs>

    <!-- API 资源 Dialog -->
    <el-dialog
      v-model="apiDialogVisible"
      :title="apiMode === 'create' ? '新增 API 资源' : '编辑 API 资源'"
      width="520px"
      destroy-on-close
    >
      <el-form ref="apiFormRef" :model="apiForm" :rules="apiFormRules" label-width="80px">
        <el-form-item label="路径" prop="api_path">
          <el-input v-model="apiForm.api_path" placeholder="/api/v1/xxx 或 /api/v1/xxx/{id}" />
        </el-form-item>
        <el-form-item label="方法" prop="api_method">
          <el-select v-model="apiForm.api_method" style="width:100%">
            <el-option v-for="m in HTTP_METHODS" :key="m" :label="m" :value="m" />
          </el-select>
        </el-form-item>
        <el-form-item label="分组">
          <el-autocomplete
            v-model="apiForm.api_group"
            :fetch-suggestions="fetchGroupSuggestions"
            placeholder="如 user / role / audit"
            style="width:100%"
          />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="apiForm.api_desc" placeholder="简短说明" />
        </el-form-item>
        <el-form-item label="状态">
          <el-radio-group v-model="apiForm.api_status">
            <el-radio :value="APIStatusActive">启用</el-radio>
            <el-radio :value="APIStatusDisabled">禁用</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="apiDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="apiSubmitting" @click="submitAPI">保存</el-button>
      </template>
    </el-dialog>

    <!-- Skip Rule Dialog -->
    <el-dialog v-model="skipDialogVisible" title="新增白名单规则" width="480px" destroy-on-close>
      <el-form ref="skipFormRef" :model="skipForm" :rules="skipFormRules" label-width="80px">
        <el-form-item label="路径" prop="api_path">
          <el-input v-model="skipForm.api_path" placeholder="/api/v1/xxx 或含 {placeholder}" />
        </el-form-item>
        <el-form-item label="方法" prop="api_method">
          <el-select v-model="skipForm.api_method" style="width:100%">
            <el-option v-for="m in HTTP_METHODS" :key="m" :label="m" :value="m" />
          </el-select>
        </el-form-item>
        <el-form-item label="原因">
          <el-input v-model="skipForm.skip_reason" placeholder="说明为何绕过鉴权" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="skipDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="skipSubmitting" @click="submitSkipRule">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.page-container {
  padding: 16px;
}
.filter-card {
  margin-bottom: 12px;
}
.filter-card :deep(.el-card__body) {
  padding-bottom: 0;
}
.table-card {
  margin-top: 0;
}
.pagination-wrapper {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}
</style>
