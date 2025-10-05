import asyncio
from typing import Any, Dict

from src.config import Config
from src.metrics import MANAGER_ITEMS_PROCESSED, QUEUE_SIZE
from src.shared import DataPipes
from src.logger import log


class Manager:
    def __init__(self, config: Config, data_pipes: DataPipes):
        self.config = config
        self.data_pipes = data_pipes
        log.info("Manager initialized.")
        self.metrics_enabled = self.config.metrics_config.get("enabled", False)

    async def _update_queue_metrics(self):
        """Periodically updates the queue size gauges."""
        while True:
            if self.metrics_enabled:
                QUEUE_SIZE.labels(queue_name="requests").set(self.data_pipes.requests.qsize())
                QUEUE_SIZE.labels(queue_name="bit_votes").set(self.data_pipes.bit_votes.qsize())
                QUEUE_SIZE.labels(queue_name="signing_policies").set(self.data_pipes.signing_policies.qsize())
            await asyncio.sleep(5)  # Update every 5 seconds

    async def run(self, ctx: Dict[str, Any]):
        """The main loop for the manager."""
        log.info("Manager is running.")

        # Create tasks to listen on each queue
        log.debug("Creating listener tasks...")
        request_task = asyncio.create_task(self._listen_for_requests())
        bit_vote_task = asyncio.create_task(self._listen_for_bit_votes())
        policy_task = asyncio.create_task(self._listen_for_policies())
        log.debug("Listener tasks created.")
        metrics_task = asyncio.create_task(self._update_queue_metrics())

        try:
            # Wait for all listener tasks to complete.
            # In a real scenario, these would run indefinitely until cancelled.
            await asyncio.gather(request_task, bit_vote_task, policy_task, metrics_task)
        except asyncio.CancelledError:
            log.info("Manager task cancelled.")
        finally:
            # Cancel the listener tasks when the manager is stopped
            log.debug("Cancelling listener tasks...")
            request_task.cancel()
            bit_vote_task.cancel()
            policy_task.cancel()
            log.info("Manager stopped.")
            metrics_task.cancel()
            print("Manager stopped.")

    async def _listen_for_requests(self):
        while True:
            request = await self.data_pipes.requests.get()
            log.info(f"Manager received request: {request}")
            log.debug(f"Processing request: {request}")
            # Process the request...
            if self.metrics_enabled:
                MANAGER_ITEMS_PROCESSED.labels(item_type="request").inc()
            self.data_pipes.requests.task_done()

    async def _listen_for_bit_votes(self):
        while True:
            bit_vote = await self.data_pipes.bit_votes.get()
            log.info(f"Manager received bit vote: {bit_vote}")
            log.debug(f"Processing bit vote: {bit_vote}")
            # Process the bit vote...
            if self.metrics_enabled:
                MANAGER_ITEMS_PROCESSED.labels(item_type="bit_vote").inc()
            self.data_pipes.bit_votes.task_done()

    async def _listen_for_policies(self):
        while True:
            policy = await self.data_pipes.signing_policies.get()
            log.info(f"Manager received signing policy: {policy}")
            log.debug(f"Processing signing policy: {policy}")
            # Process the signing policy...
            if self.metrics_enabled:
                MANAGER_ITEMS_PROCESSED.labels(item_type="policy").inc()
            self.data_pipes.signing_policies.task_done()
