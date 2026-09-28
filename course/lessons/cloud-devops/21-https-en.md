# Domain, reverse proxy, and automatic HTTPS

_Lead (summary):_ **Point DNS at the server, serve CloudLab through Caddy, and verify the TLS chain externally.**

## Lesson outcome

Point DNS at the server, serve CloudLab through Caddy, and verify the TLS chain externally. Do not move on until you can explain what every command verifies and which failure it can reveal.

## Why this matters

A reverse proxy accepts the public connection, terminates TLS, and forwards the request to the private application. Automatic certificates require the domain to resolve to the right IP and ports 80 and 443 to be reachable. Do not expose the application port publicly without need.

## Practice

Run the commands from the `course/cloud-devops-lab` root. The commands target a learning environment. Read a command in full and check the current directory before it changes a server or data.

```shell
export CLOUDLAB_DOMAIN=cloudlab.example.com
getent hosts "$CLOUDLAB_DOMAIN" 2>/dev/null || true
printf 'Caddy route:\n'
sed -n '1,80p' deploy/caddy/Caddyfile
printf 'After real DNS: curl -v https://%s/readyz\n' "$CLOUDLAB_DOMAIN"
```

[Project reference files](https://github.com/DauletBai/shanraq.org/tree/main/course/cloud-devops-lab).

## Read the result

- The app is reached through the proxy, not a public port 8080.
- The command exits with a predictable status.
- The result can be repeated in a clean environment.

![Support map 21](/static/course/cloud-devops/map-21-https-en.svg)

## Exercise

**Required.** Write a POSIX sh check for an HTTPS URL that requires valid certificate verification and HTTP 200 from `/readyz`, then prints the certificate expiry without `-k`. Submit only POSIX sh commands or a script, with no real keys or tokens.

**Optional.** Write down one possible failure and the signal that would reveal it.

[Course contents](/course/cloud-devops?lang=en)
