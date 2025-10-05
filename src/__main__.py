import asyncio
import signal
import uvicorn
from typing import Set

from src.config import Config
from src.collector import Collector
from src.manager import Manager
from src.server import RestServer
from src.shared import DataPipes
from src.logger import log

async def main():
    """Main application entrypoint."""

    tasks: Set[asyncio.Task] = set()
    log.info("FDC Client starting up...")

    # 1. Initialize all components
    try:
        log.debug("Loading configurations...")
        config = Config("configs/userConfig.toml", "configs/systemConfigs")
        data_pipes = DataPipes()
        log.debug("Configurations loaded successfully.")

        log.debug("Initializing components...")
        collector = Collector(config, data_pipes)
        manager = Manager(config, data_pipes)
        server = RestServer(data_pipes, config)
        log.debug("Components initialized.")

        # The server needs to be run in a non-blocking way.
        # We will run uvicorn programmatically.
        log.debug("Setting up UVicorn server.")
        uv_config = uvicorn.Config(
            app=server.app,
            host="0.0.0.0",
            port=8080,
            log_config=None  # We use our own logger
        )
        uv_server = uvicorn.Server(uv_config)
        log.debug("Uvicorn server configured.")

    except (FileNotFoundError, ValueError) as e:
        log.error(f"Initialization Error: {e}")
        return

    # 2. Setup graceful shutdown
    loop = asyncio.get_running_loop()
    stop_event = asyncio.Event()

    def shutdown_handler():
        log.info("Shutdown signal received.")
        stop_event.set()

    for sig in (signal.SIGINT, signal.SIGTERM):
        loop.add_signal_handler(sig, shutdown_handler)

    # 3. Create and manage tasks
    try:
        log.info("Creating and starting tasks...")
        collector_task = loop.create_task(collector.run({}))
        tasks.add(collector_task)
        log.debug("Collector task created.")

        manager_task = loop.create_task(manager.run({}))
        tasks.add(manager_task)
        log.debug("Manager task created.")

        # Run the server in a separate task
        server_task = loop.create_task(uv_server.serve())
        tasks.add(server_task)
        log.debug("Server task created.")

        log.info("Application started. Press Ctrl+C to shut down.")

        # Wait for the stop event
        await stop_event.wait()

    finally:
        log.info("Shutting down application...")

        # Gracefully shut down all tasks
        for task in tasks:
            task.cancel()

        await asyncio.gather(*tasks, return_exceptions=True)

        log.info("Application shut down gracefully.")


if __name__ == "__main__":
    asyncio.run(main())