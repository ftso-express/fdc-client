import uvicorn
from fastapi import FastAPI, Security, HTTPException, Depends
from fastapi.security import APIKeyHeader
from prometheus_client import make_asgi_app
from starlette.status import HTTP_403_FORBIDDEN
from typing import List

from src.config import Config
from src.shared import DataPipes

# API Key security scheme
api_key_header = APIKeyHeader(name="X-API-KEY", auto_error=False)

class RestServer:
    def __init__(self, data_pipes: DataPipes, config: Config):
        self.data_pipes = data_pipes
        # Keep a reference to the full config for metrics
        self.full_config = config
        self.config = config.rest_server_config
        self.app = FastAPI(
            title=self.config.get("title", "FDC Python Client"),
            version=self.config.get("version", "0.1.0"),
            docs_url=self.config.get("swagger_path", "/docs")
        )
        self._configure_routes()

    def _get_valid_api_keys(self) -> List[str]:
        return self.config.get("api_keys", [])

    def _api_key_auth(self, api_key: str = Security(api_key_header)):
        """Dependency to validate the API key."""
        if api_key not in self._get_valid_api_keys():
            raise HTTPException(
                status_code=HTTP_403_FORBIDDEN, detail="Could not validate credentials"
            )
        return api_key

    def _configure_routes(self):
        # Add metrics endpoint if enabled
        if self.full_config.metrics_config.get("enabled", False):
            metrics_app = make_asgi_app()
            self.app.mount("/metrics", metrics_app)

        # Define sub-routers as specified in the config
        fsp_router = self.config.get("fsp_sub_router_path", "/fsp")
        da_router = self.config.get("da_sub_router_path", "/da")

        @self.app.get("/", tags=["Health"])
        async def health_check():
            return {"status": "ok"}

        @self.app.get(f"{fsp_router}/data", dependencies=[Depends(self._api_key_auth)], tags=["FSP"])
        async def get_fsp_data():
            # This is where the logic to get data for the FSP client would go.
            # It would interact with self.data_pipes.rounds.
            return {"message": "FSP data placeholder"}

        @self.app.get(f"{da_router}/proof", dependencies=[Depends(self._api_key_auth)], tags=["DA"])
        async def get_da_proof():
            # This is where the logic for DA endpoints would go.
            return {"message": "DA proof placeholder"}
