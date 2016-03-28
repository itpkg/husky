package job_test

import (
	"log"
	"testing"

	"github.com/garyburd/redigo/redis"
	"github.com/itpkg/husky/web/crypto"
	"github.com/itpkg/husky/web/job"
	"github.com/op/go-logging"
)

var mg = job.Manager{
	Redis: &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "localhost:6379")
		},
	},
	Serial: &crypto.Serial{},
	Prefix: "test-jobs",
	Logger: logging.MustGetLogger("test"),
}

func TestPush(t *testing.T) {

	for i := 0; i < 10; i++ {
		if e := mg.Push("echo", "hello", i); e != nil {
			t.Fatal(e)
		}
	}
}

func TestPop(t *testing.T) {
	job.Register("echo", func(args ...interface{}) error {
		log.Printf("echo %s %d", args[0].(string), args[1].(int))
		return nil
	})
	mg.Do(2)
}
