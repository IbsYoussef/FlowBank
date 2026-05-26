import logging
import json
import asyncio
import os
import tempfile
import base64
from confluent_kafka import Consumer
from app.models import Transaction

logger = logging.getLogger(__name__)

class TransactionConsumer:
    def __init__(self, bootstrap_servers: str, topic: str, group_id: str = "fraud-detection-group"):
        self.topic = topic

        config = {
            "bootstrap.servers": bootstrap_servers,
            "group.id": group_id,
            "auto.offset.reset": "latest"
        }

        # If SASL credentials are set, use SASL_SSL for Aiven
        sasl_username = os.environ.get("KAFKA_SASL_USERNAME")
        sasl_password = os.environ.get("KAFKA_SASL_PASSWORD")
        ca_cert = os.environ.get("KAFKA_CA_CERT")

        if sasl_username and sasl_password and ca_cert:
            # # Decode base64 if needed
            try:
                decoded_cert = base64.b64decode(ca_cert).decode('utf-8')
            except Exception:
                decoded_cert = ca_cert # raw PEM
            
            ca_file = tempfile.NamedTemporaryFile(mode='w', suffix='.pem', delete=False)
            ca_file.write(decoded_cert)
            ca_file.flush()
            ca_file.close()

            config.update({
                "security.protocol": "SASL_SSL",
                "sasl.mechanism": "PLAIN",
                "sasl.username": sasl_username,
                "sasl.password": sasl_password,
                "ssl.ca.location": ca_file.name,
            })
            logger.info("Kafka consumer configured with SASL_SSL")
        
        self.consumer = Consumer(config)
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
                await asyncio.sleep(3) # Delay to allow writing transaction to database
                fraud_score = await fraud_detector.score_transaction(transaction)
                await database.insert_fraud_score(fraud_score)
                logger.info("Processed transaction %s - %s", transaction.transaction_id, fraud_score.status)
            except Exception as e:
                logger.error("Error processing message: %s", e)