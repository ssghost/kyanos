package parser

import (
	"fmt"
	"kyanos/agent/protocol"
	"kyanos/common"
	"strconv"
	"strings"
)

const (
	kSimpleStringMarker = '+'
	kErrorMarker        = '-'
	kIntegerMarker      = ':'
	kBulkStringsMarker  = '$'
	kArrayMarker        = '*'
	kTerminalSequence   = "\r\n"
	kNullSize           = -1
)

var redisCommandsMap map[string][]string

func init() {
	redisCommandsMap = make(map[string][]string)
	redisCommandsMap["ACL LOAD"] = []string{"ACL LOAD"}
	redisCommandsMap["ACL SAVE"] = []string{"ACL SAVE"}
	redisCommandsMap["ACL LIST"] = []string{"ACL LIST"}
	redisCommandsMap["ACL USERS"] = []string{"ACL USERS"}
	redisCommandsMap["ACL GETUSER"] = []string{"ACL GETUSER", "username"}
	redisCommandsMap["ACL SETUSER"] = []string{"ACL SETUSER", "username", "[rule [rule ...]]"}
	redisCommandsMap["ACL DELUSER"] = []string{"ACL DELUSER", "username [username ...]"}
	redisCommandsMap["ACL CAT"] = []string{"ACL CAT", "[categoryname]"}
	redisCommandsMap["ACL GENPASS"] = []string{"ACL GENPASS", "[bits]"}
	redisCommandsMap["ACL WHOAMI"] = []string{"ACL WHOAMI"}
	redisCommandsMap["ACL LOG"] = []string{"ACL LOG", "[count or RESET]"}
	redisCommandsMap["ACL HELP"] = []string{"ACL HELP"}
	redisCommandsMap["APPEND"] = []string{"APPEND", "key", "value"}
	redisCommandsMap["AUTH"] = []string{"AUTH", "[username]", "password"}
	redisCommandsMap["BGREWRITEAOF"] = []string{"BGREWRITEAOF"}
	redisCommandsMap["BGSAVE"] = []string{"BGSAVE", "[SCHEDULE]"}
	redisCommandsMap["BITCOUNT"] = []string{"BITCOUNT", "key", "[start end]"}
	redisCommandsMap["BITFIELD"] = []string{"BITFIELD", "key", "[GET type offset]", "[SET type offset value]", "[INCRBY type offset increment]", "[OVERFLOW WRAP|SAT|FAIL]"}
	redisCommandsMap["BITOP"] = []string{"BITOP", "operation", "destkey", "key [key ...]"}
	redisCommandsMap["BITPOS"] = []string{"BITPOS", "key", "bit", "[start]", "[end]"}
	redisCommandsMap["BLPOP"] = []string{"BLPOP", "key [key ...]", "timeout"}
	redisCommandsMap["BRPOP"] = []string{"BRPOP", "key [key ...]", "timeout"}
	redisCommandsMap["BRPOPLPUSH"] = []string{"BRPOPLPUSH", "source", "destination", "timeout"}
	redisCommandsMap["BLMOVE"] = []string{"BLMOVE", "source", "destination", "LEFT|RIGHT", "LEFT|RIGHT", "timeout"}
	redisCommandsMap["BZPOPMIN"] = []string{"BZPOPMIN", "key [key ...]", "timeout"}
	redisCommandsMap["BZPOPMAX"] = []string{"BZPOPMAX", "key [key ...]", "timeout"}
	redisCommandsMap["CLIENT CACHING"] = []string{"CLIENT CACHING", "YES|NO"}
	redisCommandsMap["CLIENT ID"] = []string{"CLIENT ID"}
	redisCommandsMap["CLIENT INFO"] = []string{"CLIENT INFO"}
	redisCommandsMap["CLIENT KILL"] = []string{"CLIENT KILL", "[ip:port]", "[ID client-id]", "[TYPE normal|master|slave|pubsub]", "[USER username]", "[ADDR ip:port]", "[SKIPME yes/no]"}
	redisCommandsMap["CLIENT LIST"] = []string{"CLIENT LIST", "[TYPE normal|master|replica|pubsub]", "[ID client-id [client-id ...]]"}
	redisCommandsMap["CLIENT GETNAME"] = []string{"CLIENT GETNAME"}
	redisCommandsMap["CLIENT GETREDIR"] = []string{"CLIENT GETREDIR"}
	redisCommandsMap["CLIENT UNPAUSE"] = []string{"CLIENT UNPAUSE"}
	redisCommandsMap["CLIENT PAUSE"] = []string{"CLIENT PAUSE", "timeout", "[WRITE|ALL]"}
	redisCommandsMap["CLIENT REPLY"] = []string{"CLIENT REPLY", "ON|OFF|SKIP"}
	redisCommandsMap["CLIENT SETNAME"] = []string{"CLIENT SETNAME", "connection-name"}
	redisCommandsMap["CLIENT TRACKING"] = []string{"CLIENT TRACKING", "ON|OFF", "[REDIRECT client-id]", "[PREFIX prefix [PREFIX prefix ...]]", "[BCAST]", "[OPTIN]", "[OPTOUT]", "[NOLOOP]"}
	redisCommandsMap["CLIENT TRACKINGINFO"] = []string{"CLIENT TRACKINGINFO"}
	redisCommandsMap["CLIENT UNBLOCK"] = []string{"CLIENT UNBLOCK", "client-id", "[TIMEOUT|ERROR]"}
	redisCommandsMap["CLUSTER ADDSLOTS"] = []string{"CLUSTER ADDSLOTS", "slot [slot ...]"}
	redisCommandsMap["CLUSTER BUMPEPOCH"] = []string{"CLUSTER BUMPEPOCH"}
	redisCommandsMap["CLUSTER COUNT-FAILURE-REPORTS"] = []string{"CLUSTER COUNT-FAILURE-REPORTS", "node-id"}
	redisCommandsMap["CLUSTER COUNTKEYSINSLOT"] = []string{"CLUSTER COUNTKEYSINSLOT", "slot"}
	redisCommandsMap["CLUSTER DELSLOTS"] = []string{"CLUSTER DELSLOTS", "slot [slot ...]"}
	redisCommandsMap["CLUSTER FAILOVER"] = []string{"CLUSTER FAILOVER", "[FORCE|TAKEOVER]"}
	redisCommandsMap["CLUSTER FLUSHSLOTS"] = []string{"CLUSTER FLUSHSLOTS"}
	redisCommandsMap["CLUSTER FORGET"] = []string{"CLUSTER FORGET", "node-id"}
	redisCommandsMap["CLUSTER GETKEYSINSLOT"] = []string{"CLUSTER GETKEYSINSLOT", "slot", "count"}
	redisCommandsMap["CLUSTER INFO"] = []string{"CLUSTER INFO"}
	redisCommandsMap["CLUSTER KEYSLOT"] = []string{"CLUSTER KEYSLOT", "key"}
	redisCommandsMap["CLUSTER MEET"] = []string{"CLUSTER MEET", "ip", "port"}
	redisCommandsMap["CLUSTER MYID"] = []string{"CLUSTER MYID"}
	redisCommandsMap["CLUSTER NODES"] = []string{"CLUSTER NODES"}
	redisCommandsMap["CLUSTER REPLICATE"] = []string{"CLUSTER REPLICATE", "node-id"}
	redisCommandsMap["CLUSTER RESET"] = []string{"CLUSTER RESET", "[HARD|SOFT]"}
	redisCommandsMap["CLUSTER SAVECONFIG"] = []string{"CLUSTER SAVECONFIG"}
	redisCommandsMap["CLUSTER SET-CONFIG-EPOCH"] = []string{"CLUSTER SET-CONFIG-EPOCH", "config-epoch"}
	redisCommandsMap["CLUSTER SETSLOT"] = []string{"CLUSTER SETSLOT", "slot", "IMPORTING|MIGRATING|STABLE|NODE", "[node-id]"}
	redisCommandsMap["CLUSTER SLAVES"] = []string{"CLUSTER SLAVES", "node-id"}
	redisCommandsMap["CLUSTER REPLICAS"] = []string{"CLUSTER REPLICAS", "node-id"}
	redisCommandsMap["CLUSTER SLOTS"] = []string{"CLUSTER SLOTS"}
	redisCommandsMap["COMMAND"] = []string{"COMMAND"}
	redisCommandsMap["COMMAND COUNT"] = []string{"COMMAND COUNT"}
	redisCommandsMap["COMMAND GETKEYS"] = []string{"COMMAND GETKEYS"}
	redisCommandsMap["COMMAND INFO"] = []string{"COMMAND INFO", "command-name [command-name ...]"}
	redisCommandsMap["CONFIG GET"] = []string{"CONFIG GET", "parameter"}
	redisCommandsMap["CONFIG REWRITE"] = []string{"CONFIG REWRITE"}
	redisCommandsMap["CONFIG SET"] = []string{"CONFIG SET", "parameter", "value"}
	redisCommandsMap["CONFIG RESETSTAT"] = []string{"CONFIG RESETSTAT"}
	redisCommandsMap["COPY"] = []string{"COPY", "source", "destination", "[DB destination-db]", "[REPLACE]"}
	redisCommandsMap["DBSIZE"] = []string{"DBSIZE"}
	redisCommandsMap["DEBUG OBJECT"] = []string{"DEBUG OBJECT", "key"}
	redisCommandsMap["DEBUG SEGFAULT"] = []string{"DEBUG SEGFAULT"}
	redisCommandsMap["DECR"] = []string{"DECR", "key"}
	redisCommandsMap["DECRBY"] = []string{"DECRBY", "key", "decrement"}
	redisCommandsMap["DEL"] = []string{"DEL", "key [key ...]"}
	redisCommandsMap["DISCARD"] = []string{"DISCARD"}
	redisCommandsMap["DUMP"] = []string{"DUMP", "key"}
	redisCommandsMap["ECHO"] = []string{"ECHO", "message"}
	redisCommandsMap["EVAL"] = []string{"EVAL", "script", "numkeys", "key [key ...]", "arg [arg ...]"}
	redisCommandsMap["EVALSHA"] = []string{"EVALSHA", "sha1", "numkeys", "key [key ...]", "arg [arg ...]"}
	redisCommandsMap["EXEC"] = []string{"EXEC"}
	redisCommandsMap["EXISTS"] = []string{"EXISTS", "key [key ...]"}
	redisCommandsMap["EXPIRE"] = []string{"EXPIRE", "key", "seconds"}
	redisCommandsMap["EXPIREAT"] = []string{"EXPIREAT", "key", "timestamp"}
	redisCommandsMap["FAILOVER"] = []string{"FAILOVER", "[TO host port [FORCE]]", "[ABORT]", "[TIMEOUT milliseconds]"}
	redisCommandsMap["FLUSHALL"] = []string{"FLUSHALL", "[ASYNC|SYNC]"}
	redisCommandsMap["FLUSHDB"] = []string{"FLUSHDB", "[ASYNC|SYNC]"}
	redisCommandsMap["GEOADD"] = []string{"GEOADD", "key", "[NX|XX]", "[CH]", "longitude latitude member [longitude latitude member ...]"}
	redisCommandsMap["GEOHASH"] = []string{"GEOHASH", "key", "member [member ...]"}
	redisCommandsMap["GEOPOS"] = []string{"GEOPOS", "key", "member [member ...]"}
	redisCommandsMap["GEODIST"] = []string{"GEODIST", "key", "member1", "member2", "[m|km|ft|mi]"}
	redisCommandsMap["GEORADIUS"] = []string{"GEORADIUS", "key", "longitude", "latitude", "radius", "m|km|ft|mi", "[WITHCOORD]", "[WITHDIST]", "[WITHHASH]", "[COUNT count [ANY]]", "[ASC|DESC]", "[STORE key]", "[STOREDIST key]"}
	redisCommandsMap["GEORADIUSBYMEMBER"] = []string{"GEORADIUSBYMEMBER", "key", "member", "radius", "m|km|ft|mi", "[WITHCOORD]", "[WITHDIST]", "[WITHHASH]", "[COUNT count [ANY]]", "[ASC|DESC]", "[STORE key]", "[STOREDIST key]"}
	redisCommandsMap["GEOSEARCH"] = []string{"GEOSEARCH", "key", "[FROMMEMBER member]", "[FROMLONLAT longitude latitude]", "[BYRADIUS radius m|km|ft|mi]", "[BYBOX width height m|km|ft|mi]", "[ASC|DESC]", "[COUNT count [ANY]]", "[WITHCOORD]", "[WITHDIST]", "[WITHHASH]"}
	redisCommandsMap["GEOSEARCHSTORE"] = []string{"GEOSEARCHSTORE", "destination", "source", "[FROMMEMBER member]", "[FROMLONLAT longitude latitude]", "[BYRADIUS radius m|km|ft|mi]", "[BYBOX width height m|km|ft|mi]", "[ASC|DESC]", "[COUNT count [ANY]]", "[WITHCOORD]", "[WITHDIST]", "[WITHHASH]", "[STOREDIST]"}
	redisCommandsMap["GET"] = []string{"GET", "key"}
	redisCommandsMap["GETBIT"] = []string{"GETBIT", "key", "offset"}
	redisCommandsMap["GETDEL"] = []string{"GETDEL", "key"}
	redisCommandsMap["GETEX"] = []string{"GETEX", "key", "[EX seconds|PX milliseconds|EXAT timestamp|PXAT milliseconds-timestamp|PERSIST]"}
	redisCommandsMap["GETRANGE"] = []string{"GETRANGE", "key", "start", "end"}
	redisCommandsMap["GETSET"] = []string{"GETSET", "key", "value"}
	redisCommandsMap["HDEL"] = []string{"HDEL", "key", "field [field ...]"}
	redisCommandsMap["HELLO"] = []string{"HELLO", "[protover [AUTH username password] [SETNAME clientname]]"}
	redisCommandsMap["HEXISTS"] = []string{"HEXISTS", "key", "field"}
	redisCommandsMap["HGET"] = []string{"HGET", "key", "field"}
	redisCommandsMap["HGETALL"] = []string{"HGETALL", "key"}
	redisCommandsMap["HINCRBY"] = []string{"HINCRBY", "key", "field", "increment"}
	redisCommandsMap["HINCRBYFLOAT"] = []string{"HINCRBYFLOAT", "key", "field", "increment"}
	redisCommandsMap["HKEYS"] = []string{"HKEYS", "key"}
	redisCommandsMap["HLEN"] = []string{"HLEN", "key"}
	redisCommandsMap["HMGET"] = []string{"HMGET", "key", "field [field ...]"}
	redisCommandsMap["HMSET"] = []string{"HMSET", "key", "field value [field value ...]"}
	redisCommandsMap["HSET"] = []string{"HSET", "key", "field value [field value ...]"}
	redisCommandsMap["HSETNX"] = []string{"HSETNX", "key", "field", "value"}
	redisCommandsMap["HRANDFIELD"] = []string{"HRANDFIELD", "key", "[count [WITHVALUES]]"}
	redisCommandsMap["HSTRLEN"] = []string{"HSTRLEN", "key", "field"}
	redisCommandsMap["HVALS"] = []string{"HVALS", "key"}
	redisCommandsMap["INCR"] = []string{"INCR", "key"}
	redisCommandsMap["INCRBY"] = []string{"INCRBY", "key", "increment"}
	redisCommandsMap["INCRBYFLOAT"] = []string{"INCRBYFLOAT", "key", "increment"}
	redisCommandsMap["INFO"] = []string{"INFO", "[section]"}
	redisCommandsMap["LOLWUT"] = []string{"LOLWUT", "[VERSION version]"}
	redisCommandsMap["KEYS"] = []string{"KEYS", "pattern"}
	redisCommandsMap["LASTSAVE"] = []string{"LASTSAVE"}
	redisCommandsMap["LINDEX"] = []string{"LINDEX", "key", "index"}
	redisCommandsMap["LINSERT"] = []string{"LINSERT", "key", "BEFORE|AFTER", "pivot", "element"}
	redisCommandsMap["LLEN"] = []string{"LLEN", "key"}
	redisCommandsMap["LPOP"] = []string{"LPOP", "key", "[count]"}
	redisCommandsMap["LPOS"] = []string{"LPOS", "key", "element", "[RANK rank]", "[COUNT num-matches]", "[MAXLEN len]"}
	redisCommandsMap["LPUSH"] = []string{"LPUSH", "key", "element [element ...]"}
	redisCommandsMap["LPUSHX"] = []string{"LPUSHX", "key", "element [element ...]"}
	redisCommandsMap["LRANGE"] = []string{"LRANGE", "key", "start", "stop"}
	redisCommandsMap["LREM"] = []string{"LREM", "key", "count", "element"}
	redisCommandsMap["LSET"] = []string{"LSET", "key", "index", "element"}
	redisCommandsMap["LTRIM"] = []string{"LTRIM", "key", "start", "stop"}
	redisCommandsMap["MEMORY DOCTOR"] = []string{"MEMORY DOCTOR"}
	redisCommandsMap["MEMORY HELP"] = []string{"MEMORY HELP"}
	redisCommandsMap["MEMORY MALLOC-STATS"] = []string{"MEMORY MALLOC-STATS"}
	redisCommandsMap["MEMORY PURGE"] = []string{"MEMORY PURGE"}
	redisCommandsMap["MEMORY STATS"] = []string{"MEMORY STATS"}
	redisCommandsMap["MEMORY USAGE"] = []string{"MEMORY USAGE", "key", "[SAMPLES count]"}
	redisCommandsMap["MGET"] = []string{"MGET", "key [key ...]"}
	redisCommandsMap["MIGRATE"] = []string{"MIGRATE", "host", "port", "(key|)", "destination-db", "timeout", "[COPY]", "[REPLACE]", "[AUTH password]", "[AUTH2 username password]", "[KEYS key [key ...]]"}
	redisCommandsMap["MODULE LIST"] = []string{"MODULE LIST"}
	redisCommandsMap["MODULE LOAD"] = []string{"MODULE LOAD", "path", "[ arg [arg ...]]"}
	redisCommandsMap["MODULE UNLOAD"] = []string{"MODULE UNLOAD", "name"}
	redisCommandsMap["MONITOR"] = []string{"MONITOR"}
	redisCommandsMap["MOVE"] = []string{"MOVE", "key", "db"}
	redisCommandsMap["MSET"] = []string{"MSET", "key value [key value ...]"}
	redisCommandsMap["MSETNX"] = []string{"MSETNX", "key value [key value ...]"}
	redisCommandsMap["MULTI"] = []string{"MULTI"}
	redisCommandsMap["OBJECT"] = []string{"OBJECT", "subcommand", "[arguments [arguments ...]]"}
	redisCommandsMap["PERSIST"] = []string{"PERSIST", "key"}
	redisCommandsMap["PEXPIRE"] = []string{"PEXPIRE", "key", "milliseconds"}
	redisCommandsMap["PEXPIREAT"] = []string{"PEXPIREAT", "key", "milliseconds-timestamp"}
	redisCommandsMap["PFADD"] = []string{"PFADD", "key", "element [element ...]"}
	redisCommandsMap["PFCOUNT"] = []string{"PFCOUNT", "key [key ...]"}
	redisCommandsMap["PFMERGE"] = []string{"PFMERGE", "destkey", "sourcekey [sourcekey ...]"}
	redisCommandsMap["PING"] = []string{"PING", "[message]"}
	redisCommandsMap["PSETEX"] = []string{"PSETEX", "key", "milliseconds", "value"}
	redisCommandsMap["PSUBSCRIBE"] = []string{"PSUBSCRIBE", "pattern [pattern ...]"}
	redisCommandsMap["PUBSUB"] = []string{"PUBSUB", "subcommand", "[argument [argument ...]]"}
	redisCommandsMap["PTTL"] = []string{"PTTL", "key"}
	redisCommandsMap["PUBLISH"] = []string{"PUBLISH", "channel", "message"}
	redisCommandsMap["PUNSUBSCRIBE"] = []string{"PUNSUBSCRIBE", "[pattern [pattern ...]]"}
	redisCommandsMap["QUIT"] = []string{"QUIT"}
	redisCommandsMap["RANDOMKEY"] = []string{"RANDOMKEY"}
	redisCommandsMap["READONLY"] = []string{"READONLY"}
	redisCommandsMap["READWRITE"] = []string{"READWRITE"}
	redisCommandsMap["RENAME"] = []string{"RENAME", "key", "newkey"}
	redisCommandsMap["RENAMENX"] = []string{"RENAMENX", "key", "newkey"}
	redisCommandsMap["RESET"] = []string{"RESET"}
	redisCommandsMap["RESTORE"] = []string{"RESTORE", "key", "ttl", "serialized-value", "[REPLACE]", "[ABSTTL]", "[IDLETIME seconds]", "[FREQ frequency]"}
	redisCommandsMap["ROLE"] = []string{"ROLE"}
	redisCommandsMap["RPOP"] = []string{"RPOP", "key", "[count]"}
	redisCommandsMap["RPOPLPUSH"] = []string{"RPOPLPUSH", "source", "destination"}
	redisCommandsMap["LMOVE"] = []string{"LMOVE", "source", "destination", "LEFT|RIGHT", "LEFT|RIGHT"}
	redisCommandsMap["RPUSH"] = []string{"RPUSH", "key", "element [element ...]"}
	redisCommandsMap["RPUSHX"] = []string{"RPUSHX", "key", "element [element ...]"}
	redisCommandsMap["SADD"] = []string{"SADD", "key", "member [member ...]"}
	redisCommandsMap["SAVE"] = []string{"SAVE"}
	redisCommandsMap["SCARD"] = []string{"SCARD", "key"}
	redisCommandsMap["SCRIPT DEBUG"] = []string{"SCRIPT DEBUG", "YES|SYNC|NO"}
	redisCommandsMap["SCRIPT EXISTS"] = []string{"SCRIPT EXISTS", "sha1 [sha1 ...]"}
	redisCommandsMap["SCRIPT FLUSH"] = []string{"SCRIPT FLUSH", "[ASYNC|SYNC]"}
	redisCommandsMap["SCRIPT KILL"] = []string{"SCRIPT KILL"}
	redisCommandsMap["SCRIPT LOAD"] = []string{"SCRIPT LOAD", "script"}
	redisCommandsMap["SDIFF"] = []string{"SDIFF", "key [key ...]"}
	redisCommandsMap["SDIFFSTORE"] = []string{"SDIFFSTORE", "destination", "key [key ...]"}
	redisCommandsMap["SELECT"] = []string{"SELECT", "index"}
	redisCommandsMap["SET"] = []string{"SET", "key", "value", "[EX seconds|PX milliseconds|EXAT timestamp|PXAT milliseconds-timestamp|KEEPTTL]", "[NX|XX]", "[GET]"}
	redisCommandsMap["SETBIT"] = []string{"SETBIT", "key", "offset", "value"}
	redisCommandsMap["SETEX"] = []string{"SETEX", "key", "seconds", "value"}
	redisCommandsMap["SETNX"] = []string{"SETNX", "key", "value"}
	redisCommandsMap["SETRANGE"] = []string{"SETRANGE", "key", "offset", "value"}
	redisCommandsMap["SHUTDOWN"] = []string{"SHUTDOWN", "[NOSAVE|SAVE]"}
	redisCommandsMap["SINTER"] = []string{"SINTER", "key [key ...]"}
	redisCommandsMap["SINTERSTORE"] = []string{"SINTERSTORE", "destination", "key [key ...]"}
	redisCommandsMap["SISMEMBER"] = []string{"SISMEMBER", "key", "member"}
	redisCommandsMap["SMISMEMBER"] = []string{"SMISMEMBER", "key", "member [member ...]"}
	redisCommandsMap["SLAVEOF"] = []string{"SLAVEOF", "host", "port"}
	redisCommandsMap["REPLICAOF"] = []string{"REPLICAOF", "host", "port"}
	redisCommandsMap["SLOWLOG"] = []string{"SLOWLOG", "subcommand", "[argument]"}
	redisCommandsMap["SMEMBERS"] = []string{"SMEMBERS", "key"}
	redisCommandsMap["SMOVE"] = []string{"SMOVE", "source", "destination", "member"}
	redisCommandsMap["SORT"] = []string{"SORT", "key", "[BY pattern]", "[LIMIT offset count]", "[GET pattern [GET pattern ...]]", "[ASC|DESC]", "[ALPHA]", "[STORE destination]"}
	redisCommandsMap["SPOP"] = []string{"SPOP", "key", "[count]"}
	redisCommandsMap["SRANDMEMBER"] = []string{"SRANDMEMBER", "key", "[count]"}
	redisCommandsMap["SREM"] = []string{"SREM", "key", "member [member ...]"}
	redisCommandsMap["STRALGO"] = []string{"STRALGO", "LCS", "algo-specific-argument [algo-specific-argument ...]"}
	redisCommandsMap["STRLEN"] = []string{"STRLEN", "key"}
	redisCommandsMap["SUBSCRIBE"] = []string{"SUBSCRIBE", "channel [channel ...]"}
	redisCommandsMap["SUNION"] = []string{"SUNION", "key [key ...]"}
	redisCommandsMap["SUNIONSTORE"] = []string{"SUNIONSTORE", "destination", "key [key ...]"}
	redisCommandsMap["SWAPDB"] = []string{"SWAPDB", "index1", "index2"}
	redisCommandsMap["SYNC"] = []string{"SYNC"}
	redisCommandsMap["PSYNC"] = []string{"PSYNC", "replicationid", "offset"}
	redisCommandsMap["TIME"] = []string{"TIME"}
	redisCommandsMap["TOUCH"] = []string{"TOUCH", "key [key ...]"}
	redisCommandsMap["TTL"] = []string{"TTL", "key"}
	redisCommandsMap["TYPE"] = []string{"TYPE", "key"}
	redisCommandsMap["UNSUBSCRIBE"] = []string{"UNSUBSCRIBE", "[channel [channel ...]]"}
	redisCommandsMap["UNLINK"] = []string{"UNLINK", "key [key ...]"}
	redisCommandsMap["UNWATCH"] = []string{"UNWATCH"}
	redisCommandsMap["WAIT"] = []string{"WAIT", "numreplicas", "timeout"}
	redisCommandsMap["WATCH"] = []string{"WATCH", "key [key ...]"}
	redisCommandsMap["ZADD"] = []string{"ZADD", "key", "[NX|XX]", "[GT|LT]", "[CH]", "[INCR]", "score member [score member ...]"}
	redisCommandsMap["ZCARD"] = []string{"ZCARD", "key"}
	redisCommandsMap["ZCOUNT"] = []string{"ZCOUNT", "key", "min", "max"}
	redisCommandsMap["ZDIFF"] = []string{"ZDIFF", "numkeys", "key [key ...]", "[WITHSCORES]"}
	redisCommandsMap["ZDIFFSTORE"] = []string{"ZDIFFSTORE", "destination", "numkeys", "key [key ...]"}
	redisCommandsMap["ZINCRBY"] = []string{"ZINCRBY", "key", "increment", "member"}
	redisCommandsMap["ZINTER"] = []string{"ZINTER", "numkeys", "key [key ...]", "[WEIGHTS weight [weight ...]]", "[AGGREGATE SUM|MIN|MAX]", "[WITHSCORES]"}
	redisCommandsMap["ZINTERSTORE"] = []string{"ZINTERSTORE", "destination", "numkeys", "key [key ...]", "[WEIGHTS weight [weight ...]]", "[AGGREGATE SUM|MIN|MAX]"}
	redisCommandsMap["ZLEXCOUNT"] = []string{"ZLEXCOUNT", "key", "min", "max"}
	redisCommandsMap["ZPOPMAX"] = []string{"ZPOPMAX", "key", "[count]"}
	redisCommandsMap["ZPOPMIN"] = []string{"ZPOPMIN", "key", "[count]"}
	redisCommandsMap["ZRANDMEMBER"] = []string{"ZRANDMEMBER", "key", "[count [WITHSCORES]]"}
	redisCommandsMap["ZRANGESTORE"] = []string{"ZRANGESTORE", "dst", "src", "min", "max", "[BYSCORE|BYLEX]", "[REV]", "[LIMIT offset count]"}
	redisCommandsMap["ZRANGE"] = []string{"ZRANGE", "key", "min", "max", "[BYSCORE|BYLEX]", "[REV]", "[LIMIT offset count]", "[WITHSCORES]"}
	redisCommandsMap["ZRANGEBYLEX"] = []string{"ZRANGEBYLEX", "key", "min", "max", "[LIMIT offset count]"}
	redisCommandsMap["ZREVRANGEBYLEX"] = []string{"ZREVRANGEBYLEX", "key", "max", "min", "[LIMIT offset count]"}
	redisCommandsMap["ZRANGEBYSCORE"] = []string{"ZRANGEBYSCORE", "key", "min", "max", "[WITHSCORES]", "[LIMIT offset count]"}
	redisCommandsMap["ZRANK"] = []string{"ZRANK", "key", "member"}
	redisCommandsMap["ZREM"] = []string{"ZREM", "key", "member [member ...]"}
	redisCommandsMap["ZREMRANGEBYLEX"] = []string{"ZREMRANGEBYLEX", "key", "min", "max"}
	redisCommandsMap["ZREMRANGEBYRANK"] = []string{"ZREMRANGEBYRANK", "key", "start", "stop"}
	redisCommandsMap["ZREMRANGEBYSCORE"] = []string{"ZREMRANGEBYSCORE", "key", "min", "max"}
	redisCommandsMap["ZREVRANGE"] = []string{"ZREVRANGE", "key", "start", "stop", "[WITHSCORES]"}
	redisCommandsMap["ZREVRANGEBYSCORE"] = []string{"ZREVRANGEBYSCORE", "key", "max", "min", "[WITHSCORES]", "[LIMIT offset count]"}
	redisCommandsMap["ZREVRANK"] = []string{"ZREVRANK", "key", "member"}
	redisCommandsMap["ZSCORE"] = []string{"ZSCORE", "key", "member"}
	redisCommandsMap["ZUNION"] = []string{"ZUNION", "numkeys", "key [key ...]", "[WEIGHTS weight [weight ...]]", "[AGGREGATE SUM|MIN|MAX]", "[WITHSCORES]"}
	redisCommandsMap["ZMSCORE"] = []string{"ZMSCORE", "key", "member [member ...]"}
	redisCommandsMap["ZUNIONSTORE"] = []string{"ZUNIONSTORE", "destination", "numkeys", "key [key ...]", "[WEIGHTS weight [weight ...]]", "[AGGREGATE SUM|MIN|MAX]"}
	redisCommandsMap["SCAN"] = []string{"SCAN", "cursor", "[MATCH pattern]", "[COUNT count]", "[TYPE type]"}
	redisCommandsMap["SSCAN"] = []string{"SSCAN", "key", "cursor", "[MATCH pattern]", "[COUNT count]"}
	redisCommandsMap["HSCAN"] = []string{"HSCAN", "key", "cursor", "[MATCH pattern]", "[COUNT count]"}
	redisCommandsMap["ZSCAN"] = []string{"ZSCAN", "key", "cursor", "[MATCH pattern]", "[COUNT count]"}
	redisCommandsMap["XINFO"] = []string{"XINFO", "[CONSUMERS key groupname]", "[GROUPS key]", "[STREAM key]", "[HELP]"}
	redisCommandsMap["XADD"] = []string{"XADD", "key", "[NOMKSTREAM]", "[MAXLEN|MINID [=|~] threshold [LIMIT count]]", "*|ID", "field value [field value ...]"}
	redisCommandsMap["XTRIM"] = []string{"XTRIM", "key", "MAXLEN|MINID [=|~] threshold [LIMIT count]"}
	redisCommandsMap["XDEL"] = []string{"XDEL", "key", "ID [ID ...]"}
	redisCommandsMap["XRANGE"] = []string{"XRANGE", "key", "start", "end", "[COUNT count]"}
	redisCommandsMap["XREVRANGE"] = []string{"XREVRANGE", "key", "end", "start", "[COUNT count]"}
	redisCommandsMap["XLEN"] = []string{"XLEN", "key"}
	redisCommandsMap["XREAD"] = []string{"XREAD", "[COUNT count]", "[BLOCK milliseconds]", "STREAMS", "key [key ...]", "ID [ID ...]"}
	redisCommandsMap["XGROUP"] = []string{"XGROUP", "[CREATE key groupname ID|$ [MKSTREAM]]", "[SETID key groupname ID|$]", "[DELCONSUMER key groupname consumername]"}
	redisCommandsMap["XREADGROUP"] = []string{"XREADGROUP", "GROUP group consumer", "[COUNT count]", "[BLOCK milliseconds]", "[NOACK]", "STREAMS", "key [key ...]", "ID [ID ...]"}
	redisCommandsMap["XACK"] = []string{"XACK", "key", "group", "ID [ID ...]"}
	redisCommandsMap["XCLAIM"] = []string{"XCLAIM", "key", "group", "consumer", "min-idle-time", "ID [ID ...]", "[IDLE ms]", "[TIME ms-unix-time]", "[RETRYCOUNT count]", "[FORCE]", "[JUSTID]"}
	redisCommandsMap["XAUTOCLAIM"] = []string{"XAUTOCLAIM", "key", "group", "consumer", "min-idle-time", "start", "[COUNT count]", "[JUSTID]"}
	redisCommandsMap["XPENDING"] = []string{"XPENDING", "key", "group", "[[IDLE min-idle-time] start end count [consumer]]"}
	redisCommandsMap["LATENCY DOCTOR"] = []string{"LATENCY DOCTOR"}
	redisCommandsMap["LATENCY GRAPH"] = []string{"LATENCY GRAPH", "event"}
	redisCommandsMap["LATENCY HISTORY"] = []string{"LATENCY HISTORY", "event"}
	redisCommandsMap["LATENCY LATEST"] = []string{"LATENCY LATEST"}
	redisCommandsMap["LATENCY RESET"] = []string{"LATENCY RESET", "[event [event ...]]"}
	redisCommandsMap["LATENCY HELP"] = []string{"LATENCY HELP"}
	redisCommandsMap["SENTINEL"] = []string{"SENTINEL"}
	redisCommandsMap["REPLCONF ACK"] = []string{"REPLCONF ACK", "offset"}
}

