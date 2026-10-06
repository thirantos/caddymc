# Caddymc
Simple minecraft server multiplexer using caddy.

_not affiliated with [caddy](https://github.com/caddyserver/caddy) or [caddy-l4](https://github.com/mholt/caddy-l4)_

## Compilation using `xcaddy`
```bash
xcaddy build \
    --with github.com/thirantos/caddymc \
    --with github.com/mholt/caddy-l4 
```

## Example
```caddy
{
	layer4 {
		:25565 {
			@a minecraft a.example.com
			@b minecraft b.example.com

			route @a {
				proxy 127.0.0.1:25564
			}
			route @b {
				proxy 127.0.0.1:25566
			}

		}
	}
}
```

## Future
- More matching criteria (protocol version, intent, etc)
- Unknown host / different message when not matched.
- Legacy ping
- Case sensitivity settings
- Other packet handlers
