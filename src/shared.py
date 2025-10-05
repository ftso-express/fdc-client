import asyncio

class DataPipes:
    """
    A class to hold the asyncio queues that act as data pipes between
    the collector, manager, and server components.
    """
    def __init__(self):
        self.requests = asyncio.Queue()
        self.bit_votes = asyncio.Queue()
        self.signing_policies = asyncio.Queue()
        # The rounds data is accessed by the server, so we need a way to share it.
        # A simple list or dict protected by a lock could work, or a queue
        # if the access pattern is producer/consumer. Based on the Go code,
        # it seems to be a shared data structure rather than a queue.
        # For now, I will represent it as a simple list.
        self.rounds = []

if __name__ == '__main__':
    # Example usage:
    async def main():
        pipes = DataPipes()

        # Simulate putting some data into the queues
        await pipes.requests.put("attestation_request_1")
        await pipes.bit_votes.put("bit_vote_for_round_123")
        await pipes.signing_policies.put("signing_policy_for_epoch_5")
        pipes.rounds.append("round_data_1")

        # Simulate getting data from the queues
        request = await pipes.requests.get()
        bit_vote = await pipes.bit_votes.get()
        policy = await pipes.signing_policies.get()
        round_data = pipes.rounds.pop(0)

        print(f"Got request: {request}")
        print(f"Got bit vote: {bit_vote}")
        print(f"Got signing policy: {policy}")
        print(f"Got round data: {round_data}")

        assert pipes.requests.empty()
        assert pipes.bit_votes.empty()
        assert pipes.signing_policies.empty()
        assert not pipes.rounds

        print("\nShared data pipes work as expected.")

    asyncio.run(main())