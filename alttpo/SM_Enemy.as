// Super Metroid enemy class
// Reads from the SM enemy pointer table at $7E:0F78 (16 entries, each 2 bytes)
// Each pointer points to a 0x20-byte enemy data block in WRAM
//
// Enemy data layout at the pointer target address:
//   0x00-0x01: enemy ID  (uint16)
//   0x02-0x03: X position (uint16)
//   0x04-0x05: Y position (uint16)
//   0x06-0x07: health     (uint16)
//   0x08-0x1F: other enemy state/properties

class SM_Enemy {
  uint8 enemy_index;
  uint8 host_index;
  bool is_active = false;
  array<uint8> enemy_data(0x20);
  uint16 pointer;
  uint16 health;
  uint16 Xpos;
  uint16 Ypos;

  SM_Enemy() {
    enemy_index = 0xFF;
    host_index = 0xFF;
    is_active = false;
    pointer = 0x0000;
    health = 0;
    Xpos = 0;
    Ypos = 0;
  }

  // Fetch enemy data for enemy at index i in the pointer table
  void fetch(uint8 i) {
    enemy_index = i;

    // read 16-bit pointer from the enemy pointer table at $7E:0F78
    pointer = bus::read_u16(0x7E0F78 + i * 2);

    if (pointer == 0x0000) {
      is_active = false;
      health = 0;
      Xpos = 0;
      Ypos = 0;
      return;
    }

    is_active = true;

    // read 0x20 bytes of enemy data from the pointer address
    bus::read_block_u8(pointer, 0, 0x20, enemy_data);

    // extract useful fields from the raw data
    Xpos   = uint16(enemy_data[0x02]) | (uint16(enemy_data[0x03]) << 8);
    Ypos   = uint16(enemy_data[0x04]) | (uint16(enemy_data[0x05]) << 8);
    health = uint16(enemy_data[0x06]) | (uint16(enemy_data[0x07]) << 8);
  }

  // Write the 0x20 bytes of enemy data back to the pointer address
  void apply() {
    if (!is_active || pointer == 0x0000) {
      return;
    }

    // write back any modified enemy data
    bus::write_block_u8(pointer, 0, 0x20, enemy_data);
  }
};
