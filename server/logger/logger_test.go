package logger

import (
	"crypto/rand"
	"encoding/base32"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

func randomFileName() string {
	buff := make([]byte, 24)
	rand.Read(buff)
	return strings.TrimRight(base32.StdEncoding.EncodeToString(buff), "=")
}

func TestLogger(t *testing.T) {
	fileName := path.Join(os.TempDir(), randomFileName())
	lgr, err := New(fileName)
	require.Nil(t, err)
	defer os.Remove(fileName)

	randInput := make([]byte, 512)
	rand.Read(randInput)

	for i := 0; i < 100; i++ {
		n, err := lgr.Write(randInput)
		require.Nil(t, err)
		assert.Equal(t, 512, n)
	}

	err = lgr.Close()
	assert.Nil(t, err)

	// can't write to a closed logger
	_, err = lgr.Write(randInput)
	assert.NotNil(t, err)

	// can't call close after logger has been closed
	err = lgr.Close()
	assert.NotNil(t, err)

	info, err := os.Stat(fileName)
	require.Nil(t, err)
	assert.Equal(t, int64(51200), info.Size())

}

func BenchmarkLogger(b *testing.B) {
	fileName := path.Join(os.TempDir(), randomFileName())
	lgr, err := New(fileName)
	if err != nil {
		b.Fatal("new failed ", err)
	}
	defer os.Remove(fileName)

	randInput := make([]byte, 512)
	rand.Read(randInput)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := lgr.Write(randInput)
		if err != nil {
			b.Fatal("write failed ", err)
		}
	}

	b.StopTimer()

	lgr.Close()
}

func BenchmarkLumberjack(b *testing.B) {
	fileName := path.Join(os.TempDir(), randomFileName())
	lgr := &lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     28, //days
	}
	defer os.Remove(fileName)

	randInput := make([]byte, 512)
	rand.Read(randInput)
	// first lumberjack write opens file so we count that as part of initialization
	// just to make sure we're comparing apples to apples with our logger
	_, err := lgr.Write(randInput)
	if err != nil {
		b.Fatal("first write failed ", err)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := lgr.Write(randInput)
		if err != nil {
			b.Fatal("write failed ", err)
		}
	}

	b.StopTimer()

	lgr.Close()
}
