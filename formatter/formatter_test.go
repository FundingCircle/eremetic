package formatter

import (
	"testing"
	"time"

	types "github.com/klarna/eremetic/types"
	. "github.com/smartystreets/goconvey/convey"
)

func TestFormatter(t *testing.T) {
	Convey("FormatTime", t, func() {
		Convey("A Valid Unix Timestamp", func() {
			t := time.Now().Unix()
			So(FormatTime(t), ShouldNotBeEmpty)
		})
	})
	Convey("ToLower", t, func() {
		Convey("A Lowercased String", func() {
			taskState := types.TaskState{"My_State"}
			So(ToLower(taskState), ShouldEqual, "my_state")
		})
	})
}
