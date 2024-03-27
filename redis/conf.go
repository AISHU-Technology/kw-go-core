package redis

import (
	"fmt"
	"github.com/AISHU-Technology/kw-go-core/utils"
	red "github.com/go-redis/redis/v8"
	"os"
	"time"
)

type RedisConf struct {
	Addrs []string `json:",default=127.0.0.1:6379"`
	//Port         int      `json:",default=6379"`
	Pass         string `json:",optional"`
	Username     string `json:",optional"`
	DB           int    `json:",default=0,range=[0:15]"`
	Type         string `json:",default=standalone,options=standalone|sentinel|cluster"`
	MasterName   string `json:",optional"`
	PoolSize     int    `json:",default=0"`
	MinIdleConns int    `json:",default=10"`
	MaxRetries   int    `json:",default=1,range=[0:3]"`
}

const (
	// ClusterType means redis cluster.
	clusterType = "cluster"
	// NodeType means redis node.
	masterSlaveType = "master-slave"
	// Nil is an alias of redis.Nil.
	sentinelType   = "sentinel"
	standaloneType = "standalone"
)

var rl *Redis

type Redis struct {
	Write red.Cmdable
	Read  red.Cmdable
}

func getRedisEnv(c RedisConf) *RedisConf {
	env := os.Getenv("REDISCLUSTERMODE")
	if env == sentinelType {
		c.Type = sentinelType
	} else if env == clusterType {
		c.Type = clusterType
	} else if env == standaloneType {
		c.Type = standaloneType
	} else if env == masterSlaveType || (utils.IsNotBlank(os.Getenv("REDISREADPASS")) && utils.IsNotBlank(os.Getenv("REDISWRITEPORT"))) {
		c.Type = masterSlaveType
	}
	return &c
}
func InitRedisConf(c RedisConf) *Redis {
	getRedisEnv(c)
	var redis *Redis
	switch c.Type {
	case clusterType:
		redis = newRedisCluster(c)
	case sentinelType:
		redis = newRedisSentinel(c)
	case standaloneType:
		redis = newRedisNode(c)
	case masterSlaveType:
		redis = newRedisWriteRead()
	}
	rl = redis
	return redis
}

func newRedisNode(rc RedisConf) *Redis {
	host := os.Getenv("REDISHOST")
	port := os.Getenv("REDISPORT")
	if utils.IsNotBlank(host) && utils.IsNotBlank(port) {
		rc.Addrs = []string{fmt.Sprintf("%s:%s", host, port)}
	}
	pass := os.Getenv("REDISREADPASS")
	if utils.IsBlank(pass) {
		pass = rc.Pass
	}
	var client = red.NewClient(&red.Options{
		Addr:            rc.Addrs[0],
		Password:        pass,
		DB:              rc.DB,
		DialTimeout:     8 * time.Second,
		PoolSize:        rc.PoolSize,
		MinIdleConns:    rc.MinIdleConns,
		MaxRetries:      rc.MaxRetries,
		MinRetryBackoff: 12 * time.Millisecond,
		MaxRetryBackoff: 512 * time.Millisecond,
	})
	redis := &Redis{
		Write: client,
		Read:  client,
	}
	return redis
}

func newRedisSentinel(rc RedisConf) *Redis {
	host := os.Getenv("REDISHOST")
	port := os.Getenv("REDISPORT")
	if utils.IsNotBlank(host) && utils.IsNotBlank(port) {
		rc.Addrs = []string{fmt.Sprintf("%s:%s", host, port)}
	}
	username := os.Getenv("REDISUSER")
	if utils.IsBlank(username) {
		username = rc.Username
	}
	pass := os.Getenv("REDISPASS")
	if utils.IsBlank(pass) {
		pass = rc.Pass
	}
	sentineluser := os.Getenv("SENTINELUSER")
	if utils.IsBlank(sentineluser) {
		sentineluser = rc.Username
	}
	sentinelpass := os.Getenv("SENTINELPASS")
	if utils.IsBlank(sentinelpass) {
		sentinelpass = rc.Pass
	}
	masterName := os.Getenv("SENTINELMASTER")
	if utils.IsBlank(masterName) {
		masterName = rc.MasterName
	}
	client := red.NewFailoverClient(&red.FailoverOptions{
		MasterName:       masterName,
		SentinelAddrs:    rc.Addrs,
		RouteByLatency:   true,
		RouteRandomly:    true,
		PoolSize:         rc.PoolSize,
		MinIdleConns:     rc.MinIdleConns,
		MaxRetries:       rc.MaxRetries,
		Password:         pass,
		Username:         username,
		SentinelUsername: sentineluser,
		SentinelPassword: sentinelpass,
		MinRetryBackoff:  12 * time.Millisecond,
		MaxRetryBackoff:  512 * time.Millisecond,
	})
	return &Redis{Write: client, Read: client}
}

func newRedisCluster(rc RedisConf) *Redis {
	host := os.Getenv("REDISHOST")
	port := os.Getenv("REDISPORT")
	if utils.IsNotBlank(host) && utils.IsNotBlank(port) {
		rc.Addrs = []string{fmt.Sprintf("%s:%s", host, port)}
	}
	pass := os.Getenv("REDISPASS")
	if utils.IsBlank(pass) {
		pass = rc.Pass
	}
	client := red.NewClusterClient(&red.ClusterOptions{
		//集群节点地址，理论上只要填一个可用的节点客户端就可以自动获取到集群的所有节点信息。但是最好多填一些节点以增加容灾能力，因为只填一个节点的话，如果这个节点出现了异常情况，则Go应用程序在启动过程中无法获取到集群信息。
		Addrs:           rc.Addrs,
		ReadOnly:        true, // 置为true则允许在从节点上执行只含读操作的命令
		RouteByLatency:  true,
		RouteRandomly:   true,
		Password:        pass,
		PoolSize:        rc.PoolSize,
		MinIdleConns:    rc.MinIdleConns,
		MaxRetries:      rc.MaxRetries,
		MinRetryBackoff: 12 * time.Millisecond,
		MaxRetryBackoff: 512 * time.Millisecond,
	})
	return &Redis{Write: client, Read: client}
}

func newRedisWriteRead() *Redis {

	redis := &Redis{
		Write: red.NewClient(&red.Options{
			Addr:            fmt.Sprintf("%s:%s", os.Getenv("REDISWRITEHOST"), os.Getenv("REDISWRITEPORT")),
			Username:        os.Getenv("REDISWRITEUSER"),
			Password:        os.Getenv("REDISWRITEPASS"),
			DialTimeout:     8 * time.Second,
			PoolSize:        0,
			MinIdleConns:    15,
			MaxRetries:      2,
			MinRetryBackoff: 12 * time.Millisecond,
			MaxRetryBackoff: 512 * time.Millisecond,
		}),
		Read: red.NewClient(&red.Options{
			Addr:            fmt.Sprintf("%s:%s", os.Getenv("REDISREADHOST"), os.Getenv("REDISREADPORT")),
			Username:        os.Getenv("REDISREADUSER"),
			Password:        os.Getenv("REDISREADPASS"),
			DialTimeout:     8 * time.Second,
			PoolSize:        0,
			MinIdleConns:    15,
			MaxRetries:      2,
			MinRetryBackoff: 12 * time.Millisecond,
			MaxRetryBackoff: 512 * time.Millisecond,
		}),
	}
	return redis
}
