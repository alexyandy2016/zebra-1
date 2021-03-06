package zuul

import (
    "fmt"
    "testing"
    . "github.com/smartystreets/goconvey/convey"
    "github.com/tietang/stats"
)

func TestRobinTheSameWeight(t *testing.T) {
    Convey("给定同样的权重", t, func() {

        var hosts1 = make([]HostInstance, 0)

        hosts1 = append(hosts1, HostInstance{Name: "a", Weight: 1})
        hosts1 = append(hosts1, HostInstance{Name: "b", Weight: 1})
        hosts1 = append(hosts1, HostInstance{Name: "c", Weight: 1})

        r := Robin{}
        var h *HostInstance
        Convey("next方法调用，分布准确", func() {
            h = r.Next(hosts1)
            So(hosts1[0].Name, ShouldEqual, h.Name)
            h = r.Next(hosts1)
            So(hosts1[1].Name, ShouldEqual, h.Name)
            h = r.Next(hosts1)
            So(hosts1[2].Name, ShouldEqual, h.Name)
            h = r.Next(hosts1)
            So(hosts1[0].Name, ShouldEqual, h.Name)
        })

        times := 100
        size := len(hosts1) * times
        counter := stats.NewCounter()

        for i := 0; i < size; i++ {
            h = r.Next(hosts1)
            counter.Incr(h.Name, 1)
        }
        Convey("调用100次后的统计", func() {
            So(int64(100), ShouldEqual, counter.Get(hosts1[0].Name).Count)
            So(int64(100), ShouldEqual, counter.Get(hosts1[1].Name).Count)
            So(int64(100), ShouldEqual, counter.Get(hosts1[2].Name).Count)

            fmt.Println(h.Name, h.Weight, hosts1)
        })

    })

}

func TestRobinDiffWeight(t *testing.T) {
    Convey("权重不同", t, func() {

        var hosts1 = make([]HostInstance, 0)

        hosts1 = append(hosts1, HostInstance{Name: "a", Weight: 1})
        hosts1 = append(hosts1, HostInstance{Name: "b", Weight: 2})
        hosts1 = append(hosts1, HostInstance{Name: "c", Weight: 3})

        r := Robin{}

        times := 100
        //总执行次数
        size := getTotalWeight(hosts1) * times
        counter := stats.NewCounter()
        var h *HostInstance
        for i := 0; i < size; i++ {
            h = r.Next(hosts1)
            counter.Incr(h.Name, 1)
        }
        Convey("各个权重命中是否正确", func() {
            //各个权重命中是否正确
            So(int64(100 * hosts1[0].Weight), ShouldEqual, counter.Get(hosts1[0].Name).Count)
            So(int64(100 * hosts1[1].Weight), ShouldEqual, counter.Get(hosts1[1].Name).Count)
            So(int64(100 * hosts1[2].Weight), ShouldEqual, counter.Get(hosts1[2].Name).Count)
        })
    })
}

func getTotalWeight(hosts []HostInstance) int {
    c := 0
    for _, v := range hosts {
        c = c + v.Weight
    }
    return c
}
