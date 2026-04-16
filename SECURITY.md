# Security Policy

## Reporting a Vulnerability

If you discover a security vulnerability in PhishGuard, please report it responsibly.

**Do NOT open a public GitHub issue for security vulnerabilities.**

### How to Report

- Email: [im@mxp.tw](mailto:im@mxp.tw)
- Or use [GitHub Security Advisories](https://github.com/nczz/phishguard/security/advisories/new) to report privately.

Please include:

- Description of the vulnerability
- Steps to reproduce
- Potential impact
- Suggested fix (if any)

### Response Timeline

| Stage | Target |
|-------|--------|
| Acknowledgment | Within 48 hours |
| Initial assessment | Within 7 days |
| Fix release | Within 30 days (critical: ASAP) |

### Disclosure Policy

We follow **coordinated disclosure**. Please allow us reasonable time to address the issue before any public disclosure. We will credit reporters in the release notes unless anonymity is requested.

## Supported Versions

| Version | Supported |
|---------|-----------|
| Latest  | ✅        |
| Older   | ❌        |

## Security Best Practices for Deployment

- Always set strong, unique values for `DB_PASS`, `JWT_SECRET`, and `ADMIN_PASSWORD`
- Use HTTPS in production (both app and tracker domains)
- Keep Docker images and dependencies up to date
- Restrict database access to internal networks only
