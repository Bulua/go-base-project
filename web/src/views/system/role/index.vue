<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import type { ElTree } from 'element-plus'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Delete,
  EditPen,
  Key,
  Plus,
  Refresh,
  Search,
} from '@element-plus/icons-vue'
import {
  assignRoleActions,
  assignRoleApis,
  assignRoleDataScopes,
  assignRoleMenus,
  createRole,
  deleteRole,
  getRoleActionIDs,
  getRoleApiIDs,
  getRoleDataScopeIDs,
  getRoleMenuIDs,
  getRoleResources,
  listRoles,
  updateRole,
} from '@/api/role'
import type {
  ActionOption,
  ApiOption,
  CreateRolePayload,
  RoleItem,
  RoleListQuery,
  RoleResources,
  UpdateRolePayload,
} from '@/types/role'
import { useDict } from '@/composables/useDict'

const STATUS_ACTIVE = 1
const STATUS_DISABLED = 2

const roleStatusDict = useDict('role_status')

const loading = ref(false)
const items = ref<RoleItem[]>([])
const total = ref(0)
const resources = ref<RoleResources>({ menus: [], actions: [], apis: [], roles: [] })

const filters = reactive<RoleListQuery>({
  page: 1,
  page_size: 10,
  keyword: '',
  role_status: 0,
})

// ── Form state ──────────────────────────────────────────────────────
const formDialogVisible = ref(false)
const formMode = ref<'create' | 'edit'>('create')
const formSubmitting = ref(false)
const editingId = ref<number | null>(null)
const formState = reactive({
  role_code: '',
  role_name: '',
  parent_role_id: 0,
  default_route: 'dashboard',
  sort_no: 0,
  role_status: STATUS_ACTIVE,
  remark: '',
})

// ── Permissions dialog state ────────────────────────────────────────
const permsDialogVisible = ref(false)
const permsLoading = ref(false)
const permsSubmitting = ref(false)
const permsTarget = ref<RoleItem | null>(null)
const permsTab = ref<'menus' | 'actions' | 'apis' | 'dataScopes'>('menus')
const menuTreeRef = ref<InstanceType<typeof ElTree> | null>(null)
const selectedMenuIDs = ref<number[]>([])
const selectedActionIDs = ref<number[]>([])
const selectedApiIDs = ref<number[]>([])
const selectedDataScopeIDs = ref<number[]>([])

// 内置角色的判定锚定在 seed 写入的 id=1，与 role_code 字段无关——
// 这样改名/改 code 不会导致保护逻辑失效。
const BUILTIN_ROLE_ID = 1
const isBuiltinRole = (row?: RoleItem | null) => Boolean(row && row.id === BUILTIN_ROLE_ID)

const formTitle = computed(() => (formMode.value === 'create' ? '新增角色' : '编辑角色'))

const parentOptions = computed(() => {
  // exclude the role being edited and its descendants to prevent loops
  if (formMode.value !== 'edit' || editingId.value == null) return resources.value.roles
  const blocked = collectDescendants(resources.value.roles, editingId.value)
  blocked.add(editingId.value)
  return resources.value.roles.filter((r) => !blocked.has(r.id))
})

function collectDescendants(roles: RoleItem[], rootID: number, acc = new Set<number>()): Set<number> {
  // Build child map first
  const childrenByParent = new Map<number, RoleItem[]>()
  roles.forEach((r) => {
    const arr = childrenByParent.get(r.parent_role_id) ?? []
    arr.push(r)
    childrenByParent.set(r.parent_role_id, arr)
  })
  const stack = [rootID]
  while (stack.length) {
    const id = stack.pop() as number
    const kids = childrenByParent.get(id) ?? []
    for (const k of kids) {
      if (!acc.has(k.id)) {
        acc.add(k.id)
        stack.push(k.id)
      }
    }
  }
  return acc
}

const apisGrouped = computed(() => {
  const groups = new Map<string, ApiOption[]>()
  for (const api of resources.value.apis) {
    const key = api.api_group || '其他'
    const arr = groups.get(key) ?? []
    arr.push(api)
    groups.set(key, arr)
  }
  return [...groups.entries()].map(([group, list]) => ({ group, list }))
})

