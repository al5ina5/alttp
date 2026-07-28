# Reference: Sprite Sync & PvP (from mysterypaintwo esync branch)

## PvP Damage System — PvPAttack.as

```angelscript
class PvPAttack {
  uint16 player_index;  // which player is being damaged
  uint8 sword_time;     // $3C value
  uint8 melee_item;     // 0=sword, 1=bugnet?, 2=hammer
  uint8 ancilla_mode;   // ancilla type (0 for melee)
  uint8 damage;         // 8 damage = 1 whole heart
  int8 recoil_dx;       // stagger direction X
  int8 recoil_dy;       // stagger direction Y
  int8 recoil_dz;       // stagger direction Z
};
```

### Damage calculation (in Projectile.as)
```angelscript
bool calc_damage(GameState@ defendant, GameState@ attacker) {
  damage = 0;
  // armor: reduce damage by bit-shifting right
  // damage >>= armor_shr;
  
  // weapon-specific:
  // sword = 4*8
  // hammer = 2*8  (reduced in esync)
  // silver arrows = reduced in esync
}
```

The esync branch nerfs PvP damage for hammer and silver arrows (commit `22363e0c`).

### How PvP triggers (in pre_frame.as):
```
if (settings.EnablePvP) {
    local.attack_pvp();
}
```
Called after `receive()` and before `send()`.

### How damage is sent:
PvPAttack struct is serialized as part of `GameState` — the attacker includes their attack data, the receiver applies damage to the target player.

---

## Sprite Sync — Sprite.as

OAM sprite handling with unique tile distribution:

### OAM Table Decoding
```angelscript
void decodeOAMTableBytes() {
    // b0-b3 = 4 standard OAM bytes
    // b4 = 5th byte (extended OAM table)
    
    x = int32(uint16(b0) | (uint16(b4 & 1) << 8));
    if (x >= 256) x -= 512;  // sign extension
    
    size = (b4 >> 1) & 1;
    y = int32(b1);
    chr = uint16(b2) | (uint16(b3 & 1) << 8);
    palette = (b3 >> 1) & 7;
    priority = (b3 >> 4) & 3;
    hflip = ((b3 >> 6) & 1);
    vflip = ((b3 >> 7) & 1);
}
```

### Unique Tile Capture (SpritesWindow.as)
```angelscript
class SpritesWindow {
    GUI::SNESCanvas @canvas;
    array<uint16> page0(0x1000);
    array<uint16> page1(0x1000);
    
    void render(const array<uint16> &palette) {
        // Read VRAM CHR pages
        ppu::vram.read_block(ppu::vram.chr_address(0x000), 0, 0x1000, page0);
        ppu::vram.read_block(ppu::vram.chr_address(0x100), 0, 0x1000, page1);
        
        // Draw as 4bpp tiles
        canvas.draw_sprite_4bpp(0, 0, 0, 128, 128, page0, palette);
        canvas.draw_sprite_4bpp(128, 0, 0, 128, 128, page1, palette);
    }
}
```

### Unique Tile Sync (ROMSprites.as / alttp-sprites.as)
The esync branch captures unique sprite tiles from VRAM and distributes them:
- Captures newly encountered unique tiles per-frame
- Sends only tiles not seen before by the remote player
- Uses NAK mechanism: if receiver doesn't have the tile, it sends NAK, sender re-broadcasts
- Works for both ALTTP + SM sprites

### VRAM capture (in pre_frame.as):
```
local.capture_sprites_vram();
```
Called before `receive()` in the sync lifecycle.

---

## Player Names — Player Labels Under Sprites

**Branch:** https://github.com/mysterypaintwo/alttpo/tree/player-names  
**Files changed:** GameState.as, LocalGameState.as, SettingsWindow.as, init.as

Simple approach:
- Each player's name is stored in `GameState`
- Serialized alongside position data
- Rendered as text beneath the player sprite (CHR-based)
- "Show Labels" checkbox in SettingsWindow

```angelscript
// LocalGameState.as additions:
// 0x04 - Name Player Mode
// TODO: read player name from SRAM
// name = bus::read_block_u8(0x7EF3D9);

void serialize_name(array<uint8> &r) {
    r.write_str(namePadded);
    // group name: 20 bytes exactly
    serialize_name(envelope);
}
```

To implement: add name field to `GameState.serialize()/deserialize()`, render text via VRAM font tiles beneath the remote player's sprite position.
