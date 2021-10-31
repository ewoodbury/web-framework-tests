import time
import os
import requests
import asyncio
import asyncpg
from dotenv import load_dotenv


load_dotenv()


class PostgresHandler:
    def __init__(self):
        self.conn = None

    async def connect(self):
        pg_user = os.getenv("PG_USER")
        pg_pass = os.getenv("PG_PASS")
        self.conn = await asyncpg.connect(
            f"postgres://{pg_user}:{pg_pass}@localhost:5432/postgres")

    async def create_signals_tables(self):
        create_statement = """
        CREATE TABLE IF NOT EXISTS cell_signals (
            id SERIAL PRIMARY KEY,
            test_id INTEGER NOT NULL,
            measured_at TIMESTAMP NOT NULL,
            cell_voltage DOUBLE PRECISION NOT NULL,
            cell_current DOUBLE PRECISION NOT NULL
        );
        """
        await self.conn.execute(create_statement)
        return

    async def reset_signals_table(self):
        reset_statement = """
        truncate cell_signals;
        delete from cell_signals;
        """
        await self.conn.execute(reset_statement)
        return


def test_local_memory(n_requests=1000):
    t0 = time.time()
    for i in range(n_requests):
        requests.post("http://127.0.0.1:8000/data/local",
                      data={"voltage": 3.4 + i*.001, "current": 1.1 * i*.001})
    t1 = time.time()
    return 1000 * (t1 - t0) / n_requests


def test_db_insert_small(n_requests=1000):
    t0 = time.time()
    for i in range(n_requests):
        requests.post("http://127.0.0.1:8000/data/db-small",
                      data={"voltage": 3.4 + i*.001, "current": 1.1 * i*.001})
    t1 = time.time()
    return 1000 * (t1 - t0) / n_requests


async def run():
    db = PostgresHandler()
    await db.connect()
    await db.create_signals_tables()
    await db.reset_signals_table()
    ms_per_request = {}
    ms_per_request["local_memory"] = test_local_memory()
    ms_per_request["small_db_insert"] = test_db_insert_small()
    print(ms_per_request)


asyncio.run(run())
