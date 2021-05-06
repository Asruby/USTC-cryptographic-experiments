package exp5

import (
	"crypto/sha256"
	"encoding/hex"
	"io/ioutil"
)

const BLOCKSIZE = 1024
func verifyVedio(path string) (string,[][]byte){
	//读取视频字节流
	videoByte:=readfromFile(path)
	videoBlock:=make([][]byte,len(videoByte)/BLOCKSIZE+1)
	//字节流分块，最后一块可能不满1024
	for i,idx:=0,0;i<len(videoByte);idx++{
		videoBlock[idx]=videoByte[i:min(i+BLOCKSIZE,len(videoByte))]
		i=i+BLOCKSIZE
	}
	//从后往前算sha256哈希
	temp:=[32]byte{}
	for i:=len(videoBlock)-1;i>=0;i--{
		if i!=len(videoBlock)-1{
			videoBlock[i]=append([]byte{},videoBlock[i]...)
			videoBlock[i]=append(videoBlock[i],temp[:]...)
		}
		temp=sha256.Sum256(videoBlock[i])
	}
	return hex.EncodeToString(temp[:]),videoBlock
}
//验证
func check(h string,videoBlock [][]byte)bool{
	for i:=0;i<len(videoBlock);i++{
		temp:=sha256.Sum256(videoBlock[i])
		if h!=hex.EncodeToString(temp[:]){
			return false
		}
		h=hex.EncodeToString(videoBlock[i][len(videoBlock[i])-32:])
	}
	return true
}
func readfromFile(path string) []byte {
	buf:=[]byte{}
	buf,err:=ioutil.ReadFile(path)
	if err!=nil{
		panic("read file error")
	}
	return buf
}
func min(a,b int)int{
	if a>b{
		return b
	}
	return a
}
