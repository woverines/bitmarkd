// Copyright (c) 2014-2016 Bitmark Inc.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package payment

import (
	// "github.com/bitmark-inc/bitmarkd/datastore"
	"github.com/bitmark-inc/bitmarkd/fault"
	"github.com/bitmark-inc/bitmarkd/transactionrecord"
)

// verifier loop
func (state *verifierData) Run(args interface{}, shutdown <-chan struct{}) {

	log := state.log

loop:
	for {
		log.Info("waiting…")
		select {
		case <-shutdown:
			break loop
		case transactions := <-state.queue:
			log.Infof("received: transactions: %x", transactions)
			state.setVerified(transactions)
		}
	}
}

// store all transactions in disk storage to await confirmation
func (state *verifierData) setVerified(transactions []byte) {

	// split transactions
	records := transactionrecord.Packed(transactions)
	for len(records) > 0 {

		// consistency check
		transaction, length, err := records.Unpack()
		fault.PanicIfError("setVerified", err) // memory buffer was corrupted, hardware problem or invalid write?

		// first item
		packed := records[:length]
		txId := packed.MakeLink()

		state.log.Infof("unpacked: %v", transaction)
		state.log.Infof("packed txid:%x data: %x", txId, packed)

		//datastore.WriteVerified(packed) // ***** FIX THIS: need to commit to disk

		// remaining items
		records = records[length:]
	}
}
