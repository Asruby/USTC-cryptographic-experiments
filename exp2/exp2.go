package main

import (
	"context"
	"fmt"
	"github.com/go-yaml/yaml"
	"io/ioutil"
	"math/big"
	"sync"
	"time"
)
//B is 2^20
var B =big.NewInt(1)
//GB is g^(2^20)
var GB = big.NewInt(1)

func init()  {
	B=(new(big.Int)).Lsh(B,20)
}
func main()  {
	ctx,cancel:=context.WithCancel(context.Background())
	wg:=&sync.WaitGroup{}
	wg.Add(1)
	go spinner(ctx,wg)
	m:=map[string]string{}
	// read from file
	err:=yaml.Unmarshal(readfromFile(`E:\master course\密码学\实验二\test.txt`),&m)
	if err!=nil{
		panic("fail to read file ")
	}
	p,g,h:=new(big.Int),new(big.Int),new(big.Int)
	p,_=p.SetString(m["p"],10)
	g,_=g.SetString(m["g"],10)
	h,_=h.SetString(m["h"],10)
	//g^x mod(p)=h  ====>    x=x0*B-x1     B=2^20, 1<=x0<=2^20, 0<=x1<=2^20
	//g^(x0*B)mod(p)=h*g^(x1)
	x0,x1:=SolutionOfDiscreteLogarithm(p,g,h)
	fmt.Printf("\rget x0=%s,x1=%s\n",x0,x1)
	//ans=x0*B-x1
	ans:=big.NewInt(1)
	if len(x0)!=0&&len(x1)!=0{
		x0,_:=new(big.Int).SetString(x0,10)
		x1,_:=new(big.Int).SetString(x1,10)
		ans.Mul(x0,B)
		ans.Sub(ans,x1)
	}
	fmt.Printf("\r%s\n",ans.String())
	cancel()
	wg.Wait()
}
func SolutionOfDiscreteLogarithm(p *big.Int,g *big.Int,h *big.Int) (string,string) {
	//计算等式右边，即h*g^x1,以h*g^x1为key，x1为value将结果保存到map中
	m:=map[string]string{}
	h.Mod(h,p)
	m[h.String()]="0"
	for i:=big.NewInt(1);i.Cmp(B)<=0;i.Add(i,big.NewInt(1)){
		h.Mul(h,g)
		h.Mod(h,p)
		m[h.String()]=i.String()
	}

	//计算GB=g^B
	for i:=big.NewInt(1);i.Cmp(B)<=0;i=i.Add(i,big.NewInt(1)){
		GB.Mul(GB,g)
		GB.Mod(GB,p)
	}
	if x1,ok:=m[GB.String()];ok{
		return "0",x1
	}
	//计算(g^B)^x1
	tempGB,_:=new(big.Int).SetString("1",10)
	for i:=big.NewInt(1);i.Cmp(B)<=0;i.Add(i,big.NewInt(1)){
		tempGB.Mul(tempGB,GB)
		tempGB.Mod(tempGB,p)
		if x1,ok:=m[tempGB.String()];ok{
			return i.String(),x1
		}
	}
	return "",""
}

func spinner(ctx context.Context,wg *sync.WaitGroup) {
	start:=time.Now().Unix()
	for  {
		select {
		case <-ctx.Done():
			fmt.Printf("\rspent time:%d\n",time.Now().Unix()-start)
			wg.Done()
			return
		default:
			for _, r := range `-\|/` {
				fmt.Printf("\r%c计算中...", r)
			}

		}
	}
}
func readfromFile(path string) []byte {
	buf:=[]byte{}
	buf,err:=ioutil.ReadFile(path)
	if err!=nil{
		panic("read file error")
	}
	return buf
}
