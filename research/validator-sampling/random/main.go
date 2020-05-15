package main

import (
	"fmt"
	"math/rand"
	"sort"
)

var (
	t   = 1000000
	xxx = []string{"kavavaloper1qyc2cfl0nw8r95dsdw534x99wq0xcj9rmxpl7z", "kavavaloper1qfy0e2w62g6j4jg5djcqd4py3zsaeqexjplj2d", "kavavaloper1q3y9qga5hf360dmzta67vp54qz25tmv4hhkk4t", "kavavaloper1qk0pta4ga5t8p5vv7me8dz32lvcrv2rp098cas", "kavavaloper1ptyzewnns2kn37ewtmv6ppsvhdnmeapvl7z9xh", "kavavaloper1pn5k9c5pxmg5f0rycpl9rrx6k6mk85scxf06zx", "kavavaloper1phd8jz25lumudc7ac7rhmupvcqcv7lg3c8dprc", "kavavaloper1pceqe8we7drpfqrutchwy3f99800hhzuw6cc84", "kavavaloper1rgcgqmnkeffks7enrv6hk5u4wg3nzfkmqlzjqd", "kavavaloper1yrm63pqvld8uyzkavz55p2cktpm2gm8jd8xxlu", "kavavaloper1yjj2wfers947l6n5pynpgsqlz7svc5n8ssl6ye", "kavavaloper1yna6lete8nwwwctsalzdg04ldqaz73gtn33ydq", "kavavaloper196anr2ycsalg806dz29afklnecuaupvkh5qz6c", "kavavaloper1xzud24na9tucauc7lf6pjk84kqsgq3eq747ta7", "kavavaloper1x9hq3rjc48t5upcsr3c209ycgekfasne3l5nkc", "kavavaloper1xftqdxvq0xkv2mu8c5y0jrsc578tak4m9u0s44", "kavavaloper1xka27j0jvmq97yunj5wp8fv242lycmax8ejlaf", "kavavaloper1xhxzmj8fvkqn76knay9x2chfra826369dhdu2c", "kavavaloper18zksjhrefqew0zahmts894p8asscufxvdfq702", "kavavaloper18s9m5d5cjf0humjv7mkq8pm47kchwm0r0369cx", "kavavaloper18cf35l7req0k6ulqapeyv830mrrucn9xj87plr", "kavavaloper1g20mhvcpjxp6gzlwhtfcphjehwcl2njqydgu7q", "kavavaloper1gtd040dmljyaty9tkhq0mlqz7saer48v4d608x", "kavavaloper1g4qpetrj59a29e4wxpe74x93q4df2czjh8r9ak", "kavavaloper1ffcujj05v6220ccxa6qdnpz3j48ng024ykh2df", "kavavaloper1fw7vjc3fphahqxpdjypddlulnltxws8g0mrds7", "kavavaloper1fmas0qlsucg4qwf8mqyrylcg3uluz2ffg8952q", "kavavaloper12qn7y04wzr5s3h4dmdtre4q9f4nvc03a9a9qsz", "kavavaloper12r77hhj6ylvvl4etjm0fmpzh07jum8u7qqd695", "kavavaloper129kkennsm7na34lu6sn4kwxp7ewes58y4fx6y9", "kavavaloper12g40q2parn5z9ewh5xpltmayv6y0q3zs6ddmdg", "kavavaloper12et238paeqxyhvk2pfs20ygfj6ct0dx2ccsdz5", "kavavaloper1t5l8ht0wxpd4lpe4cweftrpg5kyn3qp437yvnr", "kavavaloper1vyx7wt8s8dwcspdt7dy49sq4l3jwyhxyndakmm", "kavavaloper1vfawvtzvmcjkpqzhezvhk9q5tv5t7x8smz95uu", "kavavaloper1vw35vclatlrcmzuaxf2fleyuk38xa7xf4vdaq2", "kavavaloper1vw4t8ge2ephu0wuhcmclcw04ag2vzj9rpdme2x", "kavavaloper1vuylvflgy75d8zr07ta8x0gcqynvkvw70et853", "kavavaloper1dwae0ny4uuvacucakm9v8r8mxhw82zack4cn7y", "kavavaloper1dj63l3z0smqa8au5yek37nmj0xd60zs9dwmdh0", "kavavaloper1dntlhrw3jrej6ssdp64yfmkkz08ykyx4n7hphh", "kavavaloper1dede4flaq24j2g9u8f83vkqrqxe6cwzrxt5zsu", "kavavaloper1davtd9w5yadvlg5x7aw0kpkyhckk5m5ecrmnyl", "kavavaloper1wtcn3ylrsp4qp4urlhkf78kkt8f2ch7p0f0p98", "kavavaloper1wu8m65vqazssv2rh8rthv532hzggfr3h9azwz9", "kavavaloper1wacmlkst8s39fq83tqhlcuuacts2z6lvwfq8hv", "kavavaloper1006e7kwuhthdtxe0d90pc209cy803ehem0dg0q", "kavavaloper10m3hjapny44txmgr47rf277364htgqpr646cty", "kavavaloper1spkjjtwks5zt0dqexj8rv8ljwmwxe8ufkraukq", "kavavaloper1srhded4xw0krcwlvddtcyycuh70fp5ry9yvp86", "kavavaloper1s8akp0nq7z7vuf5g9agdswcwq7vm9uwudk0han", "kavavaloper13fxkk4730cqglgdv7w0mdelyx07myyq76uf9f3", "kavavaloper13vstf6ecmfe4p0gaufumk569sdawhtrf8gu56h", "kavavaloper13dfu6et0m8zm4hachudvn62w0c2zvcmrqn8cs3", "kavavaloper1jyuv7z9at27elvmnmzh2v39dc06r9kjcy59xkr", "kavavaloper1j26c4k2jj9tv95whdhva3e8v2fcm4s3dsgstd2", "kavavaloper1jhraz4ftxl2pd37knmeua7wjxghmlskwf7x8pf", "kavavaloper1jlzx4js09d9zzyuhuz8sfdweklrapzacuhsxq5", "kavavaloper1jl42l225565y3hm9dm4my33hjgdzleucqryhlx", "kavavaloper1nxgg4grsc0fwh893mks62d3x3r6uazgpj3m3cr", "kavavaloper1ndkn5rdl9n929am6q2zt9ndfhhggcxkhetna90", "kavavaloper1njvaku4qg9pxmx9jgjks36xrxfd6fyqs3tgs4d", "kavavaloper1nnwwu4km0alut2q8vhg7zjt45wyehddpwlfrmj", "kavavaloper15kwwzz908wl0qv4w66a5ee70kytpmfc9khvah6", "kavavaloper15urq2dtp9qce4fyc85m6upwm9xul3049dcs7da", "kavavaloper14fkp35j5nkvtztmxmsxh88jks6p3w8u7p76zs9", "kavavaloper140g8fnnl46mlvfhygj3zvjqlku6x0fwu6lgey7", "kavavaloper14kn0kk33szpwus9nh8n87fjel8djx0y02c7me3", "kavavaloper1kgddca7qj96z0qcxr2c45z73cfl0c75p27tsg6", "kavavaloper1kwj4l5putuymgxw9kx8emh3e5dpaca0hnf3zdy", "kavavaloper1k4kxxfkhhwvyzxxxgksefkzpxahp9wfc572esl", "kavavaloper1kc6vzheht92jwf0gtzhjk6jjht67rxhal9z04v", "kavavaloper1k760ypy9tzhp6l2rmg06sq4n74z0d3rejwwaa0", "kavavaloper1h9ulmhqv5e2373khk6s9n0wtrfc5qavre09fxl", "kavavaloper1hezl6xwva28xt0hk204dllalagenmfsnuu50j6", "kavavaloper1c9ye54e3pzwm3e0zpdlel6pnavrj9qqvh0atdq", "kavavaloper1cj9cdx9mg95lhvpquym08ncgzpjvhmnwdvm5kc", "kavavaloper1ceun2qqw65qce5la33j8zv8ltyyaqqfctl35n4", "kavavaloper168g9nn9vnamhsnjkm7uqqee9f3v07flgwwddf9", "kavavaloper1645czvg787l0jr6mawhs9dm3mljnggj003yed9", "kavavaloper16lnfpgn6llvn4fstg5nfrljj6aaxyee9z59jqd", "kavavaloper1m9y0c7j7wxyu3nqtmeevyfzzpga8jhdyqw42wx", "kavavaloper1mu78xhlr705mzgwqykcafp4xy3kgatvwzrww8z", "kavavaloper1uvz0vus7ktxt47cermscwe3k9gs7h9ag05sh6g", "kavavaloper1udg4pal8c9gffv4l7zhvza027z0y345gd6esa6", "kavavaloper1u0kfndes0pf8dstaunsgumv7scsnmy09p3ln9r", "kavavaloper1u3jf6c2f85kmjldxsncnhsdp44nac7v5j7vzpc", "kavavaloper1u3hqe2m7vm59l30tyaqd3zurz864dlsg7nq83f", "kavavaloper1u7vsaanwt4e5mdzmxuqurmccxjka3h0ns2n3f5", "kavavaloper1aj9r9ll8m72xr5pmxr9fmum88k864tqg39nkfu", "kavavaloper1ajwfalplxnhkwhsycfax36yyyxpxz3450s800x", "kavavaloper17wcggpjx007uc09s8y4hwrj8f228mlwez945ey", "kavavaloper17498ffqdj49zca4jm7mdf3eevq7uhcsgjvm0uk"}
)

