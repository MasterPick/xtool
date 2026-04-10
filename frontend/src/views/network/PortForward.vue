<template>
  <!-- 端口转发工具页面：管理 TCP 端口转发规则 -->
  <ToolPage title="端口转发" description="TCP 端口转发，支持多个规则同时运行">
    <template #actions>
      <button @click="showAddForm = !showAddForm" class="btn btn-primary">
        <Plus :size="14"/>{{ showAddForm ? '取消添加' : '添加规则' }}
      </button>
      <button @click="refreshList" class="btn btn-secondary"><RefreshCw :size="14"/>刷新</button>
    </template>

    <!-- 添加规则表单 -->
    <div v-if="showAddForm" class="add-form mb-4">
      <div class="flex items-end gap-3 flex-wrap">
        <div class="flex flex-col gap-1 flex-1 min-w-[140px]">
          <label class="text-xs opacity-60">监听端口</label>
          <input
            v-model="newRule.listenPort"
            type="text"
            class="input-field"
            placeholder="如 :8080 或 0.0.0.0:8080"
          />
        </div>
        <div class="flex flex-col gap-1 flex-1 min-w-[200px]">
          <label class="text-xs opacity-60">目标地址</label>
          <input
            v-model="newRule.targetAddr"
            type="text"
            class="input-field"
            placeholder="如 192.168.1.1:80 或 localhost:3000"
          />
        </div>
        <button @click="addForward" class="btn btn-primary">
          <Play :size="14"/>启动转发
        </button>
      </div>
    </div>

    <!-- 转发规则列表 -->
    <div class="flex-1 min-h-0 overflow-auto">
      <!-- 空状态 -->
      <div v-if="forwards.length === 0" class="empty-state">
        <ArrowLeftRight :size="40" class="opacity-20 mb-3" />
        <p class="opacity-40">暂无转发规则</p>
        <p class="text-xs opacity-30 mt-1">点击"添加规则"创建新的端口转发</p>
      </div>

      <!-- 规则卡片列表 -->
      <div v-else class="space-y-2">
        <div
          v-for="item in forwards"
          :key="item.id"
          class="forward-card"
        >
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-3">
              <!-- 状态指示灯 -->
              <span
                :class="['status-dot', item.status === 'running' ? 'status-running' : item.status === 'error' ? 'status-error' : 'status-stopped']"
              />
              <!-- 转发信息 -->
              <div>
                <div class="flex items-center gap-2 text-sm font-medium">
                  <span class="font-mono">{{ item.listenAddr }}</span>
                  <ArrowRight :size="14" class="opacity-40" />
                  <span class="font-mono">{{ item.targetAddr }}</span>
                </div>
                <div class="flex items-center gap-3 mt-1 text-xs opacity-50">
                  <span>连接数: {{ item.connCount }}</span>
                  <span>创建时间: {{ item.createdAt }}</span>
                  <span v-if="item.error" class="text-red-400">{{ item.error }}</span>
                </div>
              </div>
            </div>
            <!-- 操作按钮 -->
            <div class="flex items-center gap-2">
              <span
                :class="['badge text-xs', item.status === 'running' ? 'badge-success' : item.status === 'error' ? 'badge-error' : 'badge-info']"
              >
                {{ statusLabel(item.status) }}
              </span>
              <button
                v-if="item.status === 'running'"
                @click="stopForward(item.id)"
                class="btn btn-sm btn-danger"
              >
                <Square :size="12"/>停止
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </ToolPage>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Plus, Play, Square, RefreshCw, ArrowLeftRight, ArrowRight } from 'lucide-vue-next'
import ToolPage from '@/components/ToolPage.vue'
import { useAppStore } from '@/stores/app'
import { StartPortForward, StopPortForward, GetPortForwards } from '../../../wailsjs/go/network/NetworkTools'

const appStore = useAppStore()

// 转发规则列表
const forwards = ref<any[]>([])

