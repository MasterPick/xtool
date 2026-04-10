<template>
  <!-- 进制转换工具：支持二进制、八进制、十进制、十六进制互转 -->
  <div class="page-container">
    <div>
      <div class="page-title"><Binary :size="20" class="text-primary-400"/>进制转换</div>
      <div class="page-desc">二进制 · 八进制 · 十进制 · 十六进制 互转</div>
    </div>

    <!-- 输入区域 -->
    <div class="card mb-4">
      <div class="flex gap-3 items-end">
        <div class="flex-1">
          <div class="label mb-1">输入值</div>
          <input
            v-model="inputValue"
            class="input-field font-mono"
            placeholder="输入要转换的数值..."
            @input="convert"
            spellcheck="false"
          />
        </div>
        <div class="w-36">
          <div class="label mb-1">输入进制</div>
          <select v-model="fromBase" class="input-field" @change="convert">
            <option value="2">二进制 (BIN)</option>
            <option value="8">八进制 (OCT)</option>
            <option value="10">十进制 (DEC)</option>
            <option value="16">十六进制 (HEX)</option>
          </select>
        </div>
      </div>
      <div v-if="errorMsg" class="text-red-400 text-xs mt-2">{{ errorMsg }}</div>
    </div>

    <!-- 转换结果 -->
    <div class="space-y-3 flex-1 overflow-auto">
      <div v-for="item in resultItems" :key="item.base" class="card">
        <div class="flex items-center justify-between mb-2">
          <span class="font-mono font-bold text-sm text-primary-400">{{ item.label }}</span>
          <button @click="copyResult(item.value)" class="btn btn-secondary py-0.5 px-2 text-xs">
            <Copy :size="11"/>复制
          </button>
        </div>
        <div class="code-output text-sm break-all font-mono">
          <span v-if="!item.value" class="opacity-30">等待输入...</span>
          <span v-else>{{ item.value }}</span>
        </div>
      </div>

      <!-- 转换过程 -->
      <div v-if="conversionProcess" class="card">
        <div class="label mb-2">转换过程</div>
        <div class="code-output text-sm whitespace-pre-wrap font-mono opacity-80">{{ conversionProcess }}</div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { Binary, Copy } from 'lucide-vue-next'
import { useAppStore } from '@/stores/app'

const appStore = useAppStore()
const inputValue = ref('')
const fromBase = ref('10')
const errorMsg = ref('')

// 各进制结果
const binResult = ref('')
const octResult = ref('')
const decResult = ref('')
const hexResult = ref('')
const conversionProcess = ref('')

// 结果列表
const resultItems = computed(() => [
  { base: 2, label: '二进制 (BIN)', value: binResult.value },
  { base: 8, label: '八进制 (OCT)', value: octResult.value },
  { base: 10, label: '十进制 (DEC)', value: decResult.value },
  { base: 16, label: '十六进制 (HEX)', value: hexResult.value },
])

// 执行转换
function convert() {
  errorMsg.value = ''
  conversionProcess.value = ''
  binResult.value = ''
  octResult.value = ''
  decResult.value = ''
  hexResult.value = ''

  const val = inputValue.value.trim()
  if (!val) return

  // 验证输入
  const base = parseInt(fromBase.value)
  const pattern: Record<number, RegExp> = {
    2: /^[01]+$/,
    8: /^[0-7]+$/,
    10: /^[0-9]+$/,
    16: /^[0-9a-fA-F]+$/,
  }

  if (!pattern[base].test(val)) {
    errorMsg.value = `输入值不是有效的${base === 2 ? '二进制' : base === 8 ? '八进制' : base === 10 ? '十进制' : '十六进制'}数`
    return
  }

  try {
    // 先转为十进制
    const decimal = parseInt(val, base)
    if (isNaN(decimal)) {
      errorMsg.value = '无法解析输入值'
      return
    }

    // 生成转换过程
    const baseNames: Record<number, string> = { 2: '二进制', 8: '八进制', 10: '十进制', 16: '十六进制' }
    let process = `输入: ${val} (${baseNames[base]})\n`
    process += `转为十进制: ${val} × ${base}^各权位 = ${decimal}\n`
    process += `\n转换结果:\n`

    // 计算各进制
    binResult.value = decimal.toString(2)
    octResult.value = decimal.toString(8)
    decResult.value = decimal.toString(10)
    hexResult.value = decimal.toString(16).toUpperCase()

    process += `  二进制:   ${binResult.value}\n`
    process += `  八进制:   ${octResult.value}\n`
    process += `  十进制:   ${decResult.value}\n`
    process += `  十六进制: ${hexResult.value}`

    conversionProcess.value = process
  } catch (e) {
    errorMsg.value = '转换失败: ' + String(e)
  }
}

// 复制结果
async function copyResult(val: string) {
  if (!val) return
  try {
    await navigator.clipboard.writeText(val)
    appStore.showToast('success', '已复制到剪贴板')
  } catch {
    appStore.showToast('error', '复制失败')
  }
}
</script>
