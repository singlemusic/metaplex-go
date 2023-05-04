// This is a very simple program designed to allow metadata tagging to a given mint,
// with an update authority that can change that metadata going forward.
// Read more: https://github.com/metaplex-foundation/metaplex/tree/master/rust/token-metadata/program
// Code generated by https://github.com/gagliardetto/anchor-go. DO NOT EDIT.

package token_metadata

import (
	"bytes"
	"fmt"
	ag_spew "github.com/davecgh/go-spew/spew"
	ag_binary "github.com/gagliardetto/binary"
	ag_solanago "github.com/gagliardetto/solana-go"
	ag_text "github.com/gagliardetto/solana-go/text"
	ag_treeout "github.com/gagliardetto/treeout"
)

var ProgramID ag_solanago.PublicKey = ag_solanago.MustPublicKeyFromBase58("metaqbxxUerdq28cj1RbAWkYQm3ybzjb6a8bt518x1s")

func SetProgramID(pubkey ag_solanago.PublicKey) {
	ProgramID = pubkey
	ag_solanago.RegisterInstructionDecoder(ProgramID, registryDecodeInstruction)
}

const ProgramName = "TokenMetadata"

func init() {
	if !ProgramID.IsZero() {
		ag_solanago.RegisterInstructionDecoder(ProgramID, registryDecodeInstruction)
	}
}

const (
	// Create Metadata object.
	Instruction_CreateMetadataAccount uint8 = iota

	// Update a Metadata
	Instruction_UpdateMetadataAccount

	// Register a Metadata as a Master Edition V1, which means Editions can be minted.
	// Henceforth, no further tokens will be mintable from this primary mint. Will throw an error if more than one
	// token exists, and will throw an error if less than one token exists in this primary mint.
	Instruction_DeprecatedCreateMasterEdition

	// Given an authority token minted by the Printing mint of a master edition, and a brand new non-metadata-ed mint with one token
	// make a new Metadata + Edition that is a child of the master edition denoted by this authority token.
	Instruction_DeprecatedMintNewEditionFromMasterEditionViaPrintingToken

	// Allows updating the primary sale boolean on Metadata solely through owning an account
	// containing a token from the metadata's mint and being a signer on this transaction.
	// A sort of limited authority for limited update capability that is required for things like
	// Metaplex to work without needing full authority passing.
	Instruction_UpdatePrimarySaleHappenedViaToken

	// Reserve up to 200 editions in sequence for up to 200 addresses in an existing reservation PDA, which can then be used later by
	// redeemers who have printing tokens as a reservation to get a specific edition number
	// as opposed to whatever one is currently listed on the master edition. Used by Auction Manager
	// to guarantee printing order on bid redemption. AM will call whenever the first person redeems a
	// printing bid to reserve the whole block
	// of winners in order and then each winner when they get their token submits their mint and account
	// with the pda that was created by that first bidder - the token metadata can then cross reference
	// these people with the list and see that bidder A gets edition #2, so on and so forth.
	//
	// NOTE: If you have more than 20 addresses in a reservation list, this may be called multiple times to build up the list,
	// otherwise, it simply wont fit in one transaction. Only provide a total_reservation argument on the first call, which will
	// allocate the edition space, and in follow up calls this will specifically be unnecessary (and indeed will error.)
	Instruction_DeprecatedSetReservationList

	// Create an empty reservation list for a resource who can come back later as a signer and fill the reservation list
	// with reservations to ensure that people who come to get editions get the number they expect. See SetReservationList for more.
	Instruction_DeprecatedCreateReservationList

	// Sign a piece of metadata that has you as an unverified creator so that it is now verified.
	Instruction_SignMetadata

	// Using a one time authorization token from a master edition v1, print any number of printing tokens from the printing_mint
	// one time, burning the one time authorization token.
	Instruction_DeprecatedMintPrintingTokensViaToken

	// Using your update authority, mint printing tokens for your master edition.
	Instruction_DeprecatedMintPrintingTokens

	// Register a Metadata as a Master Edition V2, which means Edition V2s can be minted.
	// Henceforth, no further tokens will be mintable from this primary mint. Will throw an error if more than one
	// token exists, and will throw an error if less than one token exists in this primary mint.
	Instruction_CreateMasterEdition

	// Given a token account containing the master edition token to prove authority, and a brand new non-metadata-ed mint with one token
	// make a new Metadata + Edition that is a child of the master edition denoted by this authority token.
	Instruction_MintNewEditionFromMasterEditionViaToken

	// Converts the Master Edition V1 to a Master Edition V2, draining lamports from the two printing mints
	// to the owner of the token account holding the master edition token. Permissionless.
	// Can only be called if there are currenly no printing tokens or one time authorization tokens in circulation.
	Instruction_ConvertMasterEditionV1ToV2

	// Proxy Call to Mint Edition using a Store Token Account as a Vault Authority.
	Instruction_MintNewEditionFromMasterEditionViaVaultProxy

	// Puff a Metadata - make all of it's variable length fields (name/uri/symbol) a fixed length using a null character
	// so that it can be found using offset searches by the RPC to make client lookups cheaper.
	Instruction_PuffMetadata

	// Update a Metadata with is_mutable as a parameter
	Instruction_UpdateMetadataAccountV2

	// Create Metadata object.
	Instruction_CreateMetadataAccountV2

	// Register a Metadata as a Master Edition V2, which means Edition V2s can be minted.
	// Henceforth, no further tokens will be mintable from this primary mint. Will throw an error if more than one
	// token exists, and will throw an error if less than one token exists in this primary mint.
	Instruction_CreateMasterEditionV3

	// If a MetadataAccount Has a Collection allow the UpdateAuthority of the Collection to Verify the NFT Belongs in the Collection
	Instruction_VerifyCollection

	// Utilize or Use an NFT , burns the NFT and returns the lamports to the update authority if the use method is burn and its out of uses.
	// Use Authority can be the Holder of the NFT, or a Delegated Use Authority.
	Instruction_Utilize

	// Approve another account to call `utilize` on this NFT
	Instruction_ApproveUseAuthority

	// Revoke account to call `utilize` on this NFT
	Instruction_RevokeUseAuthority

	// If a MetadataAccount Has a Collection allow an Authority of the Collection to unverify an NFT in a Collection
	Instruction_UnverifyCollection

	// Approve another account to verify nfts beloging to a collection, [verify_collection] on the collection NFT
	Instruction_ApproveCollectionAuthority

	// Revoke account to call [verify_collection] on this NFT
	Instruction_RevokeCollectionAuthority

	Instruction_SetAndVerifyCollection

	Instruction_FreezeDelegatedAccount

	Instruction_ThawDelegatedAccount

	Instruction_RemoveCreatorVerification

	Instruction_BurnNft

	Instruction_VerifySizedCollectionItem

	Instruction_UnverifySizedCollectionItem

	Instruction_SetAndVerifySizedCollectionItem

	Instruction_CreateMetadataAccountV3
)

