# Underworld Tile Sync — Architecture Report & Implementation Plan

## 1. Tilemap Sync Architecture

### Memory Layout

| Region | Address Range | Format | Content |
|--------|--------------|--------|---------|
| Tile data | `$7E:2000-5FFF` | 16-bit per tile, 0x2000 entries | BG tile indices (tile number) |
| Attributes | `$7F:2000-3FFF` | 8-bit per tile, 0x2000 entries | BG tile attribute/palette bytes |

The tilemap is organized as two background layers in the underworld:
- **BG1** (layer 1): tiles `0x0000-0x0FFF` — main room tiles, walls, floors, pots, etc.
- **BG2** (layer 2): tiles `0x1000-0x1FFF` — overlay sprites, decorations, star tiles, traps

The underworld always uses a **64×64 tile grid** (0x40 × 0x40 = 0x1000 tiles per layer).

### Write Interceptors (capture layer)

Two buses are intercepted in `LocalGameState::register()` (lines 202-203):

```
bus::add_write_interceptor("7e:2000-5fff", ...this.tilemap_written);
bus::add_write_interceptor("7f:2000-3fff", ...this.attributes_written);
```

These capture *any* write by the game to the tilemap WRAM, storing them into a `TilemapChanges state[]` array (size 0x2000, `int32` per entry — bottom 16 bits = tile, top 8 bits = attribute, `-1` = unchanged).

### Capture Flow

1. **`tilemap_written()`** (line 1019) — intercepts byte writes to `$7E:2000-5FFF`. It reads the full 16-bit word, figures out which tile index was written, and captures the tile + attribute into `tilemap[i]`.

2. **`attributes_written()`** (line 978) — intercepts byte writes to `$7F:2000-3FFF`. It either samples the tile+attribute together (if this index has never been written before) or just overwrites the attribute.

3. **Safety guards** — both functions check `is_safe_to_sample_tilemap()` which filters out bad times (screen transitions, stair climbing, etc.). Notably **barrier tiles** (0x66/0x67) written during `sub_module == 0x16` are rejected because crystal switch sync already handles them.

4. **Timestamp** — every captured write sets `tilemapTimestamp = timestamp_now`.

### Send Flow (outgoing)

`send_tilemaps()` (line 1490) runs every frame:

1. **Compress** — `tilemap.compress_runs()` converts the sparse `state[]` array into horizontal/vertical runs of changed tiles (using `TilemapRun`).
2. **Packet** — timestamp + location + runs are packed into UDP packets (type `0x07`).
3. **Fragmentation** — multiple packets are sent if the runs exceed `MaxPacketSize` (1452 bytes).
4. **Envelope** — broadcast-to-sector (type `0x02`) ensures only players in the same room receive it.

### Receive Flow (incoming)

`deserialize_tilemaps()` (line 829 of GameState.as):

1. Reads timestamp, location, and runs into `tilemapRuns[]`.
2. Each run is deserialized via `TilemapRun::deserialize()`.

### Apply Flow (local integration)

`update_tilemap()` (line 2161 of LocalGameState.as) runs after `on_main_alttp`:

1. Checks `is_safe_to_write_tilemap()` — same safety checks as capture.
2. For each remote player in the same room: if the remote's `tilemapTimestamp > local tilemapTimestamp`, apply their runs via `tilemap.apply_wram()` (writes to WRAM `$7E2000+` and `$7F2000+`).
3. If `write_to_vram` is enabled, also calls `tilemap.apply_vram()` to update VRAM so the visual changes appear on screen immediately.

### TilemapRun Compression

`compress_runs()` (line 191 of TilemapChanges.as):
- Iterates the `state[0x2000]` array in scanline order, covering both BG layers.
- Produces either vertical or horizontal runs depending on which is longer.
- Each run is either "same" (all tiles identical) or "diff" (individual tiles).
- Already-changed entries are marked with `-1` to skip during compression.

---

## 2. Current Tilemap Sync — What Works vs What's Missing

### ✅ Already Working

