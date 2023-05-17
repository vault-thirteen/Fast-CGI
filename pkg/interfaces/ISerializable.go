package i

type ISerializable interface {
	ToBytes() (ba []byte)
}

type ISerializableE interface {
	ToBytes() (ba []byte, err error)
}
