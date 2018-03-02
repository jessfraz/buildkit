package control

import (
	"fmt"

	"golang.org/x/net/context"

	controlapi "github.com/moby/buildkit/api/services/control"
	"github.com/moby/buildkit/cache/cacheimport"
	"github.com/moby/buildkit/control"
	"github.com/moby/buildkit/frontend"
	"github.com/moby/buildkit/session"
	"github.com/moby/buildkit/worker"
	"google.golang.org/grpc"
)

// control client implements the controlapi.ControlClient interface directly
// on top on the controller.
type controlClient struct {
	controller *control.Controller
}

// Opt holds the options for the controller client.
type Opt struct {
	SessionManager *session.Manager
	Worker         worker.Worker
	Frontends      map[string]frontend.Frontend
	CacheExporter  *cacheimport.CacheExporter
	CacheImporter  *cacheimport.CacheImporter
}

// NewControlClient returns a new control client.
func NewControlClient(opt Opt) (controlapi.ControlClient, error) {
	// Create the worker controller.
	wc := &worker.Controller{}
	if err := wc.Add(opt.Worker); err != nil {
		return nil, err
	}

	// Create the controller.
	controller, err := control.NewController(control.Opt{
		SessionManager:   opt.SessionManager,
		WorkerController: wc,
		Frontends:        opt.Frontends,
		CacheExporter:    opt.CacheExporter,
		CacheImporter:    opt.CacheImporter,
	})
	if err != nil {
		return nil, fmt.Errorf("creating new controller failed: %v", err)
	}

	return &controlClient{
		controller: controller,
	}, nil
}

func (c *controlClient) DiskUsage(ctx context.Context, in *controlapi.DiskUsageRequest, opts ...grpc.CallOption) (*controlapi.DiskUsageResponse, error) {
	return c.controller.DiskUsage(ctx, in)
}

func (c *controlClient) Prune(ctx context.Context, in *controlapi.PruneRequest, opts ...grpc.CallOption) (controlapi.Control_PruneClient, error) {
	// no-op
	// TODO: implement this.
	return nil, nil
}

func (c *controlClient) Solve(ctx context.Context, in *controlapi.SolveRequest, opts ...grpc.CallOption) (*controlapi.SolveResponse, error) {
	return c.controller.Solve(ctx, in)
}

func (c *controlClient) Status(ctx context.Context, in *controlapi.StatusRequest, opts ...grpc.CallOption) (controlapi.Control_StatusClient, error) {
	// no-op
	// TODO: implement this.
	return nil, nil
}

func (c *controlClient) Session(ctx context.Context, opts ...grpc.CallOption) (controlapi.Control_SessionClient, error) {
	// no-op
	// TODO: implement this.
	return nil, nil
}

func (c *controlClient) ListWorkers(ctx context.Context, in *controlapi.ListWorkersRequest, opts ...grpc.CallOption) (*controlapi.ListWorkersResponse, error) {
	return c.controller.ListWorkers(ctx, in)
}
