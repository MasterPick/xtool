<template>
  <div class="page-container">
    <div class="page-title"><FilePen :size="20" class="text-primary-400"/>批量重命名</div>
    <div class="page-desc">批量重命名文件，支持前缀/后缀/序号/正则替换</div>

    <!-- 重命名模式选择 -->
    <div class="card mb-4">
      <div class="label mb-2">重命名模式</div>
      <div class="flex gap-2 flex-wrap">
        <button :class="['btn', renameMode==='regex'?'btn-primary':'btn-secondary']" @click="renameMode='regex'">正则替换</button>
        <button :class="['btn', renameMode==='prefix'?'btn-primary':'btn-secondary']" @click="renameMode='prefix'">添加前缀</button>
        <button :class="['btn', renameMode==='suffix'?'btn-primary':'btn-secondary']" @click="renameMode='suffix'">添加后缀</button>
        <button :class="['btn', renameMode==='sequence'?'btn-primary':'btn-secondary']" @click="renameMode='sequence'">序号命名</button>
        <button :class="['btn', renameMode==='case'?'btn-primary':'btn-secondary']" @click="renameMode='case'">大小写转换</button>
      </div>
    </div>

    <!-- 模式参数 -->
    <div class="card mb-4">
      <!-- 正则替换 -->
      <template v-if="renameMode==='regex'">
        <div class="flex gap-2 mb-3">
          <input v-model="pattern" class="input-field flex-1" placeholder="查找（支持正则）..."/>
          <input v-model="replacement" class="input-field flex-1" placeholder="替换为（$1 引用捕获组）..."/>
        </div>
      </template>
      <!-- 前缀 -->
      <template v-if="renameMode==='prefix'">
        <div class="flex gap-2 mb-3">
          <input v-model="prefixText" class="input-field flex-1" placeholder="输入前缀文本..."/>
        </div>
      </template>
      <!-- 后缀 -->
      <template v-if="renameMode==='suffix'">
        <div class="flex gap-2 mb-3">
          <input v-model="suffixText" class="input-field flex-1" placeholder="输入后缀文本（在扩展名之前）..."/>
        </div>
      </template>
      <!-- 序号 -->
      <template v-if="renameMode==='sequence'">
        <div class="flex gap-2 mb-3">
          <input v-model="seqBaseName" class="input-field flex-1" placeholder="基础名称（如 photo）..."/>
          <input v-model.number="seqStart" class="input-field w-24" placeholder="起始序号" type="number" min="0"/>
          <input v-model="seqSeparator" class="input-field w-16" placeholder="分隔符" value="_"/>
        </div>
      </template>
      <!-- 大小写 -->
      <template v-if="renameMode==='case'">
        <div class="flex gap-2 mb-3">
          <button :class="['btn', caseType==='upper'?'btn-primary':'btn-secondary']" @click="caseType='upper'">全部大写</button>
          <button :class="['btn', caseType==='lower'?'btn-primary':'btn-secondary']" @click="caseType='lower'">全部小写</button>
          <button :class="['btn', caseType==='title'?'btn-primary':'btn-secondary']" @click="caseType='title'">首字母大写</button>
        </div>
      </template>

      <div class="flex gap-2">
        <button @click="preview" class="btn btn-secondary">预览</button>
        <button @click="doRename" class="btn btn-primary">执行重命名</button>
        <button @click="clearFiles" class="btn btn-secondary">清空列表</button>
      </div>
    </div>

    <!-- 文件拖放区 -->
    <div class="card mb-4">
      <div class="label mb-2">拖放文件到此处（或点击选择）</div>
      <div class="drop-zone" @dragover.prevent @drop="onDrop">
        <Upload :size="32" class="mx-auto mb-2 opacity-40"/>
        <div class="text-sm opacity-40">拖放文件到这里</div>
      </div>
    </div>

    <!-- 文件列表预览 -->
    <div v-if="files.length" class="card">
      <div class="label mb-2">文件列表 ({{ files.length }} 个文件)</div>
      <div class="space-y-1 max-h-64 overflow-auto">
        <div v-for="(f,i) in files" :key="i" class="flex gap-3 text-xs py-1">
          <span class="opacity-60 flex-1 truncate">{{ f.original }}</span>
          <span class="text-primary-400">-></span>
          <span class="flex-1 truncate" :class="f.renamed !== f.original ? 'text-green-400' : ''">{{ f.renamed || f.original }}</span>
        </div>
      </div>
    </div>
  </div>
