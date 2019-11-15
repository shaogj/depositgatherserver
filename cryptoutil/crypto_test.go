package cryptoutil

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAESCBC(t *testing.T) {
	key := []byte("1234567812345678")
	data := []byte(`{"uid":123,"sid":"xxx","content":"hello,world"}`)
	encrpted, err := AESCBCEncrypt(key, nil, data)
	assert.Nil(t, err)
	t.Logf("%0x encrpted:%v, base64:%v", data, encrpted, base64.StdEncoding.EncodeToString(encrpted))
	decrypted, err := AESCBCDecrypt(key, nil, encrpted)
	assert.Nil(t, err)
	assert.Equal(t, string(data), string(decrypted))
}
