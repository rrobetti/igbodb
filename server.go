package main

import (
	"context"
	"encoding/json"
	"google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	mygrpc "igbodb/grpc"
	"log"
)

// var _ api1.Activity = (*grpcServer)(nil)

type igboDBServer struct {
	mygrpc.UnimplementedIgboDBServer
	StorageEngine StorageEngine
}

func (s *igboDBServer) Retrieve(ctx context.Context, objectKeys *mygrpc.ObjectKeys) (*mygrpc.Objects, error) {
	var objects = mygrpc.Objects{
		Items: []*mygrpc.Object{},
	}
	for _, key := range objectKeys.Keys {
		value, err := s.StorageEngine.Retrieve(key.Type, key.Id)
		if err == ErrIDNotFound {
			return nil, status.Error(codes.NotFound, "id was not found")
		}
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		object := mygrpc.Object{}
		object.Key = key
		_ = json.Unmarshal([]byte(value), &object.Attributes)
		objects.Items = append(objects.Items, &object)
	}

	return &objects, nil
}

func (s *igboDBServer) Create(ctx context.Context, objects *mygrpc.Objects) (*mygrpc.OperationResults, error) {
	var responses = mygrpc.OperationResults{
		Results: []*mygrpc.Result{},
	}
	for _, object := range objects.Items {
		json, _ := json.Marshal(object.Attributes) //TODO change to a different format to avoid using JSON
		err := s.StorageEngine.Create(object.Key.Type, object.Key.Id, string(json))
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
		json, _ := json.Marshal(object.Attributes)
		err := s.StorageEngine.Update(object.Key.Type, object.Key.Id, string(json))
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

func (s *igboDBServer) Delete(ctx context.Context, objectKeys *mygrpc.ObjectKeys) (*mygrpc.OperationResults, error) {
	var responses = mygrpc.OperationResults{
		Results: []*mygrpc.Result{},
	}
	for _, oKey := range objectKeys.Keys {
		err := s.StorageEngine.Delete(oKey.Type, oKey.Id)
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
	var engine *StorageEngineImpl
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
