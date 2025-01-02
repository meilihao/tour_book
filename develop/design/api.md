# api

# 同步/异步
1. 尽量设计为异步
1. 异步接口

	1. 文档标明异步
	1. 异步api被调用时返回该job的jobId, 如果该api是创建资源的类型, 还需返回该资源的id
1. 采用json返回

	json比xml更具表现力(比如表示array)

# 调用安全
ref:
- [创建已签名的 AWS API 请求](https://docs.aws.amazon.com/zh_cn/IAM/latest/UserGuide/create-signed-request.html)

## 问题和对策
1. 内容被嗅探 : 使用https
1. 防重放 : timestamp + nonce + sign
1. 如何检查nonce : timestamp + redis

```go
var (
DiffTimeRange  = 10
TimeRangeByMin = 5
TimeRange      = TimeRangeByMin * 60 // 5分钟的nonce放一起, 可以根据请求的速度动态调整
)

reqKeys := []string{"app_id", "nonce", "timestamp"}
us := ctx.Request.PostForm
for _, v := range reqKeys {
		if us.Get(v) == "" {
			ctx.ErrorJson(fmt.Errorf("无效参数%s", v).Error())

			return
		}
  }

// 检查顺序: 时间, nonce, sign
// 
now := time.Now().Unix()
	timestamp, _ := strconv.ParseInt(us.Get("timestamp"), 10, 64)
	if !(timestamp > 0 && (now-timestamp <= DiffTimeRange || timestamp-now <= DiffTimeRange)) { // 仅允许与服务器差DiffTimeRange 秒
		 // 超出了时间阈值

		return
	}

	// 重放检查
	if cache.SIsMember(fmt.Sprintf("replay_nonce_%d", timestamp/TimeRange), us.Get("nonce")).Val() {
		//已经使用过了

		return
	}

if ctx.GetHeader("Authorization") == "" || ctx.GetHeader("Authorization") != GenerateSign(us, appSecret) {
	// 签名错误

		return
	}
  
  // 业务

	cache.SAdd(fmt.Sprintf("replay_nonce_%d", timestamp/TimeRange), us.Get("nonce"))
```

```go
func GenerateSign(us url.Values, appSecret string) string {
	keys := make([]string, 0, len(us))
	for k := range us {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for k, v := range keys {
		keys[k] = v + "=" + us.Get(v)
	}

	mac := hmac.New(sha1.New, []byte(appSecret))
	mac.Write([]byte(strings.Join(keys, "&")))

	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}
```

```go
import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/robfig/cron"
)

var (
	GlobalCron = cron.New()
)

func RunCron() {
	GlobalCron.AddFunc(fmt.Sprintf("%d * * * *", 2*TimeRangeByMin), CleanReplayNonce)

	GlobalCron.Start()
}

// 删除超过了时间戳验证的nonce: 因为时间长的已经被时间戳的判断给拦截了，nonce的验证只需要验证时间有效期内的即可
func CleanReplayNonce() {
	now := time.Now().Unix()
	matchKey := now/TimeRange - 2 // (now-2*TimeRange)/TimeRange

	keys, err := cache.Keys("replay_nonce_*").Result()
	if err != nil {
		sugar.Error(err)

		return
	}
	sugar.Debugf("delete replay_nonce max(%d) : %v", matchKey, keys)

	var tmp int64
	for _, v := range keys {
		tmp, _ = strconv.ParseInt(strings.TrimPrefix(v, "replay_nonce_"), 10, 64)

		if tmp < matchKey {
			sugar.Info("del expired replay_nonce_* : ", v)

			cache.Del(v)
		}
	}
}
```

# 其他
- 禁止ping测试网络是否连通, 推荐使用`nc -w 5 -vz <client_ip> <client_port>`

	部分网络禁ping
- 禁止通过端口连通来测试服务是否可用, 推荐通过api接口

	部分应用端口能连接但是内部逻辑已卡死

### FAQ
1. 为什么防重放时不使用sign作为nonce: 假设实际需求要做短时间内发送重复内容因此这样请求的sign是一样的, 那么第二个请求开始就被认为是重放了. 那你认为如果使用更精确的时间比如纳秒作为时间戳, sign总不一样了吧, 这样就ok了. 这样的话, 单机,单进程/线程还行, 分布式/多进程就game over了. 