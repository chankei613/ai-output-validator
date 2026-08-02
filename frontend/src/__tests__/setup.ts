import { vi } from 'vitest'

// Node 25+ の実験的 localStorage がhappy-domのwindow.localStorageと衝突し
// `--localstorage-file` 未指定だと getItem 等がthrowする既知の非互換。
// メモリ上の簡易実装に差し替えて回避する。
class MemoryStorage implements Storage {
  private store = new Map<string, string>()
  get length() { return this.store.size }
  clear() { this.store.clear() }
  getItem(key: string) { return this.store.has(key) ? this.store.get(key)! : null }
  key(index: number) { return Array.from(this.store.keys())[index] ?? null }
  removeItem(key: string) { this.store.delete(key) }
  setItem(key: string, value: string) { this.store.set(key, String(value)) }
}
const memoryStorage = new MemoryStorage()
Object.defineProperty(globalThis, 'localStorage', { value: memoryStorage, configurable: true })
Object.defineProperty(window, 'localStorage', { value: memoryStorage, configurable: true })

// Wails ランタイムのモック — テスト環境では no-op
vi.mock('../../wailsjs/runtime/runtime', () => ({
  EventsOn: vi.fn(),
  EventsOff: vi.fn(),
  EventsEmit: vi.fn(),
  BrowserOpenURL: vi.fn(),
}))

vi.mock('../../wailsjs/go/main/App', () => ({
  GetAppVersion: vi.fn().mockResolvedValue('0.1.0'),
  GetAPIURL: vi.fn().mockResolvedValue('http://127.0.0.1:8426'),
  Quit: vi.fn().mockResolvedValue(undefined),
  ListSuites: vi.fn().mockResolvedValue([]),
  GetSuite: vi.fn(),
  CreateSuite: vi.fn(),
  DeleteSuite: vi.fn().mockResolvedValue(undefined),
  CreateCase: vi.fn(),
  ListCases: vi.fn().mockResolvedValue([]),
  DeleteCase: vi.fn().mockResolvedValue(undefined),
  RunSuite: vi.fn(),
  ListRuns: vi.fn().mockResolvedValue([]),
  GetRun: vi.fn(),
  ListKeys: vi.fn().mockResolvedValue([]),
  IssueKey: vi.fn(),
  RevokeKey: vi.fn().mockResolvedValue(undefined),
}))
