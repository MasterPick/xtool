<template>
  <!-- 文档转换工具页面：支持 Markdown/HTML/Plain 互转 -->
  <ToolPage title="文档转换" description="支持 Markdown、HTML、纯文本之间的格式转换">
    <template #actions>
      <!-- 源格式选择 -->
      <select v-model="fromFormat" class="select-field">
        <option value="markdown">Markdown</option>
        <option value="html">HTML</option>
        <option value="plain">纯文本</option>
      </select>
      <span class="opacity-40 text-sm">→</span>
      <!-- 目标格式选择 -->
      <select v-model="toFormat" class="select-field">
        <option value="markdown">Markdown</option>
        <option value="html">HTML</option>
        <option value="plain">纯文本</option>
      </select>
      <button @click="convert" class="btn btn-primary"><RefreshCw :size="14"/>转换</button>
      <button @click="copyOutput" class="btn btn-secondary"><Copy :size="14"/>复制结果</button>
      <button @click="downloadOutput" class="btn btn-secondary"><Download :size="14"/>下载</button>
      <button @click="clearAll" class="btn btn-secondary"><Trash2 :size="14"/>清空</button>
    </template>

    <!-- 双栏编辑区 -->
    <div class="two-col flex-1 min-h-0">
      <!-- 左侧：输入 -->
      <div class="flex flex-col gap-2 min-h-0">
        <div class="label">
          <span>输入（{{ fromFormatLabel }}）</span>
          <span class="text-xs opacity-50">{{ inputStats }}</span>
        </div>
        <textarea
          v-model="inputText"
          class="textarea-field flex-1 min-h-0"
          :placeholder="placeholderText"
          spellcheck="false"
        />
      </div>

      <!-- 右侧：输出 -->
      <div class="flex flex-col gap-2 min-h-0">
        <div class="label">
          <span>
            输出（{{ toFormatLabel }}）
            <!-- HTML 预览切换按钮 -->
            <button
              v-if="toFormat === 'html'"
              @click="showPreview = !showPreview"
              class="ml-2 text-xs px-2 py-0.5 rounded"
              :class="showPreview ? 'btn btn-primary' : 'btn btn-secondary'"
            >
              {{ showPreview ? '源码' : '预览' }}
            </button>
          </span>
          <span
            :class="['badge text-xs', resultStatus === 'success' ? 'badge-success' : resultStatus === 'error' ? 'badge-error' : 'badge-info']"
            v-if="resultStatus"
          >
            {{ resultLabel }}
          </span>
        </div>
        <!-- HTML 预览模式 -->
        <div
          v-if="toFormat === 'html' && showPreview && outputText"
          class="code-output flex-1 min-h-0 overflow-auto"
        >
          <div class="html-preview" v-html="outputText" />
        </div>
        <!-- 普通文本输出 -->
        <div v-else class="code-output flex-1 min-h-0 overflow-auto">
          <span v-if="!outputText" class="opacity-30">转换结果将显示在这里...</span>
          <pre v-else class="whitespace-pre-wrap">{{ outputText }}</pre>
        </div>
      </div>
    </div>
  </ToolPage>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { RefreshCw, Copy, Download, Trash2 } from 'lucide-vue-next'
import ToolPage from '@/components/ToolPage.vue'
import { useAppStore } from '@/stores/app'
import { ConvertDocument } from '../../../wailsjs/go/daily/DailyTools'

const appStore = useAppStore()
const inputText = ref('')
const outputText = ref('')
const fromFormat = ref('markdown')
const toFormat = ref('html')
const showPreview = ref(false)
const resultStatus = ref<'success' | 'error' | 'info' | ''>('')
const resultLabel = ref('')

// 格式标签映射
const formatLabels: Record<string, string> = {
  markdown: 'Markdown',
  html: 'HTML',
  plain: '纯文本',
}

// 源格式标签
const fromFormatLabel = computed(() => formatLabels[fromFormat.value] || fromFormat.value)

// 目标格式标签
const toFormatLabel = computed(() => formatLabels[toFormat.value] || toFormat.value)

// 各格式输入提示
const placeholders: Record<string, string> = {
  markdown: '# 标题\n\n这是一段 **Markdown** 文本。\n\n- 列表项 1\n- 列表项 2\n\n[链接](https://example.com)',
  html: '<h1>标题</h1>\n<p>这是一段 <strong>HTML</strong> 文本。</p>\n<ul>\n  <li>列表项 1</li>\n  <li>列表项 2</li>\n</ul>',
  plain: '这是一段纯文本。\n可以包含多行内容。\n\n空行分隔段落。',
}

