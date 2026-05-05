<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { UploadRequestOptions } from 'element-plus'
import { Delete, Download, Refresh, Search, Upload, ZoomIn } from '@element-plus/icons-vue'
import { deleteFile, downloadFileBlob, listFiles, uploadFile } from '@/api/file'
import type { FileListQuery, FileRecord } from '@/types/file'

// ── 列表状态 ───────────────────────────────────────────────────────────────

const loading = ref(false)
const items = ref<FileRecord[]>([])
const total = ref(0)

const filters = reactive<FileListQuery>({
  page: 1,
  page_size: 10,
  keyword: '',
  mime_category: '',
  start_date: '',
  end_date: '',
})

const dateRange = ref<[string, string] | null>(null)

async function load() {
  loading.value = true
  try {
    const res = await listFiles(filters)
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
  filters.mime_category = ''
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

// ── 上传 ───────────────────────────────────────────────────────────────────

const uploading = ref(false)

async function handleUpload(options: UploadRequestOptions) {
  uploading.value = true
  try {
    await uploadFile(options.file as File, (pct) => {
      options.onProgress({ percent: pct } as ProgressEvent & { percent: number })
    })
    ElMessage.success(`${(options.file as File).name} 上传成功`)
    load()
  } catch (e: unknown) {
    const msg = e instanceof Error ? e.message : '上传失败'
    ElMessage.error(msg)
    options.onError(Object.assign(new Error(msg), { status: 0, method: 'POST', url: '/api/v1/files' }))
  } finally {
    uploading.value = false
  }
}

// ── 删除 ───────────────────────────────────────────────────────────────────

async function handleDelete(row: FileRecord) {
  try {
    await ElMessageBox.confirm(`确认删除文件「${row.original_name}」？`, '删除确认', {
      confirmButtonText: '删除',
      cancelButtonText: '取消',
      type: 'warning',
      confirmButtonClass: 'el-button--danger',
    })
    await deleteFile(row.id)
    ElMessage.success('已删除')
    if (items.value.length === 1 && filters.page > 1) filters.page--
    load()
  } catch {
    // cancelled
  }
}

// ── 预览 / 下载 ────────────────────────────────────────────────────────────

const previewVisible = ref(false)
const previewSrc = ref('')
const previewLoading = ref(false)
let lastObjectUrl = ''

async function handlePreview(row: FileRecord) {
  previewLoading.value = true
  try {
    const blob = await downloadFileBlob(row.id)
    if (lastObjectUrl) URL.revokeObjectURL(lastObjectUrl)
    lastObjectUrl = URL.createObjectURL(blob)
    previewSrc.value = lastObjectUrl
    previewVisible.value = true
  } catch {
    ElMessage.error('预览失败')
  } finally {
    previewLoading.value = false
  }
}

async function handleDownload(row: FileRecord) {
  try {
    const blob = await downloadFileBlob(row.id)
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = row.original_name
    a.click()
    URL.revokeObjectURL(url)
  } catch {
    ElMessage.error('下载失败')
  }
}

// ── 工具函数 ───────────────────────────────────────────────────────────────

function isImage(mimeType: string): boolean {
  return mimeType.startsWith('image/')
}

function formatSize(bytes: number): string {
  if (bytes < 1024) return `${bytes} B`
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`
  if (bytes < 1024 * 1024 * 1024) return `${(bytes / 1024 / 1024).toFixed(1)} MB`
  return `${(bytes / 1024 / 1024 / 1024).toFixed(2)} GB`
}

function mimeLabel(mimeType: string): string {
  if (!mimeType) return 'unknown'
  const parts = mimeType.split('/')
  return parts[parts.length - 1].split(';')[0].toUpperCase()
}

function mimeTagType(mimeType: string): '' | 'success' | 'warning' | 'danger' | 'info' {
  if (mimeType.startsWith('image/')) return 'success'
  if (mimeType.startsWith('video/')) return 'warning'
  if (mimeType.startsWith('text/')) return 'info'
  if (mimeType.includes('pdf')) return 'danger'
  return ''
}

function formatDate(val: string): string {
  if (!val) return '-'
  const d = new Date(val)
  return isNaN(d.getTime()) ? val : d.toLocaleString('zh-CN', {
    year: 'numeric', month: '2-digit', day: '2-digit',
    hour: '2-digit', minute: '2-digit', second: '2-digit', hour12: false,
  })
}

onMounted(load)
</script>

<template>
  <div class="file-page">
    <!-- Filter bar -->
    <section class="bp-placeholder filter-bar">
      <el-form inline @submit.prevent="handleSearch">
        <el-form-item label="文件名">
          <el-input
            v-model="filters.keyword"
            placeholder="搜索文件名"
            clearable
            style="width: 180px"
          />
        </el-form-item>
        <el-form-item label="文件类型">
          <el-select v-model="filters.mime_category" clearable style="width: 110px">
            <el-option label="全部" value="" />
            <el-option label="图片" value="image" />
            <el-option label="视频" value="video" />
            <el-option label="音频" value="audio" />
            <el-option label="文本" value="text" />
            <el-option label="PDF" value="pdf" />
            <el-option label="其他" value="other" />
          </el-select>
        </el-form-item>
        <el-form-item label="上传时间">
          <el-date-picker
            v-model="dateRange"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            value-format="YYYY-MM-DD"
            style="width: 220px"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :icon="Search" @click="handleSearch">查询</el-button>
          <el-button :icon="Refresh" @click="handleReset">重置</el-button>
        </el-form-item>
        <el-form-item style="margin-left: auto">
          <el-upload
            multiple
            :show-file-list="false"
            :http-request="handleUpload"
            :disabled="uploading"
          >
            <el-button type="primary" :icon="Upload" :loading="uploading">上传文件</el-button>
          </el-upload>
        </el-form-item>
      </el-form>
    </section>

    <!-- File table -->
    <section class="bp-placeholder table-card" v-loading="loading">
      <el-table :data="items" stripe>
        <el-table-column label="文件名" min-width="200" show-overflow-tooltip>
          <template #default="{ row }">
            <div class="file-name-cell">
              <el-icon class="file-icon" :class="isImage(row.mime_type) ? 'is-image' : 'is-file'">
                <component :is="isImage(row.mime_type) ? 'Picture' : 'Document'" />
              </el-icon>
              <span>{{ row.original_name }}</span>
            </div>
          </template>
        </el-table-column>

        <el-table-column label="类型" width="110" align="center">
          <template #default="{ row }">
            <el-tag :type="mimeTagType(row.mime_type)" size="small">
              {{ mimeLabel(row.mime_type) }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column label="大小" width="100" align="right">
          <template #default="{ row }">{{ formatSize(row.file_size) }}</template>
        </el-table-column>

        <el-table-column label="上传时间" width="180">
          <template #default="{ row }">{{ formatDate(row.created_at) }}</template>
        </el-table-column>

        <el-table-column label="操作" width="130" align="center" fixed="right">
          <template #default="{ row }">
            <el-button
              v-if="isImage(row.mime_type)"
              text
              :icon="ZoomIn"
              :loading="previewLoading"
              title="预览"
              @click="handlePreview(row)"
            />
            <el-button
              text
              :icon="Download"
              title="下载"
              @click="handleDownload(row)"
            />
            <el-button
              text
              type="danger"
              :icon="Delete"
              title="删除"
              @click="handleDelete(row)"
            />
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

    <!-- Image preview dialog -->
    <el-dialog
      v-model="previewVisible"
      title="图片预览"
      width="80%"
      top="5vh"
      :append-to-body="true"
      destroy-on-close
    >
      <div class="preview-img-wrap">
        <img :src="previewSrc" class="preview-img" alt="preview" />
      </div>
    </el-dialog>
  </div>
</template>

<style scoped>
.file-page {
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

.file-name-cell {
  display: flex;
  align-items: center;
  gap: 8px;
}

.file-icon {
  font-size: 18px;
  flex-shrink: 0;
}
.file-icon.is-image { color: var(--el-color-success); }
.file-icon.is-file  { color: var(--el-text-color-secondary); }

.pagination-row {
  margin-top: 14px;
  display: flex;
  justify-content: flex-end;
}

.preview-img-wrap {
  display: flex;
  justify-content: center;
  max-height: 80vh;
  overflow: auto;
}

.preview-img {
  max-width: 100%;
  object-fit: contain;
}
</style>
