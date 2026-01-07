FerroDB — Redis-like Key-Value Database with Built-in Web Admin
Redis
KeyDB
Etcd (simplified)

## 📦 Version History

| Version | Description                                                    |
| ------: | -------------------------------------------------------------- |
|   0.5 | HTTP API + Web UI (Alpha)                                                   |
|   0.4.1 | Redis Compatibility Basics                                                   |
|   0.4 | RESP protocol                                                   |
|   0.3.5 | ACL Commands                                                   |
|   0.3.5 | Password hashing (bcrypt) + secure AUTH                        |
|   0.3.4 | AUTH + ACL (role-based)                                        |
|   0.3.3 | Improve Multi-DB                                               |
|   0.3.2 | Usability improvements (AUTH, Multi-DB, better error messages) |
|   0.3.1 | TTL key + `PERSIST` command                                    |
|   0.3.0 | Config file support + `INFO` command                           |
|   0.2.5 | Graceful shutdown + signal handling                            |
|   0.2.4 | Snapshot / AOF rewrite                                         |
|   0.2.3 | Persist TTL                                                    |
|   0.2.2 | TTL / `EXPIRE` command                                         |
|   0.2.1 | TCP server                                                     |
|   0.2.0 | AOF persistence                                                |

---

## How to Run

Jalankan server ferroDB:

```bash
go run ./cmd/ferrodb
```

Buka terminal lain untuk menghubungkan client menggunakan Telnet:

```bash
telnet [your ip] 6380
redis-cli -h [your ip] -p 6380
```

> Pastikan alamat IP dan port sesuai dengan konfigurasi server kamu.

---

## ⏱️ TTL Behavior

### Return Value untuk TTL Key

| Kondisi Key            | Return             |
| ---------------------- | ------------------ |
| Key tidak ada          | `-2`               |
| Key ada & tanpa TTL    | `-1`               |
| Key ada & memiliki TTL | Sisa waktu (detik) |

---

## ♻️ PERSIST Command

### Return Value untuk `PERSIST`

| Kondisi                            | Return |
| ---------------------------------- | ------ |
| TTL berhasil dihapus               | `1`    |
| Key tidak ada / tidak memiliki TTL | `0`    |

---

## 📖 Glossary

- **AOF (Append Only File)**
  Mekanisme persistence dengan mencatat setiap perintah write ke dalam file log.

- **TTL (Time To Live)**
  Waktu hidup sebuah key sebelum otomatis dihapus.

- **TCP**
  Protokol jaringan yang digunakan ferroDB untuk komunikasi client-server.

- **RESP**
