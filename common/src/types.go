package common

type Index struct {
	Entries  []IndexEntry
	Checksum string
}

type IndexEntry struct {
	Name, Host string
}

type ChecksumResponse struct {
	Checksum string
}
