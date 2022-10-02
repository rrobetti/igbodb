package server

import (
	"context"
	"encoding/json"
	"google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	engine "igbodb/engine"
	mygrpc "igbodb/grpc"
	"log"
)

// var _ api1.Activity = (*grpcServer)(nil)

type igboDBServer struct {
	mygrpc.UnimplementedIgboDBServer
	StorageEngine *engine.StorageEngine
}

func (s *igboDBServer) Retrieve(ctx context.Context, ids *mygrpc.Ids) (*mygrpc.Objects, error) {
	var objects = mygrpc.Objects{
		Items: []*mygrpc.Object{},
	}
	for _, id := range ids.Values {
		resp, err := s.StorageEngine.Retrieve(id)
		if err == engine.ErrIDNotFound {
			return nil, status.Error(codes.NotFound, "id was not found")
		}
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		object := mygrpc.Object{}
		_ = json.Unmarshal([]byte(resp.Description), &object)
		objects.Items = append(objects.Items, &object)
	}

	return &objects, nil
}

func (s *igboDBServer) Create(ctx context.Context, objects *mygrpc.Objects) (*mygrpc.OperationResults, error) {
	var responses = mygrpc.OperationResults{
		Results: []*mygrpc.Result{},
	}
	for _, object := range objects.Items {
		activity := new(mygrpc.Activity)
		activity.Id = object.Id
		json, _ := json.Marshal(object)
		activity.Description = string(json)
		_, err := s.StorageEngine.Insert(activity)
		var result = new(mygrpc.Result)
		if err != nil {
			result.Type = mygrpc.ResultType_FAILURE
			result.Message = err.Error()
		} else {
			result.Type = mygrpc.ResultType_SUCCESS
		}
		responses.Results = append(responses.Results, result)
	}

	return &responses, nil
}

func (s *igboDBServer) Update(ctx context.Context, objects *mygrpc.Objects) (*mygrpc.OperationResults, error) {
	var responses = mygrpc.OperationResults{
		Results: []*mygrpc.Result{},
	}
	for _, object := range objects.Items {
		activity := new(mygrpc.Activity)
		activity.Id = object.Id
		json, _ := json.Marshal(object)
		activity.Description = string(json)
		_, err := s.StorageEngine.Update(activity)
		var result = new(mygrpc.Result)
		if err != nil {
			result.Type = mygrpc.ResultType_FAILURE
			result.Message = err.Error()
		} else {
			result.Type = mygrpc.ResultType_SUCCESS
		}
		responses.Results = append(responses.Results, result)
	}

	return &responses, nil
}

func (s *igboDBServer) Delete(ctx context.Context, ids *mygrpc.Ids) (*mygrpc.OperationResults, error) {
	var responses = mygrpc.OperationResults{
		Results: []*mygrpc.Result{},
	}
	for _, id := range ids.Values {
		_, err := s.StorageEngine.Delete(id)
		var result = new(mygrpc.Result)
		if err != nil {
			result.Type = mygrpc.ResultType_FAILURE
			result.Message = err.Error()
		} else {
			result.Type = mygrpc.ResultType_SUCCESS
		}
		responses.Results = append(responses.Results, result)
	}

	return &responses, nil
}

func NewGRPCServer() *grpc.Server {
	var engine *engine.StorageEngine
	var err error
	if engine, err = engine.NewStorageEngine(); err != nil {
		log.Fatal(err)
	}
	gsrv := grpc.NewServer()
	srv := igboDBServer{
		StorageEngine: engine,
	}
	mygrpc.RegisterIgboDBServer(gsrv, &srv)
	return gsrv
}
