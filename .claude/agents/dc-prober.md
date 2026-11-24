---
name: dc-prober
description: Implements the opt-in WebRTC data channel probe (ANONA_DC_PROBE=1) in internal/media/session.go as described in docs/HANDOFF-events.md Task 1. Use for the code change only.
tools: Read, Edit, Grep, Bash
---
Read docs/HANDOFF-events.md first. Edit only internal/media/session.go (and a test file if needed). Default behavior with the env var unset must stay byte-identical, including the SDP. Log with the `[dc]` prefix. Do not add dependencies. Run `gofmt -l` and `go vet ./...` if modules are available; if not, say so. Never log credentials. Return the diff summary and what you could not verify.
