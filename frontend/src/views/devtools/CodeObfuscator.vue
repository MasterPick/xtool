<template>
  <ToolPage title="代码混淆" description="JavaScript 代码混淆">
<!-- 代码混淆工具 -->
  
    

    <!-- 配置区域 -->
    <div class="card mb-4">
      <div class="flex gap-3 items-end flex-wrap">
        <div class="w-36">
          <div class="label mb-1">语言</div>
          <select v-model="language" class="input-field">
            <option value="javascript">JavaScript</option>
          </select>
        </div>
        <div class="flex-1" />
        <button @click="obfuscate" class="btn btn-primary" :disabled="obfuscating || !inputCode.trim()">
          <Shield :size="14"/>{{ obfuscating ? '混淆中...' : '开始混淆' }}
        </button>
        <button @click="copyResult" class="btn btn-secondary" :disabled="!outputCode">
          <Copy :size="14"/>复制结果
        </button>
        <button @click="downloadResult" class="btn btn-secondary" :disabled="!outputCode">
          <Download :size="14"/>下载
        </button>
      </div>
    </div>

    <!-- 混淆选项 -->
    <div class="card mb-4">
      <div class="label mb-3">混淆选项</div>
      <div class="grid grid-cols-2 gap-3">
        <label class="flex items-center gap-2 text-sm cursor-pointer">
          <input type="checkbox" v-model="options.renameVariables" class="checkbox-field" />
          变量重命名
        </label>
        <label class="flex items-center gap-2 text-sm cursor-pointer">
          <input type="checkbox" v-model="options.encodeStrings" class="checkbox-field" />
          字符串编码
        </label>
        <label class="flex items-center gap-2 text-sm cursor-pointer">
          <input type="checkbox" v-model="options.controlFlow" class="checkbox-field" />
          控制流混淆
        </label>
        <label class="flex items-center gap-2 text-sm cursor-pointer">
          <input type="checkbox" v-model="options.deadCode" class="checkbox-field" />
          插入死代码
        </label>
        <label class="flex items-center gap-2 text-sm cursor-pointer">
          <input type="checkbox" v-model="options.selfDefending" class="checkbox-field" />
          自我保护
        </label>
        <label class="flex items-center gap-2 text-sm cursor-pointer">
          <input type="checkbox" v-model="options.compact" class="checkbox-field" />
          压缩输出
        </label>
      </div>
    </div>

    <!-- 双栏编辑区 -->
    <div class="two-col flex-1 min-h-0">
      <!-- 左侧：输入 -->
      <div class="flex flex-col gap-2 min-h-0">
        <div class="label">
          <span>输入代码</span>
          <span class="text-xs opacity-50">{{ inputCode.length }} 字符</span>
        </div>
        <textarea
          v-model="inputCode"
          class="textarea-field flex-1 min-h-0 font-mono text-sm"
          placeholder="在此粘贴要混淆的 JavaScript 代码..."
          spellcheck="false"
        />
      </div>

      <!-- 右侧：输出 -->
      <div class="flex flex-col gap-2 min-h-0">
        <div class="label">
          <span>混淆结果</span>
          <span v-if="outputCode" class="text-xs opacity-50">{{ outputCode.length }} 字符</span>
        </div>
        <div class="code-output flex-1 min-h-0 overflow-auto font-mono text-sm whitespace-pre-wrap">
          <span v-if="!outputCode" class="opacity-30">混淆结果将显示在这里...</span>
          <span v-else>{{ outputCode }}</span>
        </div>
      </div>
    </div>
  </ToolPage>
</template>

<script setup lang="ts">
import ToolPage from '@/components/ToolPage.vue'
import { ref, reactive } from 'vue'
import { Shield, Copy, Download } from 'lucide-vue-next'
import { useAppStore } from '@/stores/app'

const appStore = useAppStore()
const language = ref('javascript')
const inputCode = ref('')
const outputCode = ref('')
const obfuscating = ref(false)

// 混淆选项
const options = reactive({
  renameVariables: true,
  encodeStrings: false,
  controlFlow: false,
  deadCode: false,
  selfDefending: false,
  compact: true,
})

