package main

import (
	"log"
	"os"
	"time"

	console "github.com/asynkron/goconsole"
	"github.com/asynkron/protoactor-go/actor"
)

type NoInfluence string

func (NoInfluence) NotInfluenceReceiveTimeout() {}

func main() {
	log.Println("Example(Receive timeout)")

	system := actor.NewActorSystem()
	c := 0

	rootContext := system.Root
	props := actor.PropsFromFunc(func(context actor.Context) {
		log := log.New(os.Stdout, "[actor] ", log.LstdFlags|log.Lmsgprefix)
		switch msg := context.Message().(type) {
		case *actor.Started:
			context.SetReceiveTimeout(1 * time.Second)

		case *actor.Stopped:
			log.Printf("Stopped")

		case *actor.ReceiveTimeout:
			c++
			log.Printf("ReceiveTimeout: %d", c)
			context.SetReceiveTimeout(1 * time.Second)

		case string:
			log.Printf("Received '%s'", msg)
			if msg == "cancel" {
				context.CancelReceiveTimeout()
				log.Println("Canceled")
			}

		case NoInfluence:
			log.Println("Received a no-influence message")

		}
	})

	pid := rootContext.Spawn(props)

	log := log.New(os.Stdout, "[main] ", log.LstdFlags|log.Lmsgprefix)
	log.Println("Sending messages")
	for i := 0; i < 6; i++ {
		rootContext.Send(pid, "hello")
		time.Sleep(500 * time.Millisecond)
	}

	log.Println("Sending no-influence messages")
	for i := 0; i < 6; i++ {
		rootContext.Send(pid, NoInfluence("hello"))
		time.Sleep(500 * time.Millisecond)
	}

	rootContext.Send(pid, "cancel")
	log.Println("hit [return] to finish")
	console.ReadLine()
	rootContext.Stop(pid)
}
