import asyncio
from typing import Any, Dict

from src.config import Config
from src.shared import DataPipes

class Manager:
    def __init__(self, config: Config, data_pipes: DataPipes):
        self.config = config
        self.data_pipes = data_pipes

    async def run(self, ctx: Dict[str, Any]):
        """The main loop for the manager."""
        print("Manager is running.")

        # Create tasks to listen on each queue
        request_task = asyncio.create_task(self._listen_for_requests())
        bit_vote_task = asyncio.create_task(self._listen_for_bit_votes())
        policy_task = asyncio.create_task(self._listen_for_policies())

        try:
            # Wait for all listener tasks to complete.
            # In a real scenario, these would run indefinitely until cancelled.
            await asyncio.gather(request_task, bit_vote_task, policy_task)
        except asyncio.CancelledError:
            print("Manager task cancelled.")
        finally:
            # Cancel the listener tasks when the manager is stopped
            request_task.cancel()
            bit_vote_task.cancel()
            policy_task.cancel()
            print("Manager stopped.")

    async def _listen_for_requests(self):
        while True:
            request = await self.data_pipes.requests.get()
            print(f"Manager received request: {request}")
            # Process the request...
            self.data_pipes.requests.task_done()

    async def _listen_for_bit_votes(self):
        while True:
            bit_vote = await self.data_pipes.bit_votes.get()
            print(f"Manager received bit vote: {bit_vote}")
            # Process the bit vote...
            self.data_pipes.bit_votes.task_done()

    async def _listen_for_policies(self):
        while True:
            policy = await self.data_pipes.signing_policies.get()
            print(f"Manager received signing policy: {policy}")
            # Process the signing policy...
            self.data_pipes.signing_policies.task_done()