const actionsGrouped = computed(() => {
  const groups = new Map<number, { menuTitle: string; list: ActionOption[] }>()
  for (const a of resources.value.actions) {
    const entry = groups.get(a.menu_id)
    if (entry) {
      entry.list.push(a)
    } else {
      groups.set(a.menu_id, { menuTitle: a.menu_title, list: [a] })
    }
  }
  return [...groups.entries()].map(([menuID, v]) => ({ menuID, ...v }))
})

// ── Loaders ─────────────────────────────────────────────────────────
async function loadData() {
  loading.value = true
  try {
    const result = await listRoles({
      page: filters.page,
      page_size: filters.page_size,
      keyword: filters.keyword?.trim() || undefined,
      role_status: filters.role_status || undefined,
    })
    items.value = result.items ?? []
    total.value = result.total
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : '加载角色列表失败')
  } finally {
    loading.value = false
  }
}

async function loadResources() {
  try {
    resources.value = await getRoleResources()
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : '加载授权资源失败')
  }
}

onMounted(async () => {
  await Promise.all([loadResources(), loadData()])
})

// ── Filters ─────────────────────────────────────────────────────────
function applyFilters() {
  filters.page = 1
  void loadData()
}
function resetFilters() {
  filters.keyword = ''
  filters.role_status = 0
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

// ── Create / Edit ───────────────────────────────────────────────────
function resetFormState() {
  formState.role_code = ''
  formState.role_name = ''
  formState.parent_role_id = 0
  formState.default_route = 'dashboard'
  formState.sort_no = 0
  formState.role_status = STATUS_ACTIVE
  formState.remark = ''
}

function openCreate() {
  resetFormState()
  formMode.value = 'create'
  editingId.value = null
  formDialogVisible.value = true
}

function openEdit(row: RoleItem) {
  resetFormState()
  formMode.value = 'edit'
  editingId.value = row.id
  formState.role_code = row.role_code
  formState.role_name = row.role_name
  formState.parent_role_id = row.parent_role_id
  formState.default_route = row.default_route
  formState.sort_no = row.sort_no
  formState.role_status = row.role_status
  formState.remark = row.remark ?? ''
  formDialogVisible.value = true
}

async function submitForm() {
  if (formMode.value === 'create' && !formState.role_code.trim()) {
    ElMessage.warning('请填写角色编码')
    return
  }
  if (!formState.role_name.trim()) {
    ElMessage.warning('请填写角色名称')
    return
  }
  formSubmitting.value = true
  try {
    if (formMode.value === 'create') {
      const payload: CreateRolePayload = {
        role_code: formState.role_code.trim(),
        role_name: formState.role_name.trim(),
        parent_role_id: formState.parent_role_id || 0,
        default_route: formState.default_route.trim() || 'dashboard',
        sort_no: Number(formState.sort_no) || 0,
        role_status: formState.role_status,
        remark: formState.remark.trim() || null,
      }
      await createRole(payload)
      ElMessage.success('已创建角色')
    } else if (editingId.value) {
      const payload: UpdateRolePayload = {
        role_code: formState.role_code.trim(),
        role_name: formState.role_name.trim(),
        parent_role_id: formState.parent_role_id || 0,
        default_route: formState.default_route.trim() || 'dashboard',
        sort_no: Number(formState.sort_no) || 0,
        role_status: formState.role_status,
        remark: formState.remark.trim() || null,
      }
      await updateRole(editingId.value, payload)
      ElMessage.success('已更新角色')
    }
    formDialogVisible.value = false
    await Promise.all([loadData(), loadResources()])
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : '保存失败')
  } finally {
    formSubmitting.value = false
  }
}

