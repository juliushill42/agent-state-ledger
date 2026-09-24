#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")"
command -v go >/dev/null || { echo 'ERROR: Go 1.22+ required'; exit 1; }
gofmt -w cmd internal
go test ./...
go vet ./...
mkdir -p bin proof
go build -trimpath -ldflags='-s -w' -o bin/agent-state-ledger ./cmd/agent-state-ledger
rm -f proof/example.jsonl
bin/agent-state-ledger append proof/example.jsonl planner start '{"task":"release"}' >/dev/null
bin/agent-state-ledger append proof/example.jsonl builder complete '{"artifact":"bin/app"}' >/dev/null
bin/agent-state-ledger verify proof/example.jsonl > proof/verify.txt
cp proof/example.jsonl proof/tampered.jsonl
python3 - <<'PY'
from pathlib import Path
p=Path('proof/tampered.jsonl'); s=p.read_text(); p.write_text(s.replace('release','releasX',1))
PY
if bin/agent-state-ledger verify proof/tampered.jsonl >/dev/null 2>&1; then echo 'ERROR: tamper accepted'; exit 1; fi
sha256sum bin/agent-state-ledger proof/example.jsonl proof/verify.txt > proof/SHA256SUMS.txt
echo 'AGENT-STATE-LEDGER READY'
echo '":"'
