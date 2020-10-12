package emitter

import (
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
	app.Write("SET_DEPOSIT", JsDict{
		"proposal_id": id,
		"depositor":   depositor,
		"amount":      deposit.Amount.String(),
		"tx_hash":     txHash,
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
	app.emitSetDeposit(txHash, proposalId, msg.Proposer)
	extra["proposal_id"] = proposalId
}

// handleMsgDeposit implements emitter handler for MsgDeposit.
func (app *App) handleMsgDeposit(
	txHash []byte, msg gov.MsgDeposit, evMap EvMap, extra JsDict,
) {
	app.emitSetDeposit(txHash, msg.ProposalID, msg.Depositor)
	proposal, _ := app.GovKeeper.GetProposal(app.DeliverContext, msg.ProposalID)
	app.Write("UPDATE_PROPOSAL", JsDict{
		"id":              msg.ProposalID,
		"status":          int(proposal.Status),
		"total_deposit":   proposal.TotalDeposit.String(),
		"voting_time":     proposal.VotingStartTime.UnixNano(),
		"voting_end_time": proposal.VotingEndTime.UnixNano(),
	})
	extra["title"] = proposal.GetTitle()
}

// handleMsgVote implements emitter handler for MsgVote.
func (app *App) handleMsgVote(
	txHash []byte, msg gov.MsgVote, evMap EvMap, extra JsDict,
) {
	app.Write("SET_VOTE", JsDict{
		"proposal_id": msg.ProposalID,
		"voter":       msg.Voter,
		"answer":      int(msg.Option),
		"tx_hash":     txHash,
	})
	proposal, _ := app.GovKeeper.GetProposal(app.DeliverContext, msg.ProposalID)
	extra["title"] = proposal.GetTitle()
}

func (app *App) handleEventInactiveProposal(evMap EvMap) {
	app.Write("UPDATE_PROPOSAL", JsDict{
		"id":     atoi(evMap[types.EventTypeInactiveProposal+"."+types.AttributeKeyProposalID][0]),
		"status": StatusInactive,
	})
}

func (app *App) handleEventTypeActiveProposal(evMap EvMap) {
	id := uint64(atoi(evMap[types.EventTypeActiveProposal+"."+types.AttributeKeyProposalID][0]))
	proposal, _ := app.GovKeeper.GetProposal(app.DeliverContext, id)
	app.Write("UPDATE_PROPOSAL", JsDict{
		"id":     id,
		"status": int(proposal.Status),
	})
}
