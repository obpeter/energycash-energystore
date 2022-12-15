package store

import (
	"at.ourproject/energystore/store/ebow"
	"fmt"
)

func OpenStorageTest(tenant string, basedir string) (*BowStorage, error) {
	db, err := ebow.Open(fmt.Sprintf("%s/%s", basedir, tenant))
	if err != nil {
		return nil, err
	}
	return &BowStorage{db}, nil
}
