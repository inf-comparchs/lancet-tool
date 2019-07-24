/*
 * MIT License
 *
 * Copyright (c) 2019-2021 Ecole Polytechnique Federale Lausanne (EPFL)
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */
package main

// #include "../inc/lancet/coord_proto.h"
import "C"
import (
	"fmt"
	"os"
	"strings"
)

func computeStatsThroughput(replies []*C.struct_throughput_reply) *C.struct_throughput_reply {
	agg_stats := &C.struct_throughput_reply{}
	for _, r := range replies {
		agg_stats.Rx_bytes += r.Rx_bytes
		agg_stats.Tx_bytes += r.Tx_bytes
		agg_stats.Req_count += r.Req_count
		agg_stats.CorrectIAD += r.CorrectIAD
	}
	agg_stats.Duration = replies[0].Duration

	return agg_stats
}

func computeStatsLatency(replies []*C.struct_latency_reply) *C.struct_latency_reply {
	agg_stats := &C.struct_latency_reply{}
	for _, r := range replies {
		//printLatencyStats(r) FIXME
		agg_stats.Avg_lat += r.Avg_lat
		agg_stats.P50_i += r.P50_i
		agg_stats.P50 += r.P50
		agg_stats.P50_k += r.P50_k
		agg_stats.P90_i += r.P90_i
		agg_stats.P90 += r.P90
		agg_stats.P90_k += r.P90_k
		agg_stats.P95_i += r.P95_i
		agg_stats.P95 += r.P95
		agg_stats.P95_k += r.P95_k
		agg_stats.P99_i += r.P99_i
		agg_stats.P99 += r.P99
		agg_stats.P99_k += r.P99_k
		agg_stats.IsStationary += r.IsStationary
		agg_stats.IsIid += r.IsIid
	}
	agg_stats.Avg_lat /= C.uint64_t(len(replies))
	agg_stats.P50_i /= C.uint64_t(len(replies))
	agg_stats.P50 /= C.uint64_t(len(replies))
	agg_stats.P50_k /= C.uint64_t(len(replies))
	agg_stats.P90_i /= C.uint64_t(len(replies))
	agg_stats.P90 /= C.uint64_t(len(replies))
	agg_stats.P90_k /= C.uint64_t(len(replies))
	agg_stats.P95_i /= C.uint64_t(len(replies))
	agg_stats.P95 /= C.uint64_t(len(replies))
	agg_stats.P95_k /= C.uint64_t(len(replies))
	agg_stats.P99_i /= C.uint64_t(len(replies))
	agg_stats.P99 /= C.uint64_t(len(replies))
	agg_stats.P99_k /= C.uint64_t(len(replies))

	if agg_stats.IsIid == 0 {
		agg_stats.ToReduceSampling = 1000000
		for _, r := range replies {
			if r.ToReduceSampling < agg_stats.ToReduceSampling {
				agg_stats.ToReduceSampling = r.ToReduceSampling
			}
		}
	}
	return agg_stats

}

func printThroughputStats(stats *C.struct_throughput_reply, outCsv string) error {
	f, err := os.OpenFile(outCsv, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("Can't open outCsv file for writing\n")
	}
	defer f.Close()

	str := "#ReqCount\tQPS\tRxBw\tTxBw\n"
	fmt.Println( strings.TrimSuffix(str, "\n") )
	if _, err := f.WriteString(str); err != nil {
		return fmt.Errorf("Can't write to outCsv file\n")
	}
	str = fmt.Sprintf("%v\t%v\t%v\t%v\n", stats.Req_count,
	1e6*float64(stats.Req_count)/float64(stats.Duration),
	1e6*float64(stats.Rx_bytes)/float64(stats.Duration),
	1e6*float64(stats.Tx_bytes)/float64(stats.Duration))
	fmt.Printf(str)
	if _, err := f.WriteString(str); err != nil {
		return fmt.Errorf("Can't write to outCsv file\n")
	}

	return nil
}

func printLatencyStats(stats *C.struct_latency_reply, outCsv string) error {
	f, err := os.OpenFile(outCsv, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("Can't open outCsv file for writing\n")
	}
	defer f.Close()

	str := "#Avg Lat\t50th\t90th\t95th\t99th\n"
	fmt.Println( strings.TrimSuffix(str, "\n") )
	if _, err := f.WriteString(str); err != nil {
		return fmt.Errorf("Can't write to outCsv file\n")
	}
	str = fmt.Sprintf("%v\t%v(%v, %v)\t%v(%v, %v)\t%v(%v, %v)\t%v(%v, %v)\n",
	float64(stats.Avg_lat)/1e3,
	float64(stats.P50)/1e3, float64(stats.P50_i)/1e3, float64(stats.P50_k)/1e3,
	float64(stats.P90)/1e3, float64(stats.P90_i)/1e3, float64(stats.P90_k)/1e3,
	float64(stats.P95)/1e3, float64(stats.P95_i)/1e3, float64(stats.P95_k)/1e3,
	float64(stats.P99)/1e3, float64(stats.P99_i)/1e3, float64(stats.P99_k)/1e3)
	fmt.Printf(str)
	if _, err := f.WriteString(str); err != nil {
		return fmt.Errorf("Can't write to outCsv file\n")
	}

	return nil
}

func getRPS(stats *C.struct_throughput_reply) float64 {
	return 1e6 * float64(stats.Req_count) / float64(stats.Duration)
}

func getLatCISize(stats []*C.struct_latency_reply, percentile int) int {
	if percentile == 99 {
		return int(stats[0].P99_k - stats[0].P99_i)
	} else {
		panic("Unknown percentile")
	}
}
