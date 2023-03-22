package utility

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/tigerbeetledb/tigerbeetle-go/pkg/types"
)

func NewTransactionID() types.Uint128 {
	rand.Seed(time.Now().UnixNano())
	x, err := types.HexStringToUint128(fmt.Sprintf("%d", rand.Uint32()))
	if err != nil {
		log.Println(err)
	}

	return x
}

func Uint128(value string) types.Uint128 {
	x, err := types.HexStringToUint128(value)
	if err != nil {
		panic(err)
	}
	return x
}
