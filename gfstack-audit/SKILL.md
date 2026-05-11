---
name: gfstack-audit
description: "Go security and code audit skill. Systematic review checklist covering input validation, sensitive data, concurrency, panics, file operations, and stability. TRIGGER: code review, security audit, vulnerability check, code inspection, security review, audit code, review code, security scan, bug check. DO NOT TRIGGER: writing new code (see gfstack-*), architecture design (see gfstack-overview)."
---

# gfstack-audit

## Audit Checklist

### 1. Input & Validation

| Issue | Risk | Mitigation |
|-------|------|------------|
| Insufficient input validation | High | Strictly validate all user input |
| SQL injection | High | Use parameterized queries / ORM safe APIs |
| Command injection | High | Avoid string concatenation for commands; use `exec.Command` safe API |
| XSS | High | Escape HTML/JS output |
| CSRF | Medium | Use token verification |
| Path traversal | High | Sanitize paths, restrict upload directories |

### 2. Sensitive Data

| Issue | Risk | Mitigation |
|-------|------|------------|
| Log leakage of sensitive info | High | Never log passwords/keys; mask sensitive data |
| Hardcoded secrets | High | Store in config files or secure vaults |

### 3. Dependencies & Configuration

| Issue | Risk | Mitigation |
|-------|------|------------|
| Vulnerable dependencies | High | Regularly update dependencies |
| TLS/HTTPS misconfiguration | High | Properly use certificate verification and secure cipher suites |

### 4. Concurrency & Data Races

| Issue | Risk | Mitigation |
|-------|------|------------|
| Concurrency safety | High | Use mutex locks or `sync.Map` |
| Race conditions causing panics/logic errors | High | Use race detector (`go test -race`), add locks |
| Goroutine leak | Medium | Ensure goroutines can exit via context/channel |

### 5. Panic / Crash

| Issue | Risk | Mitigation |
|-------|------|------------|
| Nil pointer dereference | High | Nil checks before dereference |
| Array/slice out of bounds | High | Length checks before indexing |
| Type assertion failure | High | Use `v, ok := x.(T)` pattern |
| Division by zero | Medium | Check divisor before division |
| Channel operation panics | High | Check channel state before send/receive |
| Improper defer + panic handling | High | Full `recover()` coverage in goroutines |
| Unchecked function return values | Medium | Check returned nil/error values |
| JSON/deserialization errors | Medium | Validate struct and field types |
| Incorrect access control | High | Strict permission checks |

### 6. File & System

| Issue | Risk | Mitigation |
|-------|------|------------|
| Unsafe file operations | High | Restrict paths, types, and sizes |
| Log injection | Medium | Escape or filter user input in log entries |
| cgo security risks | High | Be careful with raw pointers and memory operations |

### 7. Environment & Configuration

| Issue | Risk | Mitigation |
|-------|------|------------|
| Environment variable leakage | High | Never output sensitive environment variables |
| Trusting external configuration | High | Validate config values before use |

### 8. Stability

| Issue | Risk | Mitigation |
|-------|------|------------|
| Infinite loops or blocking | High | Use timeouts or context cancellation |
| Incorrect time handling / timezone | Medium | Validate time values and timezone |
| Unbounded buffer size | High | Limit slice/map/channel capacity |
| Regex ReDoS | High | Use safe regex patterns, limit input length |
| Memory leak | High | Avoid excessive goroutines or unreleased resources |
| Unsafe reflection | Medium | Validate types and nil values |
| JSON/HTML template injection | High | Escape template output |
| Improper signal handling | Medium | Properly handle SIGTERM/SIGINT |