type RedisParser struct {
}

type RedisMessage struct {
	protocol.FrameBase
	payload string
	command string
}

func (m *RedisMessage) Command() string {
	return m.command
}

func (m *RedisMessage) FormatToString() string {
	return fmt.Sprintf("base=[%s] command=[%s] payload=[%s]", m.FrameBase.String(), m.command, m.payload)
}

func ParseSize(decoder *protocol.BinaryDecoder) (int, error) {
	str, err := decoder.ExtractStringUntil(kTerminalSequence)
	if err != nil {
		return 0, err
	}
	const kSizeStrMaxLen = 16
	if len(str) > kSizeStrMaxLen {
		return 0, common.NewInvalidArgument(
			fmt.Sprintf("Redis size string is longer than %d, which indicates traffic is misclassified as Redis.", kSizeStrMaxLen))
	}
	// Length could be -1, which stands for NULL, and means the value is not set.
	// That's different than an empty string, which length is 0.
	// So here we initialize the value to -2.
	size := -2
	size, err = strconv.Atoi(str)
	if err != nil {
		return 0, common.NewInvalidArgument(fmt.Sprintf("String '%s' cannot be parsed as integer", str))
	}
	if size < kNullSize {
		return 0, common.NewInvalidArgument(fmt.Sprintf("Size cannot be less than %d, got '%s'", kNullSize, str))
	}
	return size, nil
}

