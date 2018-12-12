package store

import (
	"bytes"
	"encoding/gob"
	"errors"
	"log"
	"time"

	"github.com/boltdb/bolt"
)

const (
	kOngoingTable = "ONGOINGTABLE"
	kHandledTable = "HANDLEDTABLE"
	kRequestsDB   = "requests.db"
)

type Request struct {
	Username     string
	Email        string
	TokenAddress string
	Timestamp    string
	Description  string
}

func (p *Request) toBytes() ([]byte, error) {
	var buf bytes.Buffer
	// Create an encoder and send a value.
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(p)
	if err != nil {
		log.Fatal("encode:", err)
		return nil, err
	}

	return buf.Bytes(), nil
}

func unmarshalRequest(b []byte) (*Request, error) {
	u := &Request{}
	var buf = bytes.Buffer{}
	buf.Write(b)
	// Create a decoder and receive a value.
	dec := gob.NewDecoder(&buf)
	err := dec.Decode(u)
	if err != nil {
		log.Fatal("decode:", err)
		return nil, err
	}
	return u, nil
}

func saveRequestToTable(req *Request, table string) error {
	db, e := bolt.Open(kRequestsDB, 0600, nil)
	if e != nil {
		return e
	}
	defer db.Close()

	return db.Update(func(tx *bolt.Tx) error {
		b, e := tx.CreateBucketIfNotExists([]byte(table))
		if e != nil {
			return e
		}
		v := time.Now().String()
		req.Timestamp = v
		bytes, err := req.toBytes()
		if err != nil {
			return err
		}
		return b.Put([]byte(req.TokenAddress), bytes)
	})
}

func SubmitRequest(req *Request) error {
	return saveRequestToTable(req, kOngoingTable)
}

func ApproveRequest(token string) error {
	r, e := GetOngoingRequest(token)
	if e != nil {
		return e
	}
	e2 := DeleteOngoingRequest(token)
	if e2 != nil {
		return e2
	}
	return saveRequestToTable(r, kHandledTable)
	//return nil
}

func getRequestFromTable(key string, table string) (*Request, error) {
	db, e := bolt.Open(kRequestsDB, 0600, nil)
	if e != nil {
		return nil, e
	}
	defer db.Close()
	var content []byte
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(table))
		if b != nil {
			content = b.Get([]byte(key))
		}
		return nil
	})
	if content != nil {
		u, e := unmarshalRequest(content)
		if e == nil {
			return u, nil
		}
		return nil, e
	}
	return nil, errors.New("key not exists")
}

func GetOngoingRequest(key string) (*Request, error) {
	return getRequestFromTable(key, kOngoingTable)
}

func DeleteOngoingRequest(key string) error {
	db, e := bolt.Open(kRequestsDB, 0600, nil)
	if e != nil {
		return e
	}
	defer db.Close()
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(kOngoingTable))
		if b == nil {
			return errors.New("no bucket")
		}
		return b.Delete([]byte(key))
	})
}

func GetHandledRequest(key string) (*Request, error) {
	return getRequestFromTable(key, kHandledTable)
}

func getAllRequestsFromTable(table string) ([]*Request, error) {
	db, e := bolt.Open(kRequestsDB, 0600, nil)
	if e != nil {
		return nil, e
	}
	defer db.Close()

	result := []*Request{}
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(table))
		if b == nil {
			return errors.New("error b")
		} //end if
		c := b.Cursor()
		if c == nil {
			return errors.New("nil cursor")
		}
		for k, v := c.First(); k != nil; k, v = c.Next() {
			item, ex := unmarshalRequest(v)
			if ex != nil {
				continue
			}
			if item == nil {
				continue
			}
			result = append(result, item)
		} //end for
		return nil
	})
	return result, nil
}

func GetOngoingRequests() ([]*Request, error) {
	return getAllRequestsFromTable(kOngoingTable)
}

func GetApprovedRequests() ([]*Request, error) {
	return getAllRequestsFromTable(kHandledTable)
}
