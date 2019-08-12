package fetcher

import "github.com/etclabscore/go-ethereum/core/types"

type FetcherInsertBlockEvent struct {
	Peer  string
	Block *types.Block
}
