package main

import (
	"testing"

	"github.com/cloudfoundry/cli/plugin/fakes"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	cliConn *fakes.FakeCliConnection
)

func TestNoApp(t *testing.T) {
	setup()
	Convey("checkArgs should not return error with search-logs test", t, func() {
		err := checkArgs(cliConn, []string{"search-logs", "test"})
		So(err, ShouldBeNil)
	})

	Convey("checkArgs should return error with search-logs", t, func() {
		err := checkArgs(cliConn, []string{"search-logs"})
		So(err, ShouldNotBeNil)
	})

}

func TestAppGuid(t *testing.T) {
	Convey("findAppGuid should not return nothing", t, func() {
		err := findAppGuid(cliConn, "test")
		So(err, ShouldNotBeNil)
	})
}

func TestNoRoutes(t *testing.T) {
	Convey("getUrlFromOuput should return nil if route exists", t, func() {
		input := []string{"urls: google.com"}
		out, err := getUrlFromOutput(input)
		So(err, ShouldBeNil)
		So(out[0], ShouldEqual, "http://google.com")
	})

	Convey("getUrlFromOuput should return error if no route exists", t, func() {
		input := []string{"urls: "}
		out, err := getUrlFromOutput(input)
		So(err, ShouldNotBeNil)
		So(out[0], ShouldEqual, "")
	})

	Convey("getUrlFromOuput should handle multiple routes", t, func() {
		input := []string{"urls: google.com, apple.com, github.com"}
		out, err := getUrlFromOutput(input)
		So(err, ShouldBeNil)
		So(out[0], ShouldEqual, "http://google.com")
		So(out[1], ShouldEqual, "http://apple.com")
		So(out[2], ShouldEqual, "http://github.com")
	})
}

func setup() {
	cliConn = &fakes.FakeCliConnection{}
}
