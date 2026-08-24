#!/usr/bin/env bun
/**
 * model-check.ts — Model Health Check (command-based, run on demand)
 *
 * เปรียบเทียบ agent config (model ที่ประกาศใน frontmatter) กับ registry
 * (`opencode models --verbose`) + model-mapping.json (intended allocation)
 *
 * Output:
 *   - Report ฉบับเต็ม  -> /tmp/opencode/model-check.md
 *   - สรุปสั้นๆ          -> stdout (ให้ agent/command เอาไป relay ต่อได้เลย)
 *
 * Status:
 *   🔴 RED     model หายจาก registry → เสนอ fallback จาก mapping
 *   🟢 GREEN  advisory (เช่น vision role ยังไม่ได้ใช้ mimo)
 *   ⚠️ NO_MODEL ไม่มี model ทั้งใน agent file และ mapping
 *   ✅ OK      ตรงกันหมด
 *   🆕 NEW     model ใน registry ที่ยังไม่มี agent ใช้ (พร้อม capabilities)
 *
 * Usage: bun ~/.config/opencode/scripts/model-check.ts
 */
import { $ } from "bun";
import { join } from "path";

const scriptsDir = import.meta.dir;
const configDir = join(scriptsDir, "..");
const agentsDir = join(configDir, "agents");
const modelMappingPath = join(configDir, "model-mapping.json");
const reportPath = "/tmp/opencode/model-check.md";

interface AgentInfo {
  name: string;
  model: string | null; // ประกาศใน agent file (null = ไม่มี)
  source: "file" | "mapping";
  role: string;
}

interface Mapping {
  agents: Record<string, { role: string; current: string; fallbacks: string[] }>;
  capability_hints?: {
    vision_roles?: string[];
    long_output_roles?: string[];
    long_context_roles?: string[];
  };
}

interface CheckResult {
  agent: string;
  role: string;
  model: string;
  source: "file" | "mapping";
  status: "RED" | "GREEN" | "OK" | "NO_MODEL";
  suggestion: string;
}

interface RegistryModel {
  id: string;
  name?: string;
  providerID?: string;
  limit?: { context?: number; input?: number; output?: number };
  capabilities?: {
    reasoning?: boolean;
    input?: { text?: boolean; audio?: boolean; image?: boolean; video?: boolean; pdf?: boolean };
    output?: { text?: boolean; audio?: boolean; image?: boolean; video?: boolean; pdf?: boolean };
  };
}

async function readAgentModels(): Promise<AgentInfo[]> {
  const agents: AgentInfo[] = [];
  try {
    const files = (await $`ls ${agentsDir}/*.md`.nothrow().text())
      .trim()
      .split("\n")
      .filter(Boolean);
    for (const file of files) {
      const name = file.split("/").pop()?.replace(/\.md$/, "");
      if (!name) continue;
      const content = await $`cat ${file}`.nothrow().text();
      const match = content.match(/^model:\s*(.+)$/m);
      agents.push({
        name,
        model: match ? match[1].trim() : null,
        source: "file",
        role: "",
      });
    }
  } catch {
    // agents dir อ่านไม่ได้ → คืนค่าว่าง
  }
  return agents;
}

/** Parse `opencode models --verbose` — ข้าม header line (ชื่อ model) แล้วเก็บ JSON object */
async function getModelRegistry(): Promise<Map<string, RegistryModel>> {
  const registry = new Map<string, RegistryModel>();
  try {
    const output = await $`opencode models --verbose`.nothrow().text();
    const lines = output.trim().split("\n");
    let buffer = "";
    let braceCount = 0;
    for (const line of lines) {
      // header line (เช่น "opencode/big-pickle") ข้ามไป — JSON เริ่มที่ "{"
      if (!buffer && !line.trim().startsWith("{")) continue;
      buffer += line + "\n";
      for (const char of line) {
        if (char === "{") braceCount++;
        else if (char === "}") braceCount--;
      }
      if (braceCount === 0 && buffer.trim()) {
        try {
          const model = JSON.parse(buffer.trim()) as RegistryModel;
          if (model.id) registry.set(model.id, model);
        } catch {
          // ข้าม JSON ที่ parse ไม่ได้
        }
        buffer = "";
      }
    }
  } catch {
    // registry ไม่พร้อมใช้งาน
  }
  return registry;
}

async function readMapping(): Promise<Mapping | null> {
  try {
    const content = await $`cat ${modelMappingPath}`.nothrow().text();
    return JSON.parse(content) as Mapping;
  } catch {
    return null;
  }
}

