package app

import (
	"golang.org/x/crypto/ripemd160"
	"bytes"
	"encoding/binary"
	"github.com/tendermint/go-crypto"
	"errors"
	"github.com/tendermint/go-wire"
	"reflect"
)

func Hash(data []byte) []byte {
	hasher := ripemd160.New()
	hasher.Write(data)
	hash := hasher.Sum(nil)
	return hash
}

func GetPrivateKeyFromBytes(raw []byte) (prikey crypto.PrivKeyEd25519, err error)  {
	if (len(raw) != 64) {
		err = errors.New("raw data must be 64 bytes")
	}

	buf := &bytes.Buffer{}
	err = binary.Write(buf, binary.BigEndian, raw)
	if (err != nil) {
		return
	}

	err = binary.Read(buf, binary.BigEndian, &prikey)
	if (err != nil) {
		return
	}

	return
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

func GetSignatureFromBytes(raw []byte) (sig crypto.SignatureEd25519, err error)  {
	if (len(raw) != 64) {
		err = errors.New("raw data must be 64 bytes")
	}

	buf := &bytes.Buffer{}
	err = binary.Write(buf, binary.BigEndian, raw)
	if (err != nil) {
		return
	}

	err = binary.Read(buf, binary.BigEndian, &sig)
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

/**
 http://stackoverflow.com/questions/21011023/copy-pointer-values-a-b-in-golang
 usage:
	type data struct {
		a string
		b string
	}

	func main() {
		old := &data{
			"works1",
			"works2",
		}
		var new *data = &data{}
		CloneValue(old, new)
		fmt.Println(new)
	}
 */
func CloneValue(source interface{}, destin interface{}) {
	x := reflect.ValueOf(source)
	if x.Kind() == reflect.Ptr {
		starX := x.Elem()
		y := reflect.New(starX.Type())
		starY := y.Elem()
		starY.Set(starX)
		reflect.ValueOf(destin).Elem().Set(y.Elem())
	} else {
		destin = x.Interface()
	}
}