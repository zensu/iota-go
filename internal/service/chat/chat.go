package chat

import (
	"context"
	"fmt"
	"github.com/zensu/iota-go/internal/protogen"
	"sync"
)

type Connection struct {
	protogen.UnimplementedBroadcastServer

	stream protogen.Broadcast_StartStreamServer
	id     string
	active bool
	error  chan error
}

type PoolConnections struct {
	protogen.UnimplementedBroadcastServer
	Connections []*Connection
}

func (p *PoolConnections) CreateStream(pc *protogen.Connect, stream protogen.Broadcast_StartStreamServer) error {
	conn := &Connection{
		stream: stream,
		id:     pc.User.Id,
		active: true,
		error:  make(chan error),
	}

	p.Connections = append(p.Connections, conn)

	return <-conn.error
}

func (p *PoolConnections) BroadcastMessage(ctx context.Context, msg *protogen.Message) (*protogen.Close, error) {
	wait := sync.WaitGroup{}
	done := make(chan int)

	for _, conn := range p.Connections {
		wait.Add(1)

		go func(msg *protogen.Message, conn *Connection) {
			defer wait.Done()

			if conn.active {
				err := conn.stream.Send(msg)
				fmt.Printf("Sending message to: %v from %v", conn.id, msg.Id)

				if err != nil {
					fmt.Printf("Error with Stream: %v - Error: %v\n", conn.stream, err)
					conn.active = false
					conn.error <- err
				}
			}
		}(msg, conn)

	}

	go func() {
		wait.Wait()
		close(done)
	}()

	<-done
	return &protogen.Close{}, nil
}