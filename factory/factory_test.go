package factory

import (
	"avilego.me/recent_news/env"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

func TestGetLogger(t *testing.T) {
	os.Remove(env.LogFile())
	logger := logger()

	logger.Print("something went wrong")

	logText := mustReadFile(env.LogFile())
	assert.FileExists(t, env.LogFile())
	assert.Contains(t, logText, "something went wrong")
}

func mustReadFile(file string) string {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
	return string(content)
}
