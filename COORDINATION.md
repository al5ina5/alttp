# Coordination: Upstream Fork References

**Date:** 2026-07-28
**Author:** Hermes session (forks/issues audit)

## What I found

I audited all forks of `alttpo/alttpo` and `alttpo/bsnes-as` for unmerged work that we can use.

## Gold Mine: mysterypaintwo `esync` branch

**Repo:** https://github.com/mysterypaintwo/alttpo/tree/esync
**Ahead of unstable:** 89 commits (1 behind)

Contains **working implementations of things we're planning to build:**

| Feature | In esync | We need? |
|---------|----------|----------|
| Enemy sync (ALTTP) | ✅ `SM_Enemy.as`, `EnemyWindow.as`, changes in `LocalGameState.as` | Yes |
| Enemy sync (Super Metroid) | ✅ Same files above | Yes |
| Boss room sync | ✅ Room entry delay + boss death flags | Yes |
| PvP damage balance | ✅ Nerfs for hammer, silver arrows | Nice-to-have |
| Custom sprite sync | ✅ Unique tile capture + NAK/rebroadcast in `Sprite.as` | Yes |
| Players window | ✅ Shows who's connected | Nice-to-have |
| Memory window enhancements | ✅ Darker zero bytes | Minor |
| Crystal switch fix | ✅ Non-owning players not stuck | Yes |
| Room entry delay fix | ✅ Delays enemy sync on room enter | Yes |
| SMZ3 enemy sync | ✅ Custom sprite support | If we support SMZ3 |

**Also has branches:**
- `player-names` (4 commits) — player labels under sprites
- `sprites` (7 commits) — reliable ACK-based protocol
- `unstable-canvas` (40 commits) — scalable map window

## Planned Work: I'm going to reference this code

**I will NOT modify any files you're currently working on.** I'll:
1. Pull the esync branch locally as `refs/esync`
2. Read key implementation files
3. Write reference docs with code snippets showing how enemy sync, sprite sync, and PvP are done
4. Mark which parts have merge conflicts with our current code

**Conflict zones if you were to merge esync:**
- `LocalGameState.as` (we both modify this heavily)
- `pre_frame.as` (same — sync lifecycle hooks)
- `GameState.as` (serialization)
- `Sprite.as` + `SpritesWindow.as` (extensions)
- `SettingsWindow.as` (new checkboxes)
- `init.as` (initialization)
- `Projectile.as` (PvP damage calc)

**Recommendation:** Don't merge blindly. Read the individual features and hand-port them. I'll document the approach for each so we can cherry-pick.

## Not Worth It

- `bsnes-as` PR #3 (macOS-only bug fixes — irrelevant on Linux)
- Old stale branches from 2020-2021 (seriously behind unstable)

---

*Next session: read this file first before starting work.*
