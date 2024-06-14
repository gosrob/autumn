package app

var dryRun bool = true

func SetDryRun(d bool) {
	dryRun = d
}

func GetDryun() bool {
	return dryRun
}
