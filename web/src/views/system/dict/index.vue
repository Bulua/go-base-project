<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Delete, EditPen, Plus, Refresh, List } from '@element-plus/icons-vue'
import {
  listDicts,
  createDict,
  updateDict,
  deleteDict,
  getDictItems,
  createDictItem,
  updateDictItem,
  deleteDictItem,
} from '@/api/dict'
import type { Dictionary, DictItem, DictListQuery, SaveDictPayload, SaveItemPayload } from '@/types/dict'

const STATUS_ACTIVE = 1
const STATUS_DISABLED = 2

// ── 字典列表 ───────────────────────────────────────────────────────────────
const loading = ref(false)
const items = ref<Dictionary[]>([])
const total = ref(0)

const filters = reactive<DictListQuery>({ page: 1, page_size: 10, keyword: '', dict_status: 0 })

async function loadData() {
  loading.value = true
  try {
    const result = await listDicts({
      page: filters.page,
      page_size: filters.page_size,
      keyword: filters.keyword?.trim() || undefined,
      dict_status: filters.dict_status || undefined,
    })
    items.value = result.items ?? []
    total.value = result.total
  } catch (e) {
    ElMessage.error(e instanceof Error ? e.message : '加载失败')
  } finally {
    loading.value = false
  }
}

onMounted(loadData)

function applyFilters() { filters.page = 1; void loadData() }
function resetFilters() { filters.keyword = ''; filters.dict_status = 0; filters.page = 1; void loadData() }
function onPageChange(p: number) { filters.page = p; void loadData() }
function onPageSizeChange(s: number) { filters.page = 1; filters.page_size = s; void loadData() }

// ── 字典表单 Dialog ────────────────────────────────────────────────────────
const formVisible = ref(false)
const formTitle = ref('')
const editingId = ref<number | null>(null)
const submitting = ref(false)

const EMPTY_FORM = (): SaveDictPayload => ({
  dict_name: '', dict_code: '', dict_status: STATUS_ACTIVE, parent_id: 0, remark: null,
})
const form = reactive<SaveDictPayload>(EMPTY_FORM())

function openCreate() {
  Object.assign(form, EMPTY_FORM())
  editingId.value = null
  formTitle.value = '新增字典'
  formVisible.value = true
}

function openEdit(row: Dictionary) {
  Object.assign(form, {
    dict_name: row.dict_name,
    dict_code: row.dict_code,
    dict_status: row.dict_status,
    parent_id: row.parent_id,
    remark: row.remark ?? null,
  })
  editingId.value = row.id
  formTitle.value = '编辑字典'
  formVisible.value = true
}

async function handleSubmit() {
  if (!form.dict_name.trim()) { ElMessage.warning('请填写字典名称'); return }
  if (!form.dict_code.trim()) { ElMessage.warning('请填写字典编码'); return }
  submitting.value = true
  try {
    if (editingId.value) {
      await updateDict(editingId.value, form)
      ElMessage.success('修改成功')
    } else {
      await createDict(form)
      ElMessage.success('创建成功')
    }
    formVisible.value = false
    await loadData()
  } catch (e) {
    ElMessage.error(e instanceof Error ? e.message : '保存失败')
  } finally {
    submitting.value = false
  }
}

async function handleDelete(row: Dictionary) {
  try {
    await ElMessageBox.confirm(`确定删除字典「${row.dict_name}」及其所有字典项？`, '删除确认', {
      type: 'warning', confirmButtonText: '删除', confirmButtonClass: 'el-button--danger',
    })
  } catch {
    return
  }
  try {
    await deleteDict(row.id)
    ElMessage.success('删除成功')
    await loadData()
  } catch (e) {
    ElMessage.error(e instanceof Error ? e.message : '删除失败')
  }
}

// ── 字典项 Drawer ──────────────────────────────────────────────────────────
const drawerVisible = ref(false)
const drawerDict = ref<Dictionary | null>(null)
const itemRows = ref<DictItem[]>([])
const itemLoading = ref(false)

async function openItemDrawer(row: Dictionary) {
  drawerDict.value = row
  drawerVisible.value = true
  itemLoading.value = true
  try {
    itemRows.value = await getDictItems(row.dict_code)
  } catch (e) {
    ElMessage.error(e instanceof Error ? e.message : '加载字典项失败')
  } finally {
    itemLoading.value = false
  }
}

// ── 字典项表单 Dialog ──────────────────────────────────────────────────────
const itemFormVisible = ref(false)
const itemFormTitle = ref('')
const editingItemId = ref<number | null>(null)
const itemSubmitting = ref(false)