func testRandom(vps []int64, ids []string, amount, algo int) (float64, int64, float64, int64, Validators) {
	var vals Validators
	for idx := range vps {
		vals.ValidatorSet = append(vals.ValidatorSet,
			Validator{
				ID:          []byte(ids[idx]),
				VotingPower: vps[idx],
			},
		)
	}

	sort.Slice(vals.ValidatorSet, func(i, j int) bool {
		return vals.ValidatorSet[i].VotingPower > vals.ValidatorSet[j].VotingPower
	})

	if len(vals.ValidatorSet) > 100 {
		vals.ValidatorSet = vals.ValidatorSet[:100]
	}
	for idx := range vals.ValidatorSet {
		vals.ValidatorSet[idx].Number = idx + 1
	}
	// fmt.Println("~~~~~~")
	// for idx := range vals.ValidatorSet {
	// 	fmt.Print(float64(vals.ValidatorSet[idx].VotingPower)/float64(vals.GetVotingPower()), " ")
	// }
	// fmt.Println("~~~~~~")

	// fmt.Println()

	allVp := vals.GetVotingPower()
	randString := fmt.Sprint(rand.Int())
	luckyVal, worse := randomValidators([]byte(randString), vals, amount, algo)
	sumPos := 0

	for _, val := range luckyVal.ValidatorSet {
		sumPos += val.Number
		// fmt.Println("=>", val.Number)
	}
	// fmt.Println("!!", luckyVal.GetVotingPower())

	return float64(luckyVal.GetVotingPower()) / float64(allVp), allVp, float64(sumPos) / float64(amount), worse, luckyVal
}

