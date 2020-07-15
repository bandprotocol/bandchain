package emitter

import (
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/gov/types"
)

// handleMsgSubmitProposal implements emitter handler for MsgSubmitProposal.
func (app *App) handleMsgSubmitProposal(
	txHash []byte, msg gov.MsgSubmitProposal, evMap EvMap, extra JsDict,
) {
	proposalId := uint64(atoi(evMap[types.EventTypeSubmitProposal+"."+types.AttributeKeyProposalID][0]))
	proposal, _ := app.GovKeeper.GetProposal(app.DeliverContext, proposalId)
	app.Write("SET_PROPOSAL", JsDict{
		"proposal_id":      proposalId,
		"detail":           string(app.Codec().MustMarshalJSON(proposal)),
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
	// TODO: add deposit for preposer
}
