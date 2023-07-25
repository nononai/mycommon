package uniqueid

import (
	"fmt"
	"testing"
)

func TestGenSn(t *testing.T) {
	sn := GenSn("TEST")
	fmt.Print(sn)

}