const EMPTY_ITEM = (): SaveItemPayload => ({
  item_label: '', item_value: '', item_extra: null, item_status: STATUS_ACTIVE, sort_no: 0, parent_id: 0,
})
const itemForm = reactive<SaveItemPayload>(EMPTY_ITEM())

function openItemCreate() {
  Object.assign(itemForm, EMPTY_ITEM())
  editingItemId.value = null
  itemFormTitle.value = '新增字典项'
  itemFormVisible.value = true
}

function openItemEdit(item: DictItem) {
  Object.assign(itemForm, {
    item_label: item.item_label,
    item_value: item.item_value,
    item_extra: item.item_extra ?? null,
    item_status: item.item_status,
    sort_no: item.sort_no,
    parent_id: item.parent_id,
  })
  editingItemId.value = item.id
  itemFormTitle.value = '编辑字典项'
  itemFormVisible.value = true
}

async function handleItemSubmit() {
  if (!itemForm.item_label.trim()) { ElMessage.warning('请填写标签'); return }
  if (!itemForm.item_value.trim()) { ElMessage.warning('请填写字典值'); return }
  if (!drawerDict.value) return
  itemSubmitting.value = true
  try {
    if (editingItemId.value) {
      const updated = await updateDictItem(editingItemId.value, itemForm)
      const idx = itemRows.value.findIndex((i) => i.id === editingItemId.value)
      if (idx >= 0) itemRows.value[idx] = updated
      ElMessage.success('修改成功')
    } else {
      const created = await createDictItem(drawerDict.value.id, itemForm)
      itemRows.value.push(created)
      ElMessage.success('创建成功')
    }
    itemFormVisible.value = false
  } catch (e) {
    ElMessage.error(e instanceof Error ? e.message : '保存失败')
  } finally {
    itemSubmitting.value = false
  }
}

