# agent-state-ledger `":"`

An append-only JSONL state ledger for agents and automation.

Every entry contains its own SHA-256 hash and the hash of the previous entry. `verify` walks the entire chain and fails on edits, deletion/reordering that breaks sequence, a changed watermark, or a bad hash.

```bash
./bootstrap.sh
bin/agent-state-ledger append state.jsonl planner started '{"task":"build"}'
bin/agent-state-ledger verify state.jsonl
```

The proof gate deliberately tampers with a copied ledger and requires verification to reject it.

Built by Titan Universal AI LLC.

`":"`
