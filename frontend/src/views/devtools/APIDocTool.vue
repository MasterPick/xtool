<template>
  <ToolPage title="接口文档生成" description="Markdown 格式接口文档">
<!-- 接口文档生成工具 -->
  
    

    <!-- 基本信息表单 -->
    <div class="card mb-4">
      <div class="label mb-3">基本信息</div>
      <div class="grid grid-cols-2 gap-3">
        <div>
          <div class="text-xs opacity-50 mb-1">接口名称</div>
          <input v-model="apiName" class="input-field" placeholder="如：用户登录接口" />
        </div>
        <div class="flex gap-2">
          <div class="w-32">
            <div class="text-xs opacity-50 mb-1">请求方法</div>
            <select v-model="method" class="input-field">
              <option value="GET">GET</option>
              <option value="POST">POST</option>
              <option value="PUT">PUT</option>
              <option value="DELETE">DELETE</option>
            </select>
          </div>
          <div class="flex-1">
            <div class="text-xs opacity-50 mb-1">URL</div>
            <input v-model="apiUrl" class="input-field font-mono" placeholder="/api/v1/users" />
          </div>
        </div>
      </div>
      <div class="mt-3">
        <div class="text-xs opacity-50 mb-1">接口描述</div>
        <textarea v-model="apiDesc" class="textarea-field" rows="2" placeholder="描述该接口的功能..." />
      </div>
    </div>

    <!-- 请求头 -->
    <div class="card mb-4">
      <div class="flex items-center justify-between mb-3">
        <div class="label mb-0">请求头 (Headers)</div>
        <button @click="addHeader" class="btn btn-secondary py-0.5 px-2 text-xs">
          <Plus :size="12"/>添加
        </button>
      </div>
      <div class="space-y-2">
        <div v-for="(h, i) in headers" :key="i" class="flex gap-2 items-center">
          <input v-model="h.key" class="input-field flex-1 font-mono text-sm" placeholder="Header Name" />
          <input v-model="h.value" class="input-field flex-1 font-mono text-sm" placeholder="Header Value" />
          <button @click="headers.splice(i, 1)" class="btn btn-danger py-0.5 px-2 text-xs">
            <Trash2 :size="12"/>
          </button>
        </div>
        <div v-if="headers.length === 0" class="text-xs opacity-30 text-center py-2">暂无请求头，点击"添加"按钮</div>
      </div>
    </div>

    <!-- 请求体 / 响应体 -->
    <div class="two-col mb-4">
      <!-- 请求体 -->
      <div class="card flex flex-col min-h-0">
        <div class="flex items-center justify-between mb-2">
          <div class="label mb-0">请求体 (Body)</div>
          <button @click="formatBody('request')" class="btn btn-secondary py-0.5 px-2 text-xs">
            <Wand2 :size="12"/>格式化
          </button>
        </div>
        <textarea
          v-model="requestBody"
          class="textarea-field flex-1 min-h-[200px] font-mono text-sm"
          placeholder='{"key": "value"}'
          spellcheck="false"
        />
      </div>

      <!-- 响应体 -->
      <div class="card flex flex-col min-h-0">
        <div class="flex items-center justify-between mb-2">
          <div class="label mb-0">响应体 (Response)</div>
          <button @click="formatBody('response')" class="btn btn-secondary py-0.5 px-2 text-xs">
            <Wand2 :size="12"/>格式化
          </button>
        </div>
        <textarea
          v-model="responseBody"
          class="textarea-field flex-1 min-h-[200px] font-mono text-sm"
          placeholder='{"code": 200, "data": {}}'
          spellcheck="false"
        />
      </div>
    </div>

    <!-- 操作按钮 -->
    <div class="flex gap-2 mb-4">
      <button @click="generateDoc" class="btn btn-primary"><FileText :size="14"/>生成文档</button>
      <button @click="copyDoc" class="btn btn-secondary" :disabled="!markdownDoc">
        <Copy :size="14"/>复制 Markdown
      </button>
      <button @click="exportDoc" class="btn btn-secondary" :disabled="!markdownDoc">
        <Download :size="14"/>导出 .md
      </button>
    </div>

    <!-- Markdown 预览 -->
    <div v-if="markdownDoc" class="card flex-1 overflow-auto">
      <div class="label mb-2">Markdown 预览</div>
      <div class="code-output text-sm whitespace-pre-wrap font-mono">{{ markdownDoc }}</div>
    </div>
  </ToolPage>
