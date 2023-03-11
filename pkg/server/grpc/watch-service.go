package grpc

import (
	"context"
	log "github.com/sirupsen/logrus"
	"github.com/tislib/data-handler/pkg/abs"
	"github.com/tislib/data-handler/pkg/stub"
)

type watchGrpcService struct {
	stub.WatchServiceServer
	watchService abs.WatchService
}

func (w *watchGrpcService) Watch(req *stub.WatchRequest, res stub.WatchService_WatchServer) error {
	localCtx, cancel := context.WithCancel(res.Context())
	defer func() {
		cancel()
	}()

	out := w.watchService.Watch(localCtx, abs.WatchParams{
		Namespace:  req.Namespace,
		Resource:   req.Resource,
		Query:      req.Query,
		BufferSize: 500,
	})

	for message := range out {
		err := res.Send(message)

		if err != nil {
			cancel()
			log.Error(err)
			return err
		}
	}

	return nil
}

func NewWatchServiceServer(service abs.WatchService) stub.WatchServiceServer {
	return &watchGrpcService{watchService: service}
}
