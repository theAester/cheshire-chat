package security

import (
  "fmt"
  "os"
  "io"
  "io/ioutil"
  "bytes"
  "path/filepath"
  "log"

  "golang.org/x/crypto/openpgp"
  "golang.org/x/crypto/openpgp/packet"

  "github.com/theAester/cheshire-chat/cmd/daemon/util"
)

var (
  personalInitialized bool = false
  PersonalEntity *openpgp.Entity
  Keyring openpgp.EntityList
  EmptyFingerprint Fingerprint
)

type Fingerprint [20]byte
type ByteString []byte

type ErrNoEntity struct{}
func (e *ErrNoEntity) Error() string {
  return "Personal indentity is not initialized"
}

func GeneratePersonalKeyPair() *openpgp.Entity {
  p, err := openpgp.NewEntity("hirad", "this is just for testing", "penis@man.com", &packet.Config{RSABits: 3072})
  if err != nil {
    log.Fatalf("Oops: %v", err)
    return PersonalEntity
  }

  return p
}

func InitPGP(config *util.Config) error {
  keyringDir := filepath.Join(config.WorkingDir, "keyring")
  if _, err := os.Stat(keyringDir); os.IsNotExist(err) {
    err := os.Mkdir(keyringDir, 0700)
    if err != nil {
      return err
    }
  }

  // Personal
  keyringFile := filepath.Join(keyringDir, "personal.kr")

  if _, err := os.Stat(keyringFile); os.IsNotExist(err) {
    //err := generateKeyRing(keyringFile)
    //if err != nil {
      //return err
    //}

    f, err := os.Create(keyringFile)
    if err != nil {
      return fmt.Errorf("Cannot create keyring file")
    }

    PersonalEntity = GeneratePersonalKeyPair()
    PersonalEntity.SerializePrivate(f, nil)
    f.Close()

    //return &ErrNoEntity{}
  }

  keyringData, err := ioutil.ReadFile(keyringFile)
  if err != nil {
    return fmt.Errorf("Error while reading personal keyring: %v", err)
  }

  tempList, err := openpgp.ReadKeyRing(bytes.NewReader(keyringData))
  if err != nil {
    return fmt.Errorf("Error while parsing personal keyring: %v", err)
  }

  PersonalEntity = tempList[0] // there is only one entity in this file
  personalInitialized = true



  // Contacts
  if !util.ContactsInitialized {
    err := util.InitContacts(config)
    return fmt.Errorf("Error while initializing contacts: %v", err)
  }

  contactList := util.ContactList
  for _, contact := range contactList {
    Keyring = append(Keyring, contact.Entity)
  }

  return nil
}

func Fp2String(fp Fingerprint) string {
  retString := ""
  for i := 0; i<10; i++ {
    firstByte := fp[2*i]
    secondByte := fp[2*i+1]
    
    stringPart := fmt.Sprintf("%x%x", firstByte, secondByte)
    if i != 9 {
      retString = retString + stringPart + " "
    } else {
      retString = retString + stringPart
    }
  }
  return retString
}

func GetContact(fingerprint Fingerprint) *openpgp.Entity {
  for _, entity := range Keyring {
    if bytes.Equal(entity.PrimaryKey.Fingerprint[:], fingerprint[:]) {
      return entity
    }
  }
  return nil
}

func GetContactAll(recepients []Fingerprint) ([]*openpgp.Entity, error) {
  var recepientsMap map[Fingerprint]*openpgp.Entity
  retList := make([]*openpgp.Entity, len(recepients))

  for _, entity := range Keyring {
    recepientsMap[entity.PrimaryKey.Fingerprint] = entity
  }

  idx := 0
  for _, fp := range recepients {
    if recepientsMap[fp] == nil {
      return nil, fmt.Errorf("Fingerprint %s does not exist in keyring", fp)
    }else{
      retList[idx] = recepientsMap[fp]
      idx++
    }
  }
  return retList, nil
}

func EncryptAndSign(message ByteString, recepients []Fingerprint) (ByteString, error) {
  var cipherText bytes.Buffer

  entities, err := GetContactAll(recepients)
  if err != nil {
    return nil, fmt.Errorf("Contact doesnt exist error: %v", err)
  }

                            // writer       to        signer      hint config
  w, err := openpgp.Encrypt(&cipherText, entities, PersonalEntity, nil, nil)
  if err != nil{
    return nil, fmt.Errorf("Error while encrypting message: %v", err)
  }

  _, err = w.Write(message)
  if err != nil{
    return nil, fmt.Errorf("Error while encrypting message: %v", err)
  }

  err = w.Close()
  if err != nil{
    return nil, fmt.Errorf("Error while encrypting message: %v", err)
  }
  
  return cipherText.Bytes(), nil
}

func DecryptAndVerify(message ByteString, expectedFingerprint Fingerprint) (ByteString, Fingerprint, error) {
  md, err := openpgp.ReadMessage(bytes.NewReader(message), Keyring, nil, nil)
  if err != nil {
    return nil, EmptyFingerprint, fmt.Errorf("Error while reading pgp message: %v", err)
  }

  if !md.IsEncrypted || !md.IsSigned {
    return nil, EmptyFingerprint, fmt.Errorf("Bad message received.")
  }

  actualFingerprint := md.SignedBy.PublicKey.Fingerprint

  r := md.UnverifiedBody

  plainText, err := io.ReadAll(r)
  if err != nil {
    return nil, EmptyFingerprint, fmt.Errorf("Error encountered while reading message: %v")
  }

  if md.SignatureError != nil {
    return nil, EmptyFingerprint, fmt.Errorf("There is a problem with the message signature: %v", md.SignatureError)
  }

  return plainText, actualFingerprint, nil
}