func dup(vps []int64, ids []string, amount, algo int) {
	maxVp := float64(0)
	minVp := float64(999999999999999)
	sumVp := float64(0)
	sumPos := float64(0)
	maxDiff := int64(999999999999999)
	maxDiffMore := int64(0)
	maxDiffLess := int64(0)
	sumMore := float64(0)
	sumLess := float64(0)
	var allVp, worse int64
	var pos, luckyVp float64
	var luckyVal, worseVal, avgVal Validators
	_ = worseVal
	_ = avgVal
	freq := make([]int, 101)
	// level := make([][]int, amount)
	// for i := range level {
	// 	level[i] = make([]int, 1)
	// }

	avgVal.ValidatorSet = make([]Validator, amount)

	for i := 0; i < t; i++ {
		luckyVp, allVp, pos, worse, luckyVal = testRandom(vps, ids, amount, algo)
		sumVp += luckyVp
		if maxVp < luckyVp {
			maxVp = luckyVp
		}

		if minVp > luckyVp {
			minVp = luckyVp
		}
		sumPos += pos

		sort.Slice(luckyVal.ValidatorSet, func(i, j int) bool {
			return luckyVal.ValidatorSet[i].VotingPower > luckyVal.ValidatorSet[j].VotingPower
		})

		for _, val := range luckyVal.ValidatorSet {
			freq[val.Number]++
		}

		more := int64(0)
		for idx := 0; idx < amount/2; idx++ {
			more += luckyVal.ValidatorSet[idx].VotingPower
		}
		less := int64(0)
		for idx := amount / 2; idx < amount; idx++ {
			less += luckyVal.ValidatorSet[idx].VotingPower
		}

		for idx := 0; idx < amount; idx++ {
			avgVal.ValidatorSet[idx].VotingPower += luckyVal.ValidatorSet[idx].VotingPower
			avgVal.ValidatorSet[idx].Number += luckyVal.ValidatorSet[idx].Number

		}

		if less < maxDiff {
			maxDiff = less
			maxDiffMore = more
			maxDiffLess = less
			worseVal = luckyVal
		}

		sumMore += float64(more)
		sumLess += float64(less)

	}

	// for idx := 0; idx < amount; idx++ {
	// 	for x := 0; x < 94; x++ {
	// 		fmt.Print(level[idx][x], " ")
	// 	}
	// 	fmt.Println()
	// }

	// for _, val := range worseVal.ValidatorSet {
	// 	fmt.Print(" ", float64(val.VotingPower)/float64(allVp))
	// }
	// fmt.Println()
	// // for _, val := range worseVal.ValidatorSet {
	// // 	fmt.Print(" ", val.Number)
	// // }
	// // fmt.Println()
	// for _, val := range avgVal.ValidatorSet {
	// 	fmt.Print(" ", float64(val.VotingPower)/float64(allVp)/float64(t))
	// }
	// // fmt.Println()
	// // for _, val := range avgVal.ValidatorSet {
	// // 	fmt.Print(" ", float64(val.Number)/float64(t))
	// // }
	// fmt.Println()

	avgVpPercent := float64(sumVp) / float64(t) * float64(100)
	minVpPercent := float64(minVp) * float64(100)
	// maxVpPercent := float64(maxVp) * float64(100)
	// avgPos := float64(sumPos) / float64(t)
	worsePercent := float64(worse) / float64(allVp) * float64(100)

	morePercent := sumMore / float64(allVp) * float64(100) / float64(t)
	lessPercent := sumLess / float64(allVp) * float64(100) / float64(t)
	maxDiffMorePercent := float64(maxDiffMore) / float64(allVp) * float64(100)
	maxDiffLessPercent := float64(maxDiffLess) / float64(allVp) * float64(100)

	_ = avgVpPercent
	_ = minVpPercent
	_ = worsePercent
	_ = morePercent
	_ = lessPercent

	// fmt.Println("avg voting power percent= ", avgVpPercent)
	// fmt.Println("min voting power percent= ", minVpPercent)
	// fmt.Println("max voting power percent= ", maxVpPercent)
	// fmt.Println("avg position of validator set= ", avgPos)
	// fmt.Println("worse case voting power percent= ", worsePercent)
	// fmt.Println("avg more = ", morePercent)
	// fmt.Println("avg less = ", lessPercent)
	// fmt.Println("max diff more = ", float64(maxDiffMore)/float64(maxDiffMore+maxDiffLess))
	// fmt.Println("max diff less = ", float64(maxDiffLess)/float64(maxDiffMore+maxDiffLess))
	if maxDiffMorePercent < maxDiffLessPercent {
		kk := maxDiffMorePercent
		maxDiffMorePercent = maxDiffLessPercent
		maxDiffLessPercent = kk
	}
	// fmt.Println("==>", maxDiffMorePercent, maxDiffLessPercent)
	fmt.Println(avgVpPercent, minVpPercent, worsePercent, morePercent+lessPercent, lessPercent, maxDiffMorePercent+maxDiffLessPercent, maxDiffLessPercent)

	// for idx := range freq {
	// 	if idx != 0 {
	// 		fmt.Print(freq[idx], " ")
	// 	}
	// }
	// fmt.Println()
	// fmt.Println(freq)

}