| Feature | How it works |
|---------|-------------|
| **Overworld tile changes** | `determine_vram_bounds_overworld()` + VRAM writes via `ow_tilemap_to_vram_address()` conversion. Cross-team sync allowed. |
| **Underworld room transitions** | On room change, `tilemap.reset(0x40)` clears state, `tilemapLocation` updates, `tilemapTimestamp` resets to 0 so incoming packets are accepted. |
| **Crystal switch / barrier tiles** | Separate `SyncableByte` at `$7E:C172` syncs the crystal switch state across all players in the same dungeon. Barrier attribute writes (0x66/0x67) during `sub_module 0x16` are explicitly ignored by the attribute interceptor. |
| **Tilemap apply to WRAM** | `apply_wram()` writes both tile data (`$7E2000+`) and attributes (`$7F2000+`) back to the local game state. |
| **Tilemap apply to VRAM** | `apply_vram()` converts tilemap indices to VRAM addresses and writes directly to PPU VRAM for immediate visual feedback. |
| **Pot items (contents)** | `DoorRandomizerMapping` syncs pot pickup bits via extra SRAM at `$7F:6018`/`$7F:6268` using `update_extras()`. This syncs *whether the pot has been looted*, but not the visual tile change. |

### ❌ Missing / Broken

| Feature | Problem | Root Cause |
|---------|---------|-----------|
| **Pot shattering (visual)** | When a pot is shattered by dashing, the pot tile changes to a broken pot tile or floor. These tilemap writes ARE captured by the interceptor, but the tile data changes are small and may be missed if timing is tight. The visual change does not propagate to other players correctly. | The write interceptor catches the tile data write, but the room layout reloads on entering — pots that were "whole" for the remote player's initial load remain whole because the base tilemap is read from the ROM room layout data, not from the synced tilemap state. |
| **Star tiles** | Star tile puzzles (e.g., Turtle Rock) change tile attributes/visuals when activated. These writes go to `$7E:2000+` / `$7F:2000+` and ARE intercepted. However, star tiles use **BG2** for their overlay, and the system may not handle BG2 tile changes properly in all cases. | The `compress_runs()` iterates over `0x0000-0x1FFF` (both layers), so it should capture both. The issue is more likely that star tile attribute writes are not properly captured when they happen during certain submodules that `is_safe_to_sample_tilemap()` rejects. Also, star tiles rely on precise attribute bytes that change quickly. |
| **Trap floors (floor traps that drop you)** | When triggered, the floor tile changes from solid to hole. This is a tilemap write to `$7E:2000+`. But the write may happen during a submodule that is explicitly blocked by the safety check (e.g., during fall animation). | The submodule guard blocks capture during transitions. When a trap triggers, the module/submodule state may be in a transition that is filtered out. |
| **Trap doors (switch-activated)** | Sub_module `0x17` runs when stepping on a switch to open trap doors. This is a known transition but the write guard may block it. | `is_safe_to_sample_tilemap()` does NOT list `sub_module == 0x17` as a blocked state for `module == 0x07`. But `is_safe_to_write_tilemap()` also does not explicitly handle `sub_module == 0x17`. The issue might be that the tilemap write happens during the submodule transition and the receiving player doesn't re-apply it. |
| **Pots lifted by remote player** | When another player lifts a pot, the tile underneath appears. The pot entity becomes a sprite (held by Link), but the base tile becomes visible. The tile data write for "pot gone → floor visible" should be captured. | Possible issue: the tile underneath is often part of the base room tilemap loaded from ROM — it's not explicitly "written" by the game as a change. The game might just clear the sprite overlay for the pot without modifying the tilemap. |
| **Initial room load sync** | When entering a room, `tilemap.reset()` clears all state to `-1`. The tilemap timestamp goes to 0. Remote players who have been in the room longer have a later timestamp, so their changes should be applied. But the reverse — changes made by a player who JUST entered — might not propagate to established players in the room. | Established players have `tilemapTimestamp > 0`. When a new player enters, their timestamp is 0, so their packets are rejected by established players' `if (remote.tilemapTimestamp > tilemapTimestamp)` check (line 2264). |

---

## 3. Detailed Analysis of Specific Underworld Tile Systems

### 3.1 Pots

