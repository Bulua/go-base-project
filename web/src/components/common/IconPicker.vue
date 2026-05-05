<script setup lang="ts">
import { computed, ref } from 'vue'
import * as Icons from '@element-plus/icons-vue'

const props = defineProps<{ modelValue: string | null }>()
const emit = defineEmits<{ (e: 'update:modelValue', v: string | null): void }>()

const visible = ref(false)
const keyword = ref('')

const iconNames = Object.keys(Icons).sort()

const filtered = computed(() => {
  const kw = keyword.value.trim().toLowerCase()
  return kw ? iconNames.filter((n) => n.toLowerCase().includes(kw)) : iconNames
})

const currentIcon = computed(() =>
  props.modelValue ? Icons[props.modelValue as keyof typeof Icons] ?? null : null,
)

function select(name: string) {
  emit('update:modelValue', name)
  visible.value = false
  keyword.value = ''
}

function clear() {
  emit('update:modelValue', null)
  visible.value = false
}

function openPicker() {
  visible.value = true
}
</script>

<template>
  <el-popover
    v-model:visible="visible"
    trigger="manual"
    placement="bottom-start"
    width="460"
    :teleported="true"
  >
    <template #reference>
      <div class="icon-trigger" @click="openPicker">
        <el-icon v-if="currentIcon" class="trigger-icon">
          <component :is="currentIcon" />
        </el-icon>
        <span class="trigger-label" :class="{ 'is-placeholder': !modelValue }">
          {{ modelValue || '点击选择图标' }}
        </span>
        <el-icon v-if="modelValue" class="trigger-clear" @click.stop="clear">
          <CircleClose />
        </el-icon>
      </div>
    </template>

    <div class="icon-picker-panel">
      <div class="picker-header">
        <el-input
          v-model="keyword"
          placeholder="搜索图标名"
          clearable
          size="small"
          style="flex: 1"
        />
        <el-button size="small" @click="visible = false">关闭</el-button>
      </div>
      <div class="icon-grid">
        <div
          v-for="name in filtered"
          :key="name"
          class="icon-item"
          :class="{ 'is-active': modelValue === name }"
          :title="name"
          @click="select(name)"
        >
          <el-icon :size="18">
            <component :is="Icons[name as keyof typeof Icons]" />
          </el-icon>
          <span class="icon-name">{{ name }}</span>
        </div>
      </div>
      <div v-if="filtered.length === 0" class="icon-empty">无匹配图标</div>
    </div>
  </el-popover>
</template>

<style scoped>
.icon-trigger {
  display: flex;
  align-items: center;
  gap: 6px;
  height: 32px;
  padding: 0 11px;
  border: 1px solid var(--el-border-color);
  border-radius: var(--el-border-radius-base);
  background: var(--el-fill-color-blank);
  cursor: pointer;
  user-select: none;
  transition: border-color 0.2s;
  width: 100%;
  box-sizing: border-box;
}

.icon-trigger:hover {
  border-color: var(--el-border-color-hover);
}

.trigger-icon {
  color: var(--el-text-color-regular);
  flex-shrink: 0;
}

.trigger-label {
  flex: 1;
  font-size: 14px;
  color: var(--el-text-color-regular);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.trigger-label.is-placeholder {
  color: var(--el-text-color-placeholder);
}

.trigger-clear {
  color: var(--el-text-color-placeholder);
  flex-shrink: 0;
  border-radius: 50%;
}

.trigger-clear:hover {
  color: var(--el-text-color-regular);
}

.icon-picker-panel {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.picker-header {
  display: flex;
  gap: 8px;
  align-items: center;
}

.icon-grid {
  display: grid;
  grid-template-columns: repeat(8, 1fr);
  gap: 2px;
  max-height: 300px;
  overflow-y: auto;
}

.icon-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 6px 2px;
  border-radius: 4px;
  cursor: pointer;
  gap: 3px;
  transition: background 0.15s;
}

.icon-item:hover {
  background: var(--el-fill-color-light);
}

.icon-item.is-active {
  background: var(--el-color-primary-light-9);
  color: var(--el-color-primary);
}

.icon-item.is-active .icon-name {
  color: var(--el-color-primary);
}

.icon-name {
  font-size: 9px;
  color: var(--el-text-color-secondary);
  text-align: center;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  width: 100%;
  max-width: 46px;
}

.icon-empty {
  text-align: center;
  color: var(--el-text-color-placeholder);
  font-size: 13px;
  padding: 16px 0;
}
</style>
