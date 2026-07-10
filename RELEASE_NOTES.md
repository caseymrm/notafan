## v1.1.0 — signed & notarized

First release signed with a Developer ID certificate and notarized by
Apple. No feature changes — the headline is distribution quality.

### What changed
- **Signed + notarized**: hardened runtime, secure timestamp, stapled
  notarization ticket. Downloads open without the "unidentified
  developer" warning.
- **Throttle notifications now work**: `UNUserNotificationCenter`
  requires a Developer ID signature; previous ad-hoc builds silently
  dropped them.
- Picks up [menuet v2.10.4](https://github.com/caseymrm/menuet/releases/tag/v2.10.4)
  (snapshot mode for the [menuet.app](https://menuet.app/apps/notafan/)
  showcase, hardened-runtime signing in menuet.mk).
