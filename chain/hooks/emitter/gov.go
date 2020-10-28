package emitter

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/bandprotocol/bandchain/chain/hooks/common"
)

var (
	EventTypeInactiveProposal = types.EventTypeInactiveProposal
	EventTypeActiveProposal   = types.EventTypeActiveProposal
	StatusInactive            = 6
)

func (h *Hook) emitSetDeposit(ctx sdk.Context, txHash []byte, id uint64, depositor sdk.AccAddress) {
	deposit, _ := h.govKeeper.GetDeposit(ctx, id, depositor)
	h.Write("SET_DEPOSIT", common.JsDict{
		"proposal_id": id,
		"depositor":   depositor,
		"amount":      deposit.Amount.String(),
		"tx_hash":     txHash,
	})
}

func (h *Hook) emitUpdateProposalAfterDeposit(ctx sdk.Context, id uint64) {
	proposal, _ := h.govKeeper.GetProposal(ctx, id)
	h.Write("UPDATE_PROPOSAL", common.JsDict{
		"id":              id,
		"status":          int(proposal.Status),
		"total_deposit":   proposal.TotalDeposit.String(),
		"voting_time":     proposal.VotingStartTime.UnixNano(),
		"voting_end_time": proposal.VotingEndTime.UnixNano(),
	})
}

// handleMsgSubmitProposal implements emitter handler for MsgSubmitProposal.
func (app *Hook) handleMsgSubmitProposal(
	ctx sdk.Context, txHash []byte, msg gov.MsgSubmitProposal, evMap common.EvMap,
) {
	proposalId := uint64(common.Atoi(evMap[types.EventTypeSubmitProposal+"."+types.AttributeKeyProposalID][0]))
	proposal, _ := app.govKeeper.GetProposal(ctx, proposalId)
	app.Write("NEW_PROPOSAL", common.JsDict{
		"id":               proposalId,
		"proposer":         msg.Proposer,
		"type":             msg.Content.ProposalType(),
		"title":            msg.Content.GetTitle(),
		"description":      msg.Content.GetDescription(),
		"proposal_route":   msg.Content.ProposalRoute(),
		"status":           int(proposal.Status),
		"submit_time":      proposal.SubmitTime.UnixNano(),
		"deposit_end_time": proposal.DepositEndTime.UnixNano(),
		"total_deposit":    proposal.TotalDeposit.String(),
		"voting_time":      proposal.VotingStartTime.UnixNano(),
		"voting_end_time":  proposal.VotingEndTime.UnixNano(),
	})
	app.emitSetDeposit(ctx, txHash, proposalId, msg.Proposer)
}

// handleMsgDeposit implements emitter handler for MsgDeposit.
func (h *Hook) handleMsgDeposit(
	ctx sdk.Context, txHash []byte, msg gov.MsgDeposit,
) {
	h.emitSetDeposit(ctx, txHash, msg.ProposalID, msg.Depositor)
	h.emitUpdateProposalAfterDeposit(ctx, msg.ProposalID)
}

// handleMsgVote implements emitter handler for MsgVote.
func (h *Hook) handleMsgVote(
	txHash []byte, msg gov.MsgVote,
) {
	h.Write("SET_VOTE", common.JsDict{
		"proposal_id": msg.ProposalID,
		"voter":       msg.Voter,
		"answer":      int(msg.Option),
		"tx_hash":     txHash,
	})
}

func (h *Hook) handleEventInactiveProposal(evMap common.EvMap) {
	h.Write("UPDATE_PROPOSAL", common.JsDict{
		"id":     common.Atoi(evMap[types.EventTypeInactiveProposal+"."+types.AttributeKeyProposalID][0]),
		"status": StatusInactive,
	})
}

func (h *Hook) handleEventTypeActiveProposal(ctx sdk.Context, evMap common.EvMap) {
	id := uint64(common.Atoi(evMap[types.EventTypeActiveProposal+"."+types.AttributeKeyProposalID][0]))
	proposal, _ := h.govKeeper.GetProposal(ctx, id)
	h.Write("UPDATE_PROPOSAL", common.JsDict{
		"id":     id,
		"status": int(proposal.Status),
	})
}
