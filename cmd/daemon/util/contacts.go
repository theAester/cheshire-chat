package util

import (
  "fmt"
  "os"
  "io/ioutil"
  "bytes"
  "path/filepath"
  "encoding/json"
  "encoding/base64"

  "golang.org/x/crypto/openpgp"
)

var ContactsInitialized bool = false
var ContactList []ContactData

type ContactSerial struct {
  Nickname string `json:"nickname"`
  Fingerprint [20]byte `json:"fingerprint"`
  Description string `json:"description"`
  AsciiArt string `json:"ascii_art"`
  Entity string `json:"entity"`
}

type ContactData struct {
  Nickname string 
  Fingerprint [20]byte 
  Description string 
  AsciiArt string
  Entity *openpgp.Entity
}

func InitContacts(config *Config) error {
  contactsDir := filepath.Join(config.WorkingDir, "contacts")
  if _, err := os.Stat(contactsDir); os.IsNotExist(err) {
    err := os.Mkdir(contactsDir, 0700)
    if err != nil {
      return err
    }
  }

  // Personal
  contactsFile := filepath.Join(contactsDir, "contacts.dat")

  if _, err := os.Stat(contactsFile); os.IsNotExist(err) {
    //err := generatecontacts(contactsFile)
    //if err != nil {
      //return err
    //}

    err := ioutil.WriteFile(contactsFile, []byte("[]"), 0644) 
    if err != nil {
      return fmt.Errorf("Cannot initialize contacts file.")
    }
  }
  contactsRaw, err := ioutil.ReadFile(contactsFile)
  if err != nil {
    return fmt.Errorf("Cannot read contacts file.")
  }
  
  var contactListSerial []ContactSerial
  err = json.Unmarshal(contactsRaw, &contactListSerial)
  if err != nil {
    return fmt.Errorf("Error while parsing contact data: %v", err)
  }

  err = deserializeContactList(&contactListSerial)
  ContactsInitialized = true

  return nil
}

func deserializeContactList(contactListSerial *[]ContactSerial) error {
  for _, contactSerial := range *contactListSerial {
    var contactData ContactData
    contactData.Nickname = contactSerial.Nickname
    contactData.Fingerprint = contactSerial.Fingerprint
    contactData.Description = contactSerial.Description
    entityString := contactSerial.Entity
    entityBytes, err := base64.StdEncoding.DecodeString(entityString)
    if err != nil {
      return fmt.Errorf("Error while decoding entity data for %s: %v",
                          contactData.Nickname, err)
    }
    entityKr, err := openpgp.ReadKeyRing(bytes.NewReader(entityBytes)) 
    if err != nil {
      return fmt.Errorf("Error while reading entity data for %s: %v",
                          contactData.Nickname, err)
    }
    entity := entityKr[0] // theres only one Entity
    contactData.Entity = entity
    ContactList = append(ContactList, contactData) 
  }
  return nil
}

func serializeContactList(contactListSerial *[]ContactSerial) error {
  for _, contactData := range ContactList {
    var contactSerial ContactSerial
    contactSerial.Nickname = contactData.Nickname
    contactSerial.Fingerprint = contactData.Fingerprint
    contactSerial.Description = contactData.Description
    entity := contactData.Entity
    var buf bytes.Buffer
    err := entity.Serialize(&buf)
    if err != nil {
      return fmt.Errorf("Error while encoding entity data for %s: %v",
                          contactData.Nickname, err)
    }
    entityBytes := buf.Bytes()
    entityString := base64.StdEncoding.EncodeToString(entityBytes)
    contactSerial.Entity = entityString
    *contactListSerial = append(*contactListSerial, contactSerial) 
  }
  return nil
}
