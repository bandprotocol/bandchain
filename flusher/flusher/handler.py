import base64 as b64
from datetime import datetime
from sqlalchemy import select
from sqlalchemy.dialects.postgresql import insert

from .db import (
    blocks,
    transactions,
    accounts,
    data_sources,
    oracle_scripts,
    requests,
    raw_requests,
    val_requests,
    reports,
    raw_reports,
    validators,
    delegations,
    validator_votes,
    unbonding_delegations,
    redelegations,
    account_transactions,
    proposals,
    deposits,
    votes,
    historical_bonded_token_on_validators,
)


class Handler(object):
    def __init__(self, conn):
        self.conn = conn

    def get_transaction_id(self, tx_hash):
        return self.conn.execute(
            select([transactions.c.id]).where(transactions.c.hash == tx_hash)
        ).scalar()

    def get_validator_id(self, val):
        return self.conn.execute(
            select([validators.c.id]).where(validators.c.operator_address == val)
        ).scalar()

    def get_account_id(self, address):
        return self.conn.execute(
            select([accounts.c.id]).where(accounts.c.address == address)
        ).scalar()

    def handle_new_block(self, msg):
        self.conn.execute(blocks.insert(), msg)

    def handle_new_transaction(self, msg):
        related_tx_accounts = msg["related_accounts"]
        del msg["related_accounts"]
        self.conn.execute(transactions.insert(), msg)
        tx_id = self.get_transaction_id(msg["hash"])
        for account in related_tx_accounts:
            self.conn.execute(
                account_transactions.insert(),
                {"transaction_id": tx_id, "account_id": self.get_account_id(account)},
            )

    def handle_set_account(self, msg):
        if self.get_account_id(msg["address"]) is None:
            self.conn.execute(accounts.insert(), msg)
        else:
            condition = True
            for col in accounts.primary_key.columns.values():
                condition = (col == msg[col.name]) & condition
            self.conn.execute(accounts.update().where(condition).values(**msg))

    def handle_set_data_source(self, msg):
        if msg["tx_hash"] is not None:
            msg["transaction_id"] = self.get_transaction_id(msg["tx_hash"])
        else:
            msg["transaction_id"] = None
        del msg["tx_hash"]
        self.conn.execute(
            insert(data_sources)
            .values(**msg)
            .on_conflict_do_update(constraint="data_sources_pkey", set_=msg)
        )

    def handle_set_oracle_script(self, msg):
        if msg["tx_hash"] is not None:
            msg["transaction_id"] = self.get_transaction_id(msg["tx_hash"])
        else:
            msg["transaction_id"] = None
        del msg["tx_hash"]
        self.conn.execute(
            insert(oracle_scripts)
            .values(**msg)
            .on_conflict_do_update(constraint="oracle_scripts_pkey", set_=msg)
        )

    def handle_new_request(self, msg):
        msg["transaction_id"] = self.get_transaction_id(msg["tx_hash"])
        del msg["tx_hash"]
        self.handle_set_request_count_per_days({"date": msg["timestamp"]})
        del msg["timestamp"]
        self.conn.execute(requests.insert(), msg)
        self.handle_set_oracle_script_request({"oracle_script_id": msg["oracle_script_id"]})

    def handle_update_request(self, msg):
        condition = True
        for col in requests.primary_key.columns.values():
            condition = (col == msg[col.name]) & condition
        self.conn.execute(requests.update().where(condition).values(**msg))

    def handle_new_val_request(self, msg):
        msg["validator_id"] = self.get_validator_id(msg["validator"])
        del msg["validator"]
        self.conn.execute(val_requests.insert(), msg)

    def handle_new_report(self, msg):
        msg["transaction_id"] = self.get_transaction_id(msg["tx_hash"])
        del msg["tx_hash"]
        msg["validator_id"] = self.get_validator_id(msg["validator"])
        del msg["validator"]
        msg["reporter_id"] = self.get_account_id(msg["reporter"])
        del msg["reporter"]
        self.conn.execute(reports.insert(), msg)

    def handle_new_raw_report(self, msg):
        msg["validator_id"] = self.get_validator_id(msg["validator"])
        del msg["validator"]
        self.conn.execute(raw_reports.insert(), msg)

    def handle_set_validator(self, msg):
        last_update = msg["last_update"]
        del msg["last_update"]
        msg["account_id"] = self.get_account_id(msg["delegator_address"])
        del msg["delegator_address"]
        if self.get_validator_id(msg["operator_address"]) is None:
            self.conn.execute(validators.insert(), msg)
        else:
            condition = True
            for col in validators.primary_key.columns.values():
                condition = (col == msg[col.name]) & condition
            self.conn.execute(validators.update().where(condition).values(**msg))
        self.handle_new_historical_bonded_token_on_validator(
            {
                "validator_id": self.get_validator_id(msg["operator_address"]),
                "bonded_tokens": msg["tokens"],
                "timestamp": last_update,
            }
        )

    def handle_update_validator(self, msg):
        if "tokens" in msg and "last_update" in msg:
            self.handle_new_historical_bonded_token_on_validator(
                {
                    "validator_id": self.get_validator_id(msg["operator_address"]),
                    "bonded_tokens": msg["tokens"],
                    "timestamp": msg["last_update"],
                }
            )
            del msg["last_update"]
        self.conn.execute(
            validators.update()
            .where(validators.c.operator_address == msg["operator_address"])
            .values(**msg)
        )

    def handle_set_delegation(self, msg):
        msg["delegator_id"] = self.get_account_id(msg["delegator_address"])
        del msg["delegator_address"]
        msg["validator_id"] = self.get_validator_id(msg["operator_address"])
        del msg["operator_address"]
        self.conn.execute(
            insert(delegations)
            .values(**msg)
            .on_conflict_do_update(constraint="delegations_pkey", set_=msg)
        )

    def handle_update_delegation(self, msg):
        msg["delegator_id"] = self.get_account_id(msg["delegator_address"])
        del msg["delegator_address"]
        msg["validator_id"] = self.get_validator_id(msg["operator_address"])
        del msg["operator_address"]
        condition = True
        for col in delegations.primary_key.columns.values():
            condition = (col == msg[col.name]) & condition
        self.conn.execute(delegations.update().where(condition).values(**msg))

    def handle_remove_delegation(self, msg):
        msg["delegator_id"] = self.get_account_id(msg["delegator_address"])
        del msg["delegator_address"]
        msg["validator_id"] = self.get_validator_id(msg["operator_address"])
        del msg["operator_address"]
        condition = True
        for col in delegations.primary_key.columns.values():
            condition = (col == msg[col.name]) & condition
        self.conn.execute(delegations.delete().where(condition))

    def handle_new_validator_vote(self, msg):
        self.conn.execute(
            insert(validator_votes)
            .values(**msg)
            .on_conflict_do_update(constraint="validator_votes_pkey", set_=msg)
        )

    def handle_new_unbonding_delegation(self, msg):
        msg["delegator_id"] = self.get_account_id(msg["delegator_address"])
        del msg["delegator_address"]
        msg["validator_id"] = self.get_validator_id(msg["operator_address"])
        del msg["operator_address"]
        self.conn.execute(insert(unbonding_delegations).values(**msg))

    def handle_new_redelegation(self, msg):
        msg["delegator_id"] = self.get_account_id(msg["delegator_address"])
        del msg["delegator_address"]
        msg["validator_src_id"] = self.get_validator_id(msg["operator_src_address"])
        del msg["operator_src_address"]
        msg["validator_dst_id"] = self.get_validator_id(msg["operator_dst_address"])
        del msg["operator_dst_address"]
        self.conn.execute(insert(redelegations).values(**msg))

    def handle_new_proposal(self, msg):
        msg["proposer_id"] = self.get_account_id(msg["proposer"])
        del msg["proposer"]
        self.conn.execute(proposals.insert(), msg)

    def handle_set_deposit(self, msg):
        msg["depositor_id"] = self.get_account_id(msg["depositor"])
        del msg["depositor"]
        msg["tx_id"] = self.get_transaction_id(msg["tx_hash"])
        del msg["tx_hash"]
        self.conn.execute(
            insert(deposits)
            .values(**msg)
            .on_conflict_do_update(constraint="deposits_pkey", set_=msg)
        )

    def handle_set_vote(self, msg):
        msg["voter_id"] = self.get_account_id(msg["voter"])
        del msg["voter"]
        msg["tx_id"] = self.get_transaction_id(msg["tx_hash"])
        del msg["tx_hash"]
        self.conn.execute(
            insert(votes).values(**msg).on_conflict_do_update(constraint="votes_pkey", set_=msg)
        )

    def handle_update_proposal(self, msg):
        condition = True
        for col in proposals.primary_key.columns.values():
            condition = (col == msg[col.name]) & condition
        self.conn.execute(proposals.update().where(condition).values(**msg))

    def handle_new_historical_bonded_token_on_validator(self, msg):
        self.conn.execute(
            insert(historical_bonded_token_on_validators)
            .values(**msg)
            .on_conflict_do_update(
                constraint="historical_bonded_token_on_validators_pkey", set_=msg
            )
        )
