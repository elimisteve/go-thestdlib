package main

import (
    "crypto/cipher"
    "crypto/des"
    "crypto/rand"
    "flag"
    "io/ioutil"
    "log"
)

var (
    IV     = []byte("superman") // 8 bytes
    triple = flag.Bool("triple", false, "Use TripleDES")
    do     = flag.String("do", "enc", "Which operation to perform: enc (encryption, default) or dec (decryption)")
)

func MakeKey() []byte {
    size := 8
    if *triple {
        size *= 3
    }
    key := make([]byte, size)
    n, err := rand.Read(key)
    if err != nil {
        log.Fatalf("Failed to read new random key: %s", err)
    }
    if n < size {
        log.Fatalf("Failed to read entire key, only read %d out of %d", n, size)
    }
    return key
}

func Key() []byte {
    file := "des.key"
    key, err := ioutil.ReadFile(file)
    if err != nil {
        log.Println("Failed reading keyfile, making a new one...")
        key = MakeKey()
        err = ioutil.WriteFile(file, key, 0644)
        if err != nil {
            log.Fatalf("Failed saving key to %s: %s", file, err)
        }
    }
    return key
}

func MakeCipher() cipher.Block {
    var c cipher.Block
    var err error
    if *triple {
        c, err = des.NewTripleDESCipher(Key())
    } else {
        c, err = des.NewCipher(Key())
    }
    if err != nil {
        log.Fatalf("Failed making the DES cipher: %s", err)
    }
    return c
}

func Crypt(input, output string) {
    blockCipher := MakeCipher()
    stream := cipher.NewCTR(blockCipher, IV)
    bytes, err := ioutil.ReadFile(input)
    if err != nil {
        log.Fatalf("Failed reading input file: %s", err)
    }
    // Look Ma! No extra memory!
    stream.XORKeyStream(bytes, bytes)
    err = ioutil.WriteFile(output, bytes, 0644)
    if err != nil {
        log.Fatalf("Failed writing output file: %s", err)
    }
}

func Encrypt() {
    Crypt("des.go", "des.go.enc")
}

func Decrypt() {
    Crypt("des.go.enc", "des.go.dec")
}

func main() {
    flag.Parse()

    switch *do {
    case "enc":
        Encrypt()
    case "dec":
        Decrypt()
    default:
        log.Fatalf("%s is not a valid operation. Must be one of enc or dec", *do)
    }
}
