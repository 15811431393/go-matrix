// Copyright (c) 2018-2019 The MATRIX Authors
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php

package keystore

import (
	"bufio"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/matrix/go-matrix/log"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/matrix/go-matrix/common"
	"github.com/matrix/go-matrix/crypto"
)

func tmpKeyStoreIface(t *testing.T, encrypted bool) (dir string, ks keyStore) {
	d, err := ioutil.TempDir("", "gman-keystore-test")
	if err != nil {
		t.Fatal(err)
	}
	if encrypted {
		ks = &keyStorePassphrase{d, veryLightScryptN, veryLightScryptP}
	} else {
		ks = &keyStorePlain{d}
	}
	return d, ks
}

func LoadKeyStoreIface(d string) (dir string, ks keyStore) {

	ks = &keyStorePassphrase{d, veryLightScryptN, veryLightScryptP}

	return d, ks
}

func TestKeyStorePlain(t *testing.T) {
	dir, ks := tmpKeyStoreIface(t, false)
	defer os.RemoveAll(dir)

	pass := "" // not used but required by API
	k1, account, err := storeNewKey(ks, rand.Reader, pass)
	if err != nil {
		t.Fatal(err)
	}
	k2, err := ks.GetKey(k1.Address, account.URL.Path, pass)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(k1.Address, k2.Address) {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(k1.PrivateKey, k2.PrivateKey) {
		t.Fatal(err)
	}
}

func TestKeyStorePassphrase(t *testing.T) {
	dir, ks := tmpKeyStoreIface(t, true)
	defer os.RemoveAll(dir)

	pass := "foo"
	k1, account, err := storeNewKey(ks, rand.Reader, pass)
	if err != nil {
		t.Fatal(err)
	}
	k2, err := ks.GetKey(k1.Address, account.URL.Path, pass)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(k1.Address, k2.Address) {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(k1.PrivateKey, k2.PrivateKey) {
		t.Fatal(err)
	}
}
func TestKeyStorePassphraseVersion(t *testing.T) {
	_, ks := LoadKeyStoreIface("keystore")

	pass := "xxx"

	k2, err := ks.GetKey(common.HexToAddress("e0b98f47c977267581df784de664074cad88c736"), ".\\keystore\\UTC--2018-11-06T07-06-28.309593000Z--e0b98f47c977267581df784de664074cad88c736", pass)

	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(common.HexToAddress("e0b98f47c977267581df784de664074cad88c736"), k2.Address) {
		t.Fatal(err)
	}
	addr0 := common.HexToAddress("e0b98f47c977267581df784de664074cad88c736")
	sig, error := crypto.SignWithVersion(common.HexToHash("1.0.1-stable").Bytes(), k2.PrivateKey)
	fmt.Println("sig", sig)
	if nil != error {
		fmt.Println("Sign Version Error:%v", sig)
	}
	addr, error := crypto.VerifySignWithVersion(common.HexToHash("1.0.0-stable").Bytes(), sig)
	if nil != error {
		fmt.Println("Verify Sign Version Error:%v", sig)
	}
	if !addr0.Equal(addr) {
		fmt.Errorf("Verify Sign Version error")
	}
}

func TestKeyStorePassphraseHeader(t *testing.T) {

	//pass := "xxx"
	in := bufio.NewReader(os.Stdin)
	fmt.Println(" please input path")
	filename, _, err := in.ReadLine()
	if err != nil {
		log.Crit("Failed to read user input", "err", err)
	}
	keyjson, err := ioutil.ReadFile(string(filename))
	if err != nil {
		log.Crit("Failed to read user input", "err", err)
	}
	fmt.Println(" please input password")
	pass, err := in.ReadString('\n')
	if err != nil {
		log.Crit("Failed to read user input", "err", err)
	}

	key, err := DecryptKey(keyjson, pass)
	if err != nil {
		log.Crit("Failed to read user input", "err", err)
	}
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(" please input sigh hash")

	addr0 := common.HexToAddress("e0b98f47c977267581df784de664074cad88c736")
	sig, error := crypto.SignWithVersion(common.HexToHash("1.0.1-stable").Bytes(), key.PrivateKey)
	fmt.Println("sig", sig)
	if nil != error {
		fmt.Println("Sign Version Error:%v", sig)
	}
	addr, error := crypto.VerifySignWithVersion(common.HexToHash("1.0.0-stable").Bytes(), sig)
	if nil != error {
		fmt.Println("Verify Sign Version Error:%v", sig)
	}
	if !addr0.Equal(addr) {
		fmt.Errorf("Verify Sign Version error")
	}
}

func TestKeyStorePassphraseDecryptionFail(t *testing.T) {
	dir, ks := tmpKeyStoreIface(t, true)
	defer os.RemoveAll(dir)

	pass := "foo"
	k1, account, err := storeNewKey(ks, rand.Reader, pass)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = ks.GetKey(k1.Address, account.URL.Path, "bar"); err != ErrDecrypt {
		t.Fatalf("wrong error for invalid passphrase\ngot %q\nwant %q", err, ErrDecrypt)
	}
}

func TestImportPreSaleKey(t *testing.T) {
	dir, ks := tmpKeyStoreIface(t, true)
	defer os.RemoveAll(dir)

	// file content of a presale key file generated with:
	// python pyethsaletool.py genwallet
	// with password "foo"
	fileContent := "{\"encseed\": \"26d87f5f2bf9835f9a47eefae571bc09f9107bb13d54ff12a4ec095d01f83897494cf34f7bed2ed34126ecba9db7b62de56c9d7cd136520a0427bfb11b8954ba7ac39b90d4650d3448e31185affcd74226a68f1e94b1108e6e0a4a91cdd83eba\", \"manaddr\": \"d4584b5f6229b7be90727b0fc8c6b91bb427821f\", \"email\": \"gustav.simonsson@gmail.com\", \"btcaddr\": \"1EVknXyFC68kKNLkh6YnKzW41svSRoaAcx\"}"
	pass := "foo"
	account, _, err := importPreSaleKey(ks, []byte(fileContent), pass)
	if err != nil {
		t.Fatal(err)
	}
	if account.Address != common.HexToAddress("d4584b5f6229b7be90727b0fc8c6b91bb427821f") {
		t.Errorf("imported account has wrong address %x", account.Address)
	}
	if !strings.HasPrefix(account.URL.Path, dir) {
		t.Errorf("imported account file not in keystore directory: %q", account.URL)
	}
}

// Test and utils for the key store tests in the Matrix JSON tests;
// testdataKeyStoreTests/basic_tests.json
type KeyStoreTestV3 struct {
	Json     encryptedKeyJSONV3
	Password string
	Priv     string
}

type KeyStoreTestV1 struct {
	Json     encryptedKeyJSONV1
	Password string
	Priv     string
}

func TestV3_PBKDF2_1(t *testing.T) {
	t.Parallel()
	tests := loadKeyStoreTestV3("testdata/v3_test_vector.json", t)
	testDecryptV3(tests["wikipage_test_vector_pbkdf2"], t)
}

var testsSubmodule = filepath.Join("..", "..", "tests", "testdata", "KeyStoreTests")

func skipIfSubmoduleMissing(t *testing.T) {
	if !common.FileExist(testsSubmodule) {
		t.Skipf("can't find JSON tests from submodule at %s", testsSubmodule)
	}
}

func TestV3_PBKDF2_2(t *testing.T) {
	skipIfSubmoduleMissing(t)
	t.Parallel()
	tests := loadKeyStoreTestV3(filepath.Join(testsSubmodule, "basic_tests.json"), t)
	testDecryptV3(tests["test1"], t)
}

func TestV3_PBKDF2_3(t *testing.T) {
	skipIfSubmoduleMissing(t)
	t.Parallel()
	tests := loadKeyStoreTestV3(filepath.Join(testsSubmodule, "basic_tests.json"), t)
	testDecryptV3(tests["python_generated_test_with_odd_iv"], t)
}

func TestV3_PBKDF2_4(t *testing.T) {
	skipIfSubmoduleMissing(t)
	t.Parallel()
	tests := loadKeyStoreTestV3(filepath.Join(testsSubmodule, "basic_tests.json"), t)
	testDecryptV3(tests["evilnonce"], t)
}

func TestV3_Scrypt_1(t *testing.T) {
	t.Parallel()
	tests := loadKeyStoreTestV3("testdata/v3_test_vector.json", t)
	testDecryptV3(tests["wikipage_test_vector_scrypt"], t)
}

func TestV3_Scrypt_2(t *testing.T) {
	skipIfSubmoduleMissing(t)
	t.Parallel()
	tests := loadKeyStoreTestV3(filepath.Join(testsSubmodule, "basic_tests.json"), t)
	testDecryptV3(tests["test2"], t)
}

func TestV1_1(t *testing.T) {
	t.Parallel()
	tests := loadKeyStoreTestV1("testdata/v1_test_vector.json", t)
	testDecryptV1(tests["test1"], t)
}

func TestV1_2(t *testing.T) {
	t.Parallel()
	ks := &keyStorePassphrase{"testdata/v1", LightScryptN, LightScryptP}
	addr := common.HexToAddress("cb61d5a9c4896fb9658090b597ef0e7be6f7b67e")
	file := "testdata/v1/cb61d5a9c4896fb9658090b597ef0e7be6f7b67e/cb61d5a9c4896fb9658090b597ef0e7be6f7b67e"
	k, err := ks.GetKey(addr, file, "g")
	if err != nil {
		t.Fatal(err)
	}
	privHex := hex.EncodeToString(crypto.FromECDSA(k.PrivateKey))
	expectedHex := "d1b1178d3529626a1a93e073f65028370d14c7eb0936eb42abef05db6f37ad7d"
	if privHex != expectedHex {
		t.Fatal(fmt.Errorf("Unexpected privkey: %v, expected %v", privHex, expectedHex))
	}
}

func testDecryptV3(test KeyStoreTestV3, t *testing.T) {
	privBytes, _, err := decryptKeyV3(&test.Json, test.Password)
	if err != nil {
		t.Fatal(err)
	}
	privHex := hex.EncodeToString(privBytes)
	if test.Priv != privHex {
		t.Fatal(fmt.Errorf("Decrypted bytes not equal to test, expected %v have %v", test.Priv, privHex))
	}
}

func testDecryptV1(test KeyStoreTestV1, t *testing.T) {
	privBytes, _, err := decryptKeyV1(&test.Json, test.Password)
	if err != nil {
		t.Fatal(err)
	}
	privHex := hex.EncodeToString(privBytes)
	if test.Priv != privHex {
		t.Fatal(fmt.Errorf("Decrypted bytes not equal to test, expected %v have %v", test.Priv, privHex))
	}
}

func loadKeyStoreTestV3(file string, t *testing.T) map[string]KeyStoreTestV3 {
	tests := make(map[string]KeyStoreTestV3)
	err := common.LoadJSON(file, &tests)
	if err != nil {
		t.Fatal(err)
	}
	return tests
}

func loadKeyStoreTestV1(file string, t *testing.T) map[string]KeyStoreTestV1 {
	tests := make(map[string]KeyStoreTestV1)
	err := common.LoadJSON(file, &tests)
	if err != nil {
		t.Fatal(err)
	}
	return tests
}

func TestKeyForDirectICAP(t *testing.T) {
	t.Parallel()
	key := NewKeyForDirectICAP(rand.Reader)
	if !strings.HasPrefix(key.Address.Hex(), "0x00") {
		t.Errorf("Expected first address byte to be zero, have: %s", key.Address.Hex())
	}
}

func TestV3_31_Byte_Key(t *testing.T) {
	t.Parallel()
	tests := loadKeyStoreTestV3("testdata/v3_test_vector.json", t)
	testDecryptV3(tests["31_byte_key"], t)
}

func TestV3_30_Byte_Key(t *testing.T) {
	t.Parallel()
	tests := loadKeyStoreTestV3("testdata/v3_test_vector.json", t)
	testDecryptV3(tests["30_byte_key"], t)
}
