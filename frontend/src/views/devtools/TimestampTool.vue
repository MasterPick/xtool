<template>
  <ToolPage title="时间戳转换" description="Unix 时间戳与日期互转">

    <!-- 当前时间 -->
    <div class="card mb-4">
      <div class="flex items-center justify-between">
        <span class="text-sm opacity-60">当前时间</span>
        <button @click="refreshNow" class="btn btn-secondary py-1 px-2 text-xs"><RefreshCw :size="12"/>刷新</button>
      </div>
      <div class="code-output mt-2 text-sm">{{ nowInfo }}</div>
    </div>

    <div class="two-col">
      <!-- 时间戳 → 日期 -->
      <div class="card">
        <div class="label mb-3">时间戳 → 日期时间</div>
        <input v-model="tsInput" class="input-field mb-2"
          placeholder="输入时间戳（秒/毫秒）..." />
        <div v-if="tsInput" class="text-xs mb-3" :class="tsInput.length > 13 ? 'text-blue-400' : 'text-green-400'">
          {{ tsInput.length > 13 ? '自动检测: 毫秒级时间戳' : tsInput.length === 10 ? '自动检测: 秒级时间戳' : '自动检测: 请输入10位(秒)或13位(毫秒)时间戳' }}
        </div>
        <button @click="tsToDate" class="btn btn-primary w-full mb-3">转换</button>
        <div class="code-output text-sm" :class="tsError ? 'text-red-400' : ''">{{ tsResult || '结果显示在这里...' }}</div>
      </div>

      <!-- 日期 → 时间戳 -->
      <div class="card">
        <div class="label mb-3">日期时间 → 时间戳</div>
        <input v-model="dateInput" class="input-field mb-3"
          placeholder="如：2024-01-01 12:00:00" />
        <button @click="dateToTs" class="btn btn-primary w-full mb-3">转换</button>
        <div class="code-output text-sm" :class="dateError ? 'text-red-400' : ''">{{ dateResult || '结果显示在这里...' }}</div>
      </div>
    </div>
  </ToolPage>
</template>

<script setup lang="ts">
import ToolPage from '@/components/ToolPage.vue'
import { ref, onMounted } from 'vue'
import { Clock, RefreshCw } from 'lucide-vue-next'
import { useAppStore } from '@/stores/app'
import { TimestampToDatetime, DatetimeToTimestamp, GetCurrentTimestamp } from '../../../wailsjs/go/devtools/DevTools'

const appStore = useAppStore()
const tsInput   = ref('')
const dateInput = ref('')
const tsResult  = ref('')
const dateResult = ref('')
const nowInfo   = ref('')
const tsError   = ref(false)
const dateError = ref(false)

async function refreshNow() {
  try {
    const res = await GetCurrentTimestamp()
    if (res.success !== false) {
      nowInfo.value = res.data
    } else {
      nowInfo.value = '获取失败: ' + (res.error || '')
    }
  } catch (e) {
    nowInfo.value = '获取失败: ' + String(e)
  }
}

async function tsToDate() {
  if (!tsInput.value) return
  tsError.value = false
  const ts = parseInt(tsInput.value)
  if (isNaN(ts)) {
    tsResult.value = '请输入有效的时间戳'
    tsError.value = true
    appStore.showToast('error', '请输入有效的时间戳')
    return
  }
  try {
    const res = await TimestampToDatetime(ts)
    if (res.success !== false) {
      tsResult.value = res.data
      tsError.value = false
    } else {
      tsResult.value = res.error || '转换失败'
      tsError.value = true
    }
  } catch (e) {
    tsResult.value = '转换异常: ' + String(e)
    tsError.value = true
    appStore.showToast('error', '转换异常: ' + String(e))
  }
}

async function dateToTs() {
  if (!dateInput.value) return
  dateError.value = false
  try {
    const res = await DatetimeToTimestamp(dateInput.value)
    if (res.success) {
      dateResult.value = res.data
      dateError.value = false
    } else {
      dateResult.value = res.error || '转换失败'
      dateError.value = true
      appStore.showToast('error', res.error || '转换失败')
    }
  } catch (e) {
    dateResult.value = '转换异常: ' + String(e)
    dateError.value = true
    appStore.showToast('error', '转换异常: ' + String(e))
  }
}

onMounted(refreshNow)
</script>
