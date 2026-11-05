# OpenID Connect login

Traggo can authenticate users against any OpenID Connect provider (authorization code flow with PKCE, `state` and `nonce`, ID token signature, issuer, audience and expiry are verified).
It is disabled unless `TRAGGO_OIDC_ISSUER` is set; without it nothing changes.
After a successful login traggo issues its regular session cookie (a device named "OIDC login", 30 days), so everything else works as before.

## Configuration

| Variable | Default | Description |
|---|---|---|
| `TRAGGO_OIDC_ISSUER` | _(empty, disabled)_ | Issuer URL, `<issuer>/.well-known/openid-configuration` is used for discovery. |
| `TRAGGO_OIDC_CLIENT_ID` | | Client id. Required. |
| `TRAGGO_OIDC_CLIENT_SECRET` | | Client secret (confidential client). Prefer a secret manager or a file next to the binary over the command line. |
| `TRAGGO_OIDC_REDIRECT_URL` | | Public callback URL, e.g. `https://traggo.example.com/auth/oidc/callback`. Required. The path is used as route, register exactly this URL at the provider. Use `https`, the flow cookie gets the `Secure` flag from the scheme. |
| `TRAGGO_OIDC_SCOPES` | `openid,profile,email` | Comma separated scopes. Must contain `openid`. |
| `TRAGGO_OIDC_USERNAME_CLAIM` | `preferred_username` | ID token claim matched against the traggo user name. |
| `TRAGGO_OIDC_ADMIN_CLAIM` / `TRAGGO_OIDC_ADMIN_VALUE` | _(empty)_ | Optional. A user is admin if the claim equals the value or is a list containing it (e.g. `groups` / `traggo-admins`). The admin flag is then synchronised on every OIDC login, in both directions. Without these the admin flag is never changed by OIDC. Set both or none. |
| `TRAGGO_OIDC_AUTO_REGISTER` | `false` | Create users which do not exist yet. They get a random password, so they can only use SSO until an admin sets one. |
| `TRAGGO_OIDC_DISABLE_PASSWORD_LOGIN` | `false` | Reject the username/password login and hide the form. Requires an issuer. Keep a way to recover (e.g. unset it) if the provider is unavailable. |

The provider is contacted on the first login attempt, traggo starts even if it is unreachable.

## Flow

1. The login page asks `GET /auth/oidc/info` and shows a "Login with SSO" button when enabled.
2. `GET /auth/oidc/login` stores a random `state`, `nonce` and PKCE verifier in a short-lived, HMAC-signed, `HttpOnly`, `SameSite=Lax` cookie and redirects to the provider.
3. The provider redirects to the callback. Traggo checks the state against the cookie, exchanges the code (with the PKCE verifier), verifies the ID token and the nonce.
4. The user is looked up by the username claim; if missing and auto register is on, it is created. Then the session cookie is set and the browser is sent to `/`.

## Security notes

- Users are matched **by name**. Anyone who can obtain a given username at the provider can log in as the local user with that name. Only enable this when the provider is trusted and usernames are unique and not self-chosen.
- The session cookie is `HttpOnly`, `SameSite=Lax` and `Secure` when the request came in over https (`X-Forwarded-Proto: https` behind a proxy is honoured, so the proxy must set or strip that header).
- Secrets are masked in the debug log of the configuration.

## Example (generic provider)

```
TRAGGO_OIDC_ISSUER=https://idp.example.com/application/o/traggo/
TRAGGO_OIDC_CLIENT_ID=traggo
TRAGGO_OIDC_CLIENT_SECRET=<from your secret manager>
TRAGGO_OIDC_REDIRECT_URL=https://traggo.example.com/auth/oidc/callback
```

Register the redirect URL as exact match and allow the scopes `openid profile email` (plus the group claim's scope if you use the admin mapping).
When developing with the UI dev server, open the backend port directly for the SSO redirect, the dev proxy does not forward browser navigations.
