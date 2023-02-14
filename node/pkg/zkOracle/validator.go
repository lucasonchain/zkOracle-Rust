package zkOracle

import (
	"context"
	"fmt"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr/mimc"
	"github.com/consensys/gnark-crypto/ecc/bn254/twistededwards/eddsa"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const CONFIRMATIONS uint64 = 5

type Validator struct {
	ethClient        *ethclient.Client
	zkOracleContract *ZKOracleContract
	privateKey       *eddsa.PrivateKey
}

func NewValidator(ethClient *ethclient.Client, zkOracleContract *ZKOracleContract, privateKey *eddsa.PrivateKey) *Validator {
	return &Validator{ethClient: ethClient, zkOracleContract: zkOracleContract, privateKey: privateKey}
}

func (v *Validator) Validate(ctx context.Context) error {
	if err := v.WatchAndHandleBlockRequestedEvent(ctx); err != nil {
		return fmt.Errorf("watch and handle block requested events: %w")
	}
	return nil
}

func (v *Validator) WatchAndHandleBlockRequestedEvent(ctx context.Context) error {
	sink := make(chan *ZKOracleContractBlockRequested)
	defer close(sink)

	sub, err := v.zkOracleContract.WatchBlockRequested(&bind.WatchOpts{
		Context: ctx,
	}, sink, nil, nil)
	if err != nil {
		return err
	}
	defer sub.Unsubscribe()

	for {
		select {
		case event := <-sink:
			if err := v.HandleBlockRequestedEvent(ctx, event); err != nil {
				logger.Err(err).
					Uint64("requestNumber", event.Request.Uint64()).
					Uint64("blockNumber", event.Number.Uint64()).
					Msg("handle block requested event")
			}
		case err = <-sub.Err():
			return err
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (v *Validator) HandleBlockRequestedEvent(ctx context.Context, event *ZKOracleContractBlockRequested) error {

	logger.Info().
		Uint64("requestNumber", event.Request.Uint64()).
		Uint64("blockNumber", event.Number.Uint64()).
		Msg("handle block requested event")

	block, err := v.ethClient.HeaderByNumber(ctx, event.Number)
	if err != nil {
		return fmt.Errorf("block by number: %w", err)
	}

	currentBlockNumber, err := v.ethClient.BlockNumber(ctx)
	if err != nil {
		return fmt.Errorf("blocknumber: %w", err)
	}

	logger.Info().
		Uint64("requestNumber", event.Request.Uint64()).
		Uint64("blockNumber", event.Number.Uint64()).
		Uint64("head", currentBlockNumber).
		Msg("check block confirmed")

	if currentBlockNumber-block.Number.Uint64() < CONFIRMATIONS {
		return fmt.Errorf("block not confirmed")
	}

	sig, err := v.privateKey.Sign(block.Hash().Bytes(), mimc.NewMiMC())
	if err != nil {
		return fmt.Errorf("sign: %w", err)
	}

	i, err := v.zkOracleContract.GetAggregator(
		&bind.CallOpts{
			Context: ctx,
		},
	)
	if err != nil {
		return fmt.Errorf("get aggregator: %w", err)
	}

	addr, err := v.zkOracleContract.GetIPAddress(&bind.CallOpts{Context: ctx}, i)
	if err != nil {
		return fmt.Errorf("get ip addr: %w", err)
	}

	logger.Info().
		Uint64("requestNumber", event.Request.Uint64()).
		Uint64("index", i.Uint64()).
		Str("ipAddr", addr).
		Msg("sending vote to aggregator")

	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("dial %s: %w", addr, err)
	}

	client := NewOracleNodeClient(conn)
	_, err = client.SendVote(ctx, &SendVoteRequest{
		Request:   event.Request.Uint64(),
		BlockHash: block.Hash().Bytes(),
		Signature: sig,
	})
	if err != nil {
		return fmt.Errorf("get ip addr: %w", err)
	}

	logger.Info().
		Uint64("requestNumber", event.Request.Uint64()).
		Msg("received response")

	return nil
}