Pot behavior in ALTTP:
- **Pot as object**: A pot is both a tile (BG1 tile entry) and a sprite object (`$0D00-0FA0` area). The tile shows a pot graphic.
- **Lifting a pot**: The pot sprite appears above Link's head. The tile underneath becomes visible (usually floor). The game does NOT typically write a new tile value — it relies on the sprite layer covering the pot tile. When the pot is lifted, the sprite is just moved.
- **Shattering a pot**: When dashed into, the pot "shatters" — the game writes a new tile value to `$7E:2000+` to show broken pot pieces or the underlying floor tile.

What's currently synced:
- **Pot items** (contents): via `DoorRandomizerMapping.update_extras()` — pot pickup bits in extra SRAM are OR-merged across players.
- **Pot tile changes**: The write interceptor should capture the `$7E` tile data write when a pot shatters. This gets sent to other players via the tilemap sync.

What's NOT synced:
- **Pot visibility**: When player A lifts a pot, player B cannot see the pot floating in the air. The pot sprite is handled by the `enableObjectSync` system (object sync at `$7E:0D00`), which does sync objects but may not work reliably for lifted pots.
- **Pot shattering visual sync**: Works in theory through tilemap sync, but the floor tile that replaces the pot is just the base room tile. If both players are in the room, the tile change should propagate. However if player A shatters a pot and player B enters the room afterward, player B's room load will display the pot as whole (from ROM base data), and the tilemap sync will only apply if player A is still in the room and still has the change in their `tilemap.state[]`.

### 3.2 Star Tiles

Star tiles appear in Turtle Rock (trinexx dungeon). They work as follows:
- The tilemap has two layers: BG1 has the base room, BG2 has the star overlay.
- Stepping on or shooting certain star tiles changes their **attribute byte** (palette/priority) to indicate activated/deactivated state.
- The game writes to `$7F:2000+` (attribute array) to toggle these tiles.

Why they might not sync:
- The attribute interceptor at `$7F:2000` captures all attribute writes — so it SHOULD capture star tile toggles.
- But the VRAM rendering for underworld in `write_vram()` (line 130-149) only writes a single 16-bit tile to VRAM: `ppu::vram[vram] = tile`. It does NOT write the attribute byte for the underworld path. Looking at the code:
  ```
  ppu::vram[vram] = tile;
  ```
  This writes a 16-bit value to VRAM. ALTTP VRAM entries are 16-bit: low byte is tile number, high byte is attribute/palette (the format is VHAPPPTT TTTTTTTT where V=vertical flip, H=horizontal flip, A=priority, PPP=palette). So actually the 16-bit tile value includes the attribute bits in the high byte.
- The `write_tilemap()` function stores `tile` as `uint16(c & 0x00ffff)` and `attr` as `uint8 ((c & 0xff0000) >> 16)`. The attribute is stored in the top 8 bits of the int32 state. But in the VRAM write, it only uses `tile` (bottom 16 bits).
- **KEY ISSUE**: The `write_tilemap()` function stores the 16-bit tile number and 8-bit attribute separately. When writing to WRAM (`$7E2000+`), it writes `tile` only (16-bit). The attribute goes to `$7F2000+`. But the VRAM format for ALTTP uses a 16-bit entry where the high byte IS the attribute. So the attribute needs to be OR'd into the 16-bit tile value before writing to VRAM.

Let me verify this by looking at `write_vram()` more carefully. The function writes `ppu::vram[vram] = tile` where `tile` is `uint16(c & 0x00ffff)` — this is only the tile number, stripped of the `attr` bits! 

**BUG FOUND**: For the underworld, `write_vram()` writes only the 16-bit tile number to VRAM, but it should OR the attribute byte into the high byte of the VRAM entry. The correct format is:
```
ppu::vram[vram] = tile | (uint16(attr) << 8);
```

This would explain why star tile attribute changes (palette changes) don't appear visually — the tile number might be correct but the attribute (which controls palette, flip, priority) is missing from the VRAM write.

### 3.3 Trap Floors & Trap Doors

Trap floors:
- Walk-over floor tiles that collapse, dropping Link into a pit.
- When triggered, the game writes a hole tile to `$7E:2000+` to replace the solid floor tile.
- The capture happens, but the trap trigger animation may set `sub_module` to a value that `is_safe_to_sample_tilemap()` rejects (e.g., the falling animation).

