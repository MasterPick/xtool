<template>
  <div class="page-container">
    <div>
      <div class="page-title"><Fingerprint :size="20" class="text-primary-400"/>UUID 生成器</div>
      <div class="page-desc">批量生成标准 UUID v4</div>
    </div>
    <div class="card mb-4 flex items-center gap-4 flex-wrap">
      <div class="label mb-0">生成数量</div>
      <input v-model.number="count" type="number" min="1" max="100"
        class="input-field w-24" />
      <div class="label mb-0">格式</div>
      <select v-model="format" class="input-field w-36">
        <option value="standard">标准 (含连字符)</option>
        <option value="uppercase">大写 (含连字符)</option>
        <option value="lowercase">小写 (含连字符)</option>
        <option value="no-hyphens">无连字符</option>
        <option value="brace">大括号包裹</option>
      </select>
      <button @click="generate" class="btn btn-primary" :disabled="generating">
        <RefreshCw :size="14"/>{{ generating ? '生成中...' : '生成' }}
      </button>
      <button @click="copyAll"  class="btn btn-secondary"><Copy :size="14"/>复制全部</button>
      <button @click="clear"    class="btn btn-secondary"><Trash2 :size="14"/>清空</button>
    </div>
    <div class="code-output flex-1 min-h-0 overflow-auto">
      <span v-if="!uuids" class="opacity-30">生成的 UUID 将显示在这里...</span>
      <div v-else>
        <div v-for="(uuid, idx) in uuidList" :key="idx" class="flex items-center gap-2 py-0.5 hover:bg-white/5 px-1">
          <span class="font-mono text-sm flex-1">{{ uuid }}</span>
          <button @click="copySingle(uuid)" class="btn btn-secondary py-0.5 px-2 text-xs opacity-0 hover:opacity-100" style="opacity:0.5">
            <Copy :size="11"/>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { Fingerprint, RefreshCw, Copy, Trash2 } from 'lucide-vue-next'
import { useAppStore } from '@/stores/app'
import { GenerateUUIDs } from '../../../wailsjs/go/devtools/DevTools'

const appStore = useAppStore()
const count = ref(10)
const uuids = ref('')
const format = ref<'standard' | 'uppercase' | 'lowercase' | 'no-hyphens' | 'brace'>('standard')
const generating = ref(false)

const uuidList = computed(() => {
  if (!uuids.value) return []
  return uuids.value.split('\n').map(u => formatUuid(u.trim())).filter(Boolean)
})

function formatUuid(uuid: string): string {
  switch (format.value) {
    case 'uppercase': return uuid.toUpperCase()
    case 'lowercase': return uuid.toLowerCase()
    case 'no-hyphens': return uuid.replace(/-/g, '')
    case 'brace': return `{${uuid}}`
    default: return uuid
  }
}

async function generate() {
  generating.value = true
  uuids.value = ''
  try {
    const res = await GenerateUUIDs(count.value)
    if (res.success !== false) {
      uuids.value = res.data
      appStore.showToast('success', `已生成 ${count.value} 个 UUID`)
    } else {
      appStore.showToast('error', res.error || '生成失败')
    }
  } catch (e) {
    appStore.showToast('error', '生成 UUID 失败: ' + String(e))
  } finally {
    generating.value = false
  }
}

async function copyAll() {
  if (!uuids.value) return
  try {
    const text = uuidList.value.join('\n')
    await navigator.clipboard.writeText(text)
    appStore.showToast('success', '已复制到剪贴板')
  } catch {
    appStore.showToast('error', '复制失败')
  }
}

async function copySingle(uuid: string) {
  try {
    await navigator.clipboard.writeText(uuid)
    appStore.showToast('success', '已复制')
  } catch {
    appStore.showToast('error', '复制失败')
  }
}

function clear() { uuids.value = '' }
</script>
