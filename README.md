# bunyan

Logviewer for Bunyan-based logs

Describe the flags:

- `--level` or `-l` - filter logs by level see github.com/gildas/go-logger for the levels
- `--filter`:
  Examples: `--filter '.field == value'`, `--filter '.field1 == .field2'`, `--filter '.field =~ /regexp/'`
  `--filter '.field1 == value && .field2 == value2'`, alse `||` and `!` are supported
