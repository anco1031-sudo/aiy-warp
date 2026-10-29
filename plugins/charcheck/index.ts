/**
 * aiy.charcheck — Pre-Send Quality Gate ข้อ 1 แบบอัตโนมัติ
 *
 * ที่มา: 29-Oct-26 อักษรต่างภาษาหลุดกลางประโยคไทย 4 ครั้งในวันเดียว
 *        (และหลุดซ้ำระหว่างที่กำลังแก้เรื่องนี้อยู่ด้วย)
 *        → ความตั้งใจดีไม่พอ ต้องมีด่านที่ไม่มีวันเหนื่อย
 *
 * ทำงาน: หลัง tool `write` / `edit` / `patch` สำเร็จ → รัน
 *        `python3 ~/.config/opencode/scripts/check-chars.py <ไฟล์>`
 *        เจออักษรต่างภาษา → แนบคำเตือนเข้า `result.content` (ส่วนที่ model เห็นจริง)
 *
 * ทำไม "เตือน" เป็นค่าเริ่มต้น ไม่ใช่ "บล็อก":
 *   false positive มีจริง (บทความภาษาจีน, ชื่อเรื่องต่างภาษา) → บล็อกแข็งทำให้ agent
 *   ติดลูปแก้ทั้งที่ถูกแล้ว · เตือน = model เห็นและตัดสินใจเองได้
 *   อยากเข้มกว่านี้ → options { "mode": "block" }
 *
 * รูปแบบ options ใน ~/.config/opencode/opencode.jsonc:
 *   "plugins": [{ "package": "./plugins/charcheck", "options": { "mode": "warn" } }]
 * ปล่อยว่าง = auto-discovered จาก ~/.config/opencode/plugins/ (โหมด warn)
 */
import { Plugin } from "@opencode/plugin"
import { spawnSync } from "node:child_process"
import { existsSync, statSync } from "node:fs"
import { homedir } from "node:os"
import { extname } from "node:path"

const DEFAULT_SCRIPT = `${homedir()}/.config/opencode/scripts/check-chars.py`
const DEFAULT_EXTS = [".md", ".txt", ".py", ".sh", ".yaml", ".yml", ".json", ".csv", ".html"]
const WRITE_TOOLS = new Set(["write", "edit", "patch", "multiedit"])

type Options = {
  mode?: "warn" | "block" | "off"
  script?: string
  exts?: string[]
  allowScripts?: string[]
  timeoutMs?: number
}

type Hit = { summary: string; details: string[] }

function pathFromInput(input: unknown): string | undefined {
  if (!input || typeof input !== "object") return undefined
  const record = input as Record<string, unknown>
  const candidate = record.path ?? record.filePath ?? record.file ?? record.filename
  return typeof candidate === "string" ? candidate : undefined
}

function runCheck(script: string, file: string, allowScripts: string[], timeout: number): Hit | undefined {
  const args = [script, file, "--max-per-file", "5"]
  for (const name of allowScripts) args.push("--allow-script", name)

  const proc = spawnSync("python3", args, { encoding: "utf8", timeout })
  if (proc.error || proc.status !== 1) return undefined // 0 = สะอาด · อื่น ๆ = สคริปต์พัง → เงียบ

  const lines = (proc.stdout ?? "").split("\n")
  const summary = lines.find((line) => line.includes("พบ") && line.includes("ตัว")) ?? "พบอักษรต่างภาษาตกค้าง"
  const details = lines
    .filter((line) => line.trim().startsWith("บรรทัด"))
    .map((line) => line.trim())
    .slice(0, 5)

  return { summary: summary.replace(/^❌\s*/, "").trim(), details }
}

function buildWarning(file: string, hit: Hit): string {
  const bar = "──────────────────────────────────────────────"
  return [
    "",
    bar,
    "🚨 [aiy.charcheck] อักษรต่างภาษาตกค้างในไฟล์ที่เพิ่งเขียน",
    `   ไฟล์: ${file}`,
    `   ${hit.summary}`,
    ...(hit.details.length ? hit.details.map((line) => `   ${line}`) : []),
    "",
    "   ⛔ กฎ (AGENTS.md · Pre-Send Quality Gate ข้อ 1): งานภาษาไทยต้องไม่มีอักษร",
    "      จีน / ญี่ปุ่น / เกาหลี / รัสเซีย / อาหรับ ตกค้าง",
    '   ✅ แก้: เขียนทับไฟล์เดิมด้วยคำไทย หรือถ้าตั้งใจอ้างอิง → เขียนเป็น code point เช่น "U+7BC4"',
    bar,
  ].join("\n")
}

export default Plugin.define({
  id: "aiy.charcheck",
  async setup(ctx) {
    const options = (ctx.options ?? {}) as Options
    const mode = options.mode ?? "warn"
    if (mode === "off") return

    const script = options.script ?? DEFAULT_SCRIPT
    const exts = new Set((options.exts ?? DEFAULT_EXTS).map((ext) => ext.toLowerCase()))
    const allowScripts = options.allowScripts ?? ["กรีก (Greek)"]
    const timeout = options.timeoutMs ?? 10_000

    if (!existsSync(script)) {
      console.warn(`[aiy.charcheck] ไม่พบสคริปต์ ${script} — ข้ามการตรวจ`)
      return
    }

    const registration = await ctx.tool.hook("execute.after", (event) => {
      if (event.status !== "completed") return
      if (!WRITE_TOOLS.has(event.tool)) return

      const file = pathFromInput(event.input)
      if (!file) return
      if (!exts.has(extname(file).toLowerCase())) return

      try {
        if (!existsSync(file) || !statSync(file).isFile()) return
      } catch {
        return
      }

      const hit = runCheck(script, file, allowScripts, timeout)
      if (!hit) return

      if (mode === "block") {
        throw new Error(
          `อักษรต่างภาษาตกค้าง: ${hit.summary} — แก้ไฟล์ ${file} แล้วเขียนใหม่ (ถ้าตั้งใจอ้างอิง ให้เขียนเป็น code point เช่น U+7BC4)`,
        )
      }

      const warning = buildWarning(file, hit)
      const content = event.result.content

      // `output` ของ tool write เป็น object โครงสร้าง (operation/target/...) — ที่ model
      // เห็นจริงคือ `content` text parts → ต้องแนบที่นั่น
      if (Array.isArray(content)) {
        const last = content[content.length - 1]
        if (last && last.type === "text" && typeof last.text === "string") last.text = `${last.text}${warning}`
        else content.push({ type: "text", text: warning })
      } else {
        console.warn(warning)
      }

      event.result.metadata = { ...(event.result.metadata ?? {}), charcheck: "dirty" }
    })

    return () => registration.dispose()
  },
})
