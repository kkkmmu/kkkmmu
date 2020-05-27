package rut

var ModelNameToAPIMap = map[string]API{
	"V2224G": V2{},
	"V2216G": V2{},
	"V2208G": V2{},
	"V2724Z": V2{},
	"V2708M": V2{},
	"V2808M": V2{},
	"V5624G": V2{},
	"V6524G": V2{},

	"V5724G":      V5{},
	"V5724G_SFU":  V5{},
	"V5724G_SFUX": V5{},

	"V8500_SFU": V8{},
	"V8106_SFU": V8{},
	"V8102_SFU": V8{},

	"M3000":  XP{},
	"M3200":  XP{},
	"V3306G": XP{},
	"M2400":  XP{},
}
