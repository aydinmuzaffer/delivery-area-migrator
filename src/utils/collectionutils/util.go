package collectionutils

func GetChunk[T any](collection *[]T, chunkSize, iteration int) *[]T {
	var howManyToPick = chunkSize
	if len(*collection) < (iteration*chunkSize + chunkSize) {
		howManyToPick = len(*collection) - (iteration * chunkSize)
	}
	from := iteration * chunkSize
	to := from + howManyToPick
	chunk := (*collection)[from:to]
	return &chunk
}
