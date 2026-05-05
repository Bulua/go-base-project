<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Delete, EditPen, Plus, Refresh } from '@element-plus/icons-vue'
import IconPicker from '@/components/common/IconPicker.vue'
import {
  type CreateMenuResult,
  createMenu,
  createMenuParam,
  createMenuAction,
  deleteMenu,
  deleteMenuParam,
  deleteMenuAction,
  getMenuTree,
  listMenuActions,
  updateMenu,
  updateMenuAction,
} from '@/api/menu'
import type { CreateParamPayload, MenuItem, MenuAction, RouteParam, SaveActionPayload, SaveMenuPayload } from '@/types/menu'

// ── 菜单类型 ───────────────────────────────────────────────────────────────
const MENU_TYPES = [
  { value: 1, label: '目录', type: 'primary' },
  { value: 2, label: '菜单', type: 'success' },
  { value: 3, label: '隐藏路由', type: 'warning' },
  { value: 4, label: '外链', type: 'info' },
] as const

function menuTypeLabel(type: number) {
  return MENU_TYPES.find((t) => t.value === type)?.label ?? '未知'
}
function menuTypeTag(type: number) {
  return MENU_TYPES.find((t) => t.value === type)?.type ?? 'info'
}

// ── 数据 ──────────────────────────────────────────────────────────────────
const loading = ref(false)
const treeData = ref<MenuItem[]>([])

async function loadTree() {
  loading.value = true
  try {
    treeData.value = await getMenuTree()
  } catch (e) {
    ElMessage.error(e instanceof Error ? e.message : '加载失败')
  } finally {
    loading.value = false
  }
}

onMounted(loadTree)

// ── 删除 ──────────────────────────────────────────────────────────────────
async function handleDelete(row: MenuItem) {
  await ElMessageBox.confirm(`确定删除菜单「${row.menu_title}」？`, '删除确认', {
    type: 'warning',
    confirmButtonText: '删除',
    cancelButtonText: '取消',
    confirmButtonClass: 'el-button--danger',
  })
  try {
    await deleteMenu(row.id)
    ElMessage.success('删除成功')
    await loadTree()
  } catch (e) {
    ElMessage.error(e instanceof Error ? e.message : '删除失败')
  }
}

// ── 表单 Dialog ───────────────────────────────────────────────────────────
const dialogVisible = ref(false)
const dialogTitle = ref('')
const editingId = ref<number | null>(null)
const submitting = ref(false)

const EMPTY_FORM = (): SaveMenuPayload => ({
  parent_id: 0,
  menu_type: 2,
  route_path: null,
  route_name: null,
  component_path: null,
  redirect_path: null,
  menu_title: '',
  menu_icon: null,
  sort_no: 0,
  is_hidden: false,
  is_keep_alive: false,
  is_affix: false,
  active_route: null,
  transition_name: null,
  external_url: null,
  menu_status: 1,
})

const form = reactive<SaveMenuPayload>(EMPTY_FORM())

function openCreate(parentId = 0) {
  Object.assign(form, EMPTY_FORM())
  form.parent_id = parentId
  editingId.value = null
  dialogTitle.value = parentId ? '新增子菜单' : '新增根菜单'
  dialogVisible.value = true
  paramRows.value = []
}

function openEdit(row: MenuItem) {
  Object.assign(form, {
    parent_id: row.parent_id,
    menu_type: row.menu_type,
    route_path: row.route_path ?? null,
    route_name: row.route_name ?? null,
    component_path: row.component_path ?? null,
    redirect_path: row.redirect_path ?? null,
    menu_title: row.menu_title,
    menu_icon: row.menu_icon ?? null,
    sort_no: row.sort_no,
    is_hidden: row.is_hidden,
    is_keep_alive: row.is_keep_alive,
    is_affix: row.is_affix,
    active_route: row.active_route ?? null,
    transition_name: row.transition_name ?? null,
    external_url: row.external_url ?? null,
    menu_status: row.menu_status,
  })
  editingId.value = row.id
  dialogTitle.value = '编辑菜单'
  dialogVisible.value = true
  paramRows.value = row.params ? [...row.params] : []
}

