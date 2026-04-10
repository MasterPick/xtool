<template>
  <div class="page-container">
    <div class="page-title"><Cpu :size="20" class="text-primary-400"/>进程管理</div>
    <div class="page-desc">查看并管理系统运行中的进程</div>
    <div class="flex gap-2 mb-4">
      <input v-model="filter" class="input-field flex-1" placeholder="按进程名过滤..."/>
      <select v-model="sortField" class="input-field w-28">
        <option value="cpu">按 CPU</option>
        <option value="memory">按内存</option>
        <option value="pid">按 PID</option>
        <option value="name">按名称</option>
      </select>
      <button @click="toggleSortOrder" class="btn btn-secondary" :title="sortAsc ? '升序' : '降序'">
        <ArrowUpDown :size="14"/>
      </button>
      <button @click="loadProcs" class="btn btn-primary" :disabled="loading">
        <RefreshCw :size="14" :class="loading?'loading-spin':''"/>刷新
      </button>
    </div>
    <div class="flex-1 overflow-auto">
      <table class="w-full text-sm">
        <thead>
          <tr class="text-left opacity-50 text-xs border-b" style="border-color:var(--border-color)">
            <th class="pb-2 pr-4">进程名</th><th class="pb-2 pr-4">PID</th>
            <th class="pb-2 pr-4">CPU%</th><th class="pb-2 pr-4">内存%</th>
            <th class="pb-2">操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="p in filteredProcs" :key="p.pid"
              class="border-b hover:bg-white/3 transition-colors"
              style="border-color:var(--border-color)">
            <td class="py-1.5 pr-4">{{ p.name }}</td>
            <td class="py-1.5 pr-4 font-mono text-xs opacity-60">{{ p.pid }}</td>
            <td class="py-1.5 pr-4">
              <span :class="p.cpu>50?'text-red-400':p.cpu>20?'text-yellow-400':'text-green-400'">
                {{ p.cpu.toFixed(1) }}%
              </span>
            </td>
            <td class="py-1.5 pr-4 opacity-70">{{ p.memory.toFixed(1) }}%</td>
            <td class="py-1.5">
              <button @click="confirmKill(p)" class="btn btn-danger py-0.5 px-2 text-xs">
                <X :size="10"/>结束
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- 结束进程确认弹窗 -->
    <div v-if="killConfirm" class="modal-overlay" @click.self="killConfirm=null">
      <div class="modal-box card w-[400px]">
        <div class="font-semibold mb-3 text-red-400">确认结束进程</div>
        <div class="text-sm mb-4">
          确定要结束进程 <span class="font-mono text-primary-400">{{ killConfirm.name }}</span>
          (PID: <span class="font-mono">{{ killConfirm.pid }}</span>) 吗？
          <div class="mt-2 text-xs opacity-50">此操作不可撤销，可能导致数据丢失。</div>
        </div>
        <div class="flex gap-2 justify-end">
          <button @click="killConfirm=null" class="btn btn-secondary">取消</button>
          <button @click="doKill" class="btn btn-danger">确认结束</button>
        </div>
      </div>
    </div>
  </div>
</template>
<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { Cpu, RefreshCw, X, ArrowUpDown } from 'lucide-vue-next'
import { useAppStore } from '@/stores/app'
import { GetProcessList, KillProcess } from '../../../wailsjs/go/sysinfo/SysInfo'
const appStore = useAppStore()
const procs = ref<any[]>([]), filter = ref(''), loading = ref(false)
const sortField = ref('cpu')
const sortAsc = ref(false)
const killConfirm = ref<any>(null)

const filteredProcs = computed(() => {
  let list = filter.value
    ? procs.value.filter(p => p.name.toLowerCase().includes(filter.value.toLowerCase()))
    : procs.value.slice(0, 500)

  // 排序
  list = [...list].sort((a, b) => {
    let va: any, vb: any
    switch (sortField.value) {
      case 'cpu': va = a.cpu; vb = b.cpu; break
      case 'memory': va = a.memory; vb = b.memory; break
      case 'pid': va = a.pid; vb = b.pid; break
      case 'name': va = a.name; vb = b.name; break
      default: va = a.cpu; vb = b.cpu
    }
    if (typeof va === 'string') {
      return sortAsc.value ? va.localeCompare(vb) : vb.localeCompare(va)
    }
    return sortAsc.value ? va - vb : vb - va
  })

  return list
})

function toggleSortOrder() {
  sortAsc.value = !sortAsc.value
}

async function loadProcs() {
  loading.value = true
  try {
    procs.value = (await GetProcessList() as any[] || []).sort((a: any, b: any) => b.cpu - a.cpu)
  } catch (e) {
    appStore.showToast('error', '加载进程列表失败: ' + String(e))
  } finally {
    loading.value = false
  }
}

function confirmKill(p: any) {
  killConfirm.value = p
}

async function doKill() {
  if (!killConfirm.value) return
  const pid = killConfirm.value.pid
  const name = killConfirm.value.name
  killConfirm.value = null
  try {
    await KillProcess(pid)
    appStore.showToast('success', `进程 ${name} (${pid}) 已终止`)
    loadProcs()
  } catch (e) {
    appStore.showToast('error', `终止进程失败: ${e}`)
  }
}

onMounted(loadProcs)
</script>
<style scoped>
.modal-overlay { position:fixed;inset:0;background:rgba(0,0,0,0.5);display:flex;align-items:center;justify-content:center;z-index:100; }
</style>
