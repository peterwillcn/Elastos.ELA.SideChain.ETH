package vm

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"

	"github.com/btcsuite/btcutil/base58"

	"github.com/elastos/Elastos.ELA.SideChain.ETH/common"
	"github.com/elastos/Elastos.ELA.SideChain.ETH/core/vm/did"
	"github.com/elastos/Elastos.ELA.SideChain.ETH/core/vm/did/base64url"
	"github.com/elastos/Elastos.ELA.SideChain.ETH/log"
	"github.com/elastos/Elastos.ELA.SideChain.ETH/params"
)

// PrecompiledContract is the basic interface for native Go contracts. The implementation
// requires a deterministic gas count based on the input size of the Run method of the
// contract.
type PrecompiledContractDID interface {
	RequiredGas(evm *EVM, input []byte) uint64              // RequiredPrice calculates the contract gas use
	Run(evm *EVM, input []byte, gas uint64) ([]byte, error) // Run runs the precompiled contract
}

var PrecompileContractsDID = map[common.Address]PrecompiledContractDID{
	common.BytesToAddress([]byte{22}): &operationDID{},
}

// RunPrecompiledContract runs and evaluates the output of a precompiled contract.
func RunPrecompiledContractDID(evm *EVM, p PrecompiledContractDID, input []byte, contract *Contract) (ret []byte, err error) {
	gas := p.RequiredGas(evm, input)
	if contract.UseGas(gas) {
		return p.Run(evm, input, contract.Gas)
	}
	log.Error("run did contract out of gas")
	return nil, ErrOutOfGas
}

type operationDID struct{}

func (j *operationDID) RequiredGas(evm *EVM, input []byte) uint64 {
	return params.DIDBaseGasCost
}

func checkPayloadSyntax(p *did.DIDPayload) error {
	// check proof
	if p.Proof.VerificationMethod == "" {
		return errors.New("proof Creator is nil")
	}
	if p.Proof.Signature == "" {
		return errors.New("proof Created is nil")
	}
	if p.DIDDoc != nil {
		if len(p.DIDDoc.Authentication) == 0 {
			return errors.New("did doc Authentication is nil")
		}
		if p.DIDDoc.Expires == "" {
			return errors.New("did doc Expires is nil")
		}

		for _, pkInfo := range p.DIDDoc.PublicKey {
			if pkInfo.ID == "" {
				return errors.New("check DIDDoc.PublicKey ID is empty")
			}
			// if  pkInfo.Type != "ECDSAsecp256r1" {
			// 	return errors.New("check DIDDoc.PublicKey Type != ECDSAsecp256r1")
			// }
			// if pkInfo.Controller == "" {
			// 	return errors.New("check DIDDoc.PublicKey Controller is empty")
			// }
			if pkInfo.PublicKeyBase58 == "" {
				return errors.New("check DIDDoc.PublicKey PublicKeyBase58 is empty")
			}
		}
	}
	return nil
}

