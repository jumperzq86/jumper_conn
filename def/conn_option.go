package def

type ConnOptions struct {
	MaxMsgSize     int64
	ReadTimeout    int64
	WriteTimeout   int64
	AsyncWriteSize int64
	Side           int8

	PongWait         int64
	PingPeriod       int64
	CloseGracePeriod int64
}

func (this *ConnOptions) CheckValid() error {
	if this.AsyncWriteSize <= 0 {
		return ErrInvalidConnParam
	}
	if this.PingPeriod != 0 && this.PingPeriod >= this.PongWait {
		return ErrInvalidConnParam
	}
	return nil
}
