# Connect Simkl in a few steps

This guide describes the Windows preview build. Sign-in, lists and settings appear directly inside the tray panel.

## Before you start

Reading Custom Lists requires a Simkl PRO or VIP account. Named custom collections are separate from built-in status categories such as Watching and Plan to Watch. [Custom Lists](https://api.simkl.org/guides/custom-lists)

## Free PRO without a paid subscription

PRO earned through Simkl's Loyalty Program also qualifies:

1. While signed in, use the Simkl website or Simkl Lists app on at least **20 days during the month**.
2. Open the [Loyalty page](https://simkl.com/vip/loyalty/) and check your progress.
3. When available, click **Claim Free PRO** to activate one month of PRO.
4. Claim it again each month; renewal is not automatic.

No paid subscription is required. The bridge does not collect activity days automatically. Check [Simkl's current terms](https://docs.simkl.org/how-to-use-simkl/account-and-billing/special-plans/offers/loyalty-program).

## 1. Register a V2 app

1. Sign in and open [Simkl Settings → Developer](https://simkl.com/settings/developer/).
2. Create a new **AUTH V2** app.
3. Select **TV, devices & command line** for this bridge. It uses code-based sign-in, including on Windows.
4. Use a name such as **List Bridge for Simkl**. Optional fields may be left blank.
5. Review, create the app, and copy its **Client ID**.

This client type needs no client secret or redirect URL. One registration covers movie, TV and anime lists. [Authentication setup](https://api.simkl.org/api-reference/auth)

## 2. Connect your account

Under Sign-in, enter the Client ID and choose Connect to Simkl. Open the provided link in your browser and approve read-only access. The bridge handles token exchange and local storage; you do not need to copy tokens or run commands.

The Client ID identifies the app. Your browser approval connects your account. Access and refresh tokens remain private. [AUTH V2](https://api.simkl.org/api-reference/auth-v2)

## 3. Choose your lists

Select a list and its Sonarr or Radarr destination. Add the local feed URL provided by the bridge as a Custom List in the destination app. A Simkl list link or numeric ID can also be pasted directly.

Check the destination app's library-cleaning setting before saving: it determines what happens to titles absent from its lists. The bridge does not change this setting itself.

## Why AUTH V1 does not work here

AUTH V1 is the old authentication system. Custom Lists already require AUTH V2. An existing V1 Client ID cannot be upgraded: register a separate V2 app and authorize it again. [Custom Lists requirements](https://api.simkl.org/guides/custom-lists), [Migration guide](https://api.simkl.org/guides/migrating-v1-to-v2)

**Checked on 3 November 2026:** Simkl expects AUTH V1 to stop working around April 2027. The exact date has not been fixed. Check the official announcement for later changes. [AUTH V1 retirement](https://api.simkl.org/api-reference/auth-v1)

## Troubleshooting

- **invalid_client:** Check that your Client ID belongs to an AUTH V2 app.
- **An old V1 connection still works elsewhere:** Custom Lists still require V2.
- **PRO/VIP required:** Check the account used when approving access.
- **Empty list:** Check its contents and confirm that you selected the intended custom collection.

Use Back to return to the main menu. Browser authorization remains separate. Personal credentials are not included in source or installation packages.

