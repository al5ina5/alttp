# ALTTP Co-op QA Smoke Test Checklist

**Testers:** User (Naples FL) + Dad  
**Purpose:** Incremental smoke test of all synced game systems in ALTTP Together  
**Prerequisites:** Both players have bsnes-as + ALTTP Together mod installed, same ROM file (SHA256 verified)

---

## (A) Setup Verification

### A1. Network Connection Test
| Step | Action | Expected Result |
|------|--------|-----------------|
| A1.1 | Both players launch bsnes-as with ALTTP Together mod | Mod loads without errors |
| A1.2 | Player 1 creates new group: opens Multiplayer menu → "Create Group" → enter name "ALTTP-TEST-001" | Group created, server shows "Waiting for players..." |
| A1.3 | Player 2 joins: Multiplayer → "Join Group" → enters "ALTTP-TEST-001" | Player 2 connects, both see each other in player list |
| A1.4 | Verify both clients show same ping/latency | Latency displayed in ms (aim for <100ms) |

**Menu Path:** `Multiplayer → Create Group` / `Multiplayer → Join Group`

### A2. ROM Checksum Verification
| Step | Action | Expected Result |
|------|--------|-----------------|
| A2.1 | Both players open Debug menu → "ROM Info" | Displays ROM name, size, checksum |
| A2.2 | Compare checksum displayed on both screens | **MUST MATCH EXACTLY** — any mismatch will cause sync issues |

**Menu Path:** `Debug → ROM Info`  
**Expected:** Same 8-character hex checksum on both clients (e.g., `6B9C2F1D`)

### A3. Initial Player Positions
| Step | Action | Expected Result |
|------|--------|-----------------|
| A3.1 | Both players start new game (Player 1 hosts) | Both spawn at Link's House, Player 1 in front, Player 2 behind |
| A3.2 | Verify Player List window shows both names | Shows "Player 1 (You)" and "Player 2" |

**Verification:** Both players should see identical overworld view

---

## (B) Overworld Tests

### B1. Basic Movement Sync
| Step | Action | Expected Result |
|------|--------|-----------------|
| B1.1 | Player 1 walks east toward forest | Player 2 sees Player 1's sprite move in real-time |
| B1.2 | Player 2 walks in different direction | Player 1 sees Player 2's sprite move |

**Expected:** < 3 frame difference in sprite positions

### B2. Rupee Pickup Sync (Economy Sync)
| Step | Action | Expected Result |
|------|--------|-----------------|
| B2.1 | Player 1 picks up a green rupee (3 rupees) on overworld | Both screens show rupee counter = 3 |
| B2.2 | Player 2 picks up another rupee nearby | Both screens show rupee counter = 6 |
| B2.3 | Player 1 opens inventory (Select button) | Both see 6 rupees in wallet |

**Menu Path:** `Select → Inventory` (rupees shown at top)  
**Expected:** Exact rupee count matches on both clients at all times

### B3. Tree Cutting Sync
| Step | Action | Expected Result |
|------|--------|-----------------|
| B3.1 | Both players stand near a cuttable tree (Lost Woods entrance) | Tree visible on both screens |
| B3.2 | Player 1 swings sword at tree | Tree falls on Player 2's screen too |
| B3.3 | Verify tree state | Tree shows "cut" state on both clients |

**Expected:** Animation plays identically on both screens

### B4. Item Pickup Sync
| Step | Action | Expected Result |
|------|--------|-----------------|
| B4.1 | Player 1 picks up a heart container on overworld | Both see heart container disappear |
| B4.2 | Player 1 opens status screen (Start → Status) | Max HP shows increased on both screens |

**Menu Path:** `Start → Status`  
**Expected:** HP/metadata identical on both clients

---

## (C) Sanctuary + Hyrule Castle Tests

