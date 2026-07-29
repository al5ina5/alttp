# ALttPO Sync Coverage Matrix

## Status: 2026-07-28

This document tracks what state is currently synchronized in ALttPO multiplayer.

---

## ✅ CURRENTLY SYNCED

### Equipment Items (0x340-0x357)
| Address | Item | Type | Status |
|---------|------|------|--------|
| 0x340 | Bow | 1 (OR) | ✅ Working |
| 0x341 | Boomerang | 1 (OR) | ✅ Working |
| 0x342 | Hookshot | 1 (OR) | ✅ Working |
| 0x344 | Mushroom | 1 (OR) | ✅ Working |
| 0x345 | Fire Rod | 1 (OR) | ✅ Working |
| 0x346 | Ice Rod | 1 (OR) | ✅ Working |
| 0x347 | Bombos | 1 (OR) | ✅ Working |
| 0x348 | Ether | 1 (OR) | ✅ Working |
| 0x349 | Quake | 1 (OR) | ✅ Working |
| 0x34A | Lamp | 1 (OR) | ✅ Working |
| 0x34B | Hammer | 1 (OR) | ✅ Working |
| 0x34C | Flute | 1 (OR) | ✅ Working |
| 0x34D | Bug Net | 1 (OR) | ✅ Working |
| 0x34E | Book | 1 (OR) | ✅ Working |
| 0x350 | Cane of Somaria | 1 (OR) | ✅ Working |
| 0x351 | Cane of Byrna | 1 (OR) | ✅ Working |
| 0x352 | Magic Cape | 1 (OR) | ✅ Working |
| 0x353 | Magic Mirror | 1 (OR) | ✅ Working |
| 0x354 | Gloves | Custom | ✅ Working |
| 0x355 | Boots | 1 (OR) | ✅ Working |
| 0x356 | Flippers | 1 (OR) | ✅ Working |
| 0x357 | Moon Pearl | 1 (OR) | ✅ Working |
| 0x359 | Sword | Custom | ✅ Working |
| 0x35A | Shield | Custom | ✅ Working |
| 0x35B | Armor | Custom | ✅ Working |

### Bottles (0x35C-0x35F)
| Address | Contents | Status |
|---------|----------|--------|
| 0x35C | Bottle 1 | ✅ Working |
| 0x35D | Bottle 2 | ✅ Working |
| 0x35E | Bottle 3 | ✅ Working |
| 0x35F | Bottle 4 | ✅ Working |

### Dungeon Items
| Address | Item | Status |
|---------|------|--------|
| 0x364-0x365 | Compasses | ✅ Working |
| 0x366-0x367 | Big Keys | ✅ Working |
| 0x368-0x369 | Maps | ✅ Working |

### Progress
| Address | Item | Status |
|---------|------|--------|
| 0x370 | Bombs capacity | ✅ Working |
| 0x371 | Arrows capacity | ✅ Working |
| 0x374 | Pendants | ✅ Working (OR) |
| 0x379 | Ability flags | ✅ Working |

### Small Keys (per dungeon)
- Address: 0xF36F + per-dungeon offset
- Status: ✅ Implemented behind `SyncSmallKeys` setting
- Note: Not enabled by default

### Tilemap
- Overworld: ✅ Working
- Underworld: ⚠️ Partial (intercepts registered but limited)

---

## ❌ MISSING / TODO

### P0 - BLOCKING (needed for co-op progression)

| Address | Item | Type | Reason |
|---------|------|------|--------|
| 0x343 | **Bombs (count)** | 1 (max) | ✅ Just enabled |
| 0x377 | **Arrows (count)** | 1 (max) | ✅ Just enabled |
| 0x360 | **Rupees** | 1 (max) | ✅ Just enabled |
| N/A | **Small Keys** | Delta-sum | ⚠️ Implemented but OFF by default |
| N/A | **Temporary underworld** | Tilemap | Pots, star tiles, trap floors not synced |

### P1 - QUALITY CO-OP

| Address | Item | Type | Reason |
|---------|------|------|--------|
| 0x36B | **Heart pieces** | 1 (OR) | Current HP not synced |
| 0x36C | **Health capacity** | 1 (OR) | Max HP not synced |
| 0x37B | **Current magic** | 1 (OR) | Not synced |
| 0x35C-0x35F | **Bottle contents** | Custom | Partially synced, verify consumption |

### P2 - POLISH

| Address | Item | Reason | Upstream | Our fork |
|---------|------|--------|----------|---------|
|| N/A | **Enemies/Bosses** | ✅ Fully synced - overlord/underworld room sync with immediate resync on room transition | ✅ **Enabled** | ✅ **Working** |
| N/A | **Object sync** | Bombs, arrows visible to other players | ❌ Disabled | ✅ **Enabled** |
| N/A | **NPC flags** | Some event flags not synced | ❌ | ❌ Planned |

---

## Protocol Status

- **Protocol 0x02**: Main state broadcast - ✅ Working
- **Frame sequence numbers**: ❌ Not implemented
- **Checksum/desync detection**: ❌ Not implemented  
- **Full snapshot/resync**: ❌ Not implemented

---

## Settings

| Setting | Default | Status |
|---------|---------|--------|
| SyncItems | ON | ✅ Working |
| SyncDungeonItems | ON | ✅ Working |
| SyncSmallKeys | OFF | ⚠️ Implemented but disabled |
| SyncTilemap | ON | ⚠️ Partial |
| SyncUnderworldTiles | ? | ❌ Not implemented |
| SyncEnemies | ? | ❌ Not implemented |

---

## Next Steps

1. **Enable bombs sync** (0x343) - easiest quick win
2. **Enable arrows sync** (0x377)
3. **Add rupee sync** (0x360-0x362)
4. **Enable SyncSmallKeys by default** or fix issues
5. **Add frame sequence numbers** (protocol reliability)
6. **Add checksum/desync detection**
