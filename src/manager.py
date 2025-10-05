import asyncio
from typing import Any, Dict

from src.config import Config
from src.shared import DataPipes
from src.logger import log

class Manager:
    def __init__(self, config: Config, data_pipes: DataPipes):
        self.config = config
        self.data_pipes = data_pipes
        log.info("Manager initialized.")

    async def run(self, ctx: Dict[str, Any]):
        """The main loop for the manager."""
        log.info("Manager is running.")

        # Create tasks to listen on each queue
        log.debug("Creating listener tasks...")
        request_task = asyncio.create_task(self._listen_for_requests())
        bit_vote_task = asyncio.create_task(self._listen_for_bit_votes())
        policy_task = asyncio.create_task(self._listen_for_policies())
        log.debug("Listener tasks created.")

        try:
            # Wait for all listener tasks to complete.
            # In a real scenario, these would run indefinitely until cancelled.
            await asyncio.gather(request_task, bit_vote_task, policy_task)
        except asyncio.CancelledError:
            log.info("Manager task cancelled.")
        finally:
            # Cancel the listener tasks when the manager is stopped
            log.debug("Cancelling listener tasks...")
            request_task.cancel()
            bit_vote_task.cancel()
            policy_task.cancel()
            log.info("Manager stopped.")

    async def _listen_for_requests(self):
        while True:
            request = await self.data_pipes.requests.get()
            log.info(f"Manager received request: {request}")
            log.debug(f"Processing request: {request}")
            # Process the request...
            self.data_pipes.requests.task_done()

    async def _listen_for_bit_votes(self):
        while True:
            bit_vote = await self.data_pipes.bit_votes.get()
            log.info(f"Manager received bit vote: {bit_vote}")
            log.debug(f"Processing bit vote: {bit_vote}")
            # Process the bit vote...
            self.data_pipes.bit_votes.task_done()

    async def _listen_for_policies(self):
        while True:
            policy = await self.data_pipes.signing_policies.get()
            log.info(f"Manager received signing policy: {policy}")
            log.debug(f"Processing signing policy: {policy}")
            # Process the signing policy...
            self.data_pipes.signing_policies.task_done()
