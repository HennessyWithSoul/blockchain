package network

type StatusMessage struct {
	Version       uint32
	CurrentHeight uint32
	ID            string
}
type GetStatusMessage struct {
}
type GetBlockMessage struct {
	From uint32
	To   uint32
}