function stripPrefix(model: string): string {
  return model.replace(/^opencode\//, "");
}

function formatLimit(n?: number): string {
  if (!n) return "-";
  return n >= 1_000_000 ? `${n / 1_000_000}M` : n >= 1_000 ? `${Math.round(n / 1000)}K` : `${n}`;
}

function summarizeCapabilities(m: RegistryModel): string {
  const caps = m.capabilities?.input ?? {};
  const modalities: string[] = [];
  if (caps.image) modalities.push("รูป");
  if (caps.video) modalities.push("วิดีโอ");
  if (caps.audio) modalities.push("เสียง");
  if (caps.pdf) modalities.push("PDF");
  const parts: string[] = [];
  parts.push(`ctx ${formatLimit(m.limit?.context)}`);
  parts.push(`out ${formatLimit(m.limit?.output)}`);
  if (m.capabilities?.reasoning) parts.push("reasoning");
  if (modalities.length) parts.push(`multimodal(${modalities.join("/")})`);
  return parts.join(" · ");
}

/** แนะนำ role ที่ model ใหม่เหมาะกับ (จาก capability hints) */
function suggestRoles(m: RegistryModel, hints: Mapping["capability_hints"]): string {
  const caps = m.capabilities?.input ?? {};
  const multimodal = caps.image || caps.video || caps.audio;
  const longContext = (m.limit?.context ?? 0) >= 512_000;
  const longOutput = (m.limit?.output ?? 0) >= 64_000;
  if (multimodal && hints?.vision_roles?.length) return `vision → ${hints.vision_roles.join("/")}`;
  if (longContext && hints?.long_context_roles?.length) return `context ใหญ่ → ${hints.long_context_roles.join("/")}`;
  if (longOutput && hints?.long_output_roles?.length) return `output ยาว → ${hints.long_output_roles.join("/")}`;
  return "กลางๆ — ดูเพิ่มเติมได้จาก report";
}

function checkAgents(
  agents: AgentInfo[],
  registry: Map<string, RegistryModel>,
  mapping: Mapping | null
): CheckResult[] {
  const results: CheckResult[] = [];
  const hints = mapping?.capability_hints ?? {};
  const visionRoles = new Set(hints.vision_roles ?? []);
  const roleMap = new Map<string, string>();
  const currentMap = new Map<string, string>();
  const fallbackMap = new Map<string, string[]>();

  for (const [name, info] of Object.entries(mapping?.agents ?? {})) {
    roleMap.set(name, info.role);
    currentMap.set(name, info.current);
    fallbackMap.set(name, info.fallbacks ?? []);
  }

  for (const agent of agents) {
    agent.role = roleMap.get(agent.name) ?? "unknown";
    // effective model = ประกาศในไฟล์ ?? mapping (intended allocation)
    let effective = agent.model;
    let source: "file" | "mapping" = "file";
    if (!effective) {
      effective = currentMap.get(agent.name) ?? null;
      source = "mapping";
      if (effective) {
        agent.model = effective;
        agent.source = source;
      }
    }

    if (!effective) {
      results.push({
        agent: agent.name,
        role: agent.role,
        model: "none",
        source,
        status: "NO_MODEL",
        suggestion: "-",
      });
      continue;
    }

    const inRegistry = registry.has(stripPrefix(effective));
    const fallback = fallbackMap.get(agent.name)?.[0] ?? currentMap.get(agent.name) ?? "-";

    if (!inRegistry) {
      results.push({
        agent: agent.name,
        role: agent.role,
        model: effective,
        source,
        status: "RED",
        suggestion: `fallback: ${fallback}`,
      });
    } else if (visionRoles.has(agent.role) && stripPrefix(effective) !== "mimo-v2.5-free") {
      results.push({
        agent: agent.name,
        role: agent.role,
        model: effective,
        source,
        status: "GREEN",
        suggestion: "upgrade: opencode/mimo-v2.5-free (vision role)",
      });
    } else {
      results.push({
        agent: agent.name,
        role: agent.role,
        model: effective,
        source,
        status: "OK",
        suggestion: "-",
      });
    }
  }

  return results;
}

function findNewModels(
  registry: Map<string, RegistryModel>,
  results: CheckResult[],
  hints: Mapping["capability_hints"]
): { model: RegistryModel; suggestion: string }[] {
  const used = new Set(results.map((r) => stripPrefix(r.model)));
  const newModels: { model: RegistryModel; suggestion: string }[] = [];
  for (const m of registry.values()) {
    if (used.has(m.id)) continue;
    // เฉพาะ opencode gateway (ตัวที่ใช้ได้จริงกับ auth ปัจจุบัน) — registry ทั้งหมดมี ~400 ตัว
    // ข้าม alias รูปแบบ "~provider/name" (dynamic latest pointer)
    if (m.providerID !== "opencode" || m.id.startsWith("~")) continue;
    newModels.push({ model: m, suggestion: suggestRoles(m, hints) });
  }
  return newModels.sort((a, b) => a.model.id.localeCompare(b.model.id));
}

async function writeReport(
  results: CheckResult[],
  newModels: { model: RegistryModel; suggestion: string }[],
  timestamp: string
): Promise<void> {
  await $`mkdir -p /tmp/opencode`.nothrow();

  const counts = { RED: 0, GREEN: 0, OK: 0, NO_MODEL: 0 };
  for (const r of results) counts[r.status]++;

  let md = `# 🧠 Model Health Check
Generated: ${timestamp}
Source: ${agentsDir} ↔ \`opencode models --verbose\` ↔ model-mapping.json
Command: \`/model-check\` หรือ \`bun ~/.config/opencode/scripts/model-check.ts\`

## Summary
- Total agents: ${results.length}
- 🔴 RED (drift): ${counts.RED}
- 🟢 GREEN (advisory): ${counts.GREEN}
- ✅ OK: ${counts.OK}
- 🆕 New models ใน registry: ${newModels.length}
`;
  if (counts.NO_MODEL > 0) md += `- ⚠️ NO MODEL: ${counts.NO_MODEL}\n`;

  md += `
## Report
| Agent | Role | Model | Source | Status | Suggestion |
|-------|------|-------|--------|--------|------------|
`;
  for (const r of results) {
    const icon =
      r.status === "RED" ? "🔴 RED" :
      r.status === "GREEN" ? "🟢 GREEN" :
      r.status === "NO_MODEL" ? "⚠️ NO MODEL" : "✅ OK";
    const src = r.source === "mapping" ? "mapping" : "file";
    md += `| ${r.agent} | ${r.role} | ${r.model} | ${src} | ${icon} | ${r.suggestion} |\n`;
  }

  if (newModels.length) {
    md += `
## 🆕 New Models (ยังไม่มี agent ใช้)
| Model | Capabilities | เหมาะกับ |
|-------|--------------|----------|
`;
    for (const { model, suggestion } of newModels) {
      md += `| ${model.id} | ${summarizeCapabilities(model)} | ${suggestion} |\n`;
    }
  }

  md += `
---
*แก้ model allocation ได้ที่ \`~/.config/opencode/model-mapping.json\` + agent files (\`~/.config/opencode/agents/*.md\`)*
`;
  await $`echo ${md} > ${reportPath}`.nothrow();
}

function printSummary(
  results: CheckResult[],
  newModels: { model: RegistryModel; suggestion: string }[],
  timestamp: string
): void {
  const counts = { RED: 0, GREEN: 0, OK: 0, NO_MODEL: 0 };
  for (const r of results) counts[r.status]++;

  console.log(`🧠 Model Health Check — ${timestamp}`);
  console.log(`Total ${results.length} agents | 🔴 RED ${counts.RED} · 🟢 GREEN ${counts.GREEN} · ✅ OK ${counts.OK} · 🆕 NEW ${newModels.length}`);

  const reds = results.filter((r) => r.status === "RED");
  if (reds.length) {
    console.log("\n🔴 Drift (model หายจาก registry):");
    for (const r of reds) console.log(`  - ${r.agent} (${r.role}): ${r.model} → ${r.suggestion}`);
  }

  const greens = results.filter((r) => r.status === "GREEN");
  if (greens.length) {
    console.log("\n🟢 Advisory:");
    for (const r of greens) console.log(`  - ${r.agent} (${r.role}): ${r.model} → ${r.suggestion}`);
  }

  if (newModels.length) {
    console.log("\n🆕 New models:");
    for (const { model, suggestion } of newModels) {
      console.log(`  - ${model.id}: ${summarizeCapabilities(model)} → ${suggestion}`);
    }
  }

  console.log(`\nReport เต็ม: ${reportPath}`);
}

async function main(): Promise<void> {
  const timestamp = new Date().toISOString();
  const [agents, registry, mapping] = await Promise.all([
    readAgentModels(),
    getModelRegistry(),
    readMapping(),
  ]);
  const results = checkAgents(agents, registry, mapping);
  const hints = mapping?.capability_hints ?? {};
  const newModels = findNewModels(registry, results, hints);
  await writeReport(results, newModels, timestamp);
  printSummary(results, newModels, timestamp);
}

main().catch((e) => {
  console.error("model-check failed:", e);
  process.exit(1);
});