// ── Delete ──────────────────────────────────────────────────────────
async function removeRole(row: RoleItem) {
  if (isBuiltinRole(row)) {
    ElMessage.warning('内置最高权限角色不可删除')
    return
  }
  const warnings: string[] = []
  if (row.user_count > 0) {
    warnings.push(`当前还有 <b>${row.user_count}</b> 个用户绑定此角色，删除后他们将自动失去这部分授权`)
  }
  const childCount = resources.value.roles.filter((r) => r.parent_role_id === row.id).length
  if (childCount > 0) {
    warnings.push(`存在 <b>${childCount}</b> 个子角色，删除后会变成顶级角色（父角色显示为 —）`)
  }
  const body = warnings.length
    ? `<p>确定删除角色「${row.role_name}」吗？该操作为软删除。</p>` +
      warnings.map((w) => `<p style="color:var(--danger,#f56c6c);margin-top:6px">· ${w}</p>`).join('')
    : `确定删除角色「${row.role_name}」吗？该操作为软删除。`
  try {
    await ElMessageBox.confirm(body, '危险操作', {
      type: 'warning',
      dangerouslyUseHTMLString: warnings.length > 0,
      confirmButtonText: '仍然删除',
      confirmButtonClass: 'el-button--danger',
    })
  } catch {
    return
  }
  try {
    await deleteRole(row.id)
    ElMessage.success('已删除角色')
    await Promise.all([loadData(), loadResources()])
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : '删除失败')
  }
}

// ── Permissions dialog ──────────────────────────────────────────────
async function openPermissions(row: RoleItem) {
  permsTarget.value = row
  permsTab.value = 'menus'
  permsDialogVisible.value = true
  permsLoading.value = true
  try {
    const [menus, actions, apis, scopes] = await Promise.all([
      getRoleMenuIDs(row.id),
      getRoleActionIDs(row.id),
      getRoleApiIDs(row.id),
      getRoleDataScopeIDs(row.id),
    ])
    selectedMenuIDs.value = menus.ids ?? []
    selectedActionIDs.value = actions.ids ?? []
    selectedApiIDs.value = apis.ids ?? []
    selectedDataScopeIDs.value = scopes.ids ?? []
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : '加载已分配权限失败')
    permsDialogVisible.value = false
  } finally {
    permsLoading.value = false
  }
}

function getCheckedMenuIDs(): number[] {
  if (!menuTreeRef.value) return selectedMenuIDs.value
  const checked = menuTreeRef.value.getCheckedKeys(false) as number[]
  const halfChecked = menuTreeRef.value.getHalfCheckedKeys() as number[]
  // Persist both fully and half-checked parents so the menu tree renders correctly.
  return [...new Set([...checked, ...halfChecked])]
}

async function submitPermissions() {
  if (!permsTarget.value) return
  permsSubmitting.value = true
  const id = permsTarget.value.id
  const errors: string[] = []
  try {
    const menuIDs = getCheckedMenuIDs()
    await Promise.allSettled([
      assignRoleMenus(id, { ids: menuIDs }).catch((e) => errors.push(`菜单: ${formatErr(e)}`)),
      assignRoleActions(id, { ids: selectedActionIDs.value }).catch((e) => errors.push(`按钮: ${formatErr(e)}`)),
      assignRoleApis(id, { ids: selectedApiIDs.value }).catch((e) => errors.push(`API: ${formatErr(e)}`)),
      assignRoleDataScopes(id, { ids: selectedDataScopeIDs.value }).catch((e) =>
        errors.push(`数据权限: ${formatErr(e)}`),
      ),
    ])
    if (errors.length) {
      ElMessage.error(`部分保存失败：${errors.join('；')}`)
    } else {
      ElMessage.success('已保存权限')
      permsDialogVisible.value = false
    }
  } finally {
    permsSubmitting.value = false
  }
}

function formatErr(e: unknown): string {
  return e instanceof Error ? e.message : String(e)
}

// ── API 分组全选 ─────────────────────────────────────────────────────
function isApiGroupAllSelected(g: (typeof apisGrouped.value)[0]): boolean {
  return g.list.length > 0 && g.list.every((api) => selectedApiIDs.value.includes(api.id))
}
function isApiGroupIndeterminate(g: (typeof apisGrouped.value)[0]): boolean {
  const n = g.list.filter((api) => selectedApiIDs.value.includes(api.id)).length
  return n > 0 && n < g.list.length
}
function toggleApiGroup(g: (typeof apisGrouped.value)[0], checked: boolean) {
  const ids = g.list.map((api) => api.id)
  if (checked) {
    selectedApiIDs.value = [...new Set([...selectedApiIDs.value, ...ids])]
  } else {
    const set = new Set(ids)
    selectedApiIDs.value = selectedApiIDs.value.filter((id) => !set.has(id))
  }
}
function selectAllApis() {
  selectedApiIDs.value = resources.value.apis.map((api) => api.id)
}
function clearAllApis() {
  selectedApiIDs.value = []
}

