# ALttPO on Couchforge

## Goal
- **aio**: hosts UDP `alttp-server` on `:4590` and auto-joins as player AIO
- **kiosk**: auto-joins `aio.local:4590` group `couchforge` as player Kiosk
- No Join window / no mouse required

## Install layout
`~/alttp-test/` (or `~/Games/alttpo/`):
- `bin/bsnes` — bsnes-as
- `bin/alttp-server` — UDP group server
- `bin/discord_game_sdk.so`
- `alttpo/` — AngelScript scripts
- `roms/alttp_usa.sfc` — legal ROM

## Launchers
`port-alttp-host` / `port-alttp-join` → `couchforge-port-adapter alttp-host|alttp-join`

Port markers: **ALttPO (Host)** / **ALttPO (Join)** (Multiworld names still alias).

## Verify
```bash
# aio
~/.local/bin/couchforge-port-adapter alttp-host
# kiosk
COUCHFORGE_ALTTP_HOST=aio.local ~/.local/bin/couchforge-port-adapter alttp-join
```
Both should show Link in bsnes-as; inventory/position sync on group `couchforge`.
