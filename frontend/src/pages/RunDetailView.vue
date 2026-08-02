<script setup lang="ts">
import { computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useSuitesStore } from '@/stores/suites'
import { useI18n } from '@/i18n'

const { t } = useI18n()
const store = useSuitesStore()
const route = useRoute()
const router = useRouter()

const runId = computed(() => route.params.id as string)

function load() {
  store.fetchRunDetail(runId.value)
}

onMounted(load)
watch(runId, load)

function fmt(v: any): string {
  const d = new Date(v)
  return isNaN(d.getTime()) ? String(v) : d.toLocaleString()
}
</script>

<template>
  <div class="p-6 space-y-6 overflow-y-auto h-full">
    <button class="text-xs text-muted-foreground hover:underline" @click="router.back()">
      &larr; {{ t('run.back') }}
    </button>

    <div v-if="store.runDetailError" class="text-sm border rounded px-3 py-2 border-red-300 text-red-600">
      {{ t('error.prefix') }}{{ store.runDetailError }}
      <button class="ml-2 underline" @click="load">{{ t('error.retry') }}</button>
    </div>

    <div v-if="store.runDetailLoading" class="text-sm text-muted-foreground">{{ t('loading') }}</div>

    <template v-else-if="store.runDetail">
      <div>
        <h2 class="text-sm font-semibold">
          {{ t('run.title') }} —
          <span :style="{ color: store.runDetail.run.passed ? '#0ca30c' : '#d03b3b' }">
            {{ store.runDetail.run.passed ? t('run.pass') : t('run.fail') }}
          </span>
        </h2>
        <p class="text-xs text-muted-foreground mt-0.5">
          {{ t('run.score') }} {{ (store.runDetail.run.score * 100).toFixed(0) }}%
          <span v-if="store.runDetail.run.source"> · {{ store.runDetail.run.source }}</span>
          · {{ fmt(store.runDetail.run.started_at) }}
        </p>
      </div>

      <div class="space-y-2 max-w-2xl">
        <div
          v-for="c in store.runDetail.results"
          :key="c.id"
          class="border border-border rounded-lg px-4 py-3 border-l-4"
          :style="{ borderLeftColor: c.passed ? '#0ca30c' : '#d03b3b' }"
        >
          <div class="text-sm font-medium flex items-center gap-2">
            <span>{{ c.passed ? '✓' : '✗' }}</span>
            <span>{{ c.case_name }}</span>
          </div>
          <pre class="text-xs font-mono bg-gray-50 border border-border rounded px-2 py-1.5 mt-2 whitespace-pre-wrap">{{ c.output }}</pre>
          <ul class="mt-2 space-y-1">
            <li v-for="(rr, i) in c.rule_results" :key="i" class="text-xs flex items-start gap-1.5">
              <span :style="{ color: rr.passed ? '#0ca30c' : '#d03b3b' }">{{ rr.passed ? '✓' : '✗' }}</span>
              <span>
                {{ t('rule.' + rr.type) }}<span v-if="rr.value"> "{{ rr.value }}"</span>
                <span v-if="!rr.passed && rr.message" class="text-muted-foreground"> — {{ rr.message }}</span>
              </span>
            </li>
          </ul>
        </div>
      </div>
    </template>
  </div>
</template>
