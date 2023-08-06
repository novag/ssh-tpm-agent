package key

import (
	"reflect"
	"testing"

	"github.com/google/go-tpm/tpm2"
	"github.com/google/go-tpm/tpm2/transport/simulator"
)

func TestECDSACreateKey(t *testing.T) {
	tpm, err := simulator.OpenSimulator()
	if err != nil {
		t.Fatal(err)
	}
	defer tpm.Close()
	k, err := CreateKey(tpm, tpm2.TPMAlgECDSA, []byte(""))
	if err != nil {
		t.Fatalf("failed key import: %v", err)
	}

	// Test if we can load the key
	// signer/signer_test.go tests the signing of the key
	_, err = LoadKey(tpm, k)
	if err != nil {
		t.Fatalf("failed loading key: %v", err)
	}
}

func TestRSACreateKey(t *testing.T) {
	tpm, err := simulator.OpenSimulator()
	if err != nil {
		t.Fatal(err)
	}
	defer tpm.Close()
	k, err := CreateKey(tpm, tpm2.TPMAlgRSA, []byte(""))
	if err != nil {
		t.Fatalf("failed key import: %v", err)
	}

	// Test if we can load the key
	// signer/signer_test.go tests the signing of the key
	_, err = LoadKey(tpm, k)
	if err != nil {
		t.Fatalf("failed loading key: %v", err)
	}
}

func mustPublic(data []byte) tpm2.TPM2BPublic {
	return tpm2.BytesAs2B[tpm2.TPMTPublic](data)
}

func mustPrivate(data []byte) tpm2.TPM2BPrivate {
	return tpm2.TPM2BPrivate{
		Buffer: data,
	}
}

func TestMarshalling(t *testing.T) {
	cases := []struct {
		k *Key
	}{
		{
			k: &Key{
				Version: 1,
				PIN:     HasPIN,
				Type:    tpm2.TPMAlgECDSA,
				Public:  mustPublic([]byte("public")),
				Private: mustPrivate([]byte("private")),
			},
		},
		{
			k: &Key{
				Version: 1,
				PIN:     NoPIN,
				Type:    tpm2.TPMAlgECDSA,
				Public:  mustPublic([]byte("public")),
				Private: mustPrivate([]byte("private")),
			},
		},
		{
			k: &Key{
				Version: 1,
				PIN:     HasPIN,
				Type:    tpm2.TPMAlgRSA,
				Public:  mustPublic([]byte("public")),
				Private: mustPrivate([]byte("private")),
			},
		},
		{
			k: &Key{
				Version: 1,
				PIN:     NoPIN,
				Type:    tpm2.TPMAlgRSA,
				Public:  mustPublic([]byte("public")),
				Private: mustPrivate([]byte("private")),
			},
		},
	}

	for _, c := range cases {
		b := EncodeKey(c.k)
		k, err := DecodeKey(b)
		if err != nil {
			t.Fatalf("test failed: %v", err)
		}

		if !reflect.DeepEqual(k, c.k) {
			t.Fatalf("keys are not the same")
		}
	}
}