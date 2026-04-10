<template>
  <ToolPage title="Hosts 编辑" description="系统 Hosts 文件编辑">
<div class="flex gap-2 mb-3">
      <button @click="loadHosts" class="btn btn-secondary"><RefreshCw :size="14"/>刷新</button>
      <button @click="saveHosts" class="btn btn-primary"><Save :size="14"/>保存</button>
      <div class="flex-1"/>
      <span class="text-xs opacity-40 self-center">保存需要管理员权限，将自动创建 .bak 备份</span>
    </div>
    <textarea v-model="content" class="textarea-field flex-1 min-h-0" spellcheck="false"
      placeholder="加载中..."/>
  </ToolPage>
</template>
<script setup lang="ts">
import ToolPage from '@/components/ToolPage.vue'
import { ref, onMounted } from 'vue'
import { Server, RefreshCw, Save } from 'lucide-vue-next'
import { useAppStore } from '@/stores/app'
import { GetHostsContent, SaveHostsContent } from '../../../wailsjs/go/network/NetworkTools'
const appStore = useAppStore()
const content = ref('')
async function loadHosts() {
  try { content.value = await GetHostsContent() as string }
  catch(e) { appStore.showToast('error', '读取失败：' + String(e)) }
}
async function saveHosts() {
  try {
    await SaveHostsContent(content.value)
    appStore.showToast('success', 'Hosts 文件已保存')
  } catch(e) { appStore.showToast('error', '保存失败：' + String(e)) }
}
onMounted(loadHosts)
</script>
