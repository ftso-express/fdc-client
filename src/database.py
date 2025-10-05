import mysql.connector
import psycopg2
from typing import Any, Dict
from src.logger import log

class DatabaseConnector:
    def __init__(self, db_config: Dict[str, Any]):
        self.db_config = db_config
        self.connection = None
        log.info("DatabaseConnector initialized.")

    def connect(self):
        """Establishes a connection to the database based on the type."""
        db_type = self.db_config.get("type")
        log.info(f"Attempting to connect to database of type: {db_type}")
        if db_type == "mysql":
            self._connect_mysql()
        elif db_type == "postgres":
            self._connect_postgres()
        else:
            log.error(f"Unsupported database type: {db_type}")
            raise ValueError(f"Unsupported database type: {db_type}")
        return self.connection

    def _connect_mysql(self):
        """Connects to a MySQL database."""
        log.debug("Connecting to MySQL...")
        try:
            self.connection = mysql.connector.connect(
                host=self.db_config["host"],
                port=self.db_config["port"],
                user=self.db_config["username"],
                password=self.db_config["password"],
                database=self.db_config["database"],
            )
            log.info("Successfully connected to MySQL database.")
        except mysql.connector.Error as err:
            log.error(f"Error connecting to MySQL database: {err}")
            raise

    def _connect_postgres(self):
        """Connects to a PostgreSQL database."""
        log.debug("Connecting to PostgreSQL...")
        try:
            self.connection = psycopg2.connect(
                host=self.db_config["host"],
                port=self.db_config["port"],
                user=self.db_config["username"],
                password=self.db_config["password"],
                dbname=self.db_config["database"],
            )
            log.info("Successfully connected to PostgreSQL database.")
        except psycopg2.Error as err:
            log.error(f"Error connecting to PostgreSQL database: {err}")
            raise

    def close(self):
        """Closes the database connection."""
        if self.connection:
            self.connection.close()
            log.info("Database connection closed.")