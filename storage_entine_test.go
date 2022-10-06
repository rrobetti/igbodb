package main

import (
	"log"
	"testing"
)

func TestCRUD(t *testing.T) {
	var engine *StorageEngineImpl
	var err error
	if engine, err = engine.NewStorageEngine(); err != nil {
		log.Fatal(err)
	}
	var storage = "storage1"
	var key = "key1"
	var value = "value1"
	var valueUpdated = "value1Updated"

	//CLEAN PREVIOUS TEST
	err = engine.Delete(storage, key)
	if err != nil {
		t.Errorf("got %v want successful deletion", err)
	}

	// CREATE
	err = engine.Create(storage, key, value)
	if err != nil {
		t.Errorf("got %v want successful creation", err)
	}
	if err != nil {
		t.Errorf("got %v want successful deletion", err)
	}

	// RETRIEVE
	var valueRet string
	valueRet, err = engine.Retrieve(storage, key)
	if err != nil {
		t.Errorf("got %v want successful retrieve", err)
	}
	if (valueRet != value) {
		t.Errorf("got %v want %v", valueRet, value)
	}

	//UPDATE
	err = engine.Update(storage, key, valueUpdated)
	if err != nil {
		t.Errorf("got %v want successful update", err)
	}
	valueRet, err = engine.Retrieve(storage, key)
	if err != nil {
		t.Errorf("got %v want successful retrieve", err)
	}
	if (valueRet != valueUpdated) {
		t.Errorf("got %v want %v", valueRet, value)
	}

	//DELETE
	err = engine.Delete(storage, key)
	if err != nil {
		t.Errorf("got %v want successful deletion", err)
	}
	valueRet, err = engine.Retrieve(storage, key)
	if err == nil {
		t.Errorf("got succesful response want %v ", ErrIDNotFound)
	}
}
