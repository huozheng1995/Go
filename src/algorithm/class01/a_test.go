package class01

import (
	"fmt"
	"testing"
)

func Test_001(t *testing.T) {
	initSwordDamage := 15
	milleLames(35, initSwordDamage, 8)
}

func milleLames(initVulnerableLevel int, initSwordDamage int, milleLamesCount int) {
	stageDamage := 0
	totalDamage := 0
	for i := 0; i < milleLamesCount; i++ {
		stageDamage = 0
		for j := 0; j < 8; j++ {
			stageDamage += initVulnerableLevel + initSwordDamage + i*8 + j
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
