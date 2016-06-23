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
	"strings"
	"time"

	"github.com/satori/go.uuid"
	"github.com/snapcore/snapd/asserts"
	"golang.org/x/crypto/openpgp"
)

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

func SignSerial(modelAssertion asserts.Assertion, targetFolder string, vaultServer string) (err error) {
	if modelAssertion.Type() != asserts.ModelType {
		err = errors.New("not a model assertion")
		return err
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
	gnupgHomedir := targetFolder + "/.gnupg/"
	err = os.MkdirAll(gnupgHomedir, 0700)
	Checkerr(err)

	genkey := []byte("Key-Type: 1\nKey-Length: 4096\nName-Real: SERIAL\n")
	err = ioutil.WriteFile("/tmp/gen-key-script", genkey, 0600)
	Checkerr(err)

	Shellexec("gpg", "--homedir="+gnupgHomedir, "--batch", "--gen-key", "/tmp/gen-key-script")

	// Read public key
	f, err := os.Open(gnupgHomedir + "/pubring.gpg")
	Checkerr(err)
	el, err := openpgp.ReadKeyRing(f)
	Checkerr(err)
	entity := getKeyByName(el, "SERIAL")
	openPGPPublicKey := asserts.OpenPGPPublicKey(entity.PrimaryKey)
	encodeKey, err := asserts.EncodePublicKey(openPGPPublicKey)
	Checkerr(err)
	key := string(encodeKey)

	// TODO: clarify the format of encodeKey
	key = strings.Replace(key, "\n", "", -1)

	product_serial, err := ioutil.ReadFile("/sys/class/dmi/id/product_serial")
	Checkerr(err)
	serial := strings.Split(string(product_serial), "\n")[0] + "-" + uuid.NewV4().String()

	content := Serial(authority, key, brand, model, revision, serial, time.Now())
	body := bytes.NewBuffer([]byte(content))

	vaultServer = strings.TrimRight(vaultServer, "/")
	log.Println("vaultServer:", vaultServer)
	r, err := http.Post(vaultServer, "application/x-www-form-urlencoded", body)
	Checkerr(err)
	response, err := ioutil.ReadAll(r.Body)
	if nil != err {
		log.Fatal("Serial Sign error:", err)
	}
	err = ioutil.WriteFile(targetFolder+"/serial.txt", response, 0600)
	Checkerr(err)

	if isJSON(string(response)) {
		log.Fatal("Serial Sign error:", string(response))
	}

	log.Println("Sign serial assertion successfully!.")
	return nil
}