### C1. Small Keys Sync
| Step | Action | Expected Result |
|------|--------|-----------------|
| C1.1 | Both enter Hyrule Castle (Sanctuary entrance) | Both screens show Castle interior |
| C1.2 | Player 1 picks up small key in first room | Key count increases to 1 on both screens |
| C1.3 | Player 2 opens door requiring key (A button) | Door unlocks on both screens |

**Menu Path:** `Start → Inventory → Keys` (shows key count)  
**Expected:** Small key count synced in real-time

### C2. Sword Pickup Sync
| Step | Action | Expected Result |
|------|--------|-----------------|
| C2.1 | Both find sword in Hyrule Castle (unlock room) | Sword on pedestal visible |
| C2.2 | Player 1 touches sword | Sword disappears, Player 1 now has sword |
| C2.3 | Player 2 opens status screen | Sword icon appears in inventory for Player 1 |

**Expected:** Weapon state (sword level) synced

### C3. Sword Swing Visibility
| Step | Action | Expected Result |
|------|--------|-----------------|
| C3.1 | Both stand in Castle hallway with enemies | Guards patrol on both screens |
| C3.2 | Player 1 swings sword (A button) | Sword swing animation visible on Player 2's screen |
| C3.3 | Player 2 swings sword | Player 1 sees Player 2's sword swing |

**Expected:** Attack animations sync in real-time

### C4. Guard Kill Sync
| Step | Action | Expected Result |
|------|--------|-----------------|
| C4.1 | Player 1 attacks a guard | Guard HP decreases |
| C4.2 | Guard dies | Guard disappears on both screens |
| C4.3 | Check both player views | No desync — guard dead on both |

**Expected:** Enemy death state identical

### C5. Item Chest Sync
| Step | Action | Expected Result |
|------|--------|-----------------|
| C5.1 | Both approach a chest in Castle | Chest visible, closed on both |
| C5.2 | Player 1 opens chest | Item appears, chest open on both |
| C5.3 | Both see same item | e.g., "Fighter Shield" acquired |

**Expected:** Chest state (open/closed) + item synced

---

## (D) Eastern Palace Boss (Armos Knights)

### D1. Boss Entry
| Step | Action | Expected Result |
|------|--------|-----------------|
| D1.1 | Both enter Eastern Palace boss room | Armos Knights arena loads |
| D1.2 | Boss cutscene plays | Both see intro animation |

**Expected:** Cutscene in sync

### D2. Boss Damage/Death Sync
| Step | Action | Expected Result |
|------|--------|-----------------|
| D2.1 | Player 1 attacks Armos Knights | Boss takes damage, HP bar decreases |
| D2.2 | Player 2 attacks same boss | Both see boss HP decrease together |
| D2.3 | Player 1 lands final blow | Boss dies on both screens |

**Expected:** Boss HP synced across clients

### D3. Boss Door & Reward Sync
| Step | Action | Expected Result |
|------|--------|-----------------|
| D3.1 | After boss death, door opens | Exit door opens on both screens |
| D3.2 | Player 1 touches pendant (Pendant of Power) | Pendant collected |
| D3.3 | Both open inventory/status screen | Pendant shows in both inventories |

**Menu Path:** `Start → Status` → shows pendants  
**Expected:** Both see same pendant count (1/3)

### D3.4. Loot Distribution
| Step | Action | Expected Result |
|------|--------|-----------------|
| D3.5 | Check Player List window for boss loot icon | Boss kill credited to both players |

---

## (E) Desync Recovery Tests

### E1. Client Crash Recovery
| Step | Action | Expected Result |
|------|--------|-----------------|
| E1.1 | Both progress to midpoint in Hyrule Castle | Both at same room, same state |
| E1.2 | **Dad (Player 2) closes emulator** | Player 1 sees "Player 2 disconnected" |
| E1.3 | Player 2 relaunches bsnes-as, rejoins group | Both reconnect |
| E1.4 | Verify game state | Both at same room, same key count, same HP |

