package msgq

import (
	"JudgerServer/config"
	"encoding/json"
)

type submitData struct {
	ID       string   `form:"id" binding:"required" json:"id"`
	Code     string   `form:"code" binding:"required" json:"code"`
	Language string   `form:"language" binding:"required" json:"language"`
	Limits   *limits  `form:"limits" binding:"required" json:"limits"`
	Data     []string `form:"data" binding:"required" json:"data"`
}

type limits struct {
	CPU    uint64 `form:"cpu" binding:"required" json:"cpu"`
	Memory uint64 `form:"memory" binding:"required" json:"memory"`
}

// ConsumerRoutine 消费者
func ConsumerRoutine() {
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}

	qname, err := config.GetConfig("RMQ_QNAME")
	if err != nil {
		panic(err)
	}
	q, err := ch.QueueDeclare(qname, true, false, false, false, nil)
	if err != nil {
		panic(err)
	}
	for true {
		msgs, err := ch.Consume(q.Name, "basic_comsume", false, false, false, false, nil)
		if err != nil {
			break
		}
		for msg := range msgs {
			data := msg.Body
			jdata := submitData{}
			err = json.Unmarshal(data, &jdata)
			if err != nil {
				if err = msg.Ack(false); err != nil {
					break
				}
				continue
			}
			// TODO 运行
			msg.Ack(false)
		}
		if err != nil {
			break
		}
	}
	conn.Close()
	panic(err)
}
