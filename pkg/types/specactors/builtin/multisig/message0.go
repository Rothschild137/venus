package multisig

import (
	"golang.org/x/xerrors"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"

	builtin0 "github.com/filecoin-project/specs-actors/actors/builtin"
	init0 "github.com/filecoin-project/specs-actors/actors/builtin/init"
	multisig0 "github.com/filecoin-project/specs-actors/actors/builtin/multisig"

	"github.com/filecoin-project/venus/pkg/types/internal"
	"github.com/filecoin-project/venus/pkg/types/specactors"
	init_ "github.com/filecoin-project/venus/pkg/types/specactors/builtin/init"
)

type message0 struct{ from address.Address }

func (m message0) Create(
	signers []address.Address, threshold uint64,
	unlockStart, unlockDuration abi.ChainEpoch,
	initialAmount abi.TokenAmount,
) (*internal.Message, error) {

	lenAddrs := uint64(len(signers))

	if lenAddrs < threshold {
		return nil, xerrors.Errorf("cannot require signing of more addresses than provided for multisig")
	}

	if threshold == 0 {
		threshold = lenAddrs
	}

	if m.from == address.Undef {
		return nil, xerrors.Errorf("must provide source address")
	}

	if unlockStart != 0 {
		return nil, xerrors.Errorf("actors v0 does not support a non-zero vesting start time")
	}

	// Set up constructor parameters for multisig
	msigParams := &multisig0.ConstructorParams{
		Signers:               signers,
		NumApprovalsThreshold: threshold,
		UnlockDuration:        unlockDuration,
	}

	enc, actErr := specactors.SerializeParams(msigParams)
	if actErr != nil {
		return nil, actErr
	}

	// new actors are created by invoking 'exec' on the init actor with the constructor params
	execParams := &init0.ExecParams{
		CodeCID:           builtin0.MultisigActorCodeID,
		ConstructorParams: enc,
	}

	enc, actErr = specactors.SerializeParams(execParams)
	if actErr != nil {
		return nil, actErr
	}

	return &internal.Message{
		To:     init_.Address,
		From:   m.from,
		Method: builtin0.MethodsInit.Exec,
		Params: enc,
		Value:  initialAmount,
	}, nil
}

func (m message0) Propose(msig, to address.Address, amt abi.TokenAmount,
	method abi.MethodNum, params []byte) (*internal.Message, error) {

	if msig == address.Undef {
		return nil, xerrors.Errorf("must provide a multisig address for proposal")
	}

	if to == address.Undef {
		return nil, xerrors.Errorf("must provide a target address for proposal")
	}

	if amt.Sign() == -1 {
		return nil, xerrors.Errorf("must provide a non-negative amount for proposed send")
	}

	if m.from == address.Undef {
		return nil, xerrors.Errorf("must provide source address")
	}

	enc, actErr := specactors.SerializeParams(&multisig0.ProposeParams{
		To:     to,
		Value:  amt,
		Method: method,
		Params: params,
	})
	if actErr != nil {
		return nil, xerrors.Errorf("failed to serialize parameters: %w", actErr)
	}

	return &internal.Message{
		To:     msig,
		From:   m.from,
		Value:  abi.NewTokenAmount(0),
		Method: builtin0.MethodsMultisig.Propose,
		Params: enc,
	}, nil
}

func (m message0) Approve(msig address.Address, txID uint64, hashData *ProposalHashData) (*internal.Message, error) {
	enc, err := txnParams(txID, hashData)
	if err != nil {
		return nil, err
	}

	return &internal.Message{
		To:     msig,
		From:   m.from,
		Value:  internal.NewInt(0),
		Method: builtin0.MethodsMultisig.Approve,
		Params: enc,
	}, nil
}

func (m message0) Cancel(msig address.Address, txID uint64, hashData *ProposalHashData) (*internal.Message, error) {
	enc, err := txnParams(txID, hashData)
	if err != nil {
		return nil, err
	}

	return &internal.Message{
		To:     msig,
		From:   m.from,
		Value:  internal.NewInt(0),
		Method: builtin0.MethodsMultisig.Cancel,
		Params: enc,
	}, nil
}
