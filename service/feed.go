package service

import "time"

func Feed() (VideoList, error) {
	return NewFeedFlow().Do()
}

func NewFeedFlow() *FeedFlow {
	return &FeedFlow{}
}

type FeedFlow struct {
	Videos   []Video
	NextTime time.Time
}

func (f *FeedFlow) Do() (VideoList, error) {

	//latest_time := time.Now()
	return nil, nil
}
