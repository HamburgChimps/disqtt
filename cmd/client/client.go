package main

import (
	"context"
	"github.com/HamburgChimps/disqtt/internal/router"
	"github.com/eclipse/paho.mqtt.golang/paho"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	server := "broker:1883"
	qos := 2
	clientid := "discovery"
	username := "testuser"
	password := "testpassword"

	// paho.SetDebugLogger(log.New(os.Stderr, "SUB: ", log.LstdFlags))
	msgChan := make(chan *paho.Publish)

	c, err := paho.NewClient(
		paho.OpenTCPConn(server),
		paho.DefaultMessageHandler(func(m *paho.Publish) {
			msgChan <- m
		}))

	cp := &paho.Connect{
		KeepAlive:  30,
		ClientID:   clientid,
		CleanStart: true,
		Username:   username,
		Password:   []byte(password),
	}

	if username != "" {
		cp.UsernameFlag = true
	}
	if password != "" {
		cp.PasswordFlag = true
	}

	ca, err := c.Connect(cp)
	if err != nil {
		log.Fatalln(err)
	}
	if ca.ReasonCode != 0 {
		log.Fatalf("Failed to connect to %s : %d - %s", server, ca.ReasonCode, ca.Properties.ReasonString)
	}

	log.Printf("Connected to %s\n", server)

	ic := make(chan os.Signal, 1)
	signal.Notify(ic, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-ic
		log.Println("signal received, exiting")
		if c != nil {
			d := &paho.Disconnect{ReasonCode: 0}
			c.Disconnect(d)
		}
		os.Exit(0)
	}()

	r := router.NewRouter(func(message *paho.Publish) {
		log.Println("DefaultHandler:", message)
		log.Println("Payload:", string(message.Payload))

	})

	r.RegisterHandler("/services/:agent/heartbeat", func(message *paho.Publish) {
		log.Println("Heartbeat received:", message)
		log.Println("Heartbeat:", string(message.Payload))

	})

	r.RegisterHandler("/services/:agent/config", func(message *paho.Publish) {
		log.Println("Configuration changed:", message)
		log.Println("Configuration:", string(message.Payload))
	})

	subscriptions := map[string]paho.SubscribeOptions{
		"/services/+/heartbeat": paho.SubscribeOptions{QoS: byte(qos)},
		"/services/+/config":    paho.SubscribeOptions{QoS: byte(qos)},
		"/default":              paho.SubscribeOptions{QoS: byte(qos)},
	}

	sa, err := c.Subscribe(context.Background(), &paho.Subscribe{
		Subscriptions: subscriptions,
	})
	if err != nil {
		log.Fatalln(err)
	}
	if sa.Reasons[0] != byte(0) {
		log.Fatalf("Failed to subscribe to %s", sa.Properties.ReasonString)
	}
	for subscription := range subscriptions {
		log.Printf("Subscribed to %s", subscription)
	}

	for m := range msgChan {
		r.Route(m)
	}
}
