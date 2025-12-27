package proto

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseCommand(t *testing.T) {
	cmd := &CommandSet{
		Key:   []byte("foo"),
		Value: []byte("bar"),
		TTL:   2,
	}
	r := bytes.NewReader(cmd.Bytes())
	pcmd, err := ParseCommand(r)
	assert.Nil(t, err)
	assert.Equal(t, cmd, pcmd)
}

func TestGetCommand(t *testing.T) {
	cmd := &CommandGet{
		Key: []byte("foo"),
	}
	r := bytes.NewReader(cmd.Bytes())
	pcmd, err := ParseCommand(r)
	assert.Nil(t, err)
	assert.Equal(t, cmd, pcmd)
}