async function handleItemDelete(item: DictItem) {
  try {
    await ElMessageBox.confirm(`确定删除字典项「${item.item_label}」？`, '删除确认', {
      type: 'warning', confirmButtonText: '删除', confirmButtonClass: 'el-button--danger',
    })
    await deleteDictItem(item.id)
    itemRows.value = itemRows.value.filter((i) => i.id !== item.id)
    ElMessage.success('删除成功')
  } catch (e) {
    if (e !== 'cancel') ElMessage.error(e instanceof Error ? e.message : '删除失败')
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
  <div class="dict-page">
    <!-- 筛选栏 -->
    <section class="bp-placeholder filter-bar">
      <el-form inline @submit.prevent>
        <el-form-item label="关键词">
          <el-input
            v-model="filters.keyword"
            placeholder="名称 / 编码"
            clearable
            style="width: 220px"
            @keyup.enter="applyFilters"
          />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="filters.dict_status" placeholder="全部" clearable style="width: 120px" @change="applyFilters">
            <el-option :value="0" label="全部" />
            <el-option :value="STATUS_ACTIVE" label="启用" />
            <el-option :value="STATUS_DISABLED" label="禁用" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="applyFilters">查询</el-button>
          <el-button :icon="Refresh" @click="resetFilters">重置</el-button>
        </el-form-item>
        <el-form-item style="margin-left: auto">
          <el-button :icon="Plus" type="primary" @click="openCreate">新增字典</el-button>
        </el-form-item>
      </el-form>
    </section>

    <!-- 字典表格 -->
    <section class="bp-placeholder table-card" v-loading="loading">
      <el-table :data="items" stripe empty-text="暂无数据">
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column prop="dict_name" label="字典名称" min-width="140" />
        <el-table-column prop="dict_code" label="字典编码" min-width="160">
          <template #default="{ row }">
            <span class="code-text">{{ row.dict_code }}</span>
          </template>
        </el-table-column>
        <el-table-column label="字典项数" width="90">
          <template #default="{ row }">
            <el-tag size="small" type="info">{{ row.item_count }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="90">
          <template #default="{ row }">
            <el-tag :type="row.dict_status === STATUS_ACTIVE ? 'success' : 'danger'" size="small">
              {{ row.dict_status === STATUS_ACTIVE ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="remark" label="备注" min-width="160" show-overflow-tooltip />
        <el-table-column label="更新时间" min-width="160">
          <template #default="{ row }">{{ formatDateTime(row.updated_at) }}</template>
        </el-table-column>
        <el-table-column label="操作" fixed="right" width="200">
          <template #default="{ row }">
            <el-button :icon="List" link size="small" @click="openItemDrawer(row)">字典项</el-button>
            <el-button :icon="EditPen" link size="small" @click="openEdit(row)">编辑</el-button>
            <el-button :icon="Delete" link size="small" type="danger" @click="handleDelete(row)">删除</el-button>
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
          :page-sizes="[10, 20, 50]"
          @current-change="onPageChange"
          @size-change="onPageSizeChange"
        />
      </div>
    </section>

    <!-- 字典表单 Dialog -->
    <el-dialog v-model="formVisible" :title="formTitle" width="480px">
      <el-form label-position="top">
        <div class="form-row">
          <el-form-item label="字典名称" required>
            <el-input v-model="form.dict_name" placeholder="例如：用户状态" />
          </el-form-item>
          <el-form-item label="字典编码" required>
            <el-input v-model="form.dict_code" placeholder="例如：user_status" :disabled="editingId !== null" />
          </el-form-item>
        </div>
        <div class="form-row">
          <el-form-item label="状态">
            <el-select v-model="form.dict_status" style="width: 100%">
              <el-option :value="STATUS_ACTIVE" label="启用" />
              <el-option :value="STATUS_DISABLED" label="禁用" />
            </el-select>
          </el-form-item>
        </div>
        <el-form-item label="备注">
          <el-input v-model="form.remark" type="textarea" :rows="2" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="formVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="handleSubmit">保存</el-button>
      </template>
    </el-dialog>

    <!-- 字典项 Drawer -->
    <el-drawer
      v-model="drawerVisible"
      :title="drawerDict ? `字典项：${drawerDict.dict_name} (${drawerDict.dict_code})` : '字典项'"
      size="680px"
      direction="rtl"
    >
      <div class="drawer-toolbar">
        <el-button :icon="Plus" type="primary" size="small" @click="openItemCreate">新增字典项</el-button>
      </div>

      <el-table :data="itemRows" v-loading="itemLoading" stripe empty-text="暂无字典项" size="small">
        <el-table-column prop="item_label" label="标签" min-width="120" />
        <el-table-column prop="item_value" label="值" min-width="120">
          <template #default="{ row }">
            <span class="code-text">{{ row.item_value }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="item_extra" label="扩展" min-width="100" show-overflow-tooltip>
          <template #default="{ row }">{{ row.item_extra ?? '—' }}</template>
        </el-table-column>
        <el-table-column prop="sort_no" label="排序" width="70" />
        <el-table-column label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="row.item_status === STATUS_ACTIVE ? 'success' : 'danger'" size="small">
              {{ row.item_status === STATUS_ACTIVE ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="120" fixed="right">
          <template #default="{ row }">
            <el-button :icon="EditPen" link size="small" @click="openItemEdit(row)">编辑</el-button>
            <el-button :icon="Delete" link size="small" type="danger" @click="handleItemDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-drawer>

    <!-- 字典项表单 Dialog -->
    <el-dialog v-model="itemFormVisible" :title="itemFormTitle" width="480px">
      <el-form label-position="top">
        <div class="form-row">
          <el-form-item label="标签" required>
            <el-input v-model="itemForm.item_label" placeholder="例如：正常" />
          </el-form-item>
          <el-form-item label="字典值" required>
            <el-input v-model="itemForm.item_value" placeholder="例如：1" />
          </el-form-item>
        </div>
        <div class="form-row">
          <el-form-item label="排序">
            <el-input-number v-model="itemForm.sort_no" :min="0" :max="9999" style="width: 100%" />
          </el-form-item>
          <el-form-item label="状态">
            <el-select v-model="itemForm.item_status" style="width: 100%">
              <el-option :value="STATUS_ACTIVE" label="启用" />
              <el-option :value="STATUS_DISABLED" label="禁用" />
            </el-select>
          </el-form-item>
        </div>
        <el-form-item label="扩展信息">
          <el-input v-model="itemForm.item_extra" placeholder="可选，如颜色标识 success / danger" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="itemFormVisible = false">取消</el-button>
        <el-button type="primary" :loading="itemSubmitting" @click="handleItemSubmit">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.dict-page {
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

.filter-bar :deep(.el-form-item) { margin: 0; }

.table-card { padding: 16px; }

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

.code-text {
  font-family: var(--font-mono, monospace);
  font-size: 13px;
}

.drawer-toolbar {
  margin-bottom: 14px;
}
</style>
