
import click

from .cli import cli

from sqlalchemy import create_engine

@cli.command()
@click.argument("password")
@click.option(
    "--db",
    help="Database URI connection string.",
    default="localhost:5432/postgres",
    show_default=True,
)
def hasura(password,db):
    engine = create_engine("postgresql+psycopg2://" + db, echo=True)

    engine.execute('''DROP ROLE IF EXISTS hasura;''')
    engine.execute(f'''CREATE USER hasura WITH PASSWORD '{password}';''')
    engine.execute('''CREATE EXTENSION IF NOT EXISTS pgcrypto;''')
    engine.execute('''CREATE SCHEMA IF NOT EXISTS hdb_catalog;''')
    engine.execute('''CREATE SCHEMA IF NOT EXISTS hdb_views;''')


    # make the user an owner of system schemas
    engine.execute('''ALTER SCHEMA hdb_catalog OWNER TO hasura;''')
    engine.execute('''ALTER SCHEMA hdb_views OWNER TO hasura;''')

    # # grant select permissions on information_schema and pg_catalog. This is
    # required for hasura to query list of available tables
    engine.execute('''GRANT SELECT ON ALL TABLES IN SCHEMA information_schema TO hasura;''')
    engine.execute('''GRANT SELECT ON ALL TABLES IN SCHEMA pg_catalog TO hasura;''')

    # Below permissions are optional. This is dependent on what access to your
    # tables/schemas - you want give to hasura. If you want expose the public
    # schema for GraphQL query then give permissions on public schema to the
    # hasura user.
    # Be careful to use these in your production db. Consult the postgres manual or
    # your DBA and give appropriate permissions.

    # grant all privileges on all tables in the public schema. This can be customised:
    # For example, if you only want to use GraphQL regular queries and not mutations,
    # then you can set: GRANT SELECT ON ALL TABLES...
    engine.execute('''GRANT USAGE ON SCHEMA public TO hasura;''')
    engine.execute('''GRANT SELECT ON ALL TABLES IN SCHEMA public TO hasura;''')
    engine.execute('''GRANT ALL ON ALL SEQUENCES IN SCHEMA public TO hasura;''')
    engine.execute('''GRANT ALL ON ALL FUNCTIONS IN SCHEMA public TO hasura;''')

    engine.execute('''GRANT SELECT ON information_schema.table_constraints, information_schema.key_column_usage, information_schema.columns, information_schema.views, information_schema.schemata, information_schema.routines TO hasura;''')
    engine.execute('''GRANT SELECT ON pg_catalog.pg_constraint, pg_catalog.pg_class, pg_catalog.pg_namespace, pg_catalog.pg_attribute, pg_catalog.pg_proc, pg_catalog.pg_available_extensions, pg_catalog.pg_statio_all_tables, pg_catalog.pg_description TO hasura;''')
    engine.execute('''ALTER ROLE hasura SET statement_timeout="5s"''')
