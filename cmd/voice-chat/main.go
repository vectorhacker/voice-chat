package main

import (
	"flag"
	"fmt"
	"log"

	zmq "github.com/pebbe/zmq4"
	chat "github.com/vectorhacker/voice-chat/pkg/voice-chat"
)

var (
	endpoint = flag.String("server", "localhost", "the server endpoint")
)

func main() {

	flag.Parse()

	pushEndpoint := fmt.Sprintf("tcp://%s:5000", *endpoint)
	subEndpoint := fmt.Sprintf("tcp://%s:6000", *endpoint)

	out, err := connect(zmq.PUSH, pushEndpoint)
	in, err := connect(zmq.SUB, subEndpoint)
	if err != nil {
		// TODO: better error handling here
		panic(err)
	}
	defer out.Close()
	defer in.Close()

	in.SetSubscribe("")

	playErrChan := chat.Play(in)
	recordErrChan := chat.Record(out)

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
