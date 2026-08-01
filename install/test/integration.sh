#!/usr/bin/env bash
# ─────────────────────────────────────────────────────────────────────────────
# Fah's P0 integration suite for `aiy warp` — SANDBOXED.
# Never touches the real ~/.config/opencode. Always HOME=/tmp/aiy-warp-test.
# Covers: fresh install, doctor=0, no-op, drift, --force, --dry-run, redaction
# gate (exit 5), export flat set, merge policy, exit codes 0/2/3/4/5.
# ─────────────────────────────────────────────────────────────────────────────
set -u

REPO_ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
SRC="$REPO_ROOT/install/src"
BIN="/tmp/aiy-warp-test/bin/aiy-warp"
SANDBOX_HOME="/tmp/aiy-warp-test"
KIT_SANDBOX="$SANDBOX_HOME/repo"          # kit copy for secret-plant tests
OUT="$SANDBOX_HOME/export-out"
DEST="$SANDBOX_HOME/.config/opencode"     # sandbox opencode dest
STATE="$SANDBOX_HOME/.config/aiy-warp"    # sandbox state dir (warp.lock)

FAILED=0
PASSED=0

pass() { PASSED=$((PASSED+1)); echo "  ✓ $1"; }
fail() { FAILED=$((FAILED+1)); echo "  ✗ $1"; }

# expect_exit <name> <expected_code> <cmd...>
expect_exit() {
  local name="$1" want="$2"; shift 2
  HOME="$SANDBOX_HOME" "$@" >/dev/null 2>&1
  local got=$?
  if [ "$got" -eq "$want" ]; then pass "$name (exit $got)"; else fail "$name: want exit $want, got $got"; fi
}

echo "── 0. Build ──"
mkdir -p "$SANDBOX_HOME/bin"
if (cd "$SRC" && go build -o "$BIN" .); then pass "go build"; else fail "go build"; exit 1; fi

echo "── 1. Fresh install (P0 happy path) ──"
rm -rf "$SANDBOX_HOME" && mkdir -p "$SANDBOX_HOME/bin" && cp "$BIN" "$SANDBOX_HOME/bin/" 2>/dev/null || true
# rebuild since we wiped sandbox
(cd "$SRC" && go build -o "$BIN" .) || exit 1
cd "$REPO_ROOT"

expect_exit "install opencode (fresh)" 0 env HOME="$SANDBOX_HOME" "$BIN" install opencode -y
AGENTS_INSTALLED=$(ls "$DEST/agents/"*.md 2>/dev/null | wc -l)
if [ "$AGENTS_INSTALLED" -eq 18 ]; then pass "18 agents installed"; else fail "agents installed = $AGENTS_INSTALLED, want 18"; fi
for f in aiy.md lin.md kwan.md; do
  [ -f "$DEST/agents/$f" ] && pass "agent file present: $f" || fail "missing agent: $f"
done
[ -f "$DEST/skills/aiy-messaging/SKILL.md" ] && pass "skill dir layout ok" || fail "skill layout wrong"
[ -f "$DEST/templates/Template_Daily_Log.md" ] && pass "templates copied" || fail "templates missing"
[ -f "$DEST/playbooks/Team_Charter.md" ] && pass "playbooks copied" || fail "playbooks missing"
[ -f "$STATE/warp.lock" ] && pass "warp.lock written" || fail "warp.lock missing"

echo "── 2. doctor = 0 (fresh install) ──"
expect_exit "doctor opencode (clean)" 0 env HOME="$SANDBOX_HOME" "$BIN" doctor opencode
DOCTOR_OUT=$(env HOME="$SANDBOX_HOME" "$BIN" doctor opencode 2>&1)
if echo "$DOCTOR_OUT" | grep -q "\[FAIL\]"; then fail "doctor reported FAIL"; else pass "doctor all PASS"; fi
if echo "$DOCTOR_OUT" | grep -q "all checks pass"; then pass "doctor summary ok"; else fail "doctor summary missing"; fi

echo "── 3. Reinstall → no-op (exit 4) ──"
expect_exit "reinstall (no-op)" 4 env HOME="$SANDBOX_HOME" "$BIN" install opencode -y

echo "── 4. Drift detection ──"
echo "# local tweak" >> "$DEST/agents/lin.md"
expect_exit "install with drift (skip)" 3 env HOME="$SANDBOX_HOME" "$BIN" install opencode -y
expect_exit "doctor after drift" 3 env HOME="$SANDBOX_HOME" "$BIN" doctor opencode
DOCTOR_OUT=$(env HOME="$SANDBOX_HOME" "$BIN" doctor opencode 2>&1)
if echo "$DOCTOR_OUT" | grep -q "\[FAIL\] 1."; then pass "doctor flags check 1 drift"; else fail "doctor check 1 not FAIL"; fi

echo "── 5. --dry-run writes nothing ──"
BEFORE=$(sha256sum "$DEST/agents/lin.md" | cut -d' ' -f1)
env HOME="$SANDBOX_HOME" "$BIN" install opencode --dry-run >/dev/null 2>&1
AFTER=$(sha256sum "$DEST/agents/lin.md" | cut -d' ' -f1)
if [ "$BEFORE" = "$AFTER" ]; then pass "dry-run did not write"; else fail "dry-run mutated files"; fi