</template>
<script setup lang="ts">
import { ref } from 'vue'
import { FilePen, Upload } from 'lucide-vue-next'
import { useAppStore } from '@/stores/app'
import { BatchRename } from '../../../wailsjs/go/sysinfo/SysInfo'
const appStore = useAppStore()

const files = ref<{original:string,renamed:string,path:string}[]>([])
const renameMode = ref<'regex'|'prefix'|'suffix'|'sequence'|'case'>('regex')
const pattern = ref(''), replacement = ref('')
const prefixText = ref(''), suffixText = ref('')
const seqBaseName = ref(''), seqStart = ref(1), seqSeparator = ref('_')
const caseType = ref<'upper'|'lower'|'title'>('upper')

function onDrop(e: DragEvent) {
  const items = e.dataTransfer?.files
  if (!items) return
  for (let i = 0; i < items.length; i++) {
    const f = items[i] as any
    files.value.push({ original: f.name, renamed: f.name, path: f.path || '' })
  }
}

function getFileNameWithoutExt(name: string): { base: string; ext: string } {
  const lastDot = name.lastIndexOf('.')
  if (lastDot === -1) return { base: name, ext: '' }
  return { base: name.substring(0, lastDot), ext: name.substring(lastDot) }
}

function preview() {
  files.value = files.value.map((f, idx) => {
    const { base, ext } = getFileNameWithoutExt(f.original)
    let newBase = base
    try {
      switch (renameMode.value) {
        case 'regex':
          if (pattern.value) {
            const re = new RegExp(pattern.value, 'g')
            newBase = base.replace(re, replacement.value)
          }
          break
        case 'prefix':
          newBase = prefixText.value + base
          break
        case 'suffix':
          newBase = base + suffixText.value
          break
        case 'sequence':
          newBase = seqBaseName.value + seqSeparator.value + String(seqStart.value + idx).padStart(3, '0')
          break
        case 'case':
          switch (caseType.value) {
            case 'upper': newBase = base.toUpperCase(); break
            case 'lower': newBase = base.toLowerCase(); break
            case 'title': newBase = base.charAt(0).toUpperCase() + base.slice(1).toLowerCase(); break
          }
          break
      }
    } catch {
      // 正则错误时保持原名
    }
    return { ...f, renamed: newBase + ext }
  })
}

async function doRename() {
  // 先确保预览是最新的
  preview()

  // 检查是否有实际变更
  const changedFiles = files.value.filter(f => f.original !== f.renamed)
  if (changedFiles.length === 0) {
    appStore.showToast('warning', '没有需要重命名的文件')
    return
  }

  // 构建重命名映射
  const renameMap: Record<string, string> = {}
  for (const f of changedFiles) {
    if (f.path) {
      const lastSep = f.path.lastIndexOf(/[/\\]/.test(f.path) ? (f.path.includes('\\') ? '\\' : '/') : '/')
      const dir = lastSep >= 0 ? f.path.substring(0, lastSep + 1) : ''
      renameMap[f.path] = dir + f.renamed
    } else {
      renameMap[f.original] = f.renamed
    }
  }

  try {
    const result = await BatchRename(renameMap) as any
    if (result && result.success !== false) {
      appStore.showToast('success', `成功重命名 ${changedFiles.length} 个文件`)
      // 更新文件列表
      files.value = files.value.map(f => ({ ...f, original: f.renamed }))
    } else {
      appStore.showToast('error', result?.error || '重命名失败')
    }
  } catch (e) {
    appStore.showToast('error', '重命名失败: ' + String(e))
  }
}

function clearFiles() {
  files.value = []
}
</script>
<style scoped>
.drop-zone {
  border: 2px dashed var(--border-color);
  border-radius: 10px; padding: 32px;
  text-align: center; cursor: pointer;
  transition: border-color 0.2s;
}
.drop-zone:hover { border-color: var(--accent); }
</style>
