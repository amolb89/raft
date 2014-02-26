package clust

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"time"
)

type Config struct {
	Servers map[string]string
}

func LoadConfig(configFile string, c *Config) error {
	fileBytes, _ := ioutil.ReadFile(configFile)
	err := json.Unmarshal(fileBytes, c)
	return err
}

//Random number generation function using given range (source : http://golangcookbook.blogspot.in/2012/11/generate-random-number-in-given-range.html)
func Random(min, max int) int {
	if max == min {
		return max
	}
	rand.Seed(time.Now().UTC().UnixNano())
	return rand.Intn(max-min) + min
}
