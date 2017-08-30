package wxpay

import (
	"fmt"
	"testing"
)

func TestPayDataBase_ToXml(t *testing.T) {
	data := NewPayDataBase()
	data.SetMchId("acb")
	data.SetTotalFee(111)

	d := NewPayDataBase()

	ss, _ := data.ToXml()
	fmt.Println(ss)
	err := d.Init([]byte(ss))
	fmt.Println(err)
	fmt.Println(d.Values)
}
