<template>
  <ToolPage title="XML 工具" description="XML 格式化与压缩">
<div class="two-col flex-1 min-h-0">
      <div class="flex flex-col gap-2 min-h-0">
        <div class="label">输入 XML</div>
        <textarea v-model="input" class="textarea-field flex-1 min-h-0" placeholder="粘贴 XML 内容..." spellcheck="false"/>
        <div class="flex gap-2">
          <button @click="format" class="btn btn-primary flex-1"><Wand2 :size="14"/>格式化</button>
          <button @click="clearAll" class="btn btn-secondary"><Trash2 :size="14"/>清空</button>
        </div>
      </div>
      <div class="flex flex-col gap-2 min-h-0">
        <div class="label"><span>格式化结果</span>
          <button @click="copy" class="btn btn-secondary py-0.5 px-2 text-xs"><Copy :size="11"/>复制</button>
        </div>
        <div class="code-output flex-1 min-h-0 overflow-auto">
          <span v-if="!output" class="opacity-30">格式化结果显示在这里...</span>
          <span v-else :class="isError ? 'text-red-400' : ''">{{ output }}</span>
        </div>
      </div>
    </div>
  </ToolPage>
</template>
<script setup lang="ts">
import ToolPage from '@/components/ToolPage.vue'
import { ref } from 'vue'
import { Code2, Wand2, Copy, Trash2 } from 'lucide-vue-next'
import { useAppStore } from '@/stores/app'
import { FormatXML } from '../../../wailsjs/go/devtools/DevTools'
const appStore = useAppStore()
const input = ref(''), output = ref(''), isError = ref(false)

async function format() {
  if (!input.value.trim()) return
  isError.value = false
  output.value = ''
  try {
    const res = await FormatXML(input.value)
    isError.value = !res.success
    output.value = res.success ? res.data : (res.error || '格式化失败')
    if (!res.success) {
      appStore.showToast('error', res.error || '格式化失败')
    }
  } catch (e) {
    isError.value = true
    output.value = '格式化异常: ' + String(e)
    appStore.showToast('error', '格式化异常: ' + String(e))
  }
}

async function copy() {
  if (!output.value) return
  try {
    await navigator.clipboard.writeText(output.value)
    appStore.showToast('success', '已复制')
  } catch {
    appStore.showToast('error', '复制失败')
  }
}

function clearAll() {
  input.value = ''
  output.value = ''
  isError.value = false
}
</script>
