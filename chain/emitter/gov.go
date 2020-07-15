package emitter

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/gov/types"
)

func (app *App) emitSetDeposit(txHash []byte, id uint64, depositor sdk.AccAddress, amount sdk.Coins) {
	app.Write("SET_DEPOSIT", JsDict{
		"proposal_id": id,
		"depositor":   depositor,
		"amount":      amount.String(),
		"tx_hash":     txHash,
	})
}

func (app *App) emitUpdateProposalAfterDeposit(id uint64) {
	proposal, _ := app.GovKeeper.GetProposal(app.DeliverContext, id)
	app.Write("UPDATE_PROPOSAL", JsDict{
		"id":              id,
		"status":          int(proposal.Status),
		"total_deposit":   proposal.TotalDeposit.String(),
		"voting_time":     proposal.VotingStartTime.UnixNano(),
		"voting_end_time": proposal.VotingEndTime.UnixNano(),
	})
}

// handleMsgSubmitProposal implements emitter handler for MsgSubmitProposal.
func (app *App) handleMsgSubmitProposal(
	txHash []byte, msg gov.MsgSubmitProposal, evMap EvMap, extra JsDict,
) {
	proposalId := uint64(atoi(evMap[types.EventTypeSubmitProposal+"."+types.AttributeKeyProposalID][0]))
	proposal, _ := app.GovKeeper.GetProposal(app.DeliverContext, proposalId)
	app.Write("NEW_PROPOSAL", JsDict{
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
	app.emitSetDeposit(txHash, proposalId, msg.Proposer, msg.InitialDeposit)
}

// handleMsgDeposit implements emitter handler for MsgDeposit.
func (app *App) handleMsgDeposit(
	txHash []byte, msg gov.MsgDeposit, evMap EvMap, extra JsDict,
) {
	app.emitSetDeposit(txHash, msg.ProposalID, msg.Depositor, msg.Amount)
	app.emitUpdateProposalAfterDeposit(msg.ProposalID)
}
