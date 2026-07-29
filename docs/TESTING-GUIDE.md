# ALttP Multiplayer Testing Guide

## Quick Start
1. Load same ROM checksum on both clients (verify in game menu)
2. Connect to same group name
3. Choose different player colors (Blue/Red for you/dad)
4. Start at Link's House (safe test zone)

---

## Sync Coverage Status

### ✅ Working (Test These First)
| System | Sync Method | Test Location |
|--------|-------------|---------------|
| Equipment (Bow, Sword, Boots, etc.) | SRAM 0x340-0x357 | Anywhere |
| Bottles (slots) | SRAM 0x35C-0x35F | Anywhere |
| Dungeon Items (Compass, Big Key, Map) | SRAM 0x364-0x369 | Dungeon entrance |
| Progress (Pendants, Ability flags) | SRAM 0x370-0x379 | Check menu |
| Bombs/Arrows/Rupees capacity | SRAM 0x343, 0x377, 0x360-0x362 | Shop purchases |
| Overworld Tilemap | SRAM 0x280-0x340 | Cut trees, pickup items |
| SRAM-based Underworld Rooms | SRAM 0x00-0x250 | Dungeons |
| Player Movement | WRAM sprite sync | Anywhere |
| Player Items (sword swings, rod use) | WRAM sprite sync | Test in open area |
| Boss States (dead/alive) | SRAM room flags bit 11 | Boss rooms |
| Enemy/Overlord Positions | WRAM 0x0B00-0x0B57 | Enemy rooms |

---

## Detailed Test Scenarios

### 1. Equipment Sync
**Test:**
1. Both players at Link's House
2. Player A picks up sword from chest
3. Verify Player B sees sword in Player A's inventory
4. Player A uses sword on signpost
5. Verify Player B sees sword swing animation

**Expected:** Both players see same inventory, sword swing visible

---

### 2. Bottle Content Sync (Known Issue)
**Current Behavior:**
- Bottle slots sync (which bottle is selected)
- Bottle contents (red/green potion, bees, fairy) sync ONCE when filled
- **Consumption does NOT sync** - if Player A drinks potion, Player B still sees full bottle

**Test:**
1. Both players have red potions in bottles
2. Player A drinks potion
3. Check inventory screen
4. Verify both still see full bottles (known limitation)

**Workaround:** Share bottles at merchant before consumption

---

### 3. Boss Death Propagation
**Test (Eastern Palace - Armos Knights):**
1. Player A enters boss room first (don't fight)
2. Player B enters boss room behind A
3. Player A kills Armos Knights
4. Verify:
   - Door opens for both players
   - Pendant appears in inventory for both
   - Boss sprite disappears for both

**Timing Note:** Boss death syncs via SRAM every 32 frames (~533ms). If Player B enters after death but before sync, they may briefly see boss alive.

**Immediate Sync Trigger:** Room transition (both entering boss room, or exiting and re-entering) triggers immediate SRAM sync.

---

### 4. Overworld Tilemap Sync
**Test:**
1. Both players in Hyrule Castle area
2. Player A cuts down tree
3. Verify Player B sees tree gone
4. Player A picks up hidden item under tree
5. Verify Player B sees item appear

**Expected:** Both players see same tilemap state

---

### 5. Underworld Room Sync
**Test (Sanctuary):**
1. Both players in Sanctuary
2. Player A opens small chest
3. Verify Player B sees chest empty
4. Player A picks up item from chest
5. Verify Player B sees item taken

**Expected:** Chest state synced (room flags bit 4-9)

---

### 6. Enemy Sync (Bosses)
**Test:**
1. Both players in boss room (e.g. Eastern Palace)
2. Player A stands in boss room
3. Player B enters boss room
4. Verify boss position/HP synced in real-time (every 8 frames)
5. Both players damage boss simultaneously
6. Verify boss dies when both have dealt enough damage

**Expected:** Real-time boss position/health sync

---

### 7. Desync Recovery
**Test (Crash Recovery):**
1. Player B crashes emulator (close bsnes)
2. Player A continues playing
3. Player B restarts emulator with same ROM
4. Reconnect to same group
5. Verify state resyncs (full snapshot transfer)

**Expected:** Player B receives full state snapshot on reconnect

---

## Problem Areas (Watch For)

| Issue | Location | Mitigation |
|-------|----------|------------|
| Boss movement not synced | Boss rooms | Use boss states instead |
| Bottle consumption not synced | Anywhere | Fill bottles before sharing |
| Save/Quit may desync | Menu → Save | Avoid during co-op |
| Mirror warps dual-layer | Misery Mire | Enter one at a time |
| Ganon's Tower elevator | Ganon's Tower | Use separately |
| Boss cutscenes | Boss rooms | Enter after cutscene |

---

## QA Checklist Reference
See `docs/QA-COOP-CHECKLIST.md` for step-by-step manual testing procedure.

---

## Quick Troubleshooting

| Symptom | Likely Cause | Fix |
|---------|--------------|-----|
| Player sees wrong boss state | SRAM sync delay | Exit and re-enter room |
| Bottle shows full but empty | Consumption not synced | Known limitation, share pre-filled |
| Room state out of sync | SRAM sync interval | Enter room together, wait 1s |
| Sprite desync | Network latency | Check connection quality |
| Full state resync triggered | Player reconnected | Normal, wait 2-3s |

---

## Test Sign-Off

| Test Area | Player A | Player B | Date |
|-----------|----------|----------|------|
| Equipment Sync | [ ] | [ ] | |
| Bottle Content | [ ] | [ ] | |
| Boss Death | [ ] | [ ] | |
| Overworld Tiles | [ ] | [ ] | |
| Underworld Rooms | [ ] | [ ] | |
| Enemy Sync | [ ] | [ ] | |
| Desync Recovery | [ ] | [ ] | |
