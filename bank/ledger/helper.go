package ledger

import tb_types "github.com/tigerbeetledb/tigerbeetle-go/pkg/types"

func uint128(value string) tb_types.Uint128 {
	x, err := tb_types.HexStringToUint128(value)
	if err != nil {
		panic(err)
	}
	return x
}

func check(flag, bitmask uint16) bool {
	return flag&bitmask != 0
}

func setFlag(flag, bitmsak uint16) uint16 {
	return flag | bitmsak
}

func isDebitBalanceFlag() uint16 {
	return tb_types.AccountFlags{CreditsMustNotExceedDebits: true}.ToUint16()
}

func isCreditBalanceFlag() uint16 {
	return tb_types.AccountFlags{DebitsMustNotExceedCredits: true}.ToUint16()
}

func isLinkedFlag() uint16 {
	return tb_types.AccountFlags{Linked: true}.ToUint16()
}
