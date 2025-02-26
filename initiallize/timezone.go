package initiallize

import "time"

func TimezoneInit() {
	var cstZone = time.FixedZone("CST", 8*3600) // 东八
	time.Local = cstZone
}
