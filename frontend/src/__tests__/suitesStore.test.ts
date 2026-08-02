import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useSuitesStore } from '@/stores/suites'
import { ListSuites, GetSuite, ListCases, ListRuns } from '../../wailsjs/go/main/App'
import { db } from '../../wailsjs/go/models'

describe('suites store error handling', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.mocked(ListSuites).mockReset()
    vi.mocked(GetSuite).mockReset()
    vi.mocked(ListCases).mockReset()
    vi.mocked(ListRuns).mockReset()
  })

  it('captures a failed fetchSuites() as store.error and clears loading', async () => {
    vi.mocked(ListSuites).mockRejectedValueOnce(new Error('network down'))
    const store = useSuitesStore()

    await store.fetchSuites()

    expect(store.loading).toBe(false)
    expect(store.error).toContain('network down')
  })

  it('clears the previous error on a successful retry', async () => {
    vi.mocked(ListSuites).mockRejectedValueOnce(new Error('network down'))
    const store = useSuitesStore()
    await store.fetchSuites()
    expect(store.error).not.toBeNull()

    vi.mocked(ListSuites).mockResolvedValueOnce([])
    await store.fetchSuites()

    expect(store.error).toBeNull()
  })

  it('fetchDetail reverses runs to show the most recent first', async () => {
    vi.mocked(GetSuite).mockResolvedValueOnce(
      db.TestSuite.createFrom({ id: 's1', name: 'greeting-bot', description: '' }),
    )
    vi.mocked(ListCases).mockResolvedValueOnce([])
    vi.mocked(ListRuns).mockResolvedValueOnce([
      db.TestRun.createFrom({ id: 'r1', suite_id: 's1', source: 'ci', passed: true, score: 1 }),
      db.TestRun.createFrom({ id: 'r2', suite_id: 's1', source: 'ci', passed: false, score: 0.5 }),
    ])

    const store = useSuitesStore()
    await store.fetchDetail('s1')

    expect(store.detailError).toBeNull()
    expect(store.runs.map((r) => r.id)).toEqual(['r2', 'r1'])
  })
})
