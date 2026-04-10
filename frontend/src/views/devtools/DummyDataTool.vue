<template>
  <ToolPage title="占位文本生成" description="批量生成测试数据">
<!-- 占位文本/测试数据生成工具 -->
  
    

    <!-- 配置区域 -->
    <div class="card mb-4">
      <div class="flex gap-3 items-end flex-wrap">
        <!-- 数据类型 -->
        <div class="flex-1 min-w-[160px]">
          <div class="label mb-1">数据类型</div>
          <select v-model="dataType" class="input-field">
            <option v-for="t in dataTypes" :key="t.id" :value="t.id">{{ t.label }}</option>
          </select>
        </div>
        <!-- 数量 -->
        <div class="w-28">
          <div class="label mb-1">数量</div>
          <input v-model.number="count" type="number" min="1" max="100" class="input-field" />
        </div>
        <!-- 语言 -->
        <div class="w-28">
          <div class="label mb-1">语言</div>
          <select v-model="lang" class="input-field">
            <option value="zh">中文</option>
            <option value="en">英文</option>
          </select>
        </div>
        <!-- 操作按钮 -->
        <button @click="generate" class="btn btn-primary" :disabled="generating">
          <Wand2 :size="14"/>{{ generating ? '生成中...' : '生成' }}
        </button>
        <button @click="copyAll" class="btn btn-secondary" :disabled="!results.length">
          <Copy :size="14"/>复制全部
        </button>
        <button @click="exportData" class="btn btn-secondary" :disabled="!results.length">
          <Download :size="14"/>导出
        </button>
      </div>
    </div>

    <!-- 生成结果 -->
    <div class="flex-1 overflow-auto">
      <div class="space-y-2">
        <div v-for="(item, i) in results" :key="i" class="card flex items-center gap-3">
          <span class="text-xs opacity-40 w-8 shrink-0 text-right">{{ i + 1 }}</span>
          <div class="flex-1 font-mono text-sm break-all">{{ item }}</div>
          <button @click="copyOne(item)" class="btn btn-secondary py-0.5 px-2 text-xs shrink-0">
            <Copy :size="11"/>
          </button>
        </div>
        <div v-if="results.length === 0" class="flex items-center justify-center opacity-30 pt-16">
          <div class="text-center">
            <Database :size="32" class="mx-auto mb-2"/>
            <div class="text-sm">选择数据类型后点击生成</div>
          </div>
        </div>
      </div>
    </div>
  </ToolPage>
</template>

<script setup lang="ts">
import ToolPage from '@/components/ToolPage.vue'
import { ref } from 'vue'
import { Database, Wand2, Copy, Download } from 'lucide-vue-next'
import { useAppStore } from '@/stores/app'

const appStore = useAppStore()
const dataType = ref('name')
const count = ref(10)
const lang = ref('zh')
const generating = ref(false)
const results = ref<string[]>([])

// 数据类型列表
const dataTypes = [
  { id: 'name', label: '姓名' },
  { id: 'phone', label: '手机号' },
  { id: 'email', label: '邮箱' },
  { id: 'address', label: '地址' },
  { id: 'company', label: '公司' },
  { id: 'ip', label: 'IP 地址' },
  { id: 'mac', label: 'MAC 地址' },
  { id: 'date', label: '日期' },
  { id: 'uuid', label: 'UUID' },
  { id: 'sentence', label: '句子' },
  { id: 'paragraph', label: '段落' },
]

