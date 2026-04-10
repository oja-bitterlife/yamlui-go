package yamlui

// **********************************************************************
// 時間を扱う
// ステートレスにしたいけど、StartFrameだけはどうしても必要
type Timer struct {
	Start    int
	Duration int
}

func NewTimer(start, duration int) Timer {
	return Timer{Start: start, Duration: duration}
}

func (t *Timer) Reset(restart int) {
	t.Start = restart
}

// ==================================================
// Getter
func (t *Timer) End() int {
	return t.Start + t.Duration
}

// ==================================================
// 経過を返す
func (t *Timer) Elapsed(current int) int {
	return current - t.Start
}

func (t *Timer) Remain(current int) int {
	elapsed := t.Elapsed(current)
	if elapsed >= t.Duration {
		return 0
	}
	return t.Duration - elapsed
}

func (t *Timer) Progress(current int, resolution int) int {
	elapsed := t.Elapsed(current)
	// 0除算注意
	if t.Duration == 0 {
		return 0
	}

	if elapsed >= t.Duration {
		return resolution
	}
	return elapsed * resolution / t.Duration
}

func (t *Timer) IsFinish(current int) bool {
	// まだ開始していない
	if t.Duration == 0 {
		return false
	}
	return current >= t.Start+t.Duration
}