func ParseBulkString(decoder *protocol.BinaryDecoder, msg *protocol.BaseProtocolMessage) (string, error) {
	const maxLen int = 512 * 1024 * 1024
	length, err := ParseSize(decoder)
	if err != nil {
		return "", err
	}
	if length > maxLen {
		return "", common.NewInvalidArgument(fmt.Sprintf("Length cannot be larger than 512MB, got '%d'", length))
	}
	if length == kNullSize {
		return "<NULL>", nil
	}
	str, err := decoder.ExtractString(length + len(kTerminalSequence))
	if err != nil {
		return "", err
	}
	str = str[:length]
	return str, nil
}

func ParseArray(decoder *protocol.BinaryDecoder, msg *protocol.BaseProtocolMessage) (*RedisMessage, error) {
	size, err := ParseSize(decoder)
	if err != nil {
		return nil, err
	}
	if size == kNullSize {
		return &RedisMessage{
			FrameBase: protocol.NewFrameBase(msg.StartTs, int(msg.TotalBytes())),
			payload:   "[NULL]",
		}, nil
	}
	msgSlice := make([]RedisMessage, 0)
	for i := 0; i < size; i++ {
		_msg, err := ParseMessage(decoder, msg)
		if err != nil {
			return nil, err
		}
		msgSlice = append(msgSlice, *_msg)
	}

	ret := &RedisMessage{
		FrameBase: protocol.NewFrameBase(msg.StartTs, int(msg.TotalBytes())),
	}
	cmd, payload := getCmdAndArgs(msgSlice)
	ret.command = cmd
	ret.payload = payload

	return ret, nil
}

