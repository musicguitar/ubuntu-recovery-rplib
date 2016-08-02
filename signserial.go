package rplib

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/satori/go.uuid"
	"github.com/snapcore/snapd/asserts"
	"golang.org/x/crypto/openpgp"
)

const KEYLENGTH = 4096
const KEYID = "SERIAL"
const SerialUnsigned = "serialUnsigned.txt"
const SerialSigned = "serial.txt"

func Serial(authority, key, brand, model, revision, serial string, t time.Time) string {
	content := fmt.Sprintf("type: serial\nauthority-id: %s\ndevice-key: %s\nbrand-id: %s\nmodel: %s\nrevision: %s\nserial: %s\ntimestamp: %s\n\n%s\n", authority, key, brand, model, revision, serial, t.UTC().Format("2006-01-02T15:04:05Z"), key)
	return content
}

func getKeyByName(keyring openpgp.EntityList, name string) *openpgp.Entity {
	for _, entity := range keyring {
		for _, ident := range entity.Identities {
			if ident.UserId.Name == name {
				return entity
			}
		}
	}

	return nil
}

func isJSON(s string) bool {
	var js map[string]interface{}
	return json.Unmarshal([]byte(s), &js) == nil

}

func SerialAssertionGen(modelAssertion asserts.Assertion, targetFolder string) (serialAssertion string, err error) {
	gnupgHomedir := filepath.Join(targetFolder, ".gnupg/")
	publicKeyFile := filepath.Join(gnupgHomedir, "pubring.gpg")

	if modelAssertion.Type() != asserts.ModelType {
		err = errors.New("not a model assertion")
		return "", err
	}

	authority := modelAssertion.Header("authority-id")
	log.Println("authority:", authority)
	brand := modelAssertion.Header("brand-id")
	log.Println("brand:", brand)
	model := modelAssertion.Header("model")
	log.Println("model:", model)
	revision := modelAssertion.Header("revision")
	log.Println("revision:", revision)

	// generate gpg key pair
	log.Println("targetFolder:", targetFolder)
	os.MkdirAll(targetFolder, 0755)

	if _, err := os.Stat(gnupgHomedir); err == nil {
		// gpg folder already exist
		err = errors.New(fmt.Sprintf("gpg folder %s already exist!", gnupgHomedir))
		return "", err
	}
	err = os.MkdirAll(gnupgHomedir, 0700)
	Checkerr(err)

	genkey := []byte(fmt.Sprintf("Key-Type: 1\nKey-Length: %d\nName-Real: %s\n", KEYLENGTH, KEYID))
	err = ioutil.WriteFile("/tmp/gen-key-script", genkey, 0600)
	Checkerr(err)

	Shellexec("gpg", "--homedir="+gnupgHomedir, "--batch", "--gen-key", "/tmp/gen-key-script")

	// Read public key
	f, err := os.Open(publicKeyFile)
	Checkerr(err)
	el, err := openpgp.ReadKeyRing(f)
	Checkerr(err)
	entity := getKeyByName(el, KEYID)
	openPGPPublicKey := asserts.OpenPGPPublicKey(entity.PrimaryKey)
	encodeKey, err := asserts.EncodePublicKey(openPGPPublicKey)
	Checkerr(err)
	key := string(encodeKey)

	// TODO: clarify the format of encodeKey
	key = strings.Replace(key, "\n", "", -1)

	product_serial_content, err := ioutil.ReadFile(SMBIOS_SERIAL)
	product_serial := strings.Split(string(product_serial_content), "\n")[0]
	Checkerr(err)
	serial := product_serial + "-" + uuid.NewV4().String()

	serialAssertion = Serial(authority, key, brand, model, revision, serial, time.Now())

	return serialAssertion, nil
}

func SignSerial(modelAssertion asserts.Assertion, targetFolder string, vaultServer string, apikey string) (err error) {
	signedFile := filepath.Join(targetFolder, SerialSigned)

	content, err := SerialAssertionGen(modelAssertion, targetFolder)
	if nil != err {
		return err
	}
	body := bytes.NewBuffer([]byte(content))

	// send http request
	vaultServer = strings.TrimRight(vaultServer, "/")
	log.Println("vaultServer:", vaultServer)
	req, err := http.NewRequest("POST", vaultServer, body)
	Checkerr(err)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("api-key", apikey)

	client := &http.Client{}
	response, err := client.Do(req)
	if nil != err {
		log.Fatal("Serial Sign error:", err)
	}
	defer response.Body.Close()

	returnBody, _ := ioutil.ReadAll(response.Body)
	if isJSON(string(returnBody)) {
		// sign server return error message in json form
		log.Fatal("Serial Sign error:", string(returnBody))
	}

	err = ioutil.WriteFile(signedFile, returnBody, 0600)
	Checkerr(err)

	log.Println("Sign serial assertion successfully!.")
	return nil
}