// 是否显示添加表单
const showAddForm = ref(false)

// 新规则表单数据
const newRule = ref({
  listenPort: '',
  targetAddr: '',
})

// 状态标签映射
function statusLabel(status: string) {
  const map: Record<string, string> = {
    running: '运行中',
    stopped: '已停止',
    error: '错误',
  }
  return map[status] || status
}

// 刷新转发规则列表
async function refreshList() {
  try {
    const list = await GetPortForwards()
    forwards.value = list || []
  } catch (e) {
    appStore.showToast('error', '获取转发列表失败: ' + String(e))
  }
}

// 添加并启动转发规则
async function addForward() {
  const listenPort = newRule.value.listenPort.trim()
  const targetAddr = newRule.value.targetAddr.trim()

  if (!listenPort) {
    appStore.showToast('error', '请输入监听端口')
    return
  }
  if (!targetAddr) {
    appStore.showToast('error', '请输入目标地址')
    return
  }

  // 如果只输入了端口号，自动加上冒号前缀
  let listenAddr = listenPort
  if (/^\d+$/.test(listenPort)) {
    listenAddr = ':' + listenPort
  }

  try {
    const res = await StartPortForward(listenAddr, targetAddr)
    if (res.success) {
      appStore.showToast('success', `转发规则已启动：${listenAddr} -> ${targetAddr}`)
      forwards.value = res.forwards || []
      // 清空表单
      newRule.value.listenPort = ''
      newRule.value.targetAddr = ''
      showAddForm.value = false
    } else {
      appStore.showToast('error', res.error || '启动失败')
    }
  } catch (e) {
    appStore.showToast('error', '启动转发失败: ' + String(e))
  }
}

// 停止转发规则
async function stopForward(id: string) {
  try {
    const res = await StopPortForward(id)
    if (res.success) {
      appStore.showToast('success', '转发规则已停止')
      forwards.value = res.forwards || []
    } else {
      appStore.showToast('error', res.error || '停止失败')
    }
  } catch (e) {
    appStore.showToast('error', '停止转发失败: ' + String(e))
  }
}

// 页面加载时获取转发列表
onMounted(() => {
  refreshList()
})
</script>

<style scoped>
/* 添加规则表单 */
.add-form {
  padding: 16px;
  background: var(--bg-hover, rgba(255, 255, 255, 0.03));
  border: 1px solid var(--border-color, rgba(255, 255, 255, 0.08));
  border-radius: var(--radius-lg, 12px);
}

/* 转发规则卡片 */
.forward-card {
  padding: 12px 16px;
  background: var(--bg-hover, rgba(255, 255, 255, 0.03));
  border: 1px solid var(--border-color, rgba(255, 255, 255, 0.08));
  border-radius: var(--radius-md, 8px);
  transition: all var(--transition-fast, 0.15s ease);
}

.forward-card:hover {
  background: var(--bg-active, rgba(255, 255, 255, 0.06));
  border-color: var(--accent, #00ADD8);
}

/* 状态指示灯 */
.status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  flex-shrink: 0;
}

.status-running {
  background: #22c55e;
  box-shadow: 0 0 6px rgba(34, 197, 94, 0.5);
  animation: pulse 2s infinite;
}

.status-stopped {
  background: #6b7280;
}

.status-error {
  background: var(--danger);
  box-shadow: 0 0 6px rgba(239, 68, 68, 0.5);
}

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.5; }
}

/* 空状态 */
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  min-height: 200px;
}

/* 危险按钮 */
.btn-danger {
  background: rgba(239, 68, 68, 0.15);
  color: #f87171;
  border: 1px solid rgba(239, 68, 68, 0.2);
}

.btn-danger:hover {
  background: rgba(239, 68, 68, 0.25);
}

/* 小按钮 */
.btn-sm {
  padding: 4px 10px;
  font-size: 12px;
}
</style>
