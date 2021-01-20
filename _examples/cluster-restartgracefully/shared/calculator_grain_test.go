package shared

import (
	fmt "fmt"
	"strconv"
	"testing"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/cluster"
	"github.com/AsynkronIT/protoactor-go/cluster/consul"
	"github.com/AsynkronIT/protoactor-go/cluster/etcd"
	"github.com/AsynkronIT/protoactor-go/log"
	"github.com/AsynkronIT/protoactor-go/remote"
	"github.com/stretchr/testify/assert"
)

var (
	system   = actor.NewActorSystem()
	_cluster *cluster.Cluster
)

func startCluster(port int) {
	var cp cluster.ClusterProvider
	var err error
	var provider = "etcd"
	switch provider {
	case "consul":
		cp, err = consul.New()
	case "etcd":
		cp, err = etcd.New()
	default:
		panic(fmt.Errorf("Invalid provider:%s", provider))
	}

	if err != nil {
		panic(err)
	}

	remoteCfg := remote.Configure("127.0.0.1", port)
	cfg := cluster.Configure("cluster-restartgracefully", cp, remoteCfg, Kind)
	_cluster = cluster.New(system, cfg)
	_cluster.Start()
}

func stopCluster() {
	_cluster.Shutdown(true)
}

func TestCalculatorGrain(t *testing.T) {
	actor.SetLogLevel(log.ErrorLevel)
	remote.SetLogLevel(log.ErrorLevel)
	cluster.SetLogLevel(log.ErrorLevel)
	startCluster(0)
	t.Cleanup(func() {
		stopCluster()
	})

	assert := assert.New(t)

	calcGrain := GetCalculatorGrainClient(_cluster, "client-1")
	resp, err := calcGrain.GetCurrent(&Void{})
	assert.NoError(err)
	expected := resp.GetNumber()

	for i := 0; i < 10000; i++ {
		// add
		{
			resp, err := calcGrain.Add(&NumberRequest{int64(i)})
			assert.NoError(err)

			expected += int64(i)
			actual := resp.GetNumber()
			assert.Equal(expected, actual)
		}

		// subtract
		{
			resp, err := calcGrain.Subtract(&NumberRequest{int64(i)})
			assert.NoError(err)

			expected -= int64(i)
			actual := resp.GetNumber()
			assert.Equal(expected, actual)
		}

		// current
		{
			resp, err := calcGrain.GetCurrent(&Void{})
			assert.NoError(err)

			actual := resp.GetNumber()
			assert.Equal(expected, actual)
		}
	}
}

func BenchmarkGetCalculatorGrainClient(b *testing.B) {
	actor.SetLogLevel(log.ErrorLevel)
	remote.SetLogLevel(log.ErrorLevel)
	cluster.SetLogLevel(log.ErrorLevel)
	SetLogLevel(log.ErrorLevel)
	startCluster(0)

	b.Cleanup(func() {
		stopCluster()
	})

	b.Run("with-cache", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			grainId := "client-" + strconv.Itoa(1)
			calcGrain := GetCalculatorGrainClient(_cluster, grainId)
			_ = calcGrain
		}
	})

	b.Run("without-cache", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			grainId := "client-" + strconv.Itoa(i)
			calcGrain := GetCalculatorGrainClient(_cluster, grainId)
			_ = calcGrain
		}
	})
}

func BenchmarkCalculatorGrain_Add(b *testing.B) {
	actor.SetLogLevel(log.ErrorLevel)
	remote.SetLogLevel(log.ErrorLevel)
	cluster.SetLogLevel(log.ErrorLevel)
	SetLogLevel(log.ErrorLevel)
	startCluster(0)

	b.Cleanup(func() {
		stopCluster()
	})

	// b.Run("with-cache", func(b *testing.B) {
	// 	for i := 0; i < b.N; i++ {
	// 		grainId := "client-" + strconv.Itoa(int(i))
	// 		GetCalculatorGrainClient(_cluster, grainId)
	// 	}
	// 	b.ResetTimer()

	// 	for i := 0; i < b.N; i++ {
	// 		grainId := "client-" + strconv.Itoa(int(i))
	// 		calcGrain := GetCalculatorGrainClient(_cluster, grainId)
	// 		resp, err := calcGrain.Add(&NumberRequest{int64(i)})
	// 		if err != nil {
	// 			b.FailNow()
	// 		}
	// 		_ = resp
	// 	}

	// })

	b.Run("1", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			grainId := "client-" + strconv.Itoa(int(i))
			calcGrain := GetCalculatorGrainClient(_cluster, grainId)
			resp, err := calcGrain.Add(&NumberRequest{int64(i)})
			if err != nil {
				b.FailNow()
			}
			_ = resp
		}
	})

	// b.Run("N", func(b *testing.B) {
	// 	var i int64
	// 	b.RunParallel(func(pb *testing.PB) {
	// 		for pb.Next() {
	// 			grainId := "client-" + strconv.Itoa(int(i))
	// 			atomic.AddInt64(&i, 1)

	// 			calcGrain := GetCalculatorGrainClient(_cluster, grainId)
	// 			resp, err := calcGrain.Add(&NumberRequest{int64(i)})
	// 			if err != nil {
	// 				b.FailNow()
	// 			}
	// 			_ = resp
	// 		}
	// 	})
	// })
}
