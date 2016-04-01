package core

import (
	"fmt"
	"os"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/chonglou/husky/api/core"
	"github.com/garyburd/redigo/redis"
	"github.com/jinzhu/gorm"
	"gopkg.in/olivere/elastic.v3"
)

//Config 配置
type Config struct {
	Env           string        `toml:"-"`
	Secrets       string        `toml:"secrets"`
	HTTP          HTTP          `toml:"http"`
	Database      Database      `toml:"database"`
	Redis         Redis         `toml:"redis"`
	ElasticSearch ElasticSearch `toml:"elastic_search"`
	Workers       Workers       `toml:"workers"`
}

//Home home url
func (p *Config) Home() string {
	if p.IsProduction() {
		if p.HTTP.Ssl {
			return fmt.Sprintf("https://%s", p.HTTP.Host)
		}
		return fmt.Sprintf("http://%s", p.HTTP.Host)

	}
	if p.HTTP.Ssl {
		return fmt.Sprintf("https://%s:%d", p.HTTP.Host, p.HTTP.Port)
	}
	return fmt.Sprintf("http://%s:%d", p.HTTP.Host, p.HTTP.Port)

}

//IsProduction production mode ?
func (p *Config) IsProduction() bool {
	return p.Env == "production"
}

//Key get key bytes
func (p *Config) Key(i, l int) ([]byte, error) {
	buf, err := core.FromBase64(p.Secrets)
	if err != nil {
		return nil, err
	}
	return buf[i : i+l], nil
}

//-----------------------------------------------------------------------------

//HTTP 配置信息
type HTTP struct {
	Host string `toml:"host"`
	Port int    `toml:"port"`
	Ssl  bool   `toml:"ssl"`
}

//-----------------------------------------------------------------------------

//Redis 配置信息
type Redis struct {
	Host string `toml:"host"`
	Port int    `toml:"port"`
	Db   int    `toml:"db"`
}

//URL 连接
func (p *Redis) URL() string {
	return fmt.Sprintf("%s:%d", p.Host, p.Port)
}

//Open 打开
func (p *Redis) Open() *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, e := redis.Dial("tcp", p.URL())
			if e != nil {
				return nil, e
			}
			if _, e = c.Do("SELECT", p.Db); e != nil {
				c.Close()
				return nil, e
			}
			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

//-----------------------------------------------------------------------------

//Workers 任务
type Workers struct {
	ID     string         `toml:"id"`
	Pool   int            `toml:"pool"`
	Queues map[string]int `toml:"queues"`
}

//-----------------------------------------------------------------------------

//Database 配置信息
type Database struct {
	Type string            `toml:"type"`
	Args map[string]string `toml:"args"`
}

//Execute 执行sql
func (p *Database) Execute(sql string) (string, []string) {
	switch p.Type {
	case "postgres":
		cmd, args := p.psql()
		args = append(args, "-c", sql)
		return cmd, args
	default:
		return p.bad()
	}
}

func (p *Database) bad() (string, []string) {
	return "echo", []string{fmt.Sprintf("Unsupported database adapter %s", p.Type)}
}
func (p *Database) psql() (string, []string) {
	args := []string{"-U", p.Args["user"]}
	if host, ok := p.Args["host"]; ok {
		args = append(args, "-h", host)
	}
	if port, ok := p.Args["port"]; ok {
		args = append(args, "-p", port)
	}
	return "psql", args
}

//Console 控制台
func (p *Database) Console() (string, []string) {
	switch p.Type {
	case "postgres":
		cmd, args := p.psql()
		args = append(args, p.Args["dbname"])
		return cmd, args
	default:
		return p.bad()
	}
}

//Open 打开连接
func (p *Database) Open() (*gorm.DB, error) {
	//postgresql: "user=%s password=%s host=%s port=%d dbname=%s sslmode=%s"
	args := ""
	for k, v := range p.Args {
		args += fmt.Sprintf(" %s=%s ", k, v)
	}
	db, err := gorm.Open(p.Type, args)
	if err != nil {
		return nil, err
	}

	if err := db.DB().Ping(); err != nil {
		return nil, err
	}
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
	return db, nil
}

//-----------------------------------------------------------------------------

//ElasticSearch 配置信息
type ElasticSearch struct {
	Host  string `toml:"host"`
	Port  int    `toml:"port"`
	Index string `toml:"index"`
}

//URL 连接地址
func (p *ElasticSearch) URL() string {
	return fmt.Sprintf("http://%s:%d", p.Host, p.Port)
}

//Open 打开连接
func (p *ElasticSearch) Open() (*elastic.Client, error) {
	cli, err := elastic.NewClient(elastic.SetURL(p.URL()))
	if err == nil {
		return nil, err
	}
	//exi, err := cli.IndexExists(p.Index)
	return cli, err
}

//-----------------------------------------------------------------------------

//Load 加载配置文件
func Load(file string, obj interface{}) error {
	_, err := toml.DecodeFile(file, obj)
	return err
}

//Store 写入配置文件
func Store(file string, obj interface{}) error {
	fi, err := os.Create(file)
	defer fi.Close()

	if err == nil {
		end := toml.NewEncoder(fi)
		err = end.Encode(obj)
	}
	return err
}
