package class01

import (
	"fmt"
	"testing"
)

func Test_001(t *testing.T) {
	initSwordDamage := 74
	milleLames(11, initSwordDamage, 8)
	//8967
	//8037 930
	//6913 1124
	//5614 1299 3353
	//4137 1477 44%
	//3068
}

func milleLames(initVulnerableLevel int, initSwordDamage int, milleLamesCount int) {
	stageDamage := 0
	totalDamage := 0
	vulnerableLevel := initVulnerableLevel
	for i := 0; i < milleLamesCount; i++ {
		stageDamage = 0
		for j := 0; j < 8; j++ {
			stageDamage += (initSwordDamage + vulnerableLevel) + (vulnerableLevel + 14)
			vulnerableLevel++
		}

		ratio := 100
		if totalDamage > 0 {
			ratio = stageDamage * 100 / totalDamage
		}
		fmt.Printf("stage: %d, stage damage: %d, Ratio, %d%%, ", i+1, stageDamage, ratio)
		totalDamage += stageDamage
		fmt.Printf("total damage: %d \n", totalDamage)
	}
}
