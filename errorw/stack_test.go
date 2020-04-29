package errorw

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"testing"
)

var e = NewMessage("testing error").
	WithField("foo", "bar").
	WithField("struct", testStruct{Value: "value"}).
	WithWrap("wrap")

var e2 = New(errors.New("testing error"))

type testStruct struct {
	Value string
}

func TestError_StackTrace(t *testing.T) {
	// e
	st := e.StackTrace()

	want := []string{
		"github.com/XSAM/go-hybrid/errorw.init\n\t.*go-hybrid/errorw/stack_test.go:11",
	}
	for i, w := range want {
		testFormatRegexp(t, i, st[i], "%+v", w)
	}

	// e2
	st = e2.StackTrace()

	want = []string{
		"github.com/XSAM/go-hybrid/errorw.init\n\t.*go-hybrid/errorw/stack_test.go:16",
	}
	for i, w := range want {
		testFormatRegexp(t, i, st[i], "%+v", w)
	}
}

func TestStack_Format(t *testing.T) {
	testCases := []struct {
		format string
		want   string
	}{
		{
			format: "%s",
			want:   "",
		},
		{
			format: "%v",
			want:   "",
		},
		{
			format: "%+v",
			want:   "\ngithub.com/XSAM/go-hybrid/errorw.init\n\t.*go-hybrid/errorw/stack_test.go:11",
		},
	}

	for i, tc := range testCases {
		t.Run(tc.format, func(t *testing.T) {
			testFormatRegexp(t, i, e.Stack, tc.format, tc.want)
		})
	}
}

func testFormatRegexp(t *testing.T, n int, arg interface{}, format, want string) {
	got := fmt.Sprintf(format, arg)
	gotLines := strings.SplitN(got, "\n", -1)
	wantLines := strings.SplitN(want, "\n", -1)

	if len(wantLines) > len(gotLines) {
		t.Errorf("test %d: wantLines(%d) > gotLines(%d):\n got: %q\nwant: %q", n+1, len(wantLines), len(gotLines), got, want)
		return
	}

	for i, w := range wantLines {
		match, err := regexp.MatchString(w, gotLines[i])
		if err != nil {
			t.Fatal(err)
		}
		if !match {
			t.Errorf("test %d: line %d: fmt.Sprintf(%q, err):\n got: %q\nwant: %q", n+1, i+1, format, got, want)
		}
	}
}
