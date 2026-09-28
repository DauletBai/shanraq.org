# Networking without magic: IP, port, DNS, HTTP, and TLS

_Lead (summary):_ **Trace a request from a domain name to a process socket and localize failures layer by layer.**

## Lesson outcome

Trace a request from a domain name to a process socket and localize failures layer by layer. Do not move on until you can explain what every command verifies and which failure it can reveal.

## Why this matters

DNS answers which address, TCP whether a port accepts a connection, HTTP whether the application understood a request, and TLS whom the secure connection reached. Test layers in order so a certificate failure does not become a vague “server is down.”

## Practice

Run the commands from the `course/cloud-devops-lab` root. The commands target a learning environment. Read a command in full and check the current directory before it changes a server or data.

```shell
getent hosts example.com 2>/dev/null || nslookup example.com
curl -sS -o /dev/null -w 'status=%{http_code} ip=%{remote_ip} tls=%{ssl_verify_result}\n' https://example.com
printf 'GET /healthz HTTP/1.1\r\nHost: localhost\r\nConnection: close\r\n\r\n' | nc 127.0.0.1 8080 || true
```

[Project reference files](https://github.com/DauletBai/shanraq.org/tree/main/course/cloud-devops-lab).

## Read the result

- DNS returns an address independently of HTTP.
- curl shows status, remote IP, and TLS verification.
- You can explain why port 443 is not itself HTTPS.

![Support map 6](/static/course/cloud-devops/map-06-networks-dns-http-en.svg)

## Exercise

**Required.** Write a POSIX sh diagnostic that accepts a URL, prints its HTTP status and remote IP with `curl`, and fails for status 400 or higher. Submit only POSIX sh commands or a script, with no real keys or tokens.

**Optional.** Write down one possible failure and the signal that would reveal it.

[Course contents](/course/cloud-devops?lang=en)