func testCosmos(amount, algo int) {
	cosmosVps := []int64{243735, 457336, 928470, 6047109, 96455, 757388, 723902, 165563, 460520, 525680, 575653, 11209236, 115926, 100686, 14367, 457219, 250002, 1072128, 10244, 4772345, 302293, 145001, 193516, 938963, 385425, 28994, 419774, 572547, 630479, 1216988, 3490106, 557176, 337852, 7314356, 685457, 1221660, 245002, 1054073, 809489, 2283342, 10331469, 929700, 1515889, 73145, 331660, 68002, 1016350, 194202, 3479842, 399992, 251809, 1814827, 606078, 10194, 117924, 686014, 191760, 2149984, 160004, 9797346, 448224, 100004, 12594, 1000040, 500001, 10436, 2161060, 720894, 148298, 5365254, 25070, 3441942, 1745837, 895229, 85872, 106391, 1604729, 94068, 1574051, 3590614, 12688106, 170902, 5614171, 654800, 931443, 61729, 11377197, 10418, 576087, 970924, 1079608, 441499, 85001, 342875, 16050, 2123703, 353702, 5236944, 331888, 89465, 731460, 520172, 3994359, 51979, 1177445, 12099, 1129720, 412338, 11107, 4536088, 2783510, 10006, 2069404, 93436, 90120, 1658363, 5856297, 2277639, 1564696, 111581, 64325, 10028, 727981, 900529, 570104}
	cosmosIds := []string{"000001E443FD237E4B616E2FA69DF4EE3D49A94F", "000AA5ABF590A815EBCBDAE070AFF50BE571EB8B", "019B9CA2944D3CC36C7C73283EF3D58E56C8A5D4", "099E2B09583331AFDE35E5FA96673D2CA7DEA316", "0D2186876D26D882F7BE50DED92BD3CB53838143", "18C78D135C9D81D74F6234DBD268C47F0F89E844", "1AE0BD432F9A5122474A646325D1AFA6068692E9", "1CED30733D1625C89AB698677606D0E37B3676A9", "1E20E074030FA93836CABAED6581D6C46AEF6471", "1E9CE94FD0BA5CFEB901F90BC658D64D85B134D2", "21223475CE86F3C7CD5E985AA88FC24A29C97813", "2199EAE894CA391FA82F01C2C614BFEB103D056C", "246808FF08F382EB00E125F5088B6FE32B978857", "24935D59FAA94E793652CBF4716C6041CD7AA400", "2745ECA0BACF8F7B1BF8B7B434629B4F18C7BA42", "2B9A55D3BF93D7375DD207B75C5ED4D2B91D9146", "2C528D2345ED6E953C1C0819E2C3A01ABFCBA557", "2C9CCC317FB283D54AC748838A64F29106039E51", "2CD6A3523F5D61573EF1E3654B421AD211CFE8DC", "31920F9BC3A39B66876CC7D6D5E589E10393BF0E", "3363E8F97B02ECC00289E72173D827543047ACDA", "33C99B2E79B887ADA2228E53FB6EA000AB983337", "34337FEEE36B03CC9D5C9F59B308C0E317DB607C", "399671C2FE4B2714EC6E87D4EE454EF15F33AA2A", "407F144D1C9DEA4EE6A8CBC2D4C022A657506B83", "412D552A7129C496AD4C3DAAA7CEF7089C73B140", "42D6705E716616B4A5442BDAA050B7C6E9FDDE43", "46A3F8B8393BAA153C40E5722EAE82EA0D48B32D", "4906F2A5334D906A4C63F9E9D61527A9F593C4EF", "49BBFB1BA1A75052E3226E8E1E0EFEB33918B8B2", "4AF69D6A5436C30E3584C1628433DE55E758BCCA", "4C92230FAC162303D981C06DD22663A4FC7622BC", "4E9CB39F4B1FA617339744A5600B62802652D69C", "51205659A717DFFB96E054F8BD1108730E17AEA7", "51DB2566204EE266427EA8A6CB719835AB170BE9", "52E1646134432BF9532B4881C6ED32E40AE5A2DD", "5731848E19257705AA28CC7EFAA8C708EE014D52", "57713BB7421C7FEB381B863FC87DED5E829AA961", "5AB353B748D45F20DFCE19D73BA89F26E1C34CF7", "671460930CCDC9B06C5D055E4D550EB8DAF2291E", "679B89785973BE94D4FDF8B66F84A929932E91C5", "68F5BBEACEF114C720EA9C98BFA2FFDE01C54FD1", "696ABC95186FD65A07050C28AB00C9358A315030", "6CB47D786B2F350C13A60BB77D398AC82E900985", "6F322BAF5D73D5C77271941901B1D866476D5C90", "708FDDCE121CDADA502F2B0252FEF13FDAA31E50", "70C5B4E6779C59A24CFD9146581E27021C2AEC26", "732CEEF54C374DDC6ADECBFD707AEFD07FEDC143", "75DAB316F4CA1367F532AB71A80B7FA65AB69039", "7674C7FE7960B486BA1052DA225F54119580EB74", "7A416D78E6C9AFCF115949CA092E61BCD9DB8A18", "7B3A2EFE5B3FCDF819FCF52607314CEFE4754BB6", "7B3D01F754DFF8474ED0E358812FD437E09389DC", "7CAB2B15C7F0238619A65D8833EEB53942552651", "7FCBACCB12563C8823E8EBEC8C0DE580DD8C7F04", "808D6B054A0B6D3FF5F5EAF0A65CFC64C543F833", "818964B4FB36D28109C3E853778B33231B27C5FC", "81965FE8A15FA8078C9202F32E4CFA72F85F2A22", "8332069BCA9D125E72BF87E4DA172C9E83B90FD4", "83F47D7747B0F633A6BA0DF49B7DCF61F90AA1B0", "846BE4F39E3122D2A2D3FE5454E2561073E95538", "874D6AA838384A79A3D4D062541F91BF7E31BDBA", "8C8802A921114169D2581CD46E3CA6853F6F2A7F", "8CE843E04C48B7864F8568FA39E90EA13FE6586F", "8E0EE37B7B1A038DD145E30F1EF97DF3619EF429", "8FA2380FA64B122427202970CF12BD4991B4C6D2", "8FE3CFFA6A07B093E441BB84DA1B6DABF53AFA2D", "91C823A744DE50F91C17A46B624EDF8F7150A7DD", "955A47C8AC8632825DD475E90913D40AB09D3FB4", "95E060D07713070FE9822F6C50BD76BCCBF9F17A", "99304AA9AEEC13FCEF0FF72DFC8953273FE559BC", "9C17C94F7313BB4D6E064287BEEDE5D3888E8855", "9D07B301D23C547266D55D1B6C5A78CA473383A1", "9DC4012099BE743189074B85E49891AE3B3FEE9B", "9DF8E338C85E879BC84B0AAA28A08B431BD5B548", "9E14352CB5293C6073D61280A197085C6748DAFA", "9EE94DBB86F72337192BF291B0E767FD2729F00A", "A03DC128D38DB0BC5F18AE1872F1CB2E1FD41157", "A4F1D5534F3FA905A4DA606E8A10834976511FF7", "A6935D877B9776C45B96EEAE526959A3B9A5AB1A", "AC2D56057CD84765E6FBE318979093E8E44AA18F", "AC885F3EE81E7ED07FE7B2E067443A855F997BA1", "B00A6323737F321EB0B8D59C6FD497A14B60938A", "B0155252D73B7EEB74D2A8CC814397E66970A839", "B0765A2F6FCC11D8AC46275FAC06DD35F54217C1", "B0C5370ED46641878800A518EB8335E42AE87673", "B1167D0437DB9DF0D533EE2ACDE48107139BDD2E", "B2BF68AD4CED6FE8F71AACAD01003436EBE0729F", "B34591DA79AAD0213534E2E915F50DE5CDBDF250", "B4999CD535E4CD32B590BEB47020A724F40B65E5", "B4E1085F1C9EBB0EA994452CB1B8124BA89BED1A", "B543A7DF48780AEFEF593A003CD060B593C4E6B5", "B61AA419909956354A25E26DDE9103B4C1A4D5CA", "B6D7360C27F1DC36DD9BFCF23037BE8B04429209", "B724CDAA69A47B90B18D0EA7FD5B046D537DA64A", "BAC33F340F3497751F124868F049EC2E8930AC2F", "BF4CB4D59D19D451CF5E7BC49349DF4AA222D78B", "C2356622B495725961B5B201A382DD57CD3305EC", "C2DDD9700CF5DEC0457DC423829B31EA8FD4F9D4", "C4903229B9EAD415C79E8FA69D2BBA6117617C41", "C52ACDB32057F5C731BBDD48460B93C3500DD324", "CA6696F5FE66A480BF21460319BE979930852DD0", "CC05882978FC5FDD6A7721687E14C0299AE004B8", "CC87F56B58621811E2B5A47F38C6166E295CE36E", "D14A542E8756C3A942D9FD8873DC2E9A7798A17F", "D3FA9465ED96B70797572D4D83678BC2699560C1", "D540AB022088612AC74B287D076DBFBC4A377A2E", "D7F7C79487C10A5CF1ABEB1DBD81E8D49757C422", "D8F0BEDE91D3FE4DD827619499EBA5CECC76FE21", "D9F8A41B782AA6A66ADC81F953923C7DCE7B6001", "DA6AAAA959C9EF88A3EB37B1F107CB2667EBBAAB", "E6A3801926BC5F83761C4ADF88263F88195FDD5B", "E800740C68C81B30345C3AE2BA638FA56FF67EEF", "E83BFC436D2CE8DCC9EC0589B2E5B735E37FB85C", "EB1DF22507B79CE700F86C4C8B13D7DF01DFDA9C", "EBED694E6CE1224FB1E8A2DD8EE63A38568B1E2B", "ED509E78097E1306A91FEDE8E85B75D06BDDF6E3", "EE73A19751D58C5EC044C11E3FB7AE685A10D2C1", "F4CAB410DE5567DB203BD56C694FB78D482479A1", "F55714243D32FB65B6D95A29D0350EA0CABBA8EA", "F59734A896A7689436BC3422244FD862AE189C5C", "F8C01C0681578AA700D736D675C9992065F65E3E", "F919902709B7482F01C030E8B57BF93B8D87043B", "FA0E5DFACCDCF74957A742144FE55BE61D433377", "FD5D54E0D9E4768FEA4C0DFFDC89FA96B6657F32"}
	dup(cosmosVps, cosmosIds, amount, algo)
}

