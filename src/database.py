import mysql.connector
import psycopg2
from typing import Any, Dict

class DatabaseConnector:
    def __init__(self, db_config: Dict[str, Any]):
        self.db_config = db_config
        self.connection = None

    def connect(self):
        """Establishes a connection to the database based on the type."""
        db_type = self.db_config.get("type")
        if db_type == "mysql":
            self._connect_mysql()
        elif db_type == "postgres":
            self._connect_postgres()
        else:
            raise ValueError(f"Unsupported database type: {db_type}")
        return self.connection

    def _connect_mysql(self):
        """Connects to a MySQL database."""
        try:
            self.connection = mysql.connector.connect(
                host=self.db_config["host"],
                port=self.db_config["port"],
                user=self.db_config["username"],
                password=self.db_config["password"],
                database=self.db_config["database"],
            )
            print("Successfully connected to MySQL database.")
        except mysql.connector.Error as err:
            print(f"Error connecting to MySQL database: {err}")
            raise

    def _connect_postgres(self):
        """Connects to a PostgreSQL database."""
        try:
            self.connection = psycopg2.connect(
                host=self.db_config["host"],
                port=self.db_config["port"],
                user=self.db_config["username"],
                password=self.db_config["password"],
                dbname=self.db_config["database"],
            )
            print("Successfully connected to PostgreSQL database.")
        except psycopg2.Error as err:
            print(f"Error connecting to PostgreSQL database: {err}")
            raise

    def close(self):
        """Closes the database connection."""
        if self.connection:
            self.connection.close()
            print("Database connection closed.")