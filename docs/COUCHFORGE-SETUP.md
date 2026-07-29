# ALttPO on Couchforge

## Goal
- **aio** (`ALttPO (Host)`): auto-joins as player AIO
- **kiosk** (`ALttPO (Join)`): auto-joins as player Kiosk
- Same group (`couchforge`), no Join window / no IP prompts

## Why public relay by default
`ufw` on these machines blocks LAN UDP (including `:4590`). Until you run
`sudo ufw allow 4590/udp` on **aio**, both clients use the public ALttPO relay
`alttp.online` with group `couchforge` so they still auto-find each other.

For true LAN later:
```bash
# on aio
sudo ufw allow 4590/udp
# then launch with:
COUCHFORGE_ALTTP_LOCAL_SERVER=1 COUCHFORGE_ALTTP_SERVER=127.0.0.1  # host
COUCHFORGE_ALTTP_SERVER=10.0.0.55                                  # join
```

## Install layout
`~/alttp-test/`:
- `bin/bsnes` — bsnes-as
- `bin/alttp-server` — optional local UDP group server
- `bin/discord_game_sdk.so`
- `alttpo/` — AngelScript scripts (must compile; no `#include`)
- `roms/alttp_usa.sfc` — legal ROM

## Launchers
`port-alttp-host` / `port-alttp-join` → `~/.local/bin/couchforge-port-adapter alttp-host|alttp-join`

Port markers: **ALttPO (Host)** / **ALttPO (Join)** only (Archipelago multiworld removed).

## Verify
```bash
# aio
~/.local/bin/couchforge-port-adapter alttp-host
# kiosk
~/.local/bin/couchforge-port-adapter alttp-join
```
Logs should show `Kiosk joined` / `AIO joined`. No "Join a Game" window.
Input: DualShock 4 (`0x054c09cc`) via `udev` + keyboard fallback.
Display: adapter forces `1920x1080` when xrandr is stuck on 8x8/320x200.
