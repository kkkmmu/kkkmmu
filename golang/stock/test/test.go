package main

import (
	"db"
	"fmt"
)

/*
在上海证券交易所上市的证券，根据上交所"证券编码实施方案"，采用6位数编制方法，前3位数为区别证券品种，具体见下表所列：
    001×××国债现货；
    201×××国债回购；
    110×××120×××企业债券；
    129×××100×××可转换债券；
    310×××国债期货；
    500×××550×××基金；
    600×××A股；
    700×××配股；
    710×××转配股；
    701×××转配股再配股；
    711×××转配股再转配股；
    720×××红利；
    730×××新股申购；
    735×××新基金申购；
    900×××B股；
    737×××新股配售。
在深圳证券交易所上市面上证券，根据深交所证券编码实施采取4位编制方法，首位证券品种区别代码，具体见下表所示：
    0×××A股；
    1×××企业债券、国债回购、国债现货；
    2×××B股及B股权证；
    3×××转配股权证；
    4×××基金；
    5×××可转换债券；
    6×××国债期货；
    7×××期权；
    8×××配股权证；
    9×××新股配售。
*/

var histFmt string = "http://quotes.money.163.com/service/chddata.html?code=%s%s&start=19900101&end=20180701&fields=TCLOSE;HIGH;LOW;TOPEN;LCLOSE;CHG;PCHG;TURNOVER;VOTURNOVER;VATURNOVER;TCAP;MCAP"

func main() {
	sts, err := db.GetStockList()
	if err != nil {
		panic(err)
	}

	fmt.Println(len(sts))
}
