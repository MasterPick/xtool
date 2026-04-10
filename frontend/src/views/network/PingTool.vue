<template>
  <div class="page-container">
    <div>
      <div class="page-title"><Signal :size="20" class="text-primary-400"/>Ping 测试</div>
      <div class="page-desc">测试主机连通性与网络延迟</div>
    </div>

    <div class="card mb-4">
      <div class="flex gap-2">
        <input v-model="host" class="input-field flex-1"
          placeholder="输入主机名或 IP，如 google.com / 8.8.8.8"
          @keyup.enter="startPing" />
        <button @click="startPing" class="btn btn-primary" :disabled="pinging">
          <Signal :size="14" :class="pinging ? 'loading-spin' : ''"/>
          {{ pinging ? '测试中...' : '开始测试' }}
        </button>
        <button v-if="pinging" @click="stopPing" class="btn btn-danger">
          <Square :size="14"/>停止
        </button>
        <button @click="clear" class="btn btn-secondary"><Trash2 :size="14"/>清空</button>
      </div>
      <div class="flex items-center gap-4 mt-2">
        <label class="flex items-center gap-2 text-sm cursor-pointer">
          <input type="checkbox" v-model="continuousPing" class="checkbox-field" />
          连续 Ping
        </label>
        <label v-if="continuousPing" class="flex items-center gap-2 text-sm">
          间隔:
          <input v-model.number="pingInterval" type="number" min="1" max="60" class="input-field w-16" />
          秒
        </label>
      </div>
    </div>

    <!-- 统计信息 -->
    <div v-if="results.length > 1" class="flex gap-3 mb-3">
      <div class="card text-center flex-1">
        <div class="text-lg font-bold text-green-400">{{ avgLatency }} ms</div>
        <div class="text-xs opacity-50">平均延迟</div>
      </div>
      <div class="card text-center flex-1">
        <div class="text-lg font-bold text-blue-400">{{ minLatency }} ms</div>
        <div class="text-xs opacity-50">最低延迟</div>
      </div>
      <div class="card text-center flex-1">
        <div class="text-lg font-bold text-red-400">{{ maxLatency }} ms</div>
        <div class="text-xs opacity-50">最高延迟</div>
      </div>
      <div class="card text-center flex-1">
        <div class="text-lg font-bold" :class="lossRate > 0 ? 'text-red-400' : 'text-green-400'">{{ lossRate }}%</div>
        <div class="text-xs opacity-50">丢包率</div>
      </div>
    </div>

    <!-- 结果列表 -->
    <div class="flex-1 overflow-auto space-y-2">
      <div v-for="(r, i) in results" :key="i"
           :class="['card flex items-center gap-3', r.alive ? 'border-green-500/20' : 'border-red-500/20']">
        <div :class="['w-2 h-2 rounded-full shrink-0', r.alive ? 'bg-green-400' : 'bg-red-400']" />
        <div class="flex-1 font-mono text-sm">{{ r.host }}</div>
        <div v-if="r.alive" class="text-green-400 text-sm font-mono">{{ r.latencyMs.toFixed(1) }} ms</div>
        <div v-else class="text-red-400 text-sm">不可达</div>
        <div :class="['badge', r.alive ? 'badge-success' : 'badge-error']">
          {{ r.alive ? '在线' : '离线' }}
        </div>
      </div>
      <div v-if="results.length === 0" class="flex-1 flex items-center justify-center opacity-30 pt-16">
        <div class="text-center">
          <Signal :size="32" class="mx-auto mb-2"/>
          <div class="text-sm">输入主机名开始测试</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onUnmounted } from 'vue'
import { Signal, Trash2, Square } from 'lucide-vue-next'
import { useAppStore } from '@/stores/app'
import { PingHost } from '../../../wailsjs/go/network/NetworkTools'

interface PingResult { host: string; alive: boolean; latencyMs: number; error: string }

const appStore = useAppStore()
const host    = ref('')
const pinging = ref(false)
const results = ref<PingResult[]>([])
const continuousPing = ref(false)
const pingInterval = ref(1)
let pingTimer: ReturnType<typeof setInterval> | null = null

const aliveResults = computed(() => results.value.filter(r => r.alive))
const avgLatency = computed(() => {
  if (!aliveResults.value.length) return '0'
  return (aliveResults.value.reduce((sum, r) => sum + r.latencyMs, 0) / aliveResults.value.length).toFixed(1)
})
const minLatency = computed(() => {
  if (!aliveResults.value.length) return '0'
  return Math.min(...aliveResults.value.map(r => r.latencyMs)).toFixed(1)
})
const maxLatency = computed(() => {
  if (!aliveResults.value.length) return '0'
  return Math.max(...aliveResults.value.map(r => r.latencyMs)).toFixed(1)
})
const lossRate = computed(() => {
  if (!results.value.length) return 0
  return Number(((results.value.filter(r => !r.alive).length / results.value.length) * 100).toFixed(1))
})

async function startPing() {
  if (!host.value.trim()) return
  pinging.value = true

  try {
    const r = await PingHost(host.value.trim()) as PingResult
    results.value.unshift(r)
  } catch (e) {
    results.value.unshift({ host: host.value.trim(), alive: false, latencyMs: 0, error: String(e) })
    appStore.showToast('error', 'Ping 失败: ' + String(e))
  }

  if (continuousPing.value && pinging.value) {
    pingTimer = setInterval(async () => {
      try {
        const r = await PingHost(host.value.trim()) as PingResult
        results.value.unshift(r)
        // 限制结果数量
        if (results.value.length > 200) results.value = results.value.slice(0, 200)
      } catch {
        // 忽略连续 ping 中的单次错误
      }
    }, pingInterval.value * 1000)
  } else {
    pinging.value = false
  }
}

function stopPing() {
  pinging.value = false
  if (pingTimer) {
    clearInterval(pingTimer)
    pingTimer = null
  }
}

function clear() {
  stopPing()
  results.value = []
}

onUnmounted(() => {
  stopPing()
})
</script>
