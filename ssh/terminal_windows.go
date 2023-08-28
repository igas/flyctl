//go:build windows
// +build windows

package ssh

import (
	"context"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/term"
)

func watchWindowSize(ctx context.Context, fd int, sess *ssh.Session) error {

	// NOTE(Ali): Windows doesn't support SIGWINCH. The closest it has is WINDOW_BUFFER_SIZE_EVENT,
	// which you only seem to be able to receive if *all* of your console input is read with ReadConsoleInput.
	// (I'm also unsure how portable this is, it *might* just be a Windows Terminal thing, I didn't research too hard)
	// That's a huge undertaking, even *if* you stubbed stdin with a pipe and had a goroutine hydrating it from
	// the ReadConsoleInput data. (getting these types into go is a nightmare given the C unions, and I'm not quite
	// sure how to force everything in flyctl down the road to know that the pipe stdin is in fact a terminal)
	//
	// Because of this, we resort to the oldest trick in the book: polling! Sorry.

	width, height, err := term.GetSize(fd)
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-time.After(200 * time.Millisecond):
		}

		newWidth, newHeight, err := term.GetSize(fd)
		if err != nil {
			return err
		}

		if newWidth == width && newHeight == height {
			continue
		}

		width = newWidth
		height = newHeight

		if err := sess.WindowChange(height, width); err != nil {
			return err
		}
	}

	return nil
}
