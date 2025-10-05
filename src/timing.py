from typing import Dict, Any

class Timing:
    """
    A class to handle timing calculations for rounds and epochs, based on the
    logic from the Go implementation.
    """
    # Default values from the Go implementation
    DEFAULT_T0 = 1658429955
    DEFAULT_T0_REWARD_DELAY = 0
    DEFAULT_REWARD_EPOCH_LENGTH = 240
    DEFAULT_COLLECT_DURATION_SEC = 90
    DEFAULT_CHOOSE_DURATION_SEC = 45

    def __init__(self, timing_config: Dict[str, Any]):
        self.t0: int = timing_config.get("t0", self.DEFAULT_T0)
        self.t0_reward_delay: int = timing_config.get("t0_reward_delay", self.DEFAULT_T0_REWARD_DELAY)
        self.reward_epoch_length: int = timing_config.get("reward_epoch_length", self.DEFAULT_REWARD_EPOCH_LENGTH)
        self.collect_duration_sec: int = timing_config.get("collect_duration_sec", self.DEFAULT_COLLECT_DURATION_SEC)
        self.choose_duration_sec: int = timing_config.get("choose_duration_sec", self.DEFAULT_CHOOSE_DURATION_SEC)

    def round_id_for_timestamp(self, t: int) -> int:
        """Calculates the round ID that is active at a given timestamp."""
        if t < self.t0:
            raise ValueError(f"Timestamp {t} is before the first round at {self.t0}")
        return (t - self.t0) // self.collect_duration_sec

    def round_start_time(self, n: int) -> int:
        """Calculates the start time of a given round ID."""
        return self.t0 + n * self.collect_duration_sec

    def choose_start_timestamp(self, n: int) -> int:
        """Calculates the start of the choose phase for a given round ID."""
        return self.round_start_time(n + 1)

    def choose_end_timestamp(self, n: int) -> int:
        """Calculates the end of the choose phase for a given round ID."""
        return self.choose_start_timestamp(n) + self.choose_duration_sec

    def next_choose_end(self, t: int) -> tuple[int, int]:
        """
        Returns the round ID and end timestamp of the next choose phase to end.
        """
        if t < self.t0 + self.choose_duration_sec + 1:
            return 0, self.choose_end_timestamp(0)

        round_id = (t - self.t0 - self.choose_duration_sec - 1) // self.collect_duration_sec
        end_timestamp = self.choose_end_timestamp(round_id)
        return round_id, end_timestamp

    def last_collect_phase_start(self, t: int) -> tuple[int, int]:
        """Returns the round ID and start timestamp of the latest round."""
        round_id = self.round_id_for_timestamp(t)
        start_timestamp = self.round_start_time(round_id)
        return round_id, start_timestamp

    def expected_reward_epoch_start_timestamp(self, reward_epoch_id: int) -> int:
        """Calculates the expected start timestamp of a given reward epoch."""
        return (self.t0 + self.t0_reward_delay +
                self.reward_epoch_length * self.collect_duration_sec * reward_epoch_id)