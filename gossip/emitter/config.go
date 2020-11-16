package emitter

import (
	"math/rand"
	"time"

	"github.com/Fantom-foundation/lachesis-base/inter/idx"
	"github.com/ethereum/go-ethereum/common"
	_params "github.com/ethereum/go-ethereum/params"

	"github.com/Fantom-foundation/go-opera/opera/params"
)

// EmitIntervals is the configuration of emit intervals.
type EmitIntervals struct {
	Min                        time.Duration
	Max                        time.Duration
	Confirming                 time.Duration // emit time when there's no txs to originate, but at least 1 tx to confirm
	ParallelInstanceProtection time.Duration
	DoublesignProtection       time.Duration
}

// Config is the configuration of events emitter.
type Config struct {
	VersionToPublish string

	Validator common.Address `json:"validator"`

	EmitIntervals EmitIntervals `json:"emitIntervals"` // event emission intervals

	MaxGasRateGrowthFactor float64 `json:"maxGasRateGrowthFactor"` // fine to use float, because no need in determinism

	MaxTxsFromSender int `json:"maxTxsFromSender"`

	EpochTailLength idx.Frame `json:"epochTailLength"` // number of frames before event is considered epoch

	MaxParents int `json:"maxParents"`

	// thresholds on GasLeft
	SmoothTpsThreshold uint64 `json:"smoothTpsThreshold"`
	NoTxsThreshold     uint64 `json:"noTxsThreshold"`
	EmergencyThreshold uint64 `json:"emergencyThreshold"`
}

// DefaultEmitterConfig returns the default configurations for the events emitter.
func DefaultEmitterConfig() Config {
	return Config{
		VersionToPublish: _params.VersionWithMeta(),

		EmitIntervals: EmitIntervals{
			Min:                        200 * time.Millisecond,
			Max:                        12 * time.Minute,
			Confirming:                 200 * time.Millisecond,
			DoublesignProtection:       30 * time.Minute, // should be at least 2x of MaxEmitInterval
			ParallelInstanceProtection: 1 * time.Minute,
		},

		MaxGasRateGrowthFactor: 3.0,
		MaxTxsFromSender:       TxTurnNonces,
		EpochTailLength:        1,

		MaxParents: 7,

		SmoothTpsThreshold: (params.EventGas + params.TxGas) * 500,
		NoTxsThreshold:     params.EventGas * 30,
		EmergencyThreshold: params.EventGas * 5,
	}
}

// RandomizeEmitTime and return new config
func (cfg *EmitIntervals) RandomizeEmitTime(r *rand.Rand) *EmitIntervals {
	config := *cfg
	// value = value - 0.1 * value + 0.1 * random value
	if config.Max > 10 {
		config.Max = config.Max - config.Max/10 + time.Duration(r.Int63n(int64(config.Max/10)))
	}
	// value = value + 0.1 * random value
	if config.DoublesignProtection > 10 {
		config.DoublesignProtection = config.DoublesignProtection + time.Duration(r.Int63n(int64(config.DoublesignProtection/10)))
	}
	return &config
}

// FakeEmitterConfig returns the testing configurations for the events emitter.
func FakeEmitterConfig() Config {
	cfg := DefaultEmitterConfig()
	cfg.EmitIntervals.Max = 10 * time.Second // don't wait long in fakenet
	cfg.EmitIntervals.DoublesignProtection = cfg.EmitIntervals.Max * 3 / 2
	return cfg
}