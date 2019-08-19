package main

import (
	zmq "github.com/pebbe/zmq4"
)

func main() {
	pub, err := zmq.NewSocket(zmq.PUB)
	pull, err := zmq.NewSocket(zmq.PULL)
	if err != nil {
		panic(err)
	}

	pull.Bind("tcp://*:5000")
	pub.Bind("tcp://*:6000")

	for {
		msg, err := pull.Recv(0)
		if err != nil {
			panic(err)
		}

		pub.SendMessage("", msg)
	}
}
