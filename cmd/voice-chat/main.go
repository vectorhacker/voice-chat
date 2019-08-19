package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"math/rand"

	zmq "github.com/pebbe/zmq4"
	chat "github.com/vectorhacker/voice-chat/pkg/voice-chat"
)

var (
	endpoint = flag.String("server", "localhost", "the server endpoint")
	speaker  = flag.String("speaker", "", "the person speaking")
)

func main() {

	flag.Parse()

	pushEndpoint := fmt.Sprintf("tcp://%s:5000", *endpoint)
	subEndpoint := fmt.Sprintf("tcp://%s:6000", *endpoint)

	if *speaker == "" {
		*speaker = randomName()
	}

	out, err := connect(zmq.PUSH, pushEndpoint)
	in, err := connect(zmq.SUB, subEndpoint)
	if err != nil {
		// TODO: better error handling here
		panic(err)
	}
	defer out.Close()
	defer in.Close()

	in.SetSubscribe("")

	log.Println("Speaking as", *speaker)

	playErrChan := chat.Play(in, *speaker)
	recordErrChan := chat.Record(out, *speaker)

	for {
		select {
		case err := <-playErrChan:
			if err != nil {
				log.Fatal(err)
			}
		case err := <-recordErrChan:
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	// > out socket, loops through samples, sends out
	// > in socket, loops through recieves, play sample
}

func connect(socketType zmq.Type, endpoint string) (socket *zmq.Socket, err error) {
	log.Println("Connecing", endpoint)
	socket, err = zmq.NewSocket(socketType)
	if err != nil {
		return
	}

	err = socket.Connect(endpoint)
	if err != nil {
		return
	}

	return
}

var (
	adverbs = [6]string{"amazing", "rad", "joyful", "pretty", "good", "salty"}
	nouns   = [6]string{"pencil", "girl", "boy", "computer", "astronaut", "space cadet"}
)

func randomName() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	return adverbs[r.Intn(6)] + " " + nouns[r.Intn(6)]
}