**Expected:** Snapshot relay brings Player 2 back in sync

### E2. Network Interruption Test
| Step | Action | Expected Result |
|------|--------|-----------------|
| E2.1 | Both playing, Player 1 has 50 rupees | Both screens show 50 |
| E2.2 | Player 2's network drops for 5 seconds | Reconnection attempts |
| E2.3 | Network restores | Both sync to same rupee count |

**Expected:** FNV-1a checksum verification passes after reconnect

### E3. Host Migration (if implemented)
| Step | Action | Expected Result |
|------|--------|-----------------|
| E3.1 | Player 1 (host) quits game gracefully | Option to transfer host |
| E3.2 | Player 2 becomes host | Game continues without restart |

---

## (F) Known Problem Areas — Watch List

### F1. Dungeon Transitions
| Issue | What to Watch | Expected Fix |
|-------|---------------|--------------|
| F1.1 | Entering/exiting dungeon | Both players in same room after transition |
| F1.2 | Dungeon map/compass sync | Items appear in inventory for both |
| F1.3 | Room-to-room loading | No desync during scroll |

**Test:** Walk through Eastern Palace → Dark Maze → boss room, verify positions

### F2. Save-and-Quit
| Issue | What to Watch | Expected Fix |
|-------|---------------|--------------|
| F2.1 | Both players Save & Quit | Saves current state |
| F2.2 | Both reload save file | Both at same position |
| F2.3 | Inventory items preserved | Keys, items, rupees all present |

**Menu Path:** `Start → Save` (in-game)  
**Test:** Save at Sanctuary, quit, reload, verify sync

### F3. Mirror Warps
| Issue | What to Watch | Expected Fix |
|-------|---------------|--------------|
| F3.1 | Using Magic Mirror in Dark World | Both teleport to same location |
| F3.2 | Returning via portal | Both in Light World at same spot |
| F3.3 | Position desync after warp | Check XY coordinates |

**Test:** Warp from Dark Sanctuary → Light World Sanctuary

### F4. Ganon's Tower Elevator
| Issue | What to Watch | Expected Fix |
|-------|---------------|--------------|
| F4.1 | Entering Ganon's Tower elevator | Both in elevator |
| F4.2 | Elevator moves | Both at same floor |
| F4.3 | Exiting elevator | Both in same room |

**Test:** Reach Ganon's Tower, use elevator, verify positions

### F5. Boss Cutscenes
| Issue | What to Watch | Expected Fix |
|-------|---------------|--------------|
| F5.1 | Any boss intro animation | Both see same frames |
| F5.2 | Boss death animation | Plays identically |
| F5.3 | Post-boss cutscene (e.g., Moldorm) | Both in sync |

**Test:** Kill each of the 8 dungeon bosses, verify cutscene sync

---

## (G) Quick Reference — Menu Paths

| Screen | Menu Path |
|--------|-----------|
| Multiplayer Create/Join | `Multiplayer → Create Group` / `Join Group` |
| ROM Info | `Debug → ROM Info` |
| Inventory | `Select → Inventory` |
| Status | `Start → Status` |
| Keys | `Select → Inventory → Keys` |
| Save Game | `Start → Save` (in-game) |
| Player List | `Multiplayer → Player List` (window) |

---

## (H) Acceptance Criteria Summary

- [ ] Both clients connect with matching ROM checksums
- [ ] Rupee count identical at all times
- [ ] Keys, items, boss kills sync correctly
- [ ] Both see each other's sword swings
- [ ] Armos Knights death triggers door + pendant for both
- [ ] Desync recovery works after Player 2 reconnects
- [ ] All 5 problem areas tested without crash/desync

---

**Checklist Complete:** Both testers sign off after each section.

| Tester | Date | Signature |
|--------|------|-----------|
| User (Naples FL) | _________ | _________ |
| Dad | _________ | _________ |

---

*Generated for ALTTP Together v0.9+ QA testing*