echo "── 6. --force restores (merge policy) ──"
# tune host model, then force reinstall: kit restored, model preserved
python3 - "$DEST/agents/lin.md" << 'PY'
import sys
p = sys.argv[1]
s = open(p).read()
s = s.replace("model: opencode/deepseek-v4-flash-free", "model: opencode/host-tuned-model", 1)
open(p, "w").write(s)
PY
expect_exit "install --force" 0 env HOME="$SANDBOX_HOME" "$BIN" install opencode --force -y
if grep -q "model: opencode/host-tuned-model" "$DEST/agents/lin.md"; then
  pass "host model preserved by merge policy"
else
  fail "host model was clobbered by --force"
fi
expect_exit "doctor after force+merge (clean)" 0 env HOME="$SANDBOX_HOME" "$BIN" doctor opencode

echo "── 7. Redaction gate (exit 5) ──"
rm -rf "$KIT_SANDBOX" && cp -r "$REPO_ROOT" "$KIT_SANDBOX" && rm -rf "$KIT_SANDBOX/.git"
printf -- "---\ndescription: \"EVIL (อีวิล) — test\"\nmode: subagent\nmodel: opencode/x\ntoken: sk-test-abcdefghijklmnopqrstuvwxyz123456\n---\nbody\n" > "$KIT_SANDBOX/agents/evil.md"
(cd "$KIT_SANDBOX" && expect_exit "install with planted secret" 5 env HOME="$SANDBOX_HOME" "$BIN" install opencode -y)
(cd "$KIT_SANDBOX" && expect_exit "export with planted secret" 5 env HOME="$SANDBOX_HOME" "$BIN" export opencode --out "$OUT-evil")
# export gate on identifiers: real kit has Discord snowflake → export blocked without allow-identifiers
(cd "$REPO_ROOT" && expect_exit "export opencode (identifiers block)" 5 env HOME="$SANDBOX_HOME" "$BIN" export opencode --out "$OUT")

echo "── 8. Export flat set + --allow-identifiers ──"
(cd "$REPO_ROOT" && expect_exit "export with --allow-identifiers" 0 env HOME="$SANDBOX_HOME" "$BIN" export opencode --out "$OUT" --allow-identifiers 1527698229347487904,1210049942192010)
[ -f "$OUT/agents/aiy.md" ] && pass "export flat agents/aiy.md" || fail "export missing agents/aiy.md"
[ -f "$OUT/skills/aiy-messaging/SKILL.md" ] && pass "export flat skill" || fail "export missing skill"
STDOUT=$(env HOME="$SANDBOX_HOME" "$BIN" export opencode --stdout --allow-identifiers 1527698229347487904,1210049942192010 2>/dev/null | head -3)
if echo "$STDOUT" | grep -q "#### agents/"; then pass "export --stdout works"; else fail "export --stdout broken: $STDOUT"; fi

echo "── 9. Export with --agent selector ──"
(cd "$REPO_ROOT" && expect_exit "export --agent aiy" 0 env HOME="$SANDBOX_HOME" "$BIN" export opencode --agent aiy --out "$OUT-aiy" --allow-identifiers 1527698229347487904,1210049942192010)
[ -f "$OUT-aiy/agents/aiy.md" ] && pass "agent bundle has aiy.md" || fail "agent bundle missing aiy.md"
[ -f "$OUT-aiy/skills/obsidian/SKILL.md" ] && pass "agent bundle includes owned skills" || fail "agent bundle missing owned skill"

echo "── 10. Usage errors (exit 2) ──"
(cd "$REPO_ROOT" && expect_exit "unknown platform" 2 env HOME="$SANDBOX_HOME" "$BIN" install chatgpt)
(cd "$REPO_ROOT" && expect_exit "missing platform" 2 env HOME="$SANDBOX_HOME" "$BIN" install)
(cd "$REPO_ROOT" && expect_exit "--agent + --team conflict" 2 env HOME="$SANDBOX_HOME" "$BIN" install opencode --agent aiy --team kwan)
(cd "$REPO_ROOT" && expect_exit "unknown agent" 2 env HOME="$SANDBOX_HOME" "$BIN" install opencode --agent zzz)
(cd "$REPO_ROOT" && expect_exit "unknown team" 2 env HOME="$SANDBOX_HOME" "$BIN" install opencode --team bogus)
(cd "$REPO_ROOT" && expect_exit "export --collapse (P1)" 2 env HOME="$SANDBOX_HOME" "$BIN" export opencode --collapse)
(cd "$REPO_ROOT" && expect_exit "unknown command" 2 env HOME="$SANDBOX_HOME" "$BIN" sync)

echo "── 11. list + doctor --json ──"
LIST=$(env HOME="$SANDBOX_HOME" "$BIN" list 2>&1)
if echo "$LIST" | grep -q "18 agents"; then pass "list shows 18 agents"; else fail "list wrong: $(echo "$LIST" | head -1)"; fi
if echo "$LIST" | grep -q "aiy.*primary.*core"; then pass "list shows org chart"; else fail "list missing org chart"; fi
JSON=$(env HOME="$SANDBOX_HOME" "$BIN" doctor opencode --json 2>&1)
if echo "$JSON" | python3 -m json.tool >/dev/null 2>&1; then pass "doctor --json parses"; else fail "doctor --json invalid: $JSON"; fi

echo "── 12. doctor before any install (not installed) ──"
rm -rf "$SANDBOX_HOME/.config"
expect_exit "doctor on bare machine" 3 env HOME="$SANDBOX_HOME" "$BIN" doctor opencode

echo
echo "══════════════════════════════════════════════"
echo "PASSED: $PASSED   FAILED: $FAILED"
echo "══════════════════════════════════════════════"
[ "$FAILED" -eq 0 ]
