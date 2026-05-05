<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Delete, Refresh, Search, View } from '@element-plus/icons-vue'
import { cleanupOperationLogs, getOperationLog, getOperationLogs } from '@/api/audit'
import type { OperationLogQuery, OperationLogRecord } from '@/types/audit'

const loading = ref(false)
const items = ref<OperationLogRecord[]>([])
const total = ref(0)

const filters = reactive<OperationLogQuery>({
  page: 1,
  page_size: 10,
  keyword: '',
  method: '',
  status_code: 0,
  start_date: '',
  end_date: '',
})

const dateRange = ref<[string, string] | null>(null)

// detail drawer
const detailVisible = ref(false)
const detailLoading = ref(false)
const detailRecord = ref<OperationLogRecord | null>(null)

async function load() {
  loading.value = true
  try {
    const res = await getOperationLogs(filters)
    items.value = res.items
    total.value = res.total
  } catch {
    ElMessage.error('加载失败')
  } finally {
    loading.value = false
  }
}

function handleSearch() {
  if (dateRange.value) {
    filters.start_date = dateRange.value[0]
    filters.end_date = dateRange.value[1]
  } else {
    filters.start_date = ''
    filters.end_date = ''
  }
  filters.page = 1
  load()
}

function handleReset() {
  filters.keyword = ''
  filters.method = ''
  filters.status_code = 0
  filters.start_date = ''
  filters.end_date = ''
  dateRange.value = null
  filters.page = 1
  load()
}

function handlePageChange(page: number) {
  filters.page = page
  load()
}

function handleSizeChange(size: number) {
  filters.page_size = size
  filters.page = 1
  load()
}

async function handleView(row: OperationLogRecord) {
  detailVisible.value = true
  detailLoading.value = true
  detailRecord.value = null
  try {
    const rec = await getOperationLog(row.id)
    detailRecord.value = rec
  } catch {
    ElMessage.error('加载详情失败')
  } finally {
    detailLoading.value = false
  }
}

async function handleCleanup() {
  try {
    await ElMessageBox.confirm('将删除 90 天前的操作日志，确认继续？', '清理确认', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    })
    const res = await cleanupOperationLogs(90)
    ElMessage.success(`已清理 ${res.deleted} 条记录`)
    load()
  } catch {
    // cancelled
  }
}

function methodTag(method: string | null) {
  const map: Record<string, string> = {
    GET: 'success',
    POST: 'primary',
    PUT: 'warning',
    DELETE: 'danger',
    PATCH: 'warning',
  }
  return map[method ?? ''] ?? 'info'
}

function statusTag(code: number | null) {
  if (!code) return 'info'
  if (code < 300) return 'success'
  if (code < 400) return 'warning'
  return 'danger'
}

function formatJson(text: string | null) {
  if (!text) return ''
  try {
    return JSON.stringify(JSON.parse(text), null, 2)
  } catch {
    return text
  }
}

onMounted(load)
</script>

<template>
  <div class="operation-log-page">
    <!-- Filter bar -->
    <section class="bp-placeholder filter-bar">
      <el-form inline @submit.prevent="handleSearch">
        <el-form-item label="关键词">
          <el-input
            v-model="filters.keyword"
            placeholder="请求路径 / IP"
            clearable
            style="width: 200px"
          />
        </el-form-item>
        <el-form-item label="请求方法">
          <el-select v-model="filters.method" clearable style="width: 110px">
            <el-option label="全部" value="" />
            <el-option label="GET" value="GET" />
            <el-option label="POST" value="POST" />
            <el-option label="PUT" value="PUT" />
            <el-option label="DELETE" value="DELETE" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态码">
          <el-input
            v-model.number="filters.status_code"
            placeholder="如 200"
            clearable
            style="width: 100px"
          />
        </el-form-item>
        <el-form-item label="时间范围">
          <el-date-picker
            v-model="dateRange"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            value-format="YYYY-MM-DD"
            style="width: 240px"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :icon="Search" @click="handleSearch">查询</el-button>
          <el-button :icon="Refresh" @click="handleReset">重置</el-button>
          <el-button type="danger" plain :icon="Delete" @click="handleCleanup">清理旧日志</el-button>
        </el-form-item>
      </el-form>
    </section>

    <!-- Table card -->
    <section class="bp-placeholder table-card" v-loading="loading">
      <el-table :data="items" stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column label="方法" width="90" align="center">
          <template #default="{ row }">
            <el-tag :type="methodTag(row.request_method)" size="small">
              {{ row.request_method }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="request_path" label="路径" min-width="120" show-overflow-tooltip />
        <el-table-column label="状态码" width="90" align="center">
          <template #default="{ row }">
            <el-tag :type="statusTag(row.status_code)" size="small">
              {{ row.status_code }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="cost_ms" label="耗时(ms)" width="100" align="right" />
        <el-table-column prop="source_ip" label="来源 IP" width="140" />
        <el-table-column prop="created_at" label="时间" width="180">
          <template #default="{ row }">
            {{ new Date(row.created_at).toLocaleString('zh-CN') }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="80" align="center" fixed="right">
          <template #default="{ row }">
            <el-button text :icon="View" @click="handleView(row)" />
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-row">
        <el-pagination
          v-model:current-page="filters.page"
          v-model:page-size="filters.page_size"
          :total="total"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next"
          @current-change="handlePageChange"
          @size-change="handleSizeChange"
        />
      </div>
    </section>

    <!-- Detail drawer -->
    <el-drawer
      v-model="detailVisible"
      title="操作日志详情"
      size="600px"
      direction="rtl"
    >
      <div v-loading="detailLoading">
        <template v-if="detailRecord">
          <el-descriptions :column="2" border>
            <el-descriptions-item label="ID">{{ detailRecord.id }}</el-descriptions-item>
            <el-descriptions-item label="用户 ID">{{ detailRecord.user_id ?? '-' }}</el-descriptions-item>
            <el-descriptions-item label="来源 IP">{{ detailRecord.source_ip ?? '-' }}</el-descriptions-item>
            <el-descriptions-item label="方法">
              <el-tag :type="methodTag(detailRecord.request_method)" size="small">
                {{ detailRecord.request_method }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="路径" :span="2">{{ detailRecord.request_path }}</el-descriptions-item>
            <el-descriptions-item label="状态码">
              <el-tag :type="statusTag(detailRecord.status_code)" size="small">
                {{ detailRecord.status_code }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="耗时">{{ detailRecord.cost_ms }} ms</el-descriptions-item>
            <el-descriptions-item label="时间" :span="2">
              {{ new Date(detailRecord.created_at).toLocaleString('zh-CN') }}
            </el-descriptions-item>
            <el-descriptions-item v-if="detailRecord.error_message" label="错误信息" :span="2">
              {{ detailRecord.error_message }}
            </el-descriptions-item>
          </el-descriptions>

          <el-divider>请求体</el-divider>
          <pre class="body-pre">{{ formatJson(detailRecord.request_body) || '-' }}</pre>

          <el-divider>响应体</el-divider>
          <pre class="body-pre">{{ formatJson(detailRecord.response_body) || '-' }}</pre>
        </template>
      </div>
    </el-drawer>
  </div>
</template>

<style scoped>
.operation-log-page {
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

.body-pre {
  background: #f5f7fa;
  border-radius: 4px;
  padding: 12px;
  font-size: 12px;
  line-height: 1.6;
  overflow-x: auto;
  white-space: pre-wrap;
  word-break: break-all;
  color: #606266;
  max-height: 300px;
  overflow-y: auto;
}
</style>
