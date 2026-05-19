import logging
import json
from confluent_kafka import Consumer
from app.models import Transaction

logger = logging.getLogger(__name__)

class TransactionConsumer:
    def __init__(self, bootstrap_servers: str, topic: str, group_id: str = "fraud-detection-group"):
        self.topic = topic
        self.consumer = Consumer({
            "bootstrap.servers": bootstrap_servers,
            "group.id": group_id,
            "auto.offset.reset": "latest"
        })
        logger.info("Kafka consumer initialised")
    
    def start(self):
        self.consumer.subscribe([self.topic])
        logger.info("Subscribed to topic: %s", self.topic)
    
    def stop(self):
        self.consumer.close()
        logger.info("Kafka consumer closed")

    async def consume(self, fraud_detector, database):
        while True:
            msg = self.consumer.poll(1.0)

            if msg is None:
                continue

            if msg.error():
                logger.error("Kafka error: %s", msg.error())
                continue

            try:
                data = json.loads(msg.value().decode("utf-8"))
                transaction = Transaction(**data)
                fraud_score = await fraud_detector.score_transaction(transaction)
                await database.insert_fraud_score(fraud_score)
                logger.info("Processed transaction %s - %s", transaction.id, fraud_score.status)
            except Exception as e:
                logger.error("Error processing message: %s", e)