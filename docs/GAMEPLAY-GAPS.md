# Game-Breaking Progression Desync Analysis — ALTTP Multiplayer

**This is NOT about code.** It's about which permanent world/story changes can break the game if they desync between players.

Each entry: **Event** → What changes → Is it synced? → Consequence if desynced.

---

## PHASE 1: LIGHT WORLD — HYLIA TO FIRST PENDANT

### 1. Rescue Zelda → Sanctuary
**World change:** Zelda disappears from throne room, appears in Sanctuary.
**Synced?** Tilemap change (throne room floor crack + hidden passage). Should sync via TilemapChanges.
**Game-breaker?** Low — visual only, doesn't gate anything.

### 2. Uncle Sword / Killing Soldiers
**World change:** Uncle dies, sword acquired. 
**Synced?** Item $7EF359 (sword level) is synced. Uncle body is animation state, not persistent.
**Game-breaker?** No — sword acquisition is item-synced.

### 3. Lost Woods → Master Sword Pedestal
**World change:** The log bridge falls (tilemap change).
**Synced?** Tilemap sync should handle this.
**Game-breaker?** Low-medium — visual pathfinding only.

### 4. Sahasrahla's Hideout (Boots)
**World change:** None permanent. He gives boots.
**Synced?** Item $7EF355 (boots) is synced.
**Game-breaker?** No — item sync.

### 5. Kakariko Village — Blind's House / Well Drain
**World change:** Drain the well ($7E:2000+ tilemap changes), Blind's house state changes.
**Synced?** Tilemap sync should trigger. BUT — the drain event is a one-time world change that only fires once. If Player A drains the well while Player B is in a dungeon, Player B's tilemap might not update until they return to the overworld.
**Game-breaker?** MEDIUM — if the well isn't drained for Player B, they can't enter and miss key items (Heart Piece, etc.). The same applies to Blind's house.

### 6. Blacksmiths → Return to Shop
**World change:** After you take the smith to his partner, both return to the shop and start forging. This changes NPC locations permanently.
**Synced?** Progress flag in $7EF3C5-$7EF3C9 range (synced). But the NPC position in the world might not update if the tilemap/sprite state doesn't reflect it.
**Game-breaker?** Medium — you can't get the tempered sword upgrade if this is wrong. Also the Super Bomb quest starts here.

### 7. Witch's Mushroom Trade
**World change:** Give mushroom → Powder is for sale in the hut. The witch changes behavior.
**Synced?** Item $7EF344 (mushroom) is synced as OR. When one player uses it, it's consumed. The other player sees the consumption.
**Game-breaker?** Low-medium — if the witch doesn't sell powder for one player, they can't transform some enemies.

### 8. King Zora Moves / Zora's Domain Water Drained
**World change:** Pay 500 rupees → King Zora scoots, revealing the passage to Zora's Domain. The water level drops when you have flippers.
**Synced?** King Zora's position is a room flag (tilemap change). The water level change is a permanent world tile change.
**Game-breaker?** ⚠️ **HIGH** — If Player A pays and King Zora moves, but Player B didn't sync that tile, they see King Zora still blocking the path. Can't enter Zora's Domain. Can't get flippers. **Softlock.**

### 9. Zora Flippers
**World change:** None permanent, item acquisition.
**Synced?** Item $7EF356 (flippers) is synced.
**Game-breaker?** No — item sync.

---

## PHASE 2: EASTERN PALACE — FIRST PENDANT

### 10. Eastern Palace — Big Chest / Boss
**World change:** Boss death flag + pendant awarded.
**Synced?** Pendant $7EF374 synced (OR). Boss death flag is a room flag — might be in TilemapChanges or room state.
**Game-breaker?** ⚠️ **HIGH** — If Player A kills the boss and gets the pendant, but Player B's game doesn't register the boss death. Player B enters the boss room and the boss is alive again (or the door is still locked). **Either way, it's broken.** The pendant/chest logic needs to be solid.

---

## PHASE 3: DESERT PALACE — SECOND PENDANT

### 11. Book of Mudora + Desert Tablet
**World change:** Book opens the Desert Palace entrance (tile change at the石碑).
**Synced?** Book ($7EF34E) synced as item. The stone door opening is a tilemap change.
**Game-breaker?** Medium — if the door doesn't open for Player B, they can't enter Desert Palace. But this is just a tilemap sync issue.

---

## PHASE 4: TOWER OF HERA — THIRD PENDANT

### 12. Tower of Hera — Big Key / Boss
**World change:** Boss death + pendant.
**Synced?** Same as Eastern Palace.
**Game-breaker?** Medium-High — progression gating.

---

## PHASE 5: MASTER SWORD → DARK WORLD

