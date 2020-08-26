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

func (app *App) emitGovModule() {
	app.GovKeeper.IterateProposals(app.DeliverContext, func(proposal types.Proposal) (stop bool) {
		app.emitNewProposal(proposal, nil)
		return false
	})
	app.GovKeeper.IterateAllDeposits(app.DeliverContext, func(deposit types.Deposit) (stop bool) {
		app.emitSetDeposit(nil, deposit.ProposalID, deposit.Depositor)
		return false
	})
	app.GovKeeper.IterateAllVotes(app.DeliverContext, func(vote types.Vote) (stop bool) {
		app.emitSetVote(nil, vote)
		return false
	})
}

func (app *App) emitNewProposal(proposal gov.Proposal, proposer sdk.AccAddress) {
	app.Write("NEW_PROPOSAL", JsDict{
		"id":               proposal.ProposalID,
		"proposer":         proposer,
		"type":             proposal.Content.ProposalType(),
		"title":            proposal.Content.GetTitle(),
		"description":      proposal.Content.GetDescription(),
		"proposal_route":   proposal.Content.ProposalRoute(),
		"status":           int(proposal.Status),
		"submit_time":      proposal.SubmitTime.UnixNano(),
		"deposit_end_time": proposal.DepositEndTime.UnixNano(),
		"total_deposit":    proposal.TotalDeposit.String(),
		"voting_time":      proposal.VotingStartTime.UnixNano(),
		"voting_end_time":  proposal.VotingEndTime.UnixNano(),
	})
}

func (app *App) emitSetDeposit(txHash []byte, id uint64, depositor sdk.AccAddress) {
	deposit, _ := app.GovKeeper.GetDeposit(app.DeliverContext, id, depositor)
	app.Write("SET_DEPOSIT", JsDict{
		"proposal_id": id,
		"depositor":   depositor,
		"amount":      deposit.Amount.String(),
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

func (app *App) emitSetVote(txHash []byte, vote types.Vote) {
	app.Write("SET_VOTE", JsDict{
		"proposal_id": vote.ProposalID,
		"voter":       vote.Voter,
		"answer":      int(vote.Option),
		"tx_hash":     txHash,
	})
}

// handleMsgSubmitProposal implements emitter handler for MsgSubmitProposal.
func (app *App) handleMsgSubmitProposal(
	txHash []byte, msg gov.MsgSubmitProposal, evMap EvMap, extra JsDict,
) {
	proposalId := uint64(atoi(evMap[types.EventTypeSubmitProposal+"."+types.AttributeKeyProposalID][0]))
	proposal, _ := app.GovKeeper.GetProposal(app.DeliverContext, proposalId)
	app.emitNewProposal(proposal, msg.Proposer)
	app.emitSetDeposit(txHash, proposalId, msg.Proposer)
}

// handleMsgDeposit implements emitter handler for MsgDeposit.
func (app *App) handleMsgDeposit(
	txHash []byte, msg gov.MsgDeposit, evMap EvMap, extra JsDict,
) {
	app.emitSetDeposit(txHash, msg.ProposalID, msg.Depositor)
	app.emitUpdateProposalAfterDeposit(msg.ProposalID)
}

// handleMsgVote implements emitter handler for MsgVote.
func (app *App) handleMsgVote(
	txHash []byte, msg gov.MsgVote, evMap EvMap, extra JsDict,
) {
	vote, _ := app.GovKeeper.GetVote(app.DeliverContext, msg.ProposalID, msg.Voter)
	app.emitSetVote(txHash, vote)
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
