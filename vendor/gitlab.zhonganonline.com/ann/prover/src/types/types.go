package types

import (
	pbtypes "gitlab.zhonganonline.com/ann/angine/protos/types"
	agtypes "gitlab.zhonganonline.com/ann/angine/types"
	"gitlab.zhonganonline.com/ann/ann-module/lib/go-crypto"
)

type (
	// Application embeds types.Application, defines application interface in Prover
	Application interface {
		agtypes.Application
		SetCore(Core)
	}

	// Core defines the interface at which an application sees its containing organization
	Core interface {
		IsValidator() bool
		GetPublicKey() (crypto.PubKeyEd25519, bool)
		GetPrivateKey() (crypto.PrivKeyEd25519, bool)
		GetChainID() string
		GetEngine() Engine
	}

	// Engine defines the consensus engine
	Engine interface {
		GetBlock(agtypes.INT) (*agtypes.BlockCache, *pbtypes.BlockMeta, error)
		GetBlockMeta(agtypes.INT) (*pbtypes.BlockMeta, error)
		GetValidators() (agtypes.INT, *agtypes.ValidatorSet)
		PrivValidator() *agtypes.PrivValidator
		BroadcastTx([]byte) error
		Query(byte, []byte) (interface{}, error)
	}

	// Broadcaster means we can deliver tx in application
	Broadcaster interface {
		BroadcastTx([]byte) error
	}

	// Serializable transforms to bytes
	Serializable interface {
		ToBytes() ([]byte, error)
	}

	// Unserializable transforms from bytes
	Unserializable interface {
		FromBytes(bs []byte)
	}

	// Hashable aliases Serializable
	Hashable interface {
		Serializable
	}
)
