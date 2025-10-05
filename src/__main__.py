import asyncio
import signal
import uvicorn
from typing import Set

from src.config import Config
from src.collector import Collector
from src.manager import Manager
from src.server import RestServer
from src.shared import DataPipes

async def main():
    """Main application entrypoint."""

    tasks: Set[asyncio.Task] = set()

    # 1. Initialize all components
    try:
        config = Config("configs/userConfig.toml", "configs/systemConfigs")
        data_pipes = DataPipes()

        # Pass data_pipes to the collector
        collector = Collector(config, data_pipes)
        manager = Manager(config, data_pipes)
        server = RestServer(data_pipes, config)

        # The server needs to be run in a non-blocking way.
        # We will run uvicorn programmatically.
        uv_config = uvicorn.Config(
            app=server.app,
            host="0.0.0.0",
            port=8080,
            log_level="info"
        )
        uv_server = uvicorn.Server(uv_config)

    except (FileNotFoundError, ValueError) as e:
        print(f"Initialization Error: {e}")
        return

    # 2. Setup graceful shutdown
    loop = asyncio.get_running_loop()
    stop_event = asyncio.Event()

    def shutdown_handler():
        print("Shutdown signal received.")
        stop_event.set()

    for sig in (signal.SIGINT, signal.SIGTERM):
        loop.add_signal_handler(sig, shutdown_handler)

    # 3. Create and manage tasks
    try:
        collector_task = loop.create_task(collector.run({}))
        tasks.add(collector_task)

        manager_task = loop.create_task(manager.run({}))
        tasks.add(manager_task)

        # Run the server in a separate task
        server_task = loop.create_task(uv_server.serve())
        tasks.add(server_task)

        print("Application started. Press Ctrl+C to shut down.")

        # Wait for the stop event
        await stop_event.wait()

    finally:
        print("Shutting down application...")

        # Gracefully shut down all tasks
        for task in tasks:
            task.cancel()

        await asyncio.gather(*tasks, return_exceptions=True)

        print("Application shut down gracefully.")


if __name__ == "__main__":
    asyncio.run(main())