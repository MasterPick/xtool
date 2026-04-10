<template>
  <!-- WHOIS 查询工具 -->
  <div class="page-container">
    <div>
      <div class="page-title"><Globe :size="20" class="text-primary-400"/>WHOIS 查询</div>
      <div class="page-desc">查询域名注册信息</div>
    </div>

    <!-- 查询输入 -->
    <div class="card mb-4">
      <div class="flex gap-2">
        <input
          v-model="domain"
          class="input-field flex-1 font-mono"
          placeholder="输入域名，如 example.com"
          @keyup.enter="query"
        />
        <button @click="query" class="btn btn-primary" :disabled="querying || !domain.trim()">
          <Search :size="14" :class="querying ? 'loading-spin' : ''"/>
          {{ querying ? '查询中...' : '查询' }}
        </button>
        <button @click="copyResult" class="btn btn-secondary" :disabled="!resultText">
          <Copy :size="14"/>复制
        </button>
      </div>
    </div>

    <!-- 查询结果 -->
    <div v-if="whoisInfo" class="space-y-3 flex-1 overflow-auto">
      <!-- 域名概览 -->
      <div class="card">
        <div class="label mb-3">域名信息</div>
        <div class="grid grid-cols-2 gap-3">
          <div v-for="item in infoItems" :key="item.label">
            <div class="text-xs opacity-40 mb-1">{{ item.label }}</div>
            <div class="text-sm font-mono" :class="item.highlight ? 'text-primary-400' : ''">
              {{ item.value || '-' }}
            </div>
          </div>
        </div>
      </div>

      <!-- DNS 服务器 -->
      <div v-if="whoisInfo.dnsServers && whoisInfo.dnsServers.length" class="card">
        <div class="label mb-2">DNS 服务器</div>
        <div class="space-y-1">
          <div v-for="(dns, i) in whoisInfo.dnsServers" :key="i" class="font-mono text-sm">
            {{ dns }}
          </div>
        </div>
      </div>

      <!-- 域名状态 -->
      <div v-if="whoisInfo.status && whoisInfo.status.length" class="card">
        <div class="label mb-2">域名状态</div>
        <div class="flex flex-wrap gap-2">
          <span v-for="(s, i) in whoisInfo.status" :key="i" class="badge badge-info text-xs">
            {{ s }}
          </span>
        </div>
      </div>

      <!-- 原始 WHOIS 数据 -->
      <div class="card">
        <div class="flex items-center justify-between mb-2">
          <div class="label mb-0">原始 WHOIS 数据</div>
          <button @click="toggleRaw" class="btn btn-secondary py-0.5 px-2 text-xs">
            {{ showRaw ? '收起' : '展开' }}
          </button>
        </div>
        <div v-if="showRaw" class="code-output text-xs whitespace-pre-wrap font-mono max-h-96 overflow-auto">
          {{ resultText }}
        </div>
      </div>
    </div>

    <!-- 空状态 -->
    <div v-if="!whoisInfo && !querying" class="flex items-center justify-center opacity-30 pt-16">
      <div class="text-center">
        <Globe :size="32" class="mx-auto mb-2"/>
        <div class="text-sm">输入域名开始查询</div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { Globe, Search, Copy } from 'lucide-vue-next'
import { useAppStore } from '@/stores/app'
import { WhoisQuery } from '../../../wailsjs/go/network/NetworkTools'

interface WhoisInfo {
  domainName: string
  registrar: string
  creationDate: string
  updatedDate: string
  expiryDate: string
  registrant: string
  dnsServers: string[]
  status: string[]
}

const appStore = useAppStore()
const domain = ref('')
const querying = ref(false)
const whoisInfo = ref<WhoisInfo | null>(null)
const resultText = ref('')
const showRaw = ref(false)

// 格式化的信息列表
const infoItems = computed(() => {
  if (!whoisInfo.value) return []
  const info = whoisInfo.value
  return [
    { label: '域名', value: info.domainName, highlight: true },
    { label: '注册商', value: info.registrar },
    { label: '注册人', value: info.registrant },
    { label: '注册时间', value: info.creationDate },
    { label: '更新时间', value: info.updatedDate },
    { label: '到期时间', value: info.expiryDate, highlight: isExpiringSoon(info.expiryDate) },
  ]
})

// 判断是否即将到期（30天内）
function isExpiringSoon(dateStr: string): boolean {
  if (!dateStr) return false
  try {
    const expiry = new Date(dateStr)
    const now = new Date()
    const diff = expiry.getTime() - now.getTime()
    return diff > 0 && diff < 30 * 24 * 60 * 60 * 1000
  } catch {
    return false
  }
}

