[![progress-banner](https://backend.codecrafters.io/progress/dns-server/0211f9f4-c27a-4443-a9c9-1dfab7d164e9)](https://app.codecrafters.io/users/alex-popov-tech?r=2qF)

# 🌐 DNS Server in Go

> A from-scratch DNS server that parses and builds DNS packets straight from raw
> bytes — no DNS libraries. It answers queries directly or recursively forwards
> them to an upstream resolver.

![Go](https://img.shields.io/badge/Go-1.26-00ADD8?logo=go&logoColor=white)
![Protocol](https://img.shields.io/badge/protocol-DNS%2FUDP-1f6feb)
![Challenge](https://img.shields.io/badge/CodeCrafters-challenge%20completed-2ea44f)

Built as the [CodeCrafters "Build Your Own DNS server"](https://app.codecrafters.io/courses/dns-server/overview)
challenge — completed end to end, including DNS message compression and a
recursive forwarding resolver.

---

## ✨ What it does

- **Speaks DNS over UDP** on `127.0.0.1:2053`, reading and writing raw datagrams.
- **Parses & serializes the full message format** — header, question, and answer
  sections — by hand, with big-endian encoding throughout.
- **Bit-packs the header flags** (`QR`, `OPCODE`, `AA`, `TC`, `RD`, `RA`, `Z`,
  `RCODE`) into the 16-bit flags field and reads them back out.
- **Resolves DNS name compression** — follows `0xC0` pointer labels so compressed
  packets parse correctly.
- **Forwards recursively** — with `--resolver`, splits multi-question packets,
  asks an upstream server one question at a time, and merges the answers back
  into a single response.

## 🧠 The DNS message, by hand

Every packet is decoded and re-encoded from raw bytes:

```
┌──────────────┐
│   Header     │  12 bytes — ID, bit-packed flags, 4 section counts
├──────────────┤
│  Question(s) │  QNAME (length-prefixed labels) + QTYPE + QCLASS
├──────────────┤
│   Answer(s)  │  NAME + TYPE + CLASS + TTL + RDLENGTH + RDATA
└──────────────┘
```

`QNAME` labels may end with a **compression pointer** (top two bits `11`) that
redirects to an earlier offset in the packet — handled recursively during parse.

## 🏗️ Architecture

```
dns-go/
├── app/
│   └── main.go              # UDP listen loop, --resolver flag
└── internal/
    ├── config/              # resolver address (netip.AddrPort)
    ├── message/             # the DNS wire format — parse + serialize
    │   ├── message.go       #   Message.Parse() / Message.Bytes()
    │   ├── header.go        #   12-byte header, bit-packed flag helpers
    │   ├── question.go      #   QNAME labels + compression pointers
    │   └── answer.go        #   resource records (A)
    └── handler/             # query handling
        └── handler.go       #   own responses + recursive forwarding
```

## ✅ Stages

| # | Stage                     | Status |
|---|---------------------------|:------:|
| 1 | Setup UDP server          |   ✅   |
| 2 | Write header section      |   ✅   |
| 3 | Write question section    |   ✅   |
| 4 | Write answer section      |   ✅   |
| 5 | Parse header section      |   ✅   |
| 6 | Parse question section    |   ✅   |
| 7 | Parse compressed packet   |   ✅   |
| 8 | Forwarding server         |   ✅   |

**Challenge completed** 🎉

## 🚀 Run it

Requires Go `1.26`.

```sh
# Build and run the server on 127.0.0.1:2053
./your_program.sh

# Recursive forwarding mode — proxy queries to an upstream resolver
./your_program.sh --resolver 8.8.8.8:53
```

Then query it with `dig`:

```sh
dig @127.0.0.1 -p 2053 example.com
```

## 📚 Concepts explored

Binary protocol parsing · big-endian wire encoding · bit-field manipulation ·
DNS name compression · UDP networking · recursive resolution.
