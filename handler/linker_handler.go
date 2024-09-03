package handler

import (
	"fmt"
	"time"

	proto00 "github.com/sicozz/project00/api/v0.0"
	"github.com/sicozz/project00/datatype"
	"github.com/sicozz/project00/utils"
)

func QueryGetServiceInfo() (*proto00.InfoRes, error) {
	programInfo := datatype.ProgramInfo{Version: utils.VERSION, Banner: utils.BANNER}
	return &proto00.InfoRes{
		Version: programInfo.Version,
		Banner:  programInfo.Banner,
	}, nil
}

func HandleSubscription(stream proto00.Linker_SubscribeServer, stmC chan string) error {
	term := 1000
	for {
		select {
		case <-stream.Context().Done():
			err := stream.Context().Err()
			utils.Info(fmt.Sprintf("Follower unsubscribed: %v", err))
			stmC <- fmt.Sprintf("%v", term)
			return err
		default:
			stream.Send(&proto00.Heartbeat{Term: fmt.Sprintf("%v", term)})
			stmC <- "evHeartbeat"
			term = term + 1
			time.Sleep(1 * time.Second)
		}
	}
}