// 简单的 JavaScript 混淆器（纯前端实现）
function obfuscateJS(code: string): string {
  let result = code

  // 1. 变量重命名
  if (options.renameVariables) {
    // 收集变量名（简单正则匹配）
    const varRegex = /\b(var|let|const)\s+([a-zA-Z_$][a-zA-Z0-9_$]*)/g
    const funcRegex = /\bfunction\s+([a-zA-Z_$][a-zA-Z0-9_$]*)/g
    const paramRegex = /\(([^)]*)\)/g
    const reserved = new Set([
      'var', 'let', 'const', 'function', 'return', 'if', 'else', 'for', 'while',
      'do', 'switch', 'case', 'break', 'continue', 'new', 'this', 'class',
      'extends', 'import', 'export', 'default', 'try', 'catch', 'finally',
      'throw', 'typeof', 'instanceof', 'in', 'of', 'void', 'delete',
      'true', 'false', 'null', 'undefined', 'NaN', 'Infinity',
      'console', 'window', 'document', 'Math', 'Array', 'Object', 'String',
      'Number', 'Boolean', 'Date', 'JSON', 'Promise', 'Error',
      'parseInt', 'parseFloat', 'isNaN', 'isFinite', 'encodeURI', 'decodeURI',
      'alert', 'log', 'length', 'push', 'pop', 'shift', 'unshift', 'map',
      'filter', 'reduce', 'forEach', 'find', 'indexOf', 'slice', 'splice',
      'join', 'split', 'replace', 'match', 'test', 'toString', 'valueOf',
      'prototype', 'constructor', 'apply', 'call', 'bind', 'then', 'catch',
    ])

    const nameMap = new Map<string, string>()
    let counter = 0

    function getObfuscatedName(): string {
      const chars = '_abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ'
      let name = '_0x'
      let n = counter++
      do {
        name += chars[n % chars.length]
        n = Math.floor(n / chars.length)
      } while (n > 0)
      return name
    }

    // 提取变量名
    const varNames = new Set<string>()
    let match
    while ((match = varRegex.exec(code)) !== null) {
      if (!reserved.has(match[2])) varNames.add(match[2])
    }
    while ((match = funcRegex.exec(code)) !== null) {
      if (!reserved.has(match[1])) varNames.add(match[1])
    }

    // 生成替换映射
    varNames.forEach(name => {
      nameMap.set(name, getObfuscatedName())
    })

    // 执行替换（使用词边界确保完整匹配）
    nameMap.forEach((newName, oldName) => {
      const regex = new RegExp(`\\b${oldName}\\b`, 'g')
      result = result.replace(regex, newName)
    })
  }

  // 2. 字符串编码
  if (options.encodeStrings) {
    // 将字符串字面量转为 Unicode 转义
    result = result.replace(/'([^']*)'/g, (_: string, str: string) => {
      return "'" + Array.from(str).map(c => {
        const code = c.charCodeAt(0)
        return code > 127 ? `\\u${code.toString(16).padStart(4, '0')}` : c
      }).join('') + "'"
    })
    result = result.replace(/"([^"]*)"/g, (_: string, str: string) => {
      return '"' + Array.from(str).map(c => {
        const code = c.charCodeAt(0)
        return code > 127 ? `\\u${code.toString(16).padStart(4, '0')}` : c
      }).join('') + '"'
    })
  }

  // 3. 控制流混淆
  if (options.controlFlow) {
    // 将简单的 if-else 转为三元运算符或 switch-case（简化版）
    result = result.replace(/if\s*\(([^)]+)\)\s*\{([^}]*)\}\s*else\s*\{([^}]*)\}/g,
      'switch($1){case true:$2;break;default:$3}')
  }

  // 4. 插入死代码
  if (options.deadCode) {
    const deadCodes = [
      'if(false){var _dx=0;_dx++;}',
      'void(0);',
      'if(typeof _undefined==="undefined")var _undefined;',
    ]
    // 在每个语句后插入死代码
    const lines = result.split('\n')
    const newLines: string[] = []
    lines.forEach(line => {
      newLines.push(line)
      if (line.trim().endsWith(';') || line.trim().endsWith('{')) {
        if (Math.random() > 0.6) {
          newLines.push('  ' + deadCodes[Math.floor(Math.random() * deadCodes.length)])
        }
      }
    })
    result = newLines.join('\n')
  }

  // 5. 自我保护
  if (options.selfDefending) {
    const protector = `;(function(){var _s=['\\x63\\x6f\\x6e\\x73\\x6f\\x6c\\x65'];var _f=function(){try{setInterval(function(){},0x1f4)}catch(_e){}};_f()})();\n`
    result = protector + result
  }

  // 6. 压缩
  if (options.compact) {
    // 移除多余空行和注释
    result = result
      .replace(/\/\/.*$/gm, '')           // 移除单行注释
      .replace(/\/\*[\s\S]*?\*\//g, '')    // 移除多行注释
      .replace(/\n\s*\n/g, '\n')           // 移除空行
      .replace(/^\s+/gm, '')               // 移除行首空格
      .trim()
  }

  return result
}

// 执行混淆
function obfuscate() {
  if (!inputCode.value.trim()) {
    appStore.showToast('error', '请输入要混淆的代码')
    return
  }

  obfuscating.value = true
  try {
    if (language.value === 'javascript') {
      outputCode.value = obfuscateJS(inputCode.value)
      appStore.showToast('success', '代码混淆完成')
    } else {
      appStore.showToast('error', '暂不支持该语言')
    }
  } catch (e) {
    appStore.showToast('error', '混淆失败: ' + String(e))
  } finally {
    obfuscating.value = false
  }
}

// 复制结果
async function copyResult() {
  if (!outputCode.value) return
  try {
    await navigator.clipboard.writeText(outputCode.value)
    appStore.showToast('success', '已复制到剪贴板')
  } catch {
    appStore.showToast('error', '复制失败')
  }
}

// 下载结果
function downloadResult() {
  if (!outputCode.value) return
  try {
    const ext = language.value === 'javascript' ? 'js' : 'txt'
    const blob = new Blob([outputCode.value], { type: 'text/plain;charset=utf-8' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `obfuscated_${Date.now()}.${ext}`
    a.click()
    URL.revokeObjectURL(url)
    appStore.showToast('success', '下载成功')
  } catch (e) {
    appStore.showToast('error', '下载失败: ' + String(e))
  }
}
</script>
