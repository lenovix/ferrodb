version     | description 
-----       | --- 
0.3.2       | Usability (AUTH, Multi-DB, Error message)
0.3.1       | TTL key + PERSIST key
0.3         | Config file + INFO command
0.2.5       | Graceful Shutdown + Signal Handling
0.2.4       | Snapshot / AOF Rewrite
0.2.3       | Persist TTL
0.2.2       | TTL / EXPIRE command
0.2.1       | TCP Server
0.2         | AOF Persistence


How to Run:
go run ./cmd/ferrodb

//in other terminal:
telnet 10.124.39.76 6380


| Kondisi TTL key     | Return     |
| ------------------- | ---------- |
| key tidak ada       | `-2`       |
| key ada & tanpa TTL | `-1`       |
| key ada & ada TTL   | sisa detik |


| Kondisi PERSIST key                   | Return     |
| ------------------------------------- | ---------- |
| TTL berhasil dihapus                  | `1`        |
| key tidak ada/tidak punya TTL         | `0`        |


Glamory:
AOL
TTL
TCP