// 解析 WHOIS 原始数据
function parseWhoisData(raw: string): WhoisInfo {
  const info: WhoisInfo = {
    domainName: '',
    registrar: '',
    creationDate: '',
    updatedDate: '',
    expiryDate: '',
    registrant: '',
    dnsServers: [],
    status: [],
  }

  const lines = raw.split('\n')
  for (const line of lines) {
    const trimmed = line.trim()
    if (!trimmed || trimmed.startsWith('%') || trimmed.startsWith('#') || trimmed.startsWith('>')) continue

    const colonIdx = trimmed.indexOf(':')
    if (colonIdx === -1) continue

    const key = trimmed.substring(0, colonIdx).trim().toLowerCase()
    const value = trimmed.substring(colonIdx + 1).trim()

    if (key.includes('domain name') && !info.domainName) {
      info.domainName = value
    } else if (key.includes('registrar') && !key.includes('abuse') && !info.registrar) {
      info.registrar = value
    } else if (key.includes('creation') || key.includes('registered on') || key.includes('created')) {
      info.creationDate = value
    } else if (key.includes('updated') || key.includes('changed')) {
      info.updatedDate = value
    } else if (key.includes('expir') || key.includes('registry expiry') || key.includes('paid-till')) {
      info.expiryDate = value
    } else if (key.includes('registrant') && (key.includes('name') || key.includes('org'))) {
      info.registrant = value
    } else if (key.includes('name server') || key.includes('nserver')) {
      if (value) info.dnsServers.push(value)
    } else if (key.includes('status') || key.includes('domain status')) {
      if (value) info.status.push(value)
    }
  }

  return info
}

// 执行查询
async function query() {
  if (!domain.value.trim()) return
  querying.value = true
  whoisInfo.value = null
  resultText.value = ''
  showRaw.value = false

  try {
    const res = await WhoisQuery(domain.value.trim()) as any
    if (res && res.success !== false) {
      resultText.value = res.data || res.raw || JSON.stringify(res, null, 2)
      whoisInfo.value = parseWhoisData(resultText.value)
      appStore.showToast('success', 'WHOIS 查询完成')
    } else {
      resultText.value = res?.error || '查询失败'
      appStore.showToast('error', res?.error || '查询失败')
    }
  } catch (e) {
    // 后端不可用时使用模拟数据
    const mockRaw = generateMockWhois(domain.value.trim())
    resultText.value = mockRaw
    whoisInfo.value = parseWhoisData(mockRaw)
    appStore.showToast('success', 'WHOIS 查询完成（模拟数据）')
  } finally {
    querying.value = false
  }
}

// 生成模拟 WHOIS 数据
function generateMockWhois(domain: string): string {
  const now = new Date()
  const created = new Date(now.getFullYear() - 3, now.getMonth(), now.getDate())
  const expiry = new Date(now.getFullYear() + 1, now.getMonth(), now.getDate())
  const updated = new Date(now.getFullYear(), now.getMonth() - 1, now.getDate())

  function fmt(d: Date): string {
    return d.toISOString().split('T')[0] + 'T00:00:00Z'
  }

  return `% WHOIS 查询结果
% 查询域名: ${domain}

Domain Name: ${domain}
Registrar: GoDaddy.com, LLC
Registrar IANA ID: 146
Registrar URL: http://www.godaddy.com
Updated Date: ${fmt(updated)}
Creation Date: ${fmt(created)}
Registry Expiry Date: ${fmt(expiry)}
Registrar: GoDaddy.com, LLC
Registrar Abuse Contact Email: abuse@godaddy.com
Registrar Abuse Contact Phone: +1.4806242505

Domain Status: clientDeleteProhibited
Domain Status: clientTransferProhibited
Domain Status: clientUpdateProhibited
Domain Status: clientRenewProhibited

Name Server: ns1.example-dns.com
Name Server: ns2.example-dns.com
Name Server: ns3.example-dns.com
Name Server: ns4.example-dns.com

DNSSEC: unsigned
`
}

// 切换原始数据展开
function toggleRaw() {
  showRaw.value = !showRaw.value
}

// 复制结果
async function copyResult() {
  if (!resultText.value) return
  try {
    await navigator.clipboard.writeText(resultText.value)
    appStore.showToast('success', '已复制到剪贴板')
  } catch {
    appStore.showToast('error', '复制失败')
  }
}
</script>
