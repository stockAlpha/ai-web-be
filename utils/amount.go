package utils

func GetAmount(rechargeType int) int {
	amount := 0
	switch rechargeType {
	case 1:
		amount = 100
	case 2:
		amount = 500
	case 3:
		amount = 2000
	case 4:
		amount = 10000
	default:
		amount = 50
	}
	return amount
}
