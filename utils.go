package mcts

import (
	"errors"
	"io"
	"math"
	"math/rand"
	"os"
	"path"
	"time"

	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

func UCB(totalReward float64, ucbConstant float64, parentVisits int64, currentVisits int64) float64 {
	exploitationValue := totalReward / float64(currentVisits)
	explorationValue := 2.0 * ucbConstant * math.Sqrt(2*math.Log(float64(parentVisits))/float64(currentVisits))
	return exploitationValue + explorationValue
}

// byValues implements sort.Interface to sort *descending* by selection score.
// Example: sort.Sort(byValues(nodes))
type byValues []*TreeNode

func (a byValues) Len() int           { return len(a) }
func (a byValues) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byValues) Less(i, j int) bool { return a[i].ucbValues > a[j].ucbValues }

// byReverseValues implements sort.Interface to sort *ascending* by selection score.
// Example: sort.Sort(byReverseValues(nodes))
type byReverseValues []*TreeNode

func (a byReverseValues) Len() int           { return len(a) }
func (a byReverseValues) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byReverseValues) Less(i, j int) bool { return a[i].ucbValues < a[j].ucbValues }

// byVisits implements sort.Interface to sort *descending* by visits.
// Example: sort.Sort(byVisits(nodes))
type byVisits []*TreeNode

func (a byVisits) Len() int           { return len(a) }
func (a byVisits) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byVisits) Less(i, j int) bool { return a[i].visits > a[j].visits }

// byVisits implements sort.Interface to sort *ascending* by AverageReward.
// Example: sort.Sort(byReverseAverageReward(nodes))
type byReverseAverageReward []*TreeNode

func (a byReverseAverageReward) Len() int      { return len(a) }
func (a byReverseAverageReward) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a byReverseAverageReward) Less(i, j int) bool {
	return a[i].totalReward/float64(a[i].visits) < a[j].totalReward/float64(a[j].visits)
}

// byVisits implements sort.Interface to sort *descending* by AverageReward.
// Example: sort.Sort(byVisits(nodes))
type byAverageReward []*TreeNode

func (a byAverageReward) Len() int      { return len(a) }
func (a byAverageReward) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a byAverageReward) Less(i, j int) bool {
	return a[i].totalReward/float64(a[i].visits) > a[j].totalReward/float64(a[j].visits)
}

// GetLogger returns the colorful logger based ion zerolog
func GetLogger(consoleEnable bool, fileEnable bool) zerolog.Logger {
	var writers []io.Writer

	if consoleEnable {
		consoleWriter := zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC1123}
		writers = append(writers, consoleWriter)
	}
	if fileEnable {
		fileWriter := NewFileWriter()
		writers = append(writers, fileWriter)
	}

	mw := io.MultiWriter(writers...)

	logger := zerolog.New(mw).With().Timestamp().Logger()
	return logger
}

func NewFileWriter() io.Writer {
	zone, _ := time.LoadLocation("Asia/Shanghai")
	currentTime := time.Now().In(zone).Format("01-02-2006 15:04:05")

	letters := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	var uuid string
	for i := 0; i < 6; i++ {
		uuid += string(letters[rand.Intn(24)])
	}

	directory := "logging"
	filename := currentTime + " " + uuid + ".log"

	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		if err := os.MkdirAll(directory, 0744); err != nil {
			return nil
		}
	}
	return &lumberjack.Logger{
		Filename:   path.Join(directory, filename),
		MaxSize:    512,
		MaxBackups: 0,
	}
}
