# Reference: Enemy Sync Implementation (from mysterypaintwo esync branch)

**Source:** https://github.com/mysterypaintwo/alttpo/tree/esync

## Architecture

Enemy sync works by:
1. Each frame, fetch enemy data from SM RAM
2. The "host" player (lowest index) owns the enemy state
3. When host attacks an enemy, the delta is broadcast to all players
4. Non-host players see the updated enemy position/health via network packets
5. On room entry, enemy sync is delayed briefly to let the room settle

## Key File: SM_Enemy.as (131 lines)

```angelscript
const uint sm_enemy_bank = 0x7E0F78;
const uint sm_enemy_count = 0x7E0E4E;
const uint sm_enemy_kills = 0x7E0E50;
const uint sm_enemy_active_distance = 30000;

class SM_Enemy {
    uint8 enemy_index;
    uint8 host_index;
    bool is_active = false;
    array<uint16> enemy_data(0x20);
    uint16 pointer;
    uint16 health;
    uint16 Xpos;
    uint16 Ypos;

    // Called each frame to read current enemy state from SM RAM
    void fetch(uint8 i) { ... }
    
    // Called after receiving network sync to write state back to SM RAM
    void apply() { ... }
}
```

### fetch() — read enemy from RAM
Reads SM enemy pointer table at `0x7E0F78 + slot*2`, then reads:
- Health at `pointer + 0x02`
- X position at `pointer + 0x0A`  
- Y position at `pointer + 0x0E`
- Full 0x20-byte enemy data block at `pointer + 0x00`
- Segment data at `0x7FFC00` (used for segmented enemies like lanmola, moldorm)

### apply() — write enemy state back to RAM
Writes all 0x20 bytes back to the enemy's pointer address. Special cases:
- Room `0xD73F`: skip (elevators)
- Room `0xD0BF`: skip (ship/gunship)
- Room `0xD07F`: skip (ship part 2)
- Rooms `0xE87F`, `0xEAFF`: bound = 24 (fewer enemies to sync)
- Others: bound = 32

Decrement `sm_enemy_count` and increment `sm_enemy_kills` when enemy health <= 0.

## Key File: EnemyWindow.as (102 lines)

Debug UI showing:
- 16 enemies displayed as rows
- Each row: host_index, health, X, Y, pointer, active flag
- Colored labels: yellow=active, dark=disabled
- Updates every frame from `local.enemyData`

## Key Changes in LocalGameState.as (+913 lines)

Added to `send()`:
- `fetch_enemy_data()` — reads all SM enemies from RAM
- `serialize_sm_enemy(i)` — packs enemy data into two packets (0x20 bytes each)
- Serializes: enemy_index, enemy_data[0x20], host_index, is_active
- Segment data for segmented enemies

Added to receive path:
- Deserialize SM enemy data from two-packet split
- Apply remote enemy state to local RAM

Added `fetch_enemy_data()` method:
- Reads enemy pointer table  
- For each active enemy, reads pointer → follows to get health/X/Y/data
- Checks distance from local player (`get_distance_from_enemy`)
- Only syncs enemies within `sm_enemy_active_distance` (30000 units)

## Room Entry Delay

The esync branch adds a delay on room entry:
```
a0aaa47b 2025-07-25 Patch enemy sync by delaying synchronization on room entry
```
This prevents desyncs caused by reading enemy data while the room is still initializing.

## Crystal Switch Fix

```
ffcc0871 2025-06-14 esync: fixed crystal switches to not freeze non-owning players
```
When a player activates a crystal switch in a boss room, non-owning players no longer get stuck with the frozen crystal animation.

---

## How to Port (recommended approach)

1. Copy `SM_Enemy.as` and `EnemyWindow.as` into `alttpo/`
2. Add `#include "SM_Enemy.as"` to `init.as`
3. Add `fetch_enemy_data()` call in `LocalGameState.send()` before serialization
4. Add `serialize_sm_enemy()` / deserialize logic to `GameState.as`
5. Add `SyncEnemies` checkbox to `SettingsWindow.as`
6. Add enemy data handling in `pre_frame.as` (after `local.fetch()`, before `local.send()`)
7. Test: two players enter a room with enemies, killing should sync
