<template>
  <!-- 图片工具页面：格式转换、压缩、尺寸调整 -->
  <div class="page-container">
    <!-- 页面标题 -->
    <div>
      <div class="page-title">
        <Image :size="20" class="text-primary-400" />
        图片工具
      </div>
      <div class="page-desc">格式转换 · 图片压缩 · 尺寸调整 · 批量处理</div>
    </div>

    <!-- 操作工具栏 -->
    <div class="toolbar mb-4">
      <button @click="selectFiles" class="btn btn-primary">
        <Upload :size="14" />
        选择图片
      </button>
      <button v-if="images.length" @click="clearAll" class="btn btn-secondary">
        <Trash2 :size="14" />
        清空
      </button>
      <div class="flex-1" />
      <!-- 格式选择 -->
      <select v-model="outputFormat" class="input-field w-28">
        <option value="png">PNG</option>
        <option value="jpeg">JPG</option>
        <option value="webp">WebP</option>
      </select>
      <!-- 压缩质量 -->
      <div class="flex items-center gap-2">
        <span class="text-xs opacity-50">质量</span>
        <input type="range" v-model="quality" min="10" max="100" step="5" class="w-20" />
        <span class="text-xs w-8">{{ quality }}%</span>
      </div>
      <!-- 尺寸 -->
      <div class="flex items-center gap-2">
        <span class="text-xs opacity-50">缩放</span>
        <input type="number" v-model="scale" min="10" max="200" class="input-field w-16 text-center" />
        <span class="text-xs">%</span>
      </div>
    </div>

    <!-- 空状态 -->
    <div v-if="!images.length" class="flex-1 flex items-center justify-center">
      <div
        class="drop-zone flex flex-col items-center justify-center gap-4 cursor-pointer"
        :class="{ 'drag-over': isDragging }"
        @click="selectFiles"
        @dragover.prevent="isDragging = true"
        @dragleave="isDragging = false"
        @drop.prevent="handleDrop"
      >
        <ImagePlus :size="48" class="opacity-30" />
        <div class="text-center">
          <div class="text-sm opacity-60">拖拽图片到此处</div>
          <div class="text-xs opacity-40 mt-1">或点击选择文件</div>
        </div>
        <div class="text-xs opacity-30">支持 PNG / JPG / WebP / GIF / BMP</div>
      </div>
    </div>

    <!-- 图片列表 -->
    <div v-else class="flex-1 flex flex-col gap-4 min-h-0">
      <!-- 统计信息 -->
      <div class="flex items-center gap-6 text-xs opacity-60">
        <span>共 {{ images.length }} 张图片</span>
        <span>原始大小: {{ formatSize(totalOriginalSize) }}</span>
        <span v-if="totalCompressedSize">压缩后: {{ formatSize(totalCompressedSize) }}</span>
        <span v-if="savedPercent" class="text-green-400">节省 {{ savedPercent }}%</span>
      </div>

      <!-- 图片网格 -->
      <div class="flex-1 overflow-auto">
        <div class="grid grid-cols-3 gap-3">
          <div
            v-for="(img, idx) in images"
            :key="idx"
            class="image-card card"
          >
            <!-- 预览图 -->
            <div class="relative aspect-video bg-black/20 rounded-lg overflow-hidden mb-2">
              <img :src="img.preview" class="w-full h-full object-contain" />
              <!-- 状态标签 -->
              <span
                v-if="img.status"
                :class="['status-badge', img.status === 'done' ? 'status-success' : img.status === 'error' ? 'status-error' : 'status-pending']"
              >
                {{ img.status === 'done' ? '完成' : img.status === 'error' ? '失败' : '处理中' }}
              </span>
            </div>
            <!-- 文件信息 -->
            <div class="flex items-center justify-between text-xs">
              <span class="truncate flex-1 opacity-70">{{ img.name }}</span>
              <button @click="removeImage(idx)" class="p-1 hover:text-red-400">
                <X :size="12" />
              </button>
            </div>
            <div class="flex items-center gap-3 text-xs opacity-50 mt-1">
              <span>{{ img.width }}×{{ img.height }}</span>
              <span>{{ formatSize(img.size) }}</span>
              <span v-if="img.newSize" class="text-green-400">→ {{ formatSize(img.newSize) }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- 底部操作 -->
      <div class="flex items-center gap-3 pt-2 border-t border-white/5">
        <button @click="processAll" :disabled="isProcessing" class="btn btn-primary flex-1">
          <Loader2 v-if="isProcessing" :size="14" class="loading-spin" />
          <Wand2 v-else :size="14" />
          {{ isProcessing ? '处理中...' : '开始处理' }}
        </button>
        <button v-if="hasResult" @click="downloadAll" class="btn btn-secondary">
          <Download :size="14" />
          全部下载
        </button>
      </div>
    </div>

    <!-- 隐藏的文件输入 -->
    <input ref="fileInput" type="file" accept="image/*" multiple class="hidden" @change="handleFiles" />
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { Image, Upload, Trash2, ImagePlus, Wand2, Download, X, Loader2 } from 'lucide-vue-next'
import { useAppStore } from '@/stores/app'

const appStore = useAppStore()

// 文件输入引用
const fileInput = ref<HTMLInputElement | null>(null)

// 图片列表
interface ImageItem {
  file: File
  name: string
  preview: string
  width: number
  height: number
  size: number
  newSize?: number
  blob?: Blob
  status?: 'pending' | 'processing' | 'done' | 'error'
}

const images = ref<ImageItem[]>([])

// 设置
const outputFormat = ref<'png' | 'jpeg' | 'webp'>('webp')
const quality = ref(80)
const scale = ref(100)

// 状态
const isDragging = ref(false)
const isProcessing = ref(false)

// 统计
const totalOriginalSize = computed(() => images.value.reduce((sum, img) => sum + img.size, 0))
const totalCompressedSize = computed(() => images.value.reduce((sum, img) => sum + (img.newSize || 0), 0))
const savedPercent = computed(() => {
  if (!totalCompressedSize.value) return 0
  const saved = 1 - totalCompressedSize.value / totalOriginalSize.value
  return Math.round(saved * 100)
})
const hasResult = computed(() => images.value.some(img => img.status === 'done'))

// 选择文件
function selectFiles() {
  fileInput.value?.click()
}

// 处理文件选择
function handleFiles(e: Event) {
  const files = (e.target as HTMLInputElement).files
  if (!files) return
  addFiles(Array.from(files))
}

// 处理拖放
function handleDrop(e: DragEvent) {
  isDragging.value = false
  const files = e.dataTransfer?.files
  if (!files) return
  const imageFiles = Array.from(files).filter(f => f.type.startsWith('image/'))
  addFiles(imageFiles)
}

// 添加文件到列表
async function addFiles(files: File[]) {
  for (const file of files) {
    try {
      const preview = URL.createObjectURL(file)
      const { width, height } = await getImageSize(preview)
      images.value.push({
        file,
        name: file.name,
        preview,
        width,
        height,
        size: file.size,
        status: 'pending'
      })
    } catch (err) {
      console.error('加载图片失败:', err)
    }
  }
}

// 获取图片尺寸
function getImageSize(src: string): Promise<{ width: number; height: number }> {
  return new Promise((resolve, reject) => {
    const img = new window.Image()
    img.onload = () => resolve({ width: img.width, height: img.height })
    img.onerror = reject
    img.src = src
  })
}

// 移除图片
function removeImage(idx: number) {
  const img = images.value[idx]
  if (img.preview) URL.revokeObjectURL(img.preview)
  images.value.splice(idx, 1)
}

// 清空所有
function clearAll() {
  images.value.forEach(img => URL.revokeObjectURL(img.preview))
  images.value = []
}

// 处理所有图片
async function processAll() {
  if (isProcessing.value) return
  isProcessing.value = true

  for (const img of images.value) {
    if (img.status === 'done') continue
    img.status = 'processing'
    try {
      const blob = await compressImage(img.file)
      img.blob = blob
      img.newSize = blob.size
      img.status = 'done'
    } catch (err) {
      img.status = 'error'
      console.error('压缩失败:', err)
    }
  }

  isProcessing.value = false
  appStore.showToast('success', `处理完成，共节省 ${savedPercent.value}% 空间`)
}

// 压缩单张图片
function compressImage(file: File): Promise<Blob> {
  return new Promise((resolve, reject) => {
    const img = new window.Image()
    img.onload = () => {
      // 计算新尺寸
      const newWidth = Math.round(img.width * scale.value / 100)
      const newHeight = Math.round(img.height * scale.value / 100)

      // 创建 canvas
      const canvas = document.createElement('canvas')
      canvas.width = newWidth
      canvas.height = newHeight
      const ctx = canvas.getContext('2d')!

      // 绘制
      ctx.drawImage(img, 0, 0, newWidth, newHeight)

      // 导出
      const mimeType = `image/${outputFormat.value}`
      canvas.toBlob(
        blob => blob ? resolve(blob) : reject(new Error('导出失败')),
        mimeType,
        quality.value / 100
      )
    }
    img.onerror = reject
    img.src = URL.createObjectURL(file)
  })
}

// 下载全部
async function downloadAll() {
  for (const img of images.value) {
    if (img.status !== 'done' || !img.blob) continue
    const link = document.createElement('a')
    const ext = outputFormat.value === 'jpeg' ? 'jpg' : outputFormat.value
    link.download = img.name.replace(/\.[^.]+$/, `.${ext}`)
    link.href = URL.createObjectURL(img.blob)
    link.click()
    URL.revokeObjectURL(link.href)
    // 间隔避免浏览器阻止
    await new Promise(r => setTimeout(r, 200))
  }
  appStore.showToast('success', '下载完成')
}

// 格式化文件大小
function formatSize(bytes: number): string {
  if (bytes < 1024) return bytes + ' B'
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB'
  return (bytes / 1024 / 1024).toFixed(2) + ' MB'
}
</script>

<style scoped>
/* 拖放区域 */
.drop-zone {
  width: 100%;
  height: 100%;
  min-height: 300px;
  border: 2px dashed rgba(255, 255, 255, 0.15);
  border-radius: 16px;
  transition: all 0.2s ease;
}

.drop-zone:hover,
.drop-zone.drag-over {
  border-color: var(--accent);
  background: rgba(99, 102, 241, 0.05);
}

/* 图片卡片 */
.image-card {
  padding: 10px;
}

/* 状态标签 */
.status-badge {
  position: absolute;
  top: 6px;
  right: 6px;
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 10px;
  font-weight: 600;
}

.status-pending { background: rgba(255, 255, 255, 0.1); }
.status-success { background: rgba(16, 185, 129, 0.2); color: #10b981; }
.status-error { background: rgba(239, 68, 68, 0.2); color: #ef4444; }

/* 质量滑块 */
input[type="range"] {
  -webkit-appearance: none;
  height: 4px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 2px;
  cursor: pointer;
}

input[type="range"]::-webkit-slider-thumb {
  -webkit-appearance: none;
  width: 14px;
  height: 14px;
  border-radius: 50%;
  background: var(--accent);
  cursor: pointer;
}
</style>
