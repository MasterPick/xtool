<template>
  <ToolPage title="密码生成" description="安全随机密码生成器">
<!-- 密码生成器 -->
  
    

    <!-- 配置区域 -->
    <div class="card mb-4">
      <!-- 密码长度 -->
      <div class="mb-4">
        <div class="flex items-center justify-between mb-2">
          <div class="label mb-0">密码长度</div>
          <span class="font-mono text-primary-400 font-bold text-lg">{{ passwordLength }}</span>
        </div>
        <input
          v-model.number="passwordLength"
          type="range"
          min="4"
          max="128"
          class="w-full h-2 bg-white/10 rounded-full appearance-none cursor-pointer accent-primary-400"
        />
        <div class="flex justify-between text-xs opacity-30 mt-1">
          <span>4</span>
          <span>128</span>
        </div>
      </div>

      <!-- 字符类型 -->
      <div class="grid grid-cols-2 gap-3 mb-4">
        <label class="flex items-center gap-2 text-sm cursor-pointer">
          <input type="checkbox" v-model="options.uppercase" class="checkbox-field" />
          大写字母 (A-Z)
        </label>
        <label class="flex items-center gap-2 text-sm cursor-pointer">
          <input type="checkbox" v-model="options.lowercase" class="checkbox-field" />
          小写字母 (a-z)
        </label>
        <label class="flex items-center gap-2 text-sm cursor-pointer">
          <input type="checkbox" v-model="options.numbers" class="checkbox-field" />
          数字 (0-9)
        </label>
        <label class="flex items-center gap-2 text-sm cursor-pointer">
          <input type="checkbox" v-model="options.symbols" class="checkbox-field" />
          特殊符号 (!@#$...)
        </label>
      </div>

      <!-- 排除字符 -->
      <div class="mb-4">
        <div class="label mb-1">排除字符</div>
        <input v-model="excludeChars" class="input-field font-mono" placeholder="输入要排除的字符，如 0OIl1..." />
      </div>

      <!-- 批量生成数量 -->
      <div class="flex gap-3 items-end mb-4">
        <div class="w-36">
          <div class="label mb-1">生成数量</div>
          <input v-model.number="batchCount" type="number" min="1" max="20" class="input-field" />
        </div>
        <button @click="generate" class="btn btn-primary">
          <Wand2 :size="14"/>生成密码
        </button>
      </div>
    </div>

    <!-- 密码强度指示器 -->
    <div v-if="passwords.length > 0" class="card mb-4">
      <div class="label mb-2">密码强度</div>
      <div class="flex items-center gap-3">
        <div class="flex-1">
          <div class="w-full h-3 bg-white/5 rounded-full overflow-hidden">
            <div
              class="h-full rounded-full transition-all duration-300"
              :class="strengthBarClass"
              :style="{ width: strengthPercent + '%' }"
            ></div>
          </div>
        </div>
        <span :class="['text-sm font-bold', strengthTextClass]">{{ strengthLabel }}</span>
      </div>
      <div class="text-xs opacity-40 mt-1">{{ strengthTip }}</div>
    </div>

    <!-- 生成结果 -->
    <div class="flex-1 overflow-auto space-y-2">
      <div v-for="(pwd, i) in passwords" :key="i" class="card flex items-center gap-3">
        <span class="text-xs opacity-40 w-8 shrink-0 text-right">{{ i + 1 }}</span>
        <div class="flex-1 font-mono text-sm break-all select-all">{{ pwd }}</div>
        <button @click="copyOne(pwd)" class="btn btn-secondary py-0.5 px-2 text-xs shrink-0">
          <Copy :size="11"/>复制
        </button>
      </div>

      <div v-if="passwords.length === 0" class="flex items-center justify-center opacity-30 pt-16">
        <div class="text-center">
          <KeyRound :size="32" class="mx-auto mb-2"/>
          <div class="text-sm">点击"生成密码"按钮</div>
        </div>
      </div>
    </div>

    <!-- 历史记录 -->
    <div v-if="history.length > 0" class="mt-4">
      <div class="flex items-center justify-between mb-2">
        <div class="label mb-0">生成历史</div>
        <button @click="history = []" class="btn btn-secondary py-0.5 px-2 text-xs">
          <Trash2 :size="11"/>清空
        </button>
      </div>
      <div class="space-y-1 max-h-40 overflow-auto">
        <div v-for="(h, i) in history" :key="i" class="flex items-center gap-2 text-xs opacity-50 px-2 py-1 rounded hover:bg-white/3">
          <span class="font-mono truncate flex-1">{{ h.password }}</span>
          <span>{{ h.time }}</span>
          <button @click="copyOne(h.password)" class="opacity-50 hover:opacity-100">
            <Copy :size="10"/>
          </button>
        </div>
      </div>
    </div>
  </ToolPage>
</template>

<script setup lang="ts">
import ToolPage from '@/components/ToolPage.vue'
import { ref, reactive, computed } from 'vue'
import { KeyRound, Wand2, Copy, Trash2 } from 'lucide-vue-next'
import { useAppStore } from '@/stores/app'

const appStore = useAppStore()
const passwordLength = ref(16)
const batchCount = ref(5)
const excludeChars = ref('')
const passwords = ref<string[]>([])
const history = ref<{ password: string; time: string }[]>([])

// 字符类型选项
const options = reactive({
  uppercase: true,
  lowercase: true,
  numbers: true,
  symbols: true,
})

// 字符集
const charSets = {
  uppercase: 'ABCDEFGHIJKLMNOPQRSTUVWXYZ',
  lowercase: 'abcdefghijklmnopqrstuvwxyz',
  numbers: '0123456789',
  symbols: '!@#$%^&*()_+-=[]{}|;:,.<>?',
}

// 密码强度评估
const strengthPercent = computed(() => {
  if (passwords.value.length === 0) return 0
  const pwd = passwords.value[0]
  let score = 0
  if (pwd.length >= 8) score += 20
  if (pwd.length >= 12) score += 10
  if (pwd.length >= 16) score += 10
  if (/[a-z]/.test(pwd)) score += 15
  if (/[A-Z]/.test(pwd)) score += 15
  if (/[0-9]/.test(pwd)) score += 15
  if (/[^a-zA-Z0-9]/.test(pwd)) score += 15
  return Math.min(100, score)
})

const strengthLabel = computed(() => {
  const p = strengthPercent.value
  if (p < 30) return '弱'
  if (p < 50) return '中'
  if (p < 80) return '强'
  return '极强'
})

const strengthBarClass = computed(() => {
  const p = strengthPercent.value
  if (p < 30) return 'bg-red-400'
  if (p < 50) return 'bg-yellow-400'
  if (p < 80) return 'bg-blue-400'
  return 'bg-green-400'
})

const strengthTextClass = computed(() => {
  const p = strengthPercent.value
  if (p < 30) return 'text-red-400'
  if (p < 50) return 'text-yellow-400'
  if (p < 80) return 'text-blue-400'
  return 'text-green-400'
})

const strengthTip = computed(() => {
  const p = strengthPercent.value
  if (p < 30) return '密码强度较弱，建议增加长度和字符类型'
  if (p < 50) return '密码强度一般，建议使用更多字符类型'
  if (p < 80) return '密码强度较好，适合一般用途'
  return '密码强度极佳，适合高安全场景'
})

// 生成单个密码
function generateOne(): string {
  let charset = ''
  if (options.uppercase) charset += charSets.uppercase
  if (options.lowercase) charset += charSets.lowercase
  if (options.numbers) charset += charSets.numbers
  if (options.symbols) charset += charSets.symbols

  // 排除指定字符
  if (excludeChars.value) {
    const excludeSet = new Set(excludeChars.value.split(''))
    charset = charset.split('').filter(c => !excludeSet.has(c)).join('')
  }

  if (!charset) {
    appStore.showToast('error', '请至少选择一种字符类型')
    return ''
  }

  // 使用 crypto API 生成安全随机数
  const array = new Uint32Array(passwordLength.value)
  crypto.getRandomValues(array)
  return Array.from(array, n => charset[n % charset.length]).join('')
}

// 批量生成
function generate() {
  const n = Math.max(1, Math.min(20, batchCount.value))
  passwords.value = []
  for (let i = 0; i < n; i++) {
    const pwd = generateOne()
    if (pwd) passwords.value.push(pwd)
  }

  // 添加到历史记录
  const now = new Date()
  const time = `${String(now.getHours()).padStart(2, '0')}:${String(now.getMinutes()).padStart(2, '0')}:${String(now.getSeconds()).padStart(2, '0')}`
  passwords.value.forEach(pwd => {
    history.value.unshift({ password: pwd, time })
  })

  // 限制历史记录数量
  if (history.value.length > 100) {
    history.value = history.value.slice(0, 100)
  }

  if (passwords.value.length > 0) {
    appStore.showToast('success', `已生成 ${passwords.value.length} 个密码`)
  }
}

// 复制单个密码
async function copyOne(pwd: string) {
  try {
    await navigator.clipboard.writeText(pwd)
    appStore.showToast('success', '已复制到剪贴板')
  } catch {
    appStore.showToast('error', '复制失败')
  }
}
</script>
