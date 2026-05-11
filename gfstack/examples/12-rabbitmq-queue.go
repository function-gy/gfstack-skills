// ================================================================
// 示例: RabbitMQ 队列消费者 (queue/TakeCard.go)
// 完整的 RabbitMQ 消费者模式：初始化连接 → Exchange/Queue 声明 → 消费消息 → 业务处理 → 异常退款
// ================================================================

package queue

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/guid"
	"github.com/shopspring/decimal"
	"github.com/streadway/amqp"

	"hotgo/internal/dao"
	"hotgo/internal/library/hgrds/lock"
	"hotgo/internal/library/rabbitmq"
	"hotgo/internal/model/entity"
	"hotgo/internal/service"
)

var LoggerTakeCard = g.Log().Path("logs/queue/TakeCard")

// ================================================================
// RunTakeCard — 消费者启动入口
// ================================================================

func RunTakeCard() {
	var (
		ctx          = gctx.New()
		QueueName    = "TakeCard"
		MQMsg        <-chan amqp.Delivery
		MQConnStruct *rabbitmq.MQConnection
		exchangeName = g.Cfg().MustGet(ctx, "queue.rabbitmq.exchangeName").String()
		queuePrefix  = g.Cfg().MustGet(ctx, "queue.rabbitmq.queuePrefix").String()
		err          error
	)

	// 设置独立的日志输出
	g.DB().SetLogger(LoggerTakeCard)

	// 1. 构建 MQ 连接信息
	MQConnStruct = &rabbitmq.MQConnection{
		Conn:         rabbitmq.Conn,
		ExchangeName: exchangeName,
		QueueName:    fmt.Sprintf("%s%s", queuePrefix, QueueName),
		RouteKey:     fmt.Sprintf("%s.%s%s", exchangeName, queuePrefix, QueueName),
	}

	// 2. 创建 Channel
	if err = MQConnStruct.Channel(); err != nil {
		goto ERR
	}

	// 3. 声明 Exchange（topic 模式）
	if err = MQConnStruct.Exchange("topic", true, false, nil); err != nil {
		goto ERR
	}

	// 4. 声明 Queue（持久化）
	if err = MQConnStruct.Queue(true, false, nil); err != nil {
		goto ERR
	}

	// 5. 绑定 Exchange → Queue
	if err = MQConnStruct.Bind(); err != nil {
		goto ERR
	}

	// 6. 消费消息（非自动确认，手动 ack）
	if MQMsg, err = MQConnStruct.Consume(guid.S(), false); err != nil {
		goto ERR
	}

	LoggerTakeCard.Info(ctx, QueueName+" Queue START SUCCESSFUL")

	// 7. 循环消费
	for msg := range MQMsg {
		ctx = gctx.New()
		if err = HandleBusiness(ctx, msg); err != nil {
			LoggerTakeCard.Error(ctx, err)
		}
		ackMessage(ctx, LoggerTakeCard, msg) // 手动确认
	}
	return

ERR:
	LoggerTakeCard.Error(ctx, err)
	panic(err)
}

// ================================================================
// HandleBusiness — 队列执行业务
// ================================================================

func HandleBusiness(ctx context.Context, msg amqp.Delivery) (err error) {
	var (
		msgData        struct{}
	)

	// 1. 解析消息
	if err = json.Unmarshal(msg.Body, &msgData); err != nil {
		return gerror.New("解析消息失败")
	}
	// TODO 执行的具体的业务逻辑
	return
}
