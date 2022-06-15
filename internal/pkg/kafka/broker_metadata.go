package kafka

import (
	"sort"
)

type BrokerMetadata struct {
	Details        *Broker
	ConsumerGroups []string
	Logs           []*LogFile
}

type aggregatedTopicSize map[string]*LogEntry

type LogEntry struct {
	Topic     string
	Permanent uint64
	Temporary uint64
}

type LogFile struct {
	Path    string
	Entries aggregatedTopicSize
}

// func newLogFile(path string) *LogFile {
// 	return &LogFile{
// 		Path:    path,
// 		Entries: make(aggregatedTopicSize),
// 	}
// }

// func (l *LogFile) set(topic string, size int64, isTemp bool) {
// 	if _, ok := l.Entries[topic]; !ok {
// 		l.Entries[topic] = &LogEntry{
// 			Topic: topic,
// 		}
// 	}
// 	if isTemp {
// 		l.Entries[topic].Temporary += uint64(size)
// 	} else {
// 		l.Entries[topic].Permanent += uint64(size)
// 	}
// }

func (l *LogFile) SortByPermanentSize() []*LogEntry {
	result := l.toSlice()
	sort.Sort(logsByPermanentSize(result))
	return result
}

func (l *LogFile) toSlice() []*LogEntry {
	result := make([]*LogEntry, len(l.Entries))
	var i int
	for _, l := range l.Entries {
		result[i] = l
		i++
	}
	return result
}

type logsByPermanentSize []*LogEntry

func (l logsByPermanentSize) Len() int {
	return len(l)
}

func (l logsByPermanentSize) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

func (l logsByPermanentSize) Less(i, j int) bool {
	return l[i].Permanent > l[j].Permanent
}