async function handleSubmit() {
  if (!form.menu_title.trim()) {
    ElMessage.warning('请填写菜单名称')
    return
  }
  submitting.value = true
  try {
    if (editingId.value) {
      await updateMenu(editingId.value, form)
      ElMessage.success('修改成功')
    } else {
      const result: CreateMenuResult = await createMenu(form)
      if (result.code_generated) {
        ElMessage.success(`创建成功，已生成组件文件：${result.code_path}`)
      } else {
        ElMessage.success('创建成功')
      }
    }
    dialogVisible.value = false
    await loadTree()
  } catch (e) {
    ElMessage.error(e instanceof Error ? e.message : '保存失败')
  } finally {
    submitting.value = false
  }
}

// ── 路由参数 ──────────────────────────────────────────────────────────────
const paramRows = ref<RouteParam[]>([])
const addingParam = ref(false)
const newParam = reactive<CreateParamPayload>({ param_mode: 'query', param_key: '', param_value: '' })

async function handleAddParam() {
  if (!newParam.param_key.trim()) {
    ElMessage.warning('请填写参数键名')
    return
  }
  if (!editingId.value) return
  addingParam.value = true
  try {
    const created = await createMenuParam(editingId.value, {
      param_mode: newParam.param_mode,
      param_key: newParam.param_key.trim(),
      param_value: newParam.param_value || null,
    })
    paramRows.value.push(created)
    newParam.param_key = ''
    newParam.param_value = ''
  } catch (e) {
    ElMessage.error(e instanceof Error ? e.message : '添加失败')
  } finally {
    addingParam.value = false
  }
}

async function handleDeleteParam(param: RouteParam) {
  try {
    await deleteMenuParam(param.id)
    paramRows.value = paramRows.value.filter((p) => p.id !== param.id)
  } catch (e) {
    ElMessage.error(e instanceof Error ? e.message : '删除失败')
  }
}

// ── 按钮权限管理 Drawer ───────────────────────────────────────────────────
const actionDrawerVisible = ref(false)
const actionMenuId = ref(0)
const actionMenuTitle = ref('')
const actionRows = ref<MenuAction[]>([])
const actionLoading = ref(false)
const actionSubmitting = ref(false)

const EMPTY_ACTION = (): SaveActionPayload => ({
  action_code: '',
  action_name: '',
  action_desc: null,
  sort_no: 0,
  action_status: 1,
})
const actionForm = reactive<SaveActionPayload>(EMPTY_ACTION())
const editingActionId = ref<number | null>(null)
const actionFormVisible = ref(false)

async function openActionDrawer(row: MenuItem) {
  actionMenuId.value = row.id
  actionMenuTitle.value = row.menu_title
  actionDrawerVisible.value = true
  actionLoading.value = true
  try {
    actionRows.value = await listMenuActions(row.id)
  } catch (e) {
    ElMessage.error(e instanceof Error ? e.message : '加载失败')
  } finally {
    actionLoading.value = false
  }
}

function openActionCreate() {
  Object.assign(actionForm, EMPTY_ACTION())
  editingActionId.value = null
  actionFormVisible.value = true
}

function openActionEdit(action: MenuAction) {
  Object.assign(actionForm, {
    action_code: action.action_code,
    action_name: action.action_name,
    action_desc: action.action_desc ?? null,
    sort_no: action.sort_no,
    action_status: action.action_status,
  })
  editingActionId.value = action.id
  actionFormVisible.value = true
}

async function handleActionSubmit() {
  if (!actionForm.action_code.trim()) { ElMessage.warning('请填写按钮编码'); return }
  if (!actionForm.action_name.trim()) { ElMessage.warning('请填写按钮名称'); return }
  actionSubmitting.value = true
  try {
    if (editingActionId.value) {
      const updated = await updateMenuAction(editingActionId.value, actionForm)
      const idx = actionRows.value.findIndex((a) => a.id === editingActionId.value)
      if (idx !== -1) actionRows.value[idx] = updated
      ElMessage.success('修改成功')
    } else {
      const created = await createMenuAction(actionMenuId.value, actionForm)
      actionRows.value.push(created)
      ElMessage.success('创建成功')
    }
    actionFormVisible.value = false
  } catch (e) {
    ElMessage.error(e instanceof Error ? e.message : '保存失败')
  } finally {
    actionSubmitting.value = false
  }
}

async function handleActionDelete(action: MenuAction) {
  await ElMessageBox.confirm(`确定删除按钮「${action.action_name}」？`, '删除确认', {
    type: 'warning', confirmButtonText: '删除', cancelButtonText: '取消',
    confirmButtonClass: 'el-button--danger',
  })
  try {
    await deleteMenuAction(action.id)
    actionRows.value = actionRows.value.filter((a) => a.id !== action.id)
    ElMessage.success('删除成功')
  } catch (e) {
    ElMessage.error(e instanceof Error ? e.message : '删除失败')
  }
}
</script>

