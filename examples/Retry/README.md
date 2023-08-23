
# Retry Example

This file has commands that fail on the first run.

Generated tests will still succeed because commands are retried until timeout.

## Run

```bash
rm -f retry-file-flag
```

```bash
[ -f retry-file-flag ] || (
cat << EOF > retry-file-flag
this file must exist before command runs
EOF
false
)
```

## Cleanup

```bash
rm retry-file-flag
```
