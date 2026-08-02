<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useSuitesStore } from '@/stores/suites'
import { useI18n } from '@/i18n'

const { t } = useI18n()
const store = useSuitesStore()
const router = useRouter()

const name = ref('')
const description = ref('')
const creating = ref(false)

onMounted(() => store.fetchSuites())

async function create() {
  if (!name.value.trim()) return
  creating.value = true
  await store.createSuite(name.value.trim(), description.value.trim())
  creating.value = false
  name.value = ''
  description.value = ''
}

async function remove(id: string) {
  if (!confirm(t('suites.card.delete.confirm'))) return
  await store.deleteSuite(id)
}

function open(id: string) {
  router.push(`/suites/${id}`)
}

function fmt(v: any): string {
  const d = new Date(v)
  return isNaN(d.getTime()) ? String(v) : d.toLocaleString()
}
</script>

<template>
  <div class="p-6 space-y-6 overflow-y-auto h-full">
    <h2 class="text-sm font-semibold">{{ t('suites.title') }}</h2>

    <div v-if="store.error" class="text-sm border rounded px-3 py-2 border-red-300 text-red-600">
      {{ t('error.prefix') }}{{ store.error }}
      <button class="ml-2 underline" @click="store.fetchSuites">{{ t('error.retry') }}</button>
    </div>

    <div class="border border-border rounded-lg p-4 space-y-3 max-w-lg">
      <h3 class="text-xs font-semibold text-muted-foreground">{{ t('suites.new') }}</h3>
      <input v-model="name" :placeholder="t('suites.new.name')" class="w-full text-sm border border-border rounded px-2 py-1.5" />
      <input v-model="description" :placeholder="t('suites.new.description')" class="w-full text-sm border border-border rounded px-2 py-1.5" />
      <button
        :disabled="creating || !name.trim()"
        class="text-sm px-3 py-1.5 rounded bg-gray-900 text-white disabled:opacity-40"
        @click="create"
      >
        {{ t('suites.new.create') }}
      </button>
    </div>

    <div v-if="store.loading" class="text-sm text-muted-foreground">{{ t('loading') }}</div>
    <div v-else-if="store.suites.length === 0" class="text-sm text-muted-foreground">{{ t('suites.empty') }}</div>

    <div v-else class="space-y-1.5">
      <div
        v-for="s in store.suites"
        :key="s.id"
        class="flex items-center gap-3 border border-border rounded-lg px-4 py-2.5 border-l-4"
        style="border-left-color: #1fb6a7"
      >
        <div class="flex-1 cursor-pointer" @click="open(s.id)">
          <div class="text-sm font-medium">{{ s.name }}</div>
          <div class="text-xs text-muted-foreground">
            <span v-if="s.description">{{ s.description }} · </span>
            {{ fmt(s.updated_at) }}
          </div>
        </div>
        <button class="text-xs px-2 py-1 border border-border rounded hover:bg-gray-50" @click="open(s.id)">
          {{ t('suites.card.open') }}
        </button>
        <button class="text-xs text-red-600 hover:underline" @click="remove(s.id)">{{ t('suites.card.delete') }}</button>
      </div>
    </div>
  </div>
</template>
