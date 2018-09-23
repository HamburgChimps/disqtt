package main

import (
	// "bufio"
	"context"
	"fmt"
	// "io"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/eclipse/paho.mqtt.golang/paho"
)

func main() {

	server := "broker:1883"
	qos := 2
	clientid := "curiosity"
	username := "testuser"
	password := "testpassword"

	c, err := paho.NewClient(paho.OpenTCPConn(server))

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

	log.Println(cp.UsernameFlag, cp.PasswordFlag)

	ca, err := c.Connect(cp)
	if err != nil {
		log.Fatalln(err)
	}
	if ca.ReasonCode != 0 {
		log.Fatalf("Failed to connect to %s : %d - %s", server, ca.ReasonCode, ca.Properties.ReasonString)
	}

	fmt.Printf("Connected to %s\n", server)

	duration := time.Duration(4 * time.Second)
	checkTicker := time.NewTicker(duration)
	go func() {
		for t := range checkTicker.C {
			fmt.Println("Heartbeat at", t)
			message := "Hello"
			if _, err = c.Publish(context.Background(), &paho.Publish{
				Topic:   "/services/12345/heartbeat",
				QoS:     byte(qos),
				Retain:  false,
				Payload: []byte(message),
			}); err != nil {
				log.Println(err)
			}
		}
	}()

	simulateConfigChange := time.Duration(10 * time.Second)
	configTicker := time.NewTicker(simulateConfigChange)
	go func() {
		for t := range configTicker.C {
			fmt.Println("Configuration changed at", t)
			message := "Hello"
			if _, err = c.Publish(context.Background(), &paho.Publish{
				Topic:   "/services/12345/config",
				QoS:     byte(qos),
				Retain:  false,
				Payload: []byte(message),
			}); err != nil {
				log.Println(err)
			}
		}
	}()

	otherTickerDuration := time.Duration(3 * time.Second)
	otherTicker := time.NewTicker(otherTickerDuration)
	go func() {
		for t := range otherTicker.C {
			fmt.Println("Other changed at", t)
			message := "Hello"
			if _, err = c.Publish(context.Background(), &paho.Publish{
				Topic:   "/default",
				QoS:     byte(qos),
				Retain:  false,
				Payload: []byte(message),
			}); err != nil {
				log.Println(err)
			}
		}
	}()

	ic := make(chan os.Signal, 1)
	signal.Notify(ic, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-ic
		fmt.Println("signal received, exiting")
		if c != nil {
			checkTicker.Stop()
			configTicker.Stop()
			otherTicker.Stop()
			d := &paho.Disconnect{ReasonCode: 0}
			c.Disconnect(d)
		}
		os.Exit(0)
	}()

	time.Sleep(1600 * time.Second)

}
