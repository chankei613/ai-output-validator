import { ref } from 'vue'

export type Locale = 'en' | 'ja'

// localStorage に保存して再起動後も維持する
const saved = localStorage.getItem('locale') as Locale | null
const locale = ref<Locale>(saved === 'en' || saved === 'ja' ? saved : 'ja')

const messages: Record<Locale, Record<string, string>> = {
  en: {
    'app.subtitle': 'AI Output Validator',
    'lang.toggle': 'JA',
    'nav.suites': 'Suites',
    'nav.help': 'Help',
    'nav.settings': 'Settings',

    'error.prefix': 'Error: ',
    'error.retry': 'Retry',
    'loading': 'Loading…',

    'suites.title': 'Test Suites',
    'suites.empty': 'No suites yet. Create one to start defining acceptance criteria.',
    'suites.new': 'New suite',
    'suites.new.name': 'Name',
    'suites.new.description': 'Description',
    'suites.new.create': 'Create',
    'suites.card.open': 'Open',
    'suites.card.delete': 'Delete',
    'suites.card.delete.confirm': 'Delete this suite, its cases, and all run history? This cannot be undone.',

    'detail.back': 'Back to suites',
    'detail.cases': 'Cases',
    'detail.cases.empty': 'No cases yet. Add one below.',
    'detail.cases.new': 'New case',
    'detail.cases.new.name': 'Case name',
    'detail.cases.new.add': 'Create case',
    'detail.cases.delete': 'Delete',
    'detail.cases.rule.type': 'Rule type',
    'detail.cases.rule.value': 'Value',
    'detail.cases.rule.add': '+ Add rule',
    'detail.cases.rule.remove': 'Remove',

    'detail.run': 'Run against this suite',
    'detail.run.empty': 'Add at least one case before running.',
    'detail.run.output': 'Output for "{name}"',
    'detail.run.source': 'Source (optional, e.g. "manual" or "ci:github-actions")',
    'detail.run.button': 'Run validation',

    'detail.history': 'Run history',
    'detail.history.empty': 'No runs yet.',
    'detail.history.open': 'View',

    'run.title': 'Run detail',
    'run.back': 'Back to suite',
    'run.pass': 'PASS',
    'run.fail': 'FAIL',
    'run.score': 'score',

    'rule.contains': 'contains',
    'rule.not_contains': 'does not contain',
    'rule.regex': 'matches regex',
    'rule.min_length': 'min length',
    'rule.max_length': 'max length',
    'rule.json_valid': 'is valid JSON',
    'rule.json_key_exists': 'JSON key exists',

    'help.title': 'Help',
    'help.intro': 'How suites, cases, rules, and runs fit together.',
    'help.what.title': 'What this app does',
    'help.what.body': 'This is a unit-test runner for AI-generated output. A Suite groups related Cases; each Case has one or more deterministic Rules (contains, regex, length, JSON checks). Submit an AI output per case and every rule is checked — no AI provider is called by this app itself, so it works with whatever generated the output, from any system.',
    'help.start.title': 'Getting started',
    'help.start.1': 'Create a suite, then add one or more cases to it.',
    'help.start.2': 'Give each case one or more rules — all rules in a case must pass for the case to pass.',
    'help.start.3': 'Paste an AI output per case and click Run to see pass/fail and a per-rule breakdown.',
    'help.start.4': 'Every run is kept in history so you can track score over time.',
    'help.stuck.title': 'Common snags',
    'help.stuck.1': 'Want to call this from CI? Use the `ovcli` command bundled in the repository — it reads a JSON file of case outputs, POSTs it to this app\'s API, and exits 1 on failure.',
    'help.stuck.2': 'A run fails immediately if it references a case_id that doesn\'t exist — nothing is silently skipped.',
    'help.stuck.3': 'This app never calls an AI provider itself. Something else must generate the output and hand it to this app for grading.',

    'settings.title': 'Settings',
    'settings.api.title': 'API endpoint',
    'settings.api.desc': 'External systems (CI, scripts) can submit outputs and read results here.',
    'settings.keys.title': 'API keys',
    'settings.keys.name': 'Key name',
    'settings.keys.issue': 'Issue key',
    'settings.keys.issued': 'Key issued — copy it now, it will not be shown again',
    'settings.keys.copy': 'Copy',
    'settings.keys.revoke': 'Revoke',
    'settings.keys.revoked': 'Revoked',
    'settings.keys.empty': 'No keys issued yet.',
    'settings.version': 'Version',
    'settings.quit': 'Quit',
    'settings.quit.confirm': 'Quit the app?',
  },
  ja: {
    'app.subtitle': 'AI Output Validator',
    'lang.toggle': 'EN',
    'nav.suites': 'Suite',
    'nav.help': 'ヘルプ',
    'nav.settings': '設定',

    'error.prefix': 'エラー: ',
    'error.retry': '再試行',
    'loading': '読み込み中…',

    'suites.title': 'Test Suite',
    'suites.empty': 'まだSuiteがありません。作成して受け入れ条件の定義を始めましょう。',
    'suites.new': '新規作成',
    'suites.new.name': '名前',
    'suites.new.description': '説明',
    'suites.new.create': '作成',
    'suites.card.open': '開く',
    'suites.card.delete': '削除',
    'suites.card.delete.confirm': 'このSuiteとケース・実行履歴を全て削除しますか？元に戻せません。',

    'detail.back': 'Suite一覧へ戻る',
    'detail.cases': 'ケース',
    'detail.cases.empty': 'まだケースがありません。下から追加してください。',
    'detail.cases.new': '新しいケース',
    'detail.cases.new.name': 'ケース名',
    'detail.cases.new.add': 'ケースを作成',
    'detail.cases.delete': '削除',
    'detail.cases.rule.type': 'ルール種別',
    'detail.cases.rule.value': '値',
    'detail.cases.rule.add': '+ ルールを追加',
    'detail.cases.rule.remove': '削除',

    'detail.run': 'このSuiteに対して実行',
    'detail.run.empty': '実行するにはまずケースを追加してください。',
    'detail.run.output': '「{name}」への出力',
    'detail.run.source': '実行元（任意、例: "manual" や "ci:github-actions"）',
    'detail.run.button': '検証を実行',

    'detail.history': '実行履歴',
    'detail.history.empty': 'まだ実行履歴がありません。',
    'detail.history.open': '詳細',

    'run.title': '実行詳細',
    'run.back': 'Suiteへ戻る',
    'run.pass': '合格',
    'run.fail': '不合格',
    'run.score': 'スコア',

    'rule.contains': '含む',
    'rule.not_contains': '含まない',
    'rule.regex': '正規表現に一致',
    'rule.min_length': '最小文字数',
    'rule.max_length': '最大文字数',
    'rule.json_valid': '妥当なJSONである',
    'rule.json_key_exists': 'JSONキーが存在する',

    'help.title': 'ヘルプ',
    'help.intro': 'Suite・ケース・ルール・実行がどう連動するかをまとめました。',
    'help.what.title': 'このアプリでできること',
    'help.what.body': 'AI生成物のためのUnit Testランナーです。Suiteは関連するケースをまとめたもの、各ケースは1つ以上の決定論的なルール（contains・正規表現・文字数・JSONチェック）を持ちます。ケースごとにAIの出力を投入すると全ルールが検証されます。本アプリ自身はAIプロバイダーを一切呼び出さないため、どのシステムが生成した出力でも扱えます。',
    'help.start.title': 'はじめに',
    'help.start.1': 'Suiteを作成し、1つ以上のケースを追加します。',
    'help.start.2': '各ケースに1つ以上のルールを設定します — ケース内の全ルールに合格して初めてケース合格です。',
    'help.start.3': 'ケースごとにAIの出力を貼り付けて「実行」を押すと、合格/不合格とルールごとの内訳が表示されます。',
    'help.start.4': '実行はすべて履歴に残るため、スコアの推移を追跡できます。',
    'help.stuck.title': 'よくある詰まりどころ',
    'help.stuck.1': 'CIから呼びたい → リポジトリ同梱の`ovcli`コマンドを使ってください。ケース出力のJSONファイルを読み込み、本アプリのAPIにPOSTし、不合格ならexit 1で終了します。',
    'help.stuck.2': '存在しないcase_idを指定すると実行全体がエラーになります — 黙ってスキップされることはありません。',
    'help.stuck.3': '本アプリはAIプロバイダーを一切呼び出しません。出力を生成する側は別のシステムで、本アプリはその採点だけを担当します。',

    'settings.title': '設定',
    'settings.api.title': 'APIエンドポイント',
    'settings.api.desc': '外部システム（CI・スクリプト等）はここで出力の投入・結果の取得ができます。',
    'settings.keys.title': 'APIキー',
    'settings.keys.name': 'キー名',
    'settings.keys.issue': 'キーを発行',
    'settings.keys.issued': 'キーを発行しました — この場では二度と表示されないので今すぐコピーしてください',
    'settings.keys.copy': 'コピー',
    'settings.keys.revoke': '失効',
    'settings.keys.revoked': '失効済み',
    'settings.keys.empty': 'まだキーがありません。',
    'settings.version': 'バージョン',
    'settings.quit': '終了',
    'settings.quit.confirm': 'アプリを終了しますか？',
  },
}

export function useI18n() {
  function t(key: string, params?: Record<string, string | number>): string {
    let msg = messages[locale.value][key] ?? messages.en[key] ?? key
    if (params) {
      for (const [k, v] of Object.entries(params)) {
        msg = msg.replace(`{${k}}`, String(v))
      }
    }
    return msg
  }

  function toggleLocale() {
    locale.value = locale.value === 'en' ? 'ja' : 'en'
    localStorage.setItem('locale', locale.value)
  }

  return { t, locale, toggleLocale }
}
