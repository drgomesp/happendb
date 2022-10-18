package core

type ReadRepository interface {
}

type WriteRepository interface {
}

type ReadWriteRepository interface {
	ReadRepository
	WriteRepository
}
