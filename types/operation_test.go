package types

import (
	//"testing"
	//"encoding/hex"
	//"strings"
	//"log"
	//"encoding/json"
	//"bytes"
	//"encoding/binary"
	//"github.com/tendermint/go-crypto"
	//"time"
)
//
//
//func TestTxBuilding(t *testing.T) {
//	private_key_hex_str := "0A041B9462CAA4A31BAC3567E0B6E6FD9100787DB2AB433D96F6D178CABFCE90009E64C1B4731BE7DF39A40D5660D84E23885FC465DB5DDAD425789C68CF1A8E"
//
//	// init private key
//	pri_byte_arr, err := hex.DecodeString(private_key_hex_str)
//	buf := bytes.Buffer{}
//	err = binary.Write(&buf, binary.BigEndian, pri_byte_arr)
//	if (err != nil) {
//		t.Error(err.Error())
//		return
//	}
//	var pri_key crypto.PrivKeyEd25519
//	err = binary.Read(&buf, binary.BigEndian, &pri_key);
//	if (err != nil) {
//		t.Error(err.Error())
//		return
//	}
//
//	// init operation
//	var opt AccountCreateOperation
//	opt.ID = "ID"
//	opt.Username = "Username"
//	opt.Pubkey = "009E64C1B4731BE7DF39A40D5660D84E23885FC465DB5DDAD425789C68CF1A8E"
//	opt.UserRegistered = "2017-01-06 09:00:28"
//	opt.DisplayName = "DisplayName"
//
//	opt_arr, err := json.Marshal(opt)
//	if (err != nil) {
//		t.Error(err.Error())
//		return
//	}
//	opt_str := strings.ToUpper(hex.EncodeToString(opt_arr))
//
//
//	// sign the operation
//	sign := pri_key.Sign([]byte(opt_str))
//	sign_str := strings.ToUpper(hex.EncodeToString(sign.Bytes()))
//	sign_str = sign_str[2:len(sign_str)]
//
//	tx := OperationEnvelope{ Type: "AccountCreateOperation",
//		Operation: opt_str,
//		Signature: sign_str,
//		Pubkey: "009E64C1B4731BE7DF39A40D5660D84E23885FC465DB5DDAD425789C68CF1A8E",
//		Fee: 0,
//	}
//
//	str_buf, err := json.Marshal(tx)
//	if (err != nil) {
//		t.Error(err.Error())
//		return
//	}
//
//	log.Println("tx", string(str_buf[:]))
//
//	/////////
//	log.Println("operation", opt_str)
//	log.Println("verify", pri_key.PubKey().VerifyBytes([]byte(opt_str), sign))
//	time.Sleep(1*time.Second) // wait everything to print out
//}
//
//func TestVerifySignature(t *testing.T)  {
//	const json_str = `{"Type":"AccountCreateOperation","Operation":"7B224944223A224944222C22557365726E616D65223A22557365726E616D65222C225075626B6579223A2230303945363443314234373331424537444633394134304435363630443834453233383835464334363544423544444144343235373839433638434631413845222C225573657252656769737465726564223A22323031372D30312D30362030393A30303A3238222C22446973706C61794E616D65223A22446973706C61794E616D65227D","Signature":"FD0110CC64B12ADB08DF97501DEA16BADD3BA05046ADA694C1F2A137394CDA1DE884DDCB8A41D9CA5C88C03D7070F27AA33150E771FD4BF670C48ABC763FB505","Pubkey":"009E64C1B4731BE7DF39A40D5660D84E23885FC465DB5DDAD425789C68CF1A8E","Fee":0}`
//	buf := bytes.Buffer{}
//
//	//var operation json.RawMessage
//	env := OperationEnvelope{
//		//Operation: &operation,
//	}
//
//
//	err := json.Unmarshal([]byte(json_str), &env)
//	if (err != nil) {
//		t.Error(err.Error())
//		return
//	}
//
//	// read public key
//	pubKey_bytearr, err := hex.DecodeString(env.Pubkey)
//	if (err != nil) {
//		t.Error(err.Error())
//		return
//	}
//
//	buf = bytes.Buffer{}
//	err = binary.Write(&buf, binary.BigEndian, pubKey_bytearr)
//	if (err != nil) {
//		t.Error(err.Error())
//		return
//	}
//
//	var pubKey crypto.PubKeyEd25519
//	err = binary.Read(&buf, binary.BigEndian, &pubKey);
//	if (err != nil) {
//		t.Error(err.Error())
//		return
//	}
//
//
//	// read signature
//	sign_bytearr, err := hex.DecodeString(env.Signature)
//	if (err != nil) {
//		t.Error(err.Error())
//		return
//	}
//	buf = bytes.Buffer{}
//	err = binary.Write(&buf, binary.BigEndian, sign_bytearr)
//	if (err != nil) {
//		t.Error(err.Error())
//		return
//	}
//
//	var sign crypto.SignatureEd25519
//	err = binary.Read(&buf, binary.BigEndian, &sign);
//	if (err != nil) {
//		t.Error(err.Error())
//		return
//	}
//
//	// verify
//	log.Println("env.Signature: ", env.Signature)
//	log.Println("pubKey: ", pubKey.KeyString())
//	log.Println("env.Operation", env.Operation)
//	log.Println("verify", pubKey.VerifyBytes([]byte(env.Operation), sign))
//}