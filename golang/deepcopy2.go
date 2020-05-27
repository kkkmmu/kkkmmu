package main

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
)

type DeepCopy struct {
	A map[string]string
	B []string
	C Cc
}

type Cc struct {
	D map[int]*string
}

func main() {
	var F string = `fffffff`
	var E = map[int]*string{0: &F}
	var cc Cc = Cc{
		D: E,
	}
	var dc DeepCopy = DeepCopy{
		A: make(map[string]string, 0),
		B: []string{`b`, `c`, `d`},
		C: cc,
	}
	dv := DeepCopy{}
	dv = dc
	fmt.Printf("dc : %p\n", &dc)
	fmt.Printf("dv: %p\n", &dv)
	dv.B = []string{`f`, `g`, `d`}
	var K string = `kkkkkkk`
	dv.C.D[0] = &K
	fmt.Println(`dc:`, dc)
	fmt.Println(`dv:`, dv)
	fmt.Println(`df.C.D:`, *(dc.C.D[0]))
	fmt.Println(`dn.C.D:`, *(dv.C.D[0]))
	fmt.Println(`=============================`)
	var df *DeepCopy = &DeepCopy{
		A: make(map[string]string, 0),
		B: []string{`b`, `c`, `d`},
		C: cc,
	}
	dn := &DeepCopy{}
	dn = df // or *dn = *df
	fmt.Printf("dn : %p\n", dn)
	fmt.Printf("df: %p\n", df)
	dn.B = []string{`f`, `g`, `d`}
	var G string = `ggggggg`
	dn.C.D[0] = &G
	fmt.Println(`df:`, df)
	fmt.Println(`dn:`, dn)
	fmt.Println(`df.C.D:`, *(df.C.D[0]))
	fmt.Println(`dn.C.D:`, *(dn.C.D[0]))
	fmt.Println(`=============================`)
	var buf bytes.Buffer
	dg := &DeepCopy{}
	json.NewEncoder(&buf).Encode(*df)
	json.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(dg)
	dg.B = []string{`h`, `j`, `k`}
	fmt.Println(dg.C.D)
	var W string = `wwwwwww`
	dg.C.D[0] = &W
	fmt.Printf("df : %p\n", df)
	fmt.Printf("dg: %p\n", dg)
	fmt.Println(`df:`, df)
	fmt.Println(`dg:`, dg)
	fmt.Println(`df.C.D:`, *(df.C.D[0]))
	fmt.Println(`dn.C.D:`, *(dg.C.D[0]))
	fmt.Println(`=============================`)
	var buff bytes.Buffer
	var FF string = `fffffff`
	var EE = map[int]*string{0: &FF}
	var cd Cc = Cc{
		D: EE,
	}
	var dx *DeepCopy = &DeepCopy{
		A: make(map[string]string, 0),
		B: []string{`b`, `c`, `d`},
		C: cd,
	}
	dj := &DeepCopy{}
	gob.NewEncoder(&buff).Encode(*dx)
	gob.NewDecoder(bytes.NewBuffer(buff.Bytes())).Decode(dj)
	dj.B = []string{`q`, `w`, `e`}
	var S string = `qqqqqqq`
	dj.C.D[0] = &S
	fmt.Printf("dx : %p\n", dx)
	fmt.Printf("dj: %p\n", dj)
	fmt.Println(`dx:`, dx)
	fmt.Println(`dj:`, dj)
	fmt.Println(`dx.C.D:`, *(dx.C.D[0]))
	fmt.Println(`dj.C.D:`, *(dj.C.D[0]))
	fmt.Println(`=============================`)
}
