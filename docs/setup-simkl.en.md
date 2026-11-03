# Set up Simkl Bridge

<strong><a href="setup-simkl.en.md"><img src="https://github.githubassets.com/images/icons/emoji/unicode/1f1ec-1f1e7.png" width="20" height="20" alt=""> English</a></strong> · <a href="setup-simkl.de.md"><img src="https://github.githubassets.com/images/icons/emoji/unicode/1f1e9-1f1ea.png" width="20" height="20" alt=""> Deutsch</a>

[README](../README.md) · [Installation](installation.en.md) · [Troubleshooting](#troubleshooting)

Follow these steps once for your first list. You can return to **Settings** later to add lists or change the port. Already have a valid AUTH V2 Client ID? Skip to [Connect your account](#2-connect-your-account).

The app initially opens in German, regardless of the installer language. Choose **Zurück** to leave the settings page, then **English** at the bottom of the main menu. Reopen **Settings → Sign-in**; the steps below use the English labels.

## Before you start

Have Sonarr and/or Radarr installed on the same Windows computer. You also need a Simkl account with **PRO or VIP** and a Custom List containing the titles you want to import. Free Loyalty PRO also qualifies.

A **Custom List** is a collection you name yourself on Simkl. It is separate from built-in status categories such as **Watching** and **Plan to Watch**. Adding a title to a status category does not automatically add it to a Custom List. [Simkl's Custom Lists guide](https://api.simkl.org/guides/custom-lists)

## Free PRO without a paid subscription

1. While signed in, use the Simkl website or Simkl Lists app on at least **20 days during the month**.
2. Open the [Loyalty page](https://simkl.com/vip/loyalty/) and check your progress.
3. When available, choose **Claim Free PRO** to activate one month of PRO.
4. Claim it again each month; renewal is not automatic.

No paid subscription is required for this route. The bridge does not collect activity days automatically. [Simkl's current terms](https://docs.simkl.org/how-to-use-simkl/account-and-billing/special-plans/offers/loyalty-program)

## 1. Register a Simkl AUTH V2 app

The registration gives your copy of the bridge a **Client ID**. It does not create another Simkl account.

1. Sign in to Simkl and open [Settings → Developer](https://simkl.com/settings/developer/).
2. Create a new app using **AUTH V2**.
3. Select **TV, devices & command line**. This option also covers scripts and Windows apps that use code-based sign-in.
4. Enter a name such as **List Bridge for Simkl**. Optional description, homepage and contact fields may be left blank.
5. Review the registration and create the app.
6. Copy the **Client ID**. You will paste it into the bridge in the next step.

One registration covers your movie, TV and anime lists. This app type needs **no Client Secret and no redirect URL**. An unverified personal app can be used. [Simkl authentication documentation](https://api.simkl.org/api-reference/auth)

## 2. Connect your account

1. Open Simkl Bridge. If the main menu is shown, choose **Settings → Sign-in**.
2. Paste the Client ID into **Your AUTH V2 app client ID**.
3. Choose **Connect to Simkl**. The panel shows a short sign-in code.
4. Choose **Open browser**. Confirm that you are signed in to the intended Simkl account, enter the code if requested, and approve read-only access.
5. Wait for the bridge to report that sign-in succeeded. Keep the settings panel open while waiting.

The bridge saves the sign-in securely in your Windows account. You do not need to copy access tokens or refresh tokens. Your Client ID identifies the app; browser approval connects your account.

## 3. Choose a list and copy its address

1. Open **Lists** and choose **Load my lists**.
2. Select a list, then choose **Selection → Sonarr** for series or **Selection → Radarr** for movies. Alternatively, paste a Simkl list link or list ID and use the matching **Link → …** button.
3. Select the added entry in the lower list and choose **Copy address**.
4. Choose **Save and back**. If the bridge is stopped, choose **Start bridge**.

The address on the clipboard is the one you will use below. Save and leave the settings page before testing it in Sonarr or Radarr: the bridge pauses its service while settings are being edited. Repeat these steps for each list and destination.

For an anime list containing both movies and other anime entries, add a Sonarr mapping and a Radarr mapping. Each address returns the matching entries; required TVDB/TMDb identifiers must be available. You can also add other users' accessible lists by link or ID.

To stop sharing a configured list, select its entry and choose **Remove mapping → Save and back**. Remove or disable its import list in Sonarr/Radarr as appropriate. Removing the bridge mapping does not delete anything on Simkl or in the destination library. List and port changes are saved with **Save and back**; **Back** leaves those unsaved changes behind.

## 4. Add the import list to Sonarr or Radarr

In the destination app:

1. Open **Settings → Import Lists** and choose **+**.
2. Select **Custom List**.
3. Give the import list a recognizable name and paste the copied bridge address into its URL field.
4. Select the root folder and quality profile you want the destination app to use. Decide whether new titles should be monitored and searched for automatically.
5. Review the destination app's library-cleaning setting before the first import. It determines how titles removed from lists are handled; the bridge does not change this setting.
6. Choose **Test**, then **Save**. Use the destination app's list refresh action or wait for its scheduled refresh.

Sonarr and Radarr use their own import schedules. The bridge provides the lists; it does not control download quality or start imports through an administrator API key. Both apps must be able to reach the bridge on the same computer.

## 5. Enable Windows startup

Right-click the S icon next to the clock and enable **Start with Windows**. You can use **Stop bridge** without closing the tray app, or **Quit** to close both. The next launch reuses your saved sign-in and lists.

## Why AUTH V1 does not work here

Custom Lists require **AUTH V2**. A V1 Client ID cannot be converted: register a separate V2 app and reconnect with its Client ID. [Migration guide](https://api.simkl.org/guides/migrating-v1-to-v2)

The documented check on **3 November 2026** recorded Simkl's estimated AUTH V1 retirement as **around April 2027**, without a fixed date. Check the [official AUTH V1 notice](https://api.simkl.org/api-reference/auth-v1) for updates.

## Troubleshooting

| What you see | What to check |
|---|---|
| `invalid_client` or a Client ID error | Use the complete Client ID of an **AUTH V2** registration. |
| Sign-in required | Open **Settings → Sign-in** and reconnect. Approve the intended account in the browser. |
| PRO/VIP required | Check the approved account's current status. Free PRO must be claimed and renewed manually. |
| An empty list | Open the Custom List on Simkl and check its contents. A built-in watchlist status is a different list. |
| Test fails in Sonarr/Radarr | Leave the bridge's settings with **Save and back**, start the bridge, and use **Copy address** again if the port changed. |
| Port already in use | Another app or an older bridge may use port 8787. Under **Settings → General**, choose an unused port, save, and copy the updated addresses. |
| Wrong destination or missing identifiers | Match series to Sonarr and movies to Radarr. Choose **Open log** to check for missing TVDB/TMDb identifiers. |

See [Installation](installation.en.md) for updates and uninstall, or [How refreshing works](refresh.md) for caching details.

If these checks do not resolve the problem, follow [Reporting a problem](../README.md#troubleshooting). Include the relevant log excerpt and steps to reproduce it in your issue; remove private information first. The log currently uses German messages regardless of the interface language.
