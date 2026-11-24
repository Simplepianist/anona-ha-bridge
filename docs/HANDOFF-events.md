# Handoff: detection events never reach Home Assistant

Branch: `debug-push-frames` on Simplepianist/anona-ha-bridge (fork of Darkhand81/anona-ha-bridge). Goal: get person/animal/motion detections into MQTT (`anona/anona_<id>/event/<type>`), then open a PR upstream.

## Setup
Raspberry Pi 4 (4 GB), Docker, host networking. Two HM1010 cameras (Küchenkamera, Wohnzimmer). HA, Mosquitto and the bridge run on the Pi. Streams are H.265 3840x2160 + Opus and work. Image build repo: Simplepianist/anona-bridge (GitHub Actions, GHCR, private). Go cannot download modules in the agent sandbox, so compile checks happen in CI.

## Finding so far
With `ANONA_PUSH_DEBUG=1` (commit d101dbc) the control WebSocket delivered only 614-byte status frames, one per camera about every 5 minutes (hex begins `08321003A203AB02...`). No other frame, and no `drop:` line, appeared while the Anona app showed detection borders and in-app notifications. Conclusion: the cloud does not send `REPORT_IPCAMERA_EVENT` on this WebSocket. The decoder is not the problem. Caveat: the capture window and the real detection time were not matched exactly.

## Hypothesis to test
The app draws its detection border from data inside its own live session, probably a WebRTC data channel (the code uses KVS-style names such as `KvsAudioTrack`, so a channel like `kvsDataChannel` is plausible). `internal/media/session.go` builds an offer with one video and one audio transceiver and no data channel, so the camera cannot open one. Unverified.

## Task 1: data channel probe (default behavior must not change)
1. In `media.Start`, when env `ANONA_DC_PROBE=1`, call `pc.CreateDataChannel("kvsDataChannel", nil)` before `CreateOffer`. Register `OnOpen`, `OnClose`, `OnMessage` and `pc.OnDataChannel` handlers. Log `[dc] open label=...`, and for each message log `[dc] msg label=%s len=%d hex=%.400x` (and the text if valid UTF-8).
2. Log every `m=` line of the answer SDP when probing is on, so we learn whether the camera accepted `m=application`.
3. Unset env means identical SDP and behavior to today. If the camera rejects the offer with the extra section, report that and stop: it is a result.
4. If messages arrive: decode them, find the detection fields, and wire them into the existing event path (`cam.ParseEventTags` / supervisor `OnEvent` / `ha` MQTT). Add a small table-driven test with captured bytes.

## Task 2: if no data channel exists
Report that and propose options: (a) forward the Anona phone app notification to HA through the HA companion app notification listener, (b) video-based detection (not viable on the Pi).

## Rules
No credentials, tokens or captured user ids in commits or logs. Work on branches, never push to main. Do not publish images publicly: upstream has no license. Keep diffs small and match existing style. Build check: run the `build` workflow in Simplepianist/anona-bridge with `ref=<branch>` (`gh workflow run build.yml -R Simplepianist/anona-bridge -f ref=<branch>`), then image tag is `up-<sha12>`.

## Test on the Pi (user runs these)
Compose env: `ANONA_PUSH_DEBUG: "1"` and `ANONA_DC_PROBE: "1"`, image `ghcr.io/simplepianist/anona-bridge:up-<sha12>`. Then: `docker logs -f --since 1m anona-bridge 2>&1 | grep --line-buffered -E '\[(dc|ctl)\]'`. Walk in front of the Küche camera until the app shows the detection border, note the clock time, wait 20 s, paste the lines.
