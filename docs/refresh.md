# How refreshing works

<strong><a href="refresh.md"><img src="https://github.githubassets.com/images/icons/emoji/unicode/1f1ec-1f1e7.png" width="20" height="20" alt=""> English</a></strong> · <a href="refresh.de.md"><img src="https://github.githubassets.com/images/icons/emoji/unicode/1f1e9-1f1ea.png" width="20" height="20" alt=""> Deutsch</a>

[README](../README.md) · [Setup guide](setup-simkl.en.md)

## In everyday use

Sonarr or Radarr requests an import list on its own schedule. The first request reads that Simkl Custom List completely. Later requests reuse cached contents when no relevant change is found. There is no constant background polling of Simkl.

The bridge always returns a **complete list**, not just newly added titles. That is the format the destination apps need. Checking for changes reduces unnecessary downloads from Simkl; it does not change the list format sent to Sonarr or Radarr.

## What triggers a new download

| List | How it is checked |
|---|---|
| Your own regular Custom List | Check Custom List activity, then refresh the owner's list overview when needed. Compare the list's update timestamp before downloading its contents. |
| Another user's list | Check that owner's overview; their changes are not reflected in your account's activity timestamp. |
| Automatic list | Recheck on demand, including a daily content refresh while its metadata stays unchanged. Missing metadata can require another content request. |

A short local cache avoids repeated checks within the same minute. Concurrent requests share the content download, and the saved cache is reused after a restart. The cache is bound to the recognized Simkl account.

## Errors and removals

Missing lists, missing required counters, inconsistent media types and incomplete API responses produce an error. They do not become a successful empty list. A saved list is replaced only after a complete, successful download. A genuinely empty Simkl list is still returned as empty.

## Custom Lists and built-in watchlists

Custom Lists are separate from Simkl's built-in status watchlists. The `/sync/all-items?date_from=` route handles changes to the built-in watchlists; it is not a promised item-by-item change feed for Custom Lists. A changed Custom List may need to be downloaded completely again.

Local status messages between the app and its background service do not make additional Simkl API calls.

[Official sync guide](https://api.simkl.org/guides/sync) · [Official Custom Lists guide](https://api.simkl.org/guides/custom-lists)
