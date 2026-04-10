<template>
  <!-- 二维码批量处理工具 -->
  <div class="page-container">
    <div>
      <div class="page-title"><QrCode :size="20" class="text-primary-400"/>二维码批量处理</div>
      <div class="page-desc">批量生成和解析二维码</div>
    </div>

    <!-- 模式切换 -->
    <div class="tab-bar mb-4">
      <button :class="['tab-item', mode === 'generate' && 'active']" @click="mode = 'generate'">
        <Plus :size="14" class="inline mr-1"/>批量生成
      </button>
      <button :class="['tab-item', mode === 'decode' && 'active']" @click="mode = 'decode'">
        <ScanLine :size="14" class="inline mr-1"/>批量解析
      </button>
    </div>

    <!-- 批量生成模式 -->
    <div v-if="mode === 'generate'">
      <!-- 配置 -->
      <div class="card mb-4">
        <div class="label mb-2">输入内容（每行一个）</div>
        <textarea
          v-model="inputTexts"
          class="textarea-field font-mono text-sm"
          rows="6"
          placeholder="https://example.com&#10;Hello World&#10;每行一个内容..."
          spellcheck="false"
        />
        <div class="flex gap-3 items-end mt-3 flex-wrap">
          <div class="w-28">
            <div class="text-xs opacity-50 mb-1">大小</div>
            <select v-model="qrSize" class="input-field">
              <option :value="128">128px</option>
              <option :value="256">256px</option>
              <option :value="512">512px</option>
            </select>
          </div>
          <div class="w-28">
            <div class="text-xs opacity-50 mb-1">前景色</div>
            <input v-model="fgColor" type="color" class="input-field h-9 p-1 cursor-pointer" />
          </div>
          <div class="w-28">
            <div class="text-xs opacity-50 mb-1">背景色</div>
            <input v-model="bgColor" type="color" class="input-field h-9 p-1 cursor-pointer" />
          </div>
          <div class="w-28">
            <div class="text-xs opacity-50 mb-1">纠错级别</div>
            <select v-model="errorLevel" class="input-field">
              <option value="L">低 (L)</option>
              <option value="M">中 (M)</option>
              <option value="Q">较高 (Q)</option>
              <option value="H">高 (H)</option>
            </select>
          </div>
          <div class="flex-1" />
          <button @click="batchGenerate" class="btn btn-primary" :disabled="!inputTexts.trim()">
            <Wand2 :size="14"/>批量生成
          </button>
          <button @click="downloadAll" class="btn btn-secondary" :disabled="!generatedCodes.length">
            <Download :size="14"/>批量下载
          </button>
        </div>
      </div>

      <!-- 生成结果 -->
      <div class="flex-1 overflow-auto">
        <div class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-3">
          <div v-for="(item, i) in generatedCodes" :key="i" class="card text-center">
            <div class="flex justify-center mb-2">
              <img :src="item.dataUrl" :alt="item.text" class="rounded bg-white p-1" :style="{ maxWidth: '150px' }" />
            </div>
            <div class="text-xs font-mono truncate mb-2" :title="item.text">{{ item.text }}</div>
            <div class="flex gap-1 justify-center">
              <button @click="downloadOne(item, i)" class="btn btn-secondary py-0.5 px-2 text-xs">
                <Download :size="10"/>
              </button>
              <button @click="copyText(item.text)" class="btn btn-secondary py-0.5 px-2 text-xs">
                <Copy :size="10"/>
              </button>
            </div>
          </div>
        </div>
        <div v-if="generatedCodes.length === 0" class="flex items-center justify-center opacity-30 pt-16">
          <div class="text-center">
            <QrCode :size="32" class="mx-auto mb-2"/>
            <div class="text-sm">输入内容后点击批量生成</div>
          </div>
        </div>
      </div>
    </div>

    <!-- 批量解析模式 -->
    <div v-if="mode === 'decode'">
      <div class="card mb-4">
        <div class="label mb-2">上传二维码图片</div>
        <div
          class="border-2 border-dashed rounded-lg p-8 text-center cursor-pointer hover:border-primary-400/50 transition-colors"
          style="border-color:var(--border-color)"
          @click="triggerFileInput"
          @dragover.prevent
          @drop.prevent="handleDrop"
        >
          <Upload :size="32" class="mx-auto mb-2 opacity-40"/>
          <div class="text-sm opacity-50">点击或拖拽上传二维码图片（支持多选）</div>
          <div class="text-xs opacity-30 mt-1">支持 PNG、JPG、GIF 格式</div>
        </div>
        <input ref="fileInputRef" type="file" accept="image/*" multiple class="hidden" @change="handleFileSelect" />
      </div>

      <!-- 解析结果 -->
      <div class="flex-1 overflow-auto space-y-2">
        <div v-for="(item, i) in decodedResults" :key="i" class="card flex items-center gap-3">
          <img :src="item.thumbnail" class="w-12 h-12 rounded object-cover bg-white p-0.5 shrink-0" />
          <div class="flex-1 min-w-0">
            <div class="font-mono text-sm break-all">{{ item.content || '无法解析' }}</div>
            <div class="text-xs opacity-40 mt-1">{{ item.filename }}</div>
          </div>
          <div v-if="item.content" class="flex gap-1 shrink-0">
            <button @click="copyText(item.content)" class="btn btn-secondary py-0.5 px-2 text-xs">
              <Copy :size="11"/>复制
            </button>
          </div>
          <div v-else class="badge badge-error text-xs">解析失败</div>
        </div>
        <div v-if="decodedResults.length === 0" class="flex items-center justify-center opacity-30 pt-16">
          <div class="text-center">
            <ScanLine :size="32" class="mx-auto mb-2"/>
            <div class="text-sm">上传二维码图片进行解析</div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { QrCode, Plus, ScanLine, Wand2, Download, Copy, Upload } from 'lucide-vue-next'
