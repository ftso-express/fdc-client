import pytest
from src.timing import Timing

@pytest.fixture
def timing_instance():
    """Provides a default Timing instance for tests."""
    return Timing({})

def test_round_id_for_timestamp(timing_instance: Timing):
    # Test timestamp before T0
    with pytest.raises(ValueError):
        timing_instance.round_id_for_timestamp(0)

    # Test cases mirroring the Go tests
    test_cases = [
        (timing_instance.t0, 0),
        (timing_instance.t0 + 10000 * timing_instance.collect_duration_sec + 2, 10000),
    ]

    for timestamp, expected_round_id in test_cases:
        round_id = timing_instance.round_id_for_timestamp(timestamp)
        assert round_id == expected_round_id

def test_times_for_rounds(timing_instance: Timing):
    test_cases = [
        {
            "round_id": 0,
            "start_time": timing_instance.t0,
            "choose_start": timing_instance.t0 + timing_instance.collect_duration_sec,
            "choose_end": (timing_instance.t0 + timing_instance.collect_duration_sec +
                           timing_instance.choose_duration_sec),
        },
        {
            "round_id": 10000,
            "start_time": timing_instance.t0 + 10000 * timing_instance.collect_duration_sec,
            "choose_start": (timing_instance.t0 + 10001 * timing_instance.collect_duration_sec),
            "choose_end": (timing_instance.t0 + 10001 * timing_instance.collect_duration_sec +
                           timing_instance.choose_duration_sec),
        },
    ]

    for case in test_cases:
        assert timing_instance.round_start_time(case["round_id"]) == case["start_time"]
        assert timing_instance.choose_start_timestamp(case["round_id"]) == case["choose_start"]
        assert timing_instance.choose_end_timestamp(case["round_id"]) == case["choose_end"]

def test_times_for_timestamps(timing_instance: Timing):
    # Test case for timestamp = 0
    round_id_choose, choose_end = timing_instance.next_choose_end(0)
    assert round_id_choose == 0
    expected_choose_end = (timing_instance.t0 + timing_instance.collect_duration_sec +
                           timing_instance.choose_duration_sec)
    assert choose_end == expected_choose_end

    with pytest.raises(ValueError):
        timing_instance.last_collect_phase_start(0)

    # Test cases mirroring the Go tests
    test_cases = [
        {
            "timestamp": timing_instance.t0,
            "round_id_choose": 0,
            "choose_end": (timing_instance.t0 + timing_instance.collect_duration_sec +
                           timing_instance.choose_duration_sec),
            "round_id_collect": 0,
            "collect_start": timing_instance.t0,
        },
        {
            "timestamp": (timing_instance.t0 + timing_instance.collect_duration_sec +
                          timing_instance.choose_duration_sec // 2),
            "round_id_choose": 0,
            "choose_end": (timing_instance.t0 + timing_instance.collect_duration_sec +
                           timing_instance.choose_duration_sec),
            "round_id_collect": 1,
            "collect_start": timing_instance.t0 + timing_instance.collect_duration_sec,
        },
    ]

    for case in test_cases:
        round_id_choose, choose_end = timing_instance.next_choose_end(case["timestamp"])
        assert round_id_choose == case["round_id_choose"]
        assert choose_end == case["choose_end"]

        round_id_collect, collect_start = timing_instance.last_collect_phase_start(case["timestamp"])
        assert round_id_collect == case["round_id_collect"]
        assert collect_start == case["collect_start"]