// Copyright (c) 2014-2016 Bitmark Inc.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package peer

import (
	"bytes"
	"encoding/hex"
	"github.com/bitmark-inc/bitmarkd/announce"
	"github.com/bitmark-inc/bitmarkd/asset"
	"github.com/bitmark-inc/bitmarkd/block"
	"github.com/bitmark-inc/bitmarkd/currency"
	"github.com/bitmark-inc/bitmarkd/fault"
	"github.com/bitmark-inc/bitmarkd/messagebus"
	"github.com/bitmark-inc/bitmarkd/mode"
	"github.com/bitmark-inc/bitmarkd/payment"
	"github.com/bitmark-inc/bitmarkd/reservoir"
	"github.com/bitmark-inc/bitmarkd/storage"
	"github.com/bitmark-inc/bitmarkd/transactionrecord"
	"github.com/bitmark-inc/bitmarkd/util"
	"github.com/bitmark-inc/bitmarkd/zmqutil"
	"github.com/bitmark-inc/logger"
	zmq "github.com/pebbe/zmq4"
)

const (
	subscriberSignal = "inproc://bitmark-subscriber-signal"
)

type subscriber struct {
	log          *logger.L
	push         *zmq.Socket
	pull         *zmq.Socket
	clients      []*zmqutil.Client
	dynamicStart int
}

// initialise the subscriber
func (sbsc *subscriber) initialise(privateKey []byte, publicKey []byte, subscribe []Connection, dynamicEnabled bool) error {

	log := logger.New("subscriber")
	if nil == log {
		return fault.ErrInvalidLoggerChannel
	}
	sbsc.log = log

	log.Info("initialising…")

	// validate connection count
	staticCount := len(subscribe) // can be zero
	if 0 == staticCount && !dynamicEnabled {
		log.Error("zero static connections and dynamic is disabled")
		return fault.ErrNoConnectionsAvailable
	}

	// signalling channel
	err := error(nil)
	sbsc.push, sbsc.pull, err = zmqutil.NewSignalPair(subscriberSignal)
	if nil != err {
		return err
	}

	// all sockets
	sbsc.clients = make([]*zmqutil.Client, staticCount+offsetCount)
	sbsc.dynamicStart = staticCount // index of first dynamic socket

	// error for goto fail
	errX := error(nil)

	// initially connect all static sockets
	for i, c := range subscribe {
		address, err := util.NewConnection(c.Address)
		if nil != err {
			log.Errorf("client[%d]=address: %q  error: %v", i, c.Address, err)
			errX = err
			goto fail
		}
		serverPublicKey, err := hex.DecodeString(c.PublicKey)
		if nil != err {
			log.Errorf("client[%d]=public: %q  error: %v", i, c.PublicKey, err)
			errX = err
			goto fail
		}

		// prevent connection to self
		if bytes.Equal(publicKey, serverPublicKey) {
			errX = fault.ErrConnectingToSelfForbidden
			log.Errorf("client[%d]=public: %q  error: %v", i, c.PublicKey, errX)
			goto fail
		}

		client, err := zmqutil.NewClient(zmq.SUB, privateKey, publicKey, 0)
		if nil != err {
			log.Errorf("client[%d]=%q  error: %v", i, c.Address, err)
			errX = err
			goto fail
		}

		sbsc.clients[i] = client

		err = client.Connect(address, serverPublicKey)
		if nil != err {
			log.Errorf("connect[%d]=%q  error: %v", i, c.Address, err)
			errX = err
			goto fail
		}
		log.Infof("public key: %x  at: %q", serverPublicKey, c.Address)
	}

	// just create sockets for dynamic clients
	for i := sbsc.dynamicStart; i < len(sbsc.clients); i += 1 {
		client, err := zmqutil.NewClient(zmq.SUB, privateKey, publicKey, 0)
		if nil != err {
			log.Errorf("client[%d]  error: %v", i, err)
			errX = err
			goto fail
		}

		sbsc.clients[i] = client
	}

	return nil

	// error handling
fail:
	zmqutil.CloseClients(sbsc.clients)
	return errX
}