Trap doors:
- Switch-activated floor panels that open/close.
- `sub_module == 0x17` is noted in comments as "Quick little submodule that runs when you step on a switch to open trap doors."
- `is_safe_to_sample_tilemap()` does NOT explicitly block module 0x07 sub_module 0x17, so these SHOULD be captured. But `is_safe_to_write_tilemap()` also doesn't block it, so they SHOULD be written too.
- The issue might be that trap doors use attribute changes (like barrier tiles) or that the game writes to a different memory region.

---

## 4. Implementation Plan

### Phase 1: Fix VRAM Attribute Bug (Highest Priority)

**Problem**: `write_vram()` for underworld does not include the attribute byte in the VRAM write.

**Location**: `TilemapChanges.as`, line 149
**Current code**:
```
ppu::vram[vram] = tile;
```
**Should be**:
```
uint16 vramEntry = tile | (uint16(attr) << 8);
ppu::vram[vram] = vramEntry;
```

But `attr` is not available in `write_vram()` currently. The function takes `(uint i, int32 c)` and extracts `tile` from `c`, but discards `attr`. Need to:
1. Extract `attr` from `c` in `write_vram()`: `uint8 attr = uint8((c & 0xff0000) >> 16);`
2. OR the attribute into the high byte of the VRAM write.

**Impact**: Fixes star tile visual sync, barrier tile visual, any tile that relies on palette/attribute changes.

### Phase 2: Capture Timing Improvements

**Problem**: Some underworld tile changes happen during submodules that are blocked by safety checks.

**Assessment**:
- `sub_module == 0x17` (trap door switch): currently NOT blocked in either `is_safe_to_sample_tilemap()` or `is_safe_to_write_tilemap()`. This is fine.
- Trap floor activation during fall animation: The write may happen before the submodule changes. Need to verify in the ROM code exactly when the tile data write happens relative to the submodule change.

**Action**: Trace the exact submodule values during:
- Pot shattering (on overworld and underworld)
- Tile switch stepping (module 07, submodule 17)
- Trap floor activation
- Star tile activation

If writes happen during a blocked submodule, we have two options:
1. **Narrow the block**: Remove the submodule guard for specific known-safe transitions.
2. **Delayed capture**: Buffer writes during blocked periods and process them when the submodule returns to 0x00.

### Phase 3: Initial Room Load Sync

**Problem**: When player B enters a room that player A has already modified (e.g., shattered pots, activated star tiles), player B's base room tilemap is loaded from ROM data, not from A's modified state. Player B only receives A's tilemap changes if A is still sending packets.

**Flow analysis**:
1. Player A is in room R, has `tilemapTimestamp = T_A_high` with various tile changes.
2. Player B enters room R: `tilemap.reset()` sets all to -1, `tilemapTimestamp = 0`.
3. Player B's `update_tilemap()` checks: `remote.tilemapTimestamp > tilemapTimestamp` → if A has T_A_high > 0, then A's changes are accepted. ✓
4. Player B broadcasts their tilemap (all -1, no changes) — correct.
5. BUT: what if player A leaves the room? Player A's tilemap gets reset. Player B now has no one to receive changes from, and the tiles show as "whole" pots etc.

**Solution**: Store per-room tilemap state in a persistent cache. When a room is re-entered, apply the cached modifications. The changes are already being sent — the problem is that when no one is in the room to send them, the state is lost.

Implementation options:
1. **Persistent room state array**: Keep an `array<TilemapChanges@>` indexed by room number. On room exit, save tilemap state. On room enter, restore from cache.
2. **SRAM-backed tilemap**: Write certain tilemap changes back to a persistent SRAM region so they survive room transitions. This is what the crystal switch already does (it syncs to SRAM at `$7E:C172`).
3. **Server-side tilemap storage**: Have players send tilemap changes to the server as authoritative room state, and new arrivals download the saved state.

**Recommended approach**: Combine (1) and (2) — store tilemap diffs in an array per room. When entering a room:
- Look up cached tilemap for this room.
- If another player is in the room and has changes, apply those instead.
- If no one is in the room, apply the cached state.

