package app

import (
	"golang.org/x/crypto/ripemd160"
	"bytes"
	"encoding/binary"
	"github.com/tendermint/go-crypto"
	"errors"
	"github.com/tendermint/go-wire"
)

func Hash(data []byte) []byte {
	hasher := ripemd160.New()
	hasher.Write(data)
	hash := hasher.Sum(nil)
	return hash
}

func GetPubKeyFromBytes(raw []byte) (pubkey crypto.PubKeyEd25519, err error)  {
	if (len(raw) != 32) {
		err = errors.New("raw data must be 32 bytes")
	}

	buf := &bytes.Buffer{}
	err = binary.Write(buf, binary.BigEndian, raw)
	if (err != nil) {
		return
	}

	err = binary.Read(buf, binary.BigEndian, &pubkey)
	if (err != nil) {
		return
	}

	return
}

func StructToBytes(o interface{}) (raw []byte, err error)  {
	buf, n := new(bytes.Buffer), new(int)
	wire.WriteBinary(o, buf, n, &err)
	if (err != nil) {
		return
	}

	raw = buf.Bytes()

	return

	//var buf bytes.Buffer        // Stand-in for a network connection
	//enc := gob.NewEncoder(&buf) // Will write to network.
	////dec := gob.NewDecoder(&buf) // Will read from network.
	//
	//err = enc.Encode(o)
	//if (err != nil) {
	//	return
	//}
	//
	//raw = buf.Bytes()
	//
	//return
}

//func BytesToStruct(o interface{}, raw []byte) (interface{}) {
//	//var err error
//	//r, n := bytes.NewReader(raw), new(int)
//	//wire.ReadBinaryPtr(&o, r, 0, n, &err)
//	//
//	//if (err != nil) {
//	//	return nil
//	//}
//	//
//	//return o
//
//	var buf bytes.Buffer
//	_, err := buf.Write(raw)
//	if (err != nil) {
//		return nil
//	}
//
//	dec := gob.NewDecoder(&buf) // Will read from network.
//	err = dec.Decode(&o)
//	if err != nil {
//		return nil
//	}
//
//	return o
//}
