# ltsv-analyzer

ltsv-analyzer makes log counter based on the selected key.
if target key had set and target key is numberic, can get sum or average too.

## Usage

```
$ ./ltsv-analyzer -h                                                                                    
How to use web-proxy log analyzer
ltsv-analyzer [OPTIONS] {PATH1} {PATH2} ... (default path : ./logs)

- OPTIONS
    --base[-B]      : set base key
    --target[-T]    : set target key
    --operation[-O] : set operation (sum, avg, cnt[count], all)
    --sort[-S]      : set sort key  (sum, avg, cnt[count])
    --rev[-r|-R]    : reverse sort direction (default: DESCending)

  --help[-h|H]      : show usage
```

#### example


- default. want to get log counter by each hosts

```
❯ ./ltsv-analyzer -B host ./testlog.log

Print LOG COUNTER by BASE KEY [host]
host            LOG COUNTER
---------------------------
bbb.com                   3
aaa.com                   2
```

- want to get total / average bytes_sent and counter by each client ips

```
# ./ltsv-analyzer -B ip -T bytes_sent -O all ./testlog.log

Print results by BASE KEY [ip] and TARGET KEY bytes_sent
 ip         SUM(bytes_sent)     AVG(bytes_sent)         LOG COUNTER
-------------------------------------------------------------------
1.1.1.1               12656                6328                   2
1.1.1.2               11729                5864                   2
1.1.1.4             1123419             1123419                   1
```

- want to get total & counter bytes_sent by each user

```
$ ./ltsv-analyzer -B user -T bytes_sent -O sum,cnt ./testlog.log

Print results by BASE KEY [user] and TARGET KEY [bytes_sent]
user     SUM(bytes_sent)         LOG COUNTER
--------------------------------------------
bbb             1135148                   3
aaa               12656                   2
```
