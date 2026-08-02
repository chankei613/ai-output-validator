<script setup lang="ts">
import { ref, reactive, computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useSuitesStore } from '@/stores/suites'
import { useI18n } from '@/i18n'
import { db } from '../../wailsjs/go/models'

const { t } = useI18n()
const store = useSuitesStore()
const route = useRoute()
const router = useRouter()

const suiteId = computed(() => route.params.id as string)

const ruleTypes = ['contains', 'not_contains', 'regex', 'min_length', 'max_length', 'json_valid', 'json_key_exists']

const newCaseName = ref('')
const newCaseRules = ref<db.Rule[]>([{ type: 'contains', value: '' } as db.Rule])
const creatingCase = ref(false)

const runSource = ref('manual')
const outputs = reactive<Record<string, string>>({})

function load() {
  store.fetchDetail(suiteId.value)
}

onMounted(load)
watch(suiteId, load)

function addRuleRow() {
  newCaseRules.value.push({ type: 'contains', value: '' } as db.Rule)
}

function removeRuleRow(i: number) {
  newCaseRules.value.splice(i, 1)
}

function needsValue(type: string) {
  return type !== 'json_valid'
}

async function createCase() {
  if (!newCaseName.value.trim() || newCaseRules.value.length === 0) return
  creatingCase.value = true
  await store.createCase(suiteId.value, newCaseName.value.trim(), newCaseRules.value)
  creatingCase.value = false
  newCaseName.value = ''
  newCaseRules.value = [{ type: 'contains', value: '' } as db.Rule]
}

async function removeCase(caseId: string) {
  if (!confirm(t('detail.cases.delete'))) return
  await store.deleteCase(suiteId.value, caseId)
}

async function runSuite() {
  const cases = store.cases.map((c) => ({ case_id: c.id, output: outputs[c.id] ?? '' }))
  await store.runSuite(suiteId.value, runSource.value.trim(), cases)
}

function openRun(id: string) {
  router.push(`/runs/${id}`)
}

function fmt(v: any): string {
  const d = new Date(v)
  return isNaN(d.getTime()) ? String(v) : d.toLocaleString()
}
</script>