// subscriber main loop
func (sbsc *subscriber) Run(args interface{}, shutdown <-chan struct{}) {

	log := sbsc.log

	log.Info("starting…")

	queue := messagebus.Bus.Subscriber.Chan()

	go func() {
		poller := zmq.NewPoller()
		for _, c := range sbsc.clients {
			c.Add(poller, zmq.POLLIN)
		}
		poller.Add(sbsc.pull, zmq.POLLIN)
	loop:
		for {
			log.Info("waiting…")
			sockets, _ := poller.Poll(-1)
			for _, socket := range sockets {
				switch s := socket.Socket; s {
				case sbsc.pull:
					data, err := s.RecvMessageBytes(0)
					if nil != err {
						log.Errorf("pull receive error: %v", err)
						break loop
					}

					switch string(data[0]) {
					case "connect":
						command := string(data[1])
						publicKey := data[2]
						broadcasts := data[3]
						connectTo(sbsc.log, sbsc.clients, sbsc.dynamicStart, command, publicKey, broadcasts)
					default:
						break loop
					}
				default:
					data, err := s.RecvMessageBytes(0)
					if nil != err {
						log.Errorf("receive error: %v", err)

					} else {
						sbsc.process(data)
					}
				}
			}
		}
		sbsc.pull.Close()
		zmqutil.CloseClients(sbsc.clients)
	}()

loop:
	for {
		log.Info("select…")

		select {
		// wait for shutdown
		case <-shutdown:
			break loop
		// wait for message
		case item := <-queue:
			sbsc.log.Infof("received: %s  public key: %x  connect: %x", item.Command, item.Parameters[0], item.Parameters[1])
			sbsc.push.SendMessage("connect", item.Command, item.Parameters[0], item.Parameters[1])
		}
	}

	sbsc.push.SendMessage("stop")
	sbsc.push.Close()
}

// process the received subscription
func (sbsc *subscriber) process(data [][]byte) {

	log := sbsc.log
	log.Info("incoming message")

	switch string(data[0]) {
	case "block":
		log.Infof("received block: %x", data[1])
		if !mode.Is(mode.Normal) {
			err := fault.ErrNotAvailableDuringSynchronise
			log.Warnf("failed assets: error: %v", err)
		} else {
			messagebus.Bus.Blockstore.Send("remote", data[1])
		}

	case "assets":
		log.Infof("received assets: %x", data[1])
		err := processAssets(data[1])
		if nil != err {
			log.Warnf("failed assets: error: %v", err)
		} else {
			messagebus.Bus.Broadcast.Send("assets", data[1])
		}

	case "issues":
		log.Infof("received issues: %x", data[1])
		err := processIssues(data[1])
		if nil != err {
			log.Warnf("failed issues: error: %v", err)
		} else {
			messagebus.Bus.Broadcast.Send("issues", data[1])
		}

	case "transfer":
		log.Infof("received transfer: %x", data[1])
		err := processTransfer(data[1])
		if nil != err {
			log.Warnf("failed transfer: error: %v", err)
		} else {
			messagebus.Bus.Broadcast.Send("transfer", data[1])
		}

	case "proof":
		log.Infof("received proof: %x", data[1])
		err := processProof(data[1])
		if nil != err {
			log.Warnf("failed proof: error: %v", err)
		} else {
			messagebus.Bus.Broadcast.Send("proof", data[1])
		}

	case "pay":
		log.Infof("received pay: %x", data[1])
		err := processPay(data[1])
		if nil != err {
			log.Warnf("failed pay: error: %v", err)
		} else {
			messagebus.Bus.Broadcast.Send("pay", data[1])
		}

	case "rpc":
		log.Infof("received rpc: fingerprint: %x  rpc: %x", data[1], data[2])
		if announce.AddRPC(data[1], data[2]) {
			messagebus.Bus.Broadcast.Send("rpc", data[1], data[2])
		}

	case "peer":
		log.Infof("received peer: %x  broadcast: %x  listener: %x", data[1], data[2], data[3])
		if announce.AddPeer(data[1], data[2], data[3]) {
			messagebus.Bus.Broadcast.Send("peer", data[1], data[2], data[3])
		}

	default:
		log.Warnf("received unhandled: %x", data)

	}
}

