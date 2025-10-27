package dbcache

import (
	"encoding/json"
	"goRelay/pkg"
	"os"
	"sync"
	"time"
)

type Caches struct {
	Cache            map[string]string `json:"cache"`
	SetCount         uint              `json:"set_count"`
	LastSetUnixTimes int64             `json:"last_set_unix_times"`
}

var (
	lock          sync.RWMutex
	saveLock      sync.Mutex
	maxSaveCount  = 600
	MaxSaveSecond = 600
	goLog         *pkg.Logger
)

var DBFilename string

func (c *Caches) Set(k, v string) {
	lock.Lock()
	defer lock.Unlock()

	c.Cache[k] = v
	c.SetCount += 1
	c.LastSetUnixTimes = time.Now().Unix()

	go func() {
		if c.SetCount >= uint(maxSaveCount) {
			c.Save()
		}
	}()

}

func (c *Caches) Get(k string) (string, bool) {
	lock.Lock()
	defer lock.Unlock()

	v, isok := c.Cache[k]
	return v, isok
}

func (c *Caches) Del(k string) {
	lock.Lock()
	defer lock.Unlock()

	delete(c.Cache, k)
}

func (c *Caches) List() map[string]string {
	lock.Lock()
	defer lock.Unlock()

	return c.Cache
}

func (c *Caches) Exists(k string) bool {
	lock.Lock()
	defer lock.Unlock()

	_, isok := c.Cache[k]
	return isok
}

func (c *Caches) Save() {
	c.SetCount = 0

	saveLock.Lock()
	defer saveLock.Unlock()

	dataBytes, err := json.Marshal(c)
	if err != nil {
		goLog.Error("dbcache save error ", err)
		return
	}

	fileInfo, err := os.Create(DBFilename)
	if err != nil {
		goLog.Error("dbcache save error ", err)
		return
	}
	defer fileInfo.Close()

	_, err = fileInfo.Write(dataBytes)
	if err != nil {
		goLog.Error("dbcache save error ", err)
		return
	}
	goLog.Info("dacache save successful")
}

func (c *Caches) load() error {
	dataBytes, err := os.ReadFile(DBFilename)
	if err != nil {
		goLog.Error("dbcache load db error ", err)
		c.Cache = make(map[string]string, 0)
		return nil
	}

	err = json.Unmarshal(dataBytes, &c)
	if err != nil {
		goLog.Error("dbcache load db error ", err)
		return err
	}
	return nil
}

func (c *Caches) AutoSave() error {

	for {
		// goLog.Debug("c.LastSetUnixTimes: ", c.LastSetUnixTimes, " c.SetCount:", c.SetCount, " auto save: ", c.LastSetUnixTimes != 0 && c.SetCount > 0 && time.Now().Unix() > c.LastSetUnixTimes+int64(MaxSaveSecond))
		if c.LastSetUnixTimes != 0 && c.SetCount > 0 && time.Now().Unix() > c.LastSetUnixTimes+int64(MaxSaveSecond) {
			c.Save()
			c.LastSetUnixTimes = time.Now().Unix()
			c.SetCount = 0
		}
		time.Sleep(3 * time.Second)
	}
}

func Init(dbFilename string) *Caches {
	var c Caches

	DBFilename = dbFilename

	c.load()

	return &c
}

func init() {
	goLog = pkg.NewLogger()
	goLog.SetLogger(pkg.ErrorLevel)
}