### 13. Master Sword Pull
**World change:** The sword disappears from the pedestal (tilemap change). The Lost Woods exit to the pedestal clearing opens.
**Synced?** Item $7EF359 (sword level 3). Tilemap for sword pedestal should sync.
**Game-breaker?** ⚠️ **HIGH** — If Player A pulls the sword but Player B didn't sync it, Player B goes to the pedestal and sees the sword still there, can't pull it again (the flag is consumed). They might try to pull it forever. Also, the Master Sword is required for Agahnim — can't progress without it.

### 14. First Agahnim Fight
**World change:** PYRAMID APPEARS. Dark World seals are broken. This is THE MOST IMPORTANT story change in the game.
**Synced?** Progress flag $7EF3C5 bits (world state). This IS synced.
**But:** The Pyramid APPERANCE on the overworld map is a tilemap change. If Player A beats Agahnim, teleports to the Pyramid, but Player B is still in the Light World — Player B sees no Pyramid. They can't enter the Pyramid Hole. They're stuck in Light World.
**Game-breaker?** ⚠️ **CRITICAL** — The post-Agahnim world state change (Pyramid, Dark World access, Ganon's Tower availability) MUST propagate to all players. If the Pyramid doesn't appear for all players, the game is **softlocked.**

---

## PHASE 6: DARK WORLD — CRYSTAL HUNT (7 DUNGEONS)

### 15. Moon Pearl — Required Item
**World change:** Required to be in Dark World without turning into a bunny.
**Synced?** Item $7EF357 synced.
**Game-breaker?** ⚠️ **HIGH** — Without moon pearl, you're a bunny (can't use items). If one player doesn't have it synced, they're stuck in bunny form in Dark World. Cannot progress.

### 16. Dark World Flute
**World change:** Playing the flute at the stump in Light World → bird appears, fast travel unlocked.
**Synced?** The flute item is synced. But the ACTIVATION/flute status (bird unlocked) is at $7EF392 according to some RAM maps. **Check if this is in the sync range.**
**Game-breaker?** Medium — without bird, travel is slower but still possible.

### 17. Dark World Death Mountain — Gloves/Mittens
**World change:** Titan's Mitt (lift dark rocks) and Power Glove (lift light rocks) unlock paths.
**Synced?** Item $7EF354 (gloves) synced as custom.
**Game-breaker?** Medium — without mittens, can't access some Dark World areas or Turtle Rock.

### 18. Hammer Pegs / Spectacle Rock Cave
**World change:** Hammering pegs opens paths. Cave at Spectacle Rock in Dark World.
**Synced?** Tilemap changes + progress flags. **Race to the cave: if Player A breaks the pegs and Player B doesn't see them broken, Player B can't enter.**
**Game-breaker?** Medium.

### 19. Misery Mire — Medallion Gate
**World change:** Requires Ether or Bombos to enter (or Quake? — depends on game mode).
**Synced?** Items $7EF347-$7EF349 synced. The medallion gate is a one-time state change.
**Game-breaker?** ⚠️ **HIGH** — Misery Mire is one of the 7 crystals. If one player gets Ether and opens the gate, but the other player doesn't have the medallion OR the gate doesn't appear open for them, they can't enter. Crystal can't be obtained.

### 20. Turtle Rock — Medallion Gate
**World change:** Requires Quake (or Bombos/Ether depending on mode).
**Synced?** Same issue as Misery Mire.
**Game-breaker?** ⚠️ **HIGH** — Same problem. One of the 7 crystals.

### 21. Swamp Palace — Water Drain
**World change:** The drain/floodgate must be opened. The water drains, changing the tiles significantly.
**Synced?** Tilemap change for the water level. But if the trigger flag doesn't propagate, Player B can't enter.
**Game-breaker?** Medium-High — crystal dungeon, required.

### 22. Thieves' Town / Blind's House Portal
**World change:** The portal opens in Blind's house in Kakariko after certain flags.
**Synced?** Progress flags + tilemap.
**Game-breaker?** Medium — if the portal doesn't appear, can't enter Thieves' Town.

### 23. Skull Woods
**World change:** Fire rod needed to burn bushes blocking entrance. Multiple exits.
**Synced?** Fire rod $7EF345 synced. But the burnt bushes are tilemap changes.
**Game-breaker?** Medium.

### 24. Ice Palace
**World change:** Fire rod needed to melt ice.
**Synced?** Fire rod synced. Ice melting is a room state per-entrance attempt.
**Game-breaker?** Low — you can re-enter and re-melt.

### 25. Gargoyle's Domain (Misery Mire)
**World change:** Switch puzzles.
**Synced?** Room state / tilemap changes.
**Game-breaker?** Medium — if switches don't sync, puzzles can't be solved.

---

## PHASE 7: ENDGAME — GANON

### 26. All 7 Crystals → Ganon's Tower
**World change:** When all 7 crystals are collected, Ganon's Tower opens in Dark World.
**Synced?** $7EF37A (crystals) synced.
**Game-breaker?** ⚠️ **CRITICAL** — If Player A has all 7 crystals but Player B doesn't (desync), Player B can't enter Ganon's Tower. **Known bug:** Issue #28 on upstream — "Unable to enter Ganon's tower on reduced crystal mode."

