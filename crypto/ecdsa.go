package main

import (
    "crypto/ecdsa"
    "crypto/elliptic"
    "crypto/rand"
    "crypto/sha1"
    "encoding/gob"
    "flag"
    "io"
    "log"
    "os"
)

const (
    KeyFile = "ecdsa.key"
)

var (
    message = flag.String("message", "Nuke the site from orbit, it's the only way to be sure.", "The message to sign")
    do      = flag.String("do", "sign", "The operation to do, sign or verify")
)

func init() {
    var key ecdsa.PrivateKey
    gob.Register(key)

    var cp elliptic.CurveParams
    gob.Register(cp)

    var c elliptic.Curve
    gob.Register(c)
}

func HashMessage() []byte {
    h := sha1.New()
    _, err := io.WriteString(h, *message)
    if err != nil {
        log.Fatalf("Failed to hash message: %s", err)
    }
    return h.Sum(nil)
}

func MakeKey() *ecdsa.PrivateKey {
    key, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
    if err != nil {
        log.Fatalf("Failed to generate key: %s", err)
    }
    return key
}

func Key() *ecdsa.PrivateKey {
    if file, err := os.Open(KeyFile); err == nil {
        defer file.Close()
        var key ecdsa.PrivateKey
        dec := gob.NewDecoder(file)
        err = dec.Decode(&key)
        if err != nil {
            log.Fatalf("Failed reading keyfile: %s", err)
        }
        return &key
    } else {
        log.Println("No keyfile, making a new one...")
        key := MakeKey()
        file, err = os.Create(KeyFile)
        if err != nil {
            log.Fatalf("Can't create file %s: %s", KeyFile, err)
        }
        defer file.Close()
        enc := gob.NewEncoder(file)
        err = enc.Encode(key)
        if err != nil {
            log.Fatalf("Failed saving key to %s: %s", KeyFile, err)
        }
        return key
    }
    panic("not reached")
}

func Sign() {
    key := Key()
    hash := HashMessage()
    r, s, err := ecdsa.Sign(rand.Reader, key, hash)
    if err != nil {
        log.Fatalf("Failed to sign message: %s", err)
    }
    log.Println(r)
    log.Println(s)
}

func Verify() {
}

func main() {
    flag.Parse()
    switch *do {
    case "sign":
        Sign()
    case "verify":
        Verify()
    default:
        log.Fatalf("%s is not a valid operation, must be one of sign or verify", *do)
    }
}
