package generator

import (
	"fmt"
	"strings"
)

func (i *Insn) genCodeVdVs2Vs1(pos int) []string {
	sha2 := i.isExtension(crypto) && strings.HasPrefix(i.Name, "vsha")

	lmuls := iff(i.isExtension(crypto), cryptoLMULS, allLMULs)
	sews  := iff(i.isExtension(crypto), []SEW{SEW(32)}, allSEWs)
	sews  = iff(i.isExtension(crypto) && sha2, []SEW{SEW(32), SEW(64)}, sews)

	combinations := i.combinations(lmuls, sews, []bool{false})
	res := make([]string, 0, len(combinations))
	for _, c := range combinations[pos:] {
		if i.isExtension(crypto) && ((c.Vl % egs(i) != 0) || (c.Vl * int(i.Option.VLEN) < 4 * int(c.SEW))) {
			res = append(res, "")
			continue
		}

		builder := strings.Builder{}
		builder.WriteString(c.comment())

		vd := int(c.LMUL1)
		vss := []int{2 * int(c.LMUL1), 3 * int(c.LMUL1)}
		builder.WriteString(i.gWriteRandomData(c.LMUL1))
		builder.WriteString(i.gLoadDataIntoRegisterGroup(vd, c.LMUL1, SEW(8)))

		for idx, vs := range vss {
			builder.WriteString(i.gWriteIntegerTestData(c.LMUL1, c.SEW, idx))
			builder.WriteString(i.gLoadDataIntoRegisterGroup(vs, c.LMUL1, c.SEW))
		}

		builder.WriteString("# -------------- TEST BEGIN --------------\n")
		builder.WriteString(i.gVsetvli(c.Vl, c.SEW, c.LMUL))
		builder.WriteString(fmt.Sprintf("%s v%d, v%d, v%d\n",
			i.Name, vd, vss[1], vss[0]))
		builder.WriteString("# -------------- TEST END   --------------\n")

		builder.WriteString(i.gResultDataAddr())
		builder.WriteString(i.gStoreRegisterGroupIntoResultData(vd, c.LMUL1, c.SEW))
		builder.WriteString(i.gMagicInsn(vd))

		res = append(res, builder.String())
	}

	return res
}
