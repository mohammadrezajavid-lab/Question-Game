package metrics

import "github.com/prometheus/client_golang/prometheus"

var EncodeToProtobufFailedCounter = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "encode_to_protobuf_failed_total",
		Help: "Total number of struct encode to protobuf failed",
	}, []string{"encoder_name"},
)

var DecodeFromProtobufFailedCounter = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "decode_from_protobuf_failed_total",
		Help: "Total number of decode protobuf to struct failed",
	}, []string{"decoder_name"},
)

func init() {
	Registry.MustRegister(
		EncodeToProtobufFailedCounter,
		DecodeFromProtobufFailedCounter,
	)
}
