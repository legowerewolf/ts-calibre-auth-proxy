# Authenticating Reverse Proxy for Calibre

Provide access to your Calibre ebook library as a separate node on your Tailnet,
including authorization based on Tailscale user identity.

I make no claims about the security of this setup; use this at your own risk.

## Setup

### Calibre

You'll find all these settings under Preferences > Sharing over the net.

#### Main tab

- Enable "Require username and password to access the Content Server"
- Make a note of the port number if it's not 8080.

#### Advanced tab

- Set "Choose the type of authentication used" to "basic"
- Set "Number of login failures for ban" to "0"

#### User accounts tab

Create a user account for each Tailscale user you want to grant access to your
server.

For the username: take their Tailscale login-name (likely their email address,
you'll be able to see this in logs) and replace all characters that aren't
letters, numbers, spaces, hyphens, or underscores with underscores.

Set the password to `tailscale-authenticated`

### Proxy

Grab an auth key from the Tailscale admin console, and set it as the environment
variable `TS_AUTHKEY`.

Run the proxy with `ts-calibre-auth-proxy` and it'll connect to your Tailnet and
start providing access to your Calibre server.

If the port number of your Calibre server wasn't 8080, add the flag
`--origin http://localhost:PORT` to the command.

If you want to provide a different name for the Tailscale node, add the flag
`--hostname NAME` to the command.

## Future work (?)

This could probably be modified or extended fairly easily to work with other
origin servers and other authentication methods. I don't have a need to; this is
good enough for my purposes.
