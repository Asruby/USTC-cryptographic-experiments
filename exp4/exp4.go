package exp4

import (
	"math/big"
)

var times = new(big.Int).Lsh(big.NewInt(1),20)
func task1(N string)(string,string){
	bigN:=new(big.Int)
	bigN.SetString(N,10)
	//A=ceil(sqrt(N))
	//Sqrt sets z to ⌊√x⌋, the largest integer such that z² ≤ x, and returns z.
	A:=new(big.Int).Sqrt(bigN)
	if new(big.Int).Mul(A,A).Cmp(bigN)==-1{
		A.Add(A,big.NewInt(1))
	}
	//x=sqrt(A^2-N)
	x:=new(big.Int).Sqrt(new(big.Int).Sub(new(big.Int).Mul(A,A),bigN))
	//p=A-x,q=A+x
	p:=new(big.Int).Sub(A,x)
	q:=new(big.Int).Add(A,x)
	return p.String(),q.String()
}

func task2(N string)(string,string)  {
	bigN:=new(big.Int)
	bigN.SetString(N,10)
	//A=ceil(sqrt(N))
	//Sqrt sets z to ⌊√x⌋, the largest integer such that z² ≤ x, and returns z.
	A:=new(big.Int).Sqrt(bigN)
	if new(big.Int).Mul(A,A).Cmp(bigN)==-1{
		A.Add(A,big.NewInt(1))
	}
	//从sqrt(N)开始向上搜索A，直到成功分解N
	for i:=big.NewInt(0);i.Cmp(times)==-1;i.Add(i,big.NewInt(1)){
		x:=new(big.Int).Sqrt(new(big.Int).Sub(new(big.Int).Mul(A,A),bigN))
		if new(big.Int).Mul(x,x).Cmp(new(big.Int).Sub(new(big.Int).Mul(A,A),bigN))==0{
			p:=new(big.Int).Sub(A,x)
			q:=new(big.Int).Add(A,x)
			return p.String(),q.String()
		}
		A.Add(A,big.NewInt(1))
	}
	return "",""
}
