<template>
  <div class="page-container">
    <div class="page-title"><Link :size="20" class="text-primary-400"/>URL 编解码</div>
    <div class="page-desc">URL 编码与解码</div>
    <div class="tab-bar mb-4">
      <button :class="['tab-item', mode==='encode'&&'active']" @click="mode='encode'">编码</button>
      <button :class="['tab-item', mode==='decode'&&'active']" @click="mode='decode'">解码</button>
    </div>
    <div class="two-col flex-1 min-h-0">
      <div class="flex flex-col gap-2 min-h-0">
        <div class="label">{{ mode==='encode' ? '原始 URL' : '编码后的 URL' }}</div>
        <textarea v-model="input" class="textarea-field flex-1 min-h-0" spellcheck="false"
          :placeholder="mode==='encode' ? '输入原始 URL...' : '输入已编码的 URL...'"/>
        <button @click="process" class="btn btn-primary w-full" :disabled="processing">
          <Wand2 :size="14"/>{{ processing ? '处理中...' : (mode==='encode'?'编码':'解码') }}
        </button>
      </div>
      <div class="flex flex-col gap-2 min-h-0">
        <div class="label"><span>结果</span>
          <button @click="copy" class="btn btn-secondary py-0.5 px-2 text-xs"><Copy :size="11"/>复制</button>
        </div>
        <div class="code-output flex-1 overflow-auto" :class="hasError ? 'border-red-500/30' : ''">
          <span v-if="!output" class="opacity-30">结果显示在这里...</span>
          <span v-else :class="hasError ? 'text-red-400' : ''">{{ output }}</span>
        </div>
      </div>
    </div>
  </div>
</template>
<script setup lang="ts">
import { ref } from 'vue'
import { Link, Wand2, Copy } from 'lucide-vue-next'
import { useAppStore } from '@/stores/app'
import { URLEncode, URLDecode } from '../../../wailsjs/go/devtools/DevTools'
const appStore = useAppStore()
const mode = ref<'encode'|'decode'>('encode')
const input = ref(''), output = ref('')
const hasError = ref(false)
const processing = ref(false)

async function process() {
  if (!input.value) return
  processing.value = true
  hasError.value = false
  output.value = ''
  try {
    const res = mode.value === 'encode' ? await URLEncode(input.value) : await URLDecode(input.value)
    if (res.success) {
      output.value = res.data
      hasError.value = false
    } else {
      output.value = res.error || '处理失败'
      hasError.value = true
      appStore.showToast('error', res.error || '处理失败')
    }
  } catch (e) {
    output.value = '处理异常: ' + String(e)
    hasError.value = true
    appStore.showToast('error', '处理异常: ' + String(e))
  } finally {
    processing.value = false
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
</script>