// ── 按钮权限分组全选 ──────────────────────────────────────────────────
function isActionGroupAllSelected(g: (typeof actionsGrouped.value)[0]): boolean {
  return g.list.length > 0 && g.list.every((a) => selectedActionIDs.value.includes(a.id))
}
function isActionGroupIndeterminate(g: (typeof actionsGrouped.value)[0]): boolean {
  const n = g.list.filter((a) => selectedActionIDs.value.includes(a.id)).length
  return n > 0 && n < g.list.length
}
function toggleActionGroup(g: (typeof actionsGrouped.value)[0], checked: boolean) {
  const ids = g.list.map((a) => a.id)
  if (checked) {
    selectedActionIDs.value = [...new Set([...selectedActionIDs.value, ...ids])]
  } else {
    const set = new Set(ids)
    selectedActionIDs.value = selectedActionIDs.value.filter((id) => !set.has(id))
  }
}
function selectAllActions() {
  selectedActionIDs.value = resources.value.actions.map((a) => a.id)
}
function clearAllActions() {
  selectedActionIDs.value = []
}

// ── Misc ────────────────────────────────────────────────────────────
function parentName(parentRoleID: number): string {
  if (!parentRoleID) return '—'
  const parent = resources.value.roles.find((r) => r.id === parentRoleID)
  return parent?.role_name ?? `#${parentRoleID}`
}

