package currencies

const (
	FeedHttp = DataFeedType(1) + iota
	FeedChainLink
	FeedCurve
	FeedUniswap
)

type DataFeedType uint8