// InstructionIDToName returns the name of the instruction given its ID.
func InstructionIDToName(id uint8) string {
	switch id {
	case Instruction_CreateMetadataAccount:
		return "CreateMetadataAccount"
	case Instruction_UpdateMetadataAccount:
		return "UpdateMetadataAccount"
	case Instruction_DeprecatedCreateMasterEdition:
		return "DeprecatedCreateMasterEdition"
	case Instruction_DeprecatedMintNewEditionFromMasterEditionViaPrintingToken:
		return "DeprecatedMintNewEditionFromMasterEditionViaPrintingToken"
	case Instruction_UpdatePrimarySaleHappenedViaToken:
		return "UpdatePrimarySaleHappenedViaToken"
	case Instruction_DeprecatedSetReservationList:
		return "DeprecatedSetReservationList"
	case Instruction_DeprecatedCreateReservationList:
		return "DeprecatedCreateReservationList"
	case Instruction_SignMetadata:
		return "SignMetadata"
	case Instruction_DeprecatedMintPrintingTokensViaToken:
		return "DeprecatedMintPrintingTokensViaToken"
	case Instruction_DeprecatedMintPrintingTokens:
		return "DeprecatedMintPrintingTokens"
	case Instruction_CreateMasterEdition:
		return "CreateMasterEdition"
	case Instruction_MintNewEditionFromMasterEditionViaToken:
		return "MintNewEditionFromMasterEditionViaToken"
	case Instruction_ConvertMasterEditionV1ToV2:
		return "ConvertMasterEditionV1ToV2"
	case Instruction_MintNewEditionFromMasterEditionViaVaultProxy:
		return "MintNewEditionFromMasterEditionViaVaultProxy"
	case Instruction_PuffMetadata:
		return "PuffMetadata"
	case Instruction_UpdateMetadataAccountV2:
		return "UpdateMetadataAccountV2"
	case Instruction_CreateMetadataAccountV2:
		return "CreateMetadataAccountV2"
	case Instruction_CreateMetadataAccountV3:
		return "CreateMetadataAccountV3"
	case Instruction_CreateMasterEditionV3:
		return "CreateMasterEditionV3"
	case Instruction_VerifyCollection:
		return "VerifyCollection"
	case Instruction_Utilize:
		return "Utilize"
	case Instruction_ApproveUseAuthority:
		return "ApproveUseAuthority"
	case Instruction_RevokeUseAuthority:
		return "RevokeUseAuthority"
	case Instruction_UnverifyCollection:
		return "UnverifyCollection"
	case Instruction_ApproveCollectionAuthority:
		return "ApproveCollectionAuthority"
	case Instruction_RevokeCollectionAuthority:
		return "RevokeCollectionAuthority"
	default:
		return ""
	}
}

type Instruction struct {
	ag_binary.BaseVariant
}

func (inst *Instruction) EncodeToTree(parent ag_treeout.Branches) {
	if enToTree, ok := inst.Impl.(ag_text.EncodableToTree); ok {
		enToTree.EncodeToTree(parent)
	} else {
		parent.Child(ag_spew.Sdump(inst))
	}
}

