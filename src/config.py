import tomli
from typing import Any, Dict

class Config:
    def __init__(self, user_config_path: str, system_config_dir: str):
        self.user_config = self._load_toml(user_config_path)

        protocol_id = self.user_config.get("protocol_id")
        chain = self.user_config.get("chain")

        if not all([protocol_id, chain]):
            raise ValueError("Missing 'protocol_id' or 'chain' in user config")

        system_config_path = f"{system_config_dir}/{protocol_id}/{chain}.toml"
        self.system_config = self._load_toml(system_config_path)

    def _load_toml(self, path: str) -> Dict[str, Any]:
        """Loads a TOML file and returns its contents as a dictionary."""
        try:
            with open(path, "rb") as f:
                return tomli.load(f)
        except FileNotFoundError:
            raise FileNotFoundError(f"Configuration file not found at: {path}")
        except tomli.TOMLDecodeError as e:
            raise ValueError(f"Error decoding TOML file at {path}: {e}")

    @property
    def logging_config(self) -> Dict[str, Any]:
        return self.user_config.get("logging", {})

    @property
    def database_config(self) -> Dict[str, Any]:
        return self.user_config.get("db", {})

    @property
    def rest_server_config(self) -> Dict[str, Any]:
        return self.user_config.get("rest_server", {})

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
        print("Successfully loaded configurations.")
        print("\n--- User Config ---")
        print(config.user_config)
        print("\n--- System Config ---")
        print(config.system_config)

    except (FileNotFoundError, ValueError) as e:
        print(f"Error: {e}")