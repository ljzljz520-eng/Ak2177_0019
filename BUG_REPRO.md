# BUG_REPRO

The following failures were observed while validating the initial project state.
Each section records what failed, how to reproduce it, and the complete command output.
They are preserved intentionally; only failing build gates are omitted from the generated Dockerfile.

## Failure 1: Go test (.)

- Observed problem: `Go test (.)` failed in the initial project state.
- Working directory: `.`
- Command: `cd /app && GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off go test -count=1 ./...`
- Exit status: `1`

```text
?   	inventorychain/cmd/inventoryd	[no test files]
ok  	inventorychain/internal/api	0.011s
--- FAIL: Test2177BusinessRegression (0.00s)
    processor_test.go:15: slot order changed: got [{afternoon 2 bob} {afternoon 2 bob}] want [{morning 1 alice} {afternoon 2 bob}]
FAIL
FAIL	inventorychain/internal/flow015	0.002s
ok  	inventorychain/internal/model	0.002s
ok  	inventorychain/internal/report	0.002s
ok  	inventorychain/internal/review	0.002s
ok  	inventorychain/internal/service	0.017s
ok  	inventorychain/internal/store	0.012s
FAIL
```

## Architecture reproduction

### linux/amd64
- Go toolchain version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/inventoryd): exit `0`
### linux/arm64
- Go toolchain version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/inventoryd): exit `0`
