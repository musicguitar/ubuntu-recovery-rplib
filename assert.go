package rplib

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/snapcore/snapd/asserts"
	"github.com/snapcore/snapd/asserts/assertstest"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	"golang.org/x/crypto/openpgp/packet"
)

func isJSON(s string) bool {
	var js map[string]interface{}
	return json.Unmarshal([]byte(s), &js) == nil

}

func GenerateKey(bits int) (asserts.PrivateKey, []byte, error) {
	privKey, rsaPrivKey := assertstest.GenerateKey(bits)
	log.Println("new generated Public key id of accountPrivKey: ", privKey.PublicKey().ID())
	pgpKey := packet.NewRSAPrivateKey(time.Now(), rsaPrivKey)

	// export armored private key
	armored, err := ArmorBuffer(pgpKey)
	if err != nil {
		return privKey, armored, err
	}
	return privKey, armored, err
}

func ArmorBuffer(pgpKey *packet.PrivateKey) (ret []byte, err error) {
	buf := bytes.NewBuffer(nil)
	w, err := armor.Encode(buf, openpgp.PrivateKeyType, make(map[string]string))
	if err != nil {
		return nil, err
	}

	err = pgpKey.Serialize(w)
	if err != nil {
		return nil, err
	}
	err = w.Close()
	if err != nil {
		return nil, err
	}
	return []byte(buf.String()), nil
}

func NewModel(db assertstest.SignerDB, otherHeaders map[string]interface{}, keyID string) asserts.Assertion {
	if otherHeaders == nil {
		otherHeaders = make(map[string]interface{})
	}

	if otherHeaders["timestamp"] == nil {
		otherHeaders["timestamp"] = time.Now().Format(time.RFC3339)
	}
	m, err := db.Sign(asserts.ModelType, otherHeaders, nil, keyID)
	if err != nil {
		panic(err)
	}
	return m
}

func NewDevice(db assertstest.SignerDB, pubkey asserts.PublicKey, otherHeaders map[string]interface{}, keyID string) asserts.Assertion {
	if otherHeaders == nil {
		otherHeaders = make(map[string]interface{})
	}

	if otherHeaders["device-key"] == nil {
		otherHeaders["device-key"] = keyEncode(pubkey)
	}

	if otherHeaders["timestamp"] == nil {
		otherHeaders["timestamp"] = time.Now().Format(time.RFC3339)
	}
	m, err := db.Sign(asserts.SerialType, otherHeaders, nil, keyID)
	if err != nil {
		panic(err)
	}
	return m
}

func keyEncode(pubkey asserts.PublicKey) string {
	// TODO: clarify format of key string
	encodeKey, err := asserts.EncodePublicKey(pubkey)
	if err != nil {
		panic(err)
	}
	ret := string(encodeKey)
	return ret
}

func NewSerialRequest(modelAssertion asserts.Assertion, devicePrivKey asserts.PrivateKey, serial string, revision string, nonce string) (asserts.Assertion, error) {
	// TODO: check asserts.SerialAssertionType
	serialRequest, err := asserts.SignWithoutAuthority(asserts.SerialRequestType,
		map[string]interface{}{
			"brand-id":   modelAssertion.Header("brand-id"),
			"model":      modelAssertion.Header("model"),
			"series":     modelAssertion.Header("series"),
			"revision":   revision,
			"device-key": keyEncode(devicePrivKey.PublicKey()),
			"request-id": nonce,
		},
		// TODO: include HW-DETAILS in body
		// TODO: check if serial can be integer
		[]byte(fmt.Sprintf("serial: %s", serial)),
		devicePrivKey)
	return serialRequest, err
}

func SendSerialRequest(serialRequest asserts.Assertion, vaultServer string, apikey string) ([]byte, error) {
	// TODO: check asserts.SerialRequestType
	body := bytes.NewBuffer(asserts.Encode(serialRequest))

	// send http/https request
	vaultServer = strings.TrimRight(vaultServer, "/")
	vaultServer = vaultServer + "/serial"
	log.Println("send request to:", vaultServer)
	req, err := http.NewRequest("POST", vaultServer, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("api-key", apikey)

	client := &http.Client{}
	response, err := client.Do(req)
	if nil != err {
		return nil, err
	}
	defer response.Body.Close()

	returnBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	if isJSON(string(returnBody)) {
		// serial-vault return error message in json form
		log.Println("Serial Sign error:", string(returnBody))
		return nil, errors.New(string(returnBody))
	}

	return returnBody, nil
}
