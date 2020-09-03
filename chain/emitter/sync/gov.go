package emitter

import (
	"github.com/bandprotocol/bandchain/chain/emitter/common"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/gov/types"
)

var (
	EventTypeInactiveProposal = types.EventTypeInactiveProposal
	EventTypeActiveProposal   = types.EventTypeActiveProposal
	StatusInactive            = 6
)

func (app *App) emitSetDeposit(txHash []byte, id uint64, depositor sdk.AccAddress) {
	deposit, _ := app.GovKeeper.GetDeposit(app.DeliverContext, id, depositor)
	app.Write("SET_DEPOSIT", common.JsDict{
		"proposal_id": id,
		"depositor":   depositor,
		"amount":      deposit.Amount.String(),
		"tx_hash":     txHash,
	})
}

func (app *App) emitUpdateProposalAfterDeposit(id uint64) {
	proposal, _ := app.GovKeeper.GetProposal(app.DeliverContext, id)
	app.Write("UPDATE_PROPOSAL", common.JsDict{
		"id":              id,
		"status":          int(proposal.Status),
		"total_deposit":   proposal.TotalDeposit.String(),
		"voting_time":     proposal.VotingStartTime.UnixNano(),
		"voting_end_time": proposal.VotingEndTime.UnixNano(),
	})
}

// handleMsgSubmitProposal implements emitter handler for MsgSubmitProposal.
func (app *App) handleMsgSubmitProposal(
	txHash []byte, msg gov.MsgSubmitProposal, evMap common.EvMap, extra common.JsDict,
) {
	proposalId := uint64(common.Atoi(evMap[types.EventTypeSubmitProposal+"."+types.AttributeKeyProposalID][0]))
	proposal, _ := app.GovKeeper.GetProposal(app.DeliverContext, proposalId)
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
	app.emitSetDeposit(txHash, proposalId, msg.Proposer)
}

// handleMsgDeposit implements emitter handler for MsgDeposit.
func (app *App) handleMsgDeposit(
	txHash []byte, msg gov.MsgDeposit, evMap common.EvMap, extra common.JsDict,
) {
	app.emitSetDeposit(txHash, msg.ProposalID, msg.Depositor)
	app.emitUpdateProposalAfterDeposit(msg.ProposalID)
}

// handleMsgVote implements emitter handler for MsgVote.
func (app *App) handleMsgVote(
	txHash []byte, msg gov.MsgVote, evMap common.EvMap, extra common.JsDict,
) {
	app.Write("SET_VOTE", common.JsDict{
		"proposal_id": msg.ProposalID,
		"voter":       msg.Voter,
		"answer":      int(msg.Option),
		"tx_hash":     txHash,
	})
}

func (app *App) handleEventInactiveProposal(evMap common.EvMap) {
	app.Write("UPDATE_PROPOSAL", common.JsDict{
		"id":     common.Atoi(evMap[types.EventTypeInactiveProposal+"."+types.AttributeKeyProposalID][0]),
		"status": StatusInactive,
	})
}

func (app *App) handleEventTypeActiveProposal(evMap common.EvMap) {
	id := uint64(common.Atoi(evMap[types.EventTypeActiveProposal+"."+types.AttributeKeyProposalID][0]))
	proposal, _ := app.GovKeeper.GetProposal(app.DeliverContext, id)
	app.Write("UPDATE_PROPOSAL", common.JsDict{
		"id":     id,
		"status": int(proposal.Status),
	})
}
