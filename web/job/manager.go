package job

import (
	"fmt"

	"github.com/garyburd/redigo/redis"
	"github.com/itpkg/husky/web/crypto"
	"github.com/op/go-logging"
	"github.com/satori/go.uuid"
)

//QUEUE queue list name
const QUEUE = "queues"

//SUCCESS  success list name
const SUCCESS = "success"

//FAILED  success list name
const FAILED = "failed"

//Manager job manager
type Manager struct {
	Redis  *redis.Pool     `inject:""`
	Serial *crypto.Serial  `inject:""`
	Logger *logging.Logger `inject:""`
	Prefix string          `inject:"job.prefix"`
}

//Task model
type Task struct {
	ID   string
	Name string
	Args []interface{}
}

//Push push a task to queue
func (p *Manager) Push(name string, args ...interface{}) error {
	b, e := p.Serial.To(Task{ID: uuid.NewV4().String(), Name: name, Args: args})
	if e != nil {
		return e
	}
	c := p.Redis.Get()
	defer c.Close()

	_, e = c.Do("LPUSH", p.queue(QUEUE), b)
	return e
}

//Do run job
func (p *Manager) Do(n uint) {
	for {
		t, e := p.pop(n)
		if e == nil {
			p.run(t)
		} else {
			p.Logger.Error(e)
		}
	}
}

func (p *Manager) run(t *Task) {
	c := p.Redis.Get()
	defer c.Close()
	hd := handlers[t.Name]

	if hd == nil {
		c.Do("SADD", p.queue(FAILED), fmt.Sprintf("%s@%s: not found", t.ID, t.Name))
	} else {
		if err := hd(t.Args...); err == nil {
			c.Do("SADD", p.queue(SUCCESS), fmt.Sprintf("%s@%s", t.ID, t.Name))
		} else {
			c.Do("SADD", p.queue(FAILED), fmt.Sprintf("%s(%+v): %v", t.Name, t.Args, err))
		}
	}

}

func (p *Manager) pop(n uint) (*Task, error) {
	c := p.Redis.Get()
	defer c.Close()
	b, e := redis.ByteSlices(c.Do("BRPOP", p.queue(QUEUE), n))
	if e != nil {
		return nil, e
	}
	var t Task
	p.Serial.From(b[1], &t)
	return &t, e
}

func (p *Manager) queue(n string) string {
	return fmt.Sprintf("%s://%s", p.Prefix, n)
}