</template>

<script setup lang="ts">
import ToolPage from '@/components/ToolPage.vue'
import { ref } from 'vue'
import { FileText, Plus, Trash2, Wand2, Copy, Download } from 'lucide-vue-next'
import { useAppStore } from '@/stores/app'

const appStore = useAppStore()

// 基本信息表单
const apiName = ref('')
const apiDesc = ref('')
const method = ref('POST')
const apiUrl = ref('')

// 请求头键值对
const headers = ref<{ key: string; value: string }[]>([
  { key: 'Content-Type', value: 'application/json' },
])

// 请求体和响应体
const requestBody = ref('')
const responseBody = ref('')

// 生成的 Markdown 文档
const markdownDoc = ref('')

// 添加请求头
function addHeader() {
  headers.value.push({ key: '', value: '' })
}

// 格式化 JSON
function formatBody(type: 'request' | 'response') {
  const target = type === 'request' ? requestBody : responseBody
  try {
    if (!target.value.trim()) return
    const parsed = JSON.parse(target.value)
    target.value = JSON.stringify(parsed, null, 2)
    appStore.showToast('success', 'JSON 格式化成功')
  } catch {
    appStore.showToast('error', 'JSON 格式不正确')
  }
}

// 生成 Markdown 文档
function generateDoc() {
  if (!apiName.value.trim() || !apiUrl.value.trim()) {
    appStore.showToast('error', '请填写接口名称和 URL')
    return
  }

  try {
    let md = ''
    md += `# ${apiName.value}\n\n`
    if (apiDesc.value.trim()) {
      md += `> ${apiDesc.value}\n\n`
    }
    md += `## 基本信息\n\n`
    md += `| 项目 | 值 |\n`
    md += `| --- | --- |\n`
    md += `| 接口名称 | ${apiName.value} |\n`
    md += `| 请求方法 | \`${method.value}\` |\n`
    md += `| 请求 URL | \`${apiUrl.value}\` |\n\n`

    // 请求头
    const validHeaders = headers.value.filter(h => h.key.trim())
    if (validHeaders.length > 0) {
      md += `## 请求头\n\n`
      md += `| 名称 | 值 |\n`
      md += `| --- | --- |\n`
      validHeaders.forEach(h => {
        md += `| ${h.key} | ${h.value || '-'} |\n`
      })
      md += '\n'
    }

    // 请求体
    if (requestBody.value.trim()) {
      md += `## 请求体\n\n`
      md += '```json\n'
      try {
        md += JSON.stringify(JSON.parse(requestBody.value), null, 2)
      } catch {
        md += requestBody.value
      }
      md += '\n```\n\n'
    }

    // 响应体
    if (responseBody.value.trim()) {
      md += `## 响应体\n\n`
      md += '```json\n'
      try {
        md += JSON.stringify(JSON.parse(responseBody.value), null, 2)
      } catch {
        md += responseBody.value
      }
      md += '\n```\n\n'
    }

    markdownDoc.value = md
    appStore.showToast('success', '文档生成成功')
  } catch (e) {
    appStore.showToast('error', '生成失败: ' + String(e))
  }
}

// 复制 Markdown
async function copyDoc() {
  if (!markdownDoc.value) return
  try {
    await navigator.clipboard.writeText(markdownDoc.value)
    appStore.showToast('success', '已复制到剪贴板')
  } catch {
    appStore.showToast('error', '复制失败')
  }
}

// 导出 Markdown 文件
function exportDoc() {
  if (!markdownDoc.value) return
  try {
    const blob = new Blob([markdownDoc.value], { type: 'text/markdown;charset=utf-8' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `${apiName.value || 'api_doc'}_${Date.now()}.md`
    a.click()
    URL.revokeObjectURL(url)
    appStore.showToast('success', '导出成功')
  } catch (e) {
    appStore.showToast('error', '导出失败: ' + String(e))
  }
}
</script>
