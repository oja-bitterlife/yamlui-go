package yamlui

// **********************************************************************
// 時間を扱う
// ステートレスにしたいけど、StartFrameだけはどうしても必要
type Timer struct {
	start    int
	duration int
}

func NewTimer(start, duration int) Timer {
	return Timer{start: start, duration: duration}
}

func (t *Timer) Reset(restart int) {
	t.start = restart
}

// ==================================================
// Getter
func (t *Timer) Start() int {
	return t.start
}

func (t *Timer) Duration() int {
	return t.duration
}

func (t *Timer) End() int {
	return t.start + t.duration
}

// ==================================================
// 経過を返す
func (t *Timer) Elapsed(current int) int {
	return current - t.start
}

func (t *Timer) Remaining(current int) int {
	elapsed := t.Elapsed(current)
	if elapsed >= t.duration {
		return 0
	}
	return t.duration - elapsed
}

func (t *Timer) Progress(current int, resolution int) int {
	elapsed := t.Elapsed(current)
	// 0除算注意
	if t.duration == 0 {
		return 0
	}

	if elapsed >= t.duration {
		return resolution
	}
	return elapsed * resolution / t.duration
}

func (t *Timer) IsFinish(current int) bool {
	// まだ開始していない
	if t.duration == 0 {
		return false
	}
	return current >= t.start+t.duration
}
