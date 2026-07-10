from fastapi import FastAPI, Request
from typing import Optional
from contextlib import asynccontextmanager
from app.database import FraudDatabase
from app.detector import FraudDetector
from app.consumer import TransactionConsumer
import asyncio
import logging
import os

logger = logging.getLogger(__name__)

DATABASE_URL = os.getenv("DATABASE_URL", "postgresql://flowbank_user:flowbank_pass@localhost:5432/flowbank")
KAFKA_BOOTSTRAP_SERVERS = os.getenv("KAFKA_BOOTSTRAP_SERVERS", "localhost:9092")
KAFKA_TOPIC = os.getenv("KAFKA_TOPIC", "transactions")

@asynccontextmanager
async def lifespan(app: FastAPI):
    # Startup
    db = FraudDatabase(DATABASE_URL)
    await db.connect()

    detector = FraudDetector(db)

    consumer = TransactionConsumer(KAFKA_BOOTSTRAP_SERVERS, KAFKA_TOPIC)
    consumer.start()

    asyncio.create_task(consumer.consume(detector, db))

    app.state.db = db
    app.state.consumer = consumer

    logger.info("Fraud detection service started")
    yield

    # Shutdown
    consumer.stop()
    await db.close()
    logger.info("Fraud detection service stopped")

app = FastAPI(lifespan=lifespan)

@app.get("/")
async def root():
    return {"message": "root pathway working"}

@app.get("/health")
async def health():
    return {"health": "healthy"}

@app.get("/api/v1/fraud-scores")
async def get_fraud_scores(request: Request, limit: int = 100, status: Optional[str] = None):
    scores = await request.app.state.db.get_fraud_scores(limit=limit, status=status)
    return {"scores": scores, "count": len(scores)}

if __name__ == "__main__":
    import uvicorn
    port = int(os.environ.get("PORT", 8000))
    uvicorn.run(app, host="0.0.0.0", port=port)