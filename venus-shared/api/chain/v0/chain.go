package v0

import (
	"context"
	"time"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-bitfield"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/big"
	"github.com/filecoin-project/go-state-types/crypto"
	"github.com/filecoin-project/go-state-types/dline"
	"github.com/filecoin-project/go-state-types/network"
	"github.com/ipfs/go-cid"

	"github.com/filecoin-project/venus/venus-shared/actors/builtin/miner"
	chain2 "github.com/filecoin-project/venus/venus-shared/api/chain"
	"github.com/filecoin-project/venus/venus-shared/chain"
)

type IChain interface {
	IAccount
	IActor
	IBeacon
	IMinerState
	IChainInfo
}

type IAccount interface {
	// Rule[perm:read]
	StateAccountKey(ctx context.Context, addr address.Address, tsk chain.TipSetKey) (address.Address, error)
}

type IActor interface {
	// Rule[perm:read]
	StateGetActor(ctx context.Context, actor address.Address, tsk chain.TipSetKey) (*chain.Actor, error)
	// Rule[perm:read]
	ListActor(ctx context.Context) (map[address.Address]*chain.Actor, error)
}

type IBeacon interface {
	// Rule[perm:read]
	BeaconGetEntry(ctx context.Context, epoch abi.ChainEpoch) (*chain.BeaconEntry, error)
}

type IChainInfo interface {
	// Rule[perm:read]
	BlockTime(ctx context.Context) time.Duration
	// Rule[perm:read]
	ChainList(ctx context.Context, tsKey chain.TipSetKey, count int) ([]chain.TipSetKey, error)
	// Rule[perm:read]
	ChainHead(ctx context.Context) (*chain.TipSet, error)
	// Rule[perm:admin]
	ChainSetHead(ctx context.Context, key chain.TipSetKey) error
	// Rule[perm:read]
	ChainGetTipSet(ctx context.Context, key chain.TipSetKey) (*chain.TipSet, error)
	// Rule[perm:read]
	ChainGetTipSetByHeight(ctx context.Context, height abi.ChainEpoch, tsk chain.TipSetKey) (*chain.TipSet, error)
	// Rule[perm:read]
	ChainGetRandomnessFromBeacon(ctx context.Context, key chain.TipSetKey, personalization crypto.DomainSeparationTag, randEpoch abi.ChainEpoch, entropy []byte) (abi.Randomness, error)
	// Rule[perm:read]
	ChainGetRandomnessFromTickets(ctx context.Context, tsk chain.TipSetKey, personalization crypto.DomainSeparationTag, randEpoch abi.ChainEpoch, entropy []byte) (abi.Randomness, error)
	// Rule[perm:read]
	ChainGetBlock(ctx context.Context, id cid.Cid) (*chain.BlockHeader, error)
	// Rule[perm:read]
	ChainGetMessage(ctx context.Context, msgID cid.Cid) (*chain.Message, error)
	// Rule[perm:read]
	ChainGetBlockMessages(ctx context.Context, bid cid.Cid) (*chain2.BlockMessages, error)
	// Rule[perm:read]
	ChainGetMessagesInTipset(ctx context.Context, key chain.TipSetKey) ([]chain2.Message, error)
	// Rule[perm:read]
	ChainGetReceipts(ctx context.Context, id cid.Cid) ([]chain.MessageReceipt, error)
	// Rule[perm:read]
	ChainGetParentMessages(ctx context.Context, bcid cid.Cid) ([]chain2.Message, error)
	// Rule[perm:read]
	ChainGetParentReceipts(ctx context.Context, bcid cid.Cid) ([]*chain.MessageReceipt, error)
	//Rule[perm:read]
	StateVerifiedRegistryRootKey(ctx context.Context, tsk chain.TipSetKey) (address.Address, error)
	// Rule[perm:read]
	StateVerifierStatus(ctx context.Context, addr address.Address, tsk chain.TipSetKey) (*abi.StoragePower, error)
	// Rule[perm:read]
	ChainNotify(ctx context.Context) (<-chan []*chain2.HeadChange, error)
	// Rule[perm:read]
	GetFullBlock(ctx context.Context, id cid.Cid) (*chain.FullBlock, error)
	// Rule[perm:read]
	GetActor(ctx context.Context, addr address.Address) (*chain.Actor, error)
	// Rule[perm:read]
	GetParentStateRootActor(ctx context.Context, ts *chain.TipSet, addr address.Address) (*chain.Actor, error)
	// Rule[perm:read]
	GetEntry(ctx context.Context, height abi.ChainEpoch, round uint64) (*chain.BeaconEntry, error)
	// Rule[perm:read]
	MessageWait(ctx context.Context, msgCid cid.Cid, confidence, lookback abi.ChainEpoch) (*chain2.ChainMessage, error)
	// Rule[perm:read]
	ProtocolParameters(ctx context.Context) (*chain2.ProtocolParams, error)
	// Rule[perm:read]
	ResolveToKeyAddr(ctx context.Context, addr address.Address, ts *chain.TipSet) (address.Address, error)
	// Rule[perm:read]
	StateNetworkName(ctx context.Context) (chain2.NetworkName, error)
	// Rule[perm:read]
	StateGetReceipt(ctx context.Context, msg cid.Cid, from chain.TipSetKey) (*chain.MessageReceipt, error)
	// Rule[perm:read]
	StateSearchMsg(ctx context.Context, msg cid.Cid) (*chain2.MsgLookup, error)
	// Rule[perm:read]
	StateSearchMsgLimited(ctx context.Context, cid cid.Cid, limit abi.ChainEpoch) (*chain2.MsgLookup, error)
	// Rule[perm:read]
	StateWaitMsg(ctx context.Context, cid cid.Cid, confidence uint64) (*chain2.MsgLookup, error)
	// Rule[perm:read]
	StateWaitMsgLimited(ctx context.Context, cid cid.Cid, confidence uint64, limit abi.ChainEpoch) (*chain2.MsgLookup, error)
	// Rule[perm:read]
	StateNetworkVersion(ctx context.Context, tsk chain.TipSetKey) (network.Version, error)
	// Rule[perm:read]
	VerifyEntry(parent, child *chain.BeaconEntry, height abi.ChainEpoch) bool
	// Rule[perm:read]
	ChainExport(context.Context, abi.ChainEpoch, bool, chain.TipSetKey) (<-chan []byte, error)
	// Rule[perm:read]
	ChainGetPath(ctx context.Context, from chain.TipSetKey, to chain.TipSetKey) ([]*chain2.HeadChange, error)
}

type IMinerState interface {
	// Rule[perm:read]
	StateMinerSectorAllocated(ctx context.Context, maddr address.Address, s abi.SectorNumber, tsk chain.TipSetKey) (bool, error)
	// Rule[perm:read]
	StateSectorPreCommitInfo(ctx context.Context, maddr address.Address, n abi.SectorNumber, tsk chain.TipSetKey) (miner.SectorPreCommitOnChainInfo, error)
	// Rule[perm:read]
	StateSectorGetInfo(ctx context.Context, maddr address.Address, n abi.SectorNumber, tsk chain.TipSetKey) (*miner.SectorOnChainInfo, error)
	// Rule[perm:read]
	StateSectorPartition(ctx context.Context, maddr address.Address, sectorNumber abi.SectorNumber, tsk chain.TipSetKey) (*miner.SectorLocation, error)
	// Rule[perm:read]
	StateMinerSectorSize(ctx context.Context, maddr address.Address, tsk chain.TipSetKey) (abi.SectorSize, error)
	// Rule[perm:read]
	StateMinerInfo(ctx context.Context, maddr address.Address, tsk chain.TipSetKey) (miner.MinerInfo, error)
	// Rule[perm:read]
	StateMinerWorkerAddress(ctx context.Context, maddr address.Address, tsk chain.TipSetKey) (address.Address, error)
	// Rule[perm:read]
	StateMinerRecoveries(ctx context.Context, maddr address.Address, tsk chain.TipSetKey) (bitfield.BitField, error)
	// Rule[perm:read]
	StateMinerFaults(ctx context.Context, maddr address.Address, tsk chain.TipSetKey) (bitfield.BitField, error)
	// Rule[perm:read]
	StateMinerProvingDeadline(ctx context.Context, maddr address.Address, tsk chain.TipSetKey) (*dline.Info, error)
	// Rule[perm:read]
	StateMinerPartitions(ctx context.Context, maddr address.Address, dlIdx uint64, tsk chain.TipSetKey) ([]chain2.Partition, error)
	// Rule[perm:read]
	StateMinerDeadlines(ctx context.Context, maddr address.Address, tsk chain.TipSetKey) ([]chain2.Deadline, error)
	// Rule[perm:read]
	StateMinerSectors(ctx context.Context, maddr address.Address, sectorNos *bitfield.BitField, tsk chain.TipSetKey) ([]*miner.SectorOnChainInfo, error)
	// Rule[perm:read]
	StateMarketStorageDeal(ctx context.Context, dealID abi.DealID, tsk chain.TipSetKey) (*chain2.MarketDeal, error)
	// Rule[perm:read]
	StateMinerPreCommitDepositForPower(ctx context.Context, maddr address.Address, pci miner.SectorPreCommitInfo, tsk chain.TipSetKey) (big.Int, error)
	// Rule[perm:read]
	StateMinerInitialPledgeCollateral(ctx context.Context, maddr address.Address, pci miner.SectorPreCommitInfo, tsk chain.TipSetKey) (big.Int, error)
	// Rule[perm:read]
	StateVMCirculatingSupplyInternal(ctx context.Context, tsk chain.TipSetKey) (chain.CirculatingSupply, error)
	// Rule[perm:read]
	StateCirculatingSupply(ctx context.Context, tsk chain.TipSetKey) (abi.TokenAmount, error)
	// Rule[perm:read]
	StateMarketDeals(ctx context.Context, tsk chain.TipSetKey) (map[string]chain2.MarketDeal, error)
	// Rule[perm:read]
	StateMinerActiveSectors(ctx context.Context, maddr address.Address, tsk chain.TipSetKey) ([]*miner.SectorOnChainInfo, error)
	// Rule[perm:read]
	StateLookupID(ctx context.Context, addr address.Address, tsk chain.TipSetKey) (address.Address, error)
	// Rule[perm:read]
	StateListMiners(ctx context.Context, tsk chain.TipSetKey) ([]address.Address, error)
	// Rule[perm:read]
	StateListActors(ctx context.Context, tsk chain.TipSetKey) ([]address.Address, error)
	// Rule[perm:read]
	StateMinerPower(ctx context.Context, addr address.Address, tsk chain.TipSetKey) (*chain2.MinerPower, error)
	// Rule[perm:read]
	StateMinerAvailableBalance(ctx context.Context, maddr address.Address, tsk chain.TipSetKey) (big.Int, error)
	// Rule[perm:read]
	StateSectorExpiration(ctx context.Context, maddr address.Address, sectorNumber abi.SectorNumber, tsk chain.TipSetKey) (*miner.SectorExpiration, error)
	// Rule[perm:read]
	StateMinerSectorCount(ctx context.Context, addr address.Address, tsk chain.TipSetKey) (chain2.MinerSectors, error)
	// Rule[perm:read]
	StateMarketBalance(ctx context.Context, addr address.Address, tsk chain.TipSetKey) (chain2.MarketBalance, error)
	// Rule[perm:read]
	StateDealProviderCollateralBounds(ctx context.Context, size abi.PaddedPieceSize, verified bool, tsk chain.TipSetKey) (chain2.DealCollateralBounds, error)
	// Rule[perm:read]
	StateVerifiedClientStatus(ctx context.Context, addr address.Address, tsk chain.TipSetKey) (*abi.StoragePower, error)
}