package currencies

const (
	FeedChainLink = DataFeedType(1) + iota
	FeedCurve
	FeedUniswap
)

type DataFeedType uint8
