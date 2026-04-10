<template>
  <!-- 代码格式化工具页面：支持 JSON/XML/HTML/CSS/SQL/YAML 格式化 -->
  <ToolPage title="代码格式化" description="支持 JSON、XML、HTML、CSS、SQL、YAML 等多种语言的代码格式化">
    <template #actions>
      <!-- 语言选择下拉框 -->
      <select v-model="language" class="select-field">
        <option value="json">JSON</option>
        <option value="xml">XML</option>
        <option value="html">HTML</option>
        <option value="css">CSS</option>
        <option value="sql">SQL</option>
        <option value="yaml">YAML</option>
      </select>
      <button @click="formatCode" class="btn btn-primary"><Wand2 :size="14"/>格式化</button>
      <button @click="copyOutput" class="btn btn-secondary"><Copy :size="14"/>复制结果</button>
      <button @click="downloadOutput" class="btn btn-secondary"><Download :size="14"/>下载</button>
      <button @click="clearAll" class="btn btn-secondary"><Trash2 :size="14"/>清空</button>
    </template>

    <!-- 双栏编辑区 -->
    <div class="two-col flex-1 min-h-0">
      <!-- 左侧：输入 -->
      <div class="flex flex-col gap-2 min-h-0">
        <div class="label">
          <span>输入代码</span>
          <span class="text-xs opacity-50">{{ inputStats }}</span>
        </div>
        <textarea
          v-model="inputCode"
          class="textarea-field flex-1 min-h-0"
          :placeholder="placeholderText"
          spellcheck="false"
        />
      </div>

      <!-- 右侧：输出 -->
      <div class="flex flex-col gap-2 min-h-0">
        <div class="label">
          <span>格式化结果</span>
          <span
            :class="['badge text-xs', resultStatus === 'success' ? 'badge-success' : resultStatus === 'error' ? 'badge-error' : 'badge-info']"
            v-if="resultStatus"
          >
            {{ resultLabel }}
          </span>
        </div>
        <div class="code-output flex-1 min-h-0 overflow-auto">
          <span v-if="!outputCode" class="opacity-30">格式化结果将显示在这里...</span>
          <pre v-else class="whitespace-pre-wrap">{{ outputCode }}</pre>
        </div>
      </div>
    </div>
  </ToolPage>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { Wand2, Copy, Download, Trash2 } from 'lucide-vue-next'
import ToolPage from '@/components/ToolPage.vue'
import { useAppStore } from '@/stores/app'
import { FormatCode } from '../../../wailsjs/go/devtools/DevTools'

const appStore = useAppStore()
const inputCode = ref('')
const outputCode = ref('')
const language = ref('json')
const resultStatus = ref<'success' | 'error' | 'info' | ''>('')
const resultLabel = ref('')

// 各语言输入提示
const placeholders: Record<string, string> = {
  json: '在此粘贴 JSON 内容，例如：\n{"name":"张三","age":25}',
  xml: '在此粘贴 XML 内容，例如：\n<root><item>示例</item></root>',
  html: '在此粘贴 HTML 内容，例如：\n<div><h1>标题</h1><p>段落</p></div>',
  css: '在此粘贴 CSS 内容，例如：\nbody{margin:0;padding:0}',
  sql: '在此粘贴 SQL 内容，例如：\nSELECT id, name FROM users WHERE age > 18',
  yaml: '在此粘贴 YAML 内容，例如：\nname: 张三\nage: 25',
}

// 根据语言切换占位提示
const placeholderText = computed(() => placeholders[language.value] || '在此粘贴代码...')

// 实时统计输入字符数
const inputStats = computed(() => {
  if (!inputCode.value) return ''
  return `${inputCode.value.length} 字符`
})

// 执行格式化
async function formatCode() {
  if (!inputCode.value.trim()) {
    appStore.showToast('error', '请输入代码内容')
    return
  }
  try {
    const res = await FormatCode(inputCode.value, language.value)
    if (res.success) {
      outputCode.value = res.data
      resultStatus.value = 'success'
      resultLabel.value = '格式化成功'
      appStore.showToast('success', `${language.value.toUpperCase()} 格式化完成`)
    } else {
      outputCode.value = res.error
      resultStatus.value = 'error'
      resultLabel.value = '格式化失败'
      appStore.showToast('error', res.error || '格式化失败')
    }
  } catch (e) {
    outputCode.value = '格式化异常: ' + String(e)
    resultStatus.value = 'error'
    resultLabel.value = '格式化异常'
    appStore.showToast('error', '格式化异常: ' + String(e))
  }
}

// 复制结果到剪贴板
async function copyOutput() {
  if (!outputCode.value) return
  try {
    await navigator.clipboard.writeText(outputCode.value)
    appStore.showToast('success', '已复制到剪贴板')
  } catch {
    appStore.showToast('error', '复制失败')
  }
}

// 下载格式化结果
function downloadOutput() {
  if (!outputCode.value) return
  // 根据语言确定文件扩展名
  const extMap: Record<string, string> = {
    json: 'json', xml: 'xml', html: 'html',
    css: 'css', sql: 'sql', yaml: 'yaml',
  }
  const ext = extMap[language.value] || 'txt'
  const blob = new Blob([outputCode.value], { type: 'text/plain;charset=utf-8' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `formatted.${ext}`
  a.click()
  URL.revokeObjectURL(url)
  appStore.showToast('success', '文件已下载')
}

// 清空输入和输出
function clearAll() {
  inputCode.value = ''
  outputCode.value = ''
  resultStatus.value = ''
  resultLabel.value = ''
}
</script>
