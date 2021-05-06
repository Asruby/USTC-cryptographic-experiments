package main

import (
	"encoding/hex"
	"fmt"
	"github.com/go-yaml/yaml"
	"io/ioutil"
	"log"
	"strconv"
)
/*
	The encryption technology is to use the exclusive OR operation of the streaming Krabi technology to encrypt.
	Every encryption is used the same stream encryption key.
	Thus space^key^key^[A-Z]->[a-z],space^key^key^[a-z]->[A-Z],space^key^key^space->0，We can gusses the space position
	Then space^Ciphertext->key ,we get one bit key of keys
	If we have enough space position,we can get the whole bits of keys
 */
func main()  {
	m:=map[interface{}]interface{}{}
	//read from file and Unmarshal
	err:=yaml.Unmarshal(readfromFile(`E:\master course\密码学\第一次实验\ciphertext.txt`),&m)
	//fmt.Println(m)
	if err!=nil{
		log.Panic("Unmarshal failed")
	}
	arr:=make([][]byte,len(m))

	//change hex-string to ascii
	for i:=1;i<=len(m)-1;i++{
		arr[i-1],_=hex.DecodeString(m["Ciphertext#"+strconv.Itoa(i)].(string))
	}
	arr[len(m)-1],_=hex.DecodeString(m["TargetCiphertext(DecryptThisOne)"].(string))

	//get everyline's possible space position
	//
	space_position:=findSpance(arr)
	//fmt.Println(space_position)
	//to generate the key space^space^key==key
	key:=[200]byte{}
	for i,sp:=range space_position{
		for _,idx:=range sp{
			key[idx]=arr[i][idx]^' '
		}
	}

	result:=make([][]byte,len(m))
	for i:=0;i<len(arr);i++{
		for j:=0;j<len(arr[i]);j++{
			letter:=arr[i][j]^key[j]
			if letter>=' ' && letter<='~'{
				result[i]=append(result[i],letter)
			}else{
				result[i]=append(result[i],' ')
			}
		}
	}
	for _,v:=range result{
		fmt.Println(string(v))
	}

}
func findSpance(arr [][]byte) [][]int {
	space_possible:=make([][][]int,len(arr))
	for i,v1:=range arr{
		space_possible[i]=make([][]int,0,len(arr)-1)
		for j,v2:=range arr {
			if j == i {
				continue
			}
			len1, len2 := len(v1), len(v2)
			minlen := min(len1, len2)
			tempIndex := []int{}
			for k := 0; k < minlen; k++ {
				if (v1[k]^v2[k] >= 'A' && v1[k]^v2[k] <= 'Z') || (v1[k]^v2[k] >= 'a' && v1[k]^v2[k] <= 'z') || v1[k]^v2[k]==0{
					tempIndex = append(tempIndex, k)
				}
			}
			space_possible[i] = append(space_possible[i], tempIndex)
		}
	}
	space_position:=make([][]int,len(arr))
	for i,v:=range space_possible{
		space_position[i]=[]int{}
		for j:=0;j<len(arr[i]);j++{
			count:=0
			for _,p:=range v{
				if isExist(p,j){
					count++
				}
			}
			if count>=8{
				space_position[i]=append(space_position[i],j)
			}
		}
	}
	return space_position
}
func isExist(poss[]int,pos int) bool{
	for _,v:=range poss{
		if v==pos{
			return true
		}
	}
	return false
}
func readfromFile(path string) []byte {
	buf:=[]byte{}
	buf,err:=ioutil.ReadFile(path)
	if err!=nil{
		panic("read file error")
	}
	return buf
}
func min(a,b int) int{
	if a<b{
		return a
	}
	return b
}