func (j *operationDID) Run(evm *EVM, input []byte, gas uint64) ([]byte, error) {
	//block height from context BlockNumber. config height address from config

	configHeight := evm.chainConfig.OldDIDMigrateHeight
	configAddr := evm.chainConfig.OldDIDMigrateAddr
	senderAddr := evm.Context.Origin.String()
	log.Info("####", "configAddr", configAddr, "senderAddr", senderAddr)

	//BlockNumber <= configHeight senderAddr must be configAddr
	if evm.Context.BlockNumber.Cmp(configHeight) <= 0 {
		if senderAddr != configAddr {
			log.Info("#### BlockNumber.Cmp(configHeight) <= 0 or callerAddress.String() != configAddr")
			return false32Byte, errors.New("Befor configHeight only configAddr can send DID tx")
		}
	} else {
		if senderAddr == configAddr {
			log.Info("#### BlockNumber.Cmp(configHeight) > 0 callerAddress.String() should not configAddr")
			return false32Byte, errors.New("after configHeight  configAddr can not send migrate DID tx")
		}
	}

	data := getData(input, 32, uint64(len(input))-32)
	p := new(did.DIDPayload)
	if err := json.Unmarshal(data, p); err != nil {
		log.Error("DIDPayload input is error", "input", string(data))
		return false32Byte, err
	}
	switch p.Header.Operation {
	case did.Create_DID_Operation, did.Update_DID_Operation:
		payloadBase64, _ := base64url.DecodeString(p.Payload)
		payloadInfo := new(did.DIDDoc)
		if err := json.Unmarshal(payloadBase64, payloadInfo); err != nil {
			return false32Byte, errors.New("createDIDVerify Payload is error")
		}
		p.DIDDoc = payloadInfo
		// abnormal payload check
		if err := checkPayloadSyntax(p); err != nil {
			log.Error("checkPayloadSyntax error", "error", err, "ID", p.DIDDoc.ID)
			return false32Byte, err
		}
		var err error
		isRegisterDID := isDID(p.DIDDoc)
		if isRegisterDID {
			if err = checkRegisterDID(evm, p, gas); err != nil {
				log.Error("checkRegisterDID error", "error", err, "ID", p.DIDDoc.ID)
			}
		} else {
			if err = checkCustomizedDID(evm, p, gas); err != nil {
				log.Error("checkCustomizedDID error", "error", err, "ID", p.DIDDoc.ID)
			}
		}
		if err != nil {
			return false32Byte, err
		}
		buf := new(bytes.Buffer)
		p.Serialize(buf, did.DIDVersion)
		evm.StateDB.AddDIDLog(p.DIDDoc.ID, p.Header.Operation, buf.Bytes())
	case did.Transfer_DID_Operation:
		payloadBase64, _ := base64url.DecodeString(p.Payload)
		payloadInfo := new(did.DIDDoc)
		if err := json.Unmarshal(payloadBase64, payloadInfo); err != nil {
			return false32Byte, errors.New("createDIDVerify Payload is error")
		}
		p.DIDDoc = payloadInfo
		if err := checkCustomizedDID(evm, p, gas); err != nil {
			log.Error("checkCustomizedDID error", "error", err)
			return false32Byte, err
		}
		buf := new(bytes.Buffer)
		p.Serialize(buf, did.DIDVersion)
		evm.StateDB.AddDIDLog(p.DIDDoc.ID, p.Header.Operation, buf.Bytes())
	case did.Deactivate_DID_Operation:
		if err := checkDeactivateDID(evm, p); err != nil {
			log.Error("checkDeactivateDID error", "error", err)
			return false32Byte, err
		}
		id := p.Payload
		buf := new(bytes.Buffer)
		p.Serialize(buf, did.DIDVersion)
		evm.StateDB.AddDIDLog(id, p.Header.Operation, buf.Bytes())
	case did.Declare_Verifiable_Credential_Operation, did.Revoke_Verifiable_Credential_Operation:
		if err := checkVerifiableCredential(evm, p); err != nil {
			log.Error("checkVerifiableCredential error", "error", err)
			return false32Byte, err
		}
		buf := new(bytes.Buffer)
		p.Serialize(buf, did.DIDVersion)
		evm.StateDB.AddDIDLog(p.CredentialDoc.ID, p.Header.Operation, buf.Bytes())
	default:
		log.Error("error operation", "operation", p.Header.Operation)
		return false32Byte, errors.New("error operation:" + p.Header.Operation)
	}
	return true32Byte, nil
}

func isDID(didDoc *did.DIDDoc) bool {

	if !strings.HasPrefix(didDoc.ID, did.DID_ELASTOS_PREFIX) {
		return false
	}
	idString := did.GetDIDFromUri(didDoc.ID)

	for _, pkInfo := range didDoc.PublicKey {
		publicKey := base58.Decode(pkInfo.PublicKeyBase58)
		if did.IsMatched(publicKey, idString) {
			return true
		}
	}
	return false
}
