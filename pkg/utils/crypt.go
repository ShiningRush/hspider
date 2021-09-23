package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
)

var PublicKey = `-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDWuY4Gff8FO3BAKetyvNgGrdZM
9CMNoe45SzHMXxAPWw6E2idaEjqe5uJFjVx55JW+5LUSGO1H5MdTcgGEfh62ink/
cNjRGJpR25iVDImJlLi2izNs9zrQukncnpj6NGjZu/2z7XXfJb4XBwlrmR823hpC
umSD1WiMl1FMfbVorQIDAQAB
-----END PUBLIC KEY-----`

func Encrypt(str string) (string, error) {
	b, _ := pem.Decode([]byte(PublicKey))
	if b == nil {
		return "", fmt.Errorf("decode pem key failed")
	}

	pub, err := x509.ParsePKIXPublicKey(b.Bytes)
	if err != nil {
		return "", fmt.Errorf("parse pk failed: %w", err)
	}

	rsaPk := pub.(*rsa.PublicKey)
	bs, err := rsa.EncryptPKCS1v15(rand.Reader, rsaPk, []byte(str))
	if err != nil {
		return "", fmt.Errorf("encrypt failed: %w", err)
	}

	return base64.StdEncoding.EncodeToString(bs), nil
}