var InstructionImplDef = ag_binary.NewVariantDefinition(
	ag_binary.Uint8TypeIDEncoding,
	[]ag_binary.VariantType{
		{
			"CreateMetadataAccount", (*CreateMetadataAccount)(nil),
		},
		{
			"UpdateMetadataAccount", (*UpdateMetadataAccount)(nil),
		},
		{
			"DeprecatedCreateMasterEdition", (*DeprecatedCreateMasterEdition)(nil),
		},
		{
			"DeprecatedMintNewEditionFromMasterEditionViaPrintingToken", (*DeprecatedMintNewEditionFromMasterEditionViaPrintingToken)(nil),
		},
		{
			"UpdatePrimarySaleHappenedViaToken", (*UpdatePrimarySaleHappenedViaToken)(nil),
		},
		{
			"DeprecatedSetReservationList", (*DeprecatedSetReservationList)(nil),
		},
		{
			"DeprecatedCreateReservationList", (*DeprecatedCreateReservationList)(nil),
		},
		{
			"SignMetadata", (*SignMetadata)(nil),
		},
		{
			"DeprecatedMintPrintingTokensViaToken", (*DeprecatedMintPrintingTokensViaToken)(nil),
		},
		{
			"DeprecatedMintPrintingTokens", (*DeprecatedMintPrintingTokens)(nil),
		},
		{
			"CreateMasterEdition", (*CreateMasterEdition)(nil),
		},
		{
			"MintNewEditionFromMasterEditionViaToken", (*MintNewEditionFromMasterEditionViaToken)(nil),
		},
		{
			"ConvertMasterEditionV1ToV2", (*ConvertMasterEditionV1ToV2)(nil),
		},
		{
			"MintNewEditionFromMasterEditionViaVaultProxy", (*MintNewEditionFromMasterEditionViaVaultProxy)(nil),
		},
		{
			"PuffMetadata", (*PuffMetadata)(nil),
		},
		{
			"UpdateMetadataAccountV2", (*UpdateMetadataAccountV2)(nil),
		},
		{
			"CreateMetadataAccountV2", (*CreateMetadataAccountV2)(nil),
		},
		{
			"CreateMetadataAccountV3", (*CreateMetadataAccountV3)(nil),
		},
		{
			"CreateMasterEditionV3", (*CreateMasterEditionV3)(nil),
		},
		{
			"VerifyCollection", (*VerifyCollection)(nil),
		},
		{
			"Utilize", (*Utilize)(nil),
		},
		{
			"ApproveUseAuthority", (*ApproveUseAuthority)(nil),
		},
		{
			"RevokeUseAuthority", (*RevokeUseAuthority)(nil),
		},
		{
			"UnverifyCollection", (*UnverifyCollection)(nil),
		},
		{
			"ApproveCollectionAuthority", (*ApproveCollectionAuthority)(nil),
		},
		{
			"RevokeCollectionAuthority", (*RevokeCollectionAuthority)(nil),
		},
	},
)

func (inst *Instruction) ProgramID() ag_solanago.PublicKey {
	return ProgramID
}

func (inst *Instruction) Accounts() (out []*ag_solanago.AccountMeta) {
	return inst.Impl.(ag_solanago.AccountsGettable).GetAccounts()
}

func (inst *Instruction) Data() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := ag_binary.NewBorshEncoder(buf).Encode(inst); err != nil {
		return nil, fmt.Errorf("unable to encode instruction: %w", err)
	}
	return buf.Bytes(), nil
}

func (inst *Instruction) TextEncode(encoder *ag_text.Encoder, option *ag_text.Option) error {
	return encoder.Encode(inst.Impl, option)
}

func (inst *Instruction) UnmarshalWithDecoder(decoder *ag_binary.Decoder) error {
	return inst.BaseVariant.UnmarshalBinaryVariant(decoder, InstructionImplDef)
}

func (inst *Instruction) MarshalWithEncoder(encoder *ag_binary.Encoder) error {
	err := encoder.WriteUint8(inst.TypeID.Uint8())
	if err != nil {
		return fmt.Errorf("unable to write variant type: %w", err)
	}
	return encoder.Encode(inst.Impl)
}

func registryDecodeInstruction(accounts []*ag_solanago.AccountMeta, data []byte) (interface{}, error) {
	inst, err := DecodeInstruction(accounts, data)
	if err != nil {
		return nil, err
	}
	return inst, nil
}

func DecodeInstruction(accounts []*ag_solanago.AccountMeta, data []byte) (*Instruction, error) {
	inst := new(Instruction)
	if err := ag_binary.NewBorshDecoder(data).Decode(inst); err != nil {
		return nil, fmt.Errorf("unable to decode instruction: %w", err)
	}
	if v, ok := inst.Impl.(ag_solanago.AccountsSettable); ok {
		err := v.SetAccounts(accounts)
		if err != nil {
			return nil, fmt.Errorf("unable to set accounts for instruction: %w", err)
		}
	}
	return inst, nil
}
