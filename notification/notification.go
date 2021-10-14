package notification

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/segmentio/kafka-go"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"
)

func send(to string, msg string, token string) {
	if to == "" {
		panic(errors.New("No chat_id provided!"))
	} else if msg == "" {
		panic(errors.New("No message provided!"))
	}

	data := url.Values{
		"parse_mode": {"MarkdownV2"},
		"chat_id":    {to},
		"text":       {msg},
	}
	if resp, err := http.PostForm(fmt.Sprintf(`https://api.telegram.org/bot%s/sendMessage`, token), data); err != nil {
		panic(err)
	} else if byte, err := ioutil.ReadAll(resp.Body); err != nil {
		panic(err)
	} else {
		println(string(byte))
	}
}

func getFromHeader(headerName string, arr []kafka.Header) string {
	for _,h := range arr {
		if h.Key == headerName {
			return string(h.Value)
		}
	}
	return ""
}

func Notify(token string) {
	const cmd_text = "msg"
	const cmd_std_in = "std-in"
	const cmd_kafka = "kafka"
	const cmd_kafka_host = "kafka-host"
	const cmd_kafka_topic = "kafka-topic"
	const cmd_to = "to"
	var msg, to, kafka_host, kafka_topic string
	var isStdIn, isKafka bool
	flagSet := flag.NewFlagSet("notify", flag.ExitOnError)
	flagSet.StringVar(&msg, cmd_text, "", "Message to send")
	flagSet.StringVar(&to, cmd_to, "", "Message to whom")
	flagSet.BoolVar(&isStdIn, cmd_std_in, false, "Message from std in instead of msg var")
	flagSet.BoolVar(&isKafka, cmd_kafka, false, "Message from kafka")
	flagSet.StringVar(&kafka_host, cmd_kafka_host, "", "kafka host")
	flagSet.StringVar(&kafka_topic, cmd_kafka_topic, "", "kafka topic")
	nextArg := flag.Args()
	flagSet.Parse(nextArg[1:])

	if isStdIn {
		if bytes, err := ioutil.ReadAll(os.Stdin); err != nil {
			panic(err)
		} else {
			msg = string(bytes)
		}
	} else if isKafka {
		if kafka_host == "" {
			panic(errors.New("kafka host not defined"))
		} else if kafka_topic == "" {
			panic(errors.New("kafka topic not defined"))
		} else {
			reader := kafka.NewReader(kafka.ReaderConfig{
				Brokers:   []string{kafka_host},
				Topic:     kafka_topic,
			})
			reader.SetOffsetAt(context.Background(), time.Now())
			for {
				if m, err := reader.ReadMessage(context.Background()); err!=nil {
					break
				} else {
					fromHost := getFromHeader("from", m.Headers)
					about := getFromHeader("about", m.Headers)
					tempMsg := fmt.Sprintf(`
from: %s
about: %s
message: %s
`, fromHost, about, string(m.Value))
					tempMsg = tempMsg[1:]
					println(tempMsg)
					send(to, tempMsg, token)
				}
			}
		}

	}

}
