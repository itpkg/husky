package crypto_test

import (
	"crypto/aes"
	"testing"
	"time"

	"github.com/itpkg/husky/web/crypto"
)

const hello = "Hello, Husky!"

type S struct {
	I int
	S string
	T time.Time
}

func TestSerial(t *testing.T) {
	s := S{I: 123, S: hello, T: time.Now()}
	sr := crypto.Serial{}
	buf, err := sr.To(&s)
	if err != nil {
		t.Fatal(err)
	}
	var s1 S
	if err := sr.From(buf, &s1); err != nil {
		t.Fatal(err)
	}
	t.Logf("Get %+v", s1)
	if s1.I != s.I || s1.S != s.S {
		t.Fatalf("want %+v, get %+v", s, s1)
	}
}

func TestAes(t *testing.T) {
	key := []byte("AES256Key-32Characters1234567890")
	cip, e := aes.NewCipher(key)
	if e != nil {
		t.Fatal(e)
	}
	p := crypto.Encryptor{Cip: cip}

	if buf, err := p.Encode([]byte(hello)); err == nil {
		if s, err := p.Decode(buf); err == nil {
			t.Logf("%x, %s", buf, s)
			if string(s) != hello {
				t.Fatalf("Want %s, get %s", hello, s)
			}
		} else {
			t.Fatal(e)
		}
	} else {
		t.Fatal(e)
	}
}