func testKava(amount, algo int) {
	kavaVps := []int64{23004398541, 20069063571, 20004700000, 5508940018, 105030540201, 22999900000, 2675001000000, 23459640449, 21000000000, 100010000000, 20410000001, 23163160000, 48503759111, 23138902574, 32825178240, 37775800000, 2827139316200, 929346000000, 14040219, 20010000000, 1171728959900, 23171700000, 36645766677, 466602000000, 2301737241460, 20102000000, 22631605772, 22994500500, 10000000, 23011000000, 5085690431269, 2999500, 20625150000, 1506400478285, 23189000000, 125050100000, 20153000000, 1999980100000, 3317777023616, 20100000782, 4001049500000, 96527616407, 23000000000, 23145977490, 5005157202698, 4835289500, 101000000, 1009696987835, 22950000000, 1963559920986, 22950000000, 23213610000, 23081092090, 23144338847, 2678950000000, 3148249101016, 23443000000, 3148822472067, 62461232393, 23085863116, 55874224415, 1030886500000, 5723799001, 6118996081, 2354970117761, 23734500000, 3533085741911, 2888251903903, 2156193000000, 23230000000, 23388000000, 20028000000, 25830577485, 23225502889, 169401000000, 4121138974317, 932909664, 2425826213735, 23195513197, 23136261400, 2833230807310, 1255000000000, 21000000000, 32995000000, 1739448257662, 4360445500000, 3303951185253, 3305241541060, 1000009950010, 21000000000, 59802020969, 23092086042, 23770000000}
	kavaIds := []string{"kavavaloper1qyc2cfl0nw8r95dsdw534x99wq0xcj9rmxpl7z", "kavavaloper1qfy0e2w62g6j4jg5djcqd4py3zsaeqexjplj2d", "kavavaloper1q3y9qga5hf360dmzta67vp54qz25tmv4hhkk4t", "kavavaloper1qk0pta4ga5t8p5vv7me8dz32lvcrv2rp098cas", "kavavaloper1ptyzewnns2kn37ewtmv6ppsvhdnmeapvl7z9xh", "kavavaloper1pn5k9c5pxmg5f0rycpl9rrx6k6mk85scxf06zx", "kavavaloper1phd8jz25lumudc7ac7rhmupvcqcv7lg3c8dprc", "kavavaloper1pceqe8we7drpfqrutchwy3f99800hhzuw6cc84", "kavavaloper1rgcgqmnkeffks7enrv6hk5u4wg3nzfkmqlzjqd", "kavavaloper1yrm63pqvld8uyzkavz55p2cktpm2gm8jd8xxlu", "kavavaloper1yjj2wfers947l6n5pynpgsqlz7svc5n8ssl6ye", "kavavaloper1yna6lete8nwwwctsalzdg04ldqaz73gtn33ydq", "kavavaloper196anr2ycsalg806dz29afklnecuaupvkh5qz6c", "kavavaloper1xzud24na9tucauc7lf6pjk84kqsgq3eq747ta7", "kavavaloper1x9hq3rjc48t5upcsr3c209ycgekfasne3l5nkc", "kavavaloper1xftqdxvq0xkv2mu8c5y0jrsc578tak4m9u0s44", "kavavaloper1xka27j0jvmq97yunj5wp8fv242lycmax8ejlaf", "kavavaloper1xhxzmj8fvkqn76knay9x2chfra826369dhdu2c", "kavavaloper18zksjhrefqew0zahmts894p8asscufxvdfq702", "kavavaloper18s9m5d5cjf0humjv7mkq8pm47kchwm0r0369cx", "kavavaloper18cf35l7req0k6ulqapeyv830mrrucn9xj87plr", "kavavaloper1g20mhvcpjxp6gzlwhtfcphjehwcl2njqydgu7q", "kavavaloper1gtd040dmljyaty9tkhq0mlqz7saer48v4d608x", "kavavaloper1g4qpetrj59a29e4wxpe74x93q4df2czjh8r9ak", "kavavaloper1ffcujj05v6220ccxa6qdnpz3j48ng024ykh2df", "kavavaloper1fw7vjc3fphahqxpdjypddlulnltxws8g0mrds7", "kavavaloper1fmas0qlsucg4qwf8mqyrylcg3uluz2ffg8952q", "kavavaloper12qn7y04wzr5s3h4dmdtre4q9f4nvc03a9a9qsz", "kavavaloper12r77hhj6ylvvl4etjm0fmpzh07jum8u7qqd695", "kavavaloper129kkennsm7na34lu6sn4kwxp7ewes58y4fx6y9", "kavavaloper12g40q2parn5z9ewh5xpltmayv6y0q3zs6ddmdg", "kavavaloper12et238paeqxyhvk2pfs20ygfj6ct0dx2ccsdz5", "kavavaloper1t5l8ht0wxpd4lpe4cweftrpg5kyn3qp437yvnr", "kavavaloper1vyx7wt8s8dwcspdt7dy49sq4l3jwyhxyndakmm", "kavavaloper1vfawvtzvmcjkpqzhezvhk9q5tv5t7x8smz95uu", "kavavaloper1vw35vclatlrcmzuaxf2fleyuk38xa7xf4vdaq2", "kavavaloper1vw4t8ge2ephu0wuhcmclcw04ag2vzj9rpdme2x", "kavavaloper1vuylvflgy75d8zr07ta8x0gcqynvkvw70et853", "kavavaloper1dwae0ny4uuvacucakm9v8r8mxhw82zack4cn7y", "kavavaloper1dj63l3z0smqa8au5yek37nmj0xd60zs9dwmdh0", "kavavaloper1dntlhrw3jrej6ssdp64yfmkkz08ykyx4n7hphh", "kavavaloper1dede4flaq24j2g9u8f83vkqrqxe6cwzrxt5zsu", "kavavaloper1davtd9w5yadvlg5x7aw0kpkyhckk5m5ecrmnyl", "kavavaloper1wtcn3ylrsp4qp4urlhkf78kkt8f2ch7p0f0p98", "kavavaloper1wu8m65vqazssv2rh8rthv532hzggfr3h9azwz9", "kavavaloper1wacmlkst8s39fq83tqhlcuuacts2z6lvwfq8hv", "kavavaloper1006e7kwuhthdtxe0d90pc209cy803ehem0dg0q", "kavavaloper10m3hjapny44txmgr47rf277364htgqpr646cty", "kavavaloper1spkjjtwks5zt0dqexj8rv8ljwmwxe8ufkraukq", "kavavaloper1srhded4xw0krcwlvddtcyycuh70fp5ry9yvp86", "kavavaloper1s8akp0nq7z7vuf5g9agdswcwq7vm9uwudk0han", "kavavaloper13fxkk4730cqglgdv7w0mdelyx07myyq76uf9f3", "kavavaloper13vstf6ecmfe4p0gaufumk569sdawhtrf8gu56h", "kavavaloper13dfu6et0m8zm4hachudvn62w0c2zvcmrqn8cs3", "kavavaloper1jyuv7z9at27elvmnmzh2v39dc06r9kjcy59xkr", "kavavaloper1j26c4k2jj9tv95whdhva3e8v2fcm4s3dsgstd2", "kavavaloper1jhraz4ftxl2pd37knmeua7wjxghmlskwf7x8pf", "kavavaloper1jlzx4js09d9zzyuhuz8sfdweklrapzacuhsxq5", "kavavaloper1jl42l225565y3hm9dm4my33hjgdzleucqryhlx", "kavavaloper1nxgg4grsc0fwh893mks62d3x3r6uazgpj3m3cr", "kavavaloper1ndkn5rdl9n929am6q2zt9ndfhhggcxkhetna90", "kavavaloper1njvaku4qg9pxmx9jgjks36xrxfd6fyqs3tgs4d", "kavavaloper1nnwwu4km0alut2q8vhg7zjt45wyehddpwlfrmj", "kavavaloper15kwwzz908wl0qv4w66a5ee70kytpmfc9khvah6", "kavavaloper15urq2dtp9qce4fyc85m6upwm9xul3049dcs7da", "kavavaloper14fkp35j5nkvtztmxmsxh88jks6p3w8u7p76zs9", "kavavaloper140g8fnnl46mlvfhygj3zvjqlku6x0fwu6lgey7", "kavavaloper14kn0kk33szpwus9nh8n87fjel8djx0y02c7me3", "kavavaloper1kgddca7qj96z0qcxr2c45z73cfl0c75p27tsg6", "kavavaloper1kwj4l5putuymgxw9kx8emh3e5dpaca0hnf3zdy", "kavavaloper1k4kxxfkhhwvyzxxxgksefkzpxahp9wfc572esl", "kavavaloper1kc6vzheht92jwf0gtzhjk6jjht67rxhal9z04v", "kavavaloper1k760ypy9tzhp6l2rmg06sq4n74z0d3rejwwaa0", "kavavaloper1h9ulmhqv5e2373khk6s9n0wtrfc5qavre09fxl", "kavavaloper1hezl6xwva28xt0hk204dllalagenmfsnuu50j6", "kavavaloper1c9ye54e3pzwm3e0zpdlel6pnavrj9qqvh0atdq", "kavavaloper1cj9cdx9mg95lhvpquym08ncgzpjvhmnwdvm5kc", "kavavaloper1ceun2qqw65qce5la33j8zv8ltyyaqqfctl35n4", "kavavaloper168g9nn9vnamhsnjkm7uqqee9f3v07flgwwddf9", "kavavaloper1645czvg787l0jr6mawhs9dm3mljnggj003yed9", "kavavaloper16lnfpgn6llvn4fstg5nfrljj6aaxyee9z59jqd", "kavavaloper1m9y0c7j7wxyu3nqtmeevyfzzpga8jhdyqw42wx", "kavavaloper1mu78xhlr705mzgwqykcafp4xy3kgatvwzrww8z", "kavavaloper1uvz0vus7ktxt47cermscwe3k9gs7h9ag05sh6g", "kavavaloper1udg4pal8c9gffv4l7zhvza027z0y345gd6esa6", "kavavaloper1u0kfndes0pf8dstaunsgumv7scsnmy09p3ln9r", "kavavaloper1u3jf6c2f85kmjldxsncnhsdp44nac7v5j7vzpc", "kavavaloper1u3hqe2m7vm59l30tyaqd3zurz864dlsg7nq83f", "kavavaloper1u7vsaanwt4e5mdzmxuqurmccxjka3h0ns2n3f5", "kavavaloper1aj9r9ll8m72xr5pmxr9fmum88k864tqg39nkfu", "kavavaloper1ajwfalplxnhkwhsycfax36yyyxpxz3450s800x", "kavavaloper17wcggpjx007uc09s8y4hwrj8f228mlwez945ey", "kavavaloper17498ffqdj49zca4jm7mdf3eevq7uhcsgjvm0uk"}

	dup(kavaVps, kavaIds, amount, algo)

}

func main() {
	// x := []int{5, 7, 10, 15}
	x := []int{7}
	for algo := 4; algo <= 5; algo++ {
		for idx := 0; idx < len(x); idx++ {
			// fmt.Print(algo, amount, " ")
			fmt.Println()
			// testCosmos(x[idx], algo)
			// fmt.Println("--------")
			testKava(x[idx], algo)
		}
	}
	// testCosmos(10, 3)
}
