import os
import tomli
from typing import Any, Dict
from src.logger import log

class Config:
    def __init__(self, user_config_path: str, system_config_dir: str):
        log.info(f"Loading user configuration from: {user_config_path}")
        self.user_config = self._load_toml(user_config_path)

        protocol_id = self.user_config.get("protocol_id")
        chain = self.user_config.get("chain")
        log.debug(f"Protocol ID: {protocol_id}, Chain: {chain}")

        if not all([protocol_id, chain]):
            raise ValueError("Missing 'protocol_id' or 'chain' in user config")

        system_config_path = f"{system_config_dir}/{protocol_id}/{chain}.toml"
        log.info(f"Loading system configuration from: {system_config_path}")
        self.system_config = self._load_toml(system_config_path)
        log.info("Configurations loaded successfully.")

    def _load_toml(self, path: str) -> Dict[str, Any]:
        """Loads a TOML file and returns its contents as a dictionary."""
        log.debug(f"Attempting to load TOML file at: {path}")
        try:
            with open(path, "rb") as f:
                config = tomli.load(f)
                log.debug(f"Successfully loaded TOML file: {path}")
                return config
        except FileNotFoundError:
            log.error(f"Configuration file not found at: {path}")
            raise FileNotFoundError(f"Configuration file not found at: {path}")
        except tomli.TOMLDecodeError as e:
            log.error(f"Error decoding TOML file at {path}: {e}")
            raise ValueError(f"Error decoding TOML file at {path}: {e}")

    @property
    def logging_config(self) -> Dict[str, Any]:
        return self.user_config.get("logging", {})

    @property
    def database_config(self) -> Dict[str, Any]:
        """
        Returns the database configuration.
        It prioritizes environment variables over the TOML file config.
        """
        config = self.user_config.get("db", {})
        log.debug("Loading database configuration...")

        db_type = os.environ.get("DB_TYPE", config.get("type", "mysql"))
        config["type"] = db_type
        log.debug(f"Database type: {db_type}")

        default_port = 5432 if db_type == "postgres" else 3306

        config["host"] = os.environ.get("DB_HOST", config.get("host"))
        config["port"] = int(os.environ.get("DB_PORT", config.get("port", default_port)))
        config["database"] = os.environ.get("DB_DATABASE", config.get("database"))
        config["username"] = os.environ.get("DB_USERNAME", config.get("username"))
        config["password"] = os.environ.get("DB_PASSWORD", config.get("password"))
        log.debug("Database configuration loaded.")

        return config

    @property
    def rest_server_config(self) -> Dict[str, Any]:
        return self.user_config.get("rest_server", {})

    @property
    def metrics_config(self) -> Dict[str, Any]:
        return self.user_config.get("metrics", {})

    @property
    def attestation_provider_config(self) -> Dict[str, Any]:
        return self.user_config.get("verifiers", {})

    @property
    def queue_config(self) -> Dict[str, Any]:
        return self.user_config.get("queue", {})

    @property
    def contract_addresses(self) -> Dict[str, str]:
        return self.system_config.get("addresses", {})

    @property
    def timing_config(self) -> Dict[str, Any]:
        return self.system_config.get("timing", {})

if __name__ == "__main__":
    # Example usage:
    try:
        # These paths are relative to the project root
        config = Config("configs/userConfig.toml", "configs/systemConfigs")
        log.info("Successfully loaded configurations.")
        log.info("\n--- User Config ---")
        log.info(config.user_config)
        log.info("\n--- System Config ---")
        log.info(config.system_config)

    except (FileNotFoundError, ValueError) as e:
        log.error(f"Error: {e}")