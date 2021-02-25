package public

import (
	"testing"
)

func TestAesAuth_EncryptCFB(t *testing.T) {
	aas := NewAesAuth()
	//ens := aas.AesEncryptCFB([]byte("Servy0u"))
	des,err :=aas.AesDecryptCFB("516f474561596b4e544152424e7576461bc9c803b05c6a")
	if err!=nil{
		t.Error(err)
	}
	t.Log(des)
	//t.Log(ens)
}

//func TestAesAuth_DecrypterCFB(t *testing.T) {
//	//aas := NewAesAuth()
//	//aas.AesDecryptCFB()
//
//}