func getCmdAndArgs(payloads []RedisMessage) (string, string) {
	cmd := ""
	finalPayload := ""
	if len(payloads) >= 2 {
		candidateCmd := strings.ToUpper(payloads[0].payload + " " + payloads[1].payload)
		_, ok := redisCommandsMap[candidateCmd]
		if ok {
			payloads = payloads[2:]
			cmd = candidateCmd
			for _, each := range payloads {
				finalPayload += each.payload
				finalPayload += " "
			}
		}
	}

	candidateCmd := strings.ToUpper(payloads[0].payload)
	_, ok := redisCommandsMap[candidateCmd]
	if ok {
		cmd = candidateCmd
		payloads = payloads[1:]
		for _, each := range payloads {
			finalPayload += each.payload
			finalPayload += " "
		}
	}
	return cmd, finalPayload
}

func ParseMessage(decoder *protocol.BinaryDecoder, msg *protocol.BaseProtocolMessage) (*RedisMessage, error) {

	typeMarker, err := decoder.ExtractByte()
	if err != nil {
		return nil, err
	}

	switch typeMarker {
	case kSimpleStringMarker:
		str, err := decoder.ExtractStringUntil(kTerminalSequence)
		if err != nil {
			return nil, err
		}
		return &RedisMessage{
			FrameBase: protocol.NewFrameBase(msg.StartTs, int(msg.TotalBytes())),
			payload:   str,
		}, nil
	case kBulkStringsMarker:
		str, err := ParseBulkString(decoder, msg)
		if err != nil {
			return nil, err
		}
		return &RedisMessage{
			FrameBase: protocol.NewFrameBase(msg.StartTs, int(msg.TotalBytes())),
			payload:   str,
		}, nil
	case kErrorMarker:
		str, err := decoder.ExtractStringUntil(kTerminalSequence)
		if err != nil {
			return nil, err
		}
		return &RedisMessage{
			FrameBase: protocol.NewFrameBase(msg.StartTs, int(msg.TotalBytes())),
			payload:   "-" + str,
		}, nil
	case kIntegerMarker:
		str, err := decoder.ExtractStringUntil(kTerminalSequence)
		if err != nil {
			return nil, err
		}
		return &RedisMessage{
			FrameBase: protocol.NewFrameBase(msg.StartTs, int(msg.TotalBytes())),
			payload:   str,
		}, nil
	case kArrayMarker:
		return ParseArray(decoder, msg)
	default:
		return nil, common.NewInvalidArgument(fmt.Sprintf("Unexpected Redis type marker char (displayed as integer): %d", typeMarker))
	}
}

func (RedisParser) Parse(msg *protocol.BaseProtocolMessage) (protocol.ParsedMessage, error) {
	decoder := protocol.NewBinaryDecoder(msg.Data())
	return ParseMessage(decoder, msg)
}