### Phase 4: Pot Tile Visual Sync

**Problem**: Pot shattering tile changes work in theory but have coverage gaps.

**Analysis**:
- `enableObjectSync` syncs object data at `$7E:0D00-0FA0`. This handles pot objects.
- The tile that appears when a pot is lifted/shattered is not a separate write — it's the base tile underneath the pot becoming visible.
- The game writes to the tilemap only when the pot shatters (breaking animation), not when the pot is lifted.

**Solution**:
1. The existing `enableObjectSync` system should be extended to handle pot lift/shatter state.
2. When a pot is shattered, both the tile change (via tilemap interceptor) AND the pot object state change are synced.
3. When a player enters a room, if the pot items are already claimed (pot SRAM bits are set), the pots should visually appear as broken/shattered. This requires modifying the room load to check pot state and force tilemap writes.

This is the hardest part because it requires:
- ROM-level patching to overlay pot tiles when loading a room (like how the crystal switch ROM hook works).
- Or: after room load, walk through the pot item SRAM and apply tile changes to match.

### Phase 5: Safety Check Review

**Action**: Audit all submodule guards in both `is_safe_to_sample_tilemap()` and `is_safe_to_write_tilemap()` for the underworld module 0x07.

Current blocked states: 0x01 (scroll), 0x02 (load supertile), 0x06/0x07 (floor transition), 0x08/0x0e/0x10/0x12/0x13 (stairs), 0x15 (warp), 0x18 (crystal), 0x19 (mirror).

**Recommendations**:
1. Confirm submodule 0x17 (trap door switch) is NOT blocked — it's currently not listed, which is correct.
2. Consider if submodule 0x0a (light level / Agahnim room) or 0x0b/0x0c/0x0d (water changes) should allow tilemap capture.
3. Test if pot shattering writes happen during submodule 0x00 only (normal gameplay) — if so, they're already captured.

---

## 5. Summary of All Needed Changes

| # | Change | File | Complexity | Priority |
|---|--------|------|-----------|----------|
| 1 | Fix VRAM write to include attribute byte | `TilemapChanges.as:149` | 1 line | **Critical** |
| 2 | Persist room tilemap state on exit, restore on enter | `LocalGameState.as:894-922` | Medium | **High** |
| 3 | Ensure `apply_wram()` for underworld correctly OR's attribute | `TilemapChanges.as:89-95` | Already done* | Check |
| 4 | Audit submodule guards for underworld capture | `LocalGameState.as:925-975` | Low | Medium |
| 5 | Add pot shatter tile overlay on room load | `LocalGameState.as` + ROM hook | High | Medium |
| 6 | Test star tile attribute visual sync after Phase 1 fix | Testing | Low | High |
| 7 | Test trap floor/door visual sync | Testing | Low | High |

*Already done — `write_tilemap()` writes separate tile to `$7E2000+` and attr to `$7F2000+`. The issue is solely in the VRAM write.

---

## 6. Key File Reference

| File | Purpose |
|------|---------|
| `TilemapChanges.as` | Core tilemap state storage, compression, WRAM/VRAM write |
| `TilemapRun.as` | Horizontal/vertical run representation + serialization |
| `GameState.as` (base class) | `tilemap[]`, `tilemapTimestamp`, `tilemapLocation`, `deserialize_tilemaps()` |
| `LocalGameState.as` | `tilemap_written()`, `attributes_written()`, `send_tilemaps()`, `update_tilemap()`, `fetch_tilemap_changes()`, safety guards |
| `pre_frame.as` | `on_main_alttp()` calls `update_tilemap()` |
| `pre_nmi.as` | `fetch_module()`, `fetch_sfx()` — runs before main |
| `receive.as` | UDP receive + `process_message()` → `GameState.deserialize()` |
| `ROMMapping.as` | ROM addresses, `DoorRandomizerMapping` pot sync, `serialize_extras()` |
| `SyncableItem.as` | `SyncableUnderworldRoom` — SRAM room state sync |
| `SyncableByte.as` | Crystal switch sync at `$7E:C172` |
