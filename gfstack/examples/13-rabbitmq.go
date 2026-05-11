// ================================================================
// 示例: RabbitMQ 连接管理器 (internal/library/rabbitmq/rabbitmq.go)
// ================================================================

package rabbitmq

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/streadway/amqp"
)

var Conn *amqp.Connection

type MQConnection struct {
	Conn         *amqp.Connection
	Chan         *amqp.Channel
	ExchangeName string
	QueueName    string
	RouteKey     string
}

// init 自动连接 RabbitMQ（从 config.yaml 读取配置）
func init() {
	var (
		err       error
		MQLinkVar *gvar.Var
		MQLink    g.MapStrStr
	)
	if MQLinkVar, err = g.Cfg().Get(context.TODO(), "queue.rabbitmq"); err != nil {
		panic(err)
	}
	MQLink = MQLinkVar.MapStrStr()
	MQLinkUrl := fmt.Sprintf("amqp://%s:%s@%s:%s%s",
		MQLink["user"], MQLink["pwd"], MQLink["link"], MQLink["port"], MQLink["vhost"])
	if Conn, err = amqp.Dial(MQLinkUrl); err != nil {
		return
	}
}

// Channel 创建通信通道
func (c *MQConnection) Channel() error {
	ch, err := c.Conn.Channel()
	if err != nil {
		return err
	}
	c.Chan = ch
	return nil
}

// Exchange 声明交换机
func (c *MQConnection) Exchange(kind string, durable, autoDeleted bool, args amqp.Table) error {
	return c.Chan.ExchangeDeclare(c.ExchangeName, kind, durable, autoDeleted, false, false, args)
}

// Queue 声明队列
func (c *MQConnection) Queue(durable, autoDeleted bool, args amqp.Table) error {
	_, err := c.Chan.QueueDeclare(c.QueueName, durable, autoDeleted, false, false, args)
	return err
}

// Bind 绑定 Exchange → Queue
func (c *MQConnection) Bind() error {
	return c.Chan.QueueBind(c.QueueName, c.RouteKey, c.ExchangeName, false, nil)
}

// Publish 生产消息
func (c *MQConnection) Publish(message []byte, header amqp.Table, correlationId string) error {
	return c.Chan.Publish(c.ExchangeName, c.RouteKey, false, false, amqp.Publishing{
		Headers:       header,
		ContentType:   "text/html",
		CorrelationId: correlationId,
		Body:          message,
	})
}

// Consume 消费消息
// consumerName: 消费者标识，autoAck: true=自动确认 false=手动确认
func (c *MQConnection) Consume(consumerName string, autoAck bool) (<-chan amqp.Delivery, error) {
	return c.Chan.Consume(c.QueueName, consumerName, autoAck, false, false, false, nil)
}
