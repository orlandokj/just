package ringlist


type RingList struct {
    startRing *RingItem
    currentRing *RingItem
    maxLength int
    currentLength int
}

type RingItem struct {
    value interface{}
    next *RingItem
}

func NewRingList(maxLength int) *RingList {
    startRing := new(RingItem)
    startRing.next = startRing
    return &RingList{
        startRing: startRing,
        currentRing: startRing,
        maxLength: maxLength,
        currentLength: 0,
    }
}

func (rl *RingList) Add(value interface{}) {
    rl.currentRing.value = value
    if rl.currentLength < rl.maxLength {
        rl.currentRing.next = new(RingItem)
        rl.currentRing = rl.currentRing.next
        rl.currentRing.next = rl.startRing
        rl.currentLength++
    } else {
        rl.currentRing = rl.currentRing.next
        rl.startRing = rl.currentRing.next
    }
}

func (rl *RingList) Len() int {
    return rl.currentLength
}

func (rl *RingList) Clear() {
    rl.currentLength = 0
    rl.startRing = new(RingItem)
    rl.startRing.next = rl.startRing
    rl.currentRing = rl.startRing
}

func (rl *RingList) Do(f func(interface{})) {
    ring := rl.startRing
    for i := 0; i < rl.currentLength; i++ {
        f(ring.value)
        ring = ring.next
    }
}