// 中文数据池
const zhSurnames = ['张', '李', '王', '刘', '陈', '杨', '黄', '赵', '周', '吴', '徐', '孙', '马', '朱', '胡', '郭', '何', '林', '罗', '高']
const zhNames = ['伟', '芳', '娜', '敏', '静', '丽', '强', '磊', '军', '洋', '勇', '艳', '杰', '娟', '涛', '明', '超', '秀英', '华', '慧']
const zhCities = ['北京市', '上海市', '广州市', '深圳市', '杭州市', '成都市', '武汉市', '南京市', '重庆市', '西安市', '苏州市', '长沙市', '天津市', '郑州市', '东莞市']
const zhStreets = ['中山路', '人民路', '解放路', '建设路', '和平路', '文化路', '科技路', '创新路', '幸福路', '长安街']
const zhCompanies = ['科技有限公司', '信息技术有限公司', '网络科技有限公司', '数据服务有限公司', '智能科技有限公司', '软件开发有限公司', '云计算有限公司', '电子商务有限公司']
const zhCompanyPrefix = ['星辰', '蓝海', '云端', '极光', '飞鸟', '龙腾', '鹏程', '博远', '创新', '智慧']
const zhSentences = [
  '今天天气真好，适合出去散步。',
  '学习使人进步，知识改变命运。',
  '技术改变世界，创新引领未来。',
  '生活不止眼前的苟且，还有诗和远方。',
  '千里之行，始于足下。',
  '书山有路勤为径，学海无涯苦作舟。',
  '不积跬步，无以至千里。',
  '天行健，君子以自强不息。',
]
const zhParagraphs = [
  '在这个信息爆炸的时代，数据已经成为最重要的资源之一。如何高效地处理和分析数据，成为了每个企业都需要面对的挑战。通过合理的数据管理策略，企业可以更好地理解市场趋势，优化业务流程，提升决策质量。',
  '人工智能技术的快速发展正在深刻改变着我们的生活方式。从智能语音助手到自动驾驶汽车，从医疗诊断到金融分析，AI的应用场景越来越广泛。然而，技术的发展也带来了新的挑战，如数据隐私、算法偏见等问题需要我们认真对待。',
  '云计算作为现代IT基础设施的核心，为企业提供了弹性的计算资源和灵活的服务模式。通过云平台，企业可以快速部署应用，降低运维成本，提高业务敏捷性。随着容器化和微服务架构的普及，云计算的价值将进一步凸显。',
]

// 英文数据池
const enFirstNames = ['James', 'Mary', 'John', 'Patricia', 'Robert', 'Jennifer', 'Michael', 'Linda', 'David', 'Elizabeth', 'William', 'Barbara', 'Richard', 'Susan', 'Joseph', 'Jessica', 'Thomas', 'Sarah', 'Christopher', 'Karen']
const enLastNames = ['Smith', 'Johnson', 'Williams', 'Brown', 'Jones', 'Garcia', 'Miller', 'Davis', 'Rodriguez', 'Martinez', 'Hernandez', 'Lopez', 'Gonzalez', 'Wilson', 'Anderson', 'Thomas', 'Taylor', 'Moore', 'Jackson', 'Martin']
const enCities = ['New York', 'Los Angeles', 'Chicago', 'Houston', 'Phoenix', 'Philadelphia', 'San Antonio', 'San Diego', 'Dallas', 'San Jose']
const enStreets = ['Main St', 'Oak Ave', 'Elm St', 'Park Blvd', 'Cedar Ln', 'Maple Dr', 'Pine St', 'Washington Ave', 'Lincoln Rd', 'Jefferson Way']
const enCompanies = ['Tech Corp', 'Data Solutions Inc', 'Cloud Systems Ltd', 'Digital Innovations', 'Smart Software Co', 'Web Services Inc', 'AI Technologies', 'Cyber Dynamics']
const enSentences = [
  'The quick brown fox jumps over the lazy dog.',
  'Technology is best when it brings people together.',
  'Innovation distinguishes between a leader and a follower.',
  'The only way to do great work is to love what you do.',
  'Stay hungry, stay foolish.',
  'Code is like humor. When you have to explain it, it is bad.',
  'First, solve the problem. Then, write the code.',
  'Experience is the name everyone gives to their mistakes.',
]
const enParagraphs = [
  'The rapid advancement of technology has transformed the way businesses operate in the modern world. From cloud computing to artificial intelligence, organizations are leveraging cutting-edge tools to streamline operations, enhance productivity, and deliver better customer experiences.',
  'Data security remains one of the most critical challenges facing organizations today. With the increasing frequency and sophistication of cyber attacks, companies must invest in robust security measures, employee training, and incident response plans to protect their valuable digital assets.',
  'The future of work is evolving rapidly, driven by technological innovation and changing workforce expectations. Remote work, flexible schedules, and digital collaboration tools have become essential components of the modern workplace, requiring organizations to adapt their management strategies accordingly.',
]

// 随机工具函数
function rand<T>(arr: T[]): T {
  return arr[Math.floor(Math.random() * arr.length)]
}

function randInt(min: number, max: number): number {
  return Math.floor(Math.random() * (max - min + 1)) + min
}

