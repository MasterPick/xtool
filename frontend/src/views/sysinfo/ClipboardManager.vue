<template>
  <!-- 剪贴板管理工具 -->
  <div class="page-container">
    <div>
      <div class="page-title"><Clipboard :size="20" class="text-primary-400"/>剪贴板管理</div>
      <div class="page-desc">查看和管理剪贴板历史记录</div>
    </div>

    <!-- 工具栏 -->
    <div class="flex gap-2 mb-4 items-center">
      <input v-model="searchQuery" class="input-field flex-1" placeholder="搜索剪贴板内容..." />
      <button @click="toggleListening" :class="['btn', isListening ? 'btn-danger' : 'btn-primary']">
        <Activity :size="14" :class="isListening ? 'loading-spin' : ''"/>
        {{ isListening ? '停止监听' : '开始监听' }}
      </button>
      <button @click="clearAll" class="btn btn-secondary" :disabled="!history.length">
        <Trash2 :size="14"/>清空
      </button>
    </div>

    <!-- 监听状态 -->
    <div v-if="isListening" class="card mb-4 border-green-500/20">
      <div class="flex items-center gap-2 text-green-400 text-sm">
        <div class="w-2 h-2 rounded-full bg-green-400 animate-pulse" />
        正在监听剪贴板变化...
      </div>
    </div>

    <!-- 剪贴板历史列表 -->
    <div class="flex-1 overflow-auto space-y-2">
      <div v-for="(item, i) in filteredHistory" :key="item.id" class="card flex items-start gap-3">
        <!-- 置顶标识 -->
        <div class="shrink-0 pt-1">
          <button v-if="!item.pinned" @click="pinItem(i)" class="opacity-30 hover:opacity-100 transition-opacity" title="置顶">
            <Pin :size="14"/>
          </button>
          <button v-else @click="unpinItem(i)" class="text-primary-400" title="取消置顶">
            <Pin :size="14"/>
          </button>
        </div>

        <!-- 内容 -->
        <div class="flex-1 min-w-0">
          <div class="font-mono text-sm break-all whitespace-pre-wrap max-h-20 overflow-hidden">
            {{ item.content }}
          </div>
          <div class="flex items-center gap-3 mt-1">
            <span class="text-xs opacity-40">{{ item.time }}</span>
            <span class="text-xs opacity-40">{{ item.content.length }} 字符</span>
          </div>
        </div>

        <!-- 操作按钮 -->
        <div class="flex gap-1 shrink-0">
          <button @click="copyItem(item)" class="btn btn-secondary py-0.5 px-2 text-xs" title="复制">
            <Copy :size="11"/>
          </button>
          <button @click="deleteItem(i)" class="btn btn-danger py-0.5 px-2 text-xs" title="删除">
            <Trash2 :size="11"/>
          </button>
        </div>
      </div>

      <!-- 空状态 -->
      <div v-if="filteredHistory.length === 0" class="flex items-center justify-center opacity-30 pt-16">
        <div class="text-center">
          <Clipboard :size="32" class="mx-auto mb-2"/>
          <div class="text-sm">
            {{ searchQuery ? '没有匹配的记录' : '暂无剪贴板记录' }}
          </div>
          <div v-if="!searchQuery" class="text-xs mt-1 opacity-60">
            点击"开始监听"自动捕获剪贴板内容
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onUnmounted } from 'vue'
import { Clipboard, Activity, Trash2, Copy, Pin } from 'lucide-vue-next'
import { useAppStore } from '@/stores/app'

interface ClipboardItem {
  id: number
  content: string
  time: string
  pinned: boolean
}

const appStore = useAppStore()
const history = ref<ClipboardItem[]>([])
const searchQuery = ref('')
const isListening = ref(false)
let listenTimer: ReturnType<typeof setInterval> | null = null
let lastClipboardContent = ''
let idCounter = 0

// 过滤后的历史记录（置顶优先）
const filteredHistory = computed(() => {
  let list = history.value
  if (searchQuery.value.trim()) {
    const q = searchQuery.value.toLowerCase()
    list = list.filter(item => item.content.toLowerCase().includes(q))
  }
  // 置顶的排在前面
  return [...list].sort((a, b) => {
    if (a.pinned && !b.pinned) return -1
    if (!a.pinned && b.pinned) return 1
    return 0
  })
})

// 格式化时间
function formatTime(date: Date): string {
  const h = String(date.getHours()).padStart(2, '0')
  const m = String(date.getMinutes()).padStart(2, '0')
  const s = String(date.getSeconds()).padStart(2, '0')
  return `${h}:${m}:${s}`
}

// 添加记录
function addRecord(content: string) {
  if (!content.trim()) return
  // 避免重复
  if (history.value.length > 0 && history.value[0].content === content) return

  history.value.unshift({
    id: ++idCounter,
    content,
    time: formatTime(new Date()),
    pinned: false,
  })

  // 限制历史记录数量
  if (history.value.length > 200) {
    history.value = history.value.slice(0, 200)
  }
}

// 读取剪贴板内容
async function readClipboard(): Promise<string> {
  try {
    return await navigator.clipboard.readText()
  } catch {
    return ''
  }
}

// 切换监听状态
async function toggleListening() {
  if (isListening.value) {
    stopListening()
  } else {
    startListening()
  }
}

// 开始监听
async function startListening() {
  isListening.value = true
  try {
    // 尝试请求剪贴板权限
    lastClipboardContent = await readClipboard()
  } catch {
    // 权限可能被拒绝，继续尝试
  }

  // 使用轮询方式监听剪贴板变化
  // 注意：Wails 环境中 navigator.clipboard API 可能有限制
  listenTimer = setInterval(async () => {
    try {
      const content = await readClipboard()
      if (content && content !== lastClipboardContent) {
        lastClipboardContent = content
        addRecord(content)
      }
    } catch {
      // 静默处理读取失败
    }
  }, 1000)

  appStore.showToast('success', '已开始监听剪贴板')
}

// 停止监听
function stopListening() {
  isListening.value = false
  if (listenTimer) {
    clearInterval(listenTimer)
    listenTimer = null
  }
  appStore.showToast('info', '已停止监听剪贴板')
}

// 复制项目
async function copyItem(item: ClipboardItem) {
  try {
    await navigator.clipboard.writeText(item.content)
    lastClipboardContent = item.content
    appStore.showToast('success', '已复制到剪贴板')
  } catch {
    appStore.showToast('error', '复制失败')
  }
}

// 删除项目
function deleteItem(index: number) {
  // 找到在原始数组中的索引
  const item = filteredHistory.value[index]
  const realIndex = history.value.findIndex(h => h.id === item.id)
  if (realIndex > -1) {
    history.value.splice(realIndex, 1)
  }
}

// 置顶
function pinItem(filteredIndex: number) {
  const item = filteredHistory.value[filteredIndex]
  const realItem = history.value.find(h => h.id === item.id)
  if (realItem) realItem.pinned = true
}

// 取消置顶
function unpinItem(filteredIndex: number) {
  const item = filteredHistory.value[filteredIndex]
  const realItem = history.value.find(h => h.id === item.id)
  if (realItem) realItem.pinned = false
}

// 清空所有
function clearAll() {
  history.value = []
  appStore.showToast('info', '已清空剪贴板历史')
}

onUnmounted(() => {
  stopListening()
})
</script>
