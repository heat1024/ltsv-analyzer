# ltsv-analyzer

ltsv-analyzer makes log counter based on the selected key.
if target key had set and target key is numberic, can get sum or average too.

## Usage

```Shell
$ ./ltsv-analyzer -h                                                                                    
How to use web-proxy log analyzer
ltsv-analyzer [OPTIONS] {PATH1} {PATH2} ... (default path : ./logs)

if [PATH] not defined, use stdin.

- OPTIONS
    --base[-B]      : set base key
    --target[-T]    : set target key
    --operation[-O] : set operation (sum, avg, cnt[count], all)
    --sort[-S]      : set sort key  (sum, avg, cnt[count])
    --rev[-r|-R]    : reverse sort direction (default: DESCending)

  --help[-h|H]      : show usage
```

## example

### default operation

- want to get log counter by each hosts

```Shell
❯ ./ltsv-analyzer -B host ./testlog.log

Print LOG COUNTER by BASE KEY [host]
host            LOG COUNTER
---------------------------
bbb.com                   3
aaa.com                   2
```

- if want analize from multi files, just add log file paths (can use gzip files also)

```Shell
❯ ./ltsv-analyzer -B host ./testlog.log ./testlog.log.1 /some/paths/logfile.log /some/paths/logfile.log.1.gz /some/paths/logfile.log.2.gz
```

- can use `*` to file path

```Shell
❯ ./ltsv-analyzer -B host ./testlog.log* /some/paths/logfile.log*
```

- when file path not provided, use stdin for input.
want to end, type `done` or `exit`

```Shell
# stdin from pipe
❯ cat testlog.log | grep 'bbb' | ./ltsv-analyzer --base ip --target bytes_sent --operation all --sort sum
Print results by BASE KEY [ip] and TARGET KEY [bytes_sent]
 ip         SUM(bytes_sent)     AVG(bytes_sent)         LOG COUNTER
-------------------------------------------------------------------
1.1.1.4             1123419             1123419                   1
1.1.1.2               11729                5864                   2
```

```Shell
# stdin from keyboard
❯ ./ltsv-analyzer --base ip --target bytes_sent --operation all --sort sum --rev
user:aaa        host:aaa.com    response_time:2.012     ip:1.1.1.1      bytes_sent:10224        time:2020/08/28 13:57:32
user:bbb        host:bbb.com    response_time:1.338     ip:1.1.1.2      bytes_sent:5047 time:2020/08/28 13:57:36
user:aaa        host:aaa.com    response_time:5.132     ip:1.1.1.1      bytes_sent:2432 time:2020/08/28 13:59:59
user:bbb        host:bbb.com    response_time:1.243     ip:1.1.1.4      bytes_sent:1123419      time:2020/08/28 14:02:32
done   # for end input

Print results by BASE KEY [ip] and TARGET KEY [bytes_sent]
 ip         SUM(bytes_sent)     AVG(bytes_sent)         LOG COUNTER
-------------------------------------------------------------------
1.1.1.2                5047                5047                   1
1.1.1.1               12656                6328                   2
1.1.1.4             1123419             1123419                   1
```

### show sum of target value

- want to get total & counter bytes_sent by each user

```Shell
$ ./ltsv-analyzer -B user -T bytes_sent -O sum,cnt ./testlog.log

Print results by BASE KEY [user] and TARGET KEY [bytes_sent]
user     SUM(bytes_sent)         LOG COUNTER
--------------------------------------------
bbb             1135148                   3
aaa               12656                   2
```

### show all datas about base key and target value

- want to get whole data (total / average bytes_sent and counter) by each client ips

```Shell
# ./ltsv-analyzer -B ip -T bytes_sent -O all ./testlog.log

Print results by BASE KEY [ip] and TARGET KEY bytes_sent
 ip         SUM(bytes_sent)     AVG(bytes_sent)         LOG COUNTER
-------------------------------------------------------------------
1.1.1.1               12656                6328                   2
1.1.1.2               11729                5864                   2
1.1.1.4             1123419             1123419                   1
```

### sorting

- want to sort by sum column (default : descending)

```Shell
❯ ./ltsv-analyzer --base ip --target bytes_sent --operation all --sort sum  ./testlog.log

Print results by BASE KEY [ip] and TARGET KEY [bytes_sent]
 ip         SUM(bytes_sent)     AVG(bytes_sent)         LOG COUNTER
-------------------------------------------------------------------
1.1.1.4             1123419             1123419                   1
1.1.1.1               12656                6328                   2
1.1.1.2               11729                5864                   2
```

- if want ascending sort, just add --rev(-r|-R) option

```Shell
❯ ./ltsv-analyzer --base ip --target bytes_sent --operation all --sort sum --rev  ./testlog.log

Print results by BASE KEY [ip] and TARGET KEY [bytes_sent]
 ip         SUM(bytes_sent)     AVG(bytes_sent)         LOG COUNTER
-------------------------------------------------------------------
1.1.1.2               11729                5864                   2
1.1.1.1               12656                6328                   2
1.1.1.4             1123419             1123419                   1
```
