<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { CircleCheckFilled, CircleCloseFilled, Delete, Refresh, Search } from '@element-plus/icons-vue'
import { cleanupLoginLogs, getLoginLogs } from '@/api/audit'
import type { LoginLogQuery, LoginLogRecord } from '@/types/audit'

const loading = ref(false)
const items = ref<LoginLogRecord[]>([])
const total = ref(0)

const filters = reactive<LoginLogQuery>({
  page: 1,
  page_size: 10,
  keyword: '',
  login_success: 0,
  start_date: '',
  end_date: '',
})

const dateRange = ref<[string, string] | null>(null)

async function load() {
  loading.value = true
  try {
    const res = await getLoginLogs(filters)
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
  filters.login_success = 0
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

async function handleCleanup() {
  try {
    await ElMessageBox.confirm('将删除 90 天前的登录日志，确认继续？', '清理确认', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    })
    const res = await cleanupLoginLogs(90)
    ElMessage.success(`已清理 ${res.deleted} 条记录`)
    load()
  } catch {
    // cancelled
  }
}

onMounted(load)
</script>

<template>
  <div class="login-log-page">
    <!-- Filter bar -->
    <section class="bp-placeholder filter-bar">
      <el-form inline @submit.prevent="handleSearch">
        <el-form-item label="关键词">
          <el-input
            v-model="filters.keyword"
            placeholder="账号 / IP"
            clearable
            style="width: 180px"
          />
        </el-form-item>
        <el-form-item label="登录状态">
          <el-select v-model="filters.login_success" style="width: 120px">
            <el-option label="全部" :value="0" />
            <el-option label="成功" :value="1" />
            <el-option label="失败" :value="2" />
          </el-select>
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
        <el-table-column prop="login_name" label="账号" min-width="120" show-overflow-tooltip />
        <el-table-column prop="source_ip" label="来源 IP" width="140" />
        <el-table-column label="结果" width="90" align="center">
          <template #default="{ row }">
            <el-icon v-if="row.login_success" color="#67c23a" :size="18">
              <CircleCheckFilled />
            </el-icon>
            <el-icon v-else color="#f56c6c" :size="18">
              <CircleCloseFilled />
            </el-icon>
          </template>
        </el-table-column>
        <el-table-column prop="fail_reason" label="失败原因" min-width="160" show-overflow-tooltip />
        <el-table-column prop="user_agent" label="User-Agent" min-width="200" show-overflow-tooltip />
        <el-table-column prop="created_at" label="时间" width="180">
          <template #default="{ row }">
            {{ new Date(row.created_at).toLocaleString('zh-CN') }}
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
  </div>
</template>

<style scoped>
.login-log-page {
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
</style>
