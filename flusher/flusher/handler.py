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
    account_transcations,
)


class Handler(object):
    def __init__(self, conn):
        self.conn = conn

    def handle_new_block(self, msg):
        self.conn.execute(blocks.insert(), msg)

    def handle_new_transaction(self, msg):
        related_tx_accounts = msg["related_accounts"]
        del msg["related_accounts"]
        res = self.conn.execute(transactions.insert(), msg)
        tx_id = res.inserted_primary_key[0]
        for account in related_tx_accounts:
            self.conn.execute(
                account_transcations.insert(), {"transaction_id": tx_id, "address": account}
            )

    def handle_set_account(self, msg):
        self.conn.execute(
            insert(accounts)
            .values(**msg)
            .on_conflict_do_update(constraint="accounts_pkey", set_=msg)
        )

    def handle_set_data_source(self, msg):
        if "tx_hash" in msg:
            msg["transaction_id"] = self.conn.execute(
                select([transactions.c.id]).where(transactions.c.hash == msg["tx_hash"])
            ).scalar()
        else:
            msg["transaction_id"] = None
        del msg["tx_hash"]
        self.conn.execute(
            insert(data_sources)
            .values(**msg)
            .on_conflict_do_update(constraint="data_sources_pkey", set_=msg)
        )

    def handle_set_oracle_script(self, msg):
        if "tx_hash" in msg:
            msg["transaction_id"] = self.conn.execute(
                select([transactions.c.id]).where(transactions.c.hash == msg["tx_hash"])
            ).scalar()
        else:
            msg["transaction_id"] = None
        del msg["tx_hash"]
        self.conn.execute(
            insert(oracle_scripts)
            .values(**msg)
            .on_conflict_do_update(constraint="oracle_scripts_pkey", set_=msg)
        )

    def handle_new_request(self, msg):
        msg["transaction_id"] = self.conn.execute(
            select([transactions.c.id]).where(transactions.c.hash == msg["tx_hash"])
        ).scalar()
        del msg["tx_hash"]
        self.conn.execute(requests.insert(), msg)

    def handle_update_request(self, msg):
        condition = True
        for col in requests.primary_key.columns.values():
            condition = (col == msg[col.name]) & condition
        self.conn.execute(requests.update().where(condition).values(**msg))

    def handle_new_raw_request(self, msg):
        self.conn.execute(raw_requests.insert(), msg)

    def handle_new_val_request(self, msg):
        self.conn.execute(val_requests.insert(), msg)

    def handle_new_report(self, msg):
        msg["transaction_id"] = self.conn.execute(
            select([transactions.c.id]).where(transactions.c.hash == msg["tx_hash"])
        ).scalar()
        del msg["tx_hash"]
        self.conn.execute(reports.insert(), msg)

    def handle_new_raw_report(self, msg):
        self.conn.execute(raw_reports.insert(), msg)

    def handle_set_validator(self, msg):
        self.conn.execute(
            insert(validators)
            .values(**msg)
            .on_conflict_do_update(constraint="validators_pkey", set_=msg)
        )

    def handle_update_validator(self, msg):
        condition = True
        for col in validators.primary_key.columns.values():
            condition = (col == msg[col.name]) & condition
        self.conn.execute(validators.update().where(condition).values(**msg))

    def handle_set_delegation(self, msg):
        self.conn.execute(
            insert(delegations)
            .values(**msg)
            .on_conflict_do_update(constraint="delegations_pkey", set_=msg)
        )

    def handle_remove_delegation(self, msg):
        condition = True
        for col in delegations.primary_key.columns.values():
            condition = (col == msg[col.name]) & condition
        self.conn.execute(delegations.delete().where(condition))

    def handle_new_validator_vote(self, msg):
        self.conn.execute(insert(validator_votes).values(**msg))

    def handle_new_unbonding_delegation(self, msg):
        self.conn.execute(insert(unbonding_delegations).values(**msg))

    def handle_new_redelegation(self, msg):
        self.conn.execute(insert(redelegations).values(**msg))