function randIP(): string {
  return `${randInt(1, 255)}.${randInt(0, 255)}.${randInt(0, 255)}.${randInt(1, 254)}`
}

function randMAC(): string {
  return Array.from({ length: 6 }, () =>
    randInt(0, 255).toString(16).padStart(2, '0').toUpperCase()
  ).join(':')
}

function randUUID(): string {
  return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, c => {
    const r = (Math.random() * 16) | 0
    const v = c === 'x' ? r : (r & 0x3) | 0x8
    return v.toString(16)
  })
}

function randDate(): string {
  const y = randInt(2020, 2026)
  const m = String(randInt(1, 12)).padStart(2, '0')
  const d = String(randInt(1, 28)).padStart(2, '0')
  return `${y}-${m}-${d}`
}

function randPhone(): string {
  const prefixes = ['130', '131', '132', '133', '134', '135', '136', '137', '138', '139',
    '150', '151', '152', '153', '155', '156', '157', '158', '159',
    '170', '171', '172', '173', '175', '176', '177', '178',
    '180', '181', '182', '183', '184', '185', '186', '187', '188', '189']
  return rand(prefixes) + String(randInt(10000000, 99999999))
}

function randEmail(name: string): string {
  const domains = ['gmail.com', 'outlook.com', 'qq.com', '163.com', 'yahoo.com', 'hotmail.com']
  const cleanName = name.toLowerCase().replace(/\s/g, '.')
  return `${cleanName}${randInt(1, 999)}@${rand(domains)}`
}

// 根据类型生成单条数据
function generateOne(): string {
  const isZh = lang.value === 'zh'

  switch (dataType.value) {
    case 'name':
      return isZh ? rand(zhSurnames) + rand(zhNames) : rand(enFirstNames) + ' ' + rand(enLastNames)
    case 'phone':
      return isZh ? randPhone() : `+1-${randInt(200, 999)}-${randInt(100, 999)}-${randInt(1000, 9999)}`
    case 'email': {
      const name = isZh ? rand(zhSurnames) + rand(zhNames) : rand(enFirstNames) + rand(enLastNames)
      return randEmail(name)
    }
    case 'address':
      return isZh
        ? `${rand(zhCities)}${rand(zhStreets)}${randInt(1, 200)}号`
        : `${randInt(100, 9999)} ${rand(enStreets)}, ${rand(enCities)}`
    case 'company':
      return isZh
        ? rand(zhCompanyPrefix) + rand(zhCompanies)
        : rand(enCompanies)
    case 'ip':
      return randIP()
    case 'mac':
      return randMAC()
    case 'date':
      return randDate()
    case 'uuid':
      return randUUID()
    case 'sentence':
      return isZh ? rand(zhSentences) : rand(enSentences)
    case 'paragraph':
      return isZh ? rand(zhParagraphs) : rand(enParagraphs)
    default:
      return ''
  }
}

// 批量生成
function generate() {
  generating.value = true
  try {
    const n = Math.max(1, Math.min(100, count.value))
    results.value = Array.from({ length: n }, () => generateOne())
    appStore.showToast('success', `已生成 ${n} 条${dataTypes.find(t => t.id === dataType.value)?.label || ''}数据`)
  } catch (e) {
    appStore.showToast('error', '生成失败: ' + String(e))
  } finally {
    generating.value = false
  }
}

// 复制单条
async function copyOne(val: string) {
  try {
    await navigator.clipboard.writeText(val)
    appStore.showToast('success', '已复制到剪贴板')
  } catch {
    appStore.showToast('error', '复制失败')
  }
}

// 复制全部
async function copyAll() {
  if (!results.value.length) return
  try {
    await navigator.clipboard.writeText(results.value.join('\n'))
    appStore.showToast('success', `已复制 ${results.value.length} 条数据`)
  } catch {
    appStore.showToast('error', '复制失败')
  }
}

// 导出为文本文件
function exportData() {
  if (!results.value.length) return
  try {
    const content = results.value.join('\n')
    const blob = new Blob([content], { type: 'text/plain;charset=utf-8' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `${dataType.value}_data_${Date.now()}.txt`
    a.click()
    URL.revokeObjectURL(url)
    appStore.showToast('success', '导出成功')
  } catch (e) {
    appStore.showToast('error', '导出失败: ' + String(e))
  }
}
</script>
