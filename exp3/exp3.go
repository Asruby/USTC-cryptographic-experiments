package exp3

import (
	"bytes"
	"crypto/aes"
	"encoding/hex"
	"log"
	"math/rand"
	"time"
)
const (
	//块大小128位
	BLOCKSIZE=16
)
//AES在CBS模式下加密
func AesCBCEncrypt(pText string,key string)string{
	pTextByte:=[]byte(pText)
	keyByte,_:=hex.DecodeString(key)
	//CBC模式下需要进行填充
	pTextByte=pkcs5Padding(pTextByte)
	cTextByte:=make([]byte,len(pTextByte)+BLOCKSIZE)
	//生成IV随机值，置于ETextByte前BLOCKSIZE位置处
	IV:=cTextByte[:BLOCKSIZE]
	rand.Seed(time.Now().Unix())
	if _,err:=rand.Read(IV);err!=nil{
		log.Panic(err.Error())
	}
	//NewCipher 创建并返回一个新的cipher.Block.
	block,_:=aes.NewCipher(keyByte)
	//VI值与输入异或后进行AES加密，VI值更新位加密后的值参与下一轮循环
	for c,i:=append([]byte{},IV...),0;i<len(pTextByte);i+=BLOCKSIZE{
		for j:=0;j<BLOCKSIZE;j++{
			c[j]=c[j]^pTextByte[i+j]
		}
		block.Encrypt(cTextByte[i+BLOCKSIZE:i+2*BLOCKSIZE],c)
		c=append([]byte{},cTextByte[i+BLOCKSIZE:i+2*BLOCKSIZE]...)
	}
	return hex.EncodeToString(cTextByte)
}
//AES在CBC模式下解密
func AesCBCDncrypt(cText string,key string)string{
	cTextByte,_:=hex.DecodeString(cText)
	keyByte,_:=hex.DecodeString(key)
	//获取保存在密文前128位的随机值
	IV:=cTextByte[:BLOCKSIZE]
	pTextByte:=make([]byte,len(cTextByte)-BLOCKSIZE)
	block,_:=aes.NewCipher(keyByte)
	//密文通过AES解密后异或VI得到明文，VI更新为此轮的密文参与下一轮循环
	for c,i:=append([]byte{},IV...),BLOCKSIZE;i<len(cTextByte);i+=BLOCKSIZE{
		block.Decrypt(pTextByte[i-BLOCKSIZE:i],cTextByte[i:i+BLOCKSIZE])
		for j:=0;j<BLOCKSIZE;j++{
			pTextByte[i-BLOCKSIZE+j]=pTextByte[i-BLOCKSIZE+j]^c[j]
		}
		c=append([]byte{},cTextByte[i:i+BLOCKSIZE]...)
	}
	return string(pkcs5Unpadding(pTextByte))
}
//AES在CTR模式下加密
func AesCTREncrypt(pText string,key string)string {
	pTextByte:=[]byte(pText)
	keyByte,_:=hex.DecodeString(key)
	//CTR模式下不需要进行填充
	cTextByte:=make([]byte,len(pTextByte)+BLOCKSIZE)
	//生成IV用于初始计数，置于ETextByte前BLOCKSIZE位置处
	IV:=cTextByte[:BLOCKSIZE]
	rand.Seed(time.Now().Unix())
	if _,err:=rand.Read(IV);err!=nil{
		log.Panic(err.Error())
	}
	block,_:=aes.NewCipher(keyByte)
	//对计数VI进行AES加密，加密结果异或明文得到密文，计数加1参与下一轮循环
	for cnt,i:=append([]byte{},IV...),0;i<len(pTextByte);i+=BLOCKSIZE{
		temp:=make([]byte,BLOCKSIZE)
		block.Encrypt(temp,cnt)
		for j:=0;j<BLOCKSIZE&&i+j<len(pTextByte);j++{
			cTextByte[i+BLOCKSIZE+j]=pTextByte[i+j]^temp[j]
		}
		addOne(cnt)
	}
	return hex.EncodeToString(cTextByte)
}
//AES在CTR模式下解密
func AesCTRDncrypt(cText string,key string)string{
	cTextByte,_:=hex.DecodeString(cText)
	keyByte,_:=hex.DecodeString(key)
	IV:=cTextByte[:BLOCKSIZE]
	pTextByte:=make([]byte,len(cTextByte)-BLOCKSIZE)
	block,_:=aes.NewCipher(keyByte)
	//对计数VI进行AES加密，加密结果异或密文得到明文，计数加1参与下一轮循环
	for cnt,i:=append([]byte{},IV...),BLOCKSIZE;i<len(cTextByte);i+=BLOCKSIZE{
		temp:=make([]byte,BLOCKSIZE)
		block.Encrypt(temp,cnt)
		for j:=0;j<BLOCKSIZE&&i+j<len(cTextByte);j++{
			pTextByte[i-BLOCKSIZE+j]=cTextByte[i+j]^temp[j]
		}
		addOne(cnt)
	}
	return string(pTextByte)
}


//PKCS5填充方案：PKCS5是8字节填充的，即填充一定数量的内容，使得成为8的整数倍,填充的内容为填充的长度。
//eg:数据数 0x41
//   填充前：0x410x410x410x410x41
//   填充后：0x410x410x410x410x410x030x030x03
func pkcs5Padding(source []byte) []byte{
	sLen:=len(source)
	pLen:=BLOCKSIZE-sLen%BLOCKSIZE
	pText:=bytes.Repeat([]byte{byte(pLen)},pLen)
	source=append(source,pText...)
	return source
}
func pkcs5Unpadding(source []byte) []byte{
	return source[:len(source)-int(source[len(source)-1])]
}
func addOne(cnt []byte){
	for i:=len(cnt)-1;i>=0;i--{
		if cnt[i]<255{
			cnt[i]+=1
			return
		}else {
			cnt[i]=0
		}
	}
	//TODO:加1溢出处理
	return
}


