package utils

import "flag"

var (
	flagShortYes = flag.Bool("y", false, "Answer yes to all questions")
	flagYes      = flag.Bool("yes", false, "Answer yes to all questions")

	flagShortNo = flag.Bool("n", false, "Answer no to all questions")
	flagNo      = flag.Bool("no", false, "Answer no to all questions")

	flagShortVersion = flag.Bool("v", false, "Show the version of fpm")
	flagVersion      = flag.Bool("version", false, "Show the version of fpm")

	flagShortUpdate = flag.Bool("u", false, "Update fpm repository")
	flagUpdate      = flag.Bool("update", false, "Update fpm repository")
)

func FlagUpdate() bool {
	return *flagShortUpdate || *flagUpdate
}

func FlagVersion() bool {
	return *flagShortVersion || *flagVersion
}

func FlagYes() bool {
	return *flagShortYes || *flagYes
}

func FlagNo() bool {
	return *flagShortNo || *flagNo
}
