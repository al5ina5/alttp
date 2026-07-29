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

## Snapshot / Request Resync
The **Request Resync** button needs our forked `alttp-server` (kinds `0x03`/`0x04`).
Public `alttp.online` will ignore those — enemy/item sync still works over the relay;
only the manual snapshot button needs a local server on aio.

