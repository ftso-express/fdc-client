import asyncio
from typing import Any, Dict

from src.config import Config
from src.database import DatabaseConnector
from src.shared import DataPipes
from src.logger import log

class Collector:
    def __init__(self, config: Config, data_pipes: DataPipes):
        self.config = config
        self.data_pipes = data_pipes
        self.db_connector = None
        self.db_connection = None
        log.info("Collector initialized.")

    async def connect_to_db(self):
        """Establishes a connection to the database."""
        log.info("Connecting to the database...")
        db_config = self.config.database_config
        self.db_connector = DatabaseConnector(db_config)
        try:
            self.db_connection = self.db_connector.connect()
            log.info("Database connection established.")
        except Exception as e:
            log.error(f"Error connecting to database: {e}")
            raise

    async def run(self, ctx: Dict[str, Any]):
        """The main loop for the collector."""
        await self.connect_to_db()

        # The original Go application has a WaitForDBToSync function.
        # This is a placeholder for that logic.
        log.info("Waiting for DB to sync...")
        await asyncio.sleep(2) # Placeholder
        log.info("DB is in sync.")

        while True:
            try:
                # In a real implementation, this is where you would query
                # for attestation requests, signing policies, and bitVotes.
                log.debug("Collector is running, fetching data...")

                # Simulate putting data into the queues
                log.debug("Putting attestation request into queue...")
                await self.data_pipes.requests.put("attestation_request_from_collector")
                log.debug("Putting bit vote into queue...")
                await self.data_pipes.bit_votes.put("bit_vote_from_collector")
                log.debug("Putting signing policy into queue...")
                await self.data_pipes.signing_policies.put("signing_policy_from_collector")
                log.debug("Data successfully put into queues.")

                await asyncio.sleep(10) # Fetch data every 10 seconds
            except asyncio.CancelledError:
                log.info("Collector task cancelled.")
                if self.db_connector:
                    self.db_connector.close()
                    log.info("Database connection closed.")
                break
