package conflict

// ResolveConflictRequest is used when resolving file conflicts at client side
type ResolveConflictRequest struct {
	RequestId           uint64
	LatestHash          string
	LatestSyncTimestamp uint64
}