import { useAppStore } from '@/stores/app'

const appStore = useAppStore()
const mode = ref<'generate' | 'decode'>('generate')

// 生成模式相关
const inputTexts = ref('')
const qrSize = ref(256)
const fgColor = ref('#000000')
const bgColor = ref('#ffffff')
const errorLevel = ref('M')
const generatedCodes = ref<{ text: string; dataUrl: string }[]>([])

// 解析模式相关
const fileInputRef = ref<HTMLInputElement | null>(null)
const decodedResults = ref<{ filename: string; content: string; thumbnail: string }[]>([])

let QRCodeLib: any = null
let jsQRLib: any = null

// 初始化库
onMounted(async () => {
  try {
    const mod = await import('qrcode')
    QRCodeLib = mod.default || mod
  } catch {
    appStore.showToast('error', 'qrcode 库加载失败')
  }
  try {
    const mod = await import('jsqr')
    jsQRLib = mod.default || mod
  } catch {
    appStore.showToast('error', 'jsQR 库加载失败')
  }
})

// 批量生成二维码
async function batchGenerate() {
  if (!inputTexts.value.trim()) return
  if (!QRCodeLib) {
    appStore.showToast('error', 'qrcode 库未加载')
    return
  }

  const lines = inputTexts.value.split('\n').map(l => l.trim()).filter(l => l)
  generatedCodes.value = []

  for (const text of lines) {
    try {
      const dataUrl = await QRCodeLib.toDataURL(text, {
        width: qrSize.value,
        margin: 1,
        color: {
          dark: fgColor.value,
          light: bgColor.value,
        },
        errorCorrectionLevel: errorLevel.value,
      })
      generatedCodes.value.push({ text, dataUrl })
    } catch (e) {
      appStore.showToast('error', `生成失败: ${text} - ${String(e)}`)
    }
  }

  if (generatedCodes.value.length > 0) {
    appStore.showToast('success', `已生成 ${generatedCodes.value.length} 个二维码`)
  }
}

// 下载单个二维码
function downloadOne(item: { text: string; dataUrl: string }, index: number) {
  const a = document.createElement('a')
  a.href = item.dataUrl
  a.download = `qrcode_${index + 1}.png`
  a.click()
}

// 批量下载
function downloadAll() {
  generatedCodes.value.forEach((item, i) => {
    setTimeout(() => downloadOne(item, i), i * 200)
  })
  appStore.showToast('success', `正在下载 ${generatedCodes.value.length} 个二维码`)
}

// 触发文件选择
function triggerFileInput() {
  fileInputRef.value?.click()
}

// 处理文件选择
async function handleFileSelect(event: Event) {
  const input = event.target as HTMLInputElement
  if (input.files) {
    await decodeFiles(Array.from(input.files))
  }
  // 重置 input 以允许重复选择同一文件
  input.value = ''
}

// 处理拖拽
async function handleDrop(event: DragEvent) {
  if (event.dataTransfer?.files) {
    await decodeFiles(Array.from(event.dataTransfer.files))
  }
}

// 解析文件
async function decodeFiles(files: File[]) {
  if (!jsQRLib) {
    appStore.showToast('error', 'jsQR 库未加载')
    return
  }

  for (const file of files) {
    if (!file.type.startsWith('image/')) continue

    try {
      const thumbnail = URL.createObjectURL(file)
      const content = await decodeImage(file)
      decodedResults.value.unshift({
        filename: file.name,
        content: content || '',
        thumbnail,
      })
    } catch (e) {
      decodedResults.value.unshift({
        filename: file.name,
        content: '',
        thumbnail: '',
      })
    }
  }

  appStore.showToast('success', `已解析 ${files.length} 个图片`)
}

// 解析单张图片
function decodeImage(file: File): Promise<string | null> {
  return new Promise((resolve) => {
    const reader = new FileReader()
    reader.onload = () => {
      const img = new Image()
      img.onload = () => {
        const canvas = document.createElement('canvas')
        canvas.width = img.width
        canvas.height = img.height
        const ctx = canvas.getContext('2d')
        if (!ctx) { resolve(null); return }
        ctx.drawImage(img, 0, 0)
        const imageData = ctx.getImageData(0, 0, canvas.width, canvas.height)
        const code = jsQRLib(imageData.data, imageData.width, imageData.height)
        resolve(code ? code.data : null)
      }
      img.onerror = () => resolve(null)
      img.src = reader.result as string
    }
    reader.onerror = () => resolve(null)
    reader.readAsDataURL(file)
  })
}

// 复制文本
async function copyText(text: string) {
  try {
    await navigator.clipboard.writeText(text)
    appStore.showToast('success', '已复制到剪贴板')
  } catch {
    appStore.showToast('error', '复制失败')
  }
}
</script>
