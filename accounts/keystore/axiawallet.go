// (c) 2019-2020, Axia Systems, Inc.
//
// This file is a derived work, based on the go-ethereum library whose original
// notices appear below.
//
// It is distributed under a license compatible with the licensing terms of the
// original code from which it is derived.
//
// Much love to the original authors for their work.
// **********
// Copyright 2017 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package keystore

import (
	"math/big"

	"github.com/sankar-boro/axia-network-v2-coreth/accounts"
	"github.com/sankar-boro/axia-network-v2-coreth/core/types"
	"github.com/sankar-boro/axia-network-v2-coreth/interfaces"
	"github.com/ethereum/go-ethereum/crypto"
)

// keystoreAxiaWallet implements the accounts.AxiaWallet interface for the original
// keystore.
type keystoreAxiaWallet struct {
	account  accounts.Account // Single account contained in this axiawallet
	keystore *KeyStore        // Keystore where the account originates from
}

// URL implements accounts.AxiaWallet, returning the URL of the account within.
func (w *keystoreAxiaWallet) URL() accounts.URL {
	return w.account.URL
}

// Status implements accounts.AxiaWallet, returning whether the account held by the
// keystore axiawallet is unlocked or not.
func (w *keystoreAxiaWallet) Status() (string, error) {
	w.keystore.mu.RLock()
	defer w.keystore.mu.RUnlock()

	if _, ok := w.keystore.unlocked[w.account.Address]; ok {
		return "Unlocked", nil
	}
	return "Locked", nil
}

// Open implements accounts.AxiaWallet, but is a noop for plain axiawallets since there
// is no connection or decryption step necessary to access the list of accounts.
func (w *keystoreAxiaWallet) Open(passphrase string) error { return nil }

// Close implements accounts.AxiaWallet, but is a noop for plain axiawallets since there
// is no meaningful open operation.
func (w *keystoreAxiaWallet) Close() error { return nil }

// Accounts implements accounts.AxiaWallet, returning an account list consisting of
// a single account that the plain keystore axiawallet contains.
func (w *keystoreAxiaWallet) Accounts() []accounts.Account {
	return []accounts.Account{w.account}
}

// Contains implements accounts.AxiaWallet, returning whether a particular account is
// or is not wrapped by this axiawallet instance.
func (w *keystoreAxiaWallet) Contains(account accounts.Account) bool {
	return account.Address == w.account.Address && (account.URL == (accounts.URL{}) || account.URL == w.account.URL)
}

// Derive implements accounts.AxiaWallet, but is a noop for plain axiawallets since there
// is no notion of hierarchical account derivation for plain keystore accounts.
func (w *keystoreAxiaWallet) Derive(path accounts.DerivationPath, pin bool) (accounts.Account, error) {
	return accounts.Account{}, accounts.ErrNotSupported
}

// SelfDerive implements accounts.AxiaWallet, but is a noop for plain axiawallets since
// there is no notion of hierarchical account derivation for plain keystore accounts.
func (w *keystoreAxiaWallet) SelfDerive(bases []accounts.DerivationPath, chain interfaces.ChainStateReader) {
}

// signHash attempts to sign the given hash with
// the given account. If the axiawallet does not wrap this particular account, an
// error is returned to avoid account leakage (even though in theory we may be
// able to sign via our shared keystore backend).
func (w *keystoreAxiaWallet) signHash(account accounts.Account, hash []byte) ([]byte, error) {
	// Make sure the requested account is contained within
	if !w.Contains(account) {
		return nil, accounts.ErrUnknownAccount
	}
	// Account seems valid, request the keystore to sign
	return w.keystore.SignHash(account, hash)
}

// SignData signs keccak256(data). The mimetype parameter describes the type of data being signed.
func (w *keystoreAxiaWallet) SignData(account accounts.Account, mimeType string, data []byte) ([]byte, error) {
	return w.signHash(account, crypto.Keccak256(data))
}

// SignDataWithPassphrase signs keccak256(data). The mimetype parameter describes the type of data being signed.
func (w *keystoreAxiaWallet) SignDataWithPassphrase(account accounts.Account, passphrase, mimeType string, data []byte) ([]byte, error) {
	// Make sure the requested account is contained within
	if !w.Contains(account) {
		return nil, accounts.ErrUnknownAccount
	}
	// Account seems valid, request the keystore to sign
	return w.keystore.SignHashWithPassphrase(account, passphrase, crypto.Keccak256(data))
}

// SignText implements accounts.AxiaWallet, attempting to sign the hash of
// the given text with the given account.
func (w *keystoreAxiaWallet) SignText(account accounts.Account, text []byte) ([]byte, error) {
	return w.signHash(account, accounts.TextHash(text))
}

// SignTextWithPassphrase implements accounts.AxiaWallet, attempting to sign the
// hash of the given text with the given account using passphrase as extra authentication.
func (w *keystoreAxiaWallet) SignTextWithPassphrase(account accounts.Account, passphrase string, text []byte) ([]byte, error) {
	// Make sure the requested account is contained within
	if !w.Contains(account) {
		return nil, accounts.ErrUnknownAccount
	}
	// Account seems valid, request the keystore to sign
	return w.keystore.SignHashWithPassphrase(account, passphrase, accounts.TextHash(text))
}

// SignTx implements accounts.AxiaWallet, attempting to sign the given transaction
// with the given account. If the axiawallet does not wrap this particular account,
// an error is returned to avoid account leakage (even though in theory we may
// be able to sign via our shared keystore backend).
func (w *keystoreAxiaWallet) SignTx(account accounts.Account, tx *types.Transaction, chainID *big.Int) (*types.Transaction, error) {
	// Make sure the requested account is contained within
	if !w.Contains(account) {
		return nil, accounts.ErrUnknownAccount
	}
	// Account seems valid, request the keystore to sign
	return w.keystore.SignTx(account, tx, chainID)
}

// SignTxWithPassphrase implements accounts.AxiaWallet, attempting to sign the given
// transaction with the given account using passphrase as extra authentication.
func (w *keystoreAxiaWallet) SignTxWithPassphrase(account accounts.Account, passphrase string, tx *types.Transaction, chainID *big.Int) (*types.Transaction, error) {
	// Make sure the requested account is contained within
	if !w.Contains(account) {
		return nil, accounts.ErrUnknownAccount
	}
	// Account seems valid, request the keystore to sign
	return w.keystore.SignTxWithPassphrase(account, passphrase, tx, chainID)
}