// un pack each asset and cache them
func processAssets(packed []byte) error {

	if !mode.Is(mode.Normal) {
		return fault.ErrNotAvailableDuringSynchronise
	}

	ok := false
	for {
		transaction, n, err := transactionrecord.Packed(packed).Unpack()
		if nil != err {
			return err
		}

		switch tx := transaction.(type) {
		case *transactionrecord.AssetData:
			_, packedAsset, err := asset.Cache(tx)
			if nil != err {
				return err
			}
			if nil != packedAsset {
				ok = true
			}

		default:
			return fault.ErrTransactionIsNotAnAsset
		}
		packed = packed[n:]
	}

	if !ok {
		return fault.ErrNoNewTransactions
	}
	return nil
}

// un pack each issue and cache them
func processIssues(packed []byte) error {

	if !mode.Is(mode.Normal) {
		return fault.ErrNotAvailableDuringSynchronise
	}

	for {
		transaction, n, err := transactionrecord.Packed(packed).Unpack()
		if nil != err {
			return err
		}

		switch tx := transaction.(type) {
		case *transactionrecord.BitmarkIssue:

			packedIssue, err := tx.Pack(tx.Owner)
			if nil != err {
				return err
			}

			if !asset.Exists(tx.AssetIndex) {
				return fault.ErrAssetNotFound
			}

			txId := packedIssue.MakeLink()
			key := txId[:]

			// even a single verified/confirmed issue fails the whole block
			if storage.Pool.Transactions.Has(key) || reservoir.Has(txId) {
				return fault.ErrTransactionAlreadyExists
			}

		default:
			return fault.ErrTransactionIsNotAnIssue
		}
		packed = packed[n:]
	}
	return nil
}

// un pack transfer and process it
func processTransfer(packed []byte) error {

	if !mode.Is(mode.Normal) {
		return fault.ErrNotAvailableDuringSynchronise
	}

	transaction, _, err := transactionrecord.Packed(packed).Unpack()
	if nil != err {
		return err
	}

	switch tx := transaction.(type) {
	case *transactionrecord.BitmarkTransfer:

		_, packedTransfer, _, _, err := block.VerifyTransfer(tx)
		if nil != err {
			return err
		}
		payment.Store(currency.Bitcoin, packedTransfer, 1, false)

	default:
		return fault.ErrTransactionIsNotATransfer
	}
	return nil
}

// process proof block
func processProof(packed []byte) error {

	if !mode.Is(mode.Normal) {
		return fault.ErrNotAvailableDuringSynchronise
	}

	var payId payment.PayId
	if len(packed) > payment.NonceLength+len(payId) {
		return fault.ErrInvalidNonce
	}

	copy(payId[:], packed[:len(payId)])
	nonce := packed[len(payId):]
	flag := payment.TryProof(payId, nonce)

	if !flag {
		// pay id already processed or was invalid
		return fault.ErrPayIdAlreadyUsed
	}

	return nil
}

// process pay block
func processPay(packed []byte) error {

	if !mode.Is(mode.Normal) {
		return fault.ErrNotAvailableDuringSynchronise
	}

	var payId payment.PayId
	if len(packed) > payment.ReceiptLength+len(payId) {
		return fault.ErrInvalidNonce
	}

	copy(payId[:], packed[:len(payId)])
	receipt := string(packed[len(payId):])

	flag := payment.TrackPayment(payId, receipt, payment.RequiredConfirmations)
	if !flag {
		// pay id already processed or was invalid
		return fault.ErrPayIdAlreadyUsed
	}

	return nil
}