package main

import (
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
		want     string
	}{
		{
			name:     "1 hour, 2 minutes, 3 seconds",
			duration: time.Hour*1 + time.Minute*2 + time.Second*3,
			want:     "1h 2m 3s",
		},
		{
			name:     "2 minutes, 3 seconds",
			duration: time.Minute*2 + time.Second*3,
			want:     "2m 3s",
		},
		{
			name:     "3 seconds",
			duration: time.Second * 3,
			want:     "3s",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := formatDuration(tt.duration); got != tt.want {
				t.Errorf("formatDuration() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFindByName(t *testing.T) {

	tmpDir, err := os.MkdirTemp("", "sample")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	files := []string{"salamander.txt", "steakrecipe.txt", "randomfile.txt"}
	for _, file := range files {
		_, err := os.Create(filepath.Join(tmpDir, file))
		if err != nil {
			t.Fatal(err)
		}
	}

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	findByName("salamander", tmpDir, false)

	w.Close()
	os.Stdout = oldStdout

	out, _ := io.ReadAll(r)

	expected := "0) " + filepath.Join(tmpDir, "salamander.txt") + "\n\n1 files found in 0s\n"
	if string(out) != expected {
		t.Errorf("Expected %s, got %s", expected, out)
	}
}
