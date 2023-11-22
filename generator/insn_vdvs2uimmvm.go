package generator

import (
	"fmt"
	"math"
	"strings"
)

func (i *Insn) genCodeVdVs2UimmVm(pos int) []string {
	vs2Widening := strings.HasSuffix(i.Name, ".wi")
	sews := iff(vs2Widening, allSEWs[:len(allSEWs)-2], allSEWs[:len(allSEWs)-1])
	vs2Size := iff(vs2Widening, 2, 1)

	sews = iff(i.isExtension(crypto), allSEWs, sews)
	combinations := i.combinations(
		iff(vs2Widening, wideningMULs, allLMULs),
		sews,
		[]bool{false, true},
	)
	res := make([]string, 0, len(combinations))

	for _, c := range combinations[pos:] {
		builder := strings.Builder{}
		builder.WriteString(c.comment())

		builder.WriteString(i.gWriteRandomData(LMUL(1)))
		builder.WriteString(i.gLoadDataIntoRegisterGroup(0, LMUL(1), SEW(32)))

		vs2EMUL1 := LMUL(math.Max(float64(int(c.LMUL)*vs2Size), 1))
		vs2EEW := c.SEW * SEW(vs2Size)

		vd := int(c.LMUL1)
		vs2 := 2*int(c.LMUL1) + int(vs2EMUL1)
		builder.WriteString(i.gWriteRandomData(c.LMUL1))
		builder.WriteString(i.gLoadDataIntoRegisterGroup(vd, c.LMUL1, SEW(8)))

		builder.WriteString(i.gWriteIntegerTestData(vs2EMUL1, vs2EEW, 1))
		builder.WriteString(i.gLoadDataIntoRegisterGroup(vs2, vs2EMUL1, vs2EEW))

		cases := i.integerTestCases(c.SEW)
		for a := 0; a < len(cases); a++ {
			builder.WriteString("# -------------- TEST BEGIN --------------\n")
			builder.WriteString(i.gVsetvli(c.Vl, c.SEW, c.LMUL))
			switch c.SEW {
			case 8:
				builder.WriteString(fmt.Sprintf("%s v%d, v%d, %d%s\n",
					i.Name, vd, vs2, convNum[uint8](cases[a][0]), v0t(c.Mask)))
			case 16:
				builder.WriteString(fmt.Sprintf("%s v%d, v%d, %d%s\n",
					i.Name, vd, vs2, convNum[uint16](cases[a][0]), v0t(c.Mask)))
			case 32:
				builder.WriteString(fmt.Sprintf("%s v%d, v%d, %d%s\n",
					i.Name, vd, vs2, convNum[uint32](cases[a][0]), v0t(c.Mask)))
			case 64:
				builder.WriteString(fmt.Sprintf("%s v%d, v%d, %d%s\n",
					i.Name, vd, vs2, convNum[uint64](cases[a][0]), v0t(c.Mask)))
			}
			builder.WriteString("# -------------- TEST END   --------------\n")

			builder.WriteString(i.gResultDataAddr())
			builder.WriteString(i.gStoreRegisterGroupIntoResultData(vd, c.LMUL1, c.SEW))
			builder.WriteString(i.gMagicInsn(vd))
		}

		res = append(res, builder.String())
	}

	return res
}
