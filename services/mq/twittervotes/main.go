package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	nsq "github.com/bitly/go-nsq"
	mgo "gopkg.in/mgo.v2"
)

var db *mgo.Session

func dialdb() error {
	var err error
	log.Println("dialing mongodb: localhost")
	db, err = mgo.Dial("localhost")
	return err
}

func closedb() {
	db.Close()
	log.Println("closed database connection")
}

type poll struct {
	Options []string
}

func loadOptions() ([]string, error) {
	var (
		options []string
		p       poll
	)
	// nil -> no filter
	iter := db.DB("ballots").C("polls").Find(nil).Iter()
	for iter.Next(&p) {
		options = append(options, p.Options...)
	}
	iter.Close()
	return options, iter.Err()
}

// receive only chan
func publishVotes(votes <-chan string) <-chan struct{} {
	stopchan := make(chan struct{}, 1)
	pub, _ := nsq.NewProducer("localhost:4150", nsq.NewConfig())
	go func() {
		// Continually pull values form a channel.
		// For loop will only end when the votes channel is closed, otherwise execution will be blocked until a value comes down the line.
		for vote := range votes {
			pub.Publish("votes", []byte(vote)) // publish vote
		}
		log.Println("Publisher: Stopping")
		pub.Stop()
		log.Println("Publisher: Stopped")
		stopchan <- struct{}{}
	}()
	return stopchan
}

func main() {

	var stoplock sync.Mutex
	stop := false // allows us to access from many goroutines at the same time
	stopChan := make(chan struct{}, 1)
	signalChan := make(chan os.Signal, 1)

	// Will only run when specified signals is sent, which allows us to perfrom teardown code before exiting the code
	go func() {
		<-signalChan // blocks waiting for the signal(trying to read from chan)
		stoplock.Lock()
		stop = true
		stoplock.Unlock()
		log.Println("Stopping...")
		stopChan <- struct{}{}
		closeConn()
	}()
	// Send the signal down teh signalChan when someone tries to halt the program, with either sigint interrupt or the sigterm termination posix signals
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	if err := dialdb(); err != nil {
		log.Fatalln("failed to dial MongoDB: ", err)
	}
	defer closedb()
	// start things
	votes := make(chan string)                  // chan for votes
	publisherStoppedChan := publishVotes(votes) // makes mq connection and waits to publish
	twitterStoppedChan := startTwitterStream(stopChan, votes)
	// refresher
	go func() {
		for {
			time.Sleep(1 * time.Minute)
			closeConn()
			stoplock.Lock()
			if stop {
				stoplock.Unlock()
				break
			}
			stoplock.Unlock()
		}
	}()
	<-twitterStoppedChan   // blocking here while attempting to read from it
	close(votes)           // when stopChan sent to twitterStoppedChan then close votes chan
	<-publisherStoppedChan // closing votes chan will cause the publishers for...range loop to exit
}