### 27. Ganon's Tower — Maiden Rescue
**World change:** Rescuing maidens removes barriers. Each maiden rescue is a progress flag.
**Synced?** $7EF3C9 range (progress flags 2/2). Synced.
**Game-breaker?** ⚠️ **HIGH** — If Player A rescues a maiden and doesn't sync to Player B, Player B can't progress through Ganon's Tower. The barrier remains.

### 28. Agahnim 2 / Ganon Fight
**World change:** Pyramid Hole opens after Agahnim 2. Ganon fight begins.
**Synced?** The pyramid hole is a tilemap change + progress flag.
**Game-breaker?** ⚠️ **CRITICAL** — If only one player beats Agahnim 2, the pyramid hole opens for them. But the other player can't enter Ganon's room. **Endgame softlock.**

### 29. Silver Arrows
**World change:** Required to kill Ganon.
**Synced?** Item $7EF340 bow tier is synced.
**Game-breaker?** ⚠️ **HIGH** — If Player B doesn't have silver arrows but Player A does, Player B can't damage Ganon. Need to ensure the Bow + Silver upgrade flag propagates.

---

## PHASE 8: SIDE EVENTS / OPTIONAL PROGRESSION GATES

### 30. Ether / Bombos Tablet Sword Upgrade
**World change:** Reading the tablet with the book + medallion upgrades your sword.
**Synced?** Sword $7EF359 synced. But the tablet reading is a one-time event.
**Game-breaker?** Low-medium — tempered/gold sword isn't required for progression, but makes it much easier.

### 31. Flute Boy Death
**World change:** If you go to Dark World BEFORE getting the flute, the flute boy dies and you find the shovel at his grave. This changes the flute acquisition path permanently.
**Synced?** The flute is synced as an item. But if one player triggers the death event and the other doesn't, the ALTERNATE path (shovel → flute) might not work for the second player.
**Game-breaker?** Medium — you can still get the flute, just differently. Both end up with the same item.

### 32. Graveyard / Magic Cape
**World change:** Opening the grave with the shovel to get the Magic Cape.
**Synced?** Item $7EF352 synced. The grave opening is a tilemap change.
**Game-breaker?** Low — optional item.

### 33. Super Bomb Quest
**World change:** Complex chain: Smith → dwarf → blacksmith shop → Super Bomb → Pyramid.
**Synced?** Each step is a progress flag. Chain requires all steps to be consistent.
**Game-breaker?** Medium — Super Bomb isn't strictly necessary (can beat Ganon without silver arrows on Normal mode), but the chain is fragile.

### 34. Magic Upgrade NPCs (Sawmill, Spectacle Rock Cave)
**World change:** Once you pay, the upgrade is applied. The NPC doesn't give it again.
**Synced?** Magic meter capacity $7EF37B? Actually that's current magic. The upgrade flag might be separate.
**Game-breaker?** Low — optional, but lost upgrade.

---

## ⚠️ SUMMARY: CRITICAL GAME-BREAKERS (Priority Order)

| Priority | Event | Why It Breaks |
|----------|-------|--------------|
| 🔴 P0 | **Agahnim 1 → Pyramid** | Light/Dark world split. If only one player sees the pyramid, the other can't progress at all. |
| 🔴 P0 | **Crystal sync (all 7)** | Ganon's Tower won't open if crystals desync (Issue #28). |
| 🔴 P0 | **Boss death flags (all)** | If a boss respawns for one player but is dead for another, the door logic breaks. |
| 🔴 P0 | **Small key sync** | Already addressed by SyncSmallKeys. If it fails, players get locked out of dungeon rooms. |
| 🟠 P1 | **King Zora / NPC state** | If NPC position desyncs, entire areas become inaccessible (Zora's Domain, Smiths). |
| 🟠 P1 | **Medallion gates (Mire + TR)** | If gate doesn't open for both players, they can't both enter crystal dungeons. |
| 🟠 P1 | **Moon Pearl** | Bunny form = can't use items = can't progress in Dark World. |
| 🟠 P1 | **Silver Arrows** | Ganon invulnerable without them. |
| 🟡 P2 | **Well drain / Blind's house** | Key items locked behind area state changes. |
| 🟡 P2 | **Water level changes** | Swamp Palace, Zora's Domain drain. |
| 🟡 P2 | **Kakariko well / portal** | Thieves' Town entrance. |

## Testing Recommendation

For EACH of the P0 events, create a test scenario:
1. Two players in different locations (one in overworld, one in dungeon)
2. Player A triggers the event
3. Player B transitions to the affected area
4. Verify the world state matches

If any P0 event doesn't propagate, the game is **unplayable for multiplayer co-op** regardless of how pretty the code is.