// 根据源格式切换占位提示
const placeholderText = computed(() => placeholders[fromFormat.value] || '在此输入文本...')

// 实时统计输入字符数
const inputStats = computed(() => {
  if (!inputText.value) return ''
  return `${inputText.value.length} 字符`
})

// 执行格式转换
async function convert() {
  if (!inputText.value.trim()) {
    appStore.showToast('error', '请输入内容')
    return
  }
  if (fromFormat.value === toFormat.value) {
    appStore.showToast('error', '源格式和目标格式相同')
    return
  }
  try {
    const res = await ConvertDocument(inputText.value, fromFormat.value, toFormat.value) as any
    if (res.success) {
      outputText.value = res.data || ''
      resultStatus.value = 'success'
      resultLabel.value = '转换成功'
      appStore.showToast('success', `${fromFormatLabel.value} → ${toFormatLabel.value} 转换完成`)
    } else {
      outputText.value = res.error || '转换失败'
      resultStatus.value = 'error'
      resultLabel.value = '转换失败'
      appStore.showToast('error', res.error || '转换失败')
    }
  } catch (e) {
    outputText.value = '转换异常: ' + String(e)
    resultStatus.value = 'error'
    resultLabel.value = '转换异常'
    appStore.showToast('error', '转换异常: ' + String(e))
  }
}

// 复制结果到剪贴板
async function copyOutput() {
  if (!outputText.value) return
  try {
    await navigator.clipboard.writeText(outputText.value)
    appStore.showToast('success', '已复制到剪贴板')
  } catch {
    appStore.showToast('error', '复制失败')
  }
}

// 下载转换结果
function downloadOutput() {
  if (!outputText.value) return
  const extMap: Record<string, string> = {
    markdown: 'md',
    html: 'html',
    plain: 'txt',
  }
  const ext = extMap[toFormat.value] || 'txt'
  const blob = new Blob([outputText.value], { type: 'text/plain;charset=utf-8' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `converted.${ext}`
  a.click()
  URL.revokeObjectURL(url)
  appStore.showToast('success', '文件已下载')
}

// 清空输入和输出
function clearAll() {
  inputText.value = ''
  outputText.value = ''
  resultStatus.value = ''
  resultLabel.value = ''
  showPreview.value = false
}
</script>

<style scoped>
/* HTML 预览区域 */
.html-preview {
  padding: 12px;
  border-radius: var(--radius-md, 8px);
  background: var(--bg-primary, #fff);
  color: #1a1a2e;
  line-height: 1.6;
  font-size: 14px;
}

.html-preview :deep(h1) {
  font-size: 24px;
  font-weight: 700;
  margin: 16px 0 8px;
  border-bottom: 1px solid #e5e7eb;
  padding-bottom: 4px;
}

.html-preview :deep(h2) {
  font-size: 20px;
  font-weight: 600;
  margin: 14px 0 6px;
}

.html-preview :deep(h3) {
  font-size: 17px;
  font-weight: 600;
  margin: 12px 0 4px;
}

.html-preview :deep(p) {
  margin: 8px 0;
}

.html-preview :deep(ul),
.html-preview :deep(ol) {
  padding-left: 24px;
  margin: 8px 0;
}

.html-preview :deep(li) {
  margin: 4px 0;
}

.html-preview :deep(code) {
  background: #f1f5f9;
  padding: 2px 6px;
  border-radius: 4px;
  font-size: 13px;
  font-family: 'Courier New', monospace;
}

.html-preview :deep(pre) {
  background: #1e293b;
  color: #e2e8f0;
  padding: 12px;
  border-radius: 8px;
  overflow-x: auto;
  margin: 8px 0;
}

.html-preview :deep(pre code) {
  background: transparent;
  color: inherit;
  padding: 0;
}

.html-preview :deep(a) {
  color: var(--info);
  text-decoration: underline;
}

.html-preview :deep(strong) {
  font-weight: 700;
}

.html-preview :deep(em) {
  font-style: italic;
}

.html-preview :deep(hr) {
  border: none;
  border-top: 1px solid #e5e7eb;
  margin: 16px 0;
}

.html-preview :deep(img) {
  max-width: 100%;
  height: auto;
}
</style>
