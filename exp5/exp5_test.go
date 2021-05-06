package exp5

import (
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestExp5(t *testing.T)  {
	hash,videoBlock:=verifyVedio(`E:\master course\密码学\实验五\test.mp4`)
	log.Printf("chain hash is %s\n",hash)
	assert.True(t,check(hash,videoBlock))
	assert.Equal(t,"03c08f4ee0b576fe319338139c045c89c3e8e9409633bea29442e21425006ea8",hash)
}
