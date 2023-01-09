package utils

import "flag"

var (
	flagShortYes = flag.Bool("y", false, "Answer yes to all questions")
	flagYes      = flag.Bool("yes", false, "Answer yes to all questions")

	flagShortNo = flag.Bool("n", false, "Answer no to all questions")
	flagNo      = flag.Bool("no", false, "Answer no to all questions")
)

func FlagYes() bool {
	return *flagShortYes || *flagYes
}

func FlagNo() bool {
	return *flagShortNo || *flagNo
}
