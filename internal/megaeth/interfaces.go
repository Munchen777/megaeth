package megaeth

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	accTypes "main/pkg/types"
	"main/pkg/global"
)

type NFT struct {
	Value *big.Int
	ContractAddress string
	DisplayName string
}

type AngryMonkeysNFT struct { NFT }

type BlackholeNFT struct { NFT }

type BloomNFT struct { NFT }

type LapinNFT struct { NFT }

type MegaCatNFT struct { NFT }

type MegamafiaNFT struct { NFT }

type XyrophNFT struct { NFT }

type FunNFT struct {
	NFT
	Data string
}

type Mint interface {
	Mint(context.Context, accTypes.AccountData) (bool, error)
}

func GetInterfaceByModuleName() (Mint, error) {
	switch global.Module {
	case "Mint Megamafia NFT":
		return MegamafiaNFT{
			NFT: NFT{
				Value: big.NewInt(0),
				ContractAddress: "0xa3C89fEb775940886001E8f541f4b803AaD0a47B",
				DisplayName: "Mint Megamafia NFT",
			},
		}, nil
	case "Mint Mega Cat NFT":
		return MegaCatNFT{
			NFT: NFT{
				Value: big.NewInt(0),
				ContractAddress: "0x0837ec39d40CCdcea4b4B6bfCfb3d71E7EbFC71C",
				DisplayName: "Mint Mega Cat NFT",
			},
		}, nil
	case "Mint Blackhole NFT":
		return BlackholeNFT{
			NFT: NFT{
				Value: big.NewInt(0),
				ContractAddress: "0xcfD3dDe3A4B393a2a204ff16B112C2cA9B85abb7",
				DisplayName: "Mint Blackhole NFT",
			},
		}, nil
	case "Mint Xyroph NFT":
		return XyrophNFT{
			NFT: NFT{
				Value: big.NewInt(1550000000000000),
				ContractAddress: "0xd59522848e5429986d6fe6607aef6b8e7706aea5",
				DisplayName: "Mint Xyroph NFT",
			},
		}, nil
	case "Mint Lord Lapin NFT":
		return LapinNFT{
			NFT: NFT{
				Value: big.NewInt(0),
				ContractAddress: "0x0d7BEa5686E3c85cb018faa066AB36CF00b63eBB",
				DisplayName: "Mint Lord Lapin NFT",
			},
		}, nil
	case "Mint Angry Monkeys":
		return AngryMonkeysNFT{
			NFT: NFT{
				Value: big.NewInt(1440000000000000),
				ContractAddress: "0x8ac06714c0d417569bcc642cd74e48a64fe99504",
				DisplayName: "Mint Angry Monkeys",
			},
		}, nil
	case "Mint Bloom NFT":
		return BloomNFT{
			NFT: NFT{
				Value: big.NewInt(0),
				ContractAddress: "0xb33C085f82B253B12a9d36F8E8EdD123FFB53d31",
				DisplayName: "Mint Bloom NFT",
			},
		}, nil
	case "Mint FUN Starts NFT":
		return FunNFT{
			NFT: NFT{
				Value: big.NewInt(0),
				ContractAddress: "0xb8027dca96746f073896c45f65b720f9bd2afee7",
				DisplayName: "Mint Fun NFT",
			},
			Data: "0x1249c58b",
		}, nil
	default:
		msg := fmt.Sprintf("Struct with %s module name is not defined!", global.Module)
		return nil, errors.New(msg)
	}
}
