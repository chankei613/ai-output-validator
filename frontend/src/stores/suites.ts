import { defineStore } from 'pinia'
import {
  ListSuites,
  GetSuite,
  CreateSuite,
  DeleteSuite,
  CreateCase,
  ListCases,
  DeleteCase,
  RunSuite,
  ListRuns,
  GetRun,
} from '../../wailsjs/go/main/App'
import { db, api } from '../../wailsjs/go/models'

export const useSuitesStore = defineStore('suites', {
  state: () => ({
    suites: [] as db.TestSuite[],
    loading: false,
    error: null as string | null,

    detailSuite: null as db.TestSuite | null,
    cases: [] as db.TestCase[],
    runs: [] as db.TestRun[],
    detailLoading: false,
    detailError: null as string | null,

    lastRun: null as api.RunResult | null,
    runError: null as string | null,
    running: false,

    runDetail: null as api.RunResult | null,
    runDetailLoading: false,
    runDetailError: null as string | null,
  }),
  actions: {
    async fetchSuites() {
      this.loading = true
      this.error = null
      try {
        this.suites = (await ListSuites()) ?? []
      } catch (e) {
        this.error = String(e)
      } finally {
        this.loading = false
      }
    },
    async createSuite(name: string, description: string) {
      this.error = null
      try {
        await CreateSuite(name, description)
        await this.fetchSuites()
      } catch (e) {
        this.error = String(e)
      }
    },
    async deleteSuite(id: string) {
      this.error = null
      try {
        await DeleteSuite(id)
        await this.fetchSuites()
      } catch (e) {
        this.error = String(e)
      }
    },
    async fetchDetail(suiteId: string) {
      this.detailLoading = true
      this.detailError = null
      this.lastRun = null
      try {
        const [suite, cases, runs] = await Promise.all([GetSuite(suiteId), ListCases(suiteId), ListRuns(suiteId)])
        this.detailSuite = suite
        this.cases = cases ?? []
        this.runs = (runs ?? []).slice().reverse()
      } catch (e) {
        this.detailError = String(e)
      } finally {
        this.detailLoading = false
      }
    },
    async createCase(suiteId: string, name: string, rules: db.Rule[]) {
      this.detailError = null
      try {
        await CreateCase(suiteId, name, rules)
        await this.fetchDetail(suiteId)
      } catch (e) {
        this.detailError = String(e)
      }
    },
    async deleteCase(suiteId: string, caseId: string) {
      this.detailError = null
      try {
        await DeleteCase(caseId)
        await this.fetchDetail(suiteId)
      } catch (e) {
        this.detailError = String(e)
      }
    },
    async runSuite(suiteId: string, source: string, cases: api.RunCaseInput[]) {
      this.running = true
      this.runError = null
      try {
        this.lastRun = await RunSuite(suiteId, source, cases)
        await this.fetchDetail(suiteId)
      } catch (e) {
        this.runError = String(e)
      } finally {
        this.running = false
      }
    },
    async fetchRunDetail(runId: string) {
      this.runDetailLoading = true
      this.runDetailError = null
      try {
        this.runDetail = await GetRun(runId)
      } catch (e) {
        this.runDetailError = String(e)
      } finally {
        this.runDetailLoading = false
      }
    },
  },
})
