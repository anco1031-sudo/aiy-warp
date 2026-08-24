#!/usr/bin/env bash
# =============================================================================
# git-team — จัดการ Git identity ของทีม (Aiy Ecosystem)
#
# หลักการ: ใช้ auth เดียวของ Louis (token/SSH ที่มีอยู่) แต่แยก "ตัวตน" ด้วย
#          git identity ต่อ repo — commit จะลงชื่อ agent คนนั้นๆ
#
# Usage:
#   git-team list                          # ดู identity ทั้งหมด
#   git-team <agent>                       # ตั้ง identity ให้ repo ปัจจุบัน (cwd)
#   git-team <agent> <repo-path>           # ตั้ง identity ให้ repo ที่ระบุ
#   git-team show <repo-path>              # ดู identity ปัจจุบันของ repo
#
# ตัวอย่าง:
#   cd 01-Projects/trading-bot
#   git-team lin                           # commit ต่อจากนี้เป็น "Lin (หลิน)"
# =============================================================================

set -euo pipefail

# ===== Identity Map (18 agents) =====
# email ใช้ domain @agents.aiy.local — ไม่ใช่ email จริง (กัน commit ไป link กับ
# account จริงของ Louis ผิดคน; authorship แยกที่ชื่อ)
declare -A TEAM_IDENTITY=(
  [aiy]="Aiy (อัย)|aiy@agents.aiy.local"
  [lin]="Lin (หลิน)|lin@agents.aiy.local"
  [an]="An (อัน)|an@agents.aiy.local"
  [pao]="Pao (เปา)|pao@agents.aiy.local"
  [mint]="Mint (มิ้นท์)|mint@agents.aiy.local"
  [fah]="Fah (ฟ้า)|fah@agents.aiy.local"
  [cloud]="Cloud (คลาวด์)|cloud@agents.aiy.local"
  [jean]="Jean (จีน)|jean@agents.aiy.local"
  [cher]="Cher (เฌอ)|cher@agents.aiy.local"
  [ploy]="Ploy (พลอย)|ploy@agents.aiy.local"
  [jewel]="Jewel (จิวเวล)|jewel@agents.aiy.local"
  [sai]="Sai (ทราย)|sai@agents.aiy.local"
  [pleng]="Pleng (เพลง)|pleng@agents.aiy.local"
  [kwan]="Kwan (ขวัญ)|kwan@agents.aiy.local"
  [fon]="Fon (ฝน)|fon@agents.aiy.local"
  [june]="June (จูน)|june@agents.aiy.local"
  [bee]="Bee (บี)|bee@agents.aiy.local"
  [nam]="Nam (น้ำ)|nam@agents.aiy.local"
)

cmd_list() {
  echo "👥 Git Team Identities (18 agents)"
  echo "==================================="
  for agent in aiy lin an pao mint fah cloud jean cher ploy jewel sai pleng kwan fon june bee nam; do
    local name email
    name="${TEAM_IDENTITY[$agent]%%|*}"
    email="${TEAM_IDENTITY[$agent]##*|}"
    printf "  %-8s %-25s %s\n" "$agent" "$name" "$email"
  done
  echo "==================================="
  echo "💡 ใช้: git-team <agent> [repo-path] — ตั้ง identity ให้ repo"
}

cmd_set() {
  local agent="$1"
  local repo="${2:-.}"
  if [[ -z "${TEAM_IDENTITY[$agent]:-}" ]]; then
    echo "❌ ไม่มี agent ชื่อ '$agent' (ดูได้จาก: git-team list)"
    return 1
  fi
  if [[ ! -d "$repo/.git" ]]; then
    echo "❌ '$repo' ไม่ใช่ git repository (ไม่มี .git)"
    return 1
  fi
  local name email
  name="${TEAM_IDENTITY[$agent]%%|*}"
  email="${TEAM_IDENTITY[$agent]##*|}"
  git -C "$repo" config user.name "$name"
  git -C "$repo" config user.email "$email"
  echo "✅ ตั้ง identity '$name' <$email> ให้ $repo แล้ว"
  echo "   (repo นี้ commit ต่อจากนี้ลงชื่อ: $name)"
}

cmd_show() {
  local repo="${1:-.}"
  if [[ ! -d "$repo/.git" ]]; then
    echo "❌ '$repo' ไม่ใช่ git repository"
    return 1
  fi
  local name email
  name="$(git -C "$repo" config user.name 2>/dev/null || echo "(ยังไม่ตั้ง)")"
  email="$(git -C "$repo" config user.email 2>/dev/null || echo "(ยังไม่ตั้ง)")"
  echo "👤 $repo → $name <$email>"
}

cmd="${1:-list}"
case "$cmd" in
  list|-l|--list)   cmd_list ;;
  show|--show)      shift; cmd_show "${1:-.}" ;;
  help|-h|--help)
    sed -n '2,16p' "${BASH_SOURCE[0]}" | sed 's/^# \{0,1\}//'
    ;;
  *)                cmd_set "$cmd" "${2:-.}" ;;
esac