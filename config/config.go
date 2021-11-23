package config

import (
	"encoding/json"
	"fmt"
	"github.com/philchia/agollo/v4"

	"log"
	"os"
	"strconv"
	"strings"
	"sync"
)

var (
	defaultConfig *Config
	prefix        = "properties"
	PROJECT_ID    = "norecord"
	BRANCH_NAME   = "dev"
	ApolloProxy   = "http://admin:123456@10.11.101.196:9999/fat"
	ConfigServer  = ApolloProxy + "/" + PROJECT_ID + "/" + BRANCH_NAME
)

func init() {
	apollo := os.Getenv("APOLLO")
	branchName := os.Getenv("CI_COMMIT_REF_NAME")
	projectId := os.Getenv("CI_PROJECT_ID")
	if apollo != "" {
		ApolloProxy = apollo
	}
	if branchName != "" {
		BRANCH_NAME = branchName
	}
	if projectId != "" {
		PROJECT_ID = projectId
	}
	ConfigServer = ApolloProxy + "/" + PROJECT_ID + "/" + BRANCH_NAME
	log.Println("ConfigServer: ", ConfigServer)
}

type Config struct {
	agollo.Conf
	data     map[string]string
	Callback func()
	lock     sync.RWMutex
}

func NewConfig(conf agollo.Conf, callbacks ...func()) *Config {

	if defaultConfig == nil {
		callback := func() {}
		if len(callbacks) > 0 {
			callback = callbacks[0]
		}
		defaultConfig = &Config{
			Conf:     conf,
			Callback: callback,
			data:     make(map[string]string),
		}
	}
	return defaultConfig
}
func (c *Config) LoadConfig() {
	c.Conf.MetaAddr = ConfigServer
	agollo.Start(&c.Conf)
	agollo.OnUpdate(c.updateCallback)
	c.lock.Lock()
	defer c.lock.Unlock()
	for _, ns := range c.Conf.NameSpaceNames {
		keys := agollo.GetAllKeys(agollo.WithNamespace(ns + "." + prefix))
		for _, key := range keys {
			c.data[ns+"."+key] = agollo.GetString(key, agollo.WithNamespace(ns))
		}
	}
}

func (c *Config) Json() {
	bt, _ := json.MarshalIndent(c.data, "", "    ")
	fmt.Println(string(bt))
}

func (c *Config) String() string {
	c.lock.RLock()
	defer c.lock.RUnlock()
	retVal := make([]string, 0)
	for k, v := range c.data {
		retVal = append(retVal, k+": "+v)
	}
	return strings.Join(retVal, "\n")
}

func (c *Config) updateCallback(event *agollo.ChangeEvent) {
	log.Println("updateCallback")
	ns := event.Namespace
	c.lock.Lock()
	defer func() {
		c.lock.Unlock()
		c.Callback()
	}()
	for _, prop := range event.Changes {
		c.data[ns+"."+prop.Key] = prop.NewValue
	}

}

func GetString(key string, defaultVal ...string) string {
	defaultConfig.lock.RLock()
	defer defaultConfig.lock.RUnlock()
	v, ok := defaultConfig.data[key]
	if !ok {
		if len(defaultVal) > 0 {
			return defaultVal[0]
		}
		return ""
	}
	return v

}

func GetBool(key string, defaultVal ...bool) bool {
	defaultConfig.lock.RLock()
	defer defaultConfig.lock.RUnlock()
	v, ok := defaultConfig.data[key]
	if !ok {
		if len(defaultVal) > 0 {
			return defaultVal[0]
		}
		return false
	}
	lv := strings.ToLower(v)
	if ret, err := strconv.ParseBool(lv); err != nil {
		return false
	} else {
		return ret
	}
}

func GetInt(key string, defaultVal ...int) int {
	defaultConfig.lock.RLock()
	defer defaultConfig.lock.RUnlock()
	v, ok := defaultConfig.data[key]
	if !ok {
		if len(defaultVal) > 0 {
			return defaultVal[0]
		}
		return 0
	}

	nv, err := strconv.Atoi(v)
	if err != nil {
		return 0
	}
	return nv

}
func GetFloat64(key string, defaultVal ...float64) float64 {
	defaultConfig.lock.RLock()
	defer defaultConfig.lock.RUnlock()
	v, ok := defaultConfig.data[key]
	if !ok {
		if len(defaultVal) > 0 {
			return defaultVal[0]
		}
		return 0
	}

	nv, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return 0
	}
	return nv

}
func GetIntSlice(key string, defaultVal ...[]int) []int {
	defaultConfig.lock.RLock()
	defer defaultConfig.lock.RUnlock()
	v, ok := defaultConfig.data[key]
	if !ok {
		if len(defaultVal) > 0 {
			return defaultVal[0]
		}
		return nil
	}
	retVal := make([]int, 0)

	err := json.Unmarshal([]byte(v), &retVal)
	if err != nil {
		return nil
	}
	return retVal

}

func GetStringSlice(key string, defaultVal ...[]string) []string {
	defaultConfig.lock.RLock()
	defer defaultConfig.lock.RUnlock()
	v, ok := defaultConfig.data[key]
	if !ok {
		if len(defaultVal) > 0 {
			return defaultVal[0]
		}
		return nil
	}
	retVal := make([]string, 0)

	err := json.Unmarshal([]byte(v), &retVal)
	if err != nil {
		return nil
	}
	return retVal
}

func GetStruct(key string, output interface{}) bool {
	defaultConfig.lock.RLock()
	defer defaultConfig.lock.RUnlock()
	v, ok := defaultConfig.data[key]
	if !ok {
		return false
	}

	err := json.Unmarshal([]byte(v), output)
	if err != nil {
		return false
	}
	return true
}