function formatDateTime(value?: string | null) {
  if (!value) return '—'
  const d = new Date(value)
  if (Number.isNaN(d.getTime())) return value
  const pad = (n: number) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}`
}

// el-tree props for menu picker
const menuTreeProps = { label: 'menu_title', children: 'children' }

function methodTagType(method: string) {
  const m = method.toUpperCase()
  if (m === 'GET') return 'success'
  if (m === 'POST') return 'primary'
  if (m === 'PUT') return 'warning'
  if (m === 'DELETE') return 'danger'
  return 'info'
}
</script>

<template>
  <div class="role-page">
    <!-- Filters -->
    <section class="bp-placeholder filter-bar">
      <el-form inline @submit.prevent>
        <el-form-item label="关键词">
          <el-input
            v-model="filters.keyword"
            placeholder="编码 / 名称"
            clearable
            style="width: 220px"
            @keyup.enter="applyFilters"
          />
        </el-form-item>
        <el-form-item label="状态">
          <el-select
            v-model="filters.role_status"
            placeholder="全部"
            clearable
            style="width: 140px"
            @change="applyFilters"
          >
            <el-option :value="0" label="全部" />
            <el-option :value="STATUS_ACTIVE" label="启用" />
            <el-option :value="STATUS_DISABLED" label="禁用" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button :icon="Search" type="primary" @click="applyFilters">查询</el-button>
          <el-button :icon="Refresh" @click="resetFilters">重置</el-button>
        </el-form-item>
        <el-form-item style="margin-left: auto">
          <el-button :icon="Plus" type="primary" @click="openCreate">新增角色</el-button>
        </el-form-item>
      </el-form>
    </section>

    <!-- Table -->
    <section class="bp-placeholder table-card" v-loading="loading">
      <el-table :data="items" stripe row-key="id" empty-text="暂无数据">
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column prop="role_code" label="编码" min-width="140" />
        <el-table-column prop="role_name" label="名称" min-width="140" />
        <el-table-column label="父角色" min-width="140">
          <template #default="{ row }">
            <span>{{ parentName(row.parent_role_id) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="default_route" label="默认路由" min-width="140" />
        <el-table-column label="状态" width="90">
          <template #default="{ row }">
            <el-tag
              :type="roleStatusDict.typeOf(row.role_status)"
              :color="roleStatusDict.colorOf(row.role_status)"
              size="small"
            >
              {{ roleStatusDict.labelOf(row.role_status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="用户数" width="90">
          <template #default="{ row }">
            <el-tag size="small" type="info">{{ row.user_count }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="sort_no" label="排序" width="80" />
        <el-table-column label="创建时间" min-width="160">
          <template #default="{ row }">{{ formatDateTime(row.created_at) }}</template>
        </el-table-column>
        <el-table-column label="操作" fixed="right" width="220">
          <template #default="{ row }">
            <el-button :icon="EditPen" link size="small" @click="openEdit(row)">编辑</el-button>
            <el-button :icon="Key" link size="small" @click="openPermissions(row)">权限</el-button>
            <el-button
              :icon="Delete"
              link
              size="small"
              type="danger"
              :disabled="isBuiltinRole(row)"
              @click="removeRole(row)"
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

    <!-- Create / Edit role dialog -->
    <el-dialog v-model="formDialogVisible" :title="formTitle" width="560px">
      <el-form label-position="top">
        <div class="form-row">
          <el-form-item label="角色编码" required>
            <el-input v-model="formState.role_code" placeholder="例如：operator" />
            <div v-if="formMode === 'edit'" class="form-hint">
              修改编码会影响所有按 code 引用此角色的代码与配置，请谨慎操作。
            </div>
          </el-form-item>
          <el-form-item label="角色名称" required>
            <el-input v-model="formState.role_name" placeholder="例如：运营人员" />
          </el-form-item>
        </div>
        <div class="form-row">
          <el-form-item label="父角色">
            <el-select v-model="formState.parent_role_id" clearable placeholder="无" style="width: 100%">
              <el-option :value="0" label="无（顶级）" />
              <el-option
                v-for="role in parentOptions"
                :key="role.id"
                :value="role.id"
                :label="role.role_name"
              />
            </el-select>
          </el-form-item>
          <el-form-item label="默认路由">
            <el-input v-model="formState.default_route" placeholder="dashboard" />
          </el-form-item>
        </div>
        <div class="form-row">
          <el-form-item label="排序">
            <el-input-number v-model="formState.sort_no" :min="0" :max="9999" style="width: 100%" />
          </el-form-item>
          <el-form-item label="状态">
            <el-select v-model="formState.role_status" style="width: 100%">
              <el-option :value="STATUS_ACTIVE" label="启用" />
              <el-option :value="STATUS_DISABLED" label="禁用" />
            </el-select>
          </el-form-item>
        </div>
        <el-form-item label="备注">
          <el-input v-model="formState.remark" type="textarea" :rows="2" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="formDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="formSubmitting" @click="submitForm">保存</el-button>
      </template>
    </el-dialog>

    <!-- Permissions dialog -->
    <el-dialog
      v-model="permsDialogVisible"
      :title="permsTarget ? `授权：${permsTarget.role_name}` : '授权'"
      width="720px"
      top="6vh"
    >
      <div v-loading="permsLoading">
        <el-tabs v-model="permsTab" class="perms-tabs">
          <el-tab-pane label="菜单权限" name="menus">
            <div class="perms-pane">
              <el-tree
                ref="menuTreeRef"
                :data="resources.menus"
                show-checkbox
                node-key="id"
                :props="menuTreeProps"
                :default-checked-keys="selectedMenuIDs"
                :default-expand-all="true"
                empty-text="暂无菜单"
              />
            </div>
          </el-tab-pane>

          <el-tab-pane label="按钮权限" name="actions">
            <div class="perms-pane">
              <div v-if="actionsGrouped.length" class="perms-toolbar">
                <el-button size="small" @click="selectAllActions">全选</el-button>
                <el-button size="small" @click="clearAllActions">清空</el-button>
              </div>
              <el-empty
                v-if="!actionsGrouped.length"
                description="暂无按钮"
                :image-size="80"
              />
              <div v-for="g in actionsGrouped" :key="g.menuID" class="perms-group">
                <div class="group-header">
                  <el-checkbox
                    :model-value="isActionGroupAllSelected(g)"
                    :indeterminate="isActionGroupIndeterminate(g)"
                    @change="(v: boolean) => toggleActionGroup(g, v)"
                  >
                    <span class="group-title">{{ g.menuTitle }}</span>
                  </el-checkbox>
                </div>
                <el-checkbox-group v-model="selectedActionIDs" class="group-items">
                  <el-checkbox v-for="a in g.list" :key="a.id" :value="a.id">
                    {{ a.action_name }}
                    <span class="action-code">({{ a.action_code }})</span>
                  </el-checkbox>
                </el-checkbox-group>
              </div>
            </div>
          </el-tab-pane>

          <el-tab-pane label="API 权限" name="apis">
            <div class="perms-pane">
              <div v-if="apisGrouped.length" class="perms-toolbar">
                <el-button size="small" @click="selectAllApis">全选</el-button>
                <el-button size="small" @click="clearAllApis">清空</el-button>
              </div>
              <el-empty v-if="!apisGrouped.length" description="暂无 API" :image-size="80" />
              <div v-for="g in apisGrouped" :key="g.group" class="perms-group">
                <div class="group-header">
                  <el-checkbox
                    :model-value="isApiGroupAllSelected(g)"
                    :indeterminate="isApiGroupIndeterminate(g)"
                    @change="(v: boolean) => toggleApiGroup(g, v)"
                  >
                    <span class="group-title">{{ g.group }}</span>
                  </el-checkbox>
                </div>
                <el-checkbox-group v-model="selectedApiIDs" class="group-items">
                  <el-checkbox v-for="api in g.list" :key="api.id" :value="api.id">
                    <el-tag size="small" class="method-tag" :type="methodTagType(api.api_method)">
                      {{ api.api_method }}
                    </el-tag>
                    <span class="api-path">{{ api.api_path }}</span>
                    <span v-if="api.api_desc" class="api-desc">{{ api.api_desc }}</span>
                  </el-checkbox>
                </el-checkbox-group>
              </div>
            </div>
          </el-tab-pane>

          <el-tab-pane label="数据权限" name="dataScopes">
            <div class="perms-pane">
              <p class="hint">勾选当前角色可以查看的角色范围（用于按角色限制可见数据）。</p>
              <el-checkbox-group v-model="selectedDataScopeIDs" class="data-scope-list">
                <el-checkbox
                  v-for="role in resources.roles"
                  :key="role.id"
                  :value="role.id"
                >
                  {{ role.role_name }}
                  <span class="role-code">({{ role.role_code }})</span>
                </el-checkbox>
              </el-checkbox-group>
            </div>
          </el-tab-pane>
        </el-tabs>
      </div>
      <template #footer>
        <el-button @click="permsDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="permsSubmitting" @click="submitPermissions">保存全部</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.role-page {
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

.perms-tabs :deep(.el-tabs__content) {
  padding-top: 8px;
}

.perms-pane {
  max-height: 480px;
  overflow-y: auto;
  padding: 4px 8px;
}

.perms-toolbar {
  display: flex;
  gap: 8px;
  margin-bottom: 12px;
  padding-bottom: 10px;
  border-bottom: 1px solid var(--el-border-color-lighter);
}

.perms-group {
  margin-bottom: 18px;
}

.group-header {
  margin-bottom: 8px;
}

.group-header :deep(.el-checkbox__label) {
  font-size: 13px;
  font-weight: 600;
  color: var(--text-primary);
}

.group-items {
  padding-left: 24px;
}

.perms-group :deep(.el-checkbox-group) {
  display: flex;
  flex-wrap: wrap;
  gap: 12px 18px;
}

.action-code,
.role-code {
  margin-left: 4px;
  color: var(--text-tertiary);
  font-size: 12px;
}

.method-tag {
  margin-right: 6px;
  font-family: var(--font-mono);
}

.api-path {
  font-family: var(--font-mono);
  font-size: 13px;
  color: var(--text-primary);
}

.api-desc {
  margin-left: 8px;
  color: var(--text-tertiary);
  font-size: 12px;
}

.hint {
  margin-bottom: 12px;
  color: var(--text-secondary);
  font-size: 13px;
}

.data-scope-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.form-hint {
  margin-top: 4px;
  color: var(--text-tertiary);
  font-size: 12px;
  line-height: 1.5;
}
</style>