<template>
  <div class="p-6 space-y-6 overflow-y-auto h-full">
    <button class="text-xs text-muted-foreground hover:underline" @click="router.push('/suites')">
      &larr; {{ t('detail.back') }}
    </button>

    <div v-if="store.detailError" class="text-sm border rounded px-3 py-2 border-red-300 text-red-600">
      {{ t('error.prefix') }}{{ store.detailError }}
      <button class="ml-2 underline" @click="load">{{ t('error.retry') }}</button>
    </div>

    <div v-if="store.detailLoading" class="text-sm text-muted-foreground">{{ t('loading') }}</div>

    <template v-else-if="store.detailSuite">
      <div>
        <h2 class="text-sm font-semibold">{{ store.detailSuite.name }}</h2>
        <p v-if="store.detailSuite.description" class="text-xs text-muted-foreground mt-0.5">
          {{ store.detailSuite.description }}
        </p>
      </div>

      <section class="space-y-2">
        <h3 class="text-xs font-semibold text-muted-foreground">{{ t('detail.cases') }}</h3>
        <div v-if="store.cases.length === 0" class="text-xs text-muted-foreground">{{ t('detail.cases.empty') }}</div>
        <div v-else class="space-y-1.5">
          <div v-for="c in store.cases" :key="c.id" class="border border-border rounded-lg px-4 py-2.5 space-y-1">
            <div class="flex items-center justify-between">
              <span class="text-sm font-medium">{{ c.name }}</span>
              <button class="text-xs text-red-600 hover:underline" @click="removeCase(c.id)">{{ t('detail.cases.delete') }}</button>
            </div>
            <div class="text-xs text-muted-foreground">
              <span v-for="(r, i) in c.rules" :key="i" class="mr-3">
                {{ t('rule.' + r.type) }}<span v-if="r.value"> "{{ r.value }}"</span>
              </span>
            </div>
          </div>
        </div>
      </section>

      <section class="border border-border rounded-lg p-4 space-y-3 max-w-2xl">
        <h3 class="text-xs font-semibold text-muted-foreground">{{ t('detail.cases.new') }}</h3>
        <input v-model="newCaseName" :placeholder="t('detail.cases.new.name')" class="w-full text-sm border border-border rounded px-2 py-1.5" />

        <div v-for="(rule, i) in newCaseRules" :key="i" class="flex gap-2 items-center">
          <select v-model="rule.type" class="text-sm border border-border rounded px-2 py-1.5">
            <option v-for="rt in ruleTypes" :key="rt" :value="rt">{{ t('rule.' + rt) }}</option>
          </select>
          <input
            v-if="needsValue(rule.type)"
            v-model="rule.value"
            :placeholder="t('detail.cases.rule.value')"
            class="flex-1 text-sm border border-border rounded px-2 py-1.5"
          />
          <button v-if="newCaseRules.length > 1" class="text-xs text-red-600 hover:underline shrink-0" @click="removeRuleRow(i)">
            {{ t('detail.cases.rule.remove') }}
          </button>
        </div>
        <button class="text-xs text-muted-foreground hover:underline" @click="addRuleRow">{{ t('detail.cases.rule.add') }}</button>

        <div>
          <button
            :disabled="creatingCase || !newCaseName.trim()"
            class="text-sm px-3 py-1.5 rounded bg-gray-900 text-white disabled:opacity-40"
            @click="createCase"
          >
            {{ t('detail.cases.new.add') }}
          </button>
        </div>
      </section>

      <section class="border border-border rounded-lg p-4 space-y-3 max-w-2xl">
        <h3 class="text-xs font-semibold text-muted-foreground">{{ t('detail.run') }}</h3>
        <div v-if="store.cases.length === 0" class="text-xs text-muted-foreground">{{ t('detail.run.empty') }}</div>
        <template v-else>
          <div v-for="c in store.cases" :key="c.id" class="space-y-1">
            <label class="text-xs text-muted-foreground">{{ t('detail.run.output', { name: c.name }) }}</label>
            <textarea v-model="outputs[c.id]" rows="2" class="w-full text-sm border border-border rounded px-2 py-1.5 font-mono" />
          </div>
          <input v-model="runSource" :placeholder="t('detail.run.source')" class="w-full text-sm border border-border rounded px-2 py-1.5" />
          <div v-if="store.runError" class="text-sm border rounded px-3 py-2 border-red-300 text-red-600">
            {{ t('error.prefix') }}{{ store.runError }}
          </div>
          <button
            :disabled="store.running"
            class="text-sm px-3 py-1.5 rounded text-white disabled:opacity-40"
            style="background: #1fb6a7"
            @click="runSuite"
          >
            {{ t('detail.run.button') }}
          </button>
        </template>
      </section>

      <section class="space-y-2">
        <h3 class="text-xs font-semibold text-muted-foreground">{{ t('detail.history') }}</h3>
        <div v-if="store.runs.length === 0" class="text-xs text-muted-foreground">{{ t('detail.history.empty') }}</div>
        <div v-else class="space-y-1.5">
          <div
            v-for="r in store.runs"
            :key="r.id"
            class="flex items-center gap-3 border border-border rounded-lg px-4 py-2.5 border-l-4"
            :style="{ borderLeftColor: r.passed ? '#0ca30c' : '#d03b3b' }"
          >
            <div class="flex-1">
              <div class="text-sm font-medium">
                {{ r.passed ? t('run.pass') : t('run.fail') }} · {{ t('run.score') }} {{ (r.score * 100).toFixed(0) }}%
              </div>
              <div class="text-xs text-muted-foreground">
                <span v-if="r.source">{{ r.source }} · </span>{{ fmt(r.started_at) }}
              </div>
            </div>
            <button class="text-xs px-2 py-1 border border-border rounded hover:bg-gray-50" @click="openRun(r.id)">
              {{ t('detail.history.open') }}
            </button>
          </div>
        </div>
      </section>
    </template>
  </div>
</template>
