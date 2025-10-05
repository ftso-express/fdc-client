import asyncio
import mysql.connector
from typing import Any, Dict

from src.config import Config
from src.shared import DataPipes

class Collector:
    def __init__(self, config: Config, data_pipes: DataPipes):
        self.config = config
        self.data_pipes = data_pipes
        self.db_connection = None

    async def connect_to_db(self):
        """Establishes a connection to the database."""
        db_config = self.config.database_config
        try:
            self.db_connection = mysql.connector.connect(
                host=db_config["host"],
                port=db_config["port"],
                user=db_config["username"],
                password=db_config["password"],
                database=db_config["database"],
            )
            print("Successfully connected to the database.")
        except mysql.connector.Error as err:
            print(f"Error connecting to database: {err}")
            # In a real application, you'd want more robust error handling
            # and retry logic here.
            raise

    async def run(self, ctx: Dict[str, Any]):
        """The main loop for the collector."""
        await self.connect_to_db()

        # The original Go application has a WaitForDBToSync function.
        # This is a placeholder for that logic.
        print("Waiting for DB to sync...")
        await asyncio.sleep(2) # Placeholder
        print("DB is in sync.")

        while True:
            try:
                # In a real implementation, this is where you would query
                # for attestation requests, signing policies, and bitVotes.
                print("Collector is running, fetching data...")

                # Simulate putting data into the queues
                await self.data_pipes.requests.put("attestation_request_from_collector")
                await self.data_pipes.bit_votes.put("bit_vote_from_collector")
                await self.data_pipes.signing_policies.put("signing_policy_from_collector")

                await asyncio.sleep(10) # Fetch data every 10 seconds
            except asyncio.CancelledError:
                print("Collector task cancelled.")
                if self.db_connection and self.db_connection.is_connected():
                    self.db_connection.close()
                    print("Database connection closed.")
                break