<template>
  <div style="padding: 16px">
    <!-- 工具栏 -->
    <div style="margin-bottom: 12px; display: flex; gap: 8px">
      <el-button type="primary" :icon="Plus" @click="openCreate(0)">新增根菜单</el-button>
      <el-button :icon="Refresh" @click="loadTree">刷新</el-button>
    </div>

    <!-- 树形表格 -->
    <el-table
      v-loading="loading"
      :data="treeData"
      row-key="id"
      :tree-props="{ children: 'children' }"
      border
      default-expand-all
    >
      <el-table-column prop="menu_title" label="菜单名称" min-width="180">
        <template #default="{ row }">
          <el-icon v-if="row.menu_icon" style="margin-right: 4px; vertical-align: middle">
            <component :is="row.menu_icon" />
          </el-icon>
          {{ row.menu_title }}
        </template>
      </el-table-column>

      <el-table-column label="类型" width="90" align="center">
        <template #default="{ row }">
          <el-tag :type="menuTypeTag(row.menu_type)" size="small">
            {{ menuTypeLabel(row.menu_type) }}
          </el-tag>
        </template>
      </el-table-column>

      <el-table-column prop="route_path" label="路由路径" min-width="140" show-overflow-tooltip />

      <el-table-column prop="component_path" label="组件路径" min-width="200" show-overflow-tooltip />

      <el-table-column prop="sort_no" label="排序" width="70" align="center" />

      <el-table-column label="隐藏" width="70" align="center">
        <template #default="{ row }">
          <el-tag v-if="row.is_hidden" type="info" size="small">隐藏</el-tag>
          <span v-else style="color: #909399">—</span>
        </template>
      </el-table-column>

      <el-table-column label="状态" width="80" align="center">
        <template #default="{ row }">
          <el-tag :type="row.menu_status === 1 ? 'success' : 'danger'" size="small">
            {{ row.menu_status === 1 ? '启用' : '禁用' }}
          </el-tag>
        </template>
      </el-table-column>

      <el-table-column label="操作" width="240" fixed="right">
        <template #default="{ row }">
          <el-button
            v-if="row.menu_type === 1"
            type="primary"
            link
            size="small"
            :icon="Plus"
            @click="openCreate(row.id)"
          >子菜单</el-button>
          <el-button
            v-if="row.menu_type === 2 || row.menu_type === 3"
            type="warning"
            link
            size="small"
            @click="openActionDrawer(row)"
          >按钮</el-button>
          <el-button type="primary" link size="small" :icon="EditPen" @click="openEdit(row)">
            编辑
          </el-button>
          <el-button type="danger" link size="small" :icon="Delete" @click="handleDelete(row)">
            删除
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 编辑/新增 Dialog -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="680px"
      destroy-on-close
    >
      <el-form :model="form" label-width="100px" label-position="right">
        <el-row :gutter="16">
          <el-col :span="12">
            <el-form-item label="菜单名称" required>
              <el-input v-model="form.menu_title" placeholder="请输入菜单名称" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="菜单类型">
              <el-select v-model="form.menu_type" style="width: 100%">
                <el-option v-for="t in MENU_TYPES" :key="t.value" :label="t.label" :value="t.value" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="16">
          <el-col :span="12">
            <el-form-item label="路由路径">
              <el-input v-model="form.route_path" placeholder="如 /system 或 user" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="路由名称">
              <el-input v-model="form.route_name" placeholder="如 systemUser" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item v-if="form.menu_type !== 4" label="组件路径">
          <el-input
            v-model="form.component_path"
            placeholder="如 views/system/user/index.vue"
          />
        </el-form-item>

        <el-form-item v-if="form.menu_type === 1" label="重定向">
          <el-input v-model="form.redirect_path" placeholder="如 /dashboard/workbench" />
        </el-form-item>

        <el-form-item v-if="form.menu_type === 4" label="外链地址">
          <el-input v-model="form.external_url" placeholder="https://..." />
        </el-form-item>

        <el-row :gutter="16">
          <el-col :span="12">
            <el-form-item label="菜单图标">
              <IconPicker v-model="form.menu_icon" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="排序号">
              <el-input-number v-model="form.sort_no" :min="0" :max="9999" style="width: 100%" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="16">
          <el-col :span="8">
            <el-form-item label="状态">
              <el-select v-model="form.menu_status" style="width: 100%">
                <el-option :value="1" label="启用" />
                <el-option :value="2" label="禁用" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="16">
            <el-form-item label=" ">
              <el-checkbox v-model="form.is_hidden">隐藏路由</el-checkbox>
              <el-checkbox v-model="form.is_keep_alive">Keep Alive</el-checkbox>
              <el-checkbox v-model="form.is_affix">固定标签</el-checkbox>
            </el-form-item>
          </el-col>
        </el-row>

        <!-- 路由参数（仅编辑时显示） -->
        <template v-if="editingId">
          <el-divider content-position="left">路由参数</el-divider>

          <el-table :data="paramRows" size="small" style="margin-bottom: 8px">
            <el-table-column prop="param_mode" label="模式" width="80" />
            <el-table-column prop="param_key" label="键名" />
            <el-table-column prop="param_value" label="默认值" />
            <el-table-column label="" width="60" align="center">
              <template #default="{ row }">
                <el-button type="danger" link size="small" :icon="Delete" @click="handleDeleteParam(row)" />
              </template>
            </el-table-column>
          </el-table>

          <el-row :gutter="8" align="middle">
            <el-col :span="6">
              <el-select v-model="newParam.param_mode" size="small">
                <el-option value="query" label="query" />
                <el-option value="path" label="path" />
              </el-select>
            </el-col>
            <el-col :span="8">
              <el-input v-model="newParam.param_key" size="small" placeholder="键名" />
            </el-col>
            <el-col :span="7">
              <el-input v-model="newParam.param_value" size="small" placeholder="默认值（可空）" />
            </el-col>
            <el-col :span="3">
              <el-button size="small" :loading="addingParam" @click="handleAddParam">添加</el-button>
            </el-col>
          </el-row>
        </template>
      </el-form>

      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="handleSubmit">保存</el-button>
      </template>
    </el-dialog>
    <!-- 按钮权限 Drawer -->
    <el-drawer
      v-model="actionDrawerVisible"
      :title="`按钮权限 — ${actionMenuTitle}`"
      size="520px"
      destroy-on-close
    >
      <div style="margin-bottom: 12px">
        <el-button type="primary" size="small" :icon="Plus" @click="openActionCreate">新增按钮</el-button>
      </div>
      <el-table v-loading="actionLoading" :data="actionRows" border size="small">
        <el-table-column prop="action_code" label="编码" width="120" />
        <el-table-column prop="action_name" label="名称" />
        <el-table-column prop="sort_no" label="排序" width="60" align="center" />
        <el-table-column label="状态" width="70" align="center">
          <template #default="{ row }">
            <el-tag :type="row.action_status === 1 ? 'success' : 'danger'" size="small">
              {{ row.action_status === 1 ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="100" align="center">
          <template #default="{ row }">
            <el-button type="primary" link size="small" :icon="EditPen" @click="openActionEdit(row)" />
            <el-button type="danger" link size="small" :icon="Delete" @click="handleActionDelete(row)" />
          </template>
        </el-table-column>
      </el-table>
    </el-drawer>

    <!-- 按钮新增/编辑 Dialog -->
    <el-dialog
      v-model="actionFormVisible"
      :title="editingActionId ? '编辑按钮' : '新增按钮'"
      width="420px"
      append-to-body
      destroy-on-close
    >
      <el-form :model="actionForm" label-width="80px">
        <el-form-item label="按钮编码" required>
          <el-input v-model="actionForm.action_code" placeholder="如 add / edit / delete" />
        </el-form-item>
        <el-form-item label="按钮名称" required>
          <el-input v-model="actionForm.action_name" placeholder="如 新增 / 编辑 / 删除" />
        </el-form-item>
        <el-form-item label="说明">
          <el-input v-model="actionForm.action_desc" placeholder="可选" />
        </el-form-item>
        <el-row :gutter="16">
          <el-col :span="12">
            <el-form-item label="排序号">
              <el-input-number v-model="actionForm.sort_no" :min="0" :max="9999" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="状态">
              <el-select v-model="actionForm.action_status" style="width: 100%">
                <el-option :value="1" label="启用" />
                <el-option :value="2" label="禁用" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <template #footer>
        <el-button @click="actionFormVisible = false">取消</el-button>
        <el-button type="primary" :loading="actionSubmitting" @click="handleActionSubmit">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